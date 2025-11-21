package admin

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/req"

	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/sb"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const actionModalQueueFilterShow = "modal_queue_filter_show"

func taskQueueManager(logger slog.Logger, store taskstore.StoreInterface, layout Layout) *taskQueueManagerController {
	return &taskQueueManagerController{
		logger: logger,
		store:  store,
		layout: layout,
	}
}

type taskQueueManagerController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
	layout Layout
}

func (c *taskQueueManagerController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.layout.SetTitle("Queue Manager | Zeppelin")

	if errorMessage != "" {
		c.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.layout.Render(w, r))
	}

	if data.action == actionModalQueueFilterShow {
		return c.onModalRecordFilterShow(&data)
	}

	if data.action == actionModalQueuedTaskFilterShow {
		// return c.onModalQueuedTaskFilterShow(data)
	}

	htmxScript := `setTimeout(() => {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
		}
	}, 1000);`

	swalScript := `setTimeout(() => {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
		}
	}, 1000);`

	c.layout.SetBody(c.page(&data).ToHTML())
	c.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.layout.Render(w, r))
}

func (*taskQueueManagerController) onModalRecordFilterShow(data *taskQueueManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
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

	fieldQueueID := form.NewField(form.FieldOptions{
		Label: "Queue ID",
		Name:  "filter_queue_id",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formQueueID,
		Help:  `Find queue by reference number (ID).`,
	})

	filterForm := form.NewForm(form.FormOptions{
		ID:        "FormFilters",
		Method:    http.MethodGet,
		ActionURL: url(data.request, pathTaskQueueManager, map[string]string{}),
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "filter_status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the queue.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Canceled",
						Key:   taskstore.TaskQueueStatusCanceled,
					},
					{
						Value: "Deleted",
						Key:   taskstore.TaskQueueStatusDeleted,
					},
					{
						Value: "Failed",
						Key:   taskstore.TaskQueueStatusFailed,
					},
					{
						Value: "Paused",
						Key:   taskstore.TaskQueueStatusPaused,
					},
					{
						Value: "Queued",
						Key:   taskstore.TaskQueueStatusQueued,
					},
					{
						Value: "Running",
						Key:   taskstore.TaskQueueStatusRunning,
					},
					{
						Value: "Success",
						Key:   taskstore.TaskQueueStatusSuccess,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label: "Name",
				Name:  "filter_name",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formName,
				Help:  `Filter by name.`,
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
			fieldQueueID,
			form.NewField(form.FieldOptions{
				Label: "Path",
				Name:  "path",
				Type:  form.FORM_FIELD_TYPE_HIDDEN,
				Value: pathTaskQueueManager,
				Help:  `Path to this page.`,
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

func (controller *taskQueueManagerController) page(data *taskQueueManagerControllerData) hb.TagInterface {
	adminHeader := adminHeader(controller.store, &controller.logger, data.request)
	breadcrumbs := breadcrumbs(data.request, []Breadcrumb{
		{
			Name: "Queue Manager",
			URL:  url(data.request, pathTaskQueueManager, map[string]string{}),
		},
	})

	buttonQueueCreate := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New Task to Queue").
		HxGet(url(data.request, pathTaskQueueCreate, map[string]string{})).
		HxTarget("body").
		HxSwap("beforeend")

	title := hb.Heading1().
		HTML("Zeppelin. Queue Manager").
		Child(buttonQueueCreate)

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(title).
		Child(controller.tableRecords(data))
}

func (controller *taskQueueManagerController) tableRecords(data *taskQueueManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Name", taskstore.COLUMN_TITLE)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Alias", taskstore.COLUMN_ALIAS)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Reference", taskstore.COLUMN_ID)).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Status", taskstore.COLUMN_STATUS)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Started", "started_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Finished", "completed_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Duration", "duration")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", taskstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(queuedTask taskstore.TaskQueueInterface, _ int) hb.TagInterface {
				task, taskExists := lo.Find(data.taskList, func(t taskstore.TaskDefinitionInterface) bool {
					return t.ID() == queuedTask.TaskID()
				})

				taskName := lo.IfF(taskExists, func() string { return task.Title() }).Else("Unknown")
				taskAlias := lo.IfF(taskExists, func() string { return task.Alias() }).Else("Unknown")

				buttonDelete := hb.Button().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete task from queue").
					HxGet(url(data.request, pathTaskQueueDelete, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonParameters := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-list-stars")).
					Title("See queued task parameters").
					HxGet(url(data.request, pathTaskQueueParameters, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonDetails := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-info-circle-fill")).
					Title("See the details of the job run").
					HxGet(url(data.request, pathTaskQueueDetails, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonAddToQueue := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-database-add")).
					Title("Add as new task to the queue").
					HxGet(url(data.request, pathTaskQueueRequeue, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonRestart := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-arrow-clockwise")).
					Title("Restart this job").
					HxGet(url(data.request, pathTaskQueueTaskRestart, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				startedAtDate := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return queuedTask.StartedAtCarbon().Format("d M Y")
				}).Else("-")
				startedAtTime := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return queuedTask.StartedAtCarbon().ToTimeString()
				}).Else("-")
				completeddAtDate := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return queuedTask.CompletedAtCarbon().Format("d M Y")
				}).Else("-")
				completeddAtTime := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return queuedTask.CompletedAtCarbon().ToTimeString()
				}).Else("-")

				elapsedTime := lo.IfF(queuedTask.StartedAt() != "" && queuedTask.CompletedAt() != "", func() string {
					diffSeconds := queuedTask.CompletedAtCarbon().DiffAbsInSeconds(queuedTask.StartedAtCarbon())
					return cast.ToString(diffSeconds) + "s"
				}).Else("-")

				createdAtDate := queuedTask.CreatedAtCarbon().Format("d M Y")
				createdAtTime := queuedTask.CreatedAtCarbon().ToTimeString()

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(queuedTask.IsSuccess(), `color:green;`).
					StyleIf(queuedTask.IsRunning(), `color:silver;`).
					StyleIf(queuedTask.IsQueued(), `color:blue;`).
					StyleIf(queuedTask.IsFailed(), `color:red;`).
					HTML(queuedTask.Status())

				return hb.TR().
					// Name, Alias, Ref
					Child(hb.TD().
						Child(hb.Div().Text(taskName)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Alias: ").
							Text(taskAlias)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Ref: ").
							Text(queuedTask.ID()))).

					// Status
					Child(hb.TD().
						Child(status)).

					// Started At
					Child(hb.TD().
						Child(hb.Div().Text(startedAtDate)).
						Child(hb.Div().Text(startedAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Completed At
					Child(hb.TD().
						Child(hb.Div().Text(completeddAtDate)).
						Child(hb.Div().Text(completeddAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Duration
					Child(hb.TD().
						Child(hb.Div().Text(elapsedTime)).
						Style("white-space: nowrap;")).

					// Created At
					Child(hb.TD().
						Child(hb.Div().Text(createdAtDate)).
						Child(hb.Div().Text(createdAtTime)).
						Style("white-space: nowrap; font-size: 13px;")).

					// Actions
					Child(hb.TD().
						Style("text-align: center;").
						Child(buttonParameters).
						Child(buttonDetails).
						Child(buttonAddToQueue).
						Child(buttonRestart).
						Child(buttonDelete))
			})),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.recordCount), data.pageInt, data.perPage),
	})
}

func (controller *taskQueueManagerController) sortableColumnLabel(data *taskQueueManagerControllerData, tableLabel, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)

	if !isSelected {
		direction = sb.ASC
	}

	link := url(data.request, pathTaskQueueManager, map[string]string{
		"controller":     pathTaskQueueManager,
		"page":           "0",
		"by":             columnName,
		"sort":           direction,
		"date_from":      data.formCreatedFrom,
		"date_to":        data.formCreatedTo,
		"status":         data.formStatus,
		"filter_task_id": data.formTaskID,
		"queue_id":       data.formQueueID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *taskQueueManagerController) sortingIndicator(columnName, sortByColumnName, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}

func (controller *taskQueueManagerController) tableFilter(data *taskQueueManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(url(data.request, pathTaskQueueManager, map[string]string{
			"action":       actionModalQueueFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"queue_id":     data.formQueueID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing queues").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formQueueID != "" {
		description = append(description, hb.Span().Text("and queue id: "+data.formQueueID).ToHTML())
	}

	// 	if data.formTaskID != "" {

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *taskQueueManagerController) tablePagination(
	data *taskQueueManagerControllerData,
	count,
	page,
	perPage int,
) hb.TagInterface {
	url := url(data.request, pathTaskQueueManager, map[string]string{
		"status":       data.formStatus,
		"name":         data.formName,
		"created_from": data.formCreatedFrom,
		"created_to":   data.formCreatedTo,
		"by":           data.sortBy,
		"order":        data.sortOrder,
	})

	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (controller *taskQueueManagerController) prepareData(r *http.Request) (data taskQueueManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.queueID = req.GetStringTrimmed(r, "queue_id")

	data.page = req.GetStringTrimmed(r, "page")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(req.GetStringTrimmedOr(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = req.GetStringTrimmed(r, "sort")
	data.sortBy = req.GetStringTrimmed(r, "by")

	data.formCreatedFrom = req.GetStringTrimmed(r, "filter_created_from")
	data.formCreatedTo = req.GetStringTrimmed(r, "filter_created_to")
	data.formName = req.GetStringTrimmed(r, "filter_name")
	data.formQueueID = req.GetStringTrimmed(r, "filter_queue_id")
	data.formStatus = req.GetStringTrimmed(r, "filter_status")

	data.recordList, data.recordCount, err = controller.fetchRecordList(&data)

	if err != nil {
		controller.logger.Error("At queueManagerController > prepareData", "error", err.Error())
		return data, "error retrieving web queues"
	}

	data.taskList, err = controller.store.TaskDefinitionList(context.Background(), taskstore.TaskDefinitionQuery().
		SetOrderBy(taskstore.COLUMN_ALIAS).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		controller.logger.Error("At queueManagerController > prepareData", "error", err.Error())
		return data, "error retrieving tasks"
	}

	return data, ""
}

func (controller *taskQueueManagerController) fetchRecordList(data *taskQueueManagerControllerData) (records []taskstore.TaskQueueInterface, recordCount int64, err error) {
	queueIDs := []string{}

	if data.formQueueID != "" {
		queueIDs = append(queueIDs, data.formQueueID)
	}

	query := taskstore.TaskQueueQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if len(queueIDs) > 0 {
		query = query.SetIDIn(queueIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	// if data.formName != "" {
	// 	query = query.SetNameLike(data.formName)
	// }

	recordList, err := controller.store.TaskQueueList(context.Background(), query)

	if err != nil {
		return records, 0, err
	}

	recordCount, err = controller.store.TaskQueueCount(context.Background(), query)

	if err != nil {
		return []taskstore.TaskQueueInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type taskQueueManagerControllerData struct {
	request *http.Request
	action  string

	page      string
	pageInt   int
	perPage   int
	sortOrder string
	sortBy    string

	formStatus      string
	formName        string
	formCreatedFrom string
	formCreatedTo   string
	formQueueID     string
	formTaskID      string

	recordList  []taskstore.TaskQueueInterface
	recordCount int64

	queueID  string
	taskList []taskstore.TaskDefinitionInterface
}
