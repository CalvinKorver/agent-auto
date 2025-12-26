# Gmail OAuth Integration Setup Guide

## Overview
This guide will help you set up Gmail OAuth integration so users can send email replies from their own Gmail accounts.

## Architecture
- **Inbound emails**: Continue using Mailgun webhooks (no change)
- **Outbound emails**: Users connect their Gmail via OAuth and send replies through Gmail API
- **Security**: OAuth tokens encrypted with AES-256-GCM and stored in database

## Step 1: Configure Google Cloud Console

### 1.1 Create Google Cloud Project
1. Go to https://console.cloud.google.com
2. Click "Create Project"
3. Name: "Lolo AI" (or your preference)
4. Click "Create"

### 1.2 Enable Gmail API
1. In your project, go to "APIs & Services" → "Library"
2. Search for "Gmail API"
3. Click "Enable"

### 1.3 Configure OAuth Consent Screen
1. Go to "APIs & Services" → "OAuth consent screen"
2. Select "External" user type
3. Fill in required information:
   - **App name**: Car Buyer
   - **User support email**: Your email
   - **Developer contact**: Your email
4. Click "Save and Continue"
5. **Scopes**: Add the following scope:
   - `https://www.googleapis.com/auth/gmail.compose` (Manage drafts and send emails)
6. Click "Save and Continue"
7. **Test users**: Add your Gmail address (and any test users)
   - You can add up to 100 test users while in testing mode
8. Click "Save and Continue"
9. Status will be "Testing" (this is fine for now)

### 1.4 Create OAuth 2.0 Credentials
1. Go to "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "OAuth 2.0 Client ID"
3. **Application type**: Web application
4. **Name**: Car Buyer Web Client
5. **Authorized redirect URIs**: Add these:
   - `http://localhost:3000/oauth/callback` (for local development)
   - Add production URL later when deploying
6. Click "Create"
7. **IMPORTANT**: Copy the Client ID and Client Secret (you'll need these next)

## Step 2: Update Backend Environment Variables

1. Copy the example environment file if you haven't already:
   ```bash
   cp .env.example .env
   ```

2. Add the following to your `.env` file:

```bash
# Google OAuth (from Google Cloud Console)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000/oauth/callback

# Token Encryption Key (generated for you below)
TOKEN_ENCRYPTION_KEY=08670aa67eacd4e10dc6a6b954e09ef8e9a89fd0160b5815f125cc066ae9abde
```

Replace `your-client-id` and `your-client-secret` with values from Google Cloud Console.

**Note**: The encryption key above has been generated for you. Keep it secret!

## Step 3: Test the Backend

1. Start your backend server:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. The server should start without errors and you should see:
   ```
   Database connection established
   Running database migrations...
   Migrations completed successfully
   Server starting on http://localhost:8080
   ```

## Step 4: API Endpoints

The following endpoints are now available:

### OAuth Flow
- `GET /api/v1/gmail/connect` - Get OAuth authorization URL (protected)
- `GET /oauth/callback` - OAuth callback handler (public)
- `GET /api/v1/gmail/status` - Check if user has Gmail connected (protected)
- `POST /api/v1/gmail/disconnect` - Disconnect Gmail (protected)

### Email Sending
- `POST /api/v1/messages/{messageId}/reply-via-gmail` - Send reply via Gmail (protected)

## Step 5: Testing the OAuth Flow (Manual)

### 5.1 Register/Login a User
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'

# Login (save the token)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'
```

### 5.2 Get OAuth URL
```bash
curl -X GET http://localhost:8080/api/v1/gmail/connect \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response will include an `authUrl`. Copy and paste it into your browser.

### 5.3 Authorize in Browser
1. Visit the `authUrl` from step 5.2
2. Sign in with your Gmail account
3. Grant permissions
4. You'll be redirected to `http://localhost:3000/oauth/callback` (may show error if frontend not running, but that's OK)

### 5.4 Check Connection Status
```bash
curl -X GET http://localhost:8080/api/v1/gmail/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Should return: `{"connected":true}`

### 5.5 Send a Test Reply

First, you need a message in your inbox (received via Mailgun). Get a message ID from your database, then:

```bash
curl -X POST http://localhost:8080/api/v1/messages/{messageId}/reply-via-gmail \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Thanks for your message! I am interested in this vehicle."}'
```

Check your Gmail sent folder - the reply should be there, properly threaded!

## Step 6: Database Schema

A new table `gmail_tokens` has been created with the following structure:
- `user_id` (UUID, primary key)
- `access_token` (TEXT, encrypted)
- `refresh_token` (TEXT, encrypted)
- `token_type` (VARCHAR)
- `expiry` (TIMESTAMP)
- `gmail_email` (VARCHAR)
- `created_at`, `updated_at`

## Verification Checklist

- [ ] Google Cloud project created
- [ ] Gmail API enabled
- [ ] OAuth consent screen configured
- [ ] OAuth 2.0 credentials created
- [ ] Test user added to OAuth consent screen
- [ ] Environment variables updated in `.env`
- [ ] Backend starts without errors
- [ ] Database migration successful (`gmail_tokens` table created)
- [ ] OAuth flow tested successfully
- [ ] Gmail status returns connected
- [ ] Test email sent successfully

## Security Notes

1. **Encryption**: All OAuth tokens are encrypted with AES-256-GCM before storage
2. **Token Refresh**: Access tokens auto-refresh when expired
3. **Scopes**: Only `gmail.send` scope is requested (minimal permissions)
4. **HTTPS**: In production, ensure all OAuth redirects use HTTPS
5. **Secrets**: Never commit `.env` file to git

## Troubleshooting

### "No refresh token received"
- This happens if user has already authorized your app
- Solution: Revoke access in Google Account settings and try again
- Or: Use `oauth2.ApprovalForce` (already implemented)

### "Gmail not connected" error
- User hasn't completed OAuth flow yet
- Check `gmail_tokens` table in database

### Token decryption errors
- Encryption key changed or corrupted
- Check `TOKEN_ENCRYPTION_KEY` is 64 hex characters

### CORS errors in frontend
- Add your frontend URL to `ALLOWED_ORIGINS` in `.env`
- Update Google Cloud Console redirect URIs

## Next Steps (Frontend Integration)

When ready to add frontend UI:
1. Add "Connect Gmail" button
2. Add "Send Email" button on AI-drafted responses
3. Handle OAuth callback page
4. Show Gmail connection status

See the plan file for detailed frontend integration steps.

## Production Deployment

Before going to production:
1. Change `GOOGLE_REDIRECT_URL` to your production domain
2. Add production redirect URI in Google Cloud Console
3. Submit app for verification (if needed for public access)
4. Use strong `TOKEN_ENCRYPTION_KEY` (never reuse dev key)
5. Enable HTTPS for all endpoints
