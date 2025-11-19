package taskstore

import (
	"testing"
)

func TestTaskQuery(t *testing.T) {
	query := TaskDefinitionQuery()

	if query == nil {
		t.Fatal("TaskQuery: Expected query to be created, got nil")
	}

	// Test that it implements the interface
	var _ TaskDefinitionQueryInterface = query
}

func TestTaskQuery_Validate(t *testing.T) {
	tests := []struct {
		name        string
		setupQuery  func() TaskDefinitionQueryInterface
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid empty query",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery()
			},
			expectError: false,
		},
		{
			name: "valid query with all fields",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().
					SetAlias("test-alias").
					SetCreatedAtGte("2023-01-01 00:00:00").
					SetCreatedAtLte("2023-12-31 23:59:59").
					SetID("test-id").
					SetIDIn([]string{"id1", "id2"}).
					SetLimit(10).
					SetOffset(0).
					SetStatus("active").
					SetStatusIn([]string{"active", "canceled"})
			},
			expectError: false,
		},
		{
			name: "empty alias",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetAlias("")
			},
			expectError: true,
			errorMsg:    "task query. alias cannot be empty",
		},
		{
			name: "empty created_at_gte",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetCreatedAtGte("")
			},
			expectError: true,
			errorMsg:    "task query. created_at_gte cannot be empty",
		},
		{
			name: "empty created_at_lte",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetCreatedAtLte("")
			},
			expectError: true,
			errorMsg:    "task query. created_at_lte cannot be empty",
		},
		{
			name: "empty id",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetID("")
			},
			expectError: true,
			errorMsg:    "task query. id cannot be empty",
		},
		{
			name: "empty id_in array",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetIDIn([]string{})
			},
			expectError: true,
			errorMsg:    "task query. id_in cannot be empty array",
		},
		{
			name: "negative limit",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetLimit(-1)
			},
			expectError: true,
			errorMsg:    "task query. limit cannot be negative",
		},
		{
			name: "negative offset",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetOffset(-1)
			},
			expectError: true,
			errorMsg:    "task query. offset cannot be negative",
		},
		{
			name: "empty status",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetStatus("")
			},
			expectError: true,
			errorMsg:    "task query. status cannot be empty",
		},
		{
			name: "empty status_in array",
			setupQuery: func() TaskDefinitionQueryInterface {
				return TaskDefinitionQuery().SetStatusIn([]string{})
			},
			expectError: true,
			errorMsg:    "task query. status_in cannot be empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.setupQuery()
			err := query.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate: Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Validate: Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Validate: Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestTaskQuery_Alias(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasAlias() {
		t.Error("HasAlias: Expected false for new query")
	}
	if query.Alias() != "" {
		t.Errorf("Alias: Expected empty string, got '%s'", query.Alias())
	}

	// Test setting alias
	testAlias := "test-alias"
	result := query.SetAlias(testAlias)
	if result != query {
		t.Error("SetAlias: Expected method to return the same query instance")
	}
	if !query.HasAlias() {
		t.Error("HasAlias: Expected true after setting alias")
	}
	if query.Alias() != testAlias {
		t.Errorf("Alias: Expected '%s', got '%s'", testAlias, query.Alias())
	}
}

func TestTaskQuery_Columns(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	columns := query.Columns()
	if len(columns) != 0 {
		t.Errorf("Columns: Expected empty slice, got %v", columns)
	}

	// Test setting columns
	testColumns := []string{"id", "alias", "title"}
	result := query.SetColumns(testColumns)
	if result != query {
		t.Error("SetColumns: Expected method to return the same query instance")
	}

	retrievedColumns := query.Columns()
	if len(retrievedColumns) != len(testColumns) {
		t.Errorf("Columns: Expected %d columns, got %d", len(testColumns), len(retrievedColumns))
	}
	for i, col := range testColumns {
		if retrievedColumns[i] != col {
			t.Errorf("Columns: Expected column '%s' at index %d, got '%s'", col, i, retrievedColumns[i])
		}
	}
}

func TestTaskQuery_CountOnly(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasCountOnly() {
		t.Error("HasCountOnly: Expected false for new query")
	}
	if query.IsCountOnly() {
		t.Error("IsCountOnly: Expected false for new query")
	}

	// Test setting count only to true
	result := query.SetCountOnly(true)
	if result != query {
		t.Error("SetCountOnly: Expected method to return the same query instance")
	}
	if !query.HasCountOnly() {
		t.Error("HasCountOnly: Expected true after setting count only")
	}
	if !query.IsCountOnly() {
		t.Error("IsCountOnly: Expected true after setting count only to true")
	}

	// Test setting count only to false
	query.SetCountOnly(false)
	if !query.HasCountOnly() {
		t.Error("HasCountOnly: Expected true even when set to false")
	}
	if query.IsCountOnly() {
		t.Error("IsCountOnly: Expected false after setting count only to false")
	}
}

func TestTaskQuery_CreatedAtGte(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasCreatedAtGte() {
		t.Error("HasCreatedAtGte: Expected false for new query")
	}

	// Test setting created_at_gte
	testDate := "2023-01-01 00:00:00"
	result := query.SetCreatedAtGte(testDate)
	if result != query {
		t.Error("SetCreatedAtGte: Expected method to return the same query instance")
	}
	if !query.HasCreatedAtGte() {
		t.Error("HasCreatedAtGte: Expected true after setting created_at_gte")
	}
	if query.CreatedAtGte() != testDate {
		t.Errorf("CreatedAtGte: Expected '%s', got '%s'", testDate, query.CreatedAtGte())
	}
}

func TestTaskQuery_CreatedAtLte(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasCreatedAtLte() {
		t.Error("HasCreatedAtLte: Expected false for new query")
	}

	// Test setting created_at_lte
	testDate := "2023-12-31 23:59:59"
	result := query.SetCreatedAtLte(testDate)
	if result != query {
		t.Error("SetCreatedAtLte: Expected method to return the same query instance")
	}
	if !query.HasCreatedAtLte() {
		t.Error("HasCreatedAtLte: Expected true after setting created_at_lte")
	}
	if query.CreatedAtLte() != testDate {
		t.Errorf("CreatedAtLte: Expected '%s', got '%s'", testDate, query.CreatedAtLte())
	}
}

func TestTaskQuery_ID(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasID() {
		t.Error("HasID: Expected false for new query")
	}

	// Test setting ID
	testID := "test-id-123"
	result := query.SetID(testID)
	if result != query {
		t.Error("SetID: Expected method to return the same query instance")
	}
	if !query.HasID() {
		t.Error("HasID: Expected true after setting ID")
	}
	if query.ID() != testID {
		t.Errorf("ID: Expected '%s', got '%s'", testID, query.ID())
	}
}

func TestTaskQuery_IDIn(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasIDIn() {
		t.Error("HasIDIn: Expected false for new query")
	}

	// Test setting ID in
	testIDs := []string{"id1", "id2", "id3"}
	result := query.SetIDIn(testIDs)
	if result != query {
		t.Error("SetIDIn: Expected method to return the same query instance")
	}
	if !query.HasIDIn() {
		t.Error("HasIDIn: Expected true after setting ID in")
	}

	retrievedIDs := query.IDIn()
	if len(retrievedIDs) != len(testIDs) {
		t.Errorf("IDIn: Expected %d IDs, got %d", len(testIDs), len(retrievedIDs))
	}
	for i, id := range testIDs {
		if retrievedIDs[i] != id {
			t.Errorf("IDIn: Expected ID '%s' at index %d, got '%s'", id, i, retrievedIDs[i])
		}
	}
}

func TestTaskQuery_Limit(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasLimit() {
		t.Error("HasLimit: Expected false for new query")
	}

	// Test setting limit
	testLimit := 50
	result := query.SetLimit(testLimit)
	if result != query {
		t.Error("SetLimit: Expected method to return the same query instance")
	}
	if !query.HasLimit() {
		t.Error("HasLimit: Expected true after setting limit")
	}
	if query.Limit() != testLimit {
		t.Errorf("Limit: Expected %d, got %d", testLimit, query.Limit())
	}
}

func TestTaskQuery_Offset(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasOffset() {
		t.Error("HasOffset: Expected false for new query")
	}

	// Test setting offset
	testOffset := 25
	result := query.SetOffset(testOffset)
	if result != query {
		t.Error("SetOffset: Expected method to return the same query instance")
	}
	if !query.HasOffset() {
		t.Error("HasOffset: Expected true after setting offset")
	}
	if query.Offset() != testOffset {
		t.Errorf("Offset: Expected %d, got %d", testOffset, query.Offset())
	}
}

func TestTaskQuery_OrderBy(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasOrderBy() {
		t.Error("HasOrderBy: Expected false for new query")
	}

	// Test setting order by
	testOrderBy := "created_at"
	result := query.SetOrderBy(testOrderBy)
	if result != query {
		t.Error("SetOrderBy: Expected method to return the same query instance")
	}
	if !query.HasOrderBy() {
		t.Error("HasOrderBy: Expected true after setting order by")
	}
	if query.OrderBy() != testOrderBy {
		t.Errorf("OrderBy: Expected '%s', got '%s'", testOrderBy, query.OrderBy())
	}
}

func TestTaskQuery_SoftDeletedIncluded(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected false for new query")
	}
	if query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected false for new query")
	}

	// Test setting soft deleted included to true
	result := query.SetSoftDeletedIncluded(true)
	if result != query {
		t.Error("SetSoftDeletedIncluded: Expected method to return the same query instance")
	}
	if !query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected true after setting soft deleted included")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected true after setting to true")
	}

	// Test setting soft deleted included to false
	query.SetSoftDeletedIncluded(false)
	if !query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected true even when set to false")
	}
	if query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected false after setting to false")
	}
}

