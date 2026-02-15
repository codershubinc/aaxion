package anonymous_upload

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type UploadToken struct {
	Token         string
	TargetDir     string
	CreatedAt     time.Time
	MaxUploads    int
	UploadCount   int
	ExpiresAt     time.Time
	IsRevoked     bool
	MaxFileSize   int64 // in bytes
	AllowedTypes  []string
}

var (
	tokens     = make(map[string]*UploadToken)
	tokenMutex sync.RWMutex
)

var (
	ErrInvalidToken  = errors.New("invalid or expired token")
	ErrTokenRevoked  = errors.New("token has been revoked")
	ErrTokenExpired  = errors.New("token has expired")
	ErrMaxUploads    = errors.New("maximum uploads reached for this token")
)

// GenerateUploadToken creates a new one-time upload token
func GenerateUploadToken(targetDir string, maxUploads int, expiryHours int, maxFileSize int64) (string, error) {
	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	tokens[token] = &UploadToken{
		Token:       token,
		TargetDir:   targetDir,
		CreatedAt:   time.Now(),
		MaxUploads:  maxUploads,
		UploadCount: 0,
		ExpiresAt:   time.Now().Add(time.Hour * time.Duration(expiryHours)),
		IsRevoked:   false,
		MaxFileSize: maxFileSize,
	}

	return token, nil
}

// ValidateToken checks if token is valid and can be used
func ValidateToken(token string) (*UploadToken, error) {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	uploadToken, exists := tokens[token]
	if !exists {
		return nil, ErrInvalidToken
	}

	if uploadToken.IsRevoked {
		return nil, ErrTokenRevoked
	}

	if time.Now().After(uploadToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	if uploadToken.UploadCount >= uploadToken.MaxUploads {
		return nil, ErrMaxUploads
	}

	return uploadToken, nil
}

// IncrementUploadCount increments the upload counter for a token
func IncrementUploadCount(token string) error {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	uploadToken, exists := tokens[token]
	if !exists {
		return ErrInvalidToken
	}

	uploadToken.UploadCount++

	// Auto-revoke if max uploads reached
	if uploadToken.UploadCount >= uploadToken.MaxUploads {
		uploadToken.IsRevoked = true
	}

	return nil
}

// RevokeToken manually revokes a token
func RevokeToken(token string) error {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	uploadToken, exists := tokens[token]
	if !exists {
		return ErrInvalidToken
	}

	uploadToken.IsRevoked = true
	return nil
}

// CleanupExpiredTokens removes expired tokens from memory
func CleanupExpiredTokens() {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	now := time.Now()
	for token, uploadToken := range tokens {
		if now.After(uploadToken.ExpiresAt) || uploadToken.IsRevoked {
			delete(tokens, token)
		}
	}
}

// GetTokenInfo returns token information (for admin purposes)
func GetTokenInfo(token string) (*UploadToken, error) {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	uploadToken, exists := tokens[token]
	if !exists {
		return nil, ErrInvalidToken
	}

	return uploadToken, nil
}

// ListAllTokens returns all active tokens (for admin purposes)
func ListAllTokens() []*UploadToken {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	result := make([]*UploadToken, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, token)
	}

	return result
}

// Initialize starts the cleanup routine
func Initialize() {
	// Run cleanup every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			CleanupExpiredTokens()
		}
	}()
}
