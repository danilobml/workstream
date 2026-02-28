package services

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/dtos"
)

type IdentityService interface {
	Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error)
	Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error)
	Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error
	RequestPasswordReset(ctx context.Context, passResetReq dtos.RequestPasswordResetRequest) error
	ResetPassword(ctx context.Context, resetPassRequest dtos.ResetPasswordRequest) error
	ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error)
	RemoveUser(ctx context.Context, req dtos.RemoveUserRequest) error
	GetUser(ctx context.Context, req dtos.GetUserRequest) (dtos.ResponseUser, error)
	UpdateUser(ctx context.Context, updateUserRequest dtos.UpdateUserRequest) error
}
