package admin

import (
	"errors"
	"log/slog"
	"net/http"

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
	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
	// return nil
}
