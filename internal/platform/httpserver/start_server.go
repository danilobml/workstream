package httpx

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danilobml/workstream/internal/platform/routes"
)

func StartServer(
	serviceName, httpPortName string,
	registerServiceRoutes func(*http.ServeMux),
	isReady func() error,
) error {
	port := os.Getenv(httpPortName)
	if port == "" {
		msg := fmt.Sprintf("%s: %s variable could not be retrieved from env", serviceName, httpPortName)
		return errors.New(msg)
	}

	addr := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()
	routes.RegisterHealthRoutes(mux, serviceName, isReady)

	if registerServiceRoutes != nil {
		registerServiceRoutes(mux)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	shutdownDone := make(chan struct{})
	go func() {
		waitForShutdown(srv, 5*time.Second)
		close(shutdownDone)
	}()

	log.Printf("%s listening on port %s...", serviceName, port)

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-shutdownDone
	return nil
}

func waitForShutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	<-stop
	log.Println("\nGracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}

	log.Println("Shutdown complete.")
}
