package server
import (
	"fmt"
	"time"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/database"
	"gophkeeper/internal/models"
	"github.com/google/uuid"
)
type AuthService struct {
	db         *database.DB
	jwtManager *crypto.JWTManager
}
func NewAuthService(db *database.DB, jwtManager *crypto.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}
func (a *AuthService) Register(req *models.UserRegistrationRequest) (*models.AuthResponse, error) {
	_, err := a.db.GetUserByUsername(req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}
	_, err = a.db.GetUserByEmail(req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		ID:           generateID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	if err := a.db.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	token, err := a.jwtManager.GenerateToken(user.ID, user.Username, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	response := &models.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	return response, nil
}
func (a *AuthService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	user, err := a.db.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	valid, err := crypto.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}
	token, err := a.jwtManager.GenerateToken(user.ID, user.Username, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	response := &models.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	return response, nil
}
func generateID() string {
	return uuid.New().String()
}
