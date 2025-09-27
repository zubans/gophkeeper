// Package database provides database operations for GophKeeper.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB represents the database connection and operations.
type DB struct {
	conn *sql.DB
}

// Conn exposes the underlying *sql.DB connection.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// NewDB creates a new database connection.
func NewDB(connectionString string) (*DB, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// createTables creates the necessary database tables.
// Note: This method is deprecated in favor of migrations.
// It's kept for backward compatibility but should not be used.
func (db *DB) createTables() error {
	// Tables are now created via migrations
	// This method is kept for backward compatibility
	return nil
}
