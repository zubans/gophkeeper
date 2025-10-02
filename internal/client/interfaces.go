package client
import (
	"time"
	"gophkeeper/internal/models"
)
type Storage interface {
	SaveData(data *models.StoredData) error
	GetData(id string) (*models.StoredData, error)
	GetAllData(userID string) ([]models.StoredData, error)
	GetDataSince(userID string, since time.Time) ([]models.StoredData, error)
	DeleteData(id string) error
	GetDataHistory(id string) ([]models.DataHistory, error)
	GetLastSyncTime(userID string) (time.Time, error)
	UpdateLastSyncTime(userID string, t time.Time) error
	Close() error
}
type HTTPClient interface {
	Register(req *models.UserRegistrationRequest) (*models.AuthResponse, error)
	Login(req *models.UserLoginRequest) (*models.AuthResponse, error)
	AddData(data *models.StoredData, token string) error
	DeleteData(id, token string) error
	SyncData(req *models.DataSyncRequest, token string) (*models.DataSyncResponse, error)
}
type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}
type TokenManager interface {
	SaveToken(token string) error
	LoadToken() (string, error)
	ClearToken() error
}
type AuthService interface {
	Register(username, email, password string) (*models.AuthResponse, error)
	Login(username, password string) (*models.AuthResponse, error)
	IsAuthenticated() bool
	GetToken() string
	GetUserID() string
	Logout() error
}
type DataService interface {
	AddData(dataType, title string, data []string) error
	GetData(id string) (*models.StoredData, error)
	GetDataList() ([]models.StoredData, error)
	DeleteData(id string) error
	ShowHistory(id string) error
}
type SyncService interface {
	SyncData() error
}
