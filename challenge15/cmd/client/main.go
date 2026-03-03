package main

import (
	"log"
	"net/http"

	"go-interview-practice/challenge15/client"
	"go-interview-practice/challenge15/oauth"
)

func main() {
	config := oauth.OAuth2Config{
		AuthorizationEndpoint: "http://localhost:8080/authorize",
		TokenEndpoint:         "http://localhost:8080/token",
		IntrospectEndpoint:    "http://localhost:8080/introspect",
		RevokeEndpoint:        "http://localhost:8080/revoke",
		LogoutEndpoint:        "http://localhost:8080/logout",
		ClientID:              "demo-client",
		ClientSecret:          "secret",
		RedirectURI:           "http://localhost:8081/callback",
		Scopes:                []string{"read", "write", "profile"},
	}

	app := client.New(config)

	log.Println("OAuth2 client listening on :8081")
	log.Println("  Open http://localhost:8081 (ensure auth server is running on :8080)")
	if err := http.ListenAndServe(":8081", app.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
