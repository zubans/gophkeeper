package client

import (
	"database/sql"
	"fmt"
	cm "gophkeeper/internal/client/migrations"
	"gophkeeper/internal/models"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type ClientStorage struct {
	db *sql.DB
}

func NewClientStorage(dbPath string) (*ClientStorage, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	storage := &ClientStorage{db: db}
	goose.SetBaseFS(cm.ClientMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db, "client"); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	return storage, nil
}
func (s *ClientStorage) Close() error {
	return s.db.Close()
}
func (s *ClientStorage) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS stored_data (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			data BLOB NOT NULL,
			metadata TEXT,
			version INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_sync_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_deleted BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS data_history (
			id TEXT PRIMARY KEY,
			data_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			data BLOB NOT NULL,
			metadata TEXT,
			version INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_deleted BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (data_id) REFERENCES stored_data(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_stored_data_user_id ON stored_data(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_stored_data_updated_at ON stored_data(updated_at)`,
		`CREATE INDEX IF NOT EXISTS idx_data_history_data_id ON data_history(data_id)`,
		`CREATE INDEX IF NOT EXISTS idx_data_history_version ON data_history(version)`,
	}
	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}
	return nil
}
func (s *ClientStorage) SaveData(data *models.StoredData) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	var existingVersion int
	err = tx.QueryRow("SELECT version FROM stored_data WHERE id = ?", data.ID).Scan(&existingVersion)
	if err == sql.ErrNoRows {
		query := `INSERT INTO stored_data (id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		now := time.Now()
		_, err = tx.Exec(query, data.ID, data.UserID, data.Type, data.Title, data.Data, data.Metadata, data.Version, now, now, now, data.IsDeleted)
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	} else {
		query := `UPDATE stored_data SET type = ?, title = ?, data = ?, metadata = ?, version = ?, updated_at = ?, last_sync_at = ?, is_deleted = ? 
				  WHERE id = ?`
		now := time.Now()
		_, err = tx.Exec(query, data.Type, data.Title, data.Data, data.Metadata, data.Version, now, now, data.IsDeleted, data.ID)
		if err != nil {
			return fmt.Errorf("failed to update data: %w", err)
		}
	}
	if err := s.saveToHistory(tx, data); err != nil {
		return fmt.Errorf("failed to save to history: %w", err)
	}
	if err := s.cleanupHistory(tx, data.ID); err != nil {
		return fmt.Errorf("failed to cleanup history: %w", err)
	}
	return tx.Commit()
}
func (s *ClientStorage) GetData(id string) (*models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE id = ?`
	data := &models.StoredData{}
	err := s.db.QueryRow(query, id).Scan(
		&data.ID, &data.UserID, &data.Type, &data.Title, &data.Data, &data.Metadata,
		&data.Version, &data.CreatedAt, &data.UpdatedAt, &data.LastSyncAt, &data.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("data not found")
		}
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	return data, nil
}
func (s *ClientStorage) GetAllData(userID string) ([]models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE user_id = ? AND is_deleted = FALSE ORDER BY updated_at DESC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()
	var dataList []models.StoredData
	for rows.Next() {
		var data models.StoredData
		err := rows.Scan(
			&data.ID, &data.UserID, &data.Type, &data.Title, &data.Data, &data.Metadata,
			&data.Version, &data.CreatedAt, &data.UpdatedAt, &data.LastSyncAt, &data.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}
func (s *ClientStorage) GetDataSince(userID string, since time.Time) ([]models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE user_id = ? AND updated_at > ? ORDER BY updated_at DESC`
	rows, err := s.db.Query(query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()
	var dataList []models.StoredData
	for rows.Next() {
		var data models.StoredData
		err := rows.Scan(
			&data.ID, &data.UserID, &data.Type, &data.Title, &data.Data, &data.Metadata,
			&data.Version, &data.CreatedAt, &data.UpdatedAt, &data.LastSyncAt, &data.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}
func (s *ClientStorage) DeleteData(id string) error {
	query := `UPDATE stored_data SET is_deleted = TRUE, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}
	return nil
}
func (s *ClientStorage) GetLastSyncTime(userID string) (time.Time, error) {
	query := `SELECT MAX(last_sync_at) FROM stored_data WHERE user_id = ?`
	var lastSyncStr sql.NullString
	err := s.db.QueryRow(query, userID).Scan(&lastSyncStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get last sync time: %w", err)
	}
	if lastSyncStr.Valid && lastSyncStr.String != "" {
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02 15:04:05.999999999-07:00",
			"2006-01-02 15:04:05-07:00",
			"2006-01-02 15:04:05",
		}
		for _, format := range formats {
			if lastSync, err := time.Parse(format, lastSyncStr.String); err == nil {
				return lastSync, nil
			}
		}
		return time.Time{}, fmt.Errorf("failed to parse last sync time: %s", lastSyncStr.String)
	}
	return time.Time{}, nil
}
func (s *ClientStorage) UpdateLastSyncTime(userID string, syncTime time.Time) error {
	query := `UPDATE stored_data SET last_sync_at = ? WHERE user_id = ?`
	_, err := s.db.Exec(query, syncTime, userID)
	if err != nil {
		return fmt.Errorf("failed to update last sync time: %w", err)
	}
	return nil
}
func (s *ClientStorage) saveToHistory(tx *sql.Tx, data *models.StoredData) error {
	query := `INSERT INTO data_history (id, data_id, user_id, type, title, data, metadata, version, created_at, updated_at, is_deleted) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	historyID := fmt.Sprintf("%s_v%d", data.ID, data.Version)
	now := time.Now()
	_, err := tx.Exec(query, historyID, data.ID, data.UserID, data.Type, data.Title, data.Data, data.Metadata, data.Version, now, now, data.IsDeleted)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}
	return nil
}
func (s *ClientStorage) cleanupHistory(tx *sql.Tx, dataID string) error {
	query := `DELETE FROM data_history 
			  WHERE data_id = ? AND id NOT IN (
				  SELECT id FROM data_history 
				  WHERE data_id = ? 
				  ORDER BY version DESC 
				  LIMIT 10
			  )`
	_, err := tx.Exec(query, dataID, dataID)
	if err != nil {
		return fmt.Errorf("failed to cleanup history: %w", err)
	}
	return nil
}
func (s *ClientStorage) GetDataHistory(dataID string) ([]models.DataHistory, error) {
	query := `SELECT id, data_id, user_id, type, title, data, metadata, version, created_at, updated_at, is_deleted 
			  FROM data_history WHERE data_id = ? ORDER BY version DESC`
	rows, err := s.db.Query(query, dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()
	var history []models.DataHistory
	for rows.Next() {
		var h models.DataHistory
		err := rows.Scan(
			&h.ID, &h.DataID, &h.UserID, &h.Type, &h.Title, &h.Data, &h.Metadata,
			&h.Version, &h.CreatedAt, &h.UpdatedAt, &h.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		history = append(history, h)
	}
	return history, nil
}
