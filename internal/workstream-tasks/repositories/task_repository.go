package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-tasks/db"
)

type ITaskRepository interface {
	Create(ctx context.Context, task models.Task) (*models.Task, error)
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
	sqlStr := `
		INSERT INTO tasks (id, title, completed) 
		VALUES ($1, $2, $3)
		RETURNING id, title, completed
		`

	var newTask models.Task

	err := tr.db.QueryRow(ctx, sqlStr, task.Id, task.Title, task.Completed).Scan(&newTask.Id, &newTask.Title, &newTask.Completed)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}
