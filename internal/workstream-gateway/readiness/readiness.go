package readiness

import (
	"errors"
	"os"
)

func IsReady() error {
	grpcAddr := os.Getenv("TASKS_GRPC_ADDR")

	if grpcAddr == "" {
		return errors.New("gRPC server can't be initialized: TASKS_GRPC_ADDR variable could not be read from environment.")
	}

	return nil
}
