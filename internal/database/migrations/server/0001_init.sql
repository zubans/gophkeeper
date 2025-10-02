-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stored_data (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL,
    metadata TEXT,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_sync_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS data_history (
    id VARCHAR(64) PRIMARY KEY,
    data_id VARCHAR(36) NOT NULL REFERENCES stored_data(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL,
    metadata TEXT,
    version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_stored_data_user_id ON stored_data(user_id);
CREATE INDEX IF NOT EXISTS idx_stored_data_updated_at ON stored_data(updated_at);
CREATE INDEX IF NOT EXISTS idx_data_history_data_id ON data_history(data_id);
CREATE INDEX IF NOT EXISTS idx_data_history_version ON data_history(version);

-- +goose Down
DROP INDEX IF EXISTS idx_data_history_version;
DROP INDEX IF EXISTS idx_data_history_data_id;
DROP INDEX IF EXISTS idx_stored_data_updated_at;
DROP INDEX IF EXISTS idx_stored_data_user_id;
DROP TABLE IF EXISTS data_history;
DROP TABLE IF EXISTS stored_data;
DROP TABLE IF EXISTS users;