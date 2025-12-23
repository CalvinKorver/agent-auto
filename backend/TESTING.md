# Testing Gmail OAuth Integration

## Quick Start

Your backend is already running! Here's how to test it:

### Test 1: Check Server Health

```bash
curl http://localhost:8080/health
```

Expected: `{"status":"healthy","database":"connected"}`

### Test 2: Get an Auth Token

```bash
# Login (user was created by the test script)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'
```

Save the `token` value from the response.

### Test 3: Check Gmail Connection Status

```bash
# Replace YOUR_TOKEN with the token from step 2
curl -X GET http://localhost:8080/api/v1/gmail/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Expected: `{"connected":false}` (until you complete OAuth)

### Test 4: Get OAuth URL

```bash
curl -X GET http://localhost:8080/api/v1/gmail/connect \
  -H "Authorization: Bearer YOUR_TOKEN"
```

You'll get a response like:
```json
{
  "authUrl": "https://accounts.google.com/o/oauth2/auth?...",
  "state": "..."
}
```

### Test 5: Complete OAuth Flow

1. **Copy the `authUrl` from step 4**
2. **Paste it in your browser**
3. **Sign in with your Gmail account**
4. **Authorize the app**
5. **You'll be redirected to** `http://localhost:3000/oauth/callback?code=...&state=...`

The backend will process this callback and store your encrypted OAuth tokens!

### Test 6: Verify Connection

```bash
curl -X GET http://localhost:8080/api/v1/gmail/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Expected: `{"connected":true}` ✅

## Testing Email Reply

To test sending a reply, you need:
1. A message in your inbox (received via Mailgun)
2. The message ID from your database

### Get Inbox Messages

```bash
curl -X GET http://localhost:8080/api/v1/inbox/messages \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Send Reply

```bash
# Replace MESSAGE_ID with an actual message ID from inbox
curl -X POST http://localhost:8080/api/v1/messages/MESSAGE_ID/reply-via-gmail \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Thanks for reaching out! I am very interested in this vehicle."}'
```

**Check your Gmail sent folder** - the reply should be there! 📧

## Troubleshooting

### "Unauthorized" Error
- Make sure you're using the token from the login response
- Token format: `Bearer eyJhbGciOi...` (with "Bearer " prefix)

### "Gmail not connected"
- Complete the OAuth flow in step 5
- Check the status endpoint shows `connected: true`

### "Message not found"
- Make sure the message ID exists in your database
- Check inbox messages endpoint

### "Message was not received via email"
- The message must have `externalMessageID` and `senderEmail` fields
- Only messages received via Mailgun can be replied to

## Database Inspection

To see your stored tokens:

```sql
SELECT user_id, gmail_email, expiry, created_at
FROM gmail_tokens;
```

Note: `access_token` and `refresh_token` are encrypted!

## Complete Test Flow Summary

1. ✅ Start server (`go run cmd/server/main.go`)
2. ✅ Login to get JWT token
3. ✅ Get OAuth URL
4. ✅ Complete OAuth in browser
5. ✅ Check connection status
6. ✅ Send test reply
7. ✅ Verify email in Gmail sent folder

---

**Pro Tip:** Save your JWT token to an environment variable for easier testing:

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Now you can use it easily
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/gmail/status
```
