package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
)

func queueTaskRestart(logger slog.Logger, store taskstore.StoreInterface) *queueTaskRestartController {
	return &queueTaskRestartController{
		logger: logger,
		store:  store,
	}
}

type queueTaskRestartController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *queueTaskRestartController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modal(data)
}

func (c *queueTaskRestartController) formSubmitted(data queueTaskRestartControllerData) hb.TagInterface {
	task, err := c.store.TaskFindByID(data.queue.TaskID())

	if err != nil {
		c.logger.Error("At queueTaskRestartController > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	if task == nil {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task not found", Position: "top-right"})
	}

	task.SetStatus(taskstore.QueueStatusQueued)

	if err := c.store.TaskUpdate(task); err != nil {
		c.logger.Error("At queueTaskRestartController > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queue successfully created.", Position: "top-right"})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *queueTaskRestartController) modal(data queueTaskRestartControllerData) *hb.Tag {
	modalID := `ModalQueueRequeue`
	formID := modalID + `Form`

	fieldInfo := form.NewField(form.FieldOptions{
		Type: form.FORM_FIELD_TYPE_RAW,
		Value: hb.Wrap().
			Child(hb.Paragraph().
				Child(hb.Text(`You are about to restart this task`))).
			Child(hb.Paragraph().
				Child(hb.Text(`All the actions executed by this task will be repeated.`))).
			Child(hb.Paragraph().
				Child(hb.Text(`Are you sure you want to proceed?`))).
			ToHTML(),
		Required: true,
	})

	fieldQueueID := form.NewField(form.FieldOptions{
		Label:    "Queue ID",
		Name:     "queue_id",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.queueID,
		Required: true,
	})

	formCreate := form.NewForm(form.FormOptions{
		ID: formID,
		Fields: []form.FieldInterface{
			fieldQueueID,
			fieldInfo,
		},
	})

	modalCloseScript := `document.getElementById('` + modalID + `').remove();document.getElementById('ModalBackdrop').remove();`

	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonRequeue := hb.Button().
		Child(hb.I().Class("bi bi-database-add me-2")).
		HTML("Add to queue").
		Class("btn btn-success float-end").
		HxInclude(`#` + modalID).
		HxPost(url(data.request, pathQueueRequeue, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Requeue as New Task to Queue").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formCreate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonRequeue),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}

func (c *queueTaskRestartController) prepareData(r *http.Request) (data queueTaskRestartControllerData, err error) {
	data.request = r
	data.queueID = req.GetStringTrimmed(r, "queue_id")

	if data.queueID == "" {
		return data, errors.New("queue_id is required")
	}

	if data.queue, err = c.store.QueueFindByID(data.queueID); err != nil {
		return data, err
	}

	if data.queue == nil {
		return data, errors.New("queue not found")
	}

	return data, nil
}

type queueTaskRestartControllerData struct {
	request *http.Request
	queueID string
	queue   taskstore.QueueInterface
}
