package server

import "time"

// Client is a registered OAuth2 client application.
type Client struct {
	ClientID      string   `gorm:"primaryKey"`
	ClientSecret  string   //nolint:gosec // G117
	RedirectURIs  []string `gorm:"serializer:json"`
	AllowedScopes []string `gorm:"serializer:json"`
}

// AuthCode is an issued authorization code.
type AuthCode struct {
	Code                string `gorm:"primaryKey"`
	ClientID            string
	UserID              string
	RedirectURI         string
	Scopes              []string `gorm:"serializer:json"`
	ExpiresAt           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
}

// Token type constants.
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Token is an issued OAuth2 token (access or refresh).
type Token struct {
	Token     string `gorm:"primaryKey"`
	Type      string `gorm:"index;size:16"` // TokenTypeAccess or TokenTypeRefresh
	ClientID  string
	UserID    string
	Scopes    []string `gorm:"serializer:json"`
	ExpiresAt time.Time
}

// User represents an end-user in the system.
type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string //nolint:gosec // G117
}
