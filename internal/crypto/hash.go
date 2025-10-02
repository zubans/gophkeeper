package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	combined := make([]byte, 64)
	copy(combined[:32], salt)
	copy(combined[32:], hash)
	return base64.StdEncoding.EncodeToString(combined), nil
}
func VerifyPassword(password, hashedPassword string) (bool, error) {
	decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
	if len(decoded) != 64 {
		return false, fmt.Errorf("invalid hash format: expected 64 bytes, got %d", len(decoded))
	}
	salt := decoded[:32]
	hash := decoded[32:]
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	expectedHash := hasher.Sum(nil)
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1, nil
}
