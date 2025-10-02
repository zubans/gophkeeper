package serverapp

import (
	"context"
	"fmt"
	"gophkeeper/internal/config"
	"gophkeeper/internal/database"
	dbm "gophkeeper/internal/database/migrations"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server"
	"net/http"
	"time"

	"github.com/pressly/goose/v3"
)

type App struct {
	httpServer *http.Server
	db         *database.DB
}

func New(cfg config.ServerConfig) (*App, error) {
	if err := logger.Init(cfg.LogLevel, cfg.LogFile); err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	logger.Info("Initializing server application")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	logger.Debug("Connecting to database: host=%s port=%s dbname=%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := database.NewDB(connStr)
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	logger.Info("Database connection established")
	logger.Info("Running database migrations")
	goose.SetBaseFS(dbm.ServerMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set goose dialect: %v", err)
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db.Conn(), "server"); err != nil {
		logger.Error("Failed to run migrations: %v", err)
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	logger.Info("Database migrations completed successfully")

	logger.Info("Initializing HTTP server on port %s", cfg.Port)
	handler := server.NewServer(db, cfg.JWTSecret, cfg.EncryptionKey)
	httpSrv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	return &App{httpServer: httpSrv, db: db}, nil
}
func (a *App) Start() error {
	logger.Info("Starting server on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}
func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server")
	httpErr := a.httpServer.Shutdown(ctx)
	dbErr := a.db.Close()
	logger.Close()
	if httpErr != nil {
		return httpErr
	}
	return dbErr
}
func Run(ctx context.Context, cfg config.ServerConfig) error {
	app, err := New(cfg)
	if err != nil {
		logger.Error("Failed to create server app: %v", err)
		return err
	}
	errCh := make(chan error, 1)
	go func() { errCh <- app.Start() }()
	select {
	case <-ctx.Done():
		logger.Info("Received shutdown signal")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return app.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err != nil {
			logger.Error("Server error: %v", err)
		}
		return err
	}
}
