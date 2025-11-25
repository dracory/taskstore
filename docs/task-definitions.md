# Task Definitions

Task definitions describe **what** work can be performed. Each definition is identified by a unique alias and is backed by a Go handler that executes the work.

Task definitions are persisted in the `task_definition` table and are referenced by queued tasks and schedules.

## Core Concepts

- **Alias**
  - Unique string that identifies the task definition.
  - Used when enqueuing tasks (`TaskDefinitionEnqueueByAlias`).
- **Title & Description**
  - Human‑readable metadata to describe what the task does.
- **Status**
  - Controls whether the definition is active and can be used.
- **Handler**
  - Go type implementing `TaskHandlerInterface` (and optionally `TaskHandlerWithContext`).

## Implementing a Task Handler

A task handler encapsulates the logic for a given task definition.

```go
package tasks

import "github.com/dracory/taskstore"

type HelloWorldTask struct {
    taskstore.TaskHandlerBase
}

var _ taskstore.TaskHandlerInterface = (*HelloWorldTask)(nil)

func (task *HelloWorldTask) Alias() string {
    return "HelloWorldTask"
}

func (task *HelloWorldTask) Title() string {
    return "Hello World"
}

func (task *HelloWorldTask) Description() string {
    return "Say hello world"
}

func (task *HelloWorldTask) Handle() bool {
    name := task.GetParam("name")

    if name != "" {
        task.LogInfo("Hello " + name + "!")
    } else {
        task.LogInfo("Hello World!")
    }

    return true
}
```

Handlers may optionally support context cancellation by implementing `TaskHandlerWithContext`.

## Registering Task Definitions

Register handlers with the store so they can be discovered and persisted as task definitions.

```go
ctx := context.Background()

err := myTaskStore.TaskHandlerAdd(ctx, &HelloWorldTask{}, true)
if err != nil {
    // handle error
}
```

When `createIfMissing` is `true`, `TaskHandlerAdd` will create a corresponding task definition record if one does not already exist for the given alias.

## Enqueuing Tasks

To enqueue a task based on a definition alias:

```go
queuedTask, err := myTaskStore.TaskDefinitionEnqueueByAlias(
    taskstore.DefaultTaskQueue,
    "HelloWorldTask",
    map[string]any{
        "name": "Tom Jones",
    },
)
```

This creates a record in the task queue table. The task will be picked up by a running queue worker.

You can also provide convenience methods on your handler:

```go
func (task *HelloWorldTask) Enqueue(name string) (taskstore.TaskQueueInterface, error) {
    return myTaskStore.TaskDefinitionEnqueueByAlias(
        taskstore.DefaultTaskQueue,
        task.Alias(),
        map[string]any{"name": name},
    )
}
```

## Executing from the CLI

Task definitions can be executed directly from the command line using `TaskDefinitionExecuteCli`.

```go
// inside your main
myTaskStore.TaskDefinitionExecuteCli(os.Args[1], os.Args[1:])
```

Example:

```bash
go run . HelloWorldTask --name="Tom Jones"
```

## Best Practices

- **Use stable aliases** – avoid renaming aliases once used in production.
- **Keep handlers small** – delegate heavy logic to separate services or packages.
- **Log clearly** – use the logging helpers on `TaskHandlerBase` to record progress.
- **Validate parameters** – fail fast and log when inputs are invalid.
- **Use context where appropriate** – implement `TaskHandlerWithContext` for long‑running tasks.
