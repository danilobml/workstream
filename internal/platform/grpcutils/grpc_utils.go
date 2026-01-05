package grpcutils

import (
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
		return fmt.Errorf("%w: %v", errs.ErrBadRequest, err)

	default:
		return fmt.Errorf("%w: %v", errs.ErrServerError, err)
	}
}
