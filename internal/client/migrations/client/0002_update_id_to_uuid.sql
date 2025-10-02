-- +goose Up
CREATE TABLE IF NOT EXISTS stored_data_new (
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

CREATE TABLE IF NOT EXISTS data_history_new (
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

INSERT INTO stored_data_new (id, user_id, type, title, data, metadata, version, created_at, updated_at, last_sync_at, is_deleted)
SELECT 
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) as id,
    user_id,
    type,
    title,
    data,
    metadata,
    version,
    created_at,
    updated_at,
    last_sync_at,
    is_deleted
FROM stored_data;

INSERT INTO data_history_new (id, data_id, user_id, type, title, data, metadata, version, created_at, updated_at, is_deleted)
SELECT 
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) as id,
    data_id,
    user_id,
    type,
    title,
    data,
    metadata,
    version,
    created_at,
    updated_at,
    is_deleted
FROM data_history;

DROP TABLE IF EXISTS stored_data;
DROP TABLE IF EXISTS data_history;

ALTER TABLE stored_data_new RENAME TO stored_data;
ALTER TABLE data_history_new RENAME TO data_history;

CREATE INDEX IF NOT EXISTS idx_stored_data_user_id ON stored_data(user_id);
CREATE INDEX IF NOT EXISTS idx_stored_data_updated_at ON stored_data(updated_at);
CREATE INDEX IF NOT EXISTS idx_data_history_data_id ON data_history(data_id);
CREATE INDEX IF NOT EXISTS idx_data_history_version ON data_history(version);

-- +goose Down
DROP INDEX IF EXISTS idx_data_history_version;
DROP INDEX IF EXISTS idx_data_history_data_id;
DROP INDEX IF EXISTS idx_stored_data_updated_at;
DROP INDEX IF EXISTS idx_stored_data_user_id;

ALTER TABLE stored_data RENAME TO stored_data_old;
ALTER TABLE data_history RENAME TO data_history_old;

CREATE TABLE stored_data (
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

CREATE TABLE data_history (
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

DROP TABLE stored_data_old;
DROP TABLE data_history_old;