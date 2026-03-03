package client

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router returns an http.Handler with all client routes registered.
func (a *App) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", a.handleHome)
	r.Get("/login", a.handleLogin)
	r.Get("/callback", a.handleCallback)
	r.Get("/logout", a.handleLogout)
	r.Get("/api/me", a.handleMe)

	return r
}
