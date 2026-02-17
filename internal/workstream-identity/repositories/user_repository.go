package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	List(ctx context.Context) ([]*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
