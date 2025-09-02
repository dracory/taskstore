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

	if controller == pathQueueCreate {
		return queueCreate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueDelete {
		return queueDelete(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueDetails {
		return queueDetails(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueManager {
		return queueManager(a.logger, a.store, a.layout).ToTag(a.response, a.request)
	}

	if controller == pathQueueParameters {
		return queueParameters(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueRequeue {
		return queueRequeue(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueTaskRestart {
		return queueTaskRestart(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskCreate {
		return taskCreate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskDelete {
		return taskDelete(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskManager {
		return taskManager(a.logger, a.store, a.layout).ToTag(a.response, a.request)
	}

	if controller == pathTaskUpdate {
		return taskUpdate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueCreate {
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
