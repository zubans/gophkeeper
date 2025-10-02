package client
import (
	"fmt"
	"os"
	"path/filepath"
)
type TokenManagerImpl struct {
	configDir string
}
func NewTokenManager(configDir string) *TokenManagerImpl {
	return &TokenManagerImpl{
		configDir: configDir,
	}
}
func (t *TokenManagerImpl) SaveToken(token string) error {
	tokenFile := filepath.Join(t.configDir, "token")
	if err := os.WriteFile(tokenFile, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}
func (t *TokenManagerImpl) LoadToken() (string, error) {
	tokenFile := filepath.Join(t.configDir, "token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No token file exists, return empty string
		}
		return "", fmt.Errorf("failed to load token: %w", err)
	}
	return string(data), nil
}
func (t *TokenManagerImpl) ClearToken() error {
	tokenFile := filepath.Join(t.configDir, "token")
	if err := os.Remove(tokenFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear token: %w", err)
	}
	return nil
}
