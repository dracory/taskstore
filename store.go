package taskstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dracory/neat"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

type StoreInterface interface {
	// GetTaskDefinitionTableName returns the task definition table name
	GetTaskDefinitionTableName() string
	// SetTaskDefinitionTableName sets the task definition table name
	SetTaskDefinitionTableName(tableName string)
	// GetTaskQueueTableName returns the task queue table name
	GetTaskQueueTableName() string
	// SetTaskQueueTableName sets the task queue table name
	SetTaskQueueTableName(tableName string)
	// GetScheduleTableName returns the schedule table name
	GetScheduleTableName() string
	// SetScheduleTableName sets the schedule table name
	SetScheduleTableName(tableName string)

	// MigrateDown drops all tables
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error

	// MigrateUp creates all tables
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	// EnableDebug enables debug mode
	EnableDebug(debug bool) StoreInterface

	// SetErrorHandler sets the error handler
	SetErrorHandler(handler func(queueName, taskID string, err error)) StoreInterface

	// == TaskQueue Methods ==

	TaskQueueCount(ctx context.Context, options TaskQueueQueryInterface) (int64, error)
	TaskQueueCreate(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueDelete(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueDeleteByID(ctx context.Context, id string) error
	TaskQueueFindByID(ctx context.Context, TaskQueueID string) (TaskQueueInterface, error)
	TaskQueueList(ctx context.Context, query TaskQueueQueryInterface) ([]TaskQueueInterface, error)
	TaskQueueSoftDelete(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueSoftDeleteByID(ctx context.Context, id string) error
	TaskQueueUpdate(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueClaimNext(ctx context.Context, queueName string) (TaskQueueInterface, error)

	// Deprecated: Use NewTaskQueueRunner instead. These methods will be removed in a future version.
	// See docs/runners.md for the recommended approach.
	TaskQueueRunDefault(ctx context.Context, processSeconds int, unstuckMinutes int)
	// Deprecated: Use NewTaskQueueRunner instead. These methods will be removed in a future version.
	// See docs/runners.md for the recommended approach.
	TaskQueueRunSerial(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	// Deprecated: Use NewTaskQueueRunner instead. These methods will be removed in a future version.
	// See docs/runners.md for the recommended approach.
	TaskQueueRunConcurrent(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	// Deprecated: Use TaskQueueRunner.Stop() instead. These methods will be removed in a future version.
	// See docs/runners.md for the recommended approach.
	TaskQueueStop()
	// Deprecated: Use TaskQueueRunner.Stop() instead. These methods will be removed in a future version.
	// See docs/runners.md for the recommended approach.
	TaskQueueStopByName(queueName string)
	TaskQueueProcessTask(ctx context.Context, queuedTask TaskQueueInterface) (bool, error)

	// == TaskDefinition Methods ==

	TaskDefinitionCount(ctx context.Context, options TaskDefinitionQueryInterface) (int64, error)
	TaskDefinitionCreate(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDelete(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDeleteByID(ctx context.Context, id string) error
	TaskDefinitionFindByAlias(ctx context.Context, alias string) (TaskDefinitionInterface, error)
	TaskDefinitionFindByID(ctx context.Context, id string) (TaskDefinitionInterface, error)
	TaskDefinitionList(ctx context.Context, options TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error)
	TaskDefinitionSoftDelete(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionSoftDeleteByID(ctx context.Context, id string) error
	TaskDefinitionUpdate(ctx context.Context, TaskDefinition TaskDefinitionInterface) error

	// TaskDefinition Operations
	TaskDefinitionEnqueueByAlias(ctx context.Context, queueName string, alias string, parameters map[string]any) (TaskQueueInterface, error)
	TaskDefinitionExecuteCli(alias string, args []string) bool

	// == TaskHandler Methods ==

	TaskHandlerList() []TaskDefinitionHandlerInterface
	TaskHandlerAdd(ctx context.Context, taskHandler TaskDefinitionHandlerInterface, createIfMissing bool) error

	// == Schedule Methods ==

	ScheduleCount(ctx context.Context, options ScheduleQueryInterface) (int64, error)
	ScheduleCreate(ctx context.Context, schedule ScheduleInterface) error
	ScheduleDelete(ctx context.Context, schedule ScheduleInterface) error
	ScheduleDeleteByID(ctx context.Context, id string) error
	ScheduleFindByID(ctx context.Context, id string) (ScheduleInterface, error)
	ScheduleList(ctx context.Context, options ScheduleQueryInterface) ([]ScheduleInterface, error)
	ScheduleSoftDelete(ctx context.Context, schedule ScheduleInterface) error
	ScheduleSoftDeleteByID(ctx context.Context, id string) error
	ScheduleUpdate(ctx context.Context, schedule ScheduleInterface) error
	ScheduleRun(ctx context.Context) error
}

// Store defines a session store
type Store struct {
	taskDefinitionTableName string
	taskQueueTableName      string
	scheduleTableName       string
	taskHandlers            []TaskDefinitionHandlerInterface
	db                      *neat.Database
	automigrateEnabled      bool
	debugEnabled            bool
	queueMu                 sync.Mutex
	queueRunners            map[string]*queueRunner
	maxConcurrency          int // Max concurrent tasks in async mode (default: 10)
	errorHandler            func(queueName, taskID string, err error)
	logger                  *slog.Logger
	isSQLite                bool
}

type queueRunner struct {
	cancel         context.CancelFunc
	wg             sync.WaitGroup // Tracks the main queue loop goroutine
	taskWg         sync.WaitGroup // Tracks child task goroutines
	maxConcurrency int            // Maximum number of concurrent tasks (0 = unlimited)
	semaphore      chan struct{}  // Semaphore for concurrency control
}

var _ StoreInterface = (*Store)(nil)

// NewStoreOptions define the options for creating a new task store
type NewStoreOptions struct {
	TaskDefinitionTableName string
	TaskQueueTableName      string
	ScheduleTableName       string
	DB                      *sql.DB
	AutomigrateEnabled      bool
	DebugEnabled            bool
	MaxConcurrency          int                                       // Max concurrent tasks (default: 10, 0 = unlimited)
	ErrorHandler            func(queueName, taskID string, err error) // Optional error callback
}

// NewStore creates a new task store
func NewStore(opts NewStoreOptions) (*Store, error) {
	if opts.DB == nil {
		return nil, errors.New("task store: DB is required")
	}

	if opts.TaskDefinitionTableName == "" {
		return nil, errors.New("task store: TaskDefinitionTableName is required")
	}

	if opts.TaskQueueTableName == "" {
		return nil, errors.New("task store: TaskQueueTableName is required")
	}

	if opts.ScheduleTableName == "" {
		return nil, errors.New("task store: ScheduleTableName is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &Store{
		taskDefinitionTableName: opts.TaskDefinitionTableName,
		taskQueueTableName:      opts.TaskQueueTableName,
		scheduleTableName:       opts.ScheduleTableName,
		automigrateEnabled:      opts.AutomigrateEnabled,
		db:                      neatDB,
		debugEnabled:            opts.DebugEnabled,
		queueRunners:            map[string]*queueRunner{},
		maxConcurrency:          opts.MaxConcurrency,
		errorHandler:            opts.ErrorHandler,
		logger:                  logger,
		isSQLite:                strings.Contains(fmt.Sprintf("%T", opts.DB.Driver()), "sqlite"),
	}

	// Set default max concurrency if not specified
	if store.maxConcurrency == 0 {
		store.maxConcurrency = 10
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// MigrateUp creates all tables
func (st *Store) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.taskDefinitionTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: task_definition table already exists", "table", st.taskDefinitionTableName)
		}
	} else {
		err := st.db.Schema().Create(st.taskDefinitionTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 50)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 50)
			table.String(COLUMN_ALIAS, 100)
			table.Unique(COLUMN_ALIAS)
			table.String(COLUMN_TITLE, 255)
			table.Text(COLUMN_MEMO)
			table.String(COLUMN_DESCRIPTION, 255)
			table.Integer(COLUMN_IS_RECURRING)
			table.String(COLUMN_RECURRENCE_RULE, 500)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateUp failed for task_definition", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.taskQueueTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: task_queue table already exists", "table", st.taskQueueTableName)
		}
	} else {
		err := st.db.Schema().Create(st.taskQueueTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 50)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_QUEUE_NAME, 100)
			table.String(COLUMN_TASK_ID, 50)
			table.Text(COLUMN_PARAMETERS)
			table.String(COLUMN_STATUS, 50)
			table.Text(COLUMN_OUTPUT)
			table.Text(COLUMN_DETAILS)
			table.Integer(COLUMN_ATTEMPTS)
			table.DateTime(COLUMN_STARTED_AT)
			table.DateTime(COLUMN_COMPLETED_AT)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateUp failed for task_queue", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.scheduleTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: schedule table already exists", "table", st.scheduleTableName)
		}
	} else {
		err := st.db.Schema().Create(st.scheduleTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 50)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_NAME, 100)
			table.String(COLUMN_DESCRIPTION, 255)
			table.String(COLUMN_STATUS, 50)
			table.Text(COLUMN_RECURRENCE_RULE)
			table.String(COLUMN_QUEUE_NAME, 100)
			table.String(COLUMN_TASK_DEFINITION_ID, 50)
			table.Text(COLUMN_PARAMETERS)
			table.DateTime(COLUMN_START_AT)
			table.DateTime(COLUMN_END_AT)
			table.Integer(COLUMN_EXECUTION_COUNT)
			table.Integer(COLUMN_MAX_EXECUTION_COUNT)
			table.DateTime(COLUMN_LAST_RUN_AT)
			table.DateTime(COLUMN_NEXT_RUN_AT)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateUp failed for schedule", "error", err)
			}
			return err
		}
	}

	return nil
}

// MigrateDown drops all tables
func (st *Store) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.scheduleTableName) {
		if err := st.db.Schema().Drop(st.scheduleTableName); err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown failed for schedule", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.taskQueueTableName) {
		if err := st.db.Schema().Drop(st.taskQueueTableName); err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown failed for task_queue", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.taskDefinitionTableName) {
		if err := st.db.Schema().Drop(st.taskDefinitionTableName); err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown failed for task_definition", "error", err)
			}
			return err
		}
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debugEnabled bool) StoreInterface {
	st.debugEnabled = debugEnabled
	if debugEnabled {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return st
}

// GetDB returns the underlying *sql.DB.
func (st *Store) GetDB() *sql.DB {
	db, _ := st.db.DB()
	return db
}

// GetTaskDefinitionTableName returns the task definition table name
func (st *Store) GetTaskDefinitionTableName() string {
	return st.taskDefinitionTableName
}

// SetTaskDefinitionTableName sets the task definition table name
func (st *Store) SetTaskDefinitionTableName(tableName string) {
	st.taskDefinitionTableName = tableName
}

// GetTaskQueueTableName returns the task queue table name
func (st *Store) GetTaskQueueTableName() string {
	return st.taskQueueTableName
}

// SetTaskQueueTableName sets the task queue table name
func (st *Store) SetTaskQueueTableName(tableName string) {
	st.taskQueueTableName = tableName
}

// GetScheduleTableName returns the schedule table name
func (st *Store) GetScheduleTableName() string {
	return st.scheduleTableName
}

// SetScheduleTableName sets the schedule table name
func (st *Store) SetScheduleTableName(tableName string) {
	st.scheduleTableName = tableName
}

// SetErrorHandler - sets a custom error handler for queue processing errors
func (st *Store) SetErrorHandler(handler func(queueName, taskID string, err error)) StoreInterface {
	st.errorHandler = handler
	return st
}

// TaskQueueRunDefault starts the queue processor for the default queue.
// Equivalent to calling TaskQueueRunSerial with DefaultQueueName.
//
// Deprecated: Use NewTaskQueueRunner instead. This method will be removed in a future version.
// See docs/runners.md for the recommended approach.
func (store *Store) TaskQueueRunDefault(
	ctx context.Context,
	processSeconds int,
	unstuckMinutes int,
) {
	store.TaskQueueRunSerial(ctx, DefaultQueueName, processSeconds, unstuckMinutes)
}

// TaskQueueRunSerial starts a queue processor that handles tasks one at a time (serially).
// Each task must complete before the next one starts.
// The processor runs in a background goroutine and can be stopped via TaskQueueStopByName.
//
// Deprecated: Use NewTaskQueueRunner instead. This method will be removed in a future version.
// See docs/runners.md for the recommended approach.
func (store *Store) TaskQueueRunSerial(
	ctx context.Context,
	queueName string,
	processSeconds int,
	unstuckMinutes int,
) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return
	}

	queueName = normalizeQueueName(queueName)

	store.queueMu.Lock()
	if _, exists := store.queueRunners[queueName]; exists {
		store.queueMu.Unlock()
		return
	}

	runCtx, cancel := context.WithCancel(ctx)
	runner := &queueRunner{cancel: cancel}
	runner.wg.Add(1)
	store.queueRunners[queueName] = runner
	store.queueMu.Unlock()

	go func() {
		defer func() {
			store.queueMu.Lock()
			delete(store.queueRunners, queueName)
			store.queueMu.Unlock()
			runner.wg.Done()
		}()

		store.queueRunLoopSync(runCtx, queueName, processSeconds, unstuckMinutes)
	}()
}

