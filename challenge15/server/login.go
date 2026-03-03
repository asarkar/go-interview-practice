package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"go-interview-practice/challenge15/oauth"
	"golang.org/x/crypto/bcrypt"
)

const loginCookieName = "auth_session"
const loginSessionTTL = 24 * time.Hour
const csrfTokenTTL = 15 * time.Minute

type loginSession struct {
	userID    string
	expiresAt time.Time
}

var (
	loginSessions   = make(map[string]loginSession)
	loginSessionsMu sync.RWMutex

	csrfTokens   = make(map[string]time.Time)
	csrfTokensMu sync.RWMutex
)

func (s *OAuth2Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}
	if r.Method == http.MethodPost {
		s.handleLoginPost(w, r, redirect)
		return
	}
	s.renderLoginFormWithNewCSRF(w, redirect)
}

func (s *OAuth2Server) handleLoginPost(w http.ResponseWriter, r *http.Request, redirect string) {
	r.ParseForm()

	// Validate and consume the CSRF token (synchronizer token pattern).
	csrfToken := r.FormValue("csrf_token")
	csrfTokensMu.Lock()
	expiry, ok := csrfTokens[csrfToken]
	if ok {
		delete(csrfTokens, csrfToken)
	}
	csrfTokensMu.Unlock()
	if !ok || time.Now().After(expiry) {
		http.Error(w, "invalid CSRF token", http.StatusForbidden)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		s.renderLoginFormWithNewCSRF(w, redirect)
		return
	}

	var user oauth.User
	if err := s.db.First(&user, "username = ?", username).Error; err != nil {
		s.renderLoginFormWithNewCSRF(w, redirect)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.renderLoginFormWithNewCSRF(w, redirect)
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	loginSessionsMu.Lock()
	loginSessions[sessionID] = loginSession{
		userID:    user.ID,
		expiresAt: time.Now().Add(loginSessionTTL),
	}
	loginSessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     loginCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(loginSessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (s *OAuth2Server) renderLoginFormWithNewCSRF(w http.ResponseWriter, redirect string) {
	csrf, err := generateSessionID() // same 16-byte random hex generator
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	csrfTokensMu.Lock()
	csrfTokens[csrf] = time.Now().Add(csrfTokenTTL)
	csrfTokensMu.Unlock()
	s.renderLoginForm(w, redirect, csrf)
}

func (s *OAuth2Server) renderLoginForm(w http.ResponseWriter, redirect, csrfToken string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginTemplate.Execute(w, map[string]string{
		"Redirect":  redirect,
		"CSRFToken": csrfToken,
	})
}

// HandleLogout clears the server-side login session and redirects to
// post_logout_redirect_uri (defaults to "/").
func (s *OAuth2Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(loginCookieName); err == nil && cookie.Value != "" {
		loginSessionsMu.Lock()
		delete(loginSessions, cookie.Value)
		loginSessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     loginCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	redirect := r.URL.Query().Get("post_logout_redirect_uri")
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

// getUserFromLogin resolves the authenticated user from the login session cookie.
// There is no test-only bypass; tests inject the user via the context key instead.
func (s *OAuth2Server) getUserFromLogin(r *http.Request) string {
	cookie, err := r.Cookie(loginCookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}
	loginSessionsMu.RLock()
	sess, ok := loginSessions[cookie.Value]
	loginSessionsMu.RUnlock()
	if !ok || time.Now().After(sess.expiresAt) {
		return ""
	}
	return sess.userID
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
