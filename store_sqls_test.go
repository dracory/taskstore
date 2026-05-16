package taskstore

import (
	"testing"
)

func TestSqlCreateTaskQueueTable(t *testing.T) {
	tests := []struct {
		name       string
		driverName string
		wantErr    bool
	}{
		{
			name:       "sqlite driver",
			driverName: "sqlite",
			wantErr:    false,
		},
		{
			name:       "mysql driver",
			driverName: "mysql",
			wantErr:    false,
		},
		{
			name:       "postgres driver",
			driverName: "postgres",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Store{
				dbDriverName:         tt.driverName,
				taskQueueTableName:   "task_queue",
				taskDefinitionTableName: "task_definition",
				scheduleTableName:     "schedule",
			}
			got, err := st.SqlCreateTaskQueueTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("SqlCreateTaskQueueTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("SqlCreateTaskQueueTable() returned empty string")
			}
		})
	}
}

func TestSqlCreateTaskDefinitionTable(t *testing.T) {
	tests := []struct {
		name       string
		driverName string
		wantErr    bool
	}{
		{
			name:       "sqlite driver",
			driverName: "sqlite",
			wantErr:    false,
		},
		{
			name:       "mysql driver",
			driverName: "mysql",
			wantErr:    false,
		},
		{
			name:       "postgres driver",
			driverName: "postgres",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Store{
				dbDriverName:         tt.driverName,
				taskQueueTableName:   "task_queue",
				taskDefinitionTableName: "task_definition",
				scheduleTableName:     "schedule",
			}
			got, err := st.SqlCreateTaskDefinitionTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("SqlCreateTaskDefinitionTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("SqlCreateTaskDefinitionTable() returned empty string")
			}
		})
	}
}

func TestSqlCreateScheduleTable(t *testing.T) {
	tests := []struct {
		name       string
		driverName string
		wantErr    bool
	}{
		{
			name:       "sqlite driver",
			driverName: "sqlite",
			wantErr:    false,
		},
		{
			name:       "mysql driver",
			driverName: "mysql",
			wantErr:    false,
		},
		{
			name:       "postgres driver",
			driverName: "postgres",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Store{
				dbDriverName:         tt.driverName,
				taskQueueTableName:   "task_queue",
				taskDefinitionTableName: "task_definition",
				scheduleTableName:     "schedule",
			}
			got, err := st.SqlCreateScheduleTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("SqlCreateScheduleTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("SqlCreateScheduleTable() returned empty string")
			}
		})
	}
}
