package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

var alphanumericRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")

// ValidateURL checks if a URL is valid
func ValidateURL(urlStr string) error {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must have http or https scheme")
	}
	
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	
	return nil
}

// ValidateAlias checks if a custom alias is valid
func ValidateAlias(alias string) error {
	if len(alias) < 3 {
		return fmt.Errorf("alias must be at least 3 characters long")
	}
	
	if len(alias) > 50 {
		return fmt.Errorf("alias must not exceed 50 characters")
	}
	
	if !alphanumericRegex.MatchString(alias) {
		return fmt.Errorf("alias must contain only alphanumeric characters")
	}
	
	return nil
}