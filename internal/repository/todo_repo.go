package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"work/todo/internal/model"
)

var (
	ErrNotFound = errors.New("record not found")
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {

	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, title string) (*model.Todo, error) {
	query := `
		INSERT INTO todos (title, completed, created_at, updated_at)
		VALUES ($1, false, NOW(), NOW())
		RETURNING id, title, completed, created_at, updated_at
	`

	var todo model.Todo

	row := r.db.QueryRowContext(ctx, query, title)
	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create todo: %w", err)
	}
	return &todo, nil
}

func (r *TodoRepository) GetAll(ctx context.Context) ([]model.Todo, error) {
	query := `SELECT id, title, completed, created_at, updated_at FROM todos ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("todos cannot fetched: %w", err)
	}
	defer rows.Close()
	todos := make([]model.Todo, 0)

	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error while reading row: %w", err)
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row loop error: %w", err)
	}
	return todos, nil
}

func (r *TodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	query := `SELECT id, title, completed, created_at, updated_at FROM todos WHERE id=$1`
	var todo model.Todo

	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unable to scane todo: %w", err)
	}

	return &todo, nil
}

func (r *TodoRepository) Update(ctx context.Context, id int64, req model.UpdateTodoRequest) (*model.Todo, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}

	if req.Completed != nil {
		existing.Completed = *req.Completed
	}

	existing.UpdatedAt = time.Now()

	query := `UPDATE todos
			 SET  title=$1, completed = $2, updated_at =$3
			 WHERE id = $4
			`

	_, err = r.db.ExecContext(ctx, query, existing.Title, existing.Completed, existing.UpdatedAt, id)

	if err != nil {
		return nil, fmt.Errorf("can not update todo: %w", err)
	}

	return existing, nil
}

func (r *TodoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}
