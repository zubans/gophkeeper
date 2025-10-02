package client
import (
	"fmt"
	"gophkeeper/internal/models"
)
type SyncServiceImpl struct {
	storage     Storage
	httpClient  HTTPClient
	encryptor   Encryptor
	authService AuthService
}
func NewSyncService(storage Storage, httpClient HTTPClient, encryptor Encryptor, authService AuthService) *SyncServiceImpl {
	return &SyncServiceImpl{
		storage:     storage,
		httpClient:  httpClient,
		encryptor:   encryptor,
		authService: authService,
	}
}
func (s *SyncServiceImpl) SyncData() error {
	if !s.authService.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}
	userID := s.authService.GetUserID()
	lastSyncTime, err := s.storage.GetLastSyncTime(userID)
	if err != nil {
		return fmt.Errorf("failed to get last sync time: %w", err)
	}
	localData, err := s.storage.GetDataSince(userID, lastSyncTime)
	if err != nil {
		return fmt.Errorf("failed to get local data: %w", err)
	}
	encryptedLocalData := make([]models.StoredData, len(localData))
	for i, data := range localData {
		encryptedLocalData[i] = data
		if err := s.encryptData(&encryptedLocalData[i]); err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}
	}
	req := &models.DataSyncRequest{
		LastSyncAt: lastSyncTime,
		Data:       encryptedLocalData,
	}
	response, err := s.httpClient.SyncData(req, s.authService.GetToken())
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}
	for i := range response.Data {
		if err := s.decryptData(&response.Data[i]); err != nil {
			return fmt.Errorf("failed to decrypt server data: %w", err)
		}
	}
	for _, data := range response.Data {
		if err := s.storage.SaveData(&data); err != nil {
			return fmt.Errorf("failed to save server data: %w", err)
		}
	}
	if err := s.storage.UpdateLastSyncTime(userID, response.LastSyncAt); err != nil {
		return fmt.Errorf("failed to update last sync time: %w", err)
	}
	return nil
}
func (s *SyncServiceImpl) encryptData(data *models.StoredData) error {
	encrypted, err := s.encryptor.Encrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	data.Data = encrypted
	return nil
}
func (s *SyncServiceImpl) decryptData(data *models.StoredData) error {
	decrypted, err := s.encryptor.Decrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}
	data.Data = decrypted
	return nil
}
