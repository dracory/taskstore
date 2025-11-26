# Task Queues

Task queues represent **individual executions** of task definitions. Each queue item stores the parameters, status, output, and timing information for a single run.

Queue items are persisted in the `task_queue` table and are processed by background workers managed by the `Store`.

## Queue Item Model

A queue item typically includes:

- **ID** – unique identifier.
- **TaskID** – reference to the task definition.
- **QueueName** – logical queue name (e.g. `default`, `emails`, `reports`).
- **Status** – lifecycle state (Queued, Running, Success, Failed, Canceled, Paused).
- **Parameters** – JSON‑encoded parameters passed to the handler.
- **Output / Details** – optional logs or result payloads.
- **Attempts** – how many times the item has been attempted.
- **Timestamps** – created, started, completed, updated, deleted/soft‑deleted.

## Enqueuing Work

Tasks are usually enqueued via `TaskDefinitionEnqueueByAlias`:

```go
queuedTask, err := myTaskStore.TaskDefinitionEnqueueByAlias(
    taskstore.DefaultQueueName, // or "emails" for a specific queue
    "SendWelcomeEmail",
    map[string]any{
        "user_id": 123,
    },
)

if err != nil {
    // handle error
}
```

This creates a new queue record in the specified queue.

## Processing Queues

> [!WARNING]
> **Deprecated:** The following methods are deprecated and will be removed in a future version. Use the new `TaskQueueRunner` pattern instead. See [Runners documentation](./runners.md).

The store provides multiple processing modes:

```go
ctx := context.Background()

// 1. Default queue (serial) - DEPRECATED
myTaskStore.TaskQueueRunDefault(ctx, 10, 2) // every 10s, unstuck after 2 mins

// 2. Named queue (serial) - DEPRECATED
myTaskStore.TaskQueueRunSerial(ctx, "emails", 10, 2)

// 3. Named queue (concurrent) - DEPRECATED
myTaskStore.TaskQueueRunConcurrent(ctx, "emails", 10, 2)
```

**Recommended approach using TaskQueueRunner:**

```go
ctx := context.Background()

// Create a task queue runner
queueRunner := taskstore.NewTaskQueueRunner(myTaskStore, taskstore.TaskQueueRunnerOptions{
    IntervalSeconds: 10,
    UnstuckMinutes:  1,
    QueueName:       "emails",
    Logger:          log.Default(),
})

// Start the runner
queueRunner.Start(ctx)

// Later: gracefully stop the runner
defer queueRunner.Stop()
```

### Concurrency Control

- Controlled via `MaxConcurrency` in `NewStoreOptions` (default: 10).
- Implemented using a semaphore to prevent resource exhaustion.
- Each named queue uses its own runner with tracked goroutines for orderly shutdown.

### Graceful Shutdown

> [!WARNING]
> **Deprecated:** The following methods are deprecated. Use `TaskQueueRunner.Stop()` instead. See [Runners documentation](./runners.md).

```go
myTaskStore.TaskQueueStop()               // Stop default queue - DEPRECATED
myTaskStore.TaskQueueStopByName("emails") // Stop only the "emails" queue - DEPRECATED
```

Both methods wait for in‑flight tasks to complete before returning.

## Status Lifecycle

A typical lifecycle for a queue item:

1. **Queued** – created via `TaskQueueCreate` or `TaskDefinitionEnqueueByAlias`.
2. **Running** – claimed by a worker using an atomic DB operation (`SELECT FOR UPDATE`).
3. **Success / Failed** – updated by the worker after handler execution.
4. **Soft‑deleted** – optional, hides the item while keeping historical data.

## Inspecting and Managing Queue Items

The store exposes methods to query and manage queue items:

- `TaskQueueFindByID` – look up a specific item.
- `TaskQueueList` – list items matching query criteria.
- `TaskQueueUpdate` – update status or metadata.
- `TaskQueueDeleteByID` / `TaskQueueSoftDeleteByID` – remove or hide items.

Use these methods to build admin tools, dashboards, or observability pipelines.

## Best Practices

- **Use separate queues for distinct workloads** (e.g. `emails`, `reports`, `webhooks`).
- **Tune `MaxConcurrency` per deployment** to match database and CPU capacity.
- **Add application‑level monitoring** around failures and retries.
- **Avoid huge payloads in `Parameters` or `Output`** – store large data elsewhere and reference it by ID.
- **Keep processing idempotent** so retries do not cause duplicate side effects.
