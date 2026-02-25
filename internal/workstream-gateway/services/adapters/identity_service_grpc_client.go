package adapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/identity/v1"
	"github.com/danilobml/workstream/internal/platform/dtos"
	"github.com/danilobml/workstream/internal/platform/grpcutils"
	"github.com/google/uuid"
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

func (c *IdentityClient) ListAllUsers(ctx context.Context) (dtos.GetAllUsersResponse, error) {
	resp, err := c.pb.ListAllUsers(ctx, &pb.ListAllUsersRequest{})
	if err != nil {
		return dtos.GetAllUsersResponse{}, grpcutils.ParseGrpcError(err)
	}

	var respUsers []dtos.ResponseUser

	for _, user := range resp.GetUsers() {
		id, err := uuid.Parse(user.Id)
		if err != nil {
			return nil, err
		}

		user := dtos.ResponseUser{
			ID:       id,
			Email:    user.Email,
			Roles:    getResponseRoles(user.Roles),
			IsActive: user.IsActive,
		}

		respUsers = append(respUsers, user)
	}

	return respUsers, nil
}

func (c *IdentityClient) Unregister(ctx context.Context, unregisterRequest dtos.UnregisterRequest) error {
	_, err := c.pb.Unregister(ctx, &pb.UnregisterRequest{Id: unregisterRequest.Id.String()})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}

func (c *IdentityClient) RemoveUser(ctx context.Context, req dtos.RemoveUserRequest) error {
	_, err := c.pb.RemoveUser(ctx, &pb.RemoveUserRequest{Id: req.Id.String()})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}

func (c *IdentityClient) GetUser(ctx context.Context, req dtos.GetUserRequest) (dtos.ResponseUser, error) {
	resp, err := c.pb.GetUser(ctx, &pb.GetUserRequest{Id: req.Id.String()})
	if err != nil {
		return dtos.ResponseUser{}, grpcutils.ParseGrpcError(err)
	}

	user := resp.User

	id, err := uuid.Parse(user.Id)
	if err != nil {
		return dtos.ResponseUser{}, err
	}

	respUser := dtos.ResponseUser{
		ID:       id,
		Email:    user.Email,
		Roles:    getResponseRoles(user.Roles),
		IsActive: user.IsActive,
	}

	return respUser, nil
}

func (c *IdentityClient) RequestPasswordReset(ctx context.Context, req dtos.RequestPasswordResetRequest) error {
	_, err := c.pb.RequestPasswordReset(ctx, &pb.RequestPasswordResetRequest{Email: req.Email})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}

func (c *IdentityClient) ResetPassword(ctx context.Context, req dtos.ResetPasswordRequest) error {
	_, err := c.pb.ResetPassword(ctx, &pb.ResetPasswordRequest{Password: req.Password, ResetToken: req.ResetToken})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}

func getResponseRoles(roles []*pb.Role) []string {
	var roleStrings []string
	for _, role := range roles {
		roleStrings = append(roleStrings, role.GetName())
	}
	return roleStrings
}
