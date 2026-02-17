package main

import (
	"log"
	"net/http"
	"os"

	httpserver "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/platform/httpx/middleware"
	"github.com/danilobml/workstream/internal/platform/jwt"
	"github.com/danilobml/workstream/internal/workstream-gateway/grpc"
	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
	"github.com/danilobml/workstream/internal/workstream-gateway/readiness"
	"github.com/danilobml/workstream/internal/workstream-gateway/routes"
	services "github.com/danilobml/workstream/internal/workstream-gateway/services/adapters"
)

const (
	serviceName          = "workstream-gateway"
	httpPortName         = "GATEWAY_HTTP_PORT"
	tasksGrpcAddrName    = "TASKS_GRPC_ADDR"
	identityGrpcAddrName = "IDENTITY_GRPC_ADDR"
	identityApiKeyName   = "IDENTITY_API_KEY"
	secretJwtKeyName     = "SECRET_JWT_KEY"
)

func main() {
	tasksGrpcAddr := os.Getenv(tasksGrpcAddrName)
	if tasksGrpcAddr == "" {
		log.Fatal("unable to read TASKS_GRPC_ADDR from env")
	}

	identityGrpcAddr := os.Getenv(identityGrpcAddrName)
	if identityGrpcAddr == "" {
		log.Fatal("unable to read IDENTITY_GRPC_ADDR from env")
	}

	identityApiKey := os.Getenv(identityApiKeyName)
	if identityApiKey == "" {
		log.Fatal("unable to read IDENTITY_API_KEY from env")
	}

	secretJwtKey := os.Getenv(secretJwtKeyName)
	if secretJwtKey == "" {
		log.Fatal("unable to read SECRET_JWT_KEY from env")
	}

	tasksConn, err := grpc.CreateGrpcClient(tasksGrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer tasksConn.Close()

	identityConn, err := grpc.CreateGrpcClient(identityGrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer identityConn.Close()

	tasksService := services.NewTasksServiceClient(tasksConn)
	tasksHandler := handlers.NewTasksHandler(tasksService)
	tasksRouter := routes.RegisterTaskRoutes(tasksHandler)

	jwtManager := jwt.NewJwtManager([]byte(secretJwtKey))

	authMiddleware := middleware.Authenticate(jwtManager)
	identityService := services.NewIdentityServiceClient(identityConn)
	identityHandler := handlers.NewIdentityHandler(identityService, identityApiKey)
	identityRouter := routes.RegisterIdentityRoutes(identityHandler, authMiddleware)

	root := http.NewServeMux()
	root.Handle("/tasks/", http.StripPrefix("/tasks", tasksRouter))
	root.Handle("/identity/", http.StripPrefix("/identity", identityRouter))

	handler := middleware.ApplyMiddlewares(
		middleware.Recover,
		middleware.RequestId,
		middleware.Logger,
		middleware.Dos,
		middleware.RateLimit,
		middleware.Cors,
		middleware.Security,
	)(root)

	if err := httpserver.StartServer(
		serviceName,
		httpPortName,
		func(mux *http.ServeMux) {
			mux.Handle("/", handler)
		},
		readiness.IsReady,
	); err != nil {
		log.Fatal(err)
	}
}
