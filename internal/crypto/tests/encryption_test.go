package tests

import (
	"gophkeeper/internal/crypto"
	"testing"
)

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")
	plaintext := "Hello, World!"
	encrypted, err := encryptor.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	if len(encrypted) == 0 {
		t.Fatal("Encrypted data is empty")
	}
	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	if string(decrypted) != plaintext {
		t.Fatalf("Decrypted data doesn't match original. Expected: %s, Got: %s", plaintext, string(decrypted))
	}
}
func TestEncryptor_DifferentKeys(t *testing.T) {
	encryptor1 := crypto.NewEncryptor("key1")
	encryptor2 := crypto.NewEncryptor("key2")
	plaintext := "Hello, World!"
	encrypted, err := encryptor1.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	_, err = encryptor2.Decrypt(encrypted)
	if err == nil {
		t.Fatal("Expected decryption to fail with different key")
	}
}
func TestEncryptor_EmptyData(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")
	encrypted, err := encryptor.Encrypt([]byte(""))
	if err != nil {
		t.Fatalf("Failed to encrypt empty data: %v", err)
	}
	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt empty data: %v", err)
	}
	if string(decrypted) != "" {
		t.Fatalf("Expected empty string, got: %s", string(decrypted))
	}
}
func TestEncryptor_InvalidCiphertext(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")
	_, err := encryptor.Decrypt([]byte("invalid"))
	if err == nil {
		t.Fatal("Expected decryption to fail with invalid ciphertext")
	}
}
