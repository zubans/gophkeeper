-- Initialize GophKeeper database
-- This script runs when PostgreSQL container starts for the first time

-- Create database if it doesn't exist (already created by POSTGRES_DB)
-- CREATE DATABASE gophkeeper;

-- Connect to the database
\c gophkeeper;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The actual schema will be created by the server migrations
-- This file is just for any additional initialization if needed
