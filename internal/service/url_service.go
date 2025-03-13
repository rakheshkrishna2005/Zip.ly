package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/rakheshkrishna2005/url-shortener/internal/config"
	"github.com/rakheshkrishna2005/url-shortener/internal/models"
	"github.com/rakheshkrishna2005/url-shortener/internal/utils"
)

// URLRepository defines the interface for URL data access
type URLRepository interface {
	Store(ctx context.Context, url *models.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	FindByID(ctx context.Context, id int64) (*models.URL, error)
	FindByCustomAlias(ctx context.Context, alias string) (*models.URL, error)
	Update(ctx context.Context, url *models.URL) error
	Delete(ctx context.Context, id int64) error
	RecordClick(ctx context.Context, event *models.ClickEvent) error
	GetURLStats(ctx context.Context, urlID int64) (*models.URLStats, error)
}

// URLService handles the business logic for URL operations
type URLService struct {
	repo   URLRepository
	config *config.Config
}

// NewURLService creates a new URLService
func NewURLService(repo URLRepository, cfg *config.Config) *URLService {
	return &URLService{
		repo:   repo,
		config: cfg,
	}
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(ctx context.Context, req models.CreateURLRequest, userIP string) (*models.CreateURLResponse, error) {
	// Validate original URL
	_, err := url.ParseRequestURI(req.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	// Generate short code or use custom alias
	var shortCode string
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		// Check if alias is valid
		if !utils.IsValidCustomAlias(*req.CustomAlias) {
			return nil, fmt.Errorf("invalid custom alias: must be 3-50 alphanumeric characters")
		}
		
		// Ensure custom alias doesn't exceed database limit
		if len(*req.CustomAlias) > 10 {
			return nil, fmt.Errorf("custom alias too long: must be 10 characters or less")
		}

		// Check if alias already exists
		_, err := s.repo.FindByCustomAlias(ctx, *req.CustomAlias)
		if err == nil {
			return nil, models.ErrDuplicateAlias
		} else if err != models.ErrURLNotFound {
			return nil, err
		}

		shortCode = *req.CustomAlias
	} else {
		// Ensure short code length doesn't exceed 10
		codeLen := s.config.ShortCodeLen
		if codeLen > 10 {
			codeLen = 10
		}
		
		for i := 0; i < 5; i++ { // Try a few times in case of collision
			shortCode, err = utils.GenerateShortCode(codeLen)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short code: %w", err)
			}

			// Check if short code already exists
			_, err := s.repo.FindByShortCode(ctx, shortCode)
			if err == models.ErrURLNotFound {
				break // This is a good short code
			} else if err != nil {
				return nil, err
			}
			// If we get here, short code exists, try again
		}
	}

	// Calculate expiration time
	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(*req.ExpiresIn) * 24 * time.Hour)
		expiresAt = &exp
	} else if s.config.DefaultExpiry > 0 {
		exp := time.Now().Add(s.config.DefaultExpiry)
		expiresAt = &exp
	}

	// Create URL entity
	urlEntity := &models.URL{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		CustomAlias: req.CustomAlias,
		ExpiresAt:   expiresAt,
	}

	if userIP != "" {
		urlEntity.UserIP = &userIP
	}

	// Store in database
	err = s.repo.Store(ctx, urlEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to store URL: %w", err)
	}

	// Create response
	shortURL := fmt.Sprintf("%s/%s", s.config.BaseURL, shortCode)
	response := &models.CreateURLResponse{
		ShortURL:    shortURL,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		CustomAlias: req.CustomAlias,
		ExpiresAt:   expiresAt,
	}

	return response, nil
}

// GetURL retrieves a URL by short code
func (s *URLService) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	if !utils.IsValidShortCode(shortCode) {
		return nil, models.ErrInvalidShortCode
	}

	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Check if URL has expired
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return nil, models.ErrURLExpired
	}

	return url, nil
}

// GetURLByID retrieves a URL by ID
func (s *URLService) GetURLByID(ctx context.Context, id int64) (*models.URLDetailResponse, error) {
	url, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stats, err := s.repo.GetURLStats(ctx, id)
	if err != nil {
		// We can still return the URL even if stats retrieval fails
		stats = &models.URLStats{
			ClickCount: 0,
		}
	}

	return &models.URLDetailResponse{
		URL:   url,
		Stats: stats,
	}, nil
}

// UpdateURL updates a URL
func (s *URLService) UpdateURL(ctx context.Context, id int64, originalURL string, customAlias *string, expiresIn *int) error {
	urlModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Update original URL if provided
	if originalURL != "" {
		_, err := url.ParseRequestURI(originalURL)
		if err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}
		urlModel.OriginalURL = originalURL
	}

	// Update custom alias if provided
	if customAlias != nil {
		if *customAlias == "" {
			urlModel.CustomAlias = nil
		} else {
			if !utils.IsValidCustomAlias(*customAlias) {
				return fmt.Errorf("invalid custom alias: must be 3-50 alphanumeric characters")
			}

			// Check if alias already exists and doesn't belong to this URL
			existingURL, err := s.repo.FindByCustomAlias(ctx, *customAlias)
			if err == nil && existingURL.ID != id {
				return models.ErrDuplicateAlias
			} else if err != nil && err != models.ErrURLNotFound {
				return err
			}

			urlModel.CustomAlias = customAlias
		}
	}

	// Update expiry if provided
	if expiresIn != nil {
		if *expiresIn <= 0 {
			urlModel.ExpiresAt = nil
		} else {
			exp := time.Now().Add(time.Duration(*expiresIn) * 24 * time.Hour)
			urlModel.ExpiresAt = &exp
		}
	}

	return s.repo.Update(ctx, urlModel)
}

// DeleteURL deletes a URL by ID
func (s *URLService) DeleteURL(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// RecordClick records a click event for a URL
func (s *URLService) RecordClick(ctx context.Context, urlID int64, referer, userAgent, ipAddress string) error {
	event := &models.ClickEvent{
		URLID:      urlID,
		AccessedAt: time.Now(),
	}

	if referer != "" {
		event.Referer = &referer
	}

	if userAgent != "" {
		event.UserAgent = &userAgent
	}

	if ipAddress != "" {
		event.IPAddress = &ipAddress
	}

	return s.repo.RecordClick(ctx, event)
}
