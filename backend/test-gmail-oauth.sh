#!/bin/bash

# Gmail OAuth Testing Script
# This script tests the complete Gmail OAuth flow

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

echo "======================================"
echo "Gmail OAuth Integration Test"
echo "======================================"
echo ""

# Step 1: Register a test user
echo "Step 1: Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123"
  }')

echo "$REGISTER_RESPONSE" | jq '.' 2>/dev/null || echo "$REGISTER_RESPONSE"
echo ""

# Step 2: Login to get JWT token
echo "Step 2: Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "Failed to get token. Response:"
  echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"
  exit 1
fi

echo "✓ Logged in successfully"
echo "Token: ${TOKEN:0:50}..."
echo ""

# Step 3: Check Gmail connection status (should be false)
echo "Step 3: Checking Gmail status (should be not connected)..."
STATUS_RESPONSE=$(curl -s -X GET "$API_URL/gmail/status" \
  -H "Authorization: Bearer $TOKEN")

echo "$STATUS_RESPONSE" | jq '.' 2>/dev/null || echo "$STATUS_RESPONSE"
echo ""

# Step 4: Get OAuth authorization URL
echo "Step 4: Getting OAuth authorization URL..."
OAUTH_RESPONSE=$(curl -s -X GET "$API_URL/gmail/connect" \
  -H "Authorization: Bearer $TOKEN")

AUTH_URL=$(echo "$OAUTH_RESPONSE" | jq -r '.authUrl' 2>/dev/null)

if [ -z "$AUTH_URL" ] || [ "$AUTH_URL" = "null" ]; then
  echo "Failed to get auth URL. Response:"
  echo "$OAUTH_RESPONSE" | jq '.' 2>/dev/null || echo "$OAUTH_RESPONSE"
  exit 1
fi

echo "✓ OAuth URL generated successfully"
echo ""
echo "======================================"
echo "MANUAL STEP REQUIRED"
echo "======================================"
echo ""
echo "Please open this URL in your browser to authorize Gmail:"
echo ""
echo "$AUTH_URL"
echo ""
echo "After authorizing, you'll be redirected to:"
echo "http://localhost:3000/oauth/callback"
echo ""
echo "The callback will process and store your tokens."
echo ""
echo "======================================"
echo ""

# Wait for user to complete OAuth
echo "Waiting 30 seconds for you to complete OAuth..."
echo "(Press Ctrl+C to skip if you're testing something else)"
sleep 30

# Step 5: Check if Gmail is now connected
echo ""
echo "Step 5: Checking if Gmail was connected..."
STATUS_RESPONSE=$(curl -s -X GET "$API_URL/gmail/status" \
  -H "Authorization: Bearer $TOKEN")

CONNECTED=$(echo "$STATUS_RESPONSE" | jq -r '.connected' 2>/dev/null)

if [ "$CONNECTED" = "true" ]; then
  echo "✓ Gmail connected successfully!"
  echo "$STATUS_RESPONSE" | jq '.'
else
  echo "Gmail not yet connected. Status:"
  echo "$STATUS_RESPONSE" | jq '.' 2>/dev/null || echo "$STATUS_RESPONSE"
fi

echo ""
echo "======================================"
echo "Test Complete"
echo "======================================"
echo ""
echo "Your JWT token for manual testing:"
echo "$TOKEN"
echo ""
echo "Save this token to test the reply endpoint later!"
