package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateVerifier creates a cryptographically random PKCE code verifier.
func GenerateVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// DeriveChallenge computes the S256 code challenge for a given verifier.
func DeriveChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// VerifyChallenge checks that verifier matches challenge using the given method.
// Only "S256" is accepted; all other methods (including "plain") return false.
func VerifyChallenge(verifier, challenge, method string) bool {
	if method != "S256" {
		return false
	}
	return DeriveChallenge(verifier) == challenge
}
