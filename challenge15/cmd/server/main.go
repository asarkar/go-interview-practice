package main

import (
	"go-interview-practice/challenge15/oauth"
	"go-interview-practice/challenge15/server"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("oauth2.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	srv := server.NewOAuth2Server(db)

	for _, c := range []*oauth.Client{
		{
			ClientID:      "demo",
			ClientSecret:  "secret",
			RedirectURIs:  []string{"http://localhost/cb"},
			AllowedScopes: []string{"read", "write"},
		},
		{
			ClientID:      "demo-client",
			ClientSecret:  "secret",
			RedirectURIs:  []string{"http://localhost:8081/callback"},
			AllowedScopes: []string{"read", "write", "profile"},
		},
	} {
		if err := srv.EnsureClient(c); err != nil {
			log.Fatalf("failed to seed client %s: %v", c.ClientID, err)
		}
	}

	log.Println("OAuth2 server listening on :8080")
	log.Println("  GET  /authorize")
	log.Println("  POST /token")
	log.Println("  POST /revoke")
	log.Println("  POST /introspect")
	log.Println("  GET  /api/userinfo (protected)")
	if err := srv.StartServer(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
