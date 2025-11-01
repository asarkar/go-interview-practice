[View the Scoreboard](SCOREBOARD.md)

# Challenge 5: HTTP Authentication Middleware

In this challenge, you must implement an HTTP middleware in Go that checks each incoming request for a valid authentication token. If the token is invalid, the middleware should return an HTTP 401 Unauthorized response. If valid, it should pass the request to the next handler.

## Requirements

1. The middleware looks for an HTTP header "X-Auth-Token".  
2. If the header is present and equals a predefined "secret", the request is allowed and should pass to the final handler.  
3. Otherwise, return 401 Unauthorized.  
4. The router has two endpoints:  
   - GET /hello -> returns "Hello!"  
   - GET /secure -> returns "You are authorized!"  
5. The included test file has 10 scenarios checking correct behavior for valid tokens, invalid tokens, missing headers, etc.