// TaskQueueRunConcurrent starts a queue processor that handles multiple tasks concurrently.
// Tasks are processed in parallel up to the configured MaxConcurrency limit.
// The processor runs in a background goroutine and can be stopped via TaskQueueStopByName.
//
// Deprecated: Use NewTaskQueueRunner instead. This method will be removed in a future version.
// See docs/runners.md for the recommended approach.
func (store *Store) TaskQueueRunConcurrent(
	ctx context.Context,
	queueName string,
	processSeconds int,
	unstuckMinutes int,
) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return
	}

	queueName = normalizeQueueName(queueName)

	store.queueMu.Lock()
	if _, exists := store.queueRunners[queueName]; exists {
		store.queueMu.Unlock()
		return
	}

	runCtx, cancel := context.WithCancel(ctx)
	runner := &queueRunner{
		cancel:         cancel,
		maxConcurrency: store.maxConcurrency,
		semaphore:      make(chan struct{}, store.maxConcurrency),
	}
	runner.wg.Add(1)
	store.queueRunners[queueName] = runner
	store.queueMu.Unlock()

	go func() {
		defer func() {
			store.queueMu.Lock()
			delete(store.queueRunners, queueName)
			store.queueMu.Unlock()
			runner.wg.Done()
		}()

		store.queueRunLoopAsync(runCtx, queueName, processSeconds, unstuckMinutes, runner)
	}()
}

