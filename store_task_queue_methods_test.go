package taskstore

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dracory/sb"
)

func Test_Store_SqlCreateTaskQueueTable(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("SqlCreateTaskQueueTable: Error[%v]", err)
	}

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("SqlCreateTaskQueueTable: Unexpected Query, received [%v]", query)
	}
}

func Test_Store_TaskQueueCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueCreate: Error[%v]", err)
	}

	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueCreate: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueCreate: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueCreate: Error in Creating TaskQueue: received [%v]", err)
	}
}

func Test_Store_TaskQueueDeleteByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueList: Error[%v]", err)
	}

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueList: Table creation error: [%v]", err)
	}

	queuedTask := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1).
		SetStatus(TaskQueueStatusQueued)

	err = store.TaskQueueCreate(context.Background(), queuedTask)

	if err != nil {
		t.Fatal("TaskQueueList: Error in creating queued task:", err.Error())
	}

	foundQueuedTask, err := store.TaskQueueFindByID(context.Background(), queuedTask.ID())

	if err != nil {
		t.Fatal("TaskQueueDeletedByID: Error in creating queued task:", err.Error())
	}

	if foundQueuedTask == nil {
		t.Fatal("TaskQueueDeletedByID: queued task not found:")
	}

	err = store.TaskQueueDeleteByID(context.Background(), queuedTask.ID())

	if err != nil {
		t.Error("TaskQueueDeletedByID: Error deleting queued task:", err.Error())
	}

}

func Test_Store_TaskQueueFail(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueFail: Error[%v]", err)
	}

	queuedTask := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueFail: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueFail: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), queuedTask)
	if err != nil {
		t.Fatalf("TaskQueueFail: Error in Creating TaskQueue: received [%v]", err)
	}

	err = store.TaskQueueFail(context.Background(), queuedTask)
	if err != nil {
		t.Fatalf("TaskQueueFail: Error in Fail TaskQueue: received [%v]", err)
	}
}

func Test_Store_TaskQueueFindByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueFindByID: Error[%v]", err)
	}
	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueFindByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueFindByID: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueFindByID: Error in Creating TaskQueue: received [%v]", err)
	}

	id := task.ID()
	queue, err := store.TaskQueueFindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("TaskQueueFindByID: Error in TaskQueueFindByID: received [%v]", err)
	}

	if queue == nil {
		t.Fatalf("TaskQueueFindByID: Error in Finding TaskQueue: ID [%v]", id)
	}
	if queue.ID() != id {
		t.Fatalf("TaskQueueFindByID: ID not matching, Expected[%v], Received[%v]", id, queue.ID())
	}
}

func Test_Store_TaskQueueList(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueList: Error[%v]", err)
	}

	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1).
		SetStatus(TaskQueueStatusQueued)

	query := store.SqlCreateTaskQueueTable()

	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)

	if err != nil {
		t.Fatalf("TaskQueueList: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)

	if err != nil {
		t.Fatalf("TaskQueueList: Error in Creating TaskQueue: received [%v]", err)
	}

	list, err := store.TaskQueueList(context.Background(), TaskQueueQuery().
		SetStatus(TaskQueueStatusQueued).
		SetLimit(1).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder(ASC))

	if err != nil {
		t.Fatalf("TaskQueueList: Error[%v]", err)
	}

	if len(list) != 1 {
		t.Fatal("There must be 1 task, found: ", list)
	}
}

func Test_Store_TaskQueueFindNextQueuedTaskByQueue(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Error[%v]", err)
	}

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Table creation error: [%v]", err)
	}

	// default queue task
	defaultTask := NewTaskQueue().
		SetTaskID("TASK_DEFAULT").
		SetAttempts(1).
		SetStatus(TaskQueueStatusQueued)

	// named queue task
	namedTask := NewTaskQueue().
		SetTaskID("TASK_EMAILS").
		SetAttempts(1).
		SetStatus(TaskQueueStatusQueued).
		SetQueueName("emails")

	if err := store.TaskQueueCreate(context.Background(), defaultTask); err != nil {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Error creating default task: [%v]", err)
	}

	if err := store.TaskQueueCreate(context.Background(), namedTask); err != nil {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Error creating named task: [%v]", err)
	}

	q, err := store.TaskQueueFindNextQueuedTaskByQueue(context.Background(), "emails")
	if err != nil {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Error[%v]", err)
	}

	if q == nil {
		t.Fatal("TaskQueueFindNextQueuedTaskByQueue: Expected a queued task for 'emails' queue, got nil")
	}

	if q.TaskID() != "TASK_EMAILS" {
		t.Fatalf("TaskQueueFindNextQueuedTaskByQueue: Expected TASK_EMAILS, got %s", q.TaskID())
	}
}

func Test_Store_TaskQueueSoftDeleteByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueSoftDeleteByID: Error[%v]", err)
	}

	queuedTask := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueSoftDeleteByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueSoftDeleteByID: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), queuedTask)
	if err != nil {
		t.Fatalf("TaskQueueSoftDeleteByID: Error in Creating TaskQueue: received [%v]", err)
	}

	err = store.TaskQueueSoftDeleteByID(context.Background(), queuedTask.ID())
	if err != nil {
		t.Fatalf("TaskQueueSoftDeleteByID: Error in Fail TaskQueue: received [%v]", err)
	}

	queueFound, err := store.TaskQueueFindByID(context.Background(), queuedTask.ID())

	if err != nil {
		t.Fatal("TaskQueueSoftDeleteByID: Error in TaskQueueFindByID: received:", err)
	}

	if queueFound != nil {
		t.Fatal("TaskQueueSoftDeleteByID: TaskQueueFindByID should be nil, received:", queueFound)
	}
}

