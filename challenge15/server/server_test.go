package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go-interview-practice/challenge15/oauth"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testUserID = "user"

func newTestServer(t *testing.T) (*OAuth2Server, *httptest.Server) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}
	srv := NewOAuth2Server(db)
	ts := httptest.NewServer(injectTestUser(srv.Router()))
	t.Cleanup(ts.Close)
	return srv, ts
}

func injectTestUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uid := r.Header.Get("X-Test-User-ID"); uid != "" {
			r = r.WithContext(context.WithValue(r.Context(), userIDKey, uid))
		}
		next.ServeHTTP(w, r)
	})
}

func registerTestClient(t *testing.T, srv *OAuth2Server) {
	t.Helper()
	err := srv.RegisterClient(&Client{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	})
	if err != nil {
		t.Fatalf("failed to register test client: %v", err)
	}
}

func noFollow() *http.Client {
	return &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func authorizeAndGetCode(t *testing.T, ts *httptest.Server, clientID, challenge string) string {
	t.Helper()
	reqURL := fmt.Sprintf(
		"%s/authorize?response_type=code&client_id=%s&redirect_uri=https://client.example.com/callback&scope=read+profile&state=teststate&code_challenge=%s&code_challenge_method=S256",
		ts.URL,
		clientID,
		url.QueryEscape(challenge),
	)
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("X-Test-User-ID", testUserID)

	resp, err := noFollow().Do(req)
	if err != nil {
		t.Fatalf("authorize request failed: %v", err)
	}
	loc, _ := url.Parse(resp.Header.Get("Location"))
	code := loc.Query().Get("code")
	if code == "" {
		t.Fatalf(
			"expected authorization code in redirect, got location: %s",
			resp.Header.Get("Location"),
		)
	}
	return code
}

func exchangeCode(t *testing.T, ts *httptest.Server, code, verifier string) (string, string) {
	t.Helper()
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://client.example.com/callback"},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
		"code_verifier": {verifier},
	}
	resp, err := http.PostForm(ts.URL+"/token", form)
	if err != nil {
		t.Fatalf("token exchange failed: %v", err)
	}
	var tr struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(resp.Body).Decode(&tr)
	if tr.AccessToken == "" || tr.RefreshToken == "" {
		t.Fatalf(
			"exchangeCode: expected both tokens, got access=%q refresh=%q",
			tr.AccessToken,
			tr.RefreshToken,
		)
	}
	return tr.AccessToken, tr.RefreshToken
}

func TestClientRegistration(t *testing.T) {
	srv, _ := newTestServer(t)

	client := &Client{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	}
	if err := srv.RegisterClient(client); err != nil {
		t.Fatalf("failed to register client: %v", err)
	}

	dup := &Client{
		ClientID:      "test-client",
		ClientSecret:  "different-secret",
		RedirectURIs:  []string{"https://other.example.com/callback"},
		AllowedScopes: []string{"read"},
	}
	if err := srv.RegisterClient(dup); err == nil {
		t.Fatal("expected error for duplicate client ID, got nil")
	}
}

func TestGenerateVerifier(t *testing.T) {
	v1, err := oauth.GenerateVerifier()
	if err != nil {
		t.Fatalf("GenerateVerifier: %v", err)
	}
	if len(v1) == 0 {
		t.Fatal("expected non-empty verifier")
	}
	v2, err := oauth.GenerateVerifier()
	if err != nil {
		t.Fatalf("GenerateVerifier (2nd): %v", err)
	}
	if v1 == v2 {
		t.Error("expected two distinct verifiers, got the same string")
	}
}

func TestVerifyChallenge(t *testing.T) {
	tests := []struct {
		name      string
		verifier  string
		challenge string
		method    string
		want      bool
	}{
		{"S256 valid", "test-verifier", oauth.DeriveChallenge("test-verifier"), "S256", true},
		{"S256 wrong verifier", "test-verifier", "bad-challenge", "S256", false},
		{"plain rejected", "test-verifier", "test-verifier", "plain", false},
		{"unsupported method", "test-verifier", "test-verifier", "md5", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := oauth.VerifyChallenge(tt.verifier, tt.challenge, tt.method)
			if got != tt.want {
				t.Errorf("VerifyChallenge(%q, %q, %q) = %v, want %v",
					tt.verifier, tt.challenge, tt.method, got, tt.want)
			}
		})
	}
}

