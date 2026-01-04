package readiness

import (
	"errors"
	"os"

	"github.com/danilobml/workstream/internal/workstream-gateway/grpc"
)

func IsReady() error {
	grpcAddr := os.Getenv("TASKS_GRPC_ADDR")
	if grpcAddr == "" {
		return errors.New("missing TASKS_GRPC_ADDR")
	}

	conn, err := grpc.GetClient()
	if err != nil {
		return err
	}

	return grpc.CheckTasksHealth(conn)
}
