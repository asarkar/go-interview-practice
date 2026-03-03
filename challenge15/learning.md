# Learning Materials for OAuth2 Authentication System

## OAuth2 Overview

OAuth 2.0 is an authorization framework that enables a third-party application to obtain limited access to an HTTP service, either on behalf of a resource owner or by allowing the third-party application to obtain access on its own behalf.

### Key Concepts

- **Resource Owner**: The user who owns the data (e.g., a user who has photos on a photo-sharing site)
- **Client**: The third-party application that wants to access the user's data
- **Authorization Server**: The server that authenticates the resource owner and issues access tokens
- **Resource Server**: The server hosting the protected resources (can be the same as the authorization server)
- **Access Token**: A credential used by the client to access protected resources
- **Refresh Token**: A credential used to obtain new access tokens when they expire

### OAuth2 Flows

OAuth 2.0 defines several grant types or flows for different use cases:

1. **Authorization Code**: For server-side web applications
2. **Implicit**: For browser-based or mobile apps (less secure, now discouraged)
3. **Resource Owner Password Credentials**: For trusted applications
4. **Client Credentials**: For application access (no user involved)
5. **Refresh Token**: For getting new access tokens without re-authorization
6. **Device Code**: For devices with limited input capabilities

## Authorization Code Flow

The authorization code flow is the most secure and commonly used flow. It works as follows:

1. The client redirects the user to the authorization server with its client ID, requested scope, and redirect URI
2. The user authenticates and grants permissions
3. The authorization server redirects back to the client with an authorization code
4. The client exchanges the authorization code for access and refresh tokens
5. The client uses the access token to access protected resources

```
+----------+
| Resource |
|   Owner  |
+----------+
     ^
     |
    (B)
     |
+----|-----+          Client Identifier      +---------------+
|         -+----(A)-- & Redirection URI ---->|               |
|  Client  |                                  | Authorization |
|          |<---(C)-- Authorization Code ----|    Server     |
|          |                                  |               |
|          |----(D)-- Authorization Code ---->|               |
|          |          & Redirection URI       |               |
|          |                                  |               |
|          |<---(E)----- Access Token -------|               |
+-----------+      (w/ Optional Refresh Token)  +---------------+
```

## PKCE (Proof Key for Code Exchange)

PKCE (pronounced "pixy") is an extension to the authorization code flow that provides additional security for public clients (like mobile or SPA applications). It works as follows:

1. The client creates a code verifier (a random string)
2. The client generates a code challenge from the code verifier using a transformation method (usually SHA-256)
3. The client includes the code challenge and challenge method in the authorization request
4. When exchanging the authorization code, the client includes the original code verifier
5. The authorization server transforms the code verifier and compares it to the stored code challenge

```go
// Create a code verifier
codeVerifier, _ := GenerateRandomString(64)

// Create a code challenge using S256 method
h := sha256.New()
h.Write([]byte(codeVerifier))
codeChallenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
```

## Access Tokens

Access tokens are credentials used to access protected resources. They can be of different formats:

1. **Opaque Tokens**: Random strings that the resource server must validate with the authorization server
2. **JWT Tokens**: JSON Web Tokens that contain claims about the authentication and authorization

```go
// Example JWT token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

## Refresh Tokens

Refresh tokens are credentials used to obtain new access tokens when they expire. They typically have a longer lifetime than access tokens, but are more sensitive as they grant access for a longer period.

```go
// Using a refresh token to get a new access token
form := url.Values{}
form.Add("grant_type", "refresh_token")
form.Add("refresh_token", refreshToken)
form.Add("client_id", clientID)
form.Add("client_secret", clientSecret)

resp, err := http.PostForm(tokenEndpoint, form)
```

## Scopes

Scopes define the specific access permissions requested by the client.

```
Authorization Request with Scopes:
GET /authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz
    &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb
    &scope=read write email
```

## Token Storage

Access and refresh tokens should be securely stored:

- **Server-Side Applications**: Store in a secure database
- **Browser-Based Applications**: Store access tokens in memory, refresh tokens in HTTP-only secure cookies
- **Mobile Applications**: Use the platform's secure storage (Keychain on iOS, KeyStore on Android)

## Token Validation

When validating tokens, you should check:

1. Token signature (for JWT tokens)
2. Token expiration
3. Token issuer
4. Token audience
5. Token scope

```go
// Validating a token
token, err := server.ValidateToken(tokenString)
if err != nil {
    // Handle invalid token
    return
}

