package grpcutils

import (
	"errors"
	"fmt"

	"github.com/danilobml/workstream/internal/platform/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseGrpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errs.ErrServerError
	}

	switch st.Code() {
	case codes.NotFound:
		return errs.ErrNotFound
	case codes.InvalidArgument:
		return errs.ErrBadRequest
	case codes.Unauthenticated:
		return errs.ErrUnauthorized
	case codes.PermissionDenied:
		return errs.ErrUnauthorized
	case codes.AlreadyExists:
		return errs.ErrAlreadyExists
	default:
		return fmt.Errorf("%w: %v", errs.ErrServerError, err)
	}
}

func ParseCustomError(err error) error {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, errs.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, errs.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, "invalid token")
	case errors.Is(err, errs.ErrUnauthorized):
		return status.Error(codes.PermissionDenied, "unauthorized")
	case errors.Is(err, errs.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
