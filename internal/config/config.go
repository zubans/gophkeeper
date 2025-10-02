package config
import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"github.com/joho/godotenv"
)
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
type ClientConfig struct {
	ServerURL     string
	ConfigDir     string
	EncryptionKey string
}
func LoadEnv() {
	_ = godotenv.Load()
}
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
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
func LoadServerConfigWithFlags() ServerConfig {
	cfg := LoadServerConfig()
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
func LoadClientConfig() ClientConfig {
	LoadEnv()
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, ".gophkeeper")
	return ClientConfig{
		ServerURL:     getenv("SERVER_URL", "http://localhost:8080"),
		ConfigDir:     getenv("CLIENT_CONFIG_DIR", defaultDir),
		EncryptionKey: getenv("ENCRYPTION_KEY", "your-encryption-key"),
	}
}
func LoadClientConfigWithFlags() ClientConfig {
	cfg := LoadClientConfig()
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
