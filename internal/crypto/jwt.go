// Package crypto provides encryption and decryption functionality for GophKeeper.
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTClaims represents the claims in a JWT token.
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

// JWTManager handles JWT token creation and validation.
type JWTManager struct {
	secretKey []byte
}

// NewJWTManager creates a new JWT manager with the given secret key.
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken generates a JWT token for the given user.
func (j *JWTManager) GenerateToken(userID, username string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		ExpiresAt: now.Add(expiration).Unix(),
		IssuedAt:  now.Unix(),
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	payload := headerB64 + "." + claimsB64

	signature := j.sign(payload)
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return payload + "." + signatureB64, nil
}

// ValidateToken validates a JWT token and returns the claims.
func (j *JWTManager) ValidateToken(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	if !j.verifySignature(payload, signature) {
		return nil, fmt.Errorf("invalid signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// sign creates a signature for the given payload.
func (j *JWTManager) sign(payload string) []byte {
	h := hmac.New(sha256.New, j.secretKey)
	h.Write([]byte(payload))
	return h.Sum(nil)
}

// verifySignature verifies the signature of the given payload.
func (j *JWTManager) verifySignature(payload string, signature []byte) bool {
	expectedSignature := j.sign(payload)
	return hmac.Equal(signature, expectedSignature)
}
