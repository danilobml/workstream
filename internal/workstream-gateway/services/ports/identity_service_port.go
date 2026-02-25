package ports

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/dtos"
)

type IdentityServicePort interface {
	Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error)
	Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error
	RemoveUser(ctx context.Context, req dtos.RemoveUserRequest) error
	RequestPasswordReset(ctx context.Context, req dtos.RequestPasswordResetRequest) error
	ResetPassword(ctx context.Context, req dtos.ResetPasswordRequest) error
	GetUser(ctx context.Context, req dtos.GetUserRequest) (dtos.ResponseUser, error)
	/* 	UpdateUser(ctx context.Context, updateUserRequest dtos.UpdateUserRequest) error */
}
