// Package server implements the GophKeeper HTTP server.
package server

import (
	"fmt"
	"time"

	"gophkeeper/internal/crypto"
	"gophkeeper/internal/database"
	"gophkeeper/internal/models"
)

// DataService handles data storage and synchronization.
type DataService struct {
	db        *database.DB
	encryptor *crypto.Encryptor
}

// NewDataService creates a new data service.
func NewDataService(db *database.DB, encryptor *crypto.Encryptor) *DataService {
	return &DataService{
		db:        db,
		encryptor: encryptor,
	}
}

// GetUserData retrieves all data for a user.
func (d *DataService) GetUserData(userID string) ([]models.StoredData, error) {
	dataList, err := d.db.GetStoredDataByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	// Decrypt data
	for i := range dataList {
		if err := d.decryptData(&dataList[i]); err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
	}

	return dataList, nil
}

// CreateData creates new data for a user.
func (d *DataService) CreateData(data *models.StoredData) error {
	// Encrypt data before storing
	if err := d.encryptData(data); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Generate ID if not provided
	if data.ID == "" {
		data.ID = generateID()
	}

	// Set version
	data.Version = 1

	if err := d.db.CreateStoredData(data); err != nil {
		return fmt.Errorf("failed to create data: %w", err)
	}

	return nil
}

// UpdateData updates existing data.
func (d *DataService) UpdateData(data *models.StoredData) error {
	// Encrypt data before storing
	if err := d.encryptData(data); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if err := d.db.UpdateStoredData(data); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}

// DeleteData deletes data by ID.
func (d *DataService) DeleteData(dataID, userID string) error {
	// Verify ownership
	data, err := d.db.GetStoredDataByID(dataID)
	if err != nil {
		return fmt.Errorf("data not found: %w", err)
	}

	if data.UserID != userID {
		return fmt.Errorf("access denied")
	}

	if err := d.db.DeleteStoredData(dataID); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

// SyncData synchronizes data between client and server.
func (d *DataService) SyncData(userID string, req *models.DataSyncRequest) (*models.DataSyncResponse, error) {
	// Encrypt client data before processing
	for i := range req.Data {
		if err := d.encryptData(&req.Data[i]); err != nil {
			return nil, fmt.Errorf("failed to encrypt client data: %w", err)
		}
	}

	// Get server data since last sync
	serverData, err := d.db.GetStoredDataByUserIDSince(userID, req.LastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get server data: %w", err)
	}

	// Process client data with conflict resolution
	var conflicts []models.Conflict
	for _, clientData := range req.Data {
		clientData.UserID = userID

		// Check if data exists on server
		serverDataItem, err := d.db.GetStoredDataByID(clientData.ID)
		if err != nil && err.Error() != "stored data not found" {
			return nil, fmt.Errorf("failed to check existing data: %w", err)
		}

		if serverDataItem == nil {
			// New data, create it
			if err := d.db.CreateStoredData(&clientData); err != nil {
				return nil, fmt.Errorf("failed to create new data: %w", err)
			}
		} else {
			// Check for conflicts based on timestamp
			if clientData.UpdatedAt.After(serverDataItem.UpdatedAt) {
				// Client data is newer, update server
				if err := d.db.UpdateStoredData(&clientData); err != nil {
					return nil, fmt.Errorf("failed to update data: %w", err)
				}
			} else if serverDataItem.UpdatedAt.After(clientData.UpdatedAt) {
				// Server data is newer, add to conflicts
				conflicts = append(conflicts, models.Conflict{
					LocalData:  clientData,
					ServerData: *serverDataItem,
					Reason:     "Server has newer version",
				})
			} else {
				// Same timestamp, check version
				if clientData.Version > serverDataItem.Version {
					if err := d.db.UpdateStoredData(&clientData); err != nil {
						return nil, fmt.Errorf("failed to update data: %w", err)
					}
				} else if serverDataItem.Version > clientData.Version {
					conflicts = append(conflicts, models.Conflict{
						LocalData:  clientData,
						ServerData: *serverDataItem,
						Reason:     "Server has higher version",
					})
				}
			}
		}
	}

	// Decrypt server data before sending
	for i := range serverData {
		if err := d.decryptData(&serverData[i]); err != nil {
			return nil, fmt.Errorf("failed to decrypt server data: %w", err)
		}
	}

	response := &models.DataSyncResponse{
		Data:       serverData,
		LastSyncAt: time.Now(),
		Conflicts:  conflicts,
	}

	return response, nil
}

// encryptData encrypts the data field of StoredData.
func (d *DataService) encryptData(data *models.StoredData) error {
	encrypted, err := d.encryptor.Encrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	data.Data = encrypted
	return nil
}

// decryptData decrypts the data field of StoredData.
func (d *DataService) decryptData(data *models.StoredData) error {
	decrypted, err := d.encryptor.Decrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}
	data.Data = decrypted
	return nil
}
