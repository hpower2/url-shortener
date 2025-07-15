-- ⚠️  LEGACY SCHEMA - DO NOT USE FOR NEW INSTALLATIONS
-- This is the old basic schema without user authentication
-- For new installations, use the migrations instead:
-- psql urlshortener < migrations/001_initial_schema.sql
-- psql urlshortener < migrations/002_add_users_table.sql

-- Create database (run this first)
-- CREATE DATABASE urlshortener;

-- Use the database
-- \c urlshortener;

-- Create the URLs table
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    click_count INTEGER DEFAULT 0
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);

-- Insert some sample data (optional)
INSERT INTO urls (short_code, original_url) VALUES 
('abc123', 'https://www.google.com'),
('def456', 'https://www.github.com')
ON CONFLICT (short_code) DO NOTHING; 