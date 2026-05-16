package admin

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	urlpkg "net/url"

	"github.com/dracory/hb"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func adminHeader(store taskstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(url(r, pathHome, nil)).
		Class("nav-link")
	linkQueue := hb.Hyperlink().
		HTML("Task Queue").
		Href(url(r, pathTaskQueueManager, nil)).
		Class("nav-link")
	linkTasks := hb.Hyperlink().
		HTML("Task Definitions").
		Href(url(r, pathTaskDefinitionManager, nil)).
		Class("nav-link")

	queueCount, err := store.TaskQueueCount(context.Background(), taskstore.TaskQueueQuery())

	if err != nil {
		logger.Error(err.Error())
		queueCount = -1
	}

	taskCount, err := store.TaskDefinitionCount(context.Background(), taskstore.TaskDefinitionQuery())

	if err != nil {
		logger.Error(err.Error())
		taskCount = -1
	}

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkQueue.
			Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(queueCount)))))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkTasks.
			Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(taskCount)))))

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}

func breadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
	adminHomeURL := "/admin"

	adminHomeBreadcrumb := lo.
		If(adminHomeURL != "", Breadcrumb{
			Name: "Home",
			URL:  adminHomeURL,
		}).
		Else(Breadcrumb{})

	breadcrumbItems := []Breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "Zeppelin",
			URL:  url(r, pathHome, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	breadcrumbs := breadcrumbsUI(breadcrumbItems)

	return hb.Div().
		Child(breadcrumbs)
}

func redirect(w http.ResponseWriter, r *http.Request, url string) string {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return ""
}

func url(r *http.Request, path string, params map[string]string) string {
	endpoint := r.URL.Path

	if params == nil {
		params = map[string]string{}
	}

	params["controller"] = path

	url := endpoint + query(params)

	return url
}

// linkWithContext provides a context-based URL generation alternative
// This uses context to retrieve the endpoint, useful in complex routing scenarios
// Usage: Set endpoint in context with keyEndpoint before calling this function
func linkWithContext(r *http.Request, path string, params map[string]string) string {
	endpoint := r.URL.Path

	// Try to get endpoint from context first (if set by caller)
	if ctxEndpoint := r.Context().Value(keyEndpoint); ctxEndpoint != nil {
		if endpointStr, ok := ctxEndpoint.(string); ok {
			endpoint = endpointStr
		}
	}

	if params == nil {
		params = map[string]string{}
	}

	params["controller"] = path

	return endpoint + query(params)
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := urlpkg.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}

type Breadcrumb struct {
	Name string
	URL  string
}

func breadcrumbsUI(breadcrumbs []Breadcrumb) hb.TagInterface {

	ol := hb.OL().
		Class("breadcrumb").
		Style("margin-bottom: 0px;")

	for _, breadcrumb := range breadcrumbs {

		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(link)

		ol.AddChild(li)
	}

	nav := hb.Nav().
		Class("d-inline-block").
		Attr("aria-label", "breadcrumb").
		Child(ol)

	return nav
}

// isJSON is naive implementation for superficial, rough and fast checking for JSON
func isJSON(str string) bool {
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return true
	}

	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		return true
	}

	return false
}

// sortableColumnLabel creates a sortable column label with sorting indicator
// This is a standalone utility function for better reusability across controllers
// page parameter defaults to "0" to reset pagination on sort change
func sortableColumnLabel(r *http.Request, tableLabel, columnName, path string, sortBy, sortOrder, page string) hb.TagInterface {
	if page == "" {
		page = "0"
	}

	isSelected := strings.EqualFold(sortBy, columnName)

	direction := lo.If(sortOrder == "asc", "desc").Else("asc")

	if !isSelected {
		direction = "asc"
	}

	linkURL := url(r, path, map[string]string{
		"page": page,
		"by":   columnName,
		"sort": direction,
	})

	return hb.Hyperlink().
		HTML(tableLabel).
		Child(sortingIndicator(columnName, sortBy, sortOrder)).
		Href(linkURL)
}

// sortingIndicator returns the sorting indicator (up/down arrow) for a column
// This is a standalone utility function for better reusability across controllers
func sortingIndicator(columnName, sortByColumnName, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;")

	return sortingIndicator
}