func (store *Store) queueRunLoopSync(
	ctx context.Context,
	queueName string,
	processSeconds int,
	unstuckMinutes int,
) {
	if processSeconds <= 0 {
		processSeconds = 10
	}
	if unstuckMinutes <= 0 {
		unstuckMinutes = 1
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		store.TaskQueueUnstuckByQueue(ctx, queueName, unstuckMinutes)

		if !sleepWithContext(ctx, time.Second) {
			return
		}

		if err := store.TaskQueueProcessNextByQueue(ctx, queueName); err != nil && store.debugEnabled {
			log.Println("TaskQueueProcessNext error:", err)
		}

		if !sleepWithContext(ctx, time.Duration(processSeconds)*time.Second) {
			return
		}
	}
}

func (store *Store) queueRunLoopAsync(
	ctx context.Context,
	queueName string,
	processSeconds int,
	unstuckMinutes int,
	runner *queueRunner,
) {
	if processSeconds <= 0 {
		processSeconds = 10
	}
	if unstuckMinutes <= 0 {
		unstuckMinutes = 1
	}

	// When context is done, wait for all tasks to complete
	defer func() {
		runner.taskWg.Wait() // Wait for all child goroutines to finish
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		store.TaskQueueUnstuckByQueue(ctx, queueName, unstuckMinutes)

		if !sleepWithContext(ctx, time.Second) {
			return
		}

		// Acquire semaphore slot (blocks if at max concurrency)
		select {
		case runner.semaphore <- struct{}{}:
			// Got a slot, proceed
		case <-ctx.Done():
			return
		}

		// Get the next task
		nextTask, err := store.TaskQueueClaimNext(ctx, queueName)
		if err != nil {
			<-runner.semaphore // Release slot on error
			if store.debugEnabled {
				log.Println("TaskQueueClaimNext error:", err)
			}
			if !sleepWithContext(ctx, time.Duration(processSeconds)*time.Second) {
				return
			}
			continue
		}

		if nextTask == nil {
			<-runner.semaphore // Release slot when no task available
			if !sleepWithContext(ctx, time.Duration(processSeconds)*time.Second) {
				return
			}
			continue
		}

		// Track the goroutine
		runner.taskWg.Add(1)

		// Spawn goroutine to process the task
		go func(task TaskQueueInterface) {
			defer func() {
				<-runner.semaphore   // Release semaphore slot
				runner.taskWg.Done() // Mark goroutine as complete
			}()

			_, processErr := store.QueuedTaskProcessWithContext(ctx, task)
			if processErr != nil {
				// Call error handler if configured
				if store.errorHandler != nil {
					store.errorHandler(queueName, task.GetID(), processErr)
				} else if store.debugEnabled {
					log.Println("QueuedTaskProcess error:", processErr)
				}
			}
		}(nextTask)

		if !sleepWithContext(ctx, time.Duration(processSeconds)*time.Second) {
			return
		}
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		select {
		case <-ctx.Done():
			return false
		default:
			return true
		}
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// TaskQueueStop stops the default queue processor.
// It blocks until the worker goroutine and all tasks have fully completed.
//
// Deprecated: Use TaskQueueRunner.Stop() instead. This method will be removed in a future version.
// See docs/runners.md for the recommended approach.
func (store *Store) TaskQueueStop() {
	store.TaskQueueStopByName(DefaultQueueName)
}

// TaskQueueStopByName stops the specified queue processor.
// It cancels the context, waits for the queue loop to exit,
// and waits for all in-flight tasks to complete.
//
// Deprecated: Use TaskQueueRunner.Stop() instead. This method will be removed in a future version.
// See docs/runners.md for the recommended approach.
func (store *Store) TaskQueueStopByName(queueName string) {
	queueName = normalizeQueueName(queueName)

	store.queueMu.Lock()
	runner, exists := store.queueRunners[queueName]
	if !exists {
		store.queueMu.Unlock()
		return
	}
	delete(store.queueRunners, queueName)
	store.queueMu.Unlock()

	// Cancel the context to stop the queue loop
	runner.cancel()

	// Wait for the main queue loop to exit
	runner.wg.Wait()

	// Wait for all child task goroutines to complete
	// Note: This is important for async queues
	runner.taskWg.Wait()
}

// TaskQueueUnstuck clears the queue of tasks running for more than the
// specified wait time as most probably these have abnormally
// exited (panicked) and stop the rest of the queue from being
// processed
//
// The tasks are marked as failed. However, if they are still running
// in the background and they are successfully completed, they will
// be marked as success
//
// =================================================================
// Business Logic
// 1. Checks is there are running tasks in progress
// 2. If running for more than the specified wait minutes mark as failed
// =================================================================
func (store *Store) TaskQueueUnstuck(ctx context.Context, waitMinutes int) {
	store.TaskQueueUnstuckByQueue(ctx, "", waitMinutes)
}

func (store *Store) TaskQueueUnstuckByQueue(ctx context.Context, queueName string, waitMinutes int) {
	runningTasks := store.TaskQueueFindRunningByQueue(ctx, queueName, 3)

	if len(runningTasks) < 1 {
		return
	}

	for _, runningTask := range runningTasks {
		_ = store.QueuedTaskForceFail(ctx, runningTask, waitMinutes)
	}
}

func (store *Store) TaskQueueProcessTask(ctx context.Context, queuedTask TaskQueueInterface) (bool, error) {
	return store.QueuedTaskProcessWithContext(ctx, queuedTask)
}

// QueuedTaskProcessWithContext processes a queued task with context support.
// It checks if the handler implements TaskHandlerWithContext and uses that if available,
// otherwise falls back to the standard Handle() method for backward compatibility.
func (store *Store) QueuedTaskProcessWithContext(ctx context.Context, queuedTask TaskQueueInterface) (bool, error) {
	if queuedTask == nil {
		return false, errors.New("queued task is nil")
	}

	attempts := queuedTask.GetAttempts() + 1

	queuedTask.AppendDetails("Task started")
	queuedTask.SetStatus(TaskQueueStatusRunning)
	queuedTask.SetAttempts(attempts)
	queuedTask.SetStartedAt(carbon.Now(carbon.UTC).StdTime())

	err := store.TaskQueueUpdate(ctx, queuedTask)

	if err != nil {
		return false, err
	}

	// 2. Find task definition
	task, err := store.TaskDefinitionFindByID(ctx, queuedTask.GetTaskID())

	if err != nil {
		return false, err
	}

	if task == nil {
		queuedTask.AppendDetails("Task DOES NOT exist")
		queuedTask.SetStatus(TaskQueueStatusFailed)
		queuedTask.SetCompletedAt(carbon.Now(carbon.UTC).StdTime())
		err = store.TaskQueueUpdate(ctx, queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}

			return false, err
		}

		return false, nil
	}

	// 3. Get handler and check if it supports context
	handlerFunc := store.taskHandlerFuncWithContext(task.GetAlias(), ctx)

	result := handlerFunc(queuedTask)

	if result {
		queuedTask.AppendDetails("Task completed")
		err = store.TaskQueueSuccess(ctx, queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}
		}
	} else {
		queuedTask.AppendDetails("Task failed")
		err = store.TaskQueueFail(ctx, queuedTask)

		if err != nil {
			if store.debugEnabled {
				log.Println(err)
			}
		}
	}

	return true, nil
}

