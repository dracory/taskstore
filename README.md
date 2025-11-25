# Task Store <a href="https://gitpod.io/#https://github.com/gouniverse/taskstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>


[![Tests Status](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/taskstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/taskstore)](https://goreportcard.com/report/github.com/gouniverse/taskstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/taskstore)](https://pkg.go.dev/github.com/gouniverse/taskstore)

TaskStore is a robust, asynchronous durable task queue package designed to offload time-consuming or resource-intensive operations from your main application.

By deferring tasks to the background, you can improve application responsiveness and prevent performance bottlenecks.

TaskStore leverages a durable database (SQLite, MySQL, or PostgreSQL) to ensure reliable persistence and fault tolerance.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation

```bash
go get github.com/dracory/taskstore
```

## Queue Features

### Atomic Task Claiming
Tasks are claimed atomically using database transactions with `SELECT FOR UPDATE`, preventing race conditions where multiple workers might process the same task simultaneously.

### Concurrency Control
- **Default limit**: 10 concurrent tasks per queue
- **Configurable**: Set via `MaxConcurrency` in `NewStoreOptions`
- **Semaphore-based**: Automatic backpressure when limit is reached

```golang
store, err := taskstore.NewStore(taskstore.NewStoreOptions{
    DB:                      databaseInstance,
    TaskDefinitionTableName: "task_definition",
    TaskQueueTableName:      "task_queue",
    MaxConcurrency:          20, // Allow 20 concurrent tasks
})
```

### Graceful Shutdown
- `TaskQueueStop()` – Stop default queue and wait for all tasks to complete
- `TaskQueueStopByName(queueName)` – Stop specific queue and wait for all tasks
- Ensures no task goroutines are abandoned

```golang
// Start concurrent queue
store.TaskQueueRunConcurrent(ctx, "emails", 10, 1)

// Later: gracefully stop and wait for completion
store.TaskQueueStopByName("emails")
```

### Error Handling
Configure custom error handlers for monitoring and alerting:

```golang
store.SetErrorHandler(func(queueName, taskID string, err error) {
    log.Printf("[ERROR] Queue: %s, Task: %s, Error: %v", queue Name, taskID, err)
    // Send to monitoring system
    metrics.RecordTaskError(queueName, taskID)
})
```

### Context Propagation (Optional)
Task handlers can optionally implement `TaskHandlerWithContext` to support cancellation:

```golang
func (h *EmailHandler) HandleWithContext(ctx context.Context) bool {
    select {
    case <-ctx.Done():
        h.LogInfo("Task cancelled")
        return false
    case <-time.After(5 * time.Second):
        // Send email...
        h.LogSuccess("Email sent")
        return true
    }
}
```

**Note**: Existing handlers without `HandleWithContext` continue to work - this is fully backward compatible.

## Setup

```golang
myTaskStore = taskstore.NewStore(taskstore.NewStoreOptions{
	DB:                 databaseInstance,
	TaskDefinitionTableName: "my_task_definition",
	TaskQueueTableName:      "my_task_queue",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
```
 
## Documentation

- [Package overview](./docs/overview.md)
- [Task definitions](./docs/task-definitions.md)
- [Task queues](./docs/task-queues.md)
- [Schedules](./docs/schedules.md)

## Task Definitions

The task definition specifies a unit of work to be completed. It can be performed immediately, 
or enqueued to the database and deferred for asynchronous processing, ensuring your
application remains responsive.

Each task definition is uniquely identified by an alias and provides a human-readable title and description.

Each task definition is uniquely identified by an alias that allows the task to be easily called. 
A human-readable title and description give the user more information on the task definition.

To define a task definition, implement the TaskHandlerInterface and provide a Handle method
that contains the task's logic.

Optionally, extend the TaskHandlerBase struct for additional features like parameter
retrieval.

Task definitions can be executed directly from the command line (CLI) or as part of a background task queue.

The tasks placed in the task queue will be processed at a specified interval.

```golang
package tasks

func NewHelloWorldTask() *HelloWorldTask {
	return &HelloWorldTask{}
}

type HelloWorldTask struct {
	taskstore.TaskHandlerBase
}

var _ taskstore.TaskHandlerInterface = (*HelloWorldTask)(nil) // verify it extends the task handler interface

func (task *HelloWorldTask) Alias() string {
	return "HelloWorldTask"
}

func (task *HelloWorldTask) Title() string {
	return "Hello World"
}

func (task *HelloWorldTask) Description() string {
	return "Say hello world"
}

// Enqueue. Optional shortcut to quickly add this task to the task queue
func (task *HelloWorldTask) Enqueue(name string) (taskQueue taskstore.TaskQueueInterface, err error) {
	return myTaskStore.TaskDefinitionEnqueueByAlias(taskstore.DefaultTaskQueue, task.Alias(), map[string]any{
		"name": name,
	})
}

func (task *HelloWorldTask) Handle() bool {
	name := handler.GetParam("name")

        // Optional to allow adding the task to the task queue manually. Useful while in development
	if !task.HasQueuedTask() && task.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(name)

		if err != nil {
			task.LogError("Error enqueuing task: " + err.Error())
		} else {
			task.LogSuccess("Task enqueued.")
		}
		
		return true
	}

        if name != "" {
		task.LogInfo("Hello" + name + "!")	
	} else {
		task.LogInfo("Hello World!")
	}

	return true
}
```
## Registering Task Definitions to the TaskStore

Registering the task definition to the task store will persist it in the database.

```
myTaskStore.TaskHandlerAdd(NewHelloWorldTask(), true)
```

## Executing Task Definitions in the Terminal

To add the option to execute tasks from the terminal add the following to your main method

```
myTaskStore.TaskDefinitionExecuteCli(args[1], args[1:])
```

Example:
```
go run . HelloWorldTask --name="Tom Jones"
```

## Adding the Task to the Task Queue

To add a task to the background task queue

```
myTaskStore.TaskDefinitionEnqueueByAlias(taskstore.DefaultTaskQueue, NewHelloWorldTask.Alias(), map[string]any{
	"name": name,
})
```

Or if you have defined an Enqueue method as in the example task above.
```
NewHelloWorldTask().Enqueue("Tom Jones")
```

## Starting the Task Queue

To start the task queue, use one of the queue run methods:

```golang
ctx := context.Background()

// Option 1: Run default queue (serial processing)
myTaskStore.TaskQueueRunDefault(ctx, 10, 2) // every 10s, unstuck after 2 mins

// Option 2: Run named queue with serial processing
myTaskStore.TaskQueueRunSerial(ctx, "emails", 10, 2)

// Option 3: Run named queue with concurrent processing (respects MaxConcurrency)
myTaskStore.TaskQueueRunConcurrent(ctx, "emails", 10, 2)
```

## Store Methods

- `AutoMigrate() error` – automigrates (creates) the task definition and task queue tables
- `EnableDebug(debug bool) StoreInterface` – enables / disables the debug option

## Task Definition Methods

- `TaskDefinitionCreate(task TaskDefinitionInterface) error` – creates a task definition
- `TaskDefinitionFindByAlias(alias string) (TaskDefinitionInterface, error)` – finds a task definition by alias
- `TaskDefinitionFindByID(id string) (TaskDefinitionInterface, error)` – finds a task definition by ID
- `TaskDefinitionList(options TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error)` – lists task definitions
- `TaskDefinitionUpdate(task TaskDefinitionInterface) error` – updates a task definition
- `TaskDefinitionSoftDelete(task TaskDefinitionInterface) error` – soft deletes a task definition

## Task Queue Methods

- `TaskQueueCreate(queue TaskQueueInterface) error` – creates a new queued task
- `TaskQueueDeleteByID(id string) error` – deletes a queued task by ID
- `TaskQueueFindByID(id string) (TaskQueueInterface, error)` – finds a queued task by ID
- `TaskQueueFail(queue TaskQueueInterface) error` – marks a queued task as failed
- `TaskQueueSoftDeleteByID(id string) error` – soft deletes a queued task by ID (populates the deleted_at field)
- `TaskQueueSuccess(queue TaskQueueInterface) error` – completes a queued task successfully
- `TaskQueueList(options TaskQueueQueryInterface) ([]TaskQueueInterface, error)` – lists the queued tasks
- `TaskQueueUpdate(queue TaskQueueInterface) error` – updates a queued task

## Frequently Asked Questions (FAQ)

### 1. What is TaskStore used for?
TaskStore is a versatile tool for offloading time-consuming or resource-intensive
tasks from your main application. By deferring these tasks to the background,
you can improve application responsiveness and prevent performance bottlenecks.

It's ideal for tasks like data processing, sending emails, generating reports,
or performing batch operations.

### 2. How does TaskStore work?
TaskStore creates a durable queue in your database (SQLite, MySQL, or PostgreSQL)
to store tasks. These tasks are then processed asynchronously by a background worker.
You can define tasks using a simple interface and schedule them to be executed
at specific intervals or on demand.

### 3. What are the benefits of using TaskStore?

- Improved application performance: Offload time-consuming tasks to prevent performance bottlenecks.
- Asynchronous processing: Execute tasks independently of your main application flow.
- Reliability: Ensure tasks are completed even if your application crashes.
- Flexibility: Schedule tasks to run at specific intervals or on demand.
- Ease of use: Define tasks using a simple interface and integrate with your existing application.

### 4. How do I create a task definition in TaskStore?
To create a task definition, you'll need to implement the TaskHandlerInterface and provide a Handle method that contains the task's logic. You can also extend the TaskHandlerBase struct for additional features.

### 5. How do I schedule a task to run in the background?
Use the TaskDefinitionEnqueueByAlias method to add a task to the background task queue. You can specify the interval at which the task queue is processed using the QueueRunGoroutine method.

### 6. Can I monitor the status of tasks?
Yes, TaskStore provides methods to list tasks, check their status, and view task details.

### 7. How does TaskStore handle task failures?
If a task fails, it can be retried automatically or marked as failed. You can customize the retry logic to suit your specific needs.

### 8. Is TaskStore suitable for large-scale applications?
Yes, TaskStore is designed to handle large volumes of tasks. It can be scaled horizontally by adding more worker processes.

### 9. Does TaskStore support different database systems?
Yes, TaskStore supports SQLite, MySQL, and PostgreSQL.

### 10. Can I customize TaskStore to fit my specific needs?
Yes, TaskStore is highly customizable. You can extend and modify the code to suit your requirements.

## Similar

- https://github.com/harshadmanglani/polaris
- https://github.com/bamzi/jobrunner
- https://github.com/rk/go-cron
- https://github.com/fieldryand/goflow
- https://github.com/go-co-op/gocron
- https://github.com/exograd/eventline
- https://github.com/ajvb/kala
- https://github.com/shiblon/taskstore