func TestTaskQuery_SortOrder(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasSortOrder() {
		t.Error("HasSortOrder: Expected false for new query")
	}

	// Test setting sort order
	testSortOrder := "DESC"
	result := query.SetSortOrder(testSortOrder)
	if result != query {
		t.Error("SetSortOrder: Expected method to return the same query instance")
	}
	if !query.HasSortOrder() {
		t.Error("HasSortOrder: Expected true after setting sort order")
	}
	if query.SortOrder() != testSortOrder {
		t.Errorf("SortOrder: Expected '%s', got '%s'", testSortOrder, query.SortOrder())
	}
}

func TestTaskQuery_Status(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasStatus() {
		t.Error("HasStatus: Expected false for new query")
	}

	// Test setting status
	testStatus := "active"
	result := query.SetStatus(testStatus)
	if result != query {
		t.Error("SetStatus: Expected method to return the same query instance")
	}
	if !query.HasStatus() {
		t.Error("HasStatus: Expected true after setting status")
	}
	if query.Status() != testStatus {
		t.Errorf("Status: Expected '%s', got '%s'", testStatus, query.Status())
	}
}

func TestTaskQuery_StatusIn(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test default state
	if query.HasStatusIn() {
		t.Error("HasStatusIn: Expected false for new query")
	}

	// Test setting status in
	testStatuses := []string{"active", "canceled", "paused"}
	result := query.SetStatusIn(testStatuses)
	if result != query {
		t.Error("SetStatusIn: Expected method to return the same query instance")
	}
	if !query.HasStatusIn() {
		t.Error("HasStatusIn: Expected true after setting status in")
	}

	retrievedStatuses := query.StatusIn()
	if len(retrievedStatuses) != len(testStatuses) {
		t.Errorf("StatusIn: Expected %d statuses, got %d", len(testStatuses), len(retrievedStatuses))
	}
	for i, status := range testStatuses {
		if retrievedStatuses[i] != status {
			t.Errorf("StatusIn: Expected status '%s' at index %d, got '%s'", status, i, retrievedStatuses[i])
		}
	}
}

