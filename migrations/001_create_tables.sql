-- Create URLs table
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    custom_alias VARCHAR(50) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    user_ip VARCHAR(45)
);

-- Create index on short_code for faster lookups
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);

-- Create index on custom_alias for faster lookups
CREATE INDEX IF NOT EXISTS idx_urls_custom_alias ON urls(custom_alias) WHERE custom_alias IS NOT NULL;

-- Create analytics table
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
    accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    referer TEXT,
    user_agent TEXT,
    ip_address VARCHAR(45)
);

-- Create index on url_id for faster analytics queries
CREATE INDEX IF NOT EXISTS idx_analytics_url_id ON analytics(url_id);