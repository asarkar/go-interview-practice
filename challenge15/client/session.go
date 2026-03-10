package client

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

const (
	sessionCookieName = "demo_session"
	sessionTTL        = 24 * time.Hour
)

type sessionData struct {
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time // access token expiry
	sessionExpiresAt time.Time // server-side session TTL
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *App) getSession(r *http.Request) *sessionData {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}
	a.sessionsMu.RLock()
	sess := a.sessions[cookie.Value]
	a.sessionsMu.RUnlock()
	if sess == nil || time.Now().After(sess.sessionExpiresAt) {
		return nil
	}
	return sess
}

func (a *App) setSession(w http.ResponseWriter, data *sessionData) error {
	id, err := generateSessionID()
	if err != nil {
		return err
	}
	data.sessionExpiresAt = time.Now().Add(sessionTTL)
	a.sessionsMu.Lock()
	a.sessions[id] = data
	a.sessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    id,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

// clearSession deletes the server-side session entry and expires the cookie.
func (a *App) clearSession(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		a.sessionsMu.Lock()
		delete(a.sessions, cookie.Value)
		a.sessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
