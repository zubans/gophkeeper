// Package server implements the GophKeeper HTTP server.
package server

import (
	"fmt"
	"time"

	"gophkeeper/internal/crypto"
	"gophkeeper/internal/database"
	"gophkeeper/internal/models"

	"github.com/google/uuid"
)

// AuthService handles user authentication and authorization.
type AuthService struct {
	db         *database.DB
	jwtManager *crypto.JWTManager
}

// NewAuthService creates a new authentication service.
func NewAuthService(db *database.DB, jwtManager *crypto.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Register registers a new user.
func (a *AuthService) Register(req *models.UserRegistrationRequest) (*models.AuthResponse, error) {
	// Check if username already exists
	_, err := a.db.GetUserByUsername(req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	_, err = a.db.GetUserByEmail(req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           generateID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := a.db.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := a.jwtManager.GenerateToken(user.ID, user.Username, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create response
	response := &models.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	return response, nil
}

// Login authenticates a user.
func (a *AuthService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Get user by username
	user, err := a.db.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	valid, err := crypto.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := a.jwtManager.GenerateToken(user.ID, user.Username, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create response
	response := &models.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	return response, nil
}

// generateID generates a unique UUID v4.
func generateID() string {
	return uuid.New().String()
}
