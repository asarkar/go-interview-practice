package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-interview-practice/challenge15/oauth"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const userIDKey contextKey = "user_id"

// Token TTL constants.
const (
	accessTokenTTL  = 1 * time.Hour
	refreshTokenTTL = 24 * time.Hour
	authCodeTTL     = 10 * time.Minute
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type errResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func writeError(w http.ResponseWriter, status int, code, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResponse{Error: code, Description: desc})
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *OAuth2Server) authenticateClient(r *http.Request) (*oauth.Client, error) {
	var clientID, clientSecret string
	if id, secret, ok := r.BasicAuth(); ok {
		clientID, clientSecret = id, secret
	} else {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}
	client, err := s.GetClient(clientID)
	if err != nil {
		return nil, errors.New("invalid_client")
	}
	if subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(clientSecret)) != 1 {
		return nil, errors.New("invalid_client")
	}
	return client, nil
}

// HandleAuthorize implements GET /authorize (authorization code grant).
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	responseType := q.Get("response_type")
	scope := q.Get("scope")
	state := q.Get("state")
	codeChallenge := q.Get("code_challenge")
	codeChallengeMethod := q.Get("code_challenge_method")

	client, err := s.GetClient(clientID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_client", "unknown client")
		return
	}

	validRedirect := false
	for _, uri := range client.RedirectURIs {
		if uri == redirectURI {
			validRedirect = true
			break
		}
	}
	if !validRedirect {
		writeError(w, http.StatusBadRequest, "invalid_redirect_uri", "redirect URI not registered")
		return
	}

	if responseType != "code" {
		redirectWithError(w, r, redirectURI, "unsupported_response_type", state)
		return
	}

	var grantedScopes []string
	if scope != "" {
		allowedSet := make(map[string]bool, len(client.AllowedScopes))
		for _, sc := range client.AllowedScopes {
			allowedSet[sc] = true
		}
		for _, sc := range strings.Split(scope, " ") {
			if sc != "" && allowedSet[sc] {
				grantedScopes = append(grantedScopes, sc)
			}
		}
		if len(grantedScopes) == 0 {
			redirectWithError(w, r, redirectURI, "invalid_scope", state)
			return
		}
	}

	if codeChallenge == "" || codeChallengeMethod != "S256" {
		redirectWithError(w, r, redirectURI, "invalid_request", state)
		return
	}

	userID, _ := r.Context().Value(userIDKey).(string)
	if userID == "" {
		userID = s.getUserFromLogin(r)
	}
	if userID == "" {
		loginURL := "/login?redirect=" + url.QueryEscape(r.URL.RequestURI())
		http.Redirect(w, r, loginURL, http.StatusFound)
		return
	}

	codeBytes := make([]byte, 32)
	if _, err := rand.Read(codeBytes); err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to generate code")
		return
	}
	code := hex.EncodeToString(codeBytes)

	authCode := &oauth.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              grantedScopes,
		ExpiresAt:           time.Now().Add(authCodeTTL),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	if err := s.StoreAuthCode(authCode); err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to store code")
		return
	}

	redirectURL := redirectURI + "?code=" + code
	if state != "" {
		redirectURL += "&state=" + state
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func redirectWithError(w http.ResponseWriter, r *http.Request, redirectURI, errCode, state string) {
	u := redirectURI + "?error=" + errCode
	if state != "" {
		u += "&state=" + state
	}
	http.Redirect(w, r, u, http.StatusFound)
}

// HandleToken implements POST /token (authorization_code and refresh_token grants).
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.FormValue("grant_type") {
	case "authorization_code":
		s.handleAuthCodeGrant(w, r)
	case "refresh_token":
		s.handleRefreshTokenGrant(w, r)
	default:
		writeError(w, http.StatusBadRequest, "unsupported_grant_type", "unsupported grant type")
	}
}

func (s *OAuth2Server) handleAuthCodeGrant(w http.ResponseWriter, r *http.Request) {
	client, err := s.authenticateClient(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_client", "invalid client credentials")
		return
	}

	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	codeVerifier := r.FormValue("code_verifier")

	authCode, err := s.ConsumeAuthCode(code)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_grant", "invalid or already-used authorization code")
		return
	}

	if time.Now().After(authCode.ExpiresAt) {
		writeError(w, http.StatusBadRequest, "invalid_grant", "authorization code expired")
		return
	}

	if authCode.RedirectURI != redirectURI {
		writeError(w, http.StatusBadRequest, "invalid_grant", "redirect URI mismatch")
		return
	}

	if authCode.CodeChallenge != "" {
		if !oauth.VerifyChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
			writeError(w, http.StatusBadRequest, "invalid_grant", "invalid code_verifier")
			return
		}
	}

	atStr, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to generate access token")
		return
	}
	rtStr, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to generate refresh token")
		return
	}

	at := &oauth.AccessToken{Token: atStr, ClientID: client.ClientID, UserID: authCode.UserID, Scopes: authCode.Scopes, ExpiresAt: time.Now().Add(accessTokenTTL)}
	rt := &oauth.RefreshToken{Token: rtStr, ClientID: client.ClientID, UserID: authCode.UserID, Scopes: authCode.Scopes, ExpiresAt: time.Now().Add(refreshTokenTTL)}

	if err := s.IssueTokens(at, rt); err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to persist tokens")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{
		AccessToken:  atStr,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenTTL.Seconds()),
		RefreshToken: rtStr,
		Scope:        strings.Join(authCode.Scopes, " "),
	})
}

func (s *OAuth2Server) handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request) {
	client, err := s.authenticateClient(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_client", "invalid client credentials")
		return
	}

	rawRT := r.FormValue("refresh_token")
	rt, err := s.GetRefreshToken(rawRT)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_grant", "refresh token not found")
		return
	}

	if rt.ClientID != client.ClientID {
		writeError(w, http.StatusUnauthorized, "invalid_grant", "token does not belong to this client")
		return
	}

	if time.Now().After(rt.ExpiresAt) {
		writeError(w, http.StatusBadRequest, "invalid_grant", "refresh token expired")
		return
	}

	atStr, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to generate access token")
		return
	}
	rtStr, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to generate refresh token")
		return
	}

	newAT := &oauth.AccessToken{Token: atStr, ClientID: client.ClientID, UserID: rt.UserID, Scopes: rt.Scopes, ExpiresAt: time.Now().Add(accessTokenTTL)}
	newRT := &oauth.RefreshToken{Token: rtStr, ClientID: client.ClientID, UserID: rt.UserID, Scopes: rt.Scopes, ExpiresAt: time.Now().Add(refreshTokenTTL)}

	if err := s.RotateRefreshToken(rawRT, newRT, newAT); err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "failed to rotate tokens")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{
		AccessToken:  atStr,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenTTL.Seconds()),
		RefreshToken: rtStr,
		Scope:        strings.Join(rt.Scopes, " "),
	})
}

// HandleRevoke implements POST /revoke (RFC 7009).
func (s *OAuth2Server) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if _, err := s.authenticateClient(r); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_client", "invalid client credentials")
		return
	}

	token := r.FormValue("token")
	hint := r.FormValue("token_type_hint")
	switch hint {
	case "access_token":
		if s.DeleteAccessToken(token) != nil {
			s.DeleteRefreshToken(token) //nolint:errcheck
		}
	case "refresh_token":
		if s.DeleteRefreshToken(token) != nil {
			s.DeleteAccessToken(token) //nolint:errcheck
		}
	default:
		if s.DeleteAccessToken(token) != nil {
			s.DeleteRefreshToken(token) //nolint:errcheck
		}
	}
}

// HandleIntrospect implements POST /introspect (RFC 7662).
func (s *OAuth2Server) HandleIntrospect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if _, err := s.authenticateClient(r); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_client", "invalid client credentials")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	at, err := s.ValidateToken(r.FormValue("token"))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"active": false})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"active":     true,
		"client_id":  at.ClientID,
		"username":   at.UserID,
		"scope":      strings.Join(at.Scopes, " "),
		"exp":        at.ExpiresAt.Unix(),
		"token_type": "Bearer",
	})
}

// HandleUserinfo implements GET /api/userinfo — a protected resource that returns
// user info for a valid Bearer token.
func (s *OAuth2Server) HandleUserinfo(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		writeError(w, http.StatusUnauthorized, "invalid_request", "missing or invalid Authorization header")
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	at, err := s.ValidateToken(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_token", "token invalid or expired")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user_id":   at.UserID,
		"client_id": at.ClientID,
		"scope":     strings.Join(at.Scopes, " "),
		"exp":       at.ExpiresAt.Unix(),
	})
}
