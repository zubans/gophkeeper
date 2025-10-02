package mocks
import (
	"time"
	"gophkeeper/internal/models"
)
type MockStorage struct {
	data map[string]*models.StoredData
}
func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]*models.StoredData),
	}
}
func (m *MockStorage) SaveData(data *models.StoredData) error {
	m.data[data.ID] = data
	return nil
}
func (m *MockStorage) GetData(id string) (*models.StoredData, error) {
	if data, exists := m.data[id]; exists {
		return data, nil
	}
	return nil, nil
}
func (m *MockStorage) GetAllData(userID string) ([]models.StoredData, error) {
	var result []models.StoredData
	for _, data := range m.data {
		if data.UserID == userID {
			result = append(result, *data)
		}
	}
	return result, nil
}
func (m *MockStorage) GetDataSince(userID string, since time.Time) ([]models.StoredData, error) {
	var result []models.StoredData
	for _, data := range m.data {
		if data.UserID == userID && data.UpdatedAt.After(since) {
			result = append(result, *data)
		}
	}
	return result, nil
}
func (m *MockStorage) DeleteData(id string) error {
	delete(m.data, id)
	return nil
}
func (m *MockStorage) GetDataHistory(id string) ([]models.DataHistory, error) {
	return []models.DataHistory{}, nil
}
func (m *MockStorage) GetLastSyncTime(userID string) (time.Time, error) {
	return time.Time{}, nil
}
func (m *MockStorage) UpdateLastSyncTime(userID string, t time.Time) error {
	return nil
}
func (m *MockStorage) Close() error {
	return nil
}
type MockHTTPClient struct {
	ShouldFail bool
}
func (m *MockHTTPClient) Register(req *models.UserRegistrationRequest) (*models.AuthResponse, error) {
	if m.ShouldFail {
		return nil, models.ErrUserAlreadyExists
	}
	return &models.AuthResponse{
		Token: "mock-token",
		User: models.User{
			ID:       "user-123",
			Username: req.Username,
			Email:    req.Email,
		},
	}, nil
}
func (m *MockHTTPClient) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	if m.ShouldFail {
		return nil, models.ErrInvalidCredentials
	}
	return &models.AuthResponse{
		Token: "mock-token",
		User: models.User{
			ID:       "user-123",
			Username: req.Username,
		},
	}, nil
}
func (m *MockHTTPClient) AddData(data *models.StoredData, token string) error {
	if m.ShouldFail {
		return models.ErrUnauthorized
	}
	return nil
}
func (m *MockHTTPClient) DeleteData(id, token string) error {
	if m.ShouldFail {
		return models.ErrUnauthorized
	}
	return nil
}
func (m *MockHTTPClient) SyncData(req *models.DataSyncRequest, token string) (*models.DataSyncResponse, error) {
	if m.ShouldFail {
		return nil, models.ErrUnauthorized
	}
	return &models.DataSyncResponse{
		Data:       []models.StoredData{},
		LastSyncAt: time.Now(),
	}, nil
}
type MockEncryptor struct{}
func (m *MockEncryptor) Encrypt(data []byte) ([]byte, error) {
	return append([]byte("encrypted:"), data...), nil
}
func (m *MockEncryptor) Decrypt(data []byte) ([]byte, error) {
	if len(data) > 10 && string(data[:10]) == "encrypted:" {
		return data[10:], nil
	}
	return data, nil
}
type MockTokenManager struct {
	token string
}
func (m *MockTokenManager) SaveToken(token string) error {
	m.token = token
	return nil
}
func (m *MockTokenManager) LoadToken() (string, error) {
	return m.token, nil
}
func (m *MockTokenManager) ClearToken() error {
	m.token = ""
	return nil
}
type MockAuthService struct {
	Authenticated bool
	UserID        string
	Token         string
}
func (m *MockAuthService) Register(username, email, password string) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *MockAuthService) Login(username, password string) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *MockAuthService) IsAuthenticated() bool {
	return m.Authenticated
}
func (m *MockAuthService) GetToken() string {
	return m.Token
}
func (m *MockAuthService) GetUserID() string {
	return m.UserID
}
func (m *MockAuthService) Logout() error {
	m.Authenticated = false
	m.Token = ""
	m.UserID = ""
	return nil
}
