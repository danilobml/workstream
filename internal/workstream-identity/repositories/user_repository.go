package repositories

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/google/uuid"
)

type UserScope struct {
	OrganizationID *uuid.UUID
}

type UserRepository interface {
	List(ctx context.Context, scope UserScope) ([]*models.User, error)
	FindById(ctx context.Context, id uuid.UUID, scope UserScope) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User, scope UserScope) error
	SavePassword(ctx context.Context, userId uuid.UUID, newPassword string) error
	Delete(ctx context.Context, id uuid.UUID, scope UserScope) error
}
