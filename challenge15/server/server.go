package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-interview-practice/challenge15/oauth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// OAuth2Server is the OAuth2 authorization server.
type OAuth2Server struct {
	db *gorm.DB
}

// NewOAuth2Server initialises the server, auto-migrates the schema, and seeds a
// demo user with a bcrypt-hashed password.
func NewOAuth2Server(db *gorm.DB) *OAuth2Server {
	db.AutoMigrate(&oauth.Client{}, &oauth.AuthCode{}, &oauth.AccessToken{}, &oauth.RefreshToken{}, &oauth.User{})
	s := &OAuth2Server{db: db}
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash demo user password: %v", err)
	}
	s.db.Where("id = ?", "user").FirstOrCreate(&oauth.User{
		ID:       "user",
		Username: "user",
		Password: string(hash),
	})
	return s
}

// Router returns the chi router with all OAuth2 endpoints registered.
func (s *OAuth2Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/login", s.HandleLogin)
	r.Post("/login", s.HandleLogin)
	r.Get("/logout", s.HandleLogout)
	r.Get("/authorize", s.HandleAuthorize)
	r.Post("/token", s.HandleToken)
	r.Post("/revoke", s.HandleRevoke)
	r.Post("/introspect", s.HandleIntrospect)
	r.Get("/api/userinfo", s.HandleUserinfo)
	return r
}

// StartServer starts the HTTP server on addr (e.g. ":8080").
func (s *OAuth2Server) StartServer(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
