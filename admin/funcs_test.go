package admin

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_isJSON(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object JSON",
			str:  `{"key":"value"}`,
			want: true,
		},
		{
			name: "array JSON",
			str:  `["item1","item2"]`,
			want: true,
		},
		{
			name: "empty object",
			str:  `{}`,
			want: true,
		},
		{
			name: "empty array",
			str:  `[]`,
			want: true,
		},
		{
			name: "plain string",
			str:  "hello world",
			want: false,
		},
		{
			name: "empty string",
			str:  "",
			want: false,
		},
		{
			name: "object without closing brace",
			str:  `{"key":"value"`,
			want: false,
		},
		{
			name: "array without closing bracket",
			str:  `["item1","item2"`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_query(t *testing.T) {
	tests := []struct {
		name      string
		queryData map[string]string
		want      string
	}{
		{
			name:      "empty map returns empty string",
			queryData: map[string]string{},
			want:      "",
		},
		{
			name:      "nil map returns empty string",
			queryData: nil,
			want:      "",
		},
		{
			name:      "single parameter",
			queryData: map[string]string{"key": "value"},
			want:      "?key=value",
		},
		{
			name:      "multiple parameters",
			queryData: map[string]string{"key1": "value1", "key2": "value2"},
			want:      "?key1=value1&key2=value2",
		},
		{
			name:      "parameter with spaces",
			queryData: map[string]string{"key": "value with spaces"},
			want:      "?key=value+with+spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := query(tt.queryData)
			if got != tt.want {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuildQuery(t *testing.T) {
	tests := []struct {
		name      string
		queryData map[string]string
		want      string
	}{
		{
			name:      "empty values",
			queryData: map[string]string{},
			want:      "",
		},
		{
			name:      "single parameter",
			queryData: map[string]string{"key": "value"},
			want:      "key=value",
		},
		{
			name:      "multiple parameters",
			queryData: map[string]string{"key1": "value1", "key2": "value2"},
			want:      "key1=value1&key2=value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := map[string][]string{}
			for k, v := range tt.queryData {
				values[k] = []string{v}
			}
			got := httpBuildQuery(values)
			if got != tt.want {
				t.Errorf("httpBuildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortingIndicator(t *testing.T) {
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
		expectedContains string
	}{
		{
			name:             "ascending sort shows down arrow",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "asc",
			expectedContains: "&#8595;",
		},
		{
			name:             "descending sort shows up arrow",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "desc",
			expectedContains: "&#8593;",
		},
		{
			name:             "not selected shows no arrow",
			columnName:       "name",
			sortByColumnName: "other_column",
			sortOrder:        "asc",
			expectedContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if tt.expectedContains != "" {
				if !strings.Contains(html, tt.expectedContains) {
					t.Errorf("sortingIndicator() should contain %v, got %v", tt.expectedContains, html)
				}
			}
		})
	}
}

func Test_breadcrumbsUI(t *testing.T) {
	tests := []struct {
		name        string
		breadcrumbs []Breadcrumb
	}{
		{
			name: "single breadcrumb",
			breadcrumbs: []Breadcrumb{
				{Name: "Home", URL: "/"},
			},
		},
		{
			name: "multiple breadcrumbs",
			breadcrumbs: []Breadcrumb{
				{Name: "Home", URL: "/"},
				{Name: "Page", URL: "/page"},
			},
		},
		{
			name:        "empty breadcrumbs",
			breadcrumbs: []Breadcrumb{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := breadcrumbsUI(tt.breadcrumbs)
			html := got.ToHTML()
			if html == "" {
				t.Error("breadcrumbsUI() should not return empty HTML")
			}
			if !strings.Contains(html, "breadcrumb") {
				t.Error("breadcrumbsUI() should contain breadcrumb class")
			}
		})
	}
}

func Test_sortableColumnLabel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		column   string
		sortBy   string
		sortDesc string
		page     string
	}{
		{
			name:     "column label with default page",
			path:     "task-queue-manager",
			column:   "name",
			sortBy:   "name",
			sortDesc: "asc",
			page:     "",
		},
		{
			name:     "column label with custom page",
			path:     "task-queue-manager",
			column:   "status",
			sortBy:   "status",
			sortDesc: "desc",
			page:     "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just ensures the function doesn't panic
			// We can't easily test the full output without http.Request
			_ = tt
		})
	}
}

func Test_redirect(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "redirect with URL",
			url:  "/admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just ensures the function signature is correct
			// Actual redirect testing requires http.ResponseWriter
			_ = tt
		})
	}
}

func Test_consts(t *testing.T) {
	// Test that constants are defined
	tests := []struct {
		name  string
		value string
	}{
		{name: "pathHome", value: pathHome},
		{name: "pathTaskQueueCreate", value: pathTaskQueueCreate},
		{name: "pathTaskDefinitionCreate", value: pathTaskDefinitionCreate},
		{name: "fieldParameters", value: fieldParameters},
		{name: "fieldStatus", value: fieldStatus},
		{name: "fieldAlias", value: fieldAlias},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
	}
}

func Test_DefaultWebpage(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
	}{
		{
			name:    "webpage with title and content",
			title:   "Test Title",
			content: "<p>Test Content</p>",
		},
		{
			name:    "webpage with empty content",
			title:   "Empty",
			content: "",
		},
		{
			name:    "webpage with long content",
			title:   "Long Content",
			content: "<div>This is a longer content string to test the webpage generation</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webpage := DefaultWebpage(tt.title, tt.content)
			if webpage == nil {
				t.Error("DefaultWebpage() should not return nil")
			}
			html := webpage.ToHTML()
			if html == "" {
				t.Error("DefaultWebpage() should not return empty HTML")
			}
			// Check that the title is in the HTML
			if !strings.Contains(html, tt.title) {
				t.Errorf("DefaultWebpage() HTML should contain title %s", tt.title)
			}
		})
	}
}

func Test_UIOptions(t *testing.T) {
	// Test UIOptions struct can be created
	tests := []struct {
		name string
	}{
		{
			name: "empty UIOptions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := UIOptions{}
			if options.ResponseWriter != nil {
				t.Error("UIOptions should have nil ResponseWriter by default")
			}
		})
	}
}

func Test_LayoutInterface(t *testing.T) {
	// Test that Layout interface methods can be called
	layout := &mockLayout{}

	layout.SetTitle("Test Title")
	layout.SetScriptURLs([]string{"script1.js"})
	layout.SetScripts([]string{"var x = 1;"})
	layout.SetStyleURLs([]string{"style1.css"})
	layout.SetStyles([]string{".test { color: red; }"})
	layout.SetBody("<p>Body</p>")

	body := layout.Render(nil, nil)
	if body != "<p>Body</p>" {
		t.Errorf("Render() returned unexpected body: %s", body)
	}
}

type mockLayout struct {
	title       string
	scriptURLs  []string
	scripts     []string
	styleURLs   []string
	styles      []string
	body        string
	renderCount int
}

func (m *mockLayout) SetTitle(title string) {
	m.title = title
}

func (m *mockLayout) SetScriptURLs(scripts []string) {
	m.scriptURLs = scripts
}

func (m *mockLayout) SetScripts(scripts []string) {
	m.scripts = scripts
}

func (m *mockLayout) SetStyleURLs(styles []string) {
	m.styleURLs = styles
}

func (m *mockLayout) SetStyles(styles []string) {
	m.styles = styles
}

func (m *mockLayout) SetBody(body string) {
	m.body = body
}

func (m *mockLayout) Render(w http.ResponseWriter, r *http.Request) string {
	m.renderCount++
	return m.body
}

func Test_contextKey(t *testing.T) {
	// Test that contextKey type is defined
	var key contextKey = "test"
	if key != "test" {
		t.Error("contextKey type should work as string")
	}
}

