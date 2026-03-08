package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pjover/espigol/internal/domain/ports"
)

type server struct {
	httpServer *http.Server
	config     ports.ConfigService
}

func NewServer(config ports.ConfigService, db ports.DbService) ports.Server {
	port := config.GetString("server.port")

	mux := http.NewServeMux()

	// Basic health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register resource handlers
	NewPartnerHandler(db).RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	return &server{
		httpServer: srv,
		config:     config,
	}
}

func (s *server) Start() error {
	log.Printf("Starting HTTP server on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	log.Println("Stopping HTTP server...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}
	return nil
}
