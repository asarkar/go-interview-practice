# OAuth2 Client

Entry point for the browser-based OAuth2 demo client. Application logic lives in `challenge15/client`.

## Running

1. Start the OAuth2 auth server (from project root):

   ```bash
   mise exec -- go run ./challenge15/cmd/server
   ```

   The auth server listens on `:8080` and registers the client with redirect URI `http://localhost:8081/callback`.

2. Start the client (in a separate terminal):

   ```bash
   mise exec -- go run ./challenge15/cmd/client
   ```

   The client listens on `:8081`.

## Usage

1. Open http://localhost:8081
2. Click **Login with OAuth**
3. You will be redirected to the auth server's login page
4. Sign in with `testuser` / `password`
5. After consent, you will be redirected back to the client
6. Use **View API Response** to call `/api/me` (introspects the token)
7. Use **Logout** to clear the session

## Routes

| Route       | Method | Description                          |
|-------------|--------|--------------------------------------|
| `/`         | GET    | Home (login or user info)           |
| `/login`    | GET    | Redirects to OAuth authorize         |
| `/callback` | GET    | OAuth callback (code exchange)       |
| `/logout`   | GET    | Clears session                       |
| `/api/me`   | GET    | Returns introspected token (JSON)    |
