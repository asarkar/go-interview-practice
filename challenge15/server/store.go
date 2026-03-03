package server

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go-interview-practice/challenge15/oauth"
	"time"

	"gorm.io/gorm"
)

// hashToken returns the hex-encoded SHA-256 of rawToken. Only hashes are
// persisted in the database; raw tokens are only ever held in memory and
// returned to the client.
func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// RegisterClient persists a new client. Returns an error if the ClientID already exists.
func (s *OAuth2Server) RegisterClient(c *oauth.Client) error {
	return s.db.Create(c).Error
}

// EnsureClient inserts the client if it does not already exist, and is a no-op otherwise.
func (s *OAuth2Server) EnsureClient(c *oauth.Client) error {
	return s.db.Where("client_id = ?", c.ClientID).FirstOrCreate(c).Error
}

// GetClient fetches a client by ID.
func (s *OAuth2Server) GetClient(clientID string) (*oauth.Client, error) {
	var c oauth.Client
	if err := s.db.First(&c, "client_id = ?", clientID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// StoreAuthCode persists an authorization code.
func (s *OAuth2Server) StoreAuthCode(code *oauth.AuthCode) error {
	return s.db.Create(code).Error
}

// ConsumeAuthCode fetches and atomically deletes an authorization code (one-time use).
func (s *OAuth2Server) ConsumeAuthCode(code string) (*oauth.AuthCode, error) {
	var ac oauth.AuthCode
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&ac, "code = ?", code).Error; err != nil {
			return err
		}
		return tx.Delete(&ac).Error
	})
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

// IssueTokens atomically persists an access token and a refresh token.
func (s *OAuth2Server) IssueTokens(at *oauth.AccessToken, rt *oauth.RefreshToken) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		atRow := *at
		atRow.Token = hashToken(at.Token)
		if err := tx.Create(&atRow).Error; err != nil {
			return err
		}
		rtRow := *rt
		rtRow.Token = hashToken(rt.Token)
		return tx.Create(&rtRow).Error
	})
}

// GetAccessToken looks up an access token by its raw value.
func (s *OAuth2Server) GetAccessToken(rawToken string) (*oauth.AccessToken, error) {
	var at oauth.AccessToken
	if err := s.db.First(&at, "token = ?", hashToken(rawToken)).Error; err != nil {
		return nil, err
	}
	return &at, nil
}

// GetRefreshToken looks up a refresh token by its raw value.
func (s *OAuth2Server) GetRefreshToken(rawToken string) (*oauth.RefreshToken, error) {
	var rt oauth.RefreshToken
	if err := s.db.First(&rt, "token = ?", hashToken(rawToken)).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// DeleteAccessToken removes an access token identified by its raw value.
// It is idempotent: deleting a non-existent token is not an error.
func (s *OAuth2Server) DeleteAccessToken(rawToken string) error {
	return s.db.Where("token = ?", hashToken(rawToken)).Delete(&oauth.AccessToken{}).Error
}

// DeleteRefreshToken removes a refresh token identified by its raw value.
// It is idempotent: deleting a non-existent token is not an error.
func (s *OAuth2Server) DeleteRefreshToken(rawToken string) error {
	return s.db.Where("token = ?", hashToken(rawToken)).Delete(&oauth.RefreshToken{}).Error
}

// RotateRefreshToken atomically replaces the old refresh token with a new
// access+refresh token pair. Returns an error if the old token no longer exists,
// which prevents replay attacks from concurrent requests.
func (s *OAuth2Server) RotateRefreshToken(
	oldRaw string,
	newRT *oauth.RefreshToken,
	newAT *oauth.AccessToken,
) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("token = ?", hashToken(oldRaw)).Delete(&oauth.RefreshToken{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("refresh token already consumed")
		}
		atRow := *newAT
		atRow.Token = hashToken(newAT.Token)
		if err := tx.Create(&atRow).Error; err != nil {
			return err
		}
		rtRow := *newRT
		rtRow.Token = hashToken(newRT.Token)
		return tx.Create(&rtRow).Error
	})
}

// ValidateToken returns the access token record if the raw token exists and has
// not expired.
func (s *OAuth2Server) ValidateToken(rawToken string) (*oauth.AccessToken, error) {
	at, err := s.GetAccessToken(rawToken)
	if err != nil {
		return nil, err
	}
	if time.Now().After(at.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return at, nil
}
