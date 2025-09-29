package clientapp

import (
	"context"
	"fmt"
	"time"

	"gophkeeper/internal/client"
	"gophkeeper/internal/config"
)

// App represents the CLI client application lifecycle.
type App struct {
	cfg config.ClientConfig
	cli *client.Client
}

// New initializes the client using provided config.
func New(cfg config.ClientConfig) (*App, error) {
	cli, err := client.NewClient(cfg.ServerURL, cfg.ConfigDir, cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return &App{cfg: cfg, cli: cli}, nil
}

// Run starts the CLI app and handles all command execution.
func Run(ctx context.Context, args []string, parseCommand func([]string) (interface{ Execute(*client.Client) error }, error)) error {
	// Load config with flag parsing
	cfg := config.LoadClientConfigWithFlags()

	// Create app
	app, err := New(cfg)
	if err != nil {
		return err
	}

	// Parse and execute command
	cmd, err := parseCommand(args)
	if err != nil {
		return err
	}
	if err := cmd.Execute(app.cli); err != nil {
		return err
	}

	// small wait to flush logs/output
	time.Sleep(50 * time.Millisecond)

	// nothing long-running; return so caller may handle shutdown
	return nil
}
