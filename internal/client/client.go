package client

import (
	"fmt"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	authService AuthService
	dataService DataService
	syncService SyncService
}

func NewClient(serverURL, configDir, encryptionKey string) (*Client, error) {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	dbPath := filepath.Join(configDir, "data.db")
	storage, err := NewClientStorage(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	httpClient := NewHTTPClient(serverURL)
	tokenManager := NewTokenManager(configDir)
	encryptor := crypto.NewEncryptor(encryptionKey)
	authService := NewAuthService(httpClient, tokenManager)
	dataService := NewDataService(storage, httpClient, encryptor, authService)
	syncService := NewSyncService(storage, httpClient, encryptor, authService)
	return &Client{
		authService: authService,
		dataService: dataService,
		syncService: syncService,
	}, nil
}
func (c *Client) Register(username, email, password string) error {
	_, err := c.authService.Register(username, email, password)
	return err
}
func (c *Client) Login(username, password string) error {
	_, err := c.authService.Login(username, password)
	return err
}
func (c *Client) AddData(dataType, title string, data []string) error {
	return c.dataService.AddData(dataType, title, data)
}
func (c *Client) ListData() error {
	dataList, err := c.dataService.GetDataList()
	if err != nil {
		return err
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
func (c *Client) GetDataList() ([]models.StoredData, error) {
	return c.dataService.GetDataList()
}
func (c *Client) GetData(id string) error {
	data, err := c.dataService.GetData(id)
	if err != nil {
		return err
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
func (c *Client) DeleteData(id string) error {
	return c.dataService.DeleteData(id)
}
func (c *Client) SyncData() error {
	return c.syncService.SyncData()
}
func (c *Client) ShowHistory(id string) error {
	return c.dataService.ShowHistory(id)
}
func (c *Client) IsAuthenticated() bool {
	return c.authService.IsAuthenticated()
}
func (c *Client) Logout() error {
	return c.authService.Logout()
}
func GenerateID() string {
	return uuid.New().String()
}
