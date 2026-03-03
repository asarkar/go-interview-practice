package server

import (
	"html/template"
	"net/url"
)

var loginTemplate = template.Must( //nolint:gochecknoglobals
	template.New("login").Funcs(template.FuncMap{
		"urlquery": url.QueryEscape,
	}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - OAuth2 Server</title>
    <style>
        body { font-family: system-ui, sans-serif; max-width: 400px; margin: 3rem auto; padding: 0 1rem; }
        h1 { font-size: 1.25rem; }
        input { display: block; width: 100%; padding: 0.5rem; margin: 0.5rem 0; box-sizing: border-box; }
        button { padding: 0.5rem 1rem; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Sign in</h1>
    <form method="POST" action="/login?redirect={{urlquery .Redirect}}">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <input type="text" name="username" value="user" required autofocus>
        <input type="password" name="password" value="password" required>
        <button type="submit">Sign in</button>
    </form>
</body>
</html>`),
)
