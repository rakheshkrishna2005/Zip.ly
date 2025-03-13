-- Create URLs table
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(50) NOT NULL UNIQUE,
    custom_alias VARCHAR(50) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    user_ip VARCHAR(45)
);

-- Create Analytics table
CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
    accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    referer TEXT,
    user_agent TEXT,
    ip_address VARCHAR(45)
);

-- Create indexes for better query performance
CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_custom_alias ON urls(custom_alias);
CREATE INDEX idx_analytics_url_id ON analytics(url_id);
CREATE INDEX idx_analytics_accessed_at ON analytics(accessed_at); 