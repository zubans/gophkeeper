package tests
import (
	"testing"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/tests/mocks"
)
func TestSyncService_SyncData(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{
		Authenticated: true,
		UserID:        "user-123",
		Token:         "mock-token",
	}
	syncService := client.NewSyncService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	err := syncService.SyncData()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
func TestSyncService_SyncData_NotAuthenticated(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{Authenticated: false} // Not authenticated
	syncService := client.NewSyncService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	err := syncService.SyncData()
	if err == nil {
		t.Fatal("Expected error for unauthenticated user, got nil")
	}
	if err.Error() != "not authenticated" {
		t.Errorf("Expected 'not authenticated' error, got %s", err.Error())
	}
}
func TestSyncService_SyncData_HTTPError(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{ShouldFail: true} // HTTP client will fail
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{
		Authenticated: true,
		UserID:        "user-123",
		Token:         "mock-token",
	}
	syncService := client.NewSyncService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	err := syncService.SyncData()
	if err == nil {
		t.Fatal("Expected error when HTTP client fails, got nil")
	}
	if err.Error() != "sync failed: unauthorized" {
		t.Errorf("Expected 'sync failed: unauthorized' error, got %s", err.Error())
	}
}
