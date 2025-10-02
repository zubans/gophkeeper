package serverapp
import (
	"context"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"time"
	"gophkeeper/internal/config"
	"gophkeeper/internal/database"
	dbm "gophkeeper/internal/database/migrations"
	"gophkeeper/internal/server"
)
type App struct {
	httpServer *http.Server
	db         *database.DB
}
func New(cfg config.ServerConfig) (*App, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := database.NewDB(connStr)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	goose.SetBaseFS(dbm.ServerMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db.Conn(), "server"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	handler := server.NewServer(db, cfg.JWTSecret, cfg.EncryptionKey)
	httpSrv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	return &App{httpServer: httpSrv, db: db}, nil
}
func (a *App) Start() error {
	log.Printf("Starting server on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}
func (a *App) Shutdown(ctx context.Context) error {
	httpErr := a.httpServer.Shutdown(ctx)
	dbErr := a.db.Close()
	if httpErr != nil {
		return httpErr
	}
	return dbErr
}
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
