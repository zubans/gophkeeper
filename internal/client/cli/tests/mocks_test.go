package cli

import (
	"testing"

	"gophkeeper/internal/client"
	"gophkeeper/internal/client/cli"
	"gophkeeper/internal/models"
)

// MockClient is a mock implementation of the ClientInterface
type MockClient struct {
	RegisterFunc    func(username, email, password string) error
	LoginFunc       func(username, password string) error
	AddDataFunc     func(dataType, title string, data []string) error
	GetDataFunc     func(id string) error
	DeleteDataFunc  func(id string) error
	SyncDataFunc    func() error
	ShowHistoryFunc func(id string) error
	ListDataFunc    func() error
	GetDataListFunc func() ([]models.StoredData, error)
}

func (m *MockClient) Register(username, email, password string) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(username, email, password)
	}
	return nil
}
func (m *MockClient) Login(username, password string) error {
	if m.LoginFunc != nil {
		return m.LoginFunc(username, password)
	}
	return nil
}
func (m *MockClient) AddData(dataType, title string, data []string) error {
	if m.AddDataFunc != nil {
		return m.AddDataFunc(dataType, title, data)
	}
	return nil
}
func (m *MockClient) GetData(id string) error {
	if m.GetDataFunc != nil {
		return m.GetDataFunc(id)
	}
	return nil
}
func (m *MockClient) DeleteData(id string) error {
	if m.DeleteDataFunc != nil {
		return m.DeleteDataFunc(id)
	}
	return nil
}
func (m *MockClient) SyncData() error {
	if m.SyncDataFunc != nil {
		return m.SyncDataFunc()
	}
	return nil
}
func (m *MockClient) ShowHistory(id string) error {
	if m.ShowHistoryFunc != nil {
		return m.ShowHistoryFunc(id)
	}
	return nil
}
func (m *MockClient) ListData() error {
	if m.ListDataFunc != nil {
		return m.ListDataFunc()
	}
	return nil
}
func (m *MockClient) GetDataList() ([]models.StoredData, error) {
	if m.GetDataListFunc != nil {
		return m.GetDataListFunc()
	}
	return []models.StoredData{}, nil
}

// Execute method coverage using mocks
func TestExecuteMethods(t *testing.T) {
	t.Run("RegisterCommand_Execute", func(t *testing.T) {
		cmd := &cli.RegisterCommand{Username: "testuser", Email: "test@example.com", Password: "password123"}
		// Create a mock client.Client - this will fail validation but that's expected
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("LoginCommand_Execute", func(t *testing.T) {
		cmd := &cli.LoginCommand{Username: "u", Password: "p"}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("AddCommand_Execute", func(t *testing.T) {
		cmd := &cli.AddCommand{DataType: "text", Title: "T", Data: []string{"content"}}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("GetCommand_Execute", func(t *testing.T) {
		cmd := &cli.GetCommand{ID: "id"}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("DeleteCommand_Execute", func(t *testing.T) {
		cmd := &cli.DeleteCommand{ID: "id"}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("SyncCommand_Execute", func(t *testing.T) {
		cmd := &cli.SyncCommand{}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("HistoryCommand_Execute", func(t *testing.T) {
		cmd := &cli.HistoryCommand{ID: "id"}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("ListCommand_Execute", func(t *testing.T) {
		cmd := &cli.ListCommand{}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// We expect validation to pass but execution to fail due to nil client
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("HelpCommand_Execute", func(t *testing.T) {
		cmd := &cli.HelpCommand{}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// Help should work without client
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("VersionCommand_Execute", func(t *testing.T) {
		cmd := &cli.VersionCommand{}
		mockClient := &client.Client{}
		err := cmd.Execute(mockClient)
		// Version should work without client
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