func TestTaskQuery_ChainedSetters(t *testing.T) {
	query := TaskDefinitionQuery()

	// Test that all setters can be chained
	result := query.
		SetAlias("test-alias").
		SetColumns([]string{"id", "alias"}).
		SetCountOnly(true).
		SetCreatedAtGte("2023-01-01 00:00:00").
		SetCreatedAtLte("2023-12-31 23:59:59").
		SetID("test-id").
		SetIDIn([]string{"id1", "id2"}).
		SetLimit(10).
		SetOffset(5).
		SetOrderBy("created_at").
		SetSoftDeletedIncluded(true).
		SetSortOrder("DESC").
		SetStatus("active").
		SetStatusIn([]string{"active", "canceled"})

	if result != query {
		t.Error("ChainedSetters: Expected all setters to return the same query instance for chaining")
	}

	// Verify all values were set correctly
	if query.Alias() != "test-alias" {
		t.Error("ChainedSetters: Alias not set correctly")
	}
	if len(query.Columns()) != 2 {
		t.Error("ChainedSetters: Columns not set correctly")
	}
	if !query.IsCountOnly() {
		t.Error("ChainedSetters: CountOnly not set correctly")
	}
	if query.CreatedAtGte() != "2023-01-01 00:00:00" {
		t.Error("ChainedSetters: CreatedAtGte not set correctly")
	}
	if query.CreatedAtLte() != "2023-12-31 23:59:59" {
		t.Error("ChainedSetters: CreatedAtLte not set correctly")
	}
	if query.ID() != "test-id" {
		t.Error("ChainedSetters: ID not set correctly")
	}
	if len(query.IDIn()) != 2 {
		t.Error("ChainedSetters: IDIn not set correctly")
	}
	if query.Limit() != 10 {
		t.Error("ChainedSetters: Limit not set correctly")
	}
	if query.Offset() != 5 {
		t.Error("ChainedSetters: Offset not set correctly")
	}
	if query.OrderBy() != "created_at" {
		t.Error("ChainedSetters: OrderBy not set correctly")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("ChainedSetters: SoftDeletedIncluded not set correctly")
	}
	if query.SortOrder() != "DESC" {
		t.Error("ChainedSetters: SortOrder not set correctly")
	}
	if query.Status() != "active" {
		t.Error("ChainedSetters: Status not set correctly")
	}
	if len(query.StatusIn()) != 2 {
		t.Error("ChainedSetters: StatusIn not set correctly")
	}
}
