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
		INSERT INTO tasks (id, title, completed, organization_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, completed, organization_id
		`

	var newTask models.Task

	err := tr.db.QueryRow(ctx, query, task.Id, task.Title, task.Completed, task.OrganizationId).Scan(&newTask.Id, &newTask.Title, &newTask.Completed, &newTask.OrganizationId)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (tr *PgTaskRepository) List(ctx context.Context, scope TaskScope) ([]*models.Task, error) {
	query := `SELECT id, title, completed, organization_id
				FROM tasks
				WHERE organization_id = $1`

	rows, err := tr.db.Query(ctx, query, scope.OrganizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*models.Task{}
	for rows.Next() {
		task := new(models.Task)
		err := rows.Scan(&task.Id, &task.Title, &task.Completed, &task.OrganizationId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tr *PgTaskRepository) GetById(ctx context.Context, id string, scope TaskScope) (*models.Task, error) {
	query := `SELECT id, title, completed, organization_id
				FROM tasks
				WHERE id = $1 AND organization_id = $2`

	var task models.Task
	err := tr.db.QueryRow(ctx, query, id, scope.OrganizationId).Scan(&task.Id, &task.Title, &task.Completed, &task.OrganizationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (tr *PgTaskRepository) Update(ctx context.Context, task models.Task, scope TaskScope) (*models.Task, error) {
	query := `UPDATE tasks
				SET title = $1, completed = $2, updated_at = $3 
				WHERE id = $4 AND organization_id = $5
				RETURNING id, title, completed, organization_id;`

	err := tr.db.QueryRow(ctx, query, task.Title, task.Completed, time.Now(), task.Id, scope.OrganizationId).Scan(&task.Id, &task.Title, &task.Completed, &task.OrganizationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}
