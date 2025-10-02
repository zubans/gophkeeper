package server
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/database"
	"gophkeeper/internal/models"
)
type Server struct {
	db          *database.DB
	jwtManager  *crypto.JWTManager
	encryptor   *crypto.Encryptor
	authService *AuthService
	dataService *DataService
}
func NewServer(db *database.DB, jwtSecret, encryptionKey string) *Server {
	jwtManager := crypto.NewJWTManager(jwtSecret)
	encryptor := crypto.NewEncryptor(encryptionKey)
	authService := NewAuthService(db, jwtManager)
	dataService := NewDataService(db, encryptor)
	return &Server{
		db:          db,
		jwtManager:  jwtManager,
		encryptor:   encryptor,
		authService: authService,
		dataService: dataService,
	}
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/v1")
	switch {
	case path == "/register" && r.Method == "POST":
		s.handleRegister(w, r)
	case path == "/login" && r.Method == "POST":
		s.handleLogin(w, r)
	case path == "/data" && r.Method == "GET":
		s.handleGetData(w, r)
	case path == "/data" && r.Method == "POST":
		s.handleCreateData(w, r)
	case path == "/data" && r.Method == "PUT":
		s.handleUpdateData(w, r)
	case path == "/data" && r.Method == "DELETE":
		s.handleDeleteData(w, r)
	case path == "/sync" && r.Method == "POST":
		s.handleSyncData(w, r)
	default:
		s.writeErrorResponse(w, "Not found", http.StatusNotFound)
	}
}
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	response, err := s.authService.Register(&req)
	if err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.writeSuccessResponse(w, response)
}
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	s.writeSuccessResponse(w, response)
}
func (s *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.dataService.GetUserData(userID)
	if err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeSuccessResponse(w, data)
}
func (s *Server) handleCreateData(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var data models.StoredData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	data.UserID = userID
	if err := s.dataService.CreateData(&data); err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeSuccessResponse(w, data)
}
func (s *Server) handleUpdateData(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var data models.StoredData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	data.UserID = userID
	if err := s.dataService.UpdateData(&data); err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeSuccessResponse(w, data)
}
func (s *Server) handleDeleteData(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	dataID := r.URL.Query().Get("id")
	if dataID == "" {
		s.writeErrorResponse(w, "Data ID is required", http.StatusBadRequest)
		return
	}
	if err := s.dataService.DeleteData(dataID, userID); err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeSuccessResponse(w, map[string]string{"message": "Data deleted successfully"})
}
func (s *Server) handleSyncData(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req models.DataSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	response, err := s.dataService.SyncData(userID, &req)
	if err != nil {
		s.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeSuccessResponse(w, response)
}
func (s *Server) getUserIDFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return "", fmt.Errorf("invalid authorization header format")
	}
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return claims.UserID, nil
}
func (s *Server) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := models.NewSuccessResponse(data)
	json.NewEncoder(w).Encode(response)
}
func (s *Server) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	response := models.NewErrorResponse(message, statusCode)
	json.NewEncoder(w).Encode(response)
}
