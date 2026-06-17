-- Create urls table for storing shortened URLs
CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    access_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on short_code for fast lookups
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);

-- Create index on created_at for time-based queries
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);

-- Add comment to table
COMMENT ON TABLE urls IS 'Stores shortened URLs and their metadata';
COMMENT ON COLUMN urls.id IS 'Primary key, auto-incrementing';
COMMENT ON COLUMN urls.url IS 'Original long URL';
COMMENT ON COLUMN urls.short_code IS 'Unique short code (7 characters)';
COMMENT ON COLUMN urls.access_count IS 'Number of times this short URL has been accessed';
COMMENT ON COLUMN urls.created_at IS 'Timestamp when the URL was created';
COMMENT ON COLUMN urls.updated_at IS 'Timestamp when the URL was last updated';
