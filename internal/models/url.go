package models

import (
	"errors"
	"time"
)

// URL represents a shortened URL in the system
type URL struct {
	ID          int64      `db:"id" json:"id"`
	OriginalURL string     `db:"original_url" json:"original_url" validate:"required,url"`
	ShortCode   string     `db:"short_code" json:"short_code"`
	CustomAlias *string    `db:"custom_alias" json:"custom_alias,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	UserIP      *string    `db:"user_ip" json:"user_ip,omitempty"`
}

// URLStats represents the analytics data for a URL
type URLStats struct {
	ClickCount int       `json:"click_count"`
	LastClick  time.Time `json:"last_click"`
}

// CreateURLRequest represents the payload for creating a new shortened URL
type CreateURLRequest struct {
	OriginalURL string  `json:"original_url" validate:"required,url"`
	CustomAlias *string `json:"custom_alias,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	ExpiresIn   *int    `json:"expires_in,omitempty" validate:"omitempty,min=1"`
}

// CreateURLResponse represents the response for a create URL request
type CreateURLResponse struct {
	ShortURL    string     `json:"short_url"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// URLDetailResponse represents the response for a URL details request
type URLDetailResponse struct {
	URL   *URL      `json:"url"`
	Stats *URLStats `json:"stats,omitempty"`
}

// ClickEvent represents a click on a shortened URL
type ClickEvent struct {
	URLID      int64     `db:"url_id"`
	AccessedAt time.Time `db:"accessed_at"`
	Referer    *string   `db:"referer"`
	UserAgent  *string   `db:"user_agent"`
	IPAddress  *string   `db:"ip_address"`
}

// Common errors
var (
	ErrURLNotFound      = errors.New("url not found")
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrURLExpired       = errors.New("url has expired")
	ErrDuplicateAlias   = errors.New("custom alias already exists")
)