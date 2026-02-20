package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/danilobml/workstream/internal/gen/identity/v1"
	"github.com/danilobml/workstream/internal/platform/jwt"
	"github.com/danilobml/workstream/internal/workstream-identity/services"
	serviceadapters "github.com/danilobml/workstream/internal/workstream-identity/services/service_adapters"
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

func RegisterGrpcServer(identityService services.IdentityService, listener net.Listener, jwtManager *jwt.JwtManager, errCh chan<- error) {
	srv := grpc.NewServer()

	hs := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, hs)

	pb.RegisterIdentityServiceServer(srv, serviceadapters.NewIdentityGrpcAdapter(identityService, jwtManager))

	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	log.Printf("workstream-tasks - gRPC listening on %v", listener.Addr())

	if err := srv.Serve(listener); err != nil {
		errCh <- fmt.Errorf("workstream-tasks - gRPC serve failed: %w", err)
		return
	}
}
