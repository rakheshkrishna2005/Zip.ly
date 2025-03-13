package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	// Base62Charset contains all characters used for Base62 encoding
	Base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenerateShortCode creates a random short code of specified length
func GenerateShortCode(length int) (string, error) {
	// Hard limit to 10 characters to match database constraint
	if length <= 0 || length > 10 {
		length = 6 // Default safe length
	}

	charsetLength := big.NewInt(int64(len(Base62Charset)))
	result := strings.Builder{}
	result.Grow(length)

	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result.WriteByte(Base62Charset[randIndex.Int64()])
	}

	return result.String(), nil
}

// IsValidShortCode checks if a short code contains only valid characters
func IsValidShortCode(code string) bool {
	for _, c := range code {
		if !strings.ContainsRune(Base62Charset, c) {
			return false
		}
	}
	return len(code) > 0
}

// IsValidCustomAlias checks if a custom alias contains only valid characters
func IsValidCustomAlias(alias string) bool {
	if len(alias) < 3 || len(alias) > 10 {
		return false
	}

	for _, c := range alias {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}
