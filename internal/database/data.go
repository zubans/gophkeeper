package database
import (
	"database/sql"
	"fmt"
	"time"
	"gophkeeper/internal/models"
)
func (db *DB) CreateStoredData(data *models.StoredData) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `INSERT INTO stored_data (id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	now := time.Now()
	_, err = tx.Exec(query, data.ID, data.UserID, data.Type, data.Title, data.Data, data.Metadata, data.Version, now, now, now, data.IsDeleted)
	if err != nil {
		return fmt.Errorf("failed to create stored data: %w", err)
	}
	if err := db.saveToHistory(tx, data); err != nil {
		return fmt.Errorf("failed to save to history: %w", err)
	}
	if err := db.cleanupHistory(tx, data.ID); err != nil {
		return fmt.Errorf("failed to cleanup history: %w", err)
	}
	data.CreatedAt = now
	data.UpdatedAt = now
	data.LastSyncAt = now
	return tx.Commit()
}
func (db *DB) GetStoredDataByID(id string) (*models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE id = $1`
	data := &models.StoredData{}
	err := db.conn.QueryRow(query, id).Scan(
		&data.ID, &data.UserID, &data.Type, &data.Title, &data.Data, &data.Metadata,
		&data.Version, &data.CreatedAt, &data.UpdatedAt, &data.LastSyncAt, &data.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stored data not found")
		}
		return nil, fmt.Errorf("failed to get stored data: %w", err)
	}
	return data, nil
}
func (db *DB) GetStoredDataByUserID(userID string) ([]models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE user_id = $1 AND is_deleted = FALSE ORDER BY updated_at DESC`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query stored data: %w", err)
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
			return nil, fmt.Errorf("failed to scan stored data: %w", err)
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}
func (db *DB) GetStoredDataByUserIDSince(userID string, since time.Time) ([]models.StoredData, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted 
			  FROM stored_data WHERE user_id = $1 AND updated_at > $2 ORDER BY updated_at DESC`
	rows, err := db.conn.Query(query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query stored data: %w", err)
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
			return nil, fmt.Errorf("failed to scan stored data: %w", err)
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}
func (db *DB) UpdateStoredData(data *models.StoredData) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `UPDATE stored_data SET type = $2, title = $3, data = $4, metadata = $5, version = $6, updated_at = $7, last_sync_at = $8, is_deleted = $9 
			  WHERE id = $1`
	now := time.Now()
	data.UpdatedAt = now
	data.LastSyncAt = now
	data.Version++
	_, err = tx.Exec(query, data.ID, data.Type, data.Title, data.Data, data.Metadata, data.Version, data.UpdatedAt, data.LastSyncAt, data.IsDeleted)
	if err != nil {
		return fmt.Errorf("failed to update stored data: %w", err)
	}
	if err := db.saveToHistory(tx, data); err != nil {
		return fmt.Errorf("failed to save to history: %w", err)
	}
	if err := db.cleanupHistory(tx, data.ID); err != nil {
		return fmt.Errorf("failed to cleanup history: %w", err)
	}
	return tx.Commit()
}
func (db *DB) DeleteStoredData(id string) error {
	query := `UPDATE stored_data SET is_deleted = TRUE, updated_at = $1 WHERE id = $2`
	_, err := db.conn.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete stored data: %w", err)
	}
	return nil
}
func (db *DB) SyncStoredData(userID string, dataList []models.StoredData, lastSyncAt time.Time) ([]models.StoredData, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	serverData, err := db.GetStoredDataByUserIDSince(userID, lastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get server data: %w", err)
	}
	for _, clientData := range dataList {
		existingData, err := db.GetStoredDataByID(clientData.ID)
		if err != nil && err.Error() != "stored data not found" {
			return nil, fmt.Errorf("failed to check existing data: %w", err)
		}
		if existingData == nil {
			clientData.UserID = userID
			clientData.LastSyncAt = time.Now()
			if err := db.CreateStoredData(&clientData); err != nil {
				return nil, fmt.Errorf("failed to create new data: %w", err)
			}
		} else {
			if clientData.Version > existingData.Version {
				clientData.UserID = userID
				clientData.LastSyncAt = time.Now()
				if err := db.UpdateStoredData(&clientData); err != nil {
					return nil, fmt.Errorf("failed to update data: %w", err)
				}
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return serverData, nil
}
func (db *DB) saveToHistory(tx *sql.Tx, data *models.StoredData) error {
	query := `INSERT INTO data_history (id, data_id, user_id, type, title, data, metadata, version, created_at, updated_at, is_deleted) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	historyID := fmt.Sprintf("%s_v%d", data.ID, data.Version)
	now := time.Now()
	_, err := tx.Exec(query, historyID, data.ID, data.UserID, data.Type, data.Title, data.Data, data.Metadata, data.Version, now, now, data.IsDeleted)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}
	return nil
}
func (db *DB) cleanupHistory(tx *sql.Tx, dataID string) error {
	query := `DELETE FROM data_history 
			  WHERE data_id = $1 AND id NOT IN (
				  SELECT id FROM data_history 
				  WHERE data_id = $1 
				  ORDER BY version DESC 
				  LIMIT 10
			  )`
	_, err := tx.Exec(query, dataID)
	if err != nil {
		return fmt.Errorf("failed to cleanup history: %w", err)
	}
	return nil
}
func (db *DB) GetDataHistory(dataID string) ([]models.DataHistory, error) {
	query := `SELECT id, data_id, user_id, type, title, data, metadata, version, created_at, updated_at, is_deleted 
			  FROM data_history WHERE data_id = $1 ORDER BY version DESC`
	rows, err := db.conn.Query(query, dataID)
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
