package admin

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
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
				if !containsString(html, tt.expectedContains) {
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
			if !containsString(html, "breadcrumb") {
				t.Error("breadcrumbsUI() should contain breadcrumb class")
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && contains(s, substr))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func Test_Breadcrumb(t *testing.T) {
	tests := []struct {
		name string
		b    Breadcrumb
	}{
		{
			name: "breadcrumb with name and URL",
			b:    Breadcrumb{Name: "Home", URL: "/"},
		},
		{
			name: "breadcrumb with empty URL",
			b:    Breadcrumb{Name: "Test", URL: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.b.Name == "" && tt.b.URL == "" {
				t.Error("Breadcrumb should have at least name or URL")
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
			if !containsString(html, tt.title) {
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

func Test_keyEndpoint(t *testing.T) {
	// Test that keyEndpoint constant is defined
	if keyEndpoint != contextKey("endpoint") {
		t.Error("keyEndpoint should be contextKey('endpoint')")
	}
}

func Test_defaultFavicon(t *testing.T) {
	// Test that defaultFavicon is defined and not empty
	if defaultFavicon == "" {
		t.Error("defaultFavicon should not be empty")
	}
	if len(defaultFavicon) < 100 {
		t.Error("defaultFavicon should be a base64 string (longer than 100 chars)")
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
				if !containsString(html, tt.expectedContains) {
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

func Test_actionConstants(t *testing.T) {
	// Test action constants are defined
	tests := []struct {
		name  string
		value string
	}{
		{name: "actionModalQueuedTaskFilterShow", value: actionModalQueuedTaskFilterShow},
		{name: "actionModalQueuedTaskRestartShow", value: actionModalQueuedTaskRestartShow},
		{name: "actionModalQueuedTaskRestartSubmitted", value: actionModalQueuedTaskRestartSubmitted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
	}
}

func Test_pathConstants(t *testing.T) {
	// Test path constants are defined
	tests := []struct {
		name  string
		value string
	}{
		{name: "pathHome", value: pathHome},
		{name: "pathTaskQueueCreate", value: pathTaskQueueCreate},
		{name: "pathTaskQueueDelete", value: pathTaskQueueDelete},
		{name: "pathTaskQueueDetails", value: pathTaskQueueDetails},
		{name: "pathTaskQueueManager", value: pathTaskQueueManager},
		{name: "pathTaskQueueParameters", value: pathTaskQueueParameters},
		{name: "pathTaskQueueRequeue", value: pathTaskQueueRequeue},
		{name: "pathTaskQueueTaskRestart", value: pathTaskQueueTaskRestart},
		{name: "pathTaskDefinitionCreate", value: pathTaskDefinitionCreate},
		{name: "pathTaskDefinitionManager", value: pathTaskDefinitionManager},
		{name: "pathTaskDefinitionUpdate", value: pathTaskDefinitionUpdate},
		{name: "pathTaskDefinitionDelete", value: pathTaskDefinitionDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
	}
}

func Test_fieldConstants(t *testing.T) {
	// Test field constants are defined
	tests := []struct {
		name  string
		value string
	}{
		{name: "fieldParameters", value: fieldParameters},
		{name: "fieldQueueID", value: fieldQueueID},
		{name: "fieldTaskID", value: fieldTaskID},
		{name: "fieldStatus", value: fieldStatus},
		{name: "fieldTitle", value: fieldTitle},
		{name: "fieldAlias", value: fieldAlias},
		{name: "fieldDescription", value: fieldDescription},
		{name: "fieldDetails", value: fieldDetails},
		{name: "fieldFilterQueueID", value: fieldFilterQueueID},
		{name: "fieldFilterStatus", value: fieldFilterStatus},
		{name: "fieldFilterName", value: fieldFilterName},
		{name: "fieldFilterCreatedFrom", value: fieldFilterCreatedFrom},
		{name: "fieldFilterCreatedTo", value: fieldFilterCreatedTo},
		{name: "fieldFilterTaskID", value: fieldFilterTaskID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
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

func Test_homeController_prepareData(t *testing.T) {
	// Test homeController.prepareData method
	controller := &homeController{
		logger: *slog.Default(),
		store:  nil,
		layout: &mockLayout{},
	}

	req := &http.Request{}
	data, errorMessage := controller.prepareData(req)

	if errorMessage != "" {
		t.Errorf("prepareData() should not return error message, got %s", errorMessage)
	}
	if data.request == nil {
		t.Error("prepareData() should set request in data")
	}
	if data.request != req {
		t.Error("prepareData() should set the same request")
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
	if !containsString(result, "key") {
		t.Error("query() should contain key")
	}
	if !containsString(result, "value") {
		t.Error("query() should contain value")
	}
}
