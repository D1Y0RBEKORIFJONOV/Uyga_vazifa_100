// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	// Inserts a new author and returns it
	CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error)
	// Inserts a new task and returns it
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	// Inserts a new task-author relationship and ignores conflict
	CreateTaskAuthor(ctx context.Context, arg CreateTaskAuthorParams) error
	// Deletes a task by ID
	DeleteTask(ctx context.Context, id int32) error
	// Selects an author by ID
	GetAuthorByID(ctx context.Context, id int32) (Author, error)
	// Selects authors associated with a specific task
	GetAuthorsByTaskID(ctx context.Context, taskID int32) ([]Author, error)
	// Selects a task by ID
	GetTaskByID(ctx context.Context, id int32) (Task, error)
	// Lists all authors ordered by name
	ListAuthors(ctx context.Context) ([]Author, error)
	// Lists all tasks ordered by creation date
	ListTasks(ctx context.Context) ([]Task, error)
	// Updates a task and returns the updated task
	UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
}

var _ Querier = (*Queries)(nil)
