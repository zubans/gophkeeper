package config

import (
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
	ServerURL string
	ConfigDir string
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

// LoadClientConfig loads client configuration from env with defaults.
func LoadClientConfig() ClientConfig {
	LoadEnv()
	// Default config dir: ~/.gophkeeper
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, ".gophkeeper")
	return ClientConfig{
		ServerURL: getenv("SERVER_URL", "http://localhost:8080"),
		ConfigDir: getenv("CLIENT_CONFIG_DIR", defaultDir),
	}
}

// GetBool reads boolean env with default.
func GetBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}
