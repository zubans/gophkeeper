package serverapp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gophkeeper/internal/config"
	"gophkeeper/internal/database"
	dbm "gophkeeper/internal/database/migrations"
	"gophkeeper/internal/server"

	goose "github.com/pressly/goose/v3"
)

// App encapsulates server resources and lifecycle.
type App struct {
	httpServer *http.Server
	db         *database.DB
}

// New creates and initializes the server application (DB, migrations, HTTP server).
func New(cfg config.ServerConfig) (*App, error) {
	// Build database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Initialize database
	db, err := database.NewDB(connStr)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	// Run migrations using goose with embedded FS
	goose.SetBaseFS(dbm.ServerMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db.Conn(), "server"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	// Initialize HTTP server
	handler := server.NewServer(db, cfg.JWTSecret, cfg.EncryptionKey)
	httpSrv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}

	return &App{httpServer: httpSrv, db: db}, nil
}

// Start runs the HTTP server in a blocking manner.
func (a *App) Start() error {
	log.Printf("Starting server on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server and closes DB connection.
func (a *App) Shutdown(ctx context.Context) error {
	// Gracefully shutdown HTTP server
	httpErr := a.httpServer.Shutdown(ctx)

	// Ensure DB is closed after server stops accepting new connections
	dbErr := a.db.Close()

	// Prefer returning HTTP error if present
	if httpErr != nil {
		return httpErr
	}
	return dbErr
}

// Run convenience helper that starts the server and handles context cancellation.
func Run(ctx context.Context, cfg config.ServerConfig) error {
	app, err := New(cfg)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() { errCh <- app.Start() }()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return app.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