func TestAuthorizationEndpoint(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)
	hc := noFollow()

	t.Run("ValidRequest", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		reqURL := fmt.Sprintf(
			"%s/authorize?response_type=code&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz123&code_challenge=%s&code_challenge_method=S256",
			ts.URL,
			url.QueryEscape(challenge),
		)
		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header.Set("X-Test-User-ID", testUserID)

		resp, err := hc.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected 302, got %d", resp.StatusCode)
		}
		loc, _ := url.Parse(resp.Header.Get("Location"))
		if loc.Query().Get("code") == "" {
			t.Error("expected code in redirect location")
		}
		if loc.Query().Get("state") != "xyz123" {
			t.Errorf("expected state=xyz123, got %q", loc.Query().Get("state"))
		}
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/authorize?response_type=code&client_id=unknown&redirect_uri=https://client.example.com/callback&scope=read&state=xyz",
			nil,
		)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidRedirectURI", func(t *testing.T) {
		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/authorize?response_type=code&client_id=test-client&redirect_uri=https://attacker.example.com/callback&scope=read&state=xyz",
			nil,
		)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidResponseType", func(t *testing.T) {
		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/authorize?response_type=token&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz",
			nil,
		)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected 302, got %d", resp.StatusCode)
		}
		loc, _ := url.Parse(resp.Header.Get("Location"))
		if loc.Query().Get("error") != "unsupported_response_type" {
			t.Errorf("expected error=unsupported_response_type, got %q", loc.Query().Get("error"))
		}
	})

	t.Run("InvalidScope", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		reqURL := fmt.Sprintf(
			"%s/authorize?response_type=code&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=admin&state=xyz&code_challenge=%s&code_challenge_method=S256",
			ts.URL,
			url.QueryEscape(challenge),
		)
		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected 302, got %d", resp.StatusCode)
		}
		loc, _ := url.Parse(resp.Header.Get("Location"))
		if loc.Query().Get("error") != "invalid_scope" {
			t.Errorf("expected error=invalid_scope, got %q", loc.Query().Get("error"))
		}
	})

	t.Run("MissingPKCE", func(t *testing.T) {
		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/authorize?response_type=code&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz",
			nil,
		)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected 302, got %d", resp.StatusCode)
		}
		loc, _ := url.Parse(resp.Header.Get("Location"))
		if loc.Query().Get("error") != "invalid_request" {
			t.Errorf(
				"expected error=invalid_request for missing PKCE, got %q",
				loc.Query().Get("error"),
			)
		}
	})

	t.Run("PlainPKCEMethod", func(t *testing.T) {
		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/authorize?response_type=code&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz&code_challenge=abc&code_challenge_method=plain",
			nil,
		)
		req.Header.Set("X-Test-User-ID", testUserID)
		resp, _ := hc.Do(req)
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected 302, got %d", resp.StatusCode)
		}
		loc, _ := url.Parse(resp.Header.Get("Location"))
		if loc.Query().Get("error") != "invalid_request" {
			t.Errorf(
				"expected error=invalid_request for plain PKCE, got %q",
				loc.Query().Get("error"),
			)
		}
	})
}

