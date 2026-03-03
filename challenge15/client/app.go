package client

import (
	"go-interview-practice/challenge15/oauth"
	"net/http"
	"sync"
)

// App is the OAuth2 demo client application.
type App struct {
	config         oauth.OAuth2Config
	oauthClient    *oauth.OAuth2Client
	httpClient     *http.Client
	templates      templates
	pendingAuths   map[string]*pendingAuth
	pendingAuthsMu sync.RWMutex
	sessions       map[string]*sessionData
	sessionsMu     sync.RWMutex
}

// New creates a new App with the given OAuth2 config.
func New(config oauth.OAuth2Config) *App {
	return &App{
		config:       config,
		oauthClient:  oauth.NewOAuth2Client(config),
		httpClient:   &http.Client{},
		templates:    parseTemplates(),
		pendingAuths: make(map[string]*pendingAuth),
		sessions:     make(map[string]*sessionData),
	}
}
