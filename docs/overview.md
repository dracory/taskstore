# TaskStore Package Overview

## Introduction
`dracory/taskstore` is a robust, asynchronous durable task queue package designed to offload time-consuming or resource-intensive operations from the main application. It leverages a durable database (SQLite, MySQL, or PostgreSQL) for persistence.

## Core Concepts

### 1. Store
The `Store` is the central component that manages the connection to the database and provides methods to interact with task definitions and task queues. It handles:
- Database connection and migration (`AutoMigrate`).
- Task definition management (Create, Update, Delete, Find).
- Task queue management (Enqueue, Process, Status updates).

### 2. Task Definition
A task definition represents a unit of work. It is identified by a unique **Alias**.
- **Properties**: Alias, Title, Description, Status (Active/Canceled).
- **Handler**: Each task definition is associated with a `TaskHandlerInterface` implementation that defines the actual logic (`Handle` method).

### 3. Task Queue
A task queue item represents a specific instance of a task to be executed.
- **Properties**:
    - `ID`: Unique identifier.
    - `TaskID`: Reference to the parent Task.
    - `Status`: Current state (Queued, Running, Success, Failed, Canceled, Paused).
    - `Parameters`: JSON-encoded parameters for the task execution.
    - `Output`: Result or logs from the execution.
    - `Attempts`: Number of execution attempts.
    - `Timestamps`: CreatedAt, StartedAt, CompletedAt.

### 4. Recurrence Rules
The package supports recurring tasks via `RecurrenceRule`.
- **Task Queues Table**: Stores task execution instances (the task queue).

```mermaid
        text details
        int attempts
        datetime started_at
        datetime completed_at
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
```
