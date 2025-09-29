package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

const (
	serverURL = "http://localhost:8080"
)

var (
	testUser  = fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	testEmail = fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	testPass  = "testpass123"
)

func TestIntegration(t *testing.T) {
	// Wait for server to be ready
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}

	// Test user registration
	t.Run("UserRegistration", testUserRegistration)

	// Test user login
	t.Run("UserLogin", testUserLogin)

	// Test data operations
	t.Run("DataOperations", testDataOperations)

	// Test synchronization
	t.Run("Synchronization", testSynchronization)
}

func waitForServer(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(serverURL + "/api/v1/register")
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func testUserRegistration(t *testing.T) {
	req := models.UserRegistrationRequest{
		Username: testUser,
		Email:    testEmail,
		Password: testPass,
	}

	resp, err := makeRequest("POST", "/api/v1/register", req)
	if err != nil {
		t.Fatalf("Registration request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !apiResp.Success {
		t.Fatalf("Registration failed: %s", apiResp.Error)
	}

	// Verify response contains token
	var authResp models.AuthResponse
	dataBytes, _ := json.Marshal(apiResp.Data)
	if err := json.Unmarshal(dataBytes, &authResp); err != nil {
		t.Fatalf("Failed to unmarshal auth response: %v", err)
	}

	if authResp.Token == "" {
		t.Fatal("Token is empty in registration response")
	}

	if authResp.User.Username != testUser {
		t.Fatalf("Expected username %s, got %s", testUser, authResp.User.Username)
	}
}

func testUserLogin(t *testing.T) {
	req := models.UserLoginRequest{
		Username: testUser,
		Password: testPass,
	}

	resp, err := makeRequest("POST", "/api/v1/login", req)
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !apiResp.Success {
		t.Fatalf("Login failed: %s", apiResp.Error)
	}

	// Verify response contains token
	var authResp models.AuthResponse
	dataBytes, _ := json.Marshal(apiResp.Data)
	if err := json.Unmarshal(dataBytes, &authResp); err != nil {
		t.Fatalf("Failed to unmarshal auth response: %v", err)
	}

	if authResp.Token == "" {
		t.Fatal("Token is empty in login response")
	}

	if authResp.User.Username != testUser {
		t.Fatalf("Expected username %s, got %s", testUser, authResp.User.Username)
	}
}

func testDataOperations(t *testing.T) {
	// First login to get token
	token := loginAndGetToken(t)

	// Test creating login/password data
	loginData := models.StoredData{
		ID:       uuid.New().String(),
		Type:     models.DataTypeLoginPassword,
		Title:    "Test Website",
		Data:     []byte(`{"username":"testuser","password":"testpass","url":"https://example.com"}`),
		Metadata: "Test metadata",
		Version:  1,
	}

	resp, err := makeRequestWithAuth("POST", "/api/v1/data", loginData, token)
	if err != nil {
		t.Fatalf("Create data request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create data failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Test getting all data
	resp, err = makeRequestWithAuth("GET", "/api/v1/data", nil, token)
	if err != nil {
		t.Fatalf("Get data request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Get data failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !apiResp.Success {
		t.Fatalf("Get data failed: %s", apiResp.Error)
	}

	// Verify data was created
	var dataList []models.StoredData
	dataBytes, _ := json.Marshal(apiResp.Data)
	if err := json.Unmarshal(dataBytes, &dataList); err != nil {
		t.Fatalf("Failed to unmarshal data list: %v", err)
	}

	if len(dataList) == 0 {
		t.Fatal("No data returned")
	}

	found := false
	for _, data := range dataList {
		if data.Title == "Test Website" {
			found = true
			if data.Type != models.DataTypeLoginPassword {
				t.Fatalf("Expected type %s, got %s", models.DataTypeLoginPassword, data.Type)
			}
			break
		}
	}

	if !found {
		t.Fatal("Created data not found in response")
	}
}

func testSynchronization(t *testing.T) {
	token := loginAndGetToken(t)

	// Create sync request
	syncReq := models.DataSyncRequest{
		LastSyncAt: time.Now().Add(-1 * time.Hour),
		Data:       []models.StoredData{},
	}

	resp, err := makeRequestWithAuth("POST", "/api/v1/sync", syncReq, token)
	if err != nil {
		t.Fatalf("Sync request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Sync failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !apiResp.Success {
		t.Fatalf("Sync failed: %s", apiResp.Error)
	}

	// Verify sync response
	var syncResp models.DataSyncResponse
	dataBytes, _ := json.Marshal(apiResp.Data)
	if err := json.Unmarshal(dataBytes, &syncResp); err != nil {
		t.Fatalf("Failed to unmarshal sync response: %v", err)
	}

	if syncResp.LastSyncAt.IsZero() {
		t.Fatal("LastSyncAt is zero")
	}
}

func loginAndGetToken(t *testing.T) string {
	req := models.UserLoginRequest{
		Username: testUser,
		Password: testPass,
	}

	resp, err := makeRequest("POST", "/api/v1/login", req)
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !apiResp.Success {
		t.Fatalf("Login failed: %s", apiResp.Error)
	}

	var authResp models.AuthResponse
	dataBytes, _ := json.Marshal(apiResp.Data)
	if err := json.Unmarshal(dataBytes, &authResp); err != nil {
		t.Fatalf("Failed to unmarshal auth response: %v", err)
	}

	return authResp.Token
}

func makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, serverURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func makeRequestWithAuth(method, path string, body interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, serverURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
