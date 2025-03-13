package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rakheshkrishna2005/url-shortener/internal/models"
)

// URLRepository handles database operations for URLs
type URLRepository struct {
	db *sqlx.DB
}

// NewURLRepository creates a new URLRepository
func NewURLRepository(db *sqlx.DB) *URLRepository {
	return &URLRepository{db: db}
}

// Store saves a URL to the database
func (r *URLRepository) Store(ctx context.Context, url *models.URL) error {
	query := `
		INSERT INTO urls (original_url, short_code, custom_alias, expires_at, user_ip)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		url.OriginalURL,
		url.ShortCode,
		url.CustomAlias,
		url.ExpiresAt,
		url.UserIP,
	).Scan(&url.ID, &url.CreatedAt)

	return err
}

// FindByShortCode retrieves a URL by its short code
func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `
		SELECT id, original_url, short_code, custom_alias, created_at, expires_at, user_ip
		FROM urls
		WHERE short_code = $1
	`

	url := &models.URL{}
	err := r.db.GetContext(ctx, url, query, shortCode)
	if err == sql.ErrNoRows {
		return nil, models.ErrURLNotFound
	}
	return url, err
}

// FindByID retrieves a URL by its ID
func (r *URLRepository) FindByID(ctx context.Context, id int64) (*models.URL, error) {
	query := `
		SELECT id, original_url, short_code, custom_alias, created_at, expires_at, user_ip
		FROM urls
		WHERE id = $1
	`

	url := &models.URL{}
	err := r.db.GetContext(ctx, url, query, id)
	if err == sql.ErrNoRows {
		return nil, models.ErrURLNotFound
	}
	return url, err
}

// FindByCustomAlias checks if a custom alias already exists
func (r *URLRepository) FindByCustomAlias(ctx context.Context, alias string) (*models.URL, error) {
	query := `
		SELECT id, original_url, short_code, custom_alias, created_at, expires_at, user_ip
		FROM urls
		WHERE custom_alias = $1
	`

	url := &models.URL{}
	err := r.db.GetContext(ctx, url, query, alias)
	if err == sql.ErrNoRows {
		return nil, models.ErrURLNotFound
	}
	return url, err
}

// Update updates a URL record
func (r *URLRepository) Update(ctx context.Context, url *models.URL) error {
	query := `
		UPDATE urls
		SET original_url = $1, custom_alias = $2, expires_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		url.OriginalURL,
		url.CustomAlias,
		url.ExpiresAt,
		url.ID,
	)
	return err
}

// Delete removes a URL record by ID
func (r *URLRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM urls WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// RecordClick adds a click event for analytics
func (r *URLRepository) RecordClick(ctx context.Context, event *models.ClickEvent) error {
	query := `
		INSERT INTO analytics (url_id, referer, user_agent, ip_address)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.URLID,
		event.Referer,
		event.UserAgent,
		event.IPAddress,
	)
	return err
}

// GetURLStats retrieves analytics for a URL
func (r *URLRepository) GetURLStats(ctx context.Context, urlID int64) (*models.URLStats, error) {
	query := `
		SELECT COUNT(*) as click_count, MAX(accessed_at) as last_click
		FROM analytics
		WHERE url_id = $1
	`

	var stats models.URLStats
	var lastClick sql.NullTime

	err := r.db.QueryRowContext(ctx, query, urlID).Scan(&stats.ClickCount, &lastClick)
	if err != nil {
		return nil, err
	}

	if lastClick.Valid {
		stats.LastClick = lastClick.Time
	} else {
		stats.LastClick = time.Time{}
	}

	return &stats, nil
}