package tests
import (
	"gophkeeper/internal/models"
	"testing"
)
func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"message": "test"}
	response := models.NewSuccessResponse(data)
	if !response.Success {
		t.Fatal("Success should be true")
	}
	if response.Data == nil {
		t.Fatal("Data should not be nil")
	}
	if response.Error != "" {
		t.Fatal("Error should be empty")
	}
}
func TestNewErrorResponse(t *testing.T) {
	message := "test error"
	code := 400
	response := models.NewErrorResponse(message, code)
	if response.Success {
		t.Fatal("Success should be false")
	}
	if response.Error != message {
		t.Fatalf("Error message mismatch. Expected: %s, Got: %s", message, response.Error)
	}
	if response.Code != code {
		t.Fatalf("Error code mismatch. Expected: %d, Got: %d", code, response.Code)
	}
}
func TestNewSuccessResponse_NilData(t *testing.T) {
	response := models.NewSuccessResponse(nil)
	if !response.Success {
		t.Fatal("Success should be true")
	}
	if response.Data != nil {
		t.Fatal("Data should be nil")
	}
}
