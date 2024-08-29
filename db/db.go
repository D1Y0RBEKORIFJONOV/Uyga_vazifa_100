package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

type GetTaskParams struct {
	ID int32
}

type GetAuthorParams struct {
	ID int32
}

type AssignAuthorToTaskParams struct {
	TaskID   int32
	AuthorID int32
}

type Queries struct {
	db DBTX
}

func (q *Queries) GetTask(ctx context.Context, params GetTaskParams) (interface{}, error) {
	var task string
	query := "SELECT task_name FROM tasks WHERE id = $1"
	row := q.db.QueryRowContext(ctx, query, params.ID)
	if err := row.Scan(&task); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", params.ID)
		}
		return nil, fmt.Errorf("error querying task: %v", err)
	}
	return task, nil
}

func (q *Queries) GetAuthor(ctx context.Context, params GetAuthorParams) (interface{}, error) {
	var author string
	query := "SELECT author_name FROM authors WHERE id = $1"
	row := q.db.QueryRowContext(ctx, query, params.ID)
	if err := row.Scan(&author); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("author with ID %d not found", params.ID)
		}
		return nil, fmt.Errorf("error querying author: %v", err)
	}
	return author, nil
}

func (q *Queries) AssignAuthorToTask(ctx context.Context, params AssignAuthorToTaskParams) error {
	taskCheckQuery := "SELECT 1 FROM tasks WHERE id = $1"
	authorCheckQuery := "SELECT 1 FROM authors WHERE id = $1"

	var exists int
	if err := q.db.QueryRowContext(ctx, taskCheckQuery, params.TaskID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found", params.TaskID)
		}
		return fmt.Errorf("error checking task existence: %v", err)
	}

	if err := q.db.QueryRowContext(ctx, authorCheckQuery, params.AuthorID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("author with ID %d not found", params.AuthorID)
		}
		return fmt.Errorf("error checking author existence: %v", err)
	}

	assignQuery := "INSERT INTO task_assignments (task_id, author_id) VALUES ($1, $2) ON CONFLICT (task_id) DO UPDATE SET author_id = EXCLUDED.author_id"
	_, err := q.db.ExecContext(ctx, assignQuery, params.TaskID, params.AuthorID)
	if err != nil {
		return fmt.Errorf("error assigning author to task: %v", err)
	}
	return nil
}
