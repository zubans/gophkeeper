package client

import (
	"encoding/json"
	"fmt"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/models"
	"strings"
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
		if userID, err := service.extractUserIDFromToken(token); err == nil {
			service.userID = userID
		}
	}
	return service
}
func (a *AuthServiceImpl) Register(username, email, password string) (*models.AuthResponse, error) {
	logger.Info("Registering user: %s", username)
	req := &models.UserRegistrationRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	response, err := a.httpClient.Register(req)
	if err != nil {
		logger.Error("Registration failed for user %s: %v", username, err)
		return nil, err
	}
	a.token = response.Token
	a.userID = response.User.ID
	logger.Debug("User registered successfully, user_id: %s", a.userID)
	if err := a.tokenManager.SaveToken(a.token); err != nil {
		logger.Error("Failed to save token: %v", err)
		return nil, fmt.Errorf("failed to save token: %w", err)
	}
	logger.Info("User %s registered and authenticated successfully", username)
	return response, nil
}
func (a *AuthServiceImpl) Login(username, password string) (*models.AuthResponse, error) {
	logger.Info("Logging in user: %s", username)
	req := &models.UserLoginRequest{
		Username: username,
		Password: password,
	}
	response, err := a.httpClient.Login(req)
	if err != nil {
		logger.Error("Login failed for user %s: %v", username, err)
		return nil, err
	}
	a.token = response.Token
	a.userID = response.User.ID
	logger.Debug("User logged in successfully, user_id: %s", a.userID)
	if err := a.tokenManager.SaveToken(a.token); err != nil {
		logger.Error("Failed to save token: %v", err)
		return nil, fmt.Errorf("failed to save token: %w", err)
	}
	logger.Info("User %s logged in successfully", username)
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
	logger.Info("Logging out user")
	a.token = ""
	a.userID = ""
	if err := a.tokenManager.ClearToken(); err != nil {
		logger.Error("Failed to clear token: %v", err)
		return err
	}
	logger.Info("User logged out successfully")
	return nil
}

func (a *AuthServiceImpl) extractUserIDFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	payload := parts[1]
	for len(payload)%4 != 0 {
		payload += "="
	}

	decoded, err := a.base64DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(decoded, &claims); err != nil {
		return "", fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims.UserID, nil
}

func (a *AuthServiceImpl) base64DecodeString(s string) ([]byte, error) {
	const base64URLChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	const base64StdChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")

	decoded := make([]byte, 0, len(s)*3/4)

	for i := 0; i < len(s); i += 4 {
		var chunk [4]byte
		var validChars int

		for j := 0; j < 4 && i+j < len(s); j++ {
			c := s[i+j]
			if c == '=' {
				break
			}

			var val byte
			if c >= 'A' && c <= 'Z' {
				val = c - 'A'
			} else if c >= 'a' && c <= 'z' {
				val = c - 'a' + 26
			} else if c >= '0' && c <= '9' {
				val = c - '0' + 52
			} else if c == '+' {
				val = 62
			} else if c == '/' {
				val = 63
			} else {
				return nil, fmt.Errorf("invalid character: %c", c)
			}

			chunk[j] = val
			validChars++
		}

		if validChars >= 2 {
			decoded = append(decoded, (chunk[0]<<2)|(chunk[1]>>4))
		}
		if validChars >= 3 {
			decoded = append(decoded, (chunk[1]<<4)|(chunk[2]>>2))
		}
		if validChars >= 4 {
			decoded = append(decoded, (chunk[2]<<6)|chunk[3])
		}
	}

	return decoded, nil
}
