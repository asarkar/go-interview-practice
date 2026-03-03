package client

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go-interview-practice/challenge15/oauth"
)

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

const pendingTTL = 15 * time.Minute

type pendingAuth struct {
	Verifier string
	Expires  time.Time
}

var (
	pendingAuths   = make(map[string]*pendingAuth)
	pendingAuthsMu sync.RWMutex
)

func (a *App) startOAuthFlow(w http.ResponseWriter, r *http.Request) (verifier, state, authURL string, err error) {
	verifier, err = oauth.GenerateVerifier()
	if err != nil {
		return "", "", "", err
	}
	challenge := oauth.DeriveChallenge(verifier)
	state, err = randomHex(16)
	if err != nil {
		return "", "", "", err
	}

	pendingAuthsMu.Lock()
	pendingAuths[state] = &pendingAuth{Verifier: verifier, Expires: time.Now().Add(pendingTTL)}
	pendingAuthsMu.Unlock()

	authURL = a.oauthClient.GetAuthorizationURL(state, challenge, "S256")
	return verifier, state, authURL, nil
}

func (a *App) consumePendingAuth(state string) (verifier string, ok bool) {
	pendingAuthsMu.Lock()
	defer pendingAuthsMu.Unlock()
	p, exists := pendingAuths[state]
	if !exists || time.Now().After(p.Expires) {
		return "", false
	}
	delete(pendingAuths, state)
	return p.Verifier, true
}

// introspectToken calls the introspect endpoint configured in OAuth2Config.
func (a *App) introspectToken(accessToken string) (map[string]any, error) {
	form := url.Values{
		"token":         {accessToken},
		"client_id":     {a.config.ClientID},
		"client_secret": {a.config.ClientSecret},
	}
	resp, err := a.httpClient.PostForm(a.config.IntrospectEndpoint, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// revokeToken sends a best-effort revocation request; errors are silently ignored
// per RFC 7009 (clients cannot rely on revocation being synchronous).
func (a *App) revokeToken(token, tokenTypeHint string) {
	form := url.Values{
		"token":           {token},
		"token_type_hint": {tokenTypeHint},
		"client_id":       {a.config.ClientID},
		"client_secret":   {a.config.ClientSecret},
	}
	resp, err := a.httpClient.PostForm(a.config.RevokeEndpoint, form)
	if err == nil {
		resp.Body.Close()
	}
}

// refreshSession creates a fresh token pair from the stored refresh token using
// a local OAuth2Client so that concurrent requests do not race on shared state.
func (a *App) refreshSession(sess *sessionData) (*sessionData, error) {
	local := oauth.NewOAuth2Client(a.config)
	local.AccessToken = sess.AccessToken
	local.RefreshToken = sess.RefreshToken
	if err := local.DoRefreshToken(); err != nil {
		return nil, err
	}
	return &sessionData{
		AccessToken:  local.AccessToken,
		RefreshToken: local.RefreshToken,
		ExpiresAt:    local.TokenExpiry,
	}, nil
}

// exchangeCode exchanges an authorization code for tokens using a local client,
// avoiding mutation of any shared state.
func (a *App) exchangeCode(code, verifier string) (*sessionData, error) {
	local := oauth.NewOAuth2Client(a.config)
	if err := local.ExchangeCodeForToken(code, verifier); err != nil {
		return nil, err
	}
	return &sessionData{
		AccessToken:  local.AccessToken,
		RefreshToken: local.RefreshToken,
		ExpiresAt:    local.TokenExpiry,
	}, nil
}

