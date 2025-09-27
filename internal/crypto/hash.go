// Package crypto provides encryption and decryption functionality for GophKeeper.
package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
)

// HashPassword creates a secure hash of the password using SHA-256 with salt.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Create hash of salt + password
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	
	// Combine salt and hash
	combined := make([]byte, 64)
	copy(combined[:32], salt)
	copy(combined[32:], hash)
	
	return base64.StdEncoding.EncodeToString(combined), nil
}

// VerifyPassword verifies a password against its hash.
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

	// Create expected hash
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	expectedHash := hasher.Sum(nil)
	
	// Compare hashes using constant time comparison
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1, nil
}