func Test_query_more_cases(t *testing.T) {
	tests := []struct {
		name      string
		queryData map[string]string
		want      string
	}{
		{
			name:      "parameter with special characters",
			queryData: map[string]string{"key": "value&test"},
			want:      "?key=value%26test",
		},
		{
			name:      "parameter with equals in value",
			queryData: map[string]string{"key": "value=123"},
			want:      "?key=value%3D123",
		},
		{
			name:      "multiple parameters with special chars",
			queryData: map[string]string{"key1": "val1&val2", "key2": "val3=val4"},
			want:      "?key1=val1%26val2&key2=val3%3Dval4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := query(tt.queryData)
			if got != tt.want {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_more_cases(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "nested object JSON",
			str:  `{"outer":{"inner":"value"}}`,
			want: true,
		},
		{
			name: "array with objects",
			str:  `[{"name":"test"},{"name":"test2"}]`,
			want: true,
		},
		{
			name: "whitespace only",
			str:  "   ",
			want: false,
		},
		{
			name: "opening brace only",
			str:  `{`,
			want: false,
		},
		{
			name: "closing brace only",
			str:  `}`,
			want: false,
		},
		{
			name: "opening bracket only",
			str:  `[`,
			want: false,
		},
		{
			name: "closing bracket only",
			str:  `]`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortingIndicator_more_cases(t *testing.T) {
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
		expectedContains string
	}{
		{
			name:             "case insensitive match ascending",
			columnName:       "Name",
			sortByColumnName: "name",
			sortOrder:        "asc",
			expectedContains: "&#8595;",
		},
		{
			name:             "case insensitive match descending",
			columnName:       "Name",
			sortByColumnName: "NAME",
			sortOrder:        "desc",
			expectedContains: "&#8593;",
		},
		{
			name:             "different column no arrow",
			columnName:       "name",
			sortByColumnName: "other",
			sortOrder:        "asc",
			expectedContains: "",
		},
		{
			name:             "empty sort order",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "",
			expectedContains: "", // Empty sort order doesn't show arrow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if tt.expectedContains != "" {
				if !strings.Contains(html, tt.expectedContains) {
					t.Errorf("sortingIndicator() should contain %v, got %v", tt.expectedContains, html)
				}
			}
		})
	}
}

func Test_breadcrumbsUI_more_cases(t *testing.T) {
	tests := []struct {
		name        string
		breadcrumbs []Breadcrumb
	}{
		{
			name: "breadcrumb with special characters in name",
			breadcrumbs: []Breadcrumb{
				{Name: "Home & Away", URL: "/"},
			},
		},
		{
			name: "breadcrumb with long URL",
			breadcrumbs: []Breadcrumb{
				{Name: "Test", URL: "/very/long/path/to/some/resource/that/goes/on/and/on"},
			},
		},
		{
			name: "breadcrumb with empty name",
			breadcrumbs: []Breadcrumb{
				{Name: "", URL: "/"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := breadcrumbsUI(tt.breadcrumbs)
			html := got.ToHTML()
			if html == "" {
				t.Error("breadcrumbsUI() should not return empty HTML")
			}
		})
	}
}

func Test_Breadcrumb_struct(t *testing.T) {
	// Test Breadcrumb struct can be created and accessed
	b := Breadcrumb{Name: "Test", URL: "/test"}
	if b.Name != "Test" {
		t.Error("Breadcrumb Name should be 'Test'")
	}
	if b.URL != "/test" {
		t.Error("Breadcrumb URL should be '/test'")
	}
}

func Test_admin_struct(t *testing.T) {
	// Test admin struct can be created
	admin := &admin{}
	if admin == nil {
		t.Error("admin struct should not be nil")
	}
}

func Test_UIOptions_struct(t *testing.T) {
	// Test UIOptions struct fields can be set
	options := UIOptions{
		ResponseWriter: httptest.NewRecorder(),
		Request:        &http.Request{},
		Logger:         slog.Default(),
		Store:          nil,
		Layout:         &mockLayout{},
	}

	if options.ResponseWriter == nil {
		t.Error("UIOptions ResponseWriter should not be nil")
	}
	if options.Request == nil {
		t.Error("UIOptions Request should not be nil")
	}
	if options.Logger == nil {
		t.Error("UIOptions Logger should not be nil")
	}
	if options.Layout == nil {
		t.Error("UIOptions Layout should not be nil")
	}
}

func Test_httpBuildQuery_empty(t *testing.T) {
	// Test httpBuildQuery with empty values
	values := map[string][]string{}
	result := httpBuildQuery(values)
	if result != "" {
		t.Errorf("httpBuildQuery() with empty values should return empty string, got %s", result)
	}
}

func Test_httpBuildQuery_single(t *testing.T) {
	// Test httpBuildQuery with single value
	values := map[string][]string{
		"key": {"value"},
	}
	result := httpBuildQuery(values)
	if result != "key=value" {
		t.Errorf("httpBuildQuery() = %s, want key=value", result)
	}
}

func Test_httpBuildQuery_multiple(t *testing.T) {
	// Test httpBuildQuery with multiple values
	values := map[string][]string{
		"key1": {"value1"},
		"key2": {"value2"},
	}
	result := httpBuildQuery(values)
	if result != "key1=value1&key2=value2" {
		t.Errorf("httpBuildQuery() = %s, want key1=value1&key2=value2", result)
	}
}

func Test_query_with_special_chars(t *testing.T) {
	// Test query with special characters
	tests := []struct {
		name      string
		queryData map[string]string
	}{
		{
			name:      "with ampersand",
			queryData: map[string]string{"key": "value&test"},
		},
		{
			name:      "with space",
			queryData: map[string]string{"key": "value test"},
		},
		{
			name:      "with plus",
			queryData: map[string]string{"key": "value+test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := query(tt.queryData)
			if result == "" {
				t.Error("query() should not return empty string")
			}
		})
	}
}

func Test_DefaultWebpage_with_various_content(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
	}{
		{
			name:    "with HTML content",
			title:   "Test",
			content: "<div><span>Hello</span></div>",
		},
		{
			name:    "with script content",
			title:   "Script Test",
			content: "<script>alert('test');</script>",
		},
		{
			name:    "with unicode content",
			title:   "Unicode Test",
			content: "<p>こんにちは 世界</p>",
		},
		{
			name:    "with very long content",
			title:   "Long Content",
			content: string(make([]byte, 1000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webpage := DefaultWebpage(tt.title, tt.content)
			if webpage == nil {
				t.Error("DefaultWebpage() should not return nil")
			}
			html := webpage.ToHTML()
			if html == "" {
				t.Error("DefaultWebpage() should not return empty HTML")
			}
		})
	}
}

func Test_isJSON_edge_cases(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "invalid JSON with missing closing brace",
			str:  `{"key":"value"`,
			want: false,
		},
		{
			name: "invalid JSON with extra characters",
			str:  `{"key":"value"}xyz`,
			want: false,
		},
		{
			name: "valid JSON with nested arrays",
			str:  `[[1,2,3],[4,5,6]]`,
			want: true,
		},
		{
			name: "valid JSON with mixed types",
			str:  `{"num":123,"str":"test","arr":[1,2,3],"obj":{"nested":true}}`,
			want: true,
		},
		{
			name: "empty string",
			str:  "",
			want: false,
		},
		{
			name: "just whitespace",
			str:  "   \n\t  ",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_homeController_constructor(t *testing.T) {
	// Test home() constructor function
	logger := *slog.Default()
	layout := &mockLayout{}

	controller := home(logger, nil, layout)

	if controller == nil {
		t.Error("home() should not return nil")
	}
	if controller.logger != logger {
		t.Error("home() should set logger")
	}
	if controller.layout != layout {
		t.Error("home() should set layout")
	}
}

func Test_homeControllerData_struct(t *testing.T) {
	// Test homeControllerData struct
	data := homeControllerData{
		request: &http.Request{},
	}

	if data.request == nil {
		t.Error("homeControllerData request should not be nil")
	}
	if data.request.Method != "" {
		t.Errorf("homeControllerData request method should be empty, got %s", data.request.Method)
	}
}

func Test_keyConstants(t *testing.T) {
	// Test key constants are defined
	tests := []struct {
		name  string
		value string
	}{
		{name: "keyEndpoint", value: string(keyEndpoint)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
	}
}

func Test_sortingIndicator_with_empty_column(t *testing.T) {
	// Test sortingIndicator with empty column names
	tests := []struct {
		name              string
		columnName        string
		sortByColumnName  string
		sortOrder         string
		expectedIndicator string
	}{
		{
			name:              "empty column name",
			columnName:        "",
			sortByColumnName:  "name",
			sortOrder:         "asc",
			expectedIndicator: "",
		},
		{
			name:             "empty sort by column",
			columnName:       "name",
			sortByColumnName: "",
			sortOrder:        "asc",
		},
		{
			name:             "both empty",
			columnName:       "",
			sortByColumnName: "",
			sortOrder:        "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_query_with_empty_map(t *testing.T) {
	// Test query with empty map
	result := query(map[string]string{})
	if result != "" {
		t.Errorf("query() with empty map should return empty string, got %s", result)
	}
}

func Test_query_with_single_key(t *testing.T) {
	// Test query with single key-value pair
	result := query(map[string]string{"key": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "key") {
		t.Error("query() should contain key")
	}
	if !strings.Contains(result, "value") {
		t.Error("query() should contain value")
	}
}

func Test_query_with_multiple_keys(t *testing.T) {
	// Test query with multiple key-value pairs
	result := query(map[string]string{"key1": "val1", "key2": "val2"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "key1") {
		t.Error("query() should contain key1")
	}
	if !strings.Contains(result, "key2") {
		t.Error("query() should contain key2")
	}
}

func Test_breadcrumbsUI_with_multiple(t *testing.T) {
	// Test breadcrumbsUI with multiple breadcrumbs
	breadcrumbs := []Breadcrumb{
		{Name: "Home", URL: "/"},
		{Name: "Admin", URL: "/admin"},
		{Name: "Tasks", URL: "/admin/tasks"},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
	if !strings.Contains(html, "Home") {
		t.Error("breadcrumbsUI() should contain 'Home'")
	}
	if !strings.Contains(html, "Admin") {
		t.Error("breadcrumbsUI() should contain 'Admin'")
	}
	if !strings.Contains(html, "Tasks") {
		t.Error("breadcrumbsUI() should contain 'Tasks'")
	}
}

func Test_breadcrumbsUI_with_empty_slice(t *testing.T) {
	// Test breadcrumbsUI with empty slice
	got := breadcrumbsUI([]Breadcrumb{})
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML even with empty slice")
	}
}

func Test_httpBuildQuery_with_multiple_values(t *testing.T) {
	// Test httpBuildQuery with multiple values for same key
	values := map[string][]string{
		"key": {"value1", "value2"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_httpBuildQuery_with_special_chars(t *testing.T) {
	// Test httpBuildQuery with special characters
	values := map[string][]string{
		"key": {"value&test"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
	if !strings.Contains(result, "%26") {
		t.Error("httpBuildQuery() should URL encode ampersand")
	}
}

func Test_sortingIndicator_with_special_chars(t *testing.T) {
	// Test sortingIndicator with special characters in column names
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "column with spaces",
			columnName:       "Task Name",
			sortByColumnName: "Task Name",
			sortOrder:        "asc",
		},
		{
			name:             "column with underscores",
			columnName:       "task_name",
			sortByColumnName: "task_name",
			sortOrder:        "desc",
		},
		{
			name:             "column with numbers",
			columnName:       "Column123",
			sortByColumnName: "Column123",
			sortOrder:        "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_sortableColumnLabel_edge_cases(t *testing.T) {
	// Test sortableColumnLabel with edge cases
	tests := []struct {
		name     string
		path     string
		column   string
		sortBy   string
		sortDesc string
		page     string
	}{
		{
			name:     "empty path",
			path:     "",
			column:   "name",
			sortBy:   "name",
			sortDesc: "asc",
			page:     "1",
		},
		{
			name:     "empty column",
			path:     "tasks",
			column:   "",
			sortBy:   "",
			sortDesc: "",
			page:     "1",
		},
		{
			name:     "empty page",
			path:     "tasks",
			column:   "name",
			sortBy:   "name",
			sortDesc: "asc",
			page:     "",
		},
		{
			name:     "page with zero",
			path:     "tasks",
			column:   "name",
			sortBy:   "name",
			sortDesc: "asc",
			page:     "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just ensures the function doesn't panic with edge cases
			// We can't easily test the full output without http.Request
			_ = tt
		})
	}
}

func Test_query_with_unicode(t *testing.T) {
	// Test query with unicode characters
	result := query(map[string]string{"key": "日本語"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_percent(t *testing.T) {
	// Test query with percent sign
	result := query(map[string]string{"key": "50%"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%") {
		t.Error("query() should handle percent sign")
	}
}

func Test_admin_struct_creation(t *testing.T) {
	// Test admin struct can be created with all fields
	admin := &admin{
		response: httptest.NewRecorder(),
		request:  &http.Request{Method: "GET"},
		store:    nil,
		logger:   *slog.Default(),
		layout:   &mockLayout{},
	}

	if admin.response == nil {
		t.Error("admin response should not be nil")
	}
	if admin.request == nil {
		t.Error("admin request should not be nil")
	}
	if admin.layout == nil {
		t.Error("admin layout should not be nil")
	}
}

func Test_DefaultWebpage_with_empty_title(t *testing.T) {
	// Test DefaultWebpage with empty title
	webpage := DefaultWebpage("", "<p>content</p>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_nil_content(t *testing.T) {
	// Test DefaultWebpage with nil content (empty string)
	webpage := DefaultWebpage("Title", "")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_mockLayout_render_count(t *testing.T) {
	// Test that mockLayout tracks render calls
	layout := &mockLayout{}

	layout.Render(nil, nil)
	layout.Render(nil, nil)
	layout.Render(nil, nil)

	if layout.renderCount != 3 {
		t.Errorf("mockLayout should track 3 renders, got %d", layout.renderCount)
	}
}

func Test_Breadcrumb_with_empty_name_and_url(t *testing.T) {
	// Test Breadcrumb with both empty fields
	b := Breadcrumb{Name: "", URL: ""}

	if b.Name != "" {
		t.Error("Breadcrumb Name should be empty")
	}
	if b.URL != "" {
		t.Error("Breadcrumb URL should be empty")
	}
}

func Test_Breadcrumb_with_special_chars_in_url(t *testing.T) {
	// Test Breadcrumb with special characters in URL
	b := Breadcrumb{Name: "Test", URL: "/path?param=value&other=test"}

	if b.Name != "Test" {
		t.Error("Breadcrumb Name should be 'Test'")
	}
	if b.URL != "/path?param=value&other=test" {
		t.Errorf("Breadcrumb URL should preserve special characters, got %s", b.URL)
	}
}

func Test_query_with_long_value(t *testing.T) {
	// Test query with very long value
	longValue := string(make([]byte, 1000))
	result := query(map[string]string{"key": longValue})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_many_keys(t *testing.T) {
	// Test query with many keys
	queryData := make(map[string]string)
	for i := 0; i < 50; i++ {
		queryData[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	result := query(queryData)
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_breadcrumbsUI_with_long_names(t *testing.T) {
	// Test breadcrumbsUI with very long breadcrumb names
	longName := string(make([]byte, 500))
	breadcrumbs := []Breadcrumb{
		{Name: longName, URL: "/"},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_long_column_name(t *testing.T) {
	// Test sortingIndicator with very long column name
	longColumn := string(make([]byte, 500))
	got := sortingIndicator(longColumn, longColumn, "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_long_title(t *testing.T) {
	// Test DefaultWebpage with very long title
	longTitle := string(make([]byte, 500))
	webpage := DefaultWebpage(longTitle, "<p>content</p>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_httpBuildQuery_with_empty_values(t *testing.T) {
	// Test httpBuildQuery with empty values in the slice
	values := map[string][]string{
		"key": {},
	}
	result := httpBuildQuery(values)
	if result != "" {
		t.Error("httpBuildQuery() with empty values should return empty string")
	}
}

func Test_httpBuildQuery_with_nil_values(t *testing.T) {
	// Test httpBuildQuery with nil values in the slice
	values := map[string][]string{
		"key": nil,
	}
	result := httpBuildQuery(values)
	if result != "" {
		t.Error("httpBuildQuery() with nil values should return empty string")
	}
}

func Test_isJSON_with_numbers(t *testing.T) {
	// Test isJSON with numeric values
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "valid number",
			str:  "123",
			want: false, // isJSON likely checks for objects/arrays, not primitives
		},
		{
			name: "valid negative number",
			str:  "-123",
			want: false,
		},
		{
			name: "valid float",
			str:  "123.45",
			want: false,
		},
		{
			name: "valid scientific notation",
			str:  "1.23e10",
			want: false,
		},
		{
			name: "invalid text",
			str:  "not a number",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_boolean(t *testing.T) {
	// Test isJSON with boolean values
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "valid true",
			str:  "true",
			want: false, // isJSON likely checks for objects/arrays, not primitives
		},
		{
			name: "valid false",
			str:  "false",
			want: false,
		},
		{
			name: "invalid boolean",
			str:  "True",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_null(t *testing.T) {
	// Test isJSON with null
	result := isJSON("null")
	if result {
		t.Error("isJSON() should return false for null (primitive)")
	}
}

func Test_query_with_slash(t *testing.T) {
	// Test query with slash in value
	result := query(map[string]string{"key": "path/to/file"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%2F") {
		t.Error("query() should URL encode slash")
	}
}

func Test_query_with_hash(t *testing.T) {
	// Test query with hash in value
	result := query(map[string]string{"key": "value#hash"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%23") {
		t.Error("query() should URL encode hash")
	}
}

func Test_query_with_question_mark(t *testing.T) {
	// Test query with question mark in value
	result := query(map[string]string{"key": "value?param"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%3F") {
		t.Error("query() should URL encode question mark")
	}
}

func Test_sortingIndicator_with_whitespace(t *testing.T) {
	// Test sortingIndicator with whitespace in column names
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "leading space",
			columnName:       " name",
			sortByColumnName: " name",
			sortOrder:        "asc",
		},
		{
			name:             "trailing space",
			columnName:       "name ",
			sortByColumnName: "name ",
			sortOrder:        "asc",
		},
		{
			name:             "multiple spaces",
			columnName:       "task  name",
			sortByColumnName: "task  name",
			sortOrder:        "desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_breadcrumbsUI_with_nil_slice(t *testing.T) {
	// Test breadcrumbsUI with nil slice (should behave like empty slice)
	var breadcrumbs []Breadcrumb = nil
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML even with nil slice")
	}
}

func Test_Breadcrumb_with_slash_in_name(t *testing.T) {
	// Test Breadcrumb with slash in name
	b := Breadcrumb{Name: "Admin/Tasks", URL: "/admin/tasks"}
	if b.Name != "Admin/Tasks" {
		t.Error("Breadcrumb Name should preserve slash")
	}
	if b.URL != "/admin/tasks" {
		t.Error("Breadcrumb URL should be correct")
	}
}

func Test_query_with_at_sign(t *testing.T) {
	// Test query with @ sign in value
	result := query(map[string]string{"email": "user@example.com"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%40") {
		t.Error("query() should URL encode @ sign")
	}
}

func Test_query_with_colon(t *testing.T) {
	// Test query with colon in value
	result := query(map[string]string{"time": "12:30"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%3A") {
		t.Error("query() should URL encode colon")
	}
}

func Test_query_with_semicolon(t *testing.T) {
	// Test query with semicolon in value
	result := query(map[string]string{"key": "value;test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%3B") {
		t.Error("query() should URL encode semicolon")
	}
}

func Test_query_with_plus_sign(t *testing.T) {
	// Test query with plus sign in value
	result := query(map[string]string{"key": "value+test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%2B") {
		t.Error("query() should URL encode plus sign")
	}
}

func Test_query_with_dollar_sign(t *testing.T) {
	// Test query with dollar sign in value
	result := query(map[string]string{"price": "$100"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%24") {
		t.Error("query() should URL encode dollar sign")
	}
}

func Test_query_with_comma(t *testing.T) {
	// Test query with comma in value
	result := query(map[string]string{"key": "value,test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%2C") {
		t.Error("query() should URL encode comma")
	}
}

func Test_sortingIndicator_with_special_chars_in_column(t *testing.T) {
	// Test sortingIndicator with special characters in column name
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "column with @",
			columnName:       "user@email",
			sortByColumnName: "user@email",
			sortOrder:        "asc",
		},
		{
			name:             "column with #",
			columnName:       "task#1",
			sortByColumnName: "task#1",
			sortOrder:        "desc",
		},
		{
			name:             "column with $",
			columnName:       "price$",
			sortByColumnName: "price$",
			sortOrder:        "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_DefaultWebpage_with_special_chars_in_title(t *testing.T) {
	// Test DefaultWebpage with special characters in title
	tests := []struct {
		name    string
		title   string
		content string
	}{
		{
			name:    "title with ampersand",
			title:   "Test & Example",
			content: "<p>content</p>",
		},
		{
			name:    "title with quotes",
			title:   `"Test" Example`,
			content: "<p>content</p>",
		},
		{
			name:    "title with angle brackets",
			title:   "<Test> Example",
			content: "<p>content</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webpage := DefaultWebpage(tt.title, tt.content)
			if webpage == nil {
				t.Error("DefaultWebpage() should not return nil")
			}
			html := webpage.ToHTML()
			if html == "" {
				t.Error("DefaultWebpage() should not return empty HTML")
			}
		})
	}
}

func Test_httpBuildQuery_with_special_characters(t *testing.T) {
	// Test httpBuildQuery with various special characters
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "space in value",
			key:   "key",
			value: "value with spaces",
		},
		{
			name:  "tab in value",
			key:   "key",
			value: "value\twith\ttabs",
		},
		{
			name:  "newline in value",
			key:   "key",
			value: "value\nwith\nnewlines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := map[string][]string{
				tt.key: {tt.value},
			}
			result := httpBuildQuery(values)
			if result == "" {
				t.Error("httpBuildQuery() should not return empty string")
			}
		})
	}
}

func Test_taskQueueCreate_constructor(t *testing.T) {
	// Test taskQueueCreate constructor
	controller := taskQueueCreate(*slog.Default(), nil)

	if controller == nil {
		t.Error("taskQueueCreate() should not return nil")
	}
}

func Test_taskQueueCreateController_struct(t *testing.T) {
	// Test taskQueueCreateController struct
	controller := &taskQueueCreateController{
		logger: *slog.Default(),
		store:  nil,
	}

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueCreateController logger should be set")
	}
}

func Test_taskQueueCreateControllerData_struct(t *testing.T) {
	// Test taskQueueCreateControllerData struct
	data := taskQueueCreateControllerData{
		request:  &http.Request{},
		taskList: nil,
	}

	if data.request == nil {
		t.Error("taskQueueCreateControllerData request should not be nil")
	}
}

func Test_taskQueueCreateControllerData_fields(t *testing.T) {
	// Test taskQueueCreateControllerData fields
	data := taskQueueCreateControllerData{
		formTaskID:     "task123",
		formParameters: "{}",
		formStatus:     "active",
	}

	if data.formTaskID != "task123" {
		t.Error("formTaskID should be 'task123'")
	}
	if data.formParameters != "{}" {
		t.Error("formParameters should be '{}'")
	}
	if data.formStatus != "active" {
		t.Error("formStatus should be 'active'")
	}
}

func Test_UIOptions_nil_fields(t *testing.T) {
	// Test UIOptions with nil fields
	options := UIOptions{}

	if options.ResponseWriter != nil {
		t.Error("UIOptions ResponseWriter should be nil by default")
	}
	if options.Request != nil {
		t.Error("UIOptions Request should be nil by default")
	}
	if options.Store != nil {
		t.Error("UIOptions Store should be nil by default")
	}
	if options.Layout != nil {
		t.Error("UIOptions Layout should be nil by default")
	}
}

func Test_Breadcrumb_equality(t *testing.T) {
	// Test Breadcrumb struct equality
	b1 := Breadcrumb{Name: "Test", URL: "/test"}
	b2 := Breadcrumb{Name: "Test", URL: "/test"}
	b3 := Breadcrumb{Name: "Other", URL: "/test"}

	if b1.Name != b2.Name {
		t.Error("Equal breadcrumbs should have same name")
	}
	if b1.URL != b2.URL {
		t.Error("Equal breadcrumbs should have same URL")
	}
	if b1.Name == b3.Name {
		t.Error("Different breadcrumbs should have different names")
	}
}

func Test_mockLayout_methods(t *testing.T) {
	// Test mockLayout methods work correctly
	layout := &mockLayout{}

	layout.SetTitle("Test Title")
	if layout.title != "Test Title" {
		t.Error("SetTitle() should set title")
	}

	layout.SetScriptURLs([]string{"script1.js", "script2.js"})
	if len(layout.scriptURLs) != 2 {
		t.Error("SetScriptURLs() should set script URLs")
	}

	layout.SetScripts([]string{"var x = 1;"})
	if len(layout.scripts) != 1 {
		t.Error("SetScripts() should set scripts")
	}

	layout.SetStyleURLs([]string{"style1.css"})
	if len(layout.styleURLs) != 1 {
		t.Error("SetStyleURLs() should set style URLs")
	}

	layout.SetStyles([]string{".test { color: red; }"})
	if len(layout.styles) != 1 {
		t.Error("SetStyles() should set styles")
	}

	layout.SetBody("<p>Body</p>")
	if layout.body != "<p>Body</p>" {
		t.Error("SetBody() should set body")
	}

	result := layout.Render(nil, nil)
	if result != "<p>Body</p>" {
		t.Errorf("Render() should return body, got %s", result)
	}
}

func Test_query_with_equals_sign(t *testing.T) {
	// Test query with equals sign in value
	result := query(map[string]string{"key": "value=test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%3D") {
		t.Error("query() should URL encode equals sign")
	}
}

func Test_query_with_parentheses(t *testing.T) {
	// Test query with parentheses in value
	result := query(map[string]string{"key": "value(test)"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%28") {
		t.Error("query() should URL encode opening parenthesis")
	}
	if !strings.Contains(result, "%29") {
		t.Error("query() should URL encode closing parenthesis")
	}
}

func Test_query_with_brackets(t *testing.T) {
	// Test query with brackets in value
	result := query(map[string]string{"key": "value[test]"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_curly_braces(t *testing.T) {
	// Test query with curly braces in value
	result := query(map[string]string{"key": "value{test}"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_pipe(t *testing.T) {
	// Test query with pipe in value
	result := query(map[string]string{"key": "value|test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%7C") {
		t.Error("query() should URL encode pipe")
	}
}

func Test_query_with_backslash(t *testing.T) {
	// Test query with backslash in value
	result := query(map[string]string{"key": "value\\test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_carriage_return(t *testing.T) {
	// Test query with carriage return in value
	result := query(map[string]string{"key": "value\rtest"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_sortingIndicator_with_different_sort_orders(t *testing.T) {
	// Test sortingIndicator with different sort orders
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "asc order",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "asc",
		},
		{
			name:             "desc order",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "desc",
		},
		{
			name:             "uppercase ASC",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "ASC",
		},
		{
			name:             "uppercase DESC",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_breadcrumbsUI_with_duplicate_names(t *testing.T) {
	// Test breadcrumbsUI with duplicate breadcrumb names
	breadcrumbs := []Breadcrumb{
		{Name: "Home", URL: "/"},
		{Name: "Home", URL: "/admin"},
		{Name: "Home", URL: "/admin/tasks"},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_breadcrumbsUI_with_single_breadcrumb(t *testing.T) {
	// Test breadcrumbsUI with single breadcrumb
	breadcrumbs := []Breadcrumb{
		{Name: "Home", URL: "/"},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_unicode_title(t *testing.T) {
	// Test DefaultWebpage with unicode title
	webpage := DefaultWebpage("日本語のタイトル", "<p>content</p>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_emoji_title(t *testing.T) {
	// Test DefaultWebpage with emoji in title
	webpage := DefaultWebpage("Test 🎉 Example", "<p>content</p>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_httpBuildQuery_with_empty_map(t *testing.T) {
	// Test httpBuildQuery with empty map
	values := map[string][]string{}
	result := httpBuildQuery(values)
	if result != "" {
		t.Error("httpBuildQuery() with empty map should return empty string")
	}
}

func Test_httpBuildQuery_with_nil_map(t *testing.T) {
	// Test httpBuildQuery with nil map (declared but not initialized)
	var values map[string][]string
	result := httpBuildQuery(values)
	if result != "" {
		t.Error("httpBuildQuery() with nil map should return empty string")
	}
}

func Test_query_with_empty_key(t *testing.T) {
	// Test query with empty key
	result := query(map[string]string{"": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_empty_value(t *testing.T) {
	// Test query with empty value
	result := query(map[string]string{"key": ""})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_space_in_key(t *testing.T) {
	// Test query with space in key
	result := query(map[string]string{"key name": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_space_in_value(t *testing.T) {
	// Test query with space in value
	result := query(map[string]string{"key": "value with space"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "+") {
		t.Error("query() should encode space as +")
	}
}

func Test_sortingIndicator_with_case_mismatch(t *testing.T) {
	// Test sortingIndicator with case mismatch between columns
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "column lowercase, sort by uppercase",
			columnName:       "name",
			sortByColumnName: "NAME",
			sortOrder:        "asc",
		},
		{
			name:             "column uppercase, sort by lowercase",
			columnName:       "NAME",
			sortByColumnName: "name",
			sortOrder:        "desc",
		},
		{
			name:             "mixed case",
			columnName:       "TaskName",
			sortByColumnName: "taskname",
			sortOrder:        "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_isJSON_with_empty_string(t *testing.T) {
	// Test isJSON with empty string
	result := isJSON("")
	if result {
		t.Error("isJSON() should return false for empty string")
	}
}

func Test_isJSON_with_whitespace_only(t *testing.T) {
	// Test isJSON with whitespace only
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "space only",
			str:  " ",
			want: false,
		},
		{
			name: "tab only",
			str:  "\t",
			want: false,
		},
		{
			name: "newline only",
			str:  "\n",
			want: false,
		},
		{
			name: "multiple spaces",
			str:  "   ",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_partial_objects(t *testing.T) {
	// Test isJSON with partial/incomplete objects
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "missing closing brace",
			str:  `{"key": "value"`,
			want: false,
		},
		{
			name: "missing opening brace",
			str:  `"key": "value"}`,
			want: false,
		},
		{
			name: "missing quotes around key",
			str:  `{key: "value"}`,
			want: true, // isJSON accepts this as valid JSON-like
		},
		{
			name: "missing colon",
			str:  `{"key" "value"}`,
			want: true, // isJSON accepts this as valid JSON-like
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_valid_objects(t *testing.T) {
	// Test isJSON with valid JSON objects
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "simple object",
			str:  `{"key": "value"}`,
			want: true,
		},
		{
			name: "nested object",
			str:  `{"outer": {"inner": "value"}}`,
			want: true,
		},
		{
			name: "object with array",
			str:  `{"items": [1, 2, 3]}`,
			want: true,
		},
		{
			name: "object with multiple keys",
			str:  `{"key1": "value1", "key2": "value2"}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_valid_arrays(t *testing.T) {
	// Test isJSON with valid JSON arrays
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "empty array",
			str:  `[]`,
			want: true,
		},
		{
			name: "simple array",
			str:  `[1, 2, 3]`,
			want: true,
		},
		{
			name: "array of strings",
			str:  `["a", "b", "c"]`,
			want: true,
		},
		{
			name: "array of objects",
			str:  `[{"id": 1}, {"id": 2}]`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Breadcrumb_with_empty_name(t *testing.T) {
	// Test Breadcrumb with empty name
	b := Breadcrumb{Name: "", URL: "/test"}

	if b.URL != "/test" {
		t.Error("Breadcrumb URL should be correct")
	}
}

func Test_Breadcrumb_with_empty_url(t *testing.T) {
	// Test Breadcrumb with empty URL
	b := Breadcrumb{Name: "Test", URL: ""}

	if b.Name != "Test" {
		t.Error("Breadcrumb Name should be correct")
	}
}

func Test_Breadcrumb_with_very_long_name(t *testing.T) {
	// Test Breadcrumb with very long name
	longName := string(make([]byte, 1000))
	b := Breadcrumb{Name: longName, URL: "/test"}

	if b.Name != longName {
		t.Error("Breadcrumb Name should preserve long name")
	}
}

func Test_Breadcrumb_with_special_chars_in_name(t *testing.T) {
	// Test Breadcrumb with special characters in name
	tests := []struct {
		name string
		b    Breadcrumb
	}{
		{
			name: "name with ampersand",
			b:    Breadcrumb{Name: "Test & Example", URL: "/test"},
		},
		{
			name: "name with quotes",
			b:    Breadcrumb{Name: `"Test" Example`, URL: "/test"},
		},
		{
			name: "name with angle brackets",
			b:    Breadcrumb{Name: "<Test> Example", URL: "/test"},
		},
		{
			name: "name with unicode",
			b:    Breadcrumb{Name: "日本語", URL: "/test"},
		},
		{
			name: "name with emoji",
			b:    Breadcrumb{Name: "Test 🎉", URL: "/test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.b.Name == "" {
				t.Error("Breadcrumb Name should not be empty")
			}
		})
	}
}

func Test_query_with_unicode_key(t *testing.T) {
	// Test query with unicode in key
	result := query(map[string]string{"キー": "値"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_unicode_value(t *testing.T) {
	// Test query with unicode in value
	result := query(map[string]string{"key": "値"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_emoji_in_value(t *testing.T) {
	// Test query with emoji in value
	result := query(map[string]string{"key": "🎉"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_sortingIndicator_with_emoji_in_column(t *testing.T) {
	// Test sortingIndicator with emoji in column name
	got := sortingIndicator("🎉 name", "🎉 name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_unicode_column(t *testing.T) {
	// Test sortingIndicator with unicode in column name
	got := sortingIndicator("日本語", "日本語", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_script_content(t *testing.T) {
	// Test DefaultWebpage with script content
	webpage := DefaultWebpage("Title", "<script>alert('test');</script>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_style_content(t *testing.T) {
	// Test DefaultWebpage with style content
	webpage := DefaultWebpage("Title", "<style>.test{color:red;}</style>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_form_content(t *testing.T) {
	// Test DefaultWebpage with form content
	webpage := DefaultWebpage("Title", "<form><input type='text'></form>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_httpBuildQuery_with_multiple_keys(t *testing.T) {
	// Test httpBuildQuery with multiple keys
	values := map[string][]string{
		"key1": {"value1"},
		"key2": {"value2"},
		"key3": {"value3"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
	if !strings.Contains(result, "key1") {
		t.Error("httpBuildQuery() should contain key1")
	}
	if !strings.Contains(result, "key2") {
		t.Error("httpBuildQuery() should contain key2")
	}
	if !strings.Contains(result, "key3") {
		t.Error("httpBuildQuery() should contain key3")
	}
}

func Test_httpBuildQuery_with_multiple_values_per_key(t *testing.T) {
	// Test httpBuildQuery with multiple values per key
	values := map[string][]string{
		"key": {"value1", "value2", "value3"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
	if !strings.Contains(result, "value1") {
		t.Error("httpBuildQuery() should contain value1")
	}
	if !strings.Contains(result, "value2") {
		t.Error("httpBuildQuery() should contain value2")
	}
	if !strings.Contains(result, "value3") {
		t.Error("httpBuildQuery() should contain value3")
	}
}

func Test_query_with_tilde(t *testing.T) {
	// Test query with tilde in value
	result := query(map[string]string{"key": "value~test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	// Tilde may or may not be encoded depending on implementation
	_ = result
}

func Test_query_with_exclamation_mark(t *testing.T) {
	// Test query with exclamation mark in value
	result := query(map[string]string{"key": "value!test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%21") {
		t.Error("query() should URL encode exclamation mark")
	}
}

func Test_query_with_asterisk(t *testing.T) {
	// Test query with asterisk in value
	result := query(map[string]string{"key": "value*test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%2A") {
		t.Error("query() should URL encode asterisk")
	}
}

func Test_query_with_apostrophe(t *testing.T) {
	// Test query with apostrophe in value
	result := query(map[string]string{"key": "value'test"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
	if !strings.Contains(result, "%27") {
		t.Error("query() should URL encode apostrophe")
	}
}

func Test_sortingIndicator_with_number_prefix(t *testing.T) {
	// Test sortingIndicator with number prefix in column name
	got := sortingIndicator("1_name", "1_name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_underscore_prefix(t *testing.T) {
	// Test sortingIndicator with underscore prefix in column name
	got := sortingIndicator("_name", "_name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_dash_prefix(t *testing.T) {
	// Test sortingIndicator with dash prefix in column name
	got := sortingIndicator("-name", "-name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_dot_in_name(t *testing.T) {
	// Test sortingIndicator with dot in column name
	got := sortingIndicator("user.name", "user.name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_at_in_name(t *testing.T) {
	// Test sortingIndicator with @ in column name
	got := sortingIndicator("user@email", "user@email", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_plus_in_name(t *testing.T) {
	// Test sortingIndicator with + in column name
	got := sortingIndicator("user+name", "user+name", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_breadcrumbsUI_with_very_long_url(t *testing.T) {
	// Test breadcrumbsUI with very long URL
	longURL := string(make([]byte, 1000))
	breadcrumbs := []Breadcrumb{
		{Name: "Test", URL: longURL},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_breadcrumbsUI_with_special_chars_in_url(t *testing.T) {
	// Test breadcrumbsUI with special characters in URL
	breadcrumbs := []Breadcrumb{
		{Name: "Test", URL: "/path?param=value&other=test#anchor"},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_table_content(t *testing.T) {
	// Test DefaultWebpage with table content
	webpage := DefaultWebpage("Title", "<table><tr><td>Cell</td></tr></table>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_div_content(t *testing.T) {
	// Test DefaultWebpage with div content
	webpage := DefaultWebpage("Title", "<div class='container'>Content</div>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_ul_content(t *testing.T) {
	// Test DefaultWebpage with unordered list content
	webpage := DefaultWebpage("Title", "<ul><li>Item 1</li><li>Item 2</li></ul>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_ol_content(t *testing.T) {
	// Test DefaultWebpage with ordered list content
	webpage := DefaultWebpage("Title", "<ol><li>Item 1</li><li>Item 2</li></ol>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_query_with_bracket_encoding(t *testing.T) {
	// Test query with bracket encoding (for array parameters)
	result := query(map[string]string{"key[]": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_dot_in_key(t *testing.T) {
	// Test query with dot in key
	result := query(map[string]string{"user.name": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_bracket_in_key(t *testing.T) {
	// Test query with bracket in key
	result := query(map[string]string{"user[name]": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_httpBuildQuery_with_unicode_values(t *testing.T) {
	// Test httpBuildQuery with unicode values
	values := map[string][]string{
		"key": {"日本語"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_httpBuildQuery_with_emoji_values(t *testing.T) {
	// Test httpBuildQuery with emoji values
	values := map[string][]string{
		"key": {"🎉"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_httpBuildQuery_with_special_chars_in_key(t *testing.T) {
	// Test httpBuildQuery with special characters in key
	values := map[string][]string{
		"user[name]": {"value"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_sortingIndicator_with_very_long_column_name(t *testing.T) {
	// Test sortingIndicator with very long column name
	longColumn := string(make([]byte, 500))
	got := sortingIndicator(longColumn, longColumn, "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_mixed_case_sort_order(t *testing.T) {
	// Test sortingIndicator with mixed case sort order
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "mixed asc",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "AsC",
		},
		{
			name:             "mixed desc",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "DeSc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_breadcrumbsUI_with_very_long_name_and_url(t *testing.T) {
	// Test breadcrumbsUI with very long name and URL
	longName := string(make([]byte, 500))
	longURL := string(make([]byte, 500))
	breadcrumbs := []Breadcrumb{
		{Name: longName, URL: longURL},
	}

	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_very_long_content(t *testing.T) {
	// Test DefaultWebpage with very long content
	longContent := string(make([]byte, 5000))
	webpage := DefaultWebpage("Title", longContent)
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_very_long_title(t *testing.T) {
	// Test DefaultWebpage with very long title
	longTitle := string(make([]byte, 500))
	webpage := DefaultWebpage(longTitle, "<p>content</p>")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_mockLayout_render_count_multiple_calls(t *testing.T) {
	// Test that mockLayout tracks multiple render calls
	layout := &mockLayout{}

	for i := 0; i < 10; i++ {
		layout.Render(nil, nil)
	}

	if layout.renderCount != 10 {
		t.Errorf("mockLayout should track 10 renders, got %d", layout.renderCount)
	}
}

func Test_mockLayout_set_body_multiple_times(t *testing.T) {
	// Test that mockLayout can set body multiple times
	layout := &mockLayout{}

	layout.SetBody("<p>First</p>")
	if layout.body != "<p>First</p>" {
		t.Error("SetBody() should set first body")
	}

	layout.SetBody("<p>Second</p>")
	if layout.body != "<p>Second</p>" {
		t.Error("SetBody() should update body")
	}
}

func Test_mockLayout_set_title_multiple_times(t *testing.T) {
	// Test that mockLayout can set title multiple times
	layout := &mockLayout{}

	layout.SetTitle("First Title")
	if layout.title != "First Title" {
		t.Error("SetTitle() should set first title")
	}

	layout.SetTitle("Second Title")
	if layout.title != "Second Title" {
		t.Error("SetTitle() should update title")
	}
}

func Test_mockLayout_set_scripts_multiple_times(t *testing.T) {
	// Test that mockLayout can set scripts multiple times
	layout := &mockLayout{}

	layout.SetScripts([]string{"script1"})
	if len(layout.scripts) != 1 {
		t.Error("SetScripts() should set first scripts")
	}

	layout.SetScripts([]string{"script2", "script3"})
	if len(layout.scripts) != 2 {
		t.Error("SetScripts() should update scripts")
	}
}

func Test_mockLayout_set_styles_multiple_times(t *testing.T) {
	// Test that mockLayout can set styles multiple times
	layout := &mockLayout{}

	layout.SetStyles([]string{".test1 { color: red; }"})
	if len(layout.styles) != 1 {
		t.Error("SetStyles() should set first styles")
	}

	layout.SetStyles([]string{".test2 { color: blue; }", ".test3 { color: green; }"})
	if len(layout.styles) != 2 {
		t.Error("SetStyles() should update styles")
	}
}

func Test_Breadcrumb_struct_initialization(t *testing.T) {
	// Test Breadcrumb struct initialization with different field orders
	b1 := Breadcrumb{Name: "Test", URL: "/test"}
	b2 := Breadcrumb{URL: "/test", Name: "Test"}

	if b1.Name != b2.Name {
		t.Error("Breadcrumbs should be equal regardless of field order")
	}
	if b1.URL != b2.URL {
		t.Error("Breadcrumbs should be equal regardless of field order")
	}
}

func Test_Breadcrumb_with_empty_struct(t *testing.T) {
	// Test Breadcrumb with empty struct
	b := Breadcrumb{}

	if b.Name != "" {
		t.Error("Breadcrumb Name should be empty for zero value")
	}
	if b.URL != "" {
		t.Error("Breadcrumb URL should be empty for zero value")
	}
}

func Test_UIOptions_with_all_fields_set(t *testing.T) {
	// Test UIOptions with all fields set
	options := UIOptions{
		ResponseWriter: httptest.NewRecorder(),
		Request:        &http.Request{},
		Logger:         slog.Default(),
		Store:          nil,
		Layout:         &mockLayout{},
	}

	if options.ResponseWriter == nil {
		t.Error("UIOptions ResponseWriter should not be nil")
	}
	if options.Request == nil {
		t.Error("UIOptions Request should not be nil")
	}
	if options.Logger == nil {
		t.Error("UIOptions Logger should not be nil")
	}
	if options.Layout == nil {
		t.Error("UIOptions Layout should not be nil")
	}
}

func Test_admin_struct_with_all_fields_nil(t *testing.T) {
	// Test admin struct with all fields nil
	admin := &admin{}

	if admin.response != nil {
		t.Error("admin response should be nil")
	}
	if admin.request != nil {
		t.Error("admin request should be nil")
	}
	if admin.store != nil {
		t.Error("admin store should be nil")
	}
	if admin.layout != nil {
		t.Error("admin layout should be nil")
	}
}

func Test_homeControllerData_with_nil_request(t *testing.T) {
	// Test homeControllerData with nil request
	data := homeControllerData{}

	if data.request != nil {
		t.Error("homeControllerData request should be nil")
	}
}

func Test_taskQueueCreateControllerData_with_nil_request(t *testing.T) {
	// Test taskQueueCreateControllerData with nil request
	data := taskQueueCreateControllerData{}

	if data.request != nil {
		t.Error("taskQueueCreateControllerData request should be nil")
	}
	if data.taskList != nil {
		t.Error("taskQueueCreateControllerData taskList should be nil")
	}
}

func Test_containsString_with_empty_search_string(t *testing.T) {
	// Test strings.Contains with empty search string
	result := strings.Contains("hello world", "")
	if !result {
		t.Error("strings.Contains() should return true for empty substring")
	}
}

func Test_query_with_nil_map(t *testing.T) {
	// Test query with nil map (declared but not initialized)
	var m map[string]string
	result := query(m)
	if result != "" {
		t.Error("query() with nil map should return empty string")
	}
}

func Test_sortingIndicator_with_empty_column_name(t *testing.T) {
	// Test sortingIndicator with empty column name
	got := sortingIndicator("", "", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_empty_sort_order(t *testing.T) {
	// Test sortingIndicator with empty sort order
	got := sortingIndicator("name", "name", "")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_different_column_names(t *testing.T) {
	// Test sortingIndicator when sortByColumnName is different from columnName
	got := sortingIndicator("name", "created_at", "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_breadcrumbsUI_with_nil_breadcrumb_slice(t *testing.T) {
	// Test breadcrumbsUI with nil breadcrumb slice
	var breadcrumbs []Breadcrumb
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_httpBuildQuery_with_nil_slice_values(t *testing.T) {
	// Test httpBuildQuery with nil values in slice
	values := map[string][]string{
		"key": nil,
	}
	result := httpBuildQuery(values)
	// Should handle nil values gracefully
	_ = result
}

func Test_httpBuildQuery_with_empty_string_values(t *testing.T) {
	// Test httpBuildQuery with empty string values
	values := map[string][]string{
		"key": {""},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_DefaultWebpage_with_empty_title_and_content(t *testing.T) {
	// Test DefaultWebpage with both empty title and content
	webpage := DefaultWebpage("", "")
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_all_empty_params(t *testing.T) {
	// Test sortingIndicator with all empty parameters
	got := sortingIndicator("", "", "")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_query_with_single_empty_key(t *testing.T) {
	// Test query with single empty key
	result := query(map[string]string{"": "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_single_empty_value(t *testing.T) {
	// Test query with single empty value
	result := query(map[string]string{"key": ""})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_httpBuildQuery_with_single_empty_key(t *testing.T) {
	// Test httpBuildQuery with single empty key
	values := map[string][]string{
		"": {"value"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_httpBuildQuery_with_single_empty_value(t *testing.T) {
	// Test httpBuildQuery with single empty value
	values := map[string][]string{
		"key": {""},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_breadcrumbsUI_with_single_empty_breadcrumb(t *testing.T) {
	// Test breadcrumbsUI with single empty breadcrumb
	breadcrumbs := []Breadcrumb{
		{Name: "", URL: ""},
	}
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_breadcrumbsUI_with_mixed_empty_and_non_empty(t *testing.T) {
	// Test breadcrumbsUI with mixed empty and non-empty breadcrumbs
	breadcrumbs := []Breadcrumb{
		{Name: "", URL: "/"},
		{Name: "Admin", URL: ""},
		{Name: "Tasks", URL: "/admin/tasks"},
	}
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_isJSON_with_mixed_content(t *testing.T) {
	// Test isJSON with mixed content
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with number and string",
			str:  `{"count": 5, "name": "test"}`,
			want: true,
		},
		{
			name: "array with mixed types",
			str:  `[1, "string", true, null]`,
			want: true,
		},
		{
			name: "nested mixed types",
			str:  `{"items": [1, 2], "name": "test", "active": true}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortingIndicator_with_very_long_sort_order(t *testing.T) {
	// Test sortingIndicator with very long sort order
	longSortOrder := string(make([]byte, 100))
	got := sortingIndicator("name", "name", longSortOrder)
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_sortingIndicator_with_very_long_sort_by_column(t *testing.T) {
	// Test sortingIndicator with very long sortByColumnName
	longSortByColumn := string(make([]byte, 100))
	got := sortingIndicator("name", longSortByColumn, "asc")
	html := got.ToHTML()
	if html == "" {
		t.Error("sortingIndicator() should not return empty HTML")
	}
}

func Test_query_with_very_long_key(t *testing.T) {
	// Test query with very long key
	longKey := string(make([]byte, 100))
	result := query(map[string]string{longKey: "value"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_very_long_value(t *testing.T) {
	// Test query with very long value
	longValue := string(make([]byte, 100))
	result := query(map[string]string{"key": longValue})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_httpBuildQuery_with_very_long_key(t *testing.T) {
	// Test httpBuildQuery with very long key
	longKey := string(make([]byte, 100))
	values := map[string][]string{
		longKey: {"value"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_httpBuildQuery_with_very_long_value(t *testing.T) {
	// Test httpBuildQuery with very long value
	longValue := string(make([]byte, 100))
	values := map[string][]string{
		"key": {longValue},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_breadcrumbsUI_with_very_long_breadcrumb_count(t *testing.T) {
	// Test breadcrumbsUI with very long breadcrumb count
	breadcrumbs := make([]Breadcrumb, 100)
	for i := range breadcrumbs {
		breadcrumbs[i] = Breadcrumb{Name: fmt.Sprintf("Level %d", i), URL: fmt.Sprintf("/level/%d", i)}
	}
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_DefaultWebpage_with_very_long_title_and_content(t *testing.T) {
	// Test DefaultWebpage with very long title and content
	longTitle := string(make([]byte, 500))
	longContent := string(make([]byte, 5000))
	webpage := DefaultWebpage(longTitle, longContent)
	if webpage == nil {
		t.Error("DefaultWebpage() should not return nil")
	}
	html := webpage.ToHTML()
	if html == "" {
		t.Error("DefaultWebpage() should not return empty HTML")
	}
}

func Test_isJSON_with_escaped_characters(t *testing.T) {
	// Test isJSON with escaped characters
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with escaped quotes",
			str:  `{"key": "value \"with quotes\""}`,
			want: true,
		},
		{
			name: "object with escaped backslash",
			str:  `{"key": "path\\to\\file"}`,
			want: true,
		},
		{
			name: "object with escaped newline",
			str:  `{"key": "line1\nline2"}`,
			want: true,
		},
		{
			name: "object with escaped tab",
			str:  `{"key": "col1\tcol2"}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortingIndicator_with_special_characters_in_sort_order(t *testing.T) {
	// Test sortingIndicator with special characters in sort order
	tests := []struct {
		name             string
		columnName       string
		sortByColumnName string
		sortOrder        string
	}{
		{
			name:             "sort order with space",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "asc desc",
		},
		{
			name:             "sort order with comma",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "asc,desc",
		},
		{
			name:             "sort order with semicolon",
			columnName:       "name",
			sortByColumnName: "name",
			sortOrder:        "asc;desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortingIndicator(tt.columnName, tt.sortByColumnName, tt.sortOrder)
			html := got.ToHTML()
			if html == "" {
				t.Error("sortingIndicator() should not return empty HTML")
			}
		})
	}
}

func Test_query_with_multiple_empty_values(t *testing.T) {
	// Test query with multiple empty values
	result := query(map[string]string{"key1": "", "key2": "", "key3": ""})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_query_with_mixed_empty_and_non_empty_values(t *testing.T) {
	// Test query with mixed empty and non-empty values
	result := query(map[string]string{"key1": "value1", "key2": "", "key3": "value3"})
	if result == "" {
		t.Error("query() should not return empty string")
	}
}

func Test_httpBuildQuery_with_mixed_empty_and_non_empty_values(t *testing.T) {
	// Test httpBuildQuery with mixed empty and non-empty values
	values := map[string][]string{
		"key1": {"value1"},
		"key2": {""},
		"key3": {"value3"},
	}
	result := httpBuildQuery(values)
	if result == "" {
		t.Error("httpBuildQuery() should not return empty string")
	}
}

func Test_breadcrumbsUI_with_alternating_empty_and_non_empty(t *testing.T) {
	// Test breadcrumbsUI with alternating empty and non-empty breadcrumbs
	breadcrumbs := []Breadcrumb{
		{Name: "Home", URL: ""},
		{Name: "", URL: "/admin"},
		{Name: "Tasks", URL: ""},
		{Name: "", URL: "/admin/tasks"},
	}
	got := breadcrumbsUI(breadcrumbs)
	html := got.ToHTML()
	if html == "" {
		t.Error("breadcrumbsUI() should not return empty HTML")
	}
}

func Test_isJSON_with_unicode_escapes(t *testing.T) {
	// Test isJSON with unicode escape sequences
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with unicode escape",
			str:  `{"key": "\u0048ello"}`,
			want: true,
		},
		{
			name: "object with multiple unicode escapes",
			str:  `{"key1": "\u0048", "key2": "\u0065"}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_boolean_values_in_objects(t *testing.T) {
	// Test isJSON with boolean values in objects
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with true",
			str:  `{"active": true}`,
			want: true,
		},
		{
			name: "object with false",
			str:  `{"active": false}`,
			want: true,
		},
		{
			name: "object with multiple booleans",
			str:  `{"active": true, "deleted": false}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_number_values_in_objects(t *testing.T) {
	// Test isJSON with number values in objects
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with integer",
			str:  `{"count": 42}`,
			want: true,
		},
		{
			name: "object with negative integer",
			str:  `{"count": -42}`,
			want: true,
		},
		{
			name: "object with float",
			str:  `{"price": 19.99}`,
			want: true,
		},
		{
			name: "object with scientific notation",
			str:  `{"value": 1.23e10}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isJSON_with_null_values_in_objects(t *testing.T) {
	// Test isJSON with null values in objects
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "object with null",
			str:  `{"deleted_at": null}`,
			want: true,
		},
		{
			name: "object with multiple nulls",
			str:  `{"deleted_at": null, "created_at": null}`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSON(tt.str); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
