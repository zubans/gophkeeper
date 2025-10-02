package main
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	clientapp "gophkeeper/internal/app/client"
	icli "gophkeeper/internal/client/cli"
)
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := clientapp.Run(ctx, os.Args[1:], icli.ParseCommand); err != nil {
		fmt.Printf("Error: %v\n", err)
		stop()
		return
	}
}
