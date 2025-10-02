package tests
import (
	"gophkeeper/internal/crypto"
	"testing"
)
func TestHashPassword(t *testing.T) {
	password := "test-password"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if len(hash) == 0 {
		t.Fatal("Hash is empty")
	}
	hash2, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}
	if hash == hash2 {
		t.Fatal("Password hashes should be different each time")
	}
}
func TestVerifyPassword(t *testing.T) {
	password := "test-password"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	valid, err := crypto.VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if !valid {
		t.Fatal("Correct password should be valid")
	}
	valid, err = crypto.VerifyPassword("wrong-password", hash)
	if err != nil {
		t.Fatalf("Failed to verify wrong password: %v", err)
	}
	if valid {
		t.Fatal("Wrong password should not be valid")
	}
}
func TestVerifyPassword_InvalidHash(t *testing.T) {
	_, err := crypto.VerifyPassword("password", "invalid-hash")
	if err == nil {
		t.Fatal("Expected error for invalid hash")
	}
}
func TestVerifyPassword_EmptyPassword(t *testing.T) {
	hash, err := crypto.HashPassword("")
	if err != nil {
		t.Fatalf("Failed to hash empty password: %v", err)
	}
	valid, err := crypto.VerifyPassword("", hash)
	if err != nil {
		t.Fatalf("Failed to verify empty password: %v", err)
	}
	if !valid {
		t.Fatal("Empty password should be valid")
	}
}
