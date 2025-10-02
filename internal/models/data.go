package models
import (
	"time"
)
type DataType string
const (
	DataTypeLoginPassword DataType = "login_password"
	DataTypeText          DataType = "text"
	DataTypeBinary        DataType = "binary"
	DataTypeBankCard      DataType = "bank_card"
)
type StoredData struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Type       DataType  `json:"type" db:"type"`
	Title      string    `json:"title" db:"title"`
	Data       []byte    `json:"data" db:"data"`
	Metadata   string    `json:"metadata" db:"metadata"`
	Version    int       `json:"version" db:"version"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	LastSyncAt time.Time `json:"last_sync_at" db:"last_sync_at"`
	IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
}
type DataHistory struct {
	ID        string    `json:"id" db:"id"`
	DataID    string    `json:"data_id" db:"data_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Type      DataType  `json:"type" db:"type"`
	Title     string    `json:"title" db:"title"`
	Data      []byte    `json:"data" db:"data"`
	Metadata  string    `json:"metadata" db:"metadata"`
	Version   int       `json:"version" db:"version"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
type LoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Website  string `json:"website,omitempty"`
	Notes    string `json:"notes,omitempty"`
}
type BankCardData struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
	Cardholder string `json:"cardholder"`
	Bank       string `json:"bank,omitempty"`
	Notes      string `json:"notes,omitempty"`
}
type DataSyncRequest struct {
	LastSyncAt time.Time    `json:"last_sync_at"`
	Data       []StoredData `json:"data"`
}
type DataSyncResponse struct {
	Data       []StoredData `json:"data"`
	LastSyncAt time.Time    `json:"last_sync_at"`
	Conflicts  []Conflict   `json:"conflicts,omitempty"`
}
type Conflict struct {
	LocalData  StoredData `json:"local_data"`
	ServerData StoredData `json:"server_data"`
	Reason     string     `json:"reason"`
}
