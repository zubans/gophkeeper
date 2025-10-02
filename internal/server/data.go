package server
import (
	"fmt"
	"time"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/database"
	"gophkeeper/internal/models"
)
type DataService struct {
	db        *database.DB
	encryptor *crypto.Encryptor
}
func NewDataService(db *database.DB, encryptor *crypto.Encryptor) *DataService {
	return &DataService{
		db:        db,
		encryptor: encryptor,
	}
}
func (d *DataService) GetUserData(userID string) ([]models.StoredData, error) {
	dataList, err := d.db.GetStoredDataByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}
	for i := range dataList {
		if err := d.decryptData(&dataList[i]); err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
	}
	return dataList, nil
}
func (d *DataService) CreateData(data *models.StoredData) error {
	if err := d.encryptData(data); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	if data.ID == "" {
		data.ID = generateID()
	}
	data.Version = 1
	if err := d.db.CreateStoredData(data); err != nil {
		return fmt.Errorf("failed to create data: %w", err)
	}
	return nil
}
func (d *DataService) UpdateData(data *models.StoredData) error {
	if err := d.encryptData(data); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	if err := d.db.UpdateStoredData(data); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}
	return nil
}
func (d *DataService) DeleteData(dataID, userID string) error {
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
func (d *DataService) SyncData(userID string, req *models.DataSyncRequest) (*models.DataSyncResponse, error) {
	for i := range req.Data {
		if err := d.encryptData(&req.Data[i]); err != nil {
			return nil, fmt.Errorf("failed to encrypt client data: %w", err)
		}
	}
	serverData, err := d.db.GetStoredDataByUserIDSince(userID, req.LastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get server data: %w", err)
	}
	var conflicts []models.Conflict
	for _, clientData := range req.Data {
		clientData.UserID = userID
		serverDataItem, err := d.db.GetStoredDataByID(clientData.ID)
		if err != nil && err.Error() != "stored data not found" {
			return nil, fmt.Errorf("failed to check existing data: %w", err)
		}
		if serverDataItem == nil {
			if err := d.db.CreateStoredData(&clientData); err != nil {
				return nil, fmt.Errorf("failed to create new data: %w", err)
			}
		} else {
			if clientData.UpdatedAt.After(serverDataItem.UpdatedAt) {
				if err := d.db.UpdateStoredData(&clientData); err != nil {
					return nil, fmt.Errorf("failed to update data: %w", err)
				}
			} else if serverDataItem.UpdatedAt.After(clientData.UpdatedAt) {
				conflicts = append(conflicts, models.Conflict{
					LocalData:  clientData,
					ServerData: *serverDataItem,
					Reason:     "Server has newer version",
				})
			} else {
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
func (d *DataService) encryptData(data *models.StoredData) error {
	encrypted, err := d.encryptor.Encrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	data.Data = encrypted
	return nil
}
func (d *DataService) decryptData(data *models.StoredData) error {
	decrypted, err := d.encryptor.Decrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}
	data.Data = decrypted
	return nil
}
