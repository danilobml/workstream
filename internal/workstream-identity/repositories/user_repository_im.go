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

func (ur *UserRepositoryInMemory) List(ctx context.Context) ([]*models.User, error) {
	usersResp := make([]*models.User, 0, len(ur.data))
	for i := range ur.data {
		usersResp = append(usersResp, &ur.data[i])
	}
	if usersResp == nil {
		return []*models.User{}, nil
	}
	return usersResp, nil
}

func (ur *UserRepositoryInMemory) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for i := range ur.data {
		if ur.data[i].ID == id {
			return &ur.data[i], nil
		}
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

func (ur *UserRepositoryInMemory) Update(ctx context.Context, user models.User) error {
	existingUser, err := ur.FindById(ctx, user.ID)
	if err != nil {
		return err
	}

	existingUser.Email = user.Email
	existingUser.HashedPassword = user.HashedPassword
	existingUser.Roles = user.Roles
	existingUser.IsActive = user.IsActive

	return nil
}

func (ur *UserRepositoryInMemory) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ur.FindById(ctx, id)
	if err != nil {
		return err
	}

	ur.data = slices.DeleteFunc(ur.data, func(user models.User) bool {
		return user.ID == id
	})

	return nil
}
