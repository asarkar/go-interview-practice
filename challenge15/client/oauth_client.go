package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth2Config holds the client-side configuration for talking to an OAuth2 server.
type OAuth2Config struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
	IntrospectEndpoint    string
	RevokeEndpoint        string
	LogoutEndpoint        string
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scopes                []string
}

// OAuth2Client is a simple OAuth2 client that can drive the authorization code
// flow, refresh tokens, and make authenticated requests.
type OAuth2Client struct {
	Config       OAuth2Config
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
	httpClient   *http.Client
}

// NewOAuth2Client creates a new OAuth2Client with the given configuration.
func NewOAuth2Client(config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{Config: config, httpClient: &http.Client{}}
}

// GetAuthorizationURL builds the URL to redirect the user-agent to for authorization.
func (c *OAuth2Client) GetAuthorizationURL(state, challenge, method string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {c.Config.ClientID},
		"redirect_uri":          {c.Config.RedirectURI},
		"scope":                 {strings.Join(c.Config.Scopes, " ")},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {method},
	}
	return c.Config.AuthorizationEndpoint + "?" + params.Encode()
}

// ExchangeCodeForToken exchanges an authorization code for access/refresh tokens.
func (c *OAuth2Client) ExchangeCodeForToken(code, verifier string) error {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {c.Config.RedirectURI},
		"client_id":     {c.Config.ClientID},
		"client_secret": {c.Config.ClientSecret},
		"code_verifier": {verifier},
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		c.Config.TokenEndpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token exchange failed with status %s", resp.Status)
	}
	return c.decodeTokenResponse(resp)
}

// DoRefreshToken uses the stored refresh token to obtain a new access token.
func (c *OAuth2Client) DoRefreshToken() error {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.RefreshToken},
		"client_id":     {c.Config.ClientID},
		"client_secret": {c.Config.ClientSecret},
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		c.Config.TokenEndpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed with status %s", resp.Status)
	}
	return c.decodeTokenResponse(resp)
}

// MakeAuthenticatedRequest performs an HTTP request with a Bearer token.
func (c *OAuth2Client) MakeAuthenticatedRequest(urlStr, method string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	return c.httpClient.Do(req)
}

func (c *OAuth2Client) decodeTokenResponse(resp *http.Response) error {
	var tr struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}
	c.AccessToken = tr.AccessToken
	c.RefreshToken = tr.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return nil
}
