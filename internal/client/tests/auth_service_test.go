package tests
import (
	"testing"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/tests/mocks"
)
func TestAuthService_Register(t *testing.T) {
	mockHTTP := &mocks.MockHTTPClient{}
	mockToken := &mocks.MockTokenManager{}
	authService := client.NewAuthService(mockHTTP, mockToken)
	response, err := authService.Register("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.Token != "mock-token" {
		t.Errorf("Expected token 'mock-token', got %s", response.Token)
	}
	if !authService.IsAuthenticated() {
		t.Error("Expected user to be authenticated after registration")
	}
	savedToken, _ := mockToken.LoadToken()
	if savedToken != "mock-token" {
		t.Errorf("Expected saved token 'mock-token', got %s", savedToken)
	}
}
func TestAuthService_Login(t *testing.T) {
	mockHTTP := &mocks.MockHTTPClient{}
	mockToken := &mocks.MockTokenManager{}
	authService := client.NewAuthService(mockHTTP, mockToken)
	response, err := authService.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.Token != "mock-token" {
		t.Errorf("Expected token 'mock-token', got %s", response.Token)
	}
	if !authService.IsAuthenticated() {
		t.Error("Expected user to be authenticated after login")
	}
}
func TestAuthService_Logout(t *testing.T) {
	mockHTTP := &mocks.MockHTTPClient{}
	mockToken := &mocks.MockTokenManager{}
	authService := client.NewAuthService(mockHTTP, mockToken)
	_, _ = authService.Login("testuser", "password123")
	err := authService.Logout()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if authService.IsAuthenticated() {
		t.Error("Expected user to be unauthenticated after logout")
	}
	savedToken, _ := mockToken.LoadToken()
	if savedToken != "" {
		t.Errorf("Expected empty token after logout, got %s", savedToken)
	}
}
