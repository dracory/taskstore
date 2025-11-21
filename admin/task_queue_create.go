package admin

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func taskQueueCreate(logger slog.Logger, store taskstore.StoreInterface) *taskQueueCreateController {
	return &taskQueueCreateController{
		logger: logger,
		store:  store,
	}
}

type taskQueueCreateController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *taskQueueCreateController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modalQueueCreate(data)
}

func (c *taskQueueCreateController) formSubmitted(data taskQueueCreateControllerData) hb.TagInterface {
	if data.formTaskID == "" {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task is required.", Position: "top-right"})
	}

	if data.formParameters == "" {
		data.formParameters = "{}"
	}

	if !isJSON(data.formParameters) {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task Parameters is not valid JSON", Position: "top-right"})
	}

	task, err := c.store.TaskDefinitionFindByID(context.Background(), data.formTaskID)

	if err != nil {
		c.logger.Error("At queueCreateController > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	if task == nil {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task not found", Position: "top-right"})
	}

	taskParametersAny := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.formParameters), &taskParametersAny); err != nil {
		c.logger.Error("At queueCreateController > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	taskParametersMap := cast.ToStringMap(taskParametersAny)

	_, err = c.store.TaskDefinitionEnqueueByAlias(context.Background(), task.Alias(), taskParametersMap)

	if err != nil {
		c.logger.Error("At queueCreateController > formSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queue successfully created.", Position: "top-right"})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *taskQueueCreateController) modalQueueCreate(data taskQueueCreateControllerData) *hb.Tag {
	modalID := `ModalQueueCreate`
	formID := modalID + `Form`
	fieldParameters := form.NewField(form.FieldOptions{
		Label:    "Parameters",
		Name:     "parameters",
		Type:     form.FORM_FIELD_TYPE_TEXTAREA,
		Value:    data.formParameters,
		Help:     "The parameters of the queued task. Must be valid JSON.",
		Required: true,
	})

	fieldParametersSize := form.NewField(form.FieldOptions{
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: hb.Style(`#` + formID + ` textarea[name="parameters"] { height: 200px; }`).ToHTML(),
	})

	fieldTaskID := form.NewField(form.FieldOptions{
		Label:    "Task",
		Name:     "task_id",
		Type:     form.FORM_FIELD_TYPE_SELECT,
		Value:    data.formTaskID,
		Help:     "The task that will be added to the queue to be executed.",
		Required: true,
		Options: func() []form.FieldOption {
			options := []form.FieldOption{{
				Value: "-- select task --",
				Key:   "",
			}}
			lo.Map(data.taskList, func(task taskstore.TaskDefinitionInterface, _ int) form.FieldOption {
				options = append(options, form.FieldOption{
					Value: task.Title(),
					Key:   task.ID(),
				})
				return form.FieldOption{}
			})
			return options
		}(),
	})

	formCreate := form.NewForm(form.FormOptions{
		ID: formID,
		Fields: []form.FieldInterface{
			fieldTaskID,
			fieldParametersSize,
			fieldParameters,
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

	buttonCreate := hb.Button().
		Child(hb.I().Class("bi bi-run me-2")).
		HTML("Create").
		Class("btn btn-success float-end").
		HxInclude(`#` + modalID).
		HxPost(url(data.request, pathTaskQueueCreate, nil)).
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
							Text("Add New Task to Queue").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formCreate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonCreate),
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

func (c *taskQueueCreateController) prepareData(r *http.Request) (data taskQueueCreateControllerData, err error) {
	data.request = r
	data.formParameters = req.GetStringTrimmed(r, "parameters")
	data.formStatus = req.GetStringTrimmed(r, "status")
	data.formTaskID = req.GetStringTrimmed(r, "task_id")

	if data.taskList, err = c.store.TaskDefinitionList(context.Background(), taskstore.TaskDefinitionQuery().
		SetOrderBy(taskstore.COLUMN_TITLE).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100)); err != nil {
		return data, err
	}

	return data, nil
}

type taskQueueCreateControllerData struct {
	request  *http.Request
	taskList []taskstore.TaskDefinitionInterface

	formTaskID     string
	formParameters string
	formStatus     string
}