func Test_Store_TaskQueueSuccess(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueSuccess: Error[%v]", err)
	}

	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueSuccess: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueSuccess: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueSuccess: Error in Creating TaskQueue: received [%v]", err)
	}

	err = store.TaskQueueSuccess(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueSuccess: Error in Success TaskQueue: received [%v]", err)
	}
}

func Test_Store_TaskQueueUpdate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueUpdate: Error[%v]", err)
	}

	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueUpdate: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueUpdate: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueUpdate: Error in Creating TaskQueue: received [%v]", err)
	}

	err = store.TaskQueueUpdate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueUpdate: Error in Updating TaskQueue: received [%v]", err)
	}
}

func Test_Store_TaskQueue_AppendDetails(t *testing.T) {
	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	str := "Test1"
	task.AppendDetails(str)

	if !strings.Contains(task.Details(), str) {
		t.Fatalf("AppendDetails: Failed Details[%v]", task.Details())
	}
}

type Temp struct {
	Status     string `json:"status"`
	Limit      int    `json:"limit"`
	Sort_by    string `json:"sort_by"`
	Sort_order string `json:"sort_order"`
}

func Test_TaskQueue_ParametersMap(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("GetParameters: Error[%v]", err)
	}

	task := NewTaskQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("GetParameters: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("GetParameters: Table creation error: [%v]", err)
	}

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("GetParameters: Error in Creating TaskQueue: received [%v]", err)
	}

	u, err := json.Marshal(Temp{Status: "Bob", Limit: 10})

	if err != nil {
		t.Fatalf("%v", err)
	}

	task.SetParameters(string(u))

	err = json.Unmarshal([]byte(task.Parameters()), &Temp{})
	if err != nil {
		t.Fatalf("GetParameters: Error[%v]", err)
	}
}

// TestQueuedTaskForceFail_WithNullDateTime verifies that tasks with
// started_at set to sb.NULL_DATETIME are not incorrectly marked as failed
func TestQueuedTaskForceFail_WithNullDateTime(t *testing.T) {
	// Setup: Create a test store
	store, err := initStore()
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("Table creation error: [%v]", err)
	}

	// Test Case 1: Task with NULL_DATETIME should NOT be force-failed
	t.Run("TaskWithNullDateTime_ShouldNotBeFailed", func(t *testing.T) {
		queue := NewTaskQueue()
		queue.SetTaskID("test-task-1")
		queue.SetStatus(TaskQueueStatusQueued)
		// started_at and completed_at are already set to sb.NULL_DATETIME by NewTaskQueue

		err := store.TaskQueueCreate(context.Background(), queue)
		if err != nil {
			t.Fatalf("Failed to create queue: %v", err)
		}

		// Try to force fail with 1 minute timeout
		err = store.QueuedTaskForceFail(context.Background(), queue, 1)
		if err != nil {
			t.Errorf("QueuedTaskForceFail returned error: %v", err)
		}

		// Verify status is still queued (not failed)
		if queue.Status() != TaskQueueStatusQueued {
			t.Errorf("Expected status to remain 'queued', got '%s'", queue.Status())
		}
	})

	// Test Case 2: Task with empty started_at should NOT be force-failed
	t.Run("TaskWithEmptyStartedAt_ShouldNotBeFailed", func(t *testing.T) {
		queue := NewTaskQueue()
		queue.SetTaskID("test-task-2")
		queue.SetStatus(TaskQueueStatusQueued)
		queue.SetStartedAt("") // Empty string

		err := store.QueuedTaskForceFail(context.Background(), queue, 1)
		if err != nil {
			t.Errorf("QueuedTaskForceFail returned error: %v", err)
		}

		// Verify status is still queued (not failed)
		if queue.Status() != TaskQueueStatusQueued {
			t.Errorf("Expected status to remain 'queued', got '%s'", queue.Status())
		}
	})
}

func Test_Store_TaskQueueClaimNext(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskQueueClaimNext: Error[%v]", err)
	}

	query := store.SqlCreateTaskQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskQueueClaimNext: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskQueueClaimNext: Table creation error: [%v]", err)
	}

	// Create a queued task
	queueName := "test-queue"
	task := NewTaskQueue().
		SetTaskID("TASK_CLAIM").
		SetQueueName(queueName).
		SetStatus(TaskQueueStatusQueued).
		SetAttempts(0)

	err = store.TaskQueueCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskQueueClaimNext: Error in Creating TaskQueue: received [%v]", err)
	}

	// Claim the task
	claimedTask, err := store.TaskQueueClaimNext(context.Background(), queueName)
	if err != nil {
		t.Fatalf("TaskQueueClaimNext: Error claiming task: [%v]", err)
	}

	if claimedTask == nil {
		t.Fatal("TaskQueueClaimNext: Expected to claim a task, got nil")
	}

	// Verify the claimed task properties
	if claimedTask.ID() != task.ID() {
		t.Errorf("TaskQueueClaimNext: Expected task ID %s, got %s", task.ID(), claimedTask.ID())
	}

	if claimedTask.Status() != TaskQueueStatusRunning {
		t.Errorf("TaskQueueClaimNext: Expected status %s, got %s", TaskQueueStatusRunning, claimedTask.Status())
	}

	if claimedTask.StartedAt() == "" || claimedTask.StartedAt() == sb.NULL_DATETIME {
		t.Error("TaskQueueClaimNext: Expected StartedAt to be set")
	}

	// Verify database state
	dbTask, err := store.TaskQueueFindByID(context.Background(), task.ID())
	if err != nil {
		t.Fatalf("TaskQueueClaimNext: Error finding task: [%v]", err)
	}

	if dbTask.Status() != TaskQueueStatusRunning {
		t.Errorf("TaskQueueClaimNext: Database status expected %s, got %s", TaskQueueStatusRunning, dbTask.Status())
	}
}
