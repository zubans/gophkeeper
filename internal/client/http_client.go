package client
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"gophkeeper/internal/models"
)
type HTTPClientImpl struct {
	serverURL  string
	httpClient *http.Client
}
func NewHTTPClient(serverURL string) *HTTPClientImpl {
	return &HTTPClientImpl{
		serverURL:  serverURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}
func (h *HTTPClientImpl) Register(req *models.UserRegistrationRequest) (*models.AuthResponse, error) {
	var response models.AuthResponse
	if err := h.makeRequest("POST", "/api/v1/register", req, &response, ""); err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}
	return &response, nil
}
func (h *HTTPClientImpl) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	var response models.AuthResponse
	if err := h.makeRequest("POST", "/api/v1/login", req, &response, ""); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	return &response, nil
}
func (h *HTTPClientImpl) AddData(data *models.StoredData, token string) error {
	return h.makeRequest("POST", "/api/v1/data", data, data, token)
}
func (h *HTTPClientImpl) DeleteData(id, token string) error {
	url := fmt.Sprintf("%s/api/v1/data?id=%s", h.serverURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s", string(body))
	}
	return nil
}
func (h *HTTPClientImpl) SyncData(req *models.DataSyncRequest, token string) (*models.DataSyncResponse, error) {
	var response models.DataSyncResponse
	if err := h.makeRequest("POST", "/api/v1/sync", req, &response, token); err != nil {
		return nil, fmt.Errorf("sync failed: %w", err)
	}
	return &response, nil
}
func (h *HTTPClientImpl) makeRequest(method, path string, body interface{}, result interface{}, token string) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, h.serverURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		var errorResp models.ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			return fmt.Errorf("request failed: %s", errorResp.Error)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	if result != nil {
		var wrap models.APIResponse
		if err := json.Unmarshal(respBody, &wrap); err == nil {
			if !wrap.Success {
				if wrap.Error != "" {
					return fmt.Errorf("request failed: %s", wrap.Error)
				}
				return fmt.Errorf("request failed")
			}
			if wrap.Data != nil {
				dataBytes, err := json.Marshal(wrap.Data)
				if err != nil {
					return fmt.Errorf("failed to marshal wrapped data: %w", err)
				}
				if err := json.Unmarshal(dataBytes, result); err != nil {
					return fmt.Errorf("failed to unmarshal wrapped data: %w", err)
				}
				return nil
			}
		}
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}
	return nil
}
