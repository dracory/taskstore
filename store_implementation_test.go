package taskstore

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	dsn := ":memory:?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore() (*Store, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}
	return NewStore(NewStoreOptions{
		TaskTableName:      "task",
		QueueTableName:     "queue",
		DB:                 db,
		DbDriverName:       "sqlite",
		AutomigrateEnabled: false,
		DebugEnabled:       false,
	})
}

// func TestWithDb(t *testing.T) {
// 	db := InitDB("test.db")
// 	store, error := InitStore()

// 	f := WithDb(db)
// 	f(s)

// 	if s.db == nil {
// 		t.Fatalf("DB: Expected Initialized DB, received [%v]", s.db)
// 	}

// }

// func TestWithDefinitionTableName(t *testing.T) {
// 	s := InitStore()

// 	table_name := "test_taskTableName.db"
// 	f := WithDefinitionTableName(table_name)
// 	f(s)
// 	if s.taskDefinitionTableName != table_name {
// 		t.Fatalf("Expected DefinitionTableName [%v], received [%v]", table_name, s.taskDefinitionTableName)
// 	}
// 	table_name = "Table2"
// 	f = WithDefinitionTableName(table_name)
// 	f(s)
// 	if s.taskDefinitionTableName != table_name {
// 		t.Fatalf("Expected DefinitionTableName [%v], received [%v]", table_name, s.taskDefinitionTableName)
// 	}
// }

// func TestWithTaskTableName(t *testing.T) {
// 	s := InitStore()

// 	table_name := "test_taskTableName.db"
// 	f := WithTaskTableName(table_name)
// 	f(s)
// 	if s.taskTaskTableName != table_name {
// 		t.Fatalf("Expected TaskTableName [%v], received [%v]", table_name, s.taskTaskTableName)
// 	}
// 	table_name = "Table2"
// 	f = WithTaskTableName(table_name)
// 	f(s)
// 	if s.taskTaskTableName != table_name {
// 		t.Fatalf("Expected TaskTableName [%v], received [%v]", table_name, s.taskTaskTableName)
// 	}
// }

// func TestWithDebug(t *testing.T) {
// 	s := InitStore()

// 	b := false
// 	f := WithDebug(b)
// 	f(s)
// 	if s.debug != b {
// 		t.Fatalf("Expected Debug [%v], received [%v]", b, s.debug)
// 	}
// }

// func Test_Store_DriverName(t *testing.T) {
// 	db := InitDB("sqlite")
// 	store := InitStore()
// 	s := store.DriverName(db)
// 	if s != "sqlite" {
// 		t.Fatalf("Expected Debug [%v], received [%v]", "sqlite", s)
// 	}
// }
