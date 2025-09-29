// Package main implements the GophKeeper CLI client.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	clientapp "gophkeeper/internal/app/client"
	"gophkeeper/internal/client"
	icli "gophkeeper/internal/client/cli"
)

func main() {
	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Simple command execution - pass os.Args and parse function
	if err := clientapp.Run(ctx, os.Args[1:], func(args []string) (interface{ Execute(*client.Client) error }, error) {
		return icli.ParseCommand(args)
	}); err != nil {
		fmt.Printf("Error: %v\n", err)
		stop()
		return
	}
}
