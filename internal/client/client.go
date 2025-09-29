// Package client implements the GophKeeper CLI client.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gophkeeper/internal/crypto"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// Client represents the GophKeeper client.
type Client struct {
	serverURL  string
	configDir  string
	httpClient *http.Client
	encryptor  *crypto.Encryptor
	storage    *ClientStorage
	token      string
	userID     string
}

// NewClient creates a new client instance.
func NewClient(serverURL, configDir, encryptionKey string) (*Client, error) {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Initialize SQLite storage
	dbPath := filepath.Join(configDir, "data.db")
	storage, err := NewClientStorage(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Load existing token if available
	token, err := loadToken(configDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	return &Client{
		serverURL:  serverURL,
		configDir:  configDir,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		encryptor:  crypto.NewEncryptor(encryptionKey),
		storage:    storage,
		token:      token,
	}, nil
}

// Register registers a new user.
func (c *Client) Register(username, email, password string) error {
	req := models.UserRegistrationRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	var response models.AuthResponse
	if err := c.makeRequest("POST", "/api/v1/register", req, &response); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	c.token = response.Token
	c.userID = response.User.ID

	// Save token
	if err := c.saveToken(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// Login authenticates a user.
func (c *Client) Login(username, password string) error {
	req := models.UserLoginRequest{
		Username: username,
		Password: password,
	}

	var response models.AuthResponse
	if err := c.makeRequest("POST", "/api/v1/login", req, &response); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	c.token = response.Token
	c.userID = response.User.ID

	// Save token
	if err := c.saveToken(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// AddData adds new data to the server.
func (c *Client) AddData(dataType, title string, data []string) error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	var storedData models.StoredData
	storedData.ID = generateID()
	storedData.UserID = c.userID
	storedData.Type = models.DataType(dataType)
	storedData.Title = title
	storedData.Version = 1

	// Process data based on type
	switch models.DataType(dataType) {
	case models.DataTypeLoginPassword:
		if len(data) < 2 {
			return fmt.Errorf("login/password data requires at least login and password")
		}
		loginData := models.LoginPasswordData{
			Login:    data[0],
			Password: data[1],
		}
		if len(data) > 2 {
			loginData.Website = data[2]
		}
		if len(data) > 3 {
			loginData.Notes = data[3]
		}
		jsonData, _ := json.Marshal(loginData)
		storedData.Data = jsonData

	case models.DataTypeText:
		if len(data) == 0 {
			return fmt.Errorf("text data is required")
		}
		storedData.Data = []byte(data[0])

	case models.DataTypeBankCard:
		if len(data) < 4 {
			return fmt.Errorf("bank card data requires card number, expiry, CVV, and cardholder")
		}
		cardData := models.BankCardData{
			CardNumber: data[0],
			ExpiryDate: data[1],
			CVV:        data[2],
			Cardholder: data[3],
		}
		if len(data) > 4 {
			cardData.Bank = data[4]
		}
		if len(data) > 5 {
			cardData.Notes = data[5]
		}
		jsonData, _ := json.Marshal(cardData)
		storedData.Data = jsonData

	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}

	// Save locally first
	if err := c.storage.SaveData(&storedData); err != nil {
		return fmt.Errorf("failed to save data locally: %w", err)
	}

	// Send to server
	if err := c.makeRequest("POST", "/api/v1/data", storedData, &storedData); err != nil {
		return fmt.Errorf("failed to add data to server: %w", err)
	}

	return nil
}

// ListData lists all user data.
func (c *Client) ListData() error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	dataList, err := c.storage.GetAllData(c.userID)
	if err != nil {
		return fmt.Errorf("failed to list data: %w", err)
	}

	if len(dataList) == 0 {
		fmt.Println("No data found.")
		return nil
	}

	fmt.Printf("%-36s %-20s %-50s %-20s\n", "ID", "Type", "Title", "Updated")
	fmt.Println(strings.Repeat("-", 130))

	for _, data := range dataList {
		fmt.Printf("%-36s %-20s %-50s %-20s\n",
			data.ID, data.Type, data.Title, data.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// GetDataList returns all user data as a slice.
func (c *Client) GetDataList() ([]models.StoredData, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	dataList, err := c.storage.GetAllData(c.userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	return dataList, nil
}

// GetData retrieves specific data by ID.
func (c *Client) GetData(id string) error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	data, err := c.storage.GetData(id)
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}

	fmt.Printf("ID: %s\n", data.ID)
	fmt.Printf("Type: %s\n", data.Type)
	fmt.Printf("Title: %s\n", data.Title)
	fmt.Printf("Data: %s\n", string(data.Data))
	if data.Metadata != "" {
		fmt.Printf("Metadata: %s\n", data.Metadata)
	}
	fmt.Printf("Updated: %s\n", data.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// DeleteData deletes data by ID.
func (c *Client) DeleteData(id string) error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	// Mark as deleted locally
	if err := c.storage.DeleteData(id); err != nil {
		return fmt.Errorf("failed to delete data locally: %w", err)
	}

	// Send delete request to server
	url := fmt.Sprintf("%s/api/v1/data?id=%s", c.serverURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s", string(body))
	}

	return nil
}

// SyncData synchronizes data with the server.
func (c *Client) SyncData() error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	// Get last sync time
	lastSyncTime, err := c.storage.GetLastSyncTime(c.userID)
	if err != nil {
		return fmt.Errorf("failed to get last sync time: %w", err)
	}

	// Get local data since last sync
	localData, err := c.storage.GetDataSince(c.userID, lastSyncTime)
	if err != nil {
		return fmt.Errorf("failed to get local data: %w", err)
	}

	// Encrypt local data before sending
	for i := range localData {
		if err := c.encryptData(&localData[i]); err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}
	}

	req := models.DataSyncRequest{
		LastSyncAt: lastSyncTime,
		Data:       localData,
	}

	var response models.DataSyncResponse
	if err := c.makeRequest("POST", "/api/v1/sync", req, &response); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	// Decrypt server data before saving
	for i := range response.Data {
		if err := c.decryptData(&response.Data[i]); err != nil {
			return fmt.Errorf("failed to decrypt server data: %w", err)
		}
	}

	// Save server data locally
	for _, data := range response.Data {
		if err := c.storage.SaveData(&data); err != nil {
			return fmt.Errorf("failed to save server data: %w", err)
		}
	}

	// Update last sync time
	if err := c.storage.UpdateLastSyncTime(c.userID, response.LastSyncAt); err != nil {
		return fmt.Errorf("failed to update last sync time: %w", err)
	}

	return nil
}

// ShowHistory shows the history of a specific data item.
func (c *Client) ShowHistory(id string) error {
	if c.token == "" {
		return fmt.Errorf("not authenticated")
	}

	history, err := c.storage.GetDataHistory(id)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(history) == 0 {
		fmt.Println("No history found for this data item.")
		return nil
	}

	fmt.Printf("History for data ID: %s\n", id)
	fmt.Printf("%-8s %-20s %-50s %-20s\n", "Version", "Type", "Title", "Updated")
	fmt.Println(strings.Repeat("-", 100))

	for _, h := range history {
		status := "Active"
		if h.IsDeleted {
			status = "Deleted"
		}
		fmt.Printf("%-8d %-20s %-50s %-20s %s\n",
			h.Version, h.Type, h.Title, h.UpdatedAt.Format("2006-01-02 15:04:05"), status)
	}

	return nil
}

// generateID generates a unique UUID v4.
func generateID() string {
	return uuid.New().String()
}

// encryptData encrypts the data field of StoredData.
func (c *Client) encryptData(data *models.StoredData) error {
	encrypted, err := c.encryptor.Encrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	data.Data = encrypted
	return nil
}

// decryptData decrypts the data field of StoredData.
func (c *Client) decryptData(data *models.StoredData) error {
	decrypted, err := c.encryptor.Decrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}
	data.Data = decrypted
	return nil
}

// makeRequest makes an HTTP request to the server.
func (c *Client) makeRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.serverURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp models.ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			return fmt.Errorf("request failed: %s", errorResp.Error)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		// Try wrapped APIResponse first
		var wrap models.APIResponse
		if err := json.Unmarshal(respBody, &wrap); err == nil {
			if !wrap.Success {
				if wrap.Error != "" {
					return fmt.Errorf("request failed: %s", wrap.Error)
				}
				return fmt.Errorf("request failed")
			}
			if wrap.Data != nil {
				// Convert data to JSON and unmarshal into result
				dataBytes, err := json.Marshal(wrap.Data)
				if err != nil {
					return fmt.Errorf("failed to marshal wrapped data: %w", err)
				}
				if err := json.Unmarshal(dataBytes, result); err != nil {
					return fmt.Errorf("failed to unmarshal wrapped data: %w", err)
				}
				return nil
			}
		}
		// Fallback: direct unmarshal
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// saveToken saves the authentication token to disk.
func (c *Client) saveToken() error {
	tokenFile := filepath.Join(c.configDir, "token")
	return os.WriteFile(tokenFile, []byte(c.token), 0600)
}

// loadToken loads the authentication token from disk.
func loadToken(configDir string) (string, error) {
	tokenFile := filepath.Join(configDir, "token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// loadLocalData loads locally stored data.
func (c *Client) loadLocalData() ([]models.StoredData, error) {
	dataFile := filepath.Join(c.configDir, "data.json")
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.StoredData{}, nil
		}
		return nil, err
	}

	var dataList []models.StoredData
	if err := json.Unmarshal(data, &dataList); err != nil {
		return nil, err
	}

	return dataList, nil
}

// saveLocalData saves data locally.
func (c *Client) saveLocalData(dataList []models.StoredData) error {
	dataFile := filepath.Join(c.configDir, "data.json")
	jsonData, err := json.Marshal(dataList)
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, jsonData, 0600)
}

// getLastSyncTime gets the last synchronization time.
func (c *Client) getLastSyncTime() time.Time {
	syncFile := filepath.Join(c.configDir, "last_sync")
	data, err := os.ReadFile(syncFile)
	if err != nil {
		return time.Time{}
	}

	var lastSync time.Time
	if err := json.Unmarshal(data, &lastSync); err != nil {
		return time.Time{}
	}

	return lastSync
}

// updateLastSyncTime updates the last synchronization time.
func (c *Client) updateLastSyncTime(t time.Time) error {
	syncFile := filepath.Join(c.configDir, "last_sync")
	jsonData, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(syncFile, jsonData, 0600)
}
