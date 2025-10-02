package client
import (
	"encoding/json"
	"fmt"
	"strings"
	"gophkeeper/internal/models"
)
type DataServiceImpl struct {
	storage     Storage
	httpClient  HTTPClient
	encryptor   Encryptor
	authService AuthService
}
func NewDataService(storage Storage, httpClient HTTPClient, encryptor Encryptor, authService AuthService) *DataServiceImpl {
	return &DataServiceImpl{
		storage:     storage,
		httpClient:  httpClient,
		encryptor:   encryptor,
		authService: authService,
	}
}
func (d *DataServiceImpl) AddData(dataType, title string, data []string) error {
	if !d.authService.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}
	storedData := &models.StoredData{
		ID:      GenerateID(),
		UserID:  d.authService.GetUserID(),
		Type:    models.DataType(dataType),
		Title:   title,
		Version: 1,
	}
	if err := d.processDataByType(storedData, dataType, data); err != nil {
		return err
	}
	if err := d.storage.SaveData(storedData); err != nil {
		return fmt.Errorf("failed to save data locally: %w", err)
	}
	encryptedData := *storedData
	if err := d.encryptData(&encryptedData); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	if err := d.httpClient.AddData(&encryptedData, d.authService.GetToken()); err != nil {
		return fmt.Errorf("failed to add data to server: %w", err)
	}
	return nil
}
func (d *DataServiceImpl) GetData(id string) (*models.StoredData, error) {
	if !d.authService.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}
	data, err := d.storage.GetData(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	return data, nil
}
func (d *DataServiceImpl) GetDataList() ([]models.StoredData, error) {
	if !d.authService.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}
	dataList, err := d.storage.GetAllData(d.authService.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	return dataList, nil
}
func (d *DataServiceImpl) DeleteData(id string) error {
	if !d.authService.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}
	if err := d.storage.DeleteData(id); err != nil {
		return fmt.Errorf("failed to delete data locally: %w", err)
	}
	if err := d.httpClient.DeleteData(id, d.authService.GetToken()); err != nil {
		return fmt.Errorf("failed to delete data from server: %w", err)
	}
	return nil
}
func (d *DataServiceImpl) ShowHistory(id string) error {
	if !d.authService.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}
	history, err := d.storage.GetDataHistory(id)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}
	if len(history) == 0 {
		fmt.Println("No history found for this data item.")
		return nil
	}
	fmt.Printf("History for data ID: %s\n", id)
	fmt.Printf("%-8s %-20s %-50s %-20s\n", "Version", "Type", "Title", "Updated")
	fmt.Println(strings.Repeat("-", 100))
	for _, h := range history {
		status := "Active"
		if h.IsDeleted {
			status = "Deleted"
		}
		fmt.Printf("%-8d %-20s %-50s %-20s %s\n",
			h.Version, h.Type, h.Title, h.UpdatedAt.Format("2006-01-02 15:04:05"), status)
	}
	return nil
}
func (d *DataServiceImpl) processDataByType(storedData *models.StoredData, dataType string, data []string) error {
	switch models.DataType(dataType) {
	case models.DataTypeLoginPassword:
		if len(data) < 2 {
			return fmt.Errorf("login/password data requires at least login and password")
		}
		loginData := models.LoginPasswordData{
			Login:    data[0],
			Password: data[1],
		}
		if len(data) > 2 {
			loginData.Website = data[2]
		}
		if len(data) > 3 {
			loginData.Notes = data[3]
		}
		jsonData, _ := json.Marshal(loginData)
		storedData.Data = jsonData
	case models.DataTypeText:
		if len(data) == 0 {
			return fmt.Errorf("text data is required")
		}
		storedData.Data = []byte(data[0])
	case models.DataTypeBankCard:
		if len(data) < 4 {
			return fmt.Errorf("bank card data requires card number, expiry, CVV, and cardholder")
		}
		cardData := models.BankCardData{
			CardNumber: data[0],
			ExpiryDate: data[1],
			CVV:        data[2],
			Cardholder: data[3],
		}
		if len(data) > 4 {
			cardData.Bank = data[4]
		}
		if len(data) > 5 {
			cardData.Notes = data[5]
		}
		jsonData, _ := json.Marshal(cardData)
		storedData.Data = jsonData
	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}
	return nil
}
func (d *DataServiceImpl) encryptData(data *models.StoredData) error {
	encrypted, err := d.encryptor.Encrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}
	data.Data = encrypted
	return nil
}
func (d *DataServiceImpl) decryptData(data *models.StoredData) error {
	decrypted, err := d.encryptor.Decrypt(data.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}
	data.Data = decrypted
	return nil
}
