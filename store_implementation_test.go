package taskstore

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func initDB(filename ...string) (*sql.DB, error) {
	// Use shared cache mode to allow concurrent goroutines to access the same in-memory database
	// Note: Must use file:Name?mode=memory&cache=shared to allow sharing within the test but isolation between tests
	dsn := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared&parseTime=true", time.Now().UnixNano())
	if len(filename) > 0 {
		// For file-based databases, use WAL mode and busy timeout for concurrent access
		// Use _pragma for modernc.org/sqlite
		dsn = filename[0] + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
		// Remove the file if it exists to ensure clean state
		if err := os.Remove(filename[0]); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	// Configure connection pool for concurrent access
	// Explicitly set pragmas to be safe
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)   // Allow multiple concurrent connections
	db.SetMaxIdleConns(25)   // Keep idle connections to prevent in-memory DB loss
	db.SetConnMaxLifetime(0) // Connections never expire

	return db, nil
}

func initStore(filename ...string) (*Store, error) {
	db, err := initDB(filename...)
	if err != nil {
		return nil, err
	}
	return NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definition",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		DbDriverName:            "sqlite",
		AutomigrateEnabled:      true, // Enable automigration for tests
		DebugEnabled:            false,
	})
}
