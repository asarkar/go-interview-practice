[View the Scoreboard](SCOREBOARD.md)

# Challenge 15: OAuth2 Authentication System

In this challenge, you will implement an OAuth2 authentication system using Go. You'll create a server that supports the OAuth2 authorization code flow, allowing third-party applications to authenticate users without directly handling their credentials.

## Requirements

1. Implement an OAuth2 server that supports:
   - Client registration and management
   - Authorization endpoint for user consent
   - Token endpoint for exchanging codes for tokens
   - Token validation and introspection
   
2. Your implementation should support the following OAuth2 flows:
   - Authorization code grant
   - Refresh token flow
   
3. Implement security best practices:
   - PKCE (Proof Key for Code Exchange) support
   - Token expiration and revocation
   - Scope-based permissions
   - Secure storage of client secrets and tokens
   
4. Create a simple demo client application that:
   - Redirects users to the authorization endpoint
   - Exchanges authorization codes for tokens
   - Uses tokens to access protected resources
   - Refreshes tokens when they expire
   
5. The included test file has scenarios covering normal flows, error cases, and security edge cases 