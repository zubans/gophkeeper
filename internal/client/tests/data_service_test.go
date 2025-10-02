package tests
import (
	"testing"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/tests/mocks"
)
func TestDataService_AddData(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{
		Authenticated: true,
		UserID:        "user-123",
		Token:         "mock-token",
	}
	dataService := client.NewDataService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	err := dataService.AddData("text", "My Note", []string{"This is a test note"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	allData, _ := mockStorage.GetAllData("user-123")
	if len(allData) != 1 {
		t.Errorf("Expected 1 data item, got %d", len(allData))
	}
	if allData[0].Title != "My Note" {
		t.Errorf("Expected title 'My Note', got %s", allData[0].Title)
	}
}
func TestDataService_AddData_NotAuthenticated(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{Authenticated: false} // Not authenticated
	dataService := client.NewDataService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	err := dataService.AddData("text", "My Note", []string{"This is a test note"})
	if err == nil {
		t.Fatal("Expected error for unauthenticated user, got nil")
	}
	if err.Error() != "not authenticated" {
		t.Errorf("Expected 'not authenticated' error, got %s", err.Error())
	}
}
func TestDataService_GetDataList(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{
		Authenticated: true,
		UserID:        "user-123",
		Token:         "mock-token",
	}
	dataService := client.NewDataService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	_ = dataService.AddData("text", "Note 1", []string{"Content 1"})
	_ = dataService.AddData("text", "Note 2", []string{"Content 2"})
	dataList, err := dataService.GetDataList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(dataList) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(dataList))
	}
}
func TestDataService_DeleteData(t *testing.T) {
	mockStorage := mocks.NewMockStorage()
	mockHTTP := &mocks.MockHTTPClient{}
	mockEncryptor := &mocks.MockEncryptor{}
	mockAuth := &mocks.MockAuthService{
		Authenticated: true,
		UserID:        "user-123",
		Token:         "mock-token",
	}
	dataService := client.NewDataService(mockStorage, mockHTTP, mockEncryptor, mockAuth)
	_ = dataService.AddData("text", "Note to Delete", []string{"Content"})
	dataList, _ := dataService.GetDataList()
	dataID := dataList[0].ID
	err := dataService.DeleteData(dataID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	dataListAfter, _ := dataService.GetDataList()
	if len(dataListAfter) != 0 {
		t.Errorf("Expected 0 data items after deletion, got %d", len(dataListAfter))
	}
}
