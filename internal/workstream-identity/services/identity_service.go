package services

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/dtos"
	//"github.com/danilobml/workstream/internal/platform/models"
)

type IdentityService interface {
	// public:
	Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error

	// admin:
	ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error)
	RemoveUser(ctx context.Context, req dtos.RemoveUserRequest) error
	/* 	GetUserData(ctx context.Context) (dtos.ResponseUser, error)
	RequestPasswordReset(ctx context.Context, requestPassResetReq dtos.RequestPasswordResetRequest) error
	ResetPassword(ctx context.Context, resetPassRequest dtos.ResetPasswordRequest) error
	UpdateUserData(ctx context.Context, updateUserRequest dtos.UpdateUserRequest) error
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	CheckUser(ctx context.Context, checkUserReq dtos.CheckUserRequest) (dtos.CheckUserResponse, error) */
}
