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
	cfg := config.LoadServerConfigWithFlags()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := serverapp.Run(ctx, cfg); err != nil {
		log.Printf("Server stopped: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
}
