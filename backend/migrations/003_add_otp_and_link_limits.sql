-- Migration 003: Add OTP verification and link count tracking

-- Create OTP verification table
CREATE TABLE IF NOT EXISTS otp_verifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    purpose VARCHAR(50) NOT NULL DEFAULT 'email_verification', -- 'email_verification', 'password_reset', etc.
    is_verified BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP NULL,
    
    -- Indexes for performance
    UNIQUE(user_id, purpose, is_verified) -- Only one active OTP per user per purpose
);

-- Create indexes for OTP table
CREATE INDEX IF NOT EXISTS idx_otp_verifications_user_id ON otp_verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_otp_verifications_email ON otp_verifications(email);
CREATE INDEX IF NOT EXISTS idx_otp_verifications_otp_code ON otp_verifications(otp_code);
CREATE INDEX IF NOT EXISTS idx_otp_verifications_expires_at ON otp_verifications(expires_at);

-- Add email verification status to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMP NULL;

-- Add link count tracking to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS link_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS link_limit INTEGER DEFAULT 50;

-- Create indexes for new user columns
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified);
CREATE INDEX IF NOT EXISTS idx_users_link_count ON users(link_count);

-- Function to cleanup expired OTPs
CREATE OR REPLACE FUNCTION cleanup_expired_otps()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM otp_verifications 
    WHERE expires_at < CURRENT_TIMESTAMP 
    AND is_verified = FALSE;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Function to increment user link count
CREATE OR REPLACE FUNCTION increment_user_link_count(user_id_param INTEGER)
RETURNS VOID AS $$
BEGIN
    UPDATE users 
    SET link_count = link_count + 1 
    WHERE id = user_id_param;
END;
$$ LANGUAGE plpgsql;

-- Function to decrement user link count
CREATE OR REPLACE FUNCTION decrement_user_link_count(user_id_param INTEGER)
RETURNS VOID AS $$
BEGIN
    UPDATE users 
    SET link_count = GREATEST(link_count - 1, 0)
    WHERE id = user_id_param;
END;
$$ LANGUAGE plpgsql;

-- Function to check if user can create more links
CREATE OR REPLACE FUNCTION can_user_create_link(user_id_param INTEGER)
RETURNS BOOLEAN AS $$
DECLARE
    current_count INTEGER;
    link_limit_val INTEGER;
BEGIN
    SELECT link_count, link_limit INTO current_count, link_limit_val
    FROM users 
    WHERE id = user_id_param;
    
    IF current_count IS NULL THEN
        RETURN FALSE;
    END IF;
    
    RETURN current_count < link_limit_val;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update user link count when URLs are created
CREATE OR REPLACE FUNCTION trigger_increment_link_count()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM increment_user_link_count(NEW.user_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update user link count when URLs are deleted
CREATE OR REPLACE FUNCTION trigger_decrement_link_count()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM decrement_user_link_count(OLD.user_id);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for automatic link count management
CREATE TRIGGER urls_insert_trigger
    AFTER INSERT ON urls
    FOR EACH ROW
    EXECUTE FUNCTION trigger_increment_link_count();

CREATE TRIGGER urls_delete_trigger
    AFTER DELETE ON urls
    FOR EACH ROW
    EXECUTE FUNCTION trigger_decrement_link_count();

-- Initialize link count for existing users
UPDATE users SET link_count = (
    SELECT COUNT(*) FROM urls WHERE urls.user_id = users.id
) WHERE link_count = 0; 