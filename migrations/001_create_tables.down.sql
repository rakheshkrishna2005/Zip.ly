-- Drop indexes
DROP INDEX IF EXISTS idx_analytics_accessed_at;
DROP INDEX IF EXISTS idx_analytics_url_id;
DROP INDEX IF EXISTS idx_urls_custom_alias;
DROP INDEX IF EXISTS idx_urls_short_code;

-- Drop tables
DROP TABLE IF EXISTS analytics;
DROP TABLE IF EXISTS urls; 