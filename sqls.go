package taskstore

import "github.com/dracory/sb"

// SqlCreateQueueTable returns a SQL string for creating the Queue table
func (st *Store) SqlCreateQueueTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.queueTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_TASK_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:     COLUMN_PARAMETERS,
			Type:     sb.COLUMN_TYPE_TEXT,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_OUTPUT,
			Type:     sb.COLUMN_TYPE_TEXT,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_DETAILS,
			Type:     sb.COLUMN_TYPE_TEXT,
			Nullable: true,
		}).
		Column(sb.Column{
			Name: COLUMN_ATTEMPTS,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name:     COLUMN_STARTED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_COMPLETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
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
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	return sql
}

// SqlCreateTaskTable returns a SQL string for creating the Task table
func (st *Store) SqlCreateTaskTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.taskTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_ALIAS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
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
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	return sql
}
