// Package main implements the GophKeeper server.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	serverapp "gophkeeper/internal/app/server"
	"gophkeeper/internal/config"
)

func main() {
	// Load config with flags moved to internal/config
	cfg := config.LoadServerConfigWithFlags()

	// Create cancellable context and handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server app and block until error or shutdown signal
	if err := serverapp.Run(ctx, cfg); err != nil {
		// Ignore server closed errors on graceful shutdown
		log.Printf("Server stopped: %v", err)
	}

	// Ensure we allow some time for cleanup
	time.Sleep(100 * time.Millisecond)
}
