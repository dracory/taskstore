package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
)

type Layout interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

// DefaultWebpage provides a built-in default webpage layout with Bootstrap, Vue, jQuery, and SweetAlert2
// This can be used as a fallback or for simple use cases without requiring custom Layout implementation
//
// SECURITY WARNING: The content parameter is rendered as raw HTML without sanitization.
// Ensure content is trusted or properly sanitized before passing to this function to avoid XSS vulnerabilities.
func DefaultWebpage(title, content string) *hb.HtmlWebpage {
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(defaultFavicon)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_3_8(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_3_8(),
		cdn.Jquery_3_7_1(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_11(),
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}
	.form-select {
		display: block;
		width: 100%;
		padding: .375rem 2.25rem .375rem .75rem;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1.5;
		color: #212529;
		background-color: #fff;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23343a40' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M2 5l6 6 6-6'/%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right .75rem center;
		background-size: 16px 12px;
		border: 1px solid #ced4da;
		border-radius: .25rem;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}

type UIOptions struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Logger         *slog.Logger
	Store          taskstore.StoreInterface
	Layout         Layout
}

func UI(options UIOptions) (hb.TagInterface, error) {
	if options.ResponseWriter == nil {
		return nil, errors.New("options.ResponseWriter is required")
	}

	if options.Request == nil {
		return nil, errors.New("options.Request is required")
	}

	if options.Store == nil {
		return nil, errors.New("options.Store is required")
	}

	if options.Logger == nil {
		return nil, errors.New("options.Logger is required")
	}

	if options.Layout == nil {
		return nil, errors.New("options.Layout is required")
	}

	admin := &admin{
		response: options.ResponseWriter,
		request:  options.Request,
		store:    options.Store,
		logger:   *options.Logger,
		layout:   options.Layout,
	}
	return admin.handler(), nil
}

type admin struct {
	response http.ResponseWriter
	request  *http.Request
	store    taskstore.StoreInterface
	logger   slog.Logger
	layout   Layout
}

func (a *admin) handler() hb.TagInterface {
	controller := req.GetStringTrimmed(a.request, "controller")

	if controller == "" {
		controller = pathHome
	}

	if controller == pathTaskQueueCreate {
		return taskQueueCreate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueDelete {
		return taskQueueDelete(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueDetails {
		return taskQueueDetails(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueManager {
		return taskQueueManager(a.logger, a.store, a.layout).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueParameters {
		return taskQueueParameters(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueRequeue {
		return taskQueueRequeue(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueTaskRestart {
		return taskQueueTaskRestart(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskDefinitionCreate {
		return taskDefinitionCreate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskDefinitionDelete {
		return taskDefinitionDelete(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskDefinitionManager {
		return taskDefinitionManager(a.logger, a.store, a.layout).ToTag(a.response, a.request)
	}

	if controller == pathTaskDefinitionUpdate {
		return taskDefinitionUpdate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskQueueCreate {
		return hb.Div().Child(hb.H1().HTML(controller))
	}

	if controller == pathHome {
		return home(a.logger, a.store, a.layout).ToTag(a.response, a.request)
	}

	a.layout.SetBody(hb.H1().HTML(controller).ToHTML())
	return hb.Raw(a.layout.Render(a.response, a.request))
}
