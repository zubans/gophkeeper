package client
import (
	"fmt"
	"gophkeeper/internal/models"
)
type AuthServiceImpl struct {
	httpClient   HTTPClient
	tokenManager TokenManager
	token        string
	userID       string
}
func NewAuthService(httpClient HTTPClient, tokenManager TokenManager) *AuthServiceImpl {
	service := &AuthServiceImpl{
		httpClient:   httpClient,
		tokenManager: tokenManager,
	}
	if token, err := tokenManager.LoadToken(); err == nil && token != "" {
		service.token = token
	}
	return service
}
func (a *AuthServiceImpl) Register(username, email, password string) (*models.AuthResponse, error) {
	req := &models.UserRegistrationRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	response, err := a.httpClient.Register(req)
	if err != nil {
		return nil, err
	}
	a.token = response.Token
	a.userID = response.User.ID
	if err := a.tokenManager.SaveToken(a.token); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}
	return response, nil
}
func (a *AuthServiceImpl) Login(username, password string) (*models.AuthResponse, error) {
	req := &models.UserLoginRequest{
		Username: username,
		Password: password,
	}
	response, err := a.httpClient.Login(req)
	if err != nil {
		return nil, err
	}
	a.token = response.Token
	a.userID = response.User.ID
	if err := a.tokenManager.SaveToken(a.token); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}
	return response, nil
}
func (a *AuthServiceImpl) IsAuthenticated() bool {
	return a.token != ""
}
func (a *AuthServiceImpl) GetToken() string {
	return a.token
}
func (a *AuthServiceImpl) GetUserID() string {
	return a.userID
}
func (a *AuthServiceImpl) Logout() error {
	a.token = ""
	a.userID = ""
	return a.tokenManager.ClearToken()
}
