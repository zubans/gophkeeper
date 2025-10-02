-- +goose Up
CREATE TABLE IF NOT EXISTS stored_data (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    data BLOB NOT NULL,
    metadata TEXT,
    version INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_sync_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS data_history (
    id TEXT PRIMARY KEY,
    data_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    data BLOB NOT NULL,
    metadata TEXT,
    version INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_stored_data_user_id ON stored_data(user_id);
CREATE INDEX IF NOT EXISTS idx_stored_data_updated_at ON stored_data(updated_at);
CREATE INDEX IF NOT EXISTS idx_data_history_data_id ON data_history(data_id);
CREATE INDEX IF NOT EXISTS idx_data_history_version ON data_history(version);

CREATE TABLE IF NOT EXISTS sync_metadata (
    user_id TEXT PRIMARY KEY,
    last_sync_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS sync_metadata;
DROP INDEX IF EXISTS idx_data_history_version;
DROP INDEX IF EXISTS idx_data_history_data_id;
DROP INDEX IF EXISTS idx_stored_data_updated_at;
DROP INDEX IF EXISTS idx_stored_data_user_id;
DROP TABLE IF EXISTS data_history;
DROP TABLE IF EXISTS stored_data;