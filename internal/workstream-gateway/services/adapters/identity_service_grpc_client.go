package adapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/identity/v1"
	"github.com/danilobml/workstream/internal/platform/dtos"
	"github.com/danilobml/workstream/internal/platform/grpcutils"
	"google.golang.org/grpc"
)

type IdentityClient struct {
	pb pb.IdentityServiceClient
}

func NewIdentityServiceClient(conn grpc.ClientConnInterface) *IdentityClient {
	return &IdentityClient{pb: pb.NewIdentityServiceClient(conn)}
}

func (c *IdentityClient) Register(ctx context.Context, registerReq dtos.RegisterRequest) (dtos.RegisterResponse, error) {
	resp, err := c.pb.Register(ctx, &pb.RegisterRequest{Email: registerReq.Email, Password: registerReq.Password, Roles: registerReq.Roles})
	if err != nil {
		return dtos.RegisterResponse{}, grpcutils.ParseGrpcError(err)
	}

	registerResp := dtos.RegisterResponse{
		Token: resp.Token,
	}

	return registerResp, nil
}

func (c *IdentityClient) Login(ctx context.Context, loginReq dtos.LoginRequest) (dtos.LoginResponse, error) {
	resp, err := c.pb.Login(ctx, &pb.LoginRequest{Email: loginReq.Email, Password: loginReq.Password})
	if err != nil {
		return dtos.LoginResponse{}, grpcutils.ParseGrpcError(err)
	}

	loginResp := dtos.LoginResponse{
		Token: resp.Token,
	}

	return loginResp, nil
}
