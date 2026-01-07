package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-tasks/db"
	"github.com/jackc/pgx/v4"
)

type ITaskRepository interface {
	Create(ctx context.Context, task models.Task) (*models.Task, error)
	List(ctx context.Context) ([]*models.Task, error)
	GetById(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task models.Task) (*models.Task, error)
}

type PgTaskRepository struct {
	db db.DBInterface
}

func NewPgTaskRepository(db db.DBInterface) *PgTaskRepository {
	return &PgTaskRepository{
		db: db,
	}
}

func (tr *PgTaskRepository) Create(ctx context.Context, task models.Task) (*models.Task, error) {
	query := `
		INSERT INTO tasks (id, title, completed) 
		VALUES ($1, $2, $3)
		RETURNING id, title, completed
		`

	var newTask models.Task

	err := tr.db.QueryRow(ctx, query, task.Id, task.Title, task.Completed).Scan(&newTask.Id, &newTask.Title, &newTask.Completed)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (tr *PgTaskRepository) List(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT id, title, completed 
				FROM tasks`
	rows, err := tr.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := new(models.Task)
		err := rows.Scan(&task.Id, &task.Title, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tr *PgTaskRepository) GetById(ctx context.Context, id string) (*models.Task, error) {
	query := `SELECT id, title, completed
				FROM tasks
				WHERE id = $1`

	var task models.Task

	err := tr.db.QueryRow(ctx, query, id).Scan(&task.Id, &task.Title, &task.Completed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (tr *PgTaskRepository) Update(ctx context.Context, task models.Task) (*models.Task, error) {
	query := `UPDATE tasks
				SET title = $1, completed = $2, updated_at = $3 
				WHERE id = $4
				RETURNING id, title, completed;`

	err := tr.db.QueryRow(ctx, query, task.Title, task.Completed, time.Now(), task.Id).Scan(&task.Id, &task.Title, &task.Completed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}
