package federation

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projdocs/cli/pkg/federation/middleware"
	"github.com/projdocs/cli/pkg/federation/routes"
	"go.uber.org/zap"
)

// Server wraps the HTTP server and its dependencies.
type Server struct {
	http            *http.Server
	log             *zap.Logger
	db              *sql.DB
	shutdownTimeout time.Duration
}

// ServerConfig holds HTTP server parameters.
type ServerConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DefaultServerConfig returns sensible production defaults.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Addr:            "0.0.0.0:8000",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
}

// NewServer constructs a Server, wires middleware, and registers all routes.
func NewServer(cfg ServerConfig, db *sql.DB, _log *zap.Logger) *Server {

	log := _log.Named("api")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.NoRoute(middleware.NoRoute(log.Named("no_route")))
	r.NoMethod(middleware.NoMethod(log.Named("no_method")))
	r.Use(
		middleware.Recovery(log.Named("recovery")),
		middleware.RequestLogger(log.Named("request_logger")),
	)

	// register endpoints
	routes.RegisterAll(r.Group(""), db, log)

	return &Server{
		http: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		log: log,
		db:  db,
	}
}

// Serve begins listening. It blocks until the server stops.
func (s *Server) Serve(ctx context.Context) error {
	serveErr := make(chan error, 1)

	go func() {
		if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
		close(serveErr)
	}()

	s.log.Info("server started", zap.String("addr", s.http.Addr))

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		s.log.Info("shutdown signal received, commencing shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.http.Shutdown(shutdownCtx)
}
