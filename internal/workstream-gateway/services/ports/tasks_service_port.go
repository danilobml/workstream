package ports

import (
	"context"
	
	"github.com/danilobml/workstream/internal/workstream-gateway/models"
)

type TasksServicePort interface {
	CreateTask(ctx context.Context, title string) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
}
