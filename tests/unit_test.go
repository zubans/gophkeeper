package tests

import (
	"testing"
	"time"

	"gophkeeper/internal/crypto"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

func TestCrypto(t *testing.T) {
	t.Run("PasswordHashing", testPasswordHashing)
	t.Run("Encryption", testEncryption)
	t.Run("JWT", testJWT)
}

func testPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Test hashing
	hashed, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("Hashed password is empty")
	}

	if hashed == password {
		t.Fatal("Hashed password is the same as original")
	}

	// Test verification
	valid, err := crypto.VerifyPassword(password, hashed)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}

	if !valid {
		t.Fatal("Password verification failed")
	}

	// Test wrong password
	wrongValid, err := crypto.VerifyPassword("wrongpassword", hashed)
	if err != nil {
		t.Fatalf("VerifyPassword with wrong password failed: %v", err)
	}

	if wrongValid {
		t.Fatal("Wrong password should not be valid")
	}
}

func testEncryption(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key-32-characters-long!")

	data := []byte("test data to encrypt")

	// Test encryption
	encrypted, err := encryptor.Encrypt(data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatal("Encrypted data is empty")
	}

	if string(encrypted) == string(data) {
		t.Fatal("Encrypted data is the same as original")
	}

	// Test decryption
	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(data) {
		t.Fatalf("Decrypted data doesn't match original: got %s, want %s", string(decrypted), string(data))
	}
}

func testJWT(t *testing.T) {
	jwtManager := crypto.NewJWTManager("test-secret-key")

	userID := uuid.New().String()
	username := "testuser"
	duration := 1 * time.Hour

	// Test token generation
	token, err := jwtManager.GenerateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Test token validation
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Fatalf("Expected Username %s, got %s", username, claims.Username)
	}

	// Test invalid token
	_, err = jwtManager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Invalid token should fail validation")
	}
}

func TestModels(t *testing.T) {
	t.Run("DataTypes", testDataTypes)
	t.Run("APIResponse", testAPIResponse)
}

func testDataTypes(t *testing.T) {
	// Test DataType constants
	if models.DataTypeLoginPassword != "login_password" {
		t.Fatalf("Expected DataTypeLoginPassword to be 'login_password', got %s", models.DataTypeLoginPassword)
	}

	if models.DataTypeText != "text" {
		t.Fatalf("Expected DataTypeText to be 'text', got %s", models.DataTypeText)
	}

	if models.DataTypeBinary != "binary" {
		t.Fatalf("Expected DataTypeBinary to be 'binary', got %s", models.DataTypeBinary)
	}

	if models.DataTypeBankCard != "bank_card" {
		t.Fatalf("Expected DataTypeBankCard to be 'bank_card', got %s", models.DataTypeBankCard)
	}
}

func testAPIResponse(t *testing.T) {
	// Test success response
	data := "test data"
	successResp := models.NewSuccessResponse(data)

	if !successResp.Success {
		t.Fatal("Success response should have Success=true")
	}

	if successResp.Data != data {
		t.Fatalf("Expected data %v, got %v", data, successResp.Data)
	}

	// Test error response
	errorMsg := "test error"
	errorCode := 400
	errorResp := models.NewErrorResponse(errorMsg, errorCode)

	if errorResp.Success {
		t.Fatal("Error response should have Success=false")
	}

	if errorResp.Error != errorMsg {
		t.Fatalf("Expected error %s, got %s", errorMsg, errorResp.Error)
	}

	if errorResp.Code != errorCode {
		t.Fatalf("Expected code %d, got %d", errorCode, errorResp.Code)
	}
}

func TestStoredData(t *testing.T) {
	now := time.Now()

	testID := uuid.New().String()
	data := models.StoredData{
		ID:         testID,
		UserID:     uuid.New().String(),
		Type:       models.DataTypeLoginPassword,
		Title:      "Test Data",
		Data:       []byte("test data"),
		Metadata:   "test metadata",
		Version:    1,
		CreatedAt:  now,
		UpdatedAt:  now,
		LastSyncAt: now,
		IsDeleted:  false,
	}

	// Test basic fields
	if data.ID != testID {
		t.Fatalf("Expected ID '%s', got %s", testID, data.ID)
	}

	// UserID is now a UUID, so we just check it's not empty
	if data.UserID == "" {
		t.Fatalf("Expected UserID to be set, got empty string")
	}

	if data.Type != models.DataTypeLoginPassword {
		t.Fatalf("Expected Type %s, got %s", models.DataTypeLoginPassword, data.Type)
	}

	if data.Title != "Test Data" {
		t.Fatalf("Expected Title 'Test Data', got %s", data.Title)
	}

	if string(data.Data) != "test data" {
		t.Fatalf("Expected Data 'test data', got %s", string(data.Data))
	}

	if data.Metadata != "test metadata" {
		t.Fatalf("Expected Metadata 'test metadata', got %s", data.Metadata)
	}

	if data.Version != 1 {
		t.Fatalf("Expected Version 1, got %d", data.Version)
	}

	if data.IsDeleted {
		t.Fatal("Expected IsDeleted to be false")
	}
}
