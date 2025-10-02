package database
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)
type DB struct {
	conn *sql.DB
}
func (db *DB) Conn() *sql.DB {
	return db.conn
}
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
func (db *DB) Close() error {
	return db.conn.Close()
}
func (db *DB) createTables() error {
	return nil
}
