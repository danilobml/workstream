package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func StartGrpcListener(grpcPort string) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return nil, err
	}

	return lis, err
}

func RegisterGrpcServer(listener net.Listener, errCh chan<- error) {
	s := grpc.NewServer()

	hs := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, hs)

	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	log.Printf("workstream-tasks - gRPC listening on %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		errCh <- fmt.Errorf("workstream-tasks - gRPC serve failed: %w", err)
		return
	}
}
