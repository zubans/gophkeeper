package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// ServerConfig holds server runtime configuration.
type ServerConfig struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	EncryptionKey string
}

// ClientConfig holds client runtime configuration.
type ClientConfig struct {
	ServerURL     string
	ConfigDir     string
	EncryptionKey string
}

// LoadEnv loads .env if present (non-fatal if missing).
func LoadEnv() {
	_ = godotenv.Load()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// LoadServerConfig loads server configuration from env with defaults.
func LoadServerConfig() ServerConfig {
	LoadEnv()
	return ServerConfig{
		Port:          getenv("PORT", "8080"),
		DBHost:        getenv("DB_HOST", "localhost"),
		DBPort:        getenv("DB_PORT", "5432"),
		DBUser:        getenv("DB_USER", "gophkeeper"),
		DBPassword:    getenv("DB_PASSWORD", "password"),
		DBName:        getenv("DB_NAME", "gophkeeper"),
		JWTSecret:     getenv("JWT_SECRET", "your-secret-key"),
		EncryptionKey: getenv("ENCRYPTION_KEY", "your-encryption-key"),
	}
}

// LoadServerConfigWithFlags loads server config and applies flag overrides.
func LoadServerConfigWithFlags() ServerConfig {
	cfg := LoadServerConfig()

	// Parse flags for overrides
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

	return cfg
}

// LoadClientConfig loads client configuration from env with defaults.
func LoadClientConfig() ClientConfig {
	LoadEnv()
	// Default config dir: ~/.gophkeeper
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, ".gophkeeper")
	return ClientConfig{
		ServerURL:     getenv("SERVER_URL", "http://localhost:8080"),
		ConfigDir:     getenv("CLIENT_CONFIG_DIR", defaultDir),
		EncryptionKey: getenv("ENCRYPTION_KEY", "your-encryption-key"),
	}
}

// LoadClientConfigWithFlags loads client config and applies flag overrides.
func LoadClientConfigWithFlags() ClientConfig {
	cfg := LoadClientConfig()

	// Parse flags for overrides
	var (
		serverURL = flag.String("server", "", "Server URL (override env)")
		configDir = flag.String("config", "", "Configuration directory (override env)")
	)
	flag.Parse()

	if *serverURL != "" {
		cfg.ServerURL = *serverURL
	}
	if *configDir != "" {
		cfg.ConfigDir = *configDir
	}

	return cfg
}

func GetBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}
