// Package tests provides tests for the crypto package.
package tests

import (
	"gophkeeper/internal/crypto"
	"testing"
	"time"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := crypto.NewJWTManager("test-secret")

	userID := "user123"
	username := "testuser"
	expiration := 1 * time.Hour

	token, err := manager.GenerateToken(userID, username, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if len(token) == 0 {
		t.Fatal("Token is empty")
	}

	// Token should have 3 parts separated by dots
	parts := splitToken(token)
	if len(parts) != 3 {
		t.Fatalf("Token should have 3 parts, got %d", len(parts))
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := crypto.NewJWTManager("test-secret")

	userID := "user123"
	username := "testuser"
	expiration := 1 * time.Hour

	token, err := manager.GenerateToken(userID, username, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("UserID mismatch. Expected: %s, Got: %s", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Fatalf("Username mismatch. Expected: %s, Got: %s", username, claims.Username)
	}
}

func TestJWTManager_ValidateToken_InvalidFormat(t *testing.T) {
	manager := crypto.NewJWTManager("test-secret")

	_, err := manager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token format")
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager1 := crypto.NewJWTManager("secret1")
	manager2 := crypto.NewJWTManager("secret2")

	userID := "user123"
	username := "testuser"
	expiration := 1 * time.Hour

	token, err := manager1.GenerateToken(userID, username, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for token with wrong secret")
	}
}

func TestJWTManager_ValidateToken_Expired(t *testing.T) {
	manager := crypto.NewJWTManager("test-secret")

	userID := "user123"
	username := "testuser"
	expiration := -1 * time.Hour // Expired token

	token, err := manager.GenerateToken(userID, username, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func splitToken(token string) []string {
	parts := make([]string, 0)
	start := 0
	for i, char := range token {
		if char == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	parts = append(parts, token[start:])
	return parts
}
