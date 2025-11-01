# Learning Materials for HTTP Authentication Middleware

## HTTP Servers and Middleware in Go

Go provides excellent support for building HTTP servers through its standard library. This challenge focuses on implementing an HTTP middleware for authentication.

### HTTP Basics in Go

Go's `net/http` package provides everything needed to build HTTP servers:

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/hello", helloHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Understanding Middleware

Middleware functions sit between the client request and your application logic. They can:

1. Process incoming requests
2. Modify request objects
3. Terminate requests early
4. Modify response objects
5. Chain multiple middlewares together

```go
// Basic middleware structure
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Logic before the handler
        fmt.Println("Before handler")
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Logic after the handler
        fmt.Println("After handler")
    })
}
```

### The http.Handler Interface

At the core of Go's HTTP server is the `http.Handler` interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

And `http.HandlerFunc` adapts regular functions to this interface:

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

### Chaining Middleware

Middleware can be chained together to form a pipeline:

```go
func main() {
    // Create a new mux (router)
    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)
    
    // Wrap mux with middleware chain
    handler := loggingMiddleware(
        authenticationMiddleware(
            rateLimitMiddleware(mux),
        ),
    )
    
    http.ListenAndServe(":8080", handler)
}
```

### Authentication Middleware Patterns

#### Basic Authentication

```go
func basicAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get credentials from request header
        username, password, ok := r.BasicAuth()
        
        // Check credentials
        if !ok || !checkCredentials(username, password) {
            w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Set authenticated user in context
        ctx := context.WithValue(r.Context(), userContextKey, username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### Token Authentication

```go
func tokenAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        // Check format (Bearer token)
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }
        
        token := parts[1]
        
        // Validate token (depends on your token system)
        user, err := validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Set user in context
        ctx := context.WithValue(r.Context(), userContextKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Context in HTTP Requests

Go's `context` package allows passing request-scoped values between middleware and handlers:

```go
// Define a context key type to avoid collisions
type contextKey string

// Define specific keys
const userContextKey contextKey = "user"

// Store a value in context
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Create a new context with the value
        ctx := context.WithValue(r.Context(), userContextKey, "john")
        
        // Call the next handler with updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Retrieve a value from context
func handler(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value(userContextKey).(string)
    if !ok {
        // Handle missing value
        return
    }
    fmt.Fprintf(w, "Hello, %s", user)
}
```

### JWT Authentication

JSON Web Tokens (JWT) are a common method for authentication in web applications:

```go
// Using a popular Go JWT library (github.com/golang-jwt/jwt)
import "github.com/golang-jwt/jwt/v4"

// Create a JWT token
func createToken(username string, secret []byte) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    
    return token.SignedString(secret)
}

// Verify a JWT token
func verifyToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}

// JWT middleware
func jwtMiddleware(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get token from header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }
            
            // Extract token
            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            // Verify token
            claims, err := verifyToken(tokenString, secret)
            if err != nil {
                http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
                return
            }
            
            // Set user in context
            username := claims["username"].(string)
            ctx := context.WithValue(r.Context(), userContextKey, username)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Testing Middleware

Testing middleware and HTTP handlers in Go is straightforward using the `httptest` package:

```go
func TestAuthMiddleware(t *testing.T) {
    // Create a test handler
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value(userContextKey).(string)
        fmt.Fprintf(w, "User: %s", user)
    })
    
    // Wrap with middleware
    handler := authMiddleware(testHandler)
    
    // Create a test request
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    
    // Create a recorder to capture the response
    rr := httptest.NewRecorder()
    
    // Serve the request
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    
    // Check the response body
    expected := "User: john"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}
```

## Best Practices for HTTP Middleware

1. **Use context for request-scoped values**: Store authentication data, request IDs, etc.
2. **Keep middleware focused**: Each middleware should have a single responsibility
3. **Graceful error handling**: Return appropriate HTTP status codes and error messages
4. **Don't forget security headers**: Set headers like `X-XSS-Protection`, `Content-Security-Policy`, etc.
5. **Log requests and errors**: Include a logging middleware for debugging

## Further Reading

- [Go Web Examples: Middleware](https://gowebexamples.com/middleware/)
- [Go Web Examples: JSON Web Tokens](https://gowebexamples.com/jwt/)
- [Context Package Documentation](https://pkg.go.dev/context)
- [Effective Go: Web Servers](https://golang.org/doc/effective_go#web_server) 