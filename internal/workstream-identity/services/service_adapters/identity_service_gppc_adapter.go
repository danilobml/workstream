package serviceadapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/identity/v1"
	"github.com/danilobml/workstream/internal/platform/dtos"
	"github.com/danilobml/workstream/internal/workstream-identity/services"
)

type IdentityGrpcAdapter struct {
	pb.UnimplementedIdentityServiceServer
	svc services.IdentityService
}

func NewIdentityGrpcAdapter(svc services.IdentityService) *IdentityGrpcAdapter {
	return &IdentityGrpcAdapter{svc: svc}
}

func (a *IdentityGrpcAdapter) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	out, err := a.svc.Register(ctx, dtos.RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Roles: req.GetRoles(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Token: out.Token}, nil
}

func (a *IdentityGrpcAdapter) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	out, err := a.svc.Login(ctx, dtos.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: out.Token}, nil
}
