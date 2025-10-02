-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop foreign key constraints first
ALTER TABLE data_history DROP CONSTRAINT IF EXISTS data_history_data_id_fkey;
ALTER TABLE data_history DROP CONSTRAINT IF EXISTS data_history_user_id_fkey;
ALTER TABLE stored_data DROP CONSTRAINT IF EXISTS stored_data_user_id_fkey;

-- Update users table first
ALTER TABLE users ALTER COLUMN id TYPE UUID USING uuid_generate_v4();
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();

-- Update stored_data table
ALTER TABLE stored_data ALTER COLUMN id TYPE UUID USING uuid_generate_v4();
ALTER TABLE stored_data ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE stored_data ALTER COLUMN user_id TYPE UUID USING user_id::UUID;

-- Update data_history table
ALTER TABLE data_history ALTER COLUMN id TYPE UUID USING uuid_generate_v4();
ALTER TABLE data_history ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE data_history ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE data_history ALTER COLUMN data_id TYPE UUID USING data_id::UUID;

-- Recreate foreign key constraints
ALTER TABLE stored_data ADD CONSTRAINT stored_data_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE data_history ADD CONSTRAINT data_history_data_id_fkey 
    FOREIGN KEY (data_id) REFERENCES stored_data(id) ON DELETE CASCADE;
ALTER TABLE data_history ADD CONSTRAINT data_history_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
-- NOTE: Depending on production data, a full rollback may be unsafe.
-- This Down attempts to convert UUIDs back to text and recreate constraints.

ALTER TABLE data_history DROP CONSTRAINT IF EXISTS data_history_data_id_fkey;
ALTER TABLE data_history DROP CONSTRAINT IF EXISTS data_history_user_id_fkey;
ALTER TABLE stored_data DROP CONSTRAINT IF EXISTS stored_data_user_id_fkey;

ALTER TABLE data_history ALTER COLUMN data_id TYPE TEXT USING data_id::TEXT;
ALTER TABLE data_history ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE data_history ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE data_history ALTER COLUMN id DROP DEFAULT;

ALTER TABLE stored_data ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE stored_data ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE stored_data ALTER COLUMN id DROP DEFAULT;

ALTER TABLE users ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;

ALTER TABLE stored_data ADD CONSTRAINT stored_data_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE data_history ADD CONSTRAINT data_history_data_id_fkey 
    FOREIGN KEY (data_id) REFERENCES stored_data(id) ON DELETE CASCADE;
ALTER TABLE data_history ADD CONSTRAINT data_history_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;






