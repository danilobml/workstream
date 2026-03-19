package repositories

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
)

type UserRepositoryInMemory struct {
	data []models.User
}

func NewUserRepositoryInMemory() *UserRepositoryInMemory {
	return &UserRepositoryInMemory{
		data: make([]models.User, 0),
	}
}

func (ur *UserRepositoryInMemory) List(ctx context.Context, scope UserScope) ([]*models.User, error) {
	usersResp := make([]*models.User, 0)

	for i := range ur.data {
		user := &ur.data[i]

		if scope.OrganizationID != nil && user.OrganizationId != *scope.OrganizationID {
			continue
		}

		usersResp = append(usersResp, user)
	}

	return usersResp, nil
}

func (ur *UserRepositoryInMemory) FindById(ctx context.Context, id uuid.UUID, scope UserScope) (*models.User, error) {
	for i := range ur.data {
		user := &ur.data[i]

		if user.ID != id {
			continue
		}

		if scope.OrganizationID != nil && user.OrganizationId != *scope.OrganizationID {
			return nil, errs.ErrNotFound
		}

		return user, nil
	}

	return nil, errs.ErrNotFound
}

func (ur *UserRepositoryInMemory) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	for i := range ur.data {
		if ur.data[i].Email == email {
			return &ur.data[i], nil
		}
	}

	return nil, errs.ErrNotFound
}

func (ur *UserRepositoryInMemory) Create(ctx context.Context, user models.User) error {
	existingUser, _ := ur.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return errs.ErrAlreadyExists
	}

	ur.data = append(ur.data, user)

	return nil
}

func (ur *UserRepositoryInMemory) Update(ctx context.Context, user models.User, scope UserScope) error {
	existingUser, err := ur.FindById(ctx, user.ID, scope)
	if err != nil {
		return err
	}

	existingUser.Email = user.Email
	existingUser.HashedPassword = user.HashedPassword
	existingUser.Roles = user.Roles
	existingUser.IsActive = user.IsActive
	existingUser.OrganizationId = user.OrganizationId

	return nil
}

func (ur *UserRepositoryInMemory) SavePassword(ctx context.Context, userId uuid.UUID, newPassword string) error {
	for i := range ur.data {
		if ur.data[i].ID == userId {
			ur.data[i].HashedPassword = newPassword
			return nil
		}
	}

	return errs.ErrNotFound
}

func (ur *UserRepositoryInMemory) Delete(ctx context.Context, id uuid.UUID, scope UserScope) error {
	_, err := ur.FindById(ctx, id, scope)
	if err != nil {
		return err
	}

	ur.data = slices.DeleteFunc(ur.data, func(user models.User) bool {
		if user.ID != id {
			return false
		}

		if scope.OrganizationID != nil && user.OrganizationId != *scope.OrganizationID {
			return false
		}

		return true
	})

	return nil
}
