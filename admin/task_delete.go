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
	"github.com/spf13/cast"
)

func taskDelete(logger slog.Logger, store taskstore.StoreInterface) *taskDeleteCeontroller {
	return &taskDeleteCeontroller{
		logger: logger,
		store:  store,
	}
}

type taskDeleteCeontroller struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *taskDeleteCeontroller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modal(data)
}

func (c *taskDeleteCeontroller) formSubmitted(data taskDeleteCeontrollerData) hb.TagInterface {
	for _, queuedTask := range data.relatedQueuesToDelete {
		if err := c.store.QueueSoftDeleteByID(queuedTask.ID()); err != nil {
			return hb.Swal(hb.SwalOptions{
				Icon:              "error",
				Title:             "Error",
				Text:              err.Error(),
				Position:          "top-right",
				ShowCancelButton:  false,
				ShowConfirmButton: false,
			})
		}
	}

	if err := c.store.TaskSoftDelete(data.task); err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Success",
			Text:              "Task successfully deleted.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *taskDeleteCeontroller) modal(data taskDeleteCeontrollerData) *hb.Tag {
	fieldDanger := form.NewField(form.FieldOptions{
		Type: form.FORM_FIELD_TYPE_RAW,
		Value: hb.Wrap().
			Child(hb.Paragraph().
				Child(hb.Text(`You are about to permanently delete this task definition:`))).
			Child(hb.Paragraph().
				Style("font-weight: bold;").
				Child(hb.Text(`"` + data.task.Title() + `"`))).
			Child(hb.Paragraph().
				Child(hb.Text(`This will also delete all the tasks created from this definition, including ` + cast.ToString(data.relatedQueuedQueuesCount) + ` pending tasks.`))).
			Child(hb.Paragraph().
				Child(hb.Text(`Are you sure you want to proceed?`))).
			Child(hb.Paragraph().
				Class("text-danger").
				Child(hb.Text(`Warning: This action cannot be undone.`))).
			ToHTML(),
		Required: true,
	})

	fieldTaskID := form.NewField(form.FieldOptions{
		Label:    "Task ID",
		Name:     "task_id",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.taskID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormTaskDelete",
		Fields: []form.FieldInterface{
			fieldDanger,
			fieldTaskID,
		},
	})

	modalCloseScript := `document.getElementById('ModalTaskDelete').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonDelete := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Delete").
		Class("btn btn-danger float-end").
		HxInclude(`#ModalTaskDelete`).
		HxPost(url(data.request, pathTaskDelete, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalTaskDelete").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Delete Task").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonDelete),
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

func (c *taskDeleteCeontroller) prepareData(r *http.Request) (data taskDeleteCeontrollerData, err error) {
	data.request = r

	data.taskID = req.GetStringTrimmed(r, "task_id")

	if data.taskID == "" {
		return data, errors.New("task_id is required")
	}

	data.task, err = c.store.TaskFindByID(data.taskID)

	if err != nil {
		return data, err
	}

	if data.task == nil {
		return data, errors.New("task not found")
	}

	if r.Method == http.MethodPost {
		data.relatedQueuesToDelete, err = c.store.QueueList(taskstore.QueueQuery().
			SetTaskID(data.task.ID()).
			SetColumns([]string{taskstore.COLUMN_ID}))

		if err != nil {
			return data, err
		}
	}

	if r.Method == http.MethodGet {
		data.relatedQueuedQueuesCount, err = c.store.QueueCount(taskstore.QueueQuery().
			SetTaskID(data.task.ID()).
			SetColumns([]string{taskstore.COLUMN_ID}))

		if err != nil {
			return data, err
		}
	}

	return data, nil
}

type taskDeleteCeontrollerData struct {
	request                  *http.Request
	taskID                   string
	task                     taskstore.TaskInterface
	relatedQueuesToDelete    []taskstore.QueueInterface
	relatedQueuedQueuesCount int64
}
