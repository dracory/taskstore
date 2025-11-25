package taskstore

import "github.com/dracory/sb"

// SqlCreateTaskQueueTable - creates the task queue table
func (st *Store) SqlCreateTaskQueueTable() string {
	return sb.NewBuilder(st.dbDriverName).
		Table(st.taskQueueTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     50,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUEUE_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_TASK_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name: COLUMN_PARAMETERS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name: COLUMN_OUTPUT,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_DETAILS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_ATTEMPTS,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name: COLUMN_STARTED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_COMPLETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()
}

// SqlCreateTaskDefinitionTable - creates the task definition table
func (st *Store) SqlCreateTaskDefinitionTable() string {
	return sb.NewBuilder(st.dbDriverName).
		Table(st.taskDefinitionTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     50,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name:   COLUMN_ALIAS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
			Unique: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_DESCRIPTION,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_IS_RECURRING,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name:   COLUMN_RECURRENCE_RULE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 500,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()
}

// SqlCreateScheduleTable - creates the schedule table
func (st *Store) SqlCreateScheduleTable() string {
	return sb.NewBuilder(st.dbDriverName).
		Table(st.scheduleTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     50,
		}).
		Column(sb.Column{
			Name:   COLUMN_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_DESCRIPTION,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name:   COLUMN_RECURRENCE_RULE,
			Type:   sb.COLUMN_TYPE_TEXT,
			Length: 1000,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUEUE_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_TASK_DEFINITION_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name: COLUMN_PARAMETERS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_START_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_END_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_EXECUTION_COUNT,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name: COLUMN_MAX_EXECUTION_COUNT,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name: COLUMN_LAST_RUN_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_NEXT_RUN_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()
}
