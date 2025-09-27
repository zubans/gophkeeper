// Package database provides database operations for GophKeeper.
package database

import (
	"database/sql"
	"fmt"
	"time"

	"gophkeeper/internal/models"
)

// CreateUser creates a new user in the database.
func (db *DB) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	_, err := db.conn.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByUsername retrieves a user by username.
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE username = $1`

	user := &models.User{}
	err := db.conn.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email.
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE email = $1`

	user := &models.User{}
	err := db.conn.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID.
func (db *DB) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE id = $1`

	user := &models.User{}
	err := db.conn.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user.
func (db *DB) UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = $2, email = $3, password_hash = $4, updated_at = $5 
			  WHERE id = $1`

	user.UpdatedAt = time.Now()
	_, err := db.conn.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user by ID.
func (db *DB) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
