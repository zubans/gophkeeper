// Package main implements the GophKeeper server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gophkeeper/internal/config"
	"gophkeeper/internal/database"
	dbm "gophkeeper/internal/database/migrations"
	"gophkeeper/internal/migrate"
	"gophkeeper/internal/server"
)

func main() {
	var (
		port       = flag.String("port", "", "Server port (override env)")
		dbHost     = flag.String("db-host", "", "Database host")
		dbPort     = flag.String("db-port", "", "Database port")
		dbUser     = flag.String("db-user", "", "Database user")
		dbPassword = flag.String("db-password", "", "Database password")
		dbName     = flag.String("db-name", "", "Database name")
		jwtSecret  = flag.String("jwt-secret", "", "JWT secret key")
		encKey     = flag.String("encryption-key", "", "Data encryption key")
	)
	flag.Parse()

	cfg := config.LoadServerConfig()
	if *port != "" {
		cfg.Port = *port
	}
	if *dbHost != "" {
		cfg.DBHost = *dbHost
	}
	if *dbPort != "" {
		cfg.DBPort = *dbPort
	}
	if *dbUser != "" {
		cfg.DBUser = *dbUser
	}
	if *dbPassword != "" {
		cfg.DBPassword = *dbPassword
	}
	if *dbName != "" {
		cfg.DBName = *dbName
	}
	if *jwtSecret != "" {
		cfg.JWTSecret = *jwtSecret
	}
	if *encKey != "" {
		cfg.EncryptionKey = *encKey
	}

	// Build database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Initialize database
	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	runner := &migrate.Runner{DB: db.Conn(), FS: dbm.ServerMigrations, Dir: "server"}
	if err := runner.Run("postgres"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize HTTP server
	srv := server.NewServer(db, cfg.JWTSecret, cfg.EncryptionKey)
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}