// Check if the token has the required scope
if !containsScope(token.Scopes, requiredScope) {
    // Handle insufficient scope
    return
}
```

## Error Handling

OAuth2 defines standard error responses:

- `invalid_request`: The request is missing a parameter or is otherwise malformed
- `invalid_client`: Client authentication failed
- `invalid_grant`: The authorization grant is invalid or expired
- `unauthorized_client`: The client is not authorized to use this grant type
- `unsupported_grant_type`: The authorization server does not support this grant type
- `invalid_scope`: The requested scope is invalid or unknown

```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

## Security Considerations

1. **Use HTTPS**: All OAuth 2.0 communications should be over TLS
2. **Validate Redirect URIs**: Only redirect to pre-registered URIs
3. **Use PKCE**: For all clients, especially public clients
4. **Short-Lived Access Tokens**: Limit the lifetime of access tokens
5. **Validate State Parameter**: To prevent CSRF attacks
6. **Secure Token Storage**: Store tokens securely based on client type
7. **Token Revocation**: Implement token revocation for when tokens are compromised

## JWT Structure and Validation

JWT tokens consist of three parts:

1. **Header**: Identifies which algorithm is used to generate the signature
2. **Payload**: Contains the claims
3. **Signature**: Ensures that the token hasn't been altered

```go
// JWT Header
{
  "alg": "HS256",
  "typ": "JWT"
}

// JWT Payload
{
  "sub": "1234567890", // subject (user id)
  "name": "John Doe",
  "iat": 1516239022,   // issued at
  "exp": 1516242622,   // expiration time
  "aud": "my-api",     // audience
  "iss": "auth-server" // issuer
}
```

## OAuth2 in Go

Here are some useful libraries for implementing OAuth2 in Go:

- `golang.org/x/oauth2`: Client-side OAuth2 implementation
- `github.com/go-oauth2/oauth2`: Server-side OAuth2 implementation
- `github.com/golang-jwt/jwt`: JWT implementation

```go
// Using golang.org/x/oauth2 for client-side OAuth2
import "golang.org/x/oauth2"

conf := &oauth2.Config{
    ClientID:     "client-id",
    ClientSecret: "client-secret",
    Scopes:       []string{"read", "write"},
    RedirectURL:  "http://localhost:8080/callback",
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://provider.com/o/oauth2/auth",
        TokenURL: "https://provider.com/o/oauth2/token",
    },
}

// Generate authorization URL
url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

// Exchange authorization code for token
token, err := conf.Exchange(ctx, code)
```

## OpenID Connect

OpenID Connect (OIDC) is an identity layer built on top of OAuth 2.0. It allows clients to verify the identity of the end-user and to obtain basic profile information.

OIDC adds:

1. **ID Token**: A JWT containing user identity information
2. **UserInfo Endpoint**: For getting more user information
3. **Standard Claims**: For user information like name, email, etc.

```json
// Example ID Token payload
{
  "iss": "https://server.example.com",
  "sub": "24400320",
  "aud": "s6BhdRkqt3",
  "exp": 1311281970,
  "iat": 1311280970,
  "name": "Jane Doe",
  "email": "janedoe@example.com"
}
```

## Best Practices

1. **Follow the Specs**: Implement the OAuth2 specification correctly
2. **Use PKCE**: For all clients, even confidential ones
3. **Short-Lived Tokens**: Keep access tokens short-lived (< 1 hour)
4. **Rotate Refresh Tokens**: Issue new refresh tokens when refreshing access tokens
5. **Implement Token Revocation**: Allow users to revoke access
6. **Validate All Parameters**: Including redirect URIs and scopes
7. **Use Standard Error Codes**: Follow the OAuth2 error response format
8. **Log Authentication Events**: Log all token issuance and usage for auditing
9. **Test Security**: Include security testing in your test suite
10. **Stay Up-to-Date**: Keep your dependencies updated for security patches 