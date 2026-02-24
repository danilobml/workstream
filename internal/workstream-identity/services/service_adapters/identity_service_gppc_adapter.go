package serviceadapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/identity/v1"
	"github.com/danilobml/workstream/internal/platform/dtos"
	"github.com/danilobml/workstream/internal/platform/grpcutils"
	"github.com/danilobml/workstream/internal/platform/jwt"
	"github.com/danilobml/workstream/internal/workstream-identity/middleware"
	"github.com/danilobml/workstream/internal/workstream-identity/services"
	"github.com/google/uuid"
)

type IdentityGrpcAdapter struct {
	pb.UnimplementedIdentityServiceServer
	svc        services.IdentityService
	jwtManager *jwt.JwtManager
}

func NewIdentityGrpcAdapter(svc services.IdentityService, jwtManager *jwt.JwtManager) *IdentityGrpcAdapter {
	return &IdentityGrpcAdapter{svc: svc, jwtManager: jwtManager}
}

func (a *IdentityGrpcAdapter) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	out, err := a.svc.Register(ctx, dtos.RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Roles:    req.GetRoles(),
	})
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	return &pb.RegisterResponse{Token: out.Token}, nil
}

func (a *IdentityGrpcAdapter) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	out, err := a.svc.Login(ctx, dtos.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	return &pb.LoginResponse{Token: out.Token}, nil
}

func (a *IdentityGrpcAdapter) ListAllUsers(ctx context.Context, req *pb.ListAllUsersRequest) (*pb.UserListResponse, error) {
	ctx, err := middleware.AuthenticateGRPC(ctx, a.jwtManager)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	out, err := a.svc.ListAllUsers(ctx)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	var responseUsers []*pb.User
	for _, user := range out {

		roles := convertRoles(user.Roles)
		respUser := &pb.User{
			Id:       user.ID.String(),
			Email:    user.Email,
			Roles:    roles,
			IsActive: user.IsActive,
		}

		responseUsers = append(responseUsers, respUser)
	}

	return &pb.UserListResponse{Users: responseUsers}, nil
}

func (a *IdentityGrpcAdapter) Unregister(ctx context.Context, req *pb.UnregisterRequest) (*pb.UnregisterResponse, error) {
	ctx, err := middleware.AuthenticateGRPC(ctx, a.jwtManager)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	parsedId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	err = a.svc.Unregister(ctx, dtos.UnregisterRequest{Id: parsedId})
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	return &pb.UnregisterResponse{}, nil
}

func (a *IdentityGrpcAdapter) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	ctx, err := middleware.AuthenticateGRPC(ctx, a.jwtManager)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	parsedId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	err = a.svc.RemoveUser(ctx, dtos.RemoveUserRequest{Id: parsedId})
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	return &pb.RemoveUserResponse{}, nil
}

func (a *IdentityGrpcAdapter) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	err := a.svc.RequestPasswordReset(ctx, dtos.RequestPasswordResetRequest{Email: req.GetEmail()})
	if err != nil {
		return nil, grpcutils.ParseCustomError(err)
	}

	return &pb.RequestPasswordResetResponse{}, nil
}

func convertRoles(roles []string) []*pb.Role {
	var responseRoles []*pb.Role
	for _, role := range roles {
		respRole := &pb.Role{
			Name: role,
		}
		responseRoles = append(responseRoles, respRole)
	}
	return responseRoles
}
