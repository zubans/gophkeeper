// Package tests provides tests for the crypto package.
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

func TestEncryptor_EncryptStringDecryptString(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	plaintext := "Hello, World!"
	encrypted, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt string: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatal("Encrypted string is empty")
	}

	decrypted, err := encryptor.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt string: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("Decrypted string doesn't match original. Expected: %s, Got: %s", plaintext, decrypted)
	}
}

func TestEncryptor_DifferentKeys(t *testing.T) {
	encryptor1 := crypto.NewEncryptor("key1")
	encryptor2 := crypto.NewEncryptor("key2")

	plaintext := "Hello, World!"
	encrypted, err := encryptor1.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	_, err = encryptor2.DecryptString(encrypted)
	if err == nil {
		t.Fatal("Expected decryption with different key to fail")
	}
}

func TestEncryptor_EmptyData(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	encrypted, err := encryptor.Encrypt([]byte{})
	if err != nil {
		t.Fatalf("Failed to encrypt empty data: %v", err)
	}

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt empty data: %v", err)
	}

	if len(decrypted) != 0 {
		t.Fatal("Decrypted empty data should be empty")
	}
}
