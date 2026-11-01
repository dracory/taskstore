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

	active := r.URL.Query().Get("controller")
	if active == "" {
		active = pathHome
	}

	tab := func(label, path, key string, count int64) hb.TagInterface {
		class := "nav-link"
		if key == active {
			class += " active"
		}
		link := hb.Hyperlink().
			HTML(label).
			Href(url(r, path, nil)).
			Class(class)
		if count >= 0 {
			link.Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(count)))
		}
		return hb.LI().Class("nav-item").Child(link)
	}

	return hb.UL().Class("nav nav-tabs mb-3").Children([]hb.TagInterface{
		tab("Dashboard", pathHome, pathHome, -1),
		tab("Task Queue", pathTaskQueueManager, pathTaskQueueManager, queueCount),
		tab("Task Definitions", pathTaskDefinitionManager, pathTaskDefinitionManager, taskCount),
	})
}

// pageHeader renders the page heading with optional right-aligned action
// buttons, matching the admin section pattern (breadcrumb + H2 + actions).
func pageHeader(title string, actions ...hb.TagInterface) hb.TagInterface {
	header := hb.Div().
		Class("d-flex justify-content-between align-items-center mb-3").
		Child(hb.H2().Text(title).Class("mb-0"))
	if len(actions) > 0 {
		header.Child(hb.Div().Class("d-flex gap-2 align-items-center").Children(actions))
	}
	return header
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
			Name: "Tasks",
			URL:  url(r, pathHome, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	return breadcrumbsUI(breadcrumbItems)
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
		Style("border:none;")

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
