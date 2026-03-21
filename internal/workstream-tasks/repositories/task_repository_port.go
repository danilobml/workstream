package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/google/uuid"
)

type TaskScope struct {
	OrganizationId uuid.UUID
}

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) (*models.Task, error)
	List(ctx context.Context, scope TaskScope) ([]*models.Task, error)
	GetById(ctx context.Context, id string, scope TaskScope) (*models.Task, error)
	Update(ctx context.Context, task models.Task, scope TaskScope) (*models.Task, error)
}