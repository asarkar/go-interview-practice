package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	sess := a.getSession(r)
	if sess == nil {
		a.templates.execute(w, "index.html", map[string]any{"LoggedIn": false})
		return
	}

	info, err := a.introspectToken(sess.AccessToken)
	if err != nil || info["active"] != true {
		newSess, err := a.refreshSession(sess)
		if err != nil {
			a.clearSession(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		a.setSession(w, newSess)
		info, _ = a.introspectToken(newSess.AccessToken)
	}

	a.templates.execute(w, "index.html", map[string]any{
		"LoggedIn": true,
		"UserID":   info["username"],
		"Scope":    info["scope"],
		"ClientID": info["client_id"],
	})
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	_, _, authURL, err := a.startOAuthFlow(w, r)
	if err != nil {
		http.Error(w, "failed to start OAuth flow", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (a *App) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		a.templates.execute(w, "error.html", map[string]any{
			"Error": "Missing code or state in callback",
		})
		return
	}

	verifier, ok := a.consumePendingAuth(state)
	if !ok {
		a.templates.execute(w, "error.html", map[string]any{
			"Error": "Invalid or expired state",
		})
		return
	}

	sess, err := a.exchangeCode(code, verifier)
	if err != nil {
		a.templates.execute(w, "error.html", map[string]any{
			"Error": "Token exchange failed: " + err.Error(),
		})
		return
	}

	if err := a.setSession(w, sess); err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	if sess := a.getSession(r); sess != nil {
		// Best-effort: revoke both tokens so they cannot be reused after logout.
		a.revokeToken(sess.AccessToken, "access_token")
		a.revokeToken(sess.RefreshToken, "refresh_token")
	}
	a.clearSession(w, r)

	// Redirect to the auth server logout so its session cookie is also cleared.
	// Derive the client base URL from RedirectURI (e.g. "http://localhost:8081/callback" → "http://localhost:8081/").
	clientBase := a.config.RedirectURI[:strings.LastIndex(a.config.RedirectURI, "/")+1]
	logoutURL := a.config.LogoutEndpoint + "?post_logout_redirect_uri=" + url.QueryEscape(clientBase)
	http.Redirect(w, r, logoutURL, http.StatusFound)
}

func (a *App) handleMe(w http.ResponseWriter, r *http.Request) {
	sess := a.getSession(r)
	if sess == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	info, err := a.introspectToken(sess.AccessToken)
	if err != nil {
		http.Error(w, "introspection failed", http.StatusInternalServerError)
		return
	}
	if info["active"] != true {
		newSess, err := a.refreshSession(sess)
		if err != nil {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}
		a.setSession(w, newSess)
		info, _ = a.introspectToken(newSess.AccessToken)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
