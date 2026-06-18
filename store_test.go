package taskstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
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

func Test_sleepWithContext(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		cancel   bool
	}{
		{
			name:     "zero duration returns immediately",
			duration: 0,
			cancel:   false,
		},
		{
			name:     "negative duration returns immediately",
			duration: -1,
			cancel:   false,
		},
		{
			name:     "positive duration waits",
			duration: 10 * time.Millisecond,
			cancel:   false,
		},
		{
			name:     "context cancellation",
			duration: 100 * time.Millisecond,
			cancel:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.cancel {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
				go func() {
					time.Sleep(5 * time.Millisecond)
					cancel()
				}()
			}

			start := time.Now()
			got := sleepWithContext(ctx, tt.duration)
			elapsed := time.Since(start)

			if tt.cancel {
				if got {
					t.Errorf("sleepWithContext() should return false when context is cancelled")
				}
				if elapsed > 50*time.Millisecond {
					t.Errorf("sleepWithContext() should return quickly when cancelled, took %v", elapsed)
				}
			} else if tt.duration <= 0 {
				if !got {
					t.Errorf("sleepWithContext() should return true for zero/negative duration")
				}
				if elapsed > 10*time.Millisecond {
					t.Errorf("sleepWithContext() should return immediately for zero/negative duration, took %v", elapsed)
				}
			} else {
				if !got {
					t.Errorf("sleepWithContext() should return true for successful sleep")
				}
				if elapsed < tt.duration {
					t.Errorf("sleepWithContext() should wait at least %v, took %v", tt.duration, elapsed)
				}
			}
		})
	}
}

func Test_unifyName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes dashes",
			input: "my-task-name",
			want:  "mytaskname",
		},
		{
			name:  "removes underscores",
			input: "my_task_name",
			want:  "mytaskname",
		},
		{
			name:  "removes both dashes and underscores",
			input: "my-task_name",
			want:  "mytaskname",
		},
		{
			name:  "handles empty string",
			input: "",
			want:  "",
		},
		{
			name:  "handles string without special chars",
			input: "mytaskname",
			want:  "mytaskname",
		},
		{
			name:  "handles mixed case",
			input: "My-Task_Name",
			want:  "MyTaskName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unifyName(tt.input); got != tt.want {
				t.Errorf("unifyName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argsToMap(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string]string
	}{
		{
			name: "empty args returns empty map",
			args: []string{},
			want: map[string]string{},
		},
		{
			name: "filled arguments",
			args: []string{"--user=12", "--force"},
			want: map[string]string{"user": "12", "force": ""},
		},
		{
			name: "unfilled arguments",
			args: []string{"--force", "--verbose"},
			want: map[string]string{"force": "", "verbose": ""},
		},
		{
			name: "mixed filled and unfilled",
			args: []string{"--user=12", "--force", "--timeout=30"},
			want: map[string]string{"user": "12", "force": "", "timeout": "30"},
		},
		{
			name: "arguments with spaces",
			args: []string{"--user=12", "  --force  "},
			want: map[string]string{"user": "12", "force": ""},
		},
		{
			name: "arguments with equals in value",
			args: []string{"--url=http://example.com?a=1"},
			want: map[string]string{"url": "http://example.com?a"}, // Function splits on first =
		},
		{
			name: "ignores non-flag arguments",
			args: []string{"--user=12", "positional", "--force"},
			want: map[string]string{"user": "12", "force": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := argsToMap(tt.args)
			if len(got) != len(tt.want) {
				t.Errorf("argsToMap() length = %v, want %v", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("argsToMap() key %s = %v, want %v", k, got[k], v)
				}
			}
		})
	}
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

		AutomigrateEnabled: true, // Enable automigration for tests
		DebugEnabled:       false,
	})
}