func TestTokenEndpoint(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	t.Run("ValidAuthorizationCode", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)

		form := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {"https://client.example.com/callback"},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
			"code_verifier": {verifier},
		}
		resp, err := http.PostForm(ts.URL+"/token", form)
		if err != nil {
			t.Fatalf("token request failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		var tr struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
		}
		json.NewDecoder(resp.Body).Decode(&tr)

		if tr.AccessToken == "" {
			t.Error("expected access_token")
		}
		if tr.TokenType != "Bearer" {
			t.Errorf("expected token_type=Bearer, got %q", tr.TokenType)
		}
		if tr.ExpiresIn <= 0 {
			t.Errorf("expected positive expires_in, got %d", tr.ExpiresIn)
		}
		if tr.RefreshToken == "" {
			t.Error("expected refresh_token")
		}
		if tr.Scope != "read profile" {
			t.Errorf("expected scope=\"read profile\", got %q", tr.Scope)
		}

		resp2, _ := http.PostForm(ts.URL+"/token", form)
		if resp2.StatusCode != http.StatusBadRequest {
			t.Errorf("replayed code: expected 400, got %d", resp2.StatusCode)
		}
	})

	t.Run("InvalidClientCredentials", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)

		form := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {"https://client.example.com/callback"},
			"client_id":     {"test-client"},
			"client_secret": {"wrong-secret"},
			"code_verifier": {verifier},
		}
		resp, _ := http.PostForm(ts.URL+"/token", form)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
		var er struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&er)
		if er.Error != "invalid_client" {
			t.Errorf("expected error=invalid_client, got %q", er.Error)
		}
	})

	t.Run("InvalidCodeVerifier", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)

		form := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {"https://client.example.com/callback"},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
			"code_verifier": {"wrong-verifier"},
		}
		resp, _ := http.PostForm(ts.URL+"/token", form)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
		var er struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&er)
		if er.Error != "invalid_grant" {
			t.Errorf("expected error=invalid_grant, got %q", er.Error)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	t.Run("ValidRefreshToken", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)
		_, refreshToken := exchangeCode(t, ts, code, verifier)

		form := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {refreshToken},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/token", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		var tr struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
		}
		json.NewDecoder(resp.Body).Decode(&tr)

		if tr.AccessToken == "" {
			t.Error("expected access_token")
		}
		if tr.RefreshToken == "" {
			t.Error("expected refresh_token")
		}
		if tr.RefreshToken == refreshToken {
			t.Error("expected rotated (different) refresh token")
		}
		if tr.Scope != "read profile" {
			t.Errorf("expected scope=\"read profile\", got %q", tr.Scope)
		}

		form2 := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {refreshToken},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp2, _ := http.PostForm(ts.URL+"/token", form2)
		if resp2.StatusCode != http.StatusBadRequest {
			t.Errorf("old refresh token: expected 400, got %d", resp2.StatusCode)
		}
	})

	t.Run("InvalidClientForRefreshToken", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)
		_, refreshToken := exchangeCode(t, ts, code, verifier)

		form := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {refreshToken},
			"client_id":     {"test-client"},
			"client_secret": {"wrong-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/token", form)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}

		if _, err := srv.GetToken(refreshToken, TokenTypeRefresh); err != nil {
			t.Error("refresh token should still be valid after failed client authentication")
		}
	})
}

func TestRefreshTokenReplay(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	verifier, _ := oauth.GenerateVerifier()
	challenge := oauth.DeriveChallenge(verifier)
	code := authorizeAndGetCode(t, ts, "test-client", challenge)
	_, refreshToken := exchangeCode(t, ts, code, verifier)

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
	}

	// First use must succeed.
	resp1, _ := http.PostForm(ts.URL+"/token", form)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("first refresh: expected 200, got %d", resp1.StatusCode)
	}
	resp1.Body.Close()

	// Second use of the same (now-consumed) refresh token must fail.
	resp2, _ := http.PostForm(ts.URL+"/token", form)
	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("replay refresh: expected 400, got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
}

func TestTokenValidation(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	verifier, _ := oauth.GenerateVerifier()
	challenge := oauth.DeriveChallenge(verifier)
	code := authorizeAndGetCode(t, ts, "test-client", challenge)
	accessToken, _ := exchangeCode(t, ts, code, verifier)

	t.Run("ValidToken", func(t *testing.T) {
		at, err := srv.ValidateToken(accessToken)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}
		if at.UserID != testUserID {
			t.Errorf("expected UserID=user, got %q", at.UserID)
		}
		if at.ClientID != "test-client" {
			t.Errorf("expected ClientID=test-client, got %q", at.ClientID)
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		expired := &Token{
			Token:     hashToken("expired-raw-token"),
			Type:      TokenTypeAccess,
			ClientID:  "test-client",
			UserID:    testUserID,
			Scopes:    []string{"read"},
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		srv.db.Create(expired)

		_, err := srv.ValidateToken("expired-raw-token")
		if err == nil {
			t.Fatal("expected error for expired token, got nil")
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := srv.ValidateToken("non-existent-token")
		if err == nil {
			t.Fatal("expected error for non-existent token, got nil")
		}
	})
}

func TestTokenRevocation(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	t.Run("RevokeAccessToken", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)
		accessToken, _ := exchangeCode(t, ts, code, verifier)

		form := url.Values{
			"token":         {accessToken},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/revoke", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		if _, err := srv.GetToken(accessToken, TokenTypeAccess); err == nil {
			t.Error("access token should be revoked")
		}
	})

	t.Run("RevokeRefreshToken", func(t *testing.T) {
		verifier, _ := oauth.GenerateVerifier()
		challenge := oauth.DeriveChallenge(verifier)
		code := authorizeAndGetCode(t, ts, "test-client", challenge)
		_, refreshToken := exchangeCode(t, ts, code, verifier)

		form := url.Values{
			"token":         {refreshToken},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/revoke", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		if _, err := srv.GetToken(refreshToken, TokenTypeRefresh); err == nil {
			t.Error("refresh token should be revoked")
		}
	})

	t.Run("RevokeNonExistentToken", func(t *testing.T) {
		form := url.Values{
			"token":         {"non-existent-token"},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/revoke", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 per RFC 7009, got %d", resp.StatusCode)
		}
	})
}

