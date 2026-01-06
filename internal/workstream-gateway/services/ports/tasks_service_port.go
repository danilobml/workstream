package ports

import (
	"context"
	
	"github.com/danilobml/workstream/internal/platform/models"
)

type TasksServicePort interface {
	CreateTask(ctx context.Context, title string) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	ListTasks(ctx context.Context) ([]*models.Task, error)
	CompleteTask(ctx context.Context, id string) error
}
