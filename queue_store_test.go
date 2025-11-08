package taskstore

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test_Store_SqlCreateQueueTable(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("SqlCreateQueueTable: Error[%v]", err)
	}

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("SqlCreateQueueTable: Unexpected Query, received [%v]", query)
	}
}

func Test_Store_QueueCreate(t *testing.T) {
	store, err := InitStore("test_queue_create.db")
	if err != nil {
		t.Fatalf("QueueCreate: Error[%v]", err)
	}

	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueCreate: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueCreate: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueCreate: Error in Creating Queue: received [%v]", err)
	}
}

func Test_Store_QueueDeleteByID(t *testing.T) {
	store, err := InitStore("test_queue_delete_by_id.db")
	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueList: Table creation error: [%v]", err)
	}

	queuedTask := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1).
		SetStatus(QueueStatusQueued)

	err = store.QueueCreate(queuedTask)

	if err != nil {
		t.Fatal("QueueList: Error in creating queued task:", err.Error())
	}

	foundQueuedTask, err := store.QueueFindByID(queuedTask.ID())

	if err != nil {
		t.Fatal("QueueDeletedByID: Error in creating queued task:", err.Error())
	}

	if foundQueuedTask == nil {
		t.Fatal("QueueDeletedByID: queued task not found:")
	}

	err = store.QueueDeleteByID(queuedTask.ID())

	if err != nil {
		t.Error("QueueDeletedByID: Error deleting queued task:", err.Error())
	}

}

func Test_Store_QueueFail(t *testing.T) {
	store, err := InitStore("test_queue_fail.db")
	if err != nil {
		t.Fatalf("QueueFail: Error[%v]", err)
	}

	queuedTask := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueFail: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueFail: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(queuedTask)
	if err != nil {
		t.Fatalf("QueueFail: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueFail(queuedTask)
	if err != nil {
		t.Fatalf("QueueFail: Error in Fail Queue: received [%v]", err)
	}
}

func Test_Store_QueueFindByID(t *testing.T) {
	store, err := InitStore("test_queue_find_by_id.db")
	if err != nil {
		t.Fatalf("QueueFindByID: Error[%v]", err)
	}
	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueFindByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueFindByID: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueFindByID: Error in Creating Queue: received [%v]", err)
	}

	id := task.ID()
	queue, err := store.QueueFindByID(id)
	if err != nil {
		t.Fatalf("QueueFindByID: Error in QueueFindByID: received [%v]", err)
	}

	if queue == nil {
		t.Fatalf("QueueFindByID: Error in Finding Queue: ID [%v]", id)
	}
	if queue.ID() != id {
		t.Fatalf("QueueFindByID: ID not matching, Expected[%v], Received[%v]", id, queue.ID())
	}
}

func Test_Store_QueueList(t *testing.T) {
	store, err := InitStore("test_queue_list.db")
	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1).
		SetStatus(QueueStatusQueued)

	query := store.SqlCreateQueueTable()

	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueList: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)

	if err != nil {
		t.Fatalf("QueueList: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)

	if err != nil {
		t.Fatalf("QueueList: Error in Creating Queue: received [%v]", err)
	}

	list, err := store.QueueList(QueueQuery().
		SetStatus(QueueStatusQueued).
		SetLimit(1).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder(ASC))

	if err != nil {
		t.Fatalf("QueueList: Error[%v]", err)
	}

	if len(list) != 1 {
		t.Fatal("There must be 1 task, found: ", list)
	}
}

func Test_Store_QueueSoftDeleteByID(t *testing.T) {
	store, err := InitStore("test_queue_soft_delete_by_id.db")
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error[%v]", err)
	}

	queuedTask := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueSoftDeleteByID: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(queuedTask)
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueSoftDeleteByID(queuedTask.ID())
	if err != nil {
		t.Fatalf("QueueSoftDeleteByID: Error in Fail Queue: received [%v]", err)
	}

	queueFound, err := store.QueueFindByID(queuedTask.ID())

	if err != nil {
		t.Fatal("QueueSoftDeleteByID: Error in QueueFindByID: received:", err)
	}

	if queueFound != nil {
		t.Fatal("QueueSoftDeleteByID: QueueFindByID should be nil, received:", queueFound)
	}
}

func Test_Store_QueueSuccess(t *testing.T) {
	store, err := InitStore("test_queue_success.db")
	if err != nil {
		t.Fatalf("QueueSuccess: Error[%v]", err)
	}

	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueSuccess: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueSuccess: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueSuccess: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueSuccess(task)
	if err != nil {
		t.Fatalf("QueueSuccess: Error in Success Queue: received [%v]", err)
	}
}

func Test_Store_QueueUpdate(t *testing.T) {
	store, err := InitStore("test_queue_update.db")
	if err != nil {
		t.Fatalf("QueueUpdate: Error[%v]", err)
	}

	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueUpdate: UnExpected Query, received [%v]", query)
	}
	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueUpdate: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("QueueUpdate: Error in Creating Queue: received [%v]", err)
	}

	err = store.QueueUpdate(task)
	if err != nil {
		t.Fatalf("QueueUpdate: Error in Updating Queue: received [%v]", err)
	}
}

func Test_Store_Queue_AppendDetails(t *testing.T) {
	task := NewQueue().
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

func Test_Queue_ParametersMap(t *testing.T) {
	store, err := InitStore("test_queue_get_parameters.db")
	if err != nil {
		t.Fatalf("GetParameters: Error[%v]", err)
	}

	task := NewQueue().
		SetTaskID("TASK_01").
		SetAttempts(1)

	query := store.SqlCreateQueueTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("GetParameters: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("GetParameters: Table creation error: [%v]", err)
	}

	err = store.QueueCreate(task)
	if err != nil {
		t.Fatalf("GetParameters: Error in Creating Queue: received [%v]", err)
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
