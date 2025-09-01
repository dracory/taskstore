package admin

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/taskstore"
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

func (controller *queueManagerController) modalTaskFilters(data queueManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Queued Task Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	aliasOptions := lo.Map(data.taskList, func(task taskstore.TaskInterface, _ int) form.FieldOption {
		return form.FieldOption{
			Value: task.Title(),
			Key:   task.ID(),
		}
	})

	aliasOptions = append([]form.FieldOption{
		{
			Value: "",
			Key:   "",
		},
	}, aliasOptions...)

	filterForm := form.NewForm(form.FormOptions{
		ID:     "FormFilters",
		Method: http.MethodGet,
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "filter_status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the user.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Queued",
						Key:   taskstore.QueueStatusQueued,
					},
					{
						Value: "Running",
						Key:   taskstore.QueueStatusRunning,
					},
					{
						Value: "Canceled",
						Key:   taskstore.QueueStatusCanceled,
					},
					{
						Value: "Failed",
						Key:   taskstore.QueueStatusFailed,
					},
					{
						Value: "Success",
						Key:   taskstore.QueueStatusSuccess,
					},
					{
						Value: "Deleted",
						Key:   taskstore.QueueStatusDeleted,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label:   "Task",
				Name:    "filter_task_id",
				Type:    form.FORM_FIELD_TYPE_SELECT,
				Value:   data.formTaskID,
				Help:    `Filter by task.`,
				Options: aliasOptions,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created From",
				Name:  "filter_created_from",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedFrom,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created To",
				Name:  "filter_created_to",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedTo,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Queued Task ID",
				Name:  "filter_queue_id",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formQueueID,
				Help:  `Find user by reference number (ID).`,
			}),
		},
	}).Build()

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
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
