package clientapp
import (
	"context"
	"fmt"
	"time"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/cli"
	"gophkeeper/internal/config"
)
type App struct {
	cfg config.ClientConfig
	cli cli.ClientInterface
}
func New(cfg config.ClientConfig) (*App, error) {
	cli, err := client.NewClient(cfg.ServerURL, cfg.ConfigDir, cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return &App{cfg: cfg, cli: cli}, nil
}
func Run(_ context.Context, args []string, parseCommand func([]string) (cli.Command, error)) error {
	cfg := config.LoadClientConfigWithFlags()
	app, err := New(cfg)
	if err != nil {
		return err
	}
	cmd, err := parseCommand(args)
	if err != nil {
		return err
	}
	if err := cmd.Execute(app.cli); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	return nil
}
