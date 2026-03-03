package oauth

import "time"

// Client is a registered OAuth2 client application.
type Client struct {
	ClientID      string   `gorm:"primaryKey"`
	ClientSecret  string
	RedirectURIs  []string `gorm:"serializer:json"`
	AllowedScopes []string `gorm:"serializer:json"`
}

// AuthCode is an issued authorization code.
type AuthCode struct {
	Code                string   `gorm:"primaryKey"`
	ClientID            string
	UserID              string
	RedirectURI         string
	Scopes              []string `gorm:"serializer:json"`
	ExpiresAt           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
}

// AccessToken is an issued OAuth2 access token.
type AccessToken struct {
	Token     string   `gorm:"primaryKey"`
	ClientID  string
	UserID    string
	Scopes    []string `gorm:"serializer:json"`
	ExpiresAt time.Time
}

// RefreshToken is an issued OAuth2 refresh token.
type RefreshToken struct {
	Token     string   `gorm:"primaryKey"`
	ClientID  string
	UserID    string
	Scopes    []string `gorm:"serializer:json"`
	ExpiresAt time.Time
}

// User represents an end-user in the system.
type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string
}
