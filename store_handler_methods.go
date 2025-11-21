package taskstore

import (
	"context"
	"errors"
)

var errTaskMissing = errors.New("task not found")

func (store *Store) TaskHandlerAdd(ctx context.Context, taskHandler TaskHandlerInterface, createIfMissing bool) error {
	alias := taskHandler.Alias()
	task, err := store.TaskDefinitionFindByAlias(ctx, alias)

	if err != nil {
		return err
	}

	if task == nil && !createIfMissing {
		return errTaskMissing
	}

	if task == nil && createIfMissing {
		alias := taskHandler.Alias()
		title := taskHandler.Title()
		description := taskHandler.Description()

		task := NewTaskDefinition().
			SetStatus(TaskDefinitionStatusActive).
			SetAlias(alias).
			SetTitle(title).
			SetDescription(description)

		err := store.TaskDefinitionCreate(ctx, task)

		if err != nil {
			return err
		}
	}

	store.taskHandlers = append(store.taskHandlers, taskHandler)

	return nil
}

func (store *Store) TaskHandlerList() []TaskHandlerInterface {
	return store.taskHandlers
}