// TaskDefinitionExecuteCli - CLI tool to find a task by its alias and execute its handler
// - alias "list" is reserved. it lists all the available commands
func (store *Store) TaskDefinitionExecuteCli(alias string, args []string) bool {
	argumentsMap := argsToMap(args)
	fmt.Println("INFO: Executing task:", alias, "with arguments:", argumentsMap)

	// Lists the available tasks
	if alias == "list" {
		for index, taskHandler := range store.TaskHandlerList() {
			fmt.Println("WARNING:", cast.ToString(index+1)+". Task Alias: "+taskHandler.Alias())
			fmt.Println("INFO:     - Task Title: " + taskHandler.Title())
			fmt.Println("INFO:     - Task Description: " + taskHandler.Description())
		}

		return true
	}

	// Finds the task and executes its handler
	for _, taskHandler := range store.TaskHandlerList() {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(alias)) {
			taskHandler.SetOptions(argumentsMap)
			taskHandler.Handle()
			return true
		}
	}

	fmt.Println("ERROR: Unrecognized task alias:", alias)
	return false
}

func unifyName(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return name
}

// taskHandlerFuncWithContext finds the TaskHandler and returns a function that
// checks if the handler implements TaskHandlerWithContext. If it does, it calls
// HandleWithContext(ctx), otherwise it falls back to Handle() for backward compatibility.
func (store *Store) taskHandlerFuncWithContext(taskAlias string, ctx context.Context) func(queuedTask TaskQueueInterface) bool {
	unifyName := func(name string) string {
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "_", "")
		return name
	}

	for _, taskHandler := range store.taskHandlers {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(taskAlias)) {
			return func(queuedTask TaskQueueInterface) bool {
				taskHandler.SetQueuedTask(queuedTask)

				// Check if handler implements TaskHandlerWithContext
				if contextHandler, ok := taskHandler.(TaskHandlerWithContext); ok {
					return contextHandler.HandleWithContext(ctx)
				}

				// Fall back to standard Handle() for backward compatibility
				return taskHandler.Handle()
			}
		}
	}

	return func(queuedTask TaskQueueInterface) bool {
		queuedTask.AppendDetails("No handler for alias: " + taskAlias)
		_ = store.TaskQueueUpdate(ctx, queuedTask)
		return false
	}
}

// argsToMap converts command line arguments to a key value map
// supports filled (i.e. --user=12) and unfilled (i.e. --force) arguments
func argsToMap(args []string) map[string]string {
	kv := map[string]string{}
	for i := 0; i < len(args); i++ {
		current := args[i]
		current = strings.TrimSpace(current)

		if strings.HasPrefix(current, "--") {
			if strings.Contains(current, "=") {
				currentArray := strings.Split(current, "=")
				if len(currentArray) < 2 {
					continue
				}
				kv[currentArray[0][2:]] = currentArray[1]
			} else {
				next := ""
				if len(args) > i+1 {
					next = args[i+1]
					next = strings.TrimSpace(next)
				}

				if strings.HasPrefix(next, "--") {
					kv[current[2:]] = ""
					continue
				}
				kv[current[2:]] = next
			}
		}
	}
	return kv
}
