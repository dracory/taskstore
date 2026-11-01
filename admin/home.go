package admin

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/taskstore"
)

func home(logger slog.Logger, store taskstore.StoreInterface, layout Layout) *homeController {
	return &homeController{
		logger: logger,
		store:  store,
		layout: layout,
	}
}

type homeController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
	layout Layout
}

func (c *homeController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.layout.SetTitle("Dashboard | Tasks")

	if errorMessage != "" {
		c.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	c.layout.SetBody(c.page(data).ToHTML())
	c.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.layout.Render(w, r))
}

func (controller *homeController) page(data homeControllerData) hb.TagInterface {
	adminHeader := adminHeader(controller.store, &controller.logger, data.request)
	breadcrumbs := breadcrumbs(data.request, []Breadcrumb{
		{
			Name: "Dashboard",
			URL:  "",
		},
	})

	return hb.Div().
		Class("container py-4").
		Child(breadcrumbs).
		Child(pageHeader("Dashboard")).
		Child(adminHeader)
}

func (controller *homeController) prepareData(r *http.Request) (data homeControllerData, errorMessage string) {
	data.request = r

	return data, ""
}

type homeControllerData struct {
	request *http.Request
}