func TestIntrospect(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	verifier, _ := oauth.GenerateVerifier()
	challenge := oauth.DeriveChallenge(verifier)
	code := authorizeAndGetCode(t, ts, "test-client", challenge)
	accessToken, _ := exchangeCode(t, ts, code, verifier)

	t.Run("ActiveToken", func(t *testing.T) {
		form := url.Values{
			"token":         {accessToken},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/introspect", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)
		if body["active"] != true {
			t.Errorf("expected active=true, got %v", body["active"])
		}
		if body["client_id"] != "test-client" {
			t.Errorf("expected client_id=test-client, got %v", body["client_id"])
		}
	})

	t.Run("InactiveToken", func(t *testing.T) {
		form := url.Values{
			"token":         {"bogus-token"},
			"client_id":     {"test-client"},
			"client_secret": {"test-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/introspect", form)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)
		if body["active"] != false {
			t.Errorf("expected active=false, got %v", body["active"])
		}
	})

	t.Run("InvalidClient", func(t *testing.T) {
		form := url.Values{
			"token":         {accessToken},
			"client_id":     {"test-client"},
			"client_secret": {"wrong-secret"},
		}
		resp, _ := http.PostForm(ts.URL+"/introspect", form)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})
}

func TestUserinfo(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	verifier, _ := oauth.GenerateVerifier()
	challenge := oauth.DeriveChallenge(verifier)
	code := authorizeAndGetCode(t, ts, "test-client", challenge)
	accessToken, _ := exchangeCode(t, ts, code, verifier)

	t.Run("ValidToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ts.URL+"/api/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)
		if body["user_id"] != testUserID {
			t.Errorf("expected user_id=user, got %v", body["user_id"])
		}
		if body["client_id"] != "test-client" {
			t.Errorf("expected client_id=test-client, got %v", body["client_id"])
		}
	})

	t.Run("MissingAuth", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + "/api/userinfo")
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ts.URL+"/api/userinfo", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp, _ := http.DefaultClient.Do(req)
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})
}

func TestFullPKCEFlow(t *testing.T) {
	srv, ts := newTestServer(t)
	registerTestClient(t, srv)

	verifier, err := oauth.GenerateVerifier()
	if err != nil {
		t.Fatalf("GenerateVerifier: %v", err)
	}
	challenge := oauth.DeriveChallenge(verifier)

	code := authorizeAndGetCode(t, ts, "test-client", challenge)
	accessToken, refreshToken := exchangeCode(t, ts, code, verifier)

	at, err := srv.ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if at.UserID != testUserID {
		t.Errorf("expected UserID=user, got %q", at.UserID)
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
	}
	resp, _ := http.PostForm(ts.URL+"/token", form)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d", resp.StatusCode)
	}
	var tr struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(resp.Body).Decode(&tr)

	iForm := url.Values{
		"token":         {tr.AccessToken},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
	}
	iResp, _ := http.PostForm(ts.URL+"/introspect", iForm)
	var iBody map[string]any
	json.NewDecoder(iResp.Body).Decode(&iBody)
	if iBody["active"] != true {
		t.Errorf("expected active=true after refresh, got %v", iBody["active"])
	}

	rForm := url.Values{
		"token":         {tr.AccessToken},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
	}
	rResp, _ := http.PostForm(ts.URL+"/revoke", rForm)
	if rResp.StatusCode != http.StatusOK {
		t.Errorf("revoke: expected 200, got %d", rResp.StatusCode)
	}

	iForm2 := url.Values{
		"token":         {tr.AccessToken},
		"client_id":     {"test-client"},
		"client_secret": {"test-secret"},
	}
	iResp2, _ := http.PostForm(ts.URL+"/introspect", iForm2)
	var iBody2 map[string]any
	json.NewDecoder(iResp2.Body).Decode(&iBody2)
	if iBody2["active"] != false {
		t.Errorf("expected active=false after revoke, got %v", iBody2["active"])
	}
}
