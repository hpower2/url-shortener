-- Enhanced URL Shortener Database Schema

-- Create database (run this first)
-- CREATE DATABASE urlshortener;

-- Use the database
-- \c urlshortener;

-- Create the URLs table with enhanced fields
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    click_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP NULL,
    user_agent TEXT,
    ip_address INET,
    
    -- Indexes for performance
    CONSTRAINT urls_short_code_key UNIQUE (short_code)
);

-- Create click events table for detailed analytics
CREATE TABLE IF NOT EXISTS click_events (
    id SERIAL PRIMARY KEY,
    url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(2),
    city VARCHAR(100),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for click_events table
CREATE INDEX IF NOT EXISTS idx_click_events_url_id ON click_events(url_id);
CREATE INDEX IF NOT EXISTS idx_click_events_clicked_at ON click_events(clicked_at);
CREATE INDEX IF NOT EXISTS idx_click_events_country ON click_events(country);

-- Create analytics summary table (for performance)
CREATE TABLE IF NOT EXISTS url_analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_clicks INTEGER DEFAULT 0,
    unique_clicks INTEGER DEFAULT 0,
    top_countries JSONB,
    top_referrers JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint to prevent duplicate entries
    UNIQUE(url_id, date)
);

-- Create indexes for url_analytics table
CREATE INDEX IF NOT EXISTS idx_url_analytics_url_id ON url_analytics(url_id);
CREATE INDEX IF NOT EXISTS idx_url_analytics_date ON url_analytics(date);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);
CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls(expires_at);
CREATE INDEX IF NOT EXISTS idx_urls_is_active ON urls(is_active);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_urls_updated_at 
    BEFORE UPDATE ON urls 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_url_analytics_updated_at 
    BEFORE UPDATE ON url_analytics 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to increment click count
CREATE OR REPLACE FUNCTION increment_click_count(short_code_param VARCHAR(20))
RETURNS VOID AS $$
BEGIN
    UPDATE urls 
    SET click_count = click_count + 1 
    WHERE short_code = short_code_param;
END;
$$ LANGUAGE plpgsql;

-- Create view for URL statistics
CREATE OR REPLACE VIEW url_stats AS
SELECT 
    u.id,
    u.short_code,
    u.original_url,
    u.created_at,
    u.updated_at,
    u.click_count,
    u.is_active,
    u.expires_at,
    COUNT(ce.id) as total_events,
    COUNT(DISTINCT ce.ip_address) as unique_visitors,
    COUNT(CASE WHEN ce.clicked_at >= CURRENT_DATE THEN 1 END) as clicks_today,
    COUNT(CASE WHEN ce.clicked_at >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as clicks_this_week,
    COUNT(CASE WHEN ce.clicked_at >= CURRENT_DATE - INTERVAL '30 days' THEN 1 END) as clicks_this_month
FROM urls u
LEFT JOIN click_events ce ON u.id = ce.url_id
GROUP BY u.id, u.short_code, u.original_url, u.created_at, u.updated_at, u.click_count, u.is_active, u.expires_at;

-- Insert some sample data (optional)
INSERT INTO urls (short_code, original_url, user_agent, ip_address) VALUES 
('abc123', 'https://www.google.com', 'Mozilla/5.0 (Sample)', '127.0.0.1'),
('def456', 'https://www.github.com', 'Mozilla/5.0 (Sample)', '127.0.0.1'),
('xyz789', 'https://www.stackoverflow.com', 'Mozilla/5.0 (Sample)', '127.0.0.1')
ON CONFLICT (short_code) DO NOTHING;

-- Create function to clean up expired URLs
CREATE OR REPLACE FUNCTION cleanup_expired_urls()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM urls 
    WHERE expires_at IS NOT NULL 
    AND expires_at < CURRENT_TIMESTAMP 
    AND is_active = FALSE;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create function to get URL analytics
CREATE OR REPLACE FUNCTION get_url_analytics(url_id_param INTEGER, days_param INTEGER DEFAULT 30)
RETURNS TABLE (
    total_clicks BIGINT,
    unique_clicks BIGINT,
    clicks_today BIGINT,
    clicks_this_week BIGINT,
    top_countries JSONB,
    top_referrers JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(ce.id) as total_clicks,
        COUNT(DISTINCT ce.ip_address) as unique_clicks,
        COUNT(CASE WHEN ce.clicked_at >= CURRENT_DATE THEN 1 END) as clicks_today,
        COUNT(CASE WHEN ce.clicked_at >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as clicks_this_week,
        (
            SELECT jsonb_agg(jsonb_build_object('country', country, 'clicks', click_count))
            FROM (
                SELECT ce2.country, COUNT(*) as click_count
                FROM click_events ce2
                WHERE ce2.url_id = url_id_param
                AND ce2.clicked_at >= CURRENT_DATE - INTERVAL '1 day' * days_param
                AND ce2.country IS NOT NULL
                GROUP BY ce2.country
                ORDER BY click_count DESC
                LIMIT 10
            ) country_stats
        ) as top_countries,
        (
            SELECT jsonb_agg(jsonb_build_object('referrer', referer, 'clicks', click_count))
            FROM (
                SELECT ce3.referer, COUNT(*) as click_count
                FROM click_events ce3
                WHERE ce3.url_id = url_id_param
                AND ce3.clicked_at >= CURRENT_DATE - INTERVAL '1 day' * days_param
                AND ce3.referer IS NOT NULL
                GROUP BY ce3.referer
                ORDER BY click_count DESC
                LIMIT 10
            ) referrer_stats
        ) as top_referrers
    FROM click_events ce
    WHERE ce.url_id = url_id_param
    AND ce.clicked_at >= CURRENT_DATE - INTERVAL '1 day' * days_param;
END;
$$ LANGUAGE plpgsql; 


-- Add users table and update URLs table with user relationship

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add user_id column to urls table
ALTER TABLE urls ADD COLUMN IF NOT EXISTS user_id INTEGER;

-- Add foreign key constraint
ALTER TABLE urls ADD CONSTRAINT fk_urls_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);

-- Update trigger for users table
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert a default user for existing URLs (optional)
INSERT INTO users (email, password, first_name, last_name, is_active)
VALUES ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin', 'User', true)
ON CONFLICT (email) DO NOTHING;

-- Update existing URLs to belong to the default user
UPDATE urls SET user_id = (SELECT id FROM users WHERE email = 'admin@example.com' LIMIT 1) 
WHERE user_id IS NULL; 