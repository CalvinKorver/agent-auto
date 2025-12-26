#!/bin/bash

# Quick script to get OAuth URL for testing

echo "Getting OAuth URL for Gmail connection..."
echo ""

# Login
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Failed to login. Response:"
  echo "$LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Logged in successfully"
echo ""

# Get OAuth URL
OAUTH_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/gmail/connect" \
  -H "Authorization: Bearer $TOKEN")

AUTH_URL=$(echo "$OAUTH_RESPONSE" | jq -r '.authUrl' 2>/dev/null)

if [ -z "$AUTH_URL" ] || [ "$AUTH_URL" = "null" ]; then
  echo "❌ Failed to get OAuth URL. Response:"
  echo "$OAUTH_RESPONSE"
  exit 1
fi

echo "✅ OAuth URL generated!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Open this URL in your browser to connect Gmail:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "$AUTH_URL"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Your JWT token (save this for testing):"
echo "$TOKEN"
echo ""
