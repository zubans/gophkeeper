package tests
import (
	"os"
	"path/filepath"
	"testing"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/tests/mocks"
)
func TestClient_Integration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gophkeeper_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	dbPath := filepath.Join(tempDir, "test.db")
	storage, err := client.NewClientStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	tokenManager := client.NewTokenManager(tempDir)
	authService := client.NewAuthService(mockHTTP, tokenManager)
	dataService := client.NewDataService(storage, mockHTTP, mockEncryptor, authService)
	syncService := client.NewSyncService(storage, mockHTTP, mockEncryptor, authService)
	t.Run("Authentication", func(t *testing.T) {
		_, err := authService.Register("testuser", "test@example.com", "password123")
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}
		if !authService.IsAuthenticated() {
			t.Error("Expected to be authenticated after registration")
		}
		err = authService.Logout()
		if err != nil {
			t.Fatalf("Logout failed: %v", err)
		}
		if authService.IsAuthenticated() {
			t.Error("Expected to be unauthenticated after logout")
		}
		_, err = authService.Login("testuser", "password123")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if !authService.IsAuthenticated() {
			t.Error("Expected to be authenticated after login")
		}
	})
	t.Run("DataManagement", func(t *testing.T) {
		err := dataService.AddData("text", "Test Note", []string{"This is a test note"})
		if err != nil {
			t.Fatalf("Add data failed: %v", err)
		}
		dataList, err := dataService.GetDataList()
		if err != nil {
			t.Fatalf("Get data list failed: %v", err)
		}
		if len(dataList) != 1 {
			t.Errorf("Expected 1 data item, got %d", len(dataList))
		}
		dataID := dataList[0].ID
		data, err := dataService.GetData(dataID)
		if err != nil {
			t.Fatalf("Get data failed: %v", err)
		}
		if data.Title != "Test Note" {
			t.Errorf("Expected title 'Test Note', got %s", data.Title)
		}
		err = dataService.DeleteData(dataID)
		if err != nil {
			t.Fatalf("Delete data failed: %v", err)
		}
		dataListAfter, err := dataService.GetDataList()
		if err != nil {
			t.Fatalf("Get data list after deletion failed: %v", err)
		}
		if len(dataListAfter) != 0 {
			t.Errorf("Expected 0 data items after deletion, got %d", len(dataListAfter))
		}
	})
	t.Run("Synchronization", func(t *testing.T) {
		err := dataService.AddData("text", "Sync Test", []string{"Data to sync"})
		if err != nil {
			t.Fatalf("Add data failed: %v", err)
		}
		err = syncService.SyncData()
		if err != nil {
			t.Fatalf("Sync failed: %v", err)
		}
	})
}
func TestClient_FullIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gophkeeper_full_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockHTTP := &mocks.MockHTTPClient{}
	tokenManager := client.NewTokenManager(tempDir)
	mockEncryptor := &mocks.MockEncryptor{}
	dbPath := filepath.Join(tempDir, "test.db")
	storage, err := client.NewClientStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()
	authService := client.NewAuthService(mockHTTP, tokenManager)
	dataService := client.NewDataService(storage, mockHTTP, mockEncryptor, authService)
	syncService := client.NewSyncService(storage, mockHTTP, mockEncryptor, authService)
	t.Run("FullWorkflow", func(t *testing.T) {
		_, err := authService.Register("fulltest", "full@test.com", "password123")
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}
		testCases := []struct {
			dataType string
			title    string
			data     []string
		}{
			{"text", "My Secret Note", []string{"This is a secret note"}},
			{"login_password", "Gmail Account", []string{"user@gmail.com", "password123", "gmail.com", "Personal email"}},
			{"bank_card", "Credit Card", []string{"1234567890123456", "12/25", "123", "John Doe", "Bank of Test", "Main card"}},
		}
		for _, tc := range testCases {
			err := dataService.AddData(tc.dataType, tc.title, tc.data)
			if err != nil {
				t.Fatalf("Failed to add %s data: %v", tc.dataType, err)
			}
		}
		dataList, err := dataService.GetDataList()
		if err != nil {
			t.Fatalf("Failed to get data list: %v", err)
		}
		if len(dataList) != 3 {
			t.Errorf("Expected 3 data items, got %d", len(dataList))
		}
		err = syncService.SyncData()
		if err != nil {
			t.Fatalf("Sync failed: %v", err)
		}
		err = authService.Logout()
		if err != nil {
			t.Fatalf("Logout failed: %v", err)
		}
		_, err = authService.Login("fulltest", "password123")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		dataListAfterLogin, err := dataService.GetDataList()
		if err != nil {
			t.Fatalf("Failed to get data list after login: %v", err)
		}
		if len(dataListAfterLogin) != 3 {
			t.Errorf("Expected 3 data items after login, got %d", len(dataListAfterLogin))
		}
	})
}
func TestClient_ErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gophkeeper_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockHTTP := &mocks.MockHTTPClient{ShouldFail: true}
	tokenManager := client.NewTokenManager(tempDir)
	mockEncryptor := &mocks.MockEncryptor{}
	dbPath := filepath.Join(tempDir, "test.db")
	storage, err := client.NewClientStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()
	authService := client.NewAuthService(mockHTTP, tokenManager)
	dataService := client.NewDataService(storage, mockHTTP, mockEncryptor, authService)
	syncService := client.NewSyncService(storage, mockHTTP, mockEncryptor, authService)
	t.Run("AuthenticationErrors", func(t *testing.T) {
		_, err := authService.Register("testuser", "test@example.com", "password123")
		if err == nil {
			t.Error("Expected registration to fail with failing HTTP client")
		}
		_, err = authService.Login("testuser", "password123")
		if err == nil {
			t.Error("Expected login to fail with failing HTTP client")
		}
	})
	t.Run("DataOperationErrors", func(t *testing.T) {
		err := dataService.AddData("text", "Test", []string{"data"})
		if err == nil || err.Error() != "not authenticated" {
			t.Errorf("Expected 'not authenticated' error, got %v", err)
		}
		_, err = dataService.GetDataList()
		if err == nil || err.Error() != "not authenticated" {
			t.Errorf("Expected 'not authenticated' error, got %v", err)
		}
		err = syncService.SyncData()
		if err == nil || err.Error() != "not authenticated" {
			t.Errorf("Expected 'not authenticated' error, got %v", err)
		}
	})
}
