package clientapp

import (
	"context"
	"fmt"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/cli"
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"time"
)

type App struct {
	cfg config.ClientConfig
	cli cli.ClientInterface
}

func New(cfg config.ClientConfig) (*App, error) {
	if err := logger.Init(cfg.LogLevel, cfg.LogFile); err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	logger.Info("Initializing client application")
	logger.Debug("Server URL: %s", cfg.ServerURL)
	logger.Debug("Config directory: %s", cfg.ConfigDir)

	cli, err := client.NewClient(cfg.ServerURL, cfg.ConfigDir, cfg.EncryptionKey)
	if err != nil {
		logger.Error("Failed to create client: %v", err)
		return nil, fmt.Errorf("create client: %w", err)
	}

	logger.Info("Client initialized successfully")
	return &App{cfg: cfg, cli: cli}, nil
}
func Run(_ context.Context, args []string, parseCommand func([]string) (cli.Command, error)) error {
	cfg := config.LoadClientConfigWithFlags()
	app, err := New(cfg)
	if err != nil {
		logger.Error("Failed to create client app: %v", err)
		return err
	}

	logger.Debug("Parsing command: %v", args)
	cmd, err := parseCommand(args)
	if err != nil {
		logger.Error("Failed to parse command: %v", err)
		return err
	}

	logger.Info("Executing command: %T", cmd)
	if err := cmd.Execute(app.cli); err != nil {
		logger.Error("Command execution failed: %v", err)
		return err
	}

	logger.Info("Command executed successfully")
	logger.Close()
	time.Sleep(50 * time.Millisecond)
	return nil
}
