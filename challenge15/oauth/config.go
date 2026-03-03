package oauth

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
