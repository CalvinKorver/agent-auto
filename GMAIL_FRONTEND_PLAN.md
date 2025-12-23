# Gmail OAuth Integration - Frontend Implementation Plan

## Overview
Implement end-to-end Gmail OAuth integration in the frontend, allowing users to connect their Gmail account and send email replies directly from the app using their own Gmail account.

## Architecture Overview

### Current Flow
- **Inbound emails**: Received via Mailgun webhooks (no changes needed)
- **Outbound emails**: Currently no email sending capability
- **Messages**: AI generates responses that stay within the app

### New Flow with Gmail
1. User connects Gmail via OAuth (one-time setup)
2. User can send AI-generated responses as actual emails via Gmail
3. Emails properly threaded with original conversation
4. Tokens encrypted and stored securely in backend

---

## Phase 1: Gmail OAuth Connection Flow

### 1.1 Update API Client
**File**: [frontend/lib/api.ts](frontend/lib/api.ts)

Add Gmail API namespace:
```typescript
export const gmailAPI = {
  getAuthUrl: async (): Promise<{ authUrl: string }> => {
    const response = await api.get('/gmail/connect');
    return response.data;
  },

  getStatus: async (): Promise<{ connected: boolean; email?: string }> => {
    const response = await api.get('/gmail/status');
    return response.data;
  },

  disconnect: async (): Promise<void> => {
    await api.post('/gmail/disconnect');
  },
};
```

Update User interface to include Gmail connection status:
```typescript
export interface User {
  id: string;
  email: string;
  createdAt: string;
  inboxEmail?: string;
  preferences?: UserPreferences;
  gmailConnected?: boolean;  // NEW
  gmailEmail?: string;        // NEW
}
```

### 1.2 OAuth Callback Handling

**Important**: The backend handles the OAuth callback directly and redirects to the frontend:
- **Success**: Redirects to `/dashboard`
- **Error**: Redirects to `/gmail-error?error=<error_type>`

**New file**: `frontend/app/gmail-error/page.tsx`

Purpose: Handles error cases from OAuth flow

Key features:
- Extract `error` param from URL query
- Display user-friendly error messages based on error type:
  - `no_code` - "Authorization failed. Please try again."
  - `invalid_state` - "Invalid session. Please try connecting again."
  - `token_exchange_failed` - "Failed to connect Gmail. Please try again."
  - `store_failed` - "Failed to save Gmail connection. Please try again."
- Show "Try Again" button that redirects to dashboard
- Show "Contact Support" option if error persists

Pattern to follow:
- Use `'use client'` directive
- Use `useSearchParams()` from `next/navigation`
- Similar error page styling as other pages

### 1.3 Update NavUser Component
**File**: [frontend/components/dashboard/NavUser.tsx](frontend/components/dashboard/NavUser.tsx)

Add Gmail connection button to dropdown menu:

Changes needed:
1. Import `Mail` icon from `lucide-react`
2. Add state for Gmail connection status and loading
3. Add `useEffect` to fetch Gmail status on mount using `gmailAPI.getStatus()`
4. Add `handleGmailConnect()` function:
   - Calls `gmailAPI.getAuthUrl()` to get OAuth URL
   - Redirects user to Google OAuth: `window.location.href = authUrl`
5. Add `handleGmailDisconnect()` function:
   - Calls `gmailAPI.disconnect()`
   - Updates local state to reflect disconnection
   - Shows success toast
6. Add menu items in `DropdownMenuGroup` (after separator, before logout):
   - If not connected: "Connect Gmail" with Mail icon
   - If connected: "Gmail Connected ✓" with submenu for disconnect option
   - Show loading state while fetching status

### 1.4 Handle Dashboard Re-entry After OAuth
**File**: [frontend/app/dashboard/page.tsx](frontend/app/dashboard/page.tsx)

When user returns to `/dashboard` after successful OAuth:

Add logic to:
1. Check for URL param (optional): `?gmail_connected=true`
2. If present, show success toast: "Gmail connected successfully!"
3. This provides visual confirmation that OAuth succeeded

Alternative approach (simpler):
- Don't add any special handling
- User will see updated status in NavUser dropdown
- Backend already redirects to `/dashboard` after success

---

## Phase 2: Email Sending via Gmail

### 2.1 Install Toast Notification Library

**Action**: Add Sonner to project
```bash
npm install sonner
```

**File**: [frontend/app/layout.tsx](frontend/app/layout.tsx)

Add Toaster component to root layout:
```typescript
import { Toaster } from 'sonner';

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <AuthProvider>
          <ThemeProvider>
            {children}
            <Toaster position="top-right" />
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
```

### 2.2 Add Gmail API to Message Endpoints
**File**: [frontend/lib/api.ts](frontend/lib/api.ts)

Add to `messageAPI` namespace:
```typescript
export const messageAPI = {
  // ... existing methods ...

  replyViaGmail: async (messageId: string, content: string): Promise<void> => {
    await api.post(`/messages/${messageId}/reply-via-gmail`, { content });
  },
};
```

### 2.3 Update ChatPane Component
**File**: [frontend/components/dashboard/ChatPane.tsx](frontend/components/dashboard/ChatPane.tsx)

Major changes to support Gmail sending:

#### A. Add Gmail Send Handler
```typescript
const handleSendViaGmail = async (messageId: string, content: string) => {
  try {
    setSendingViaGmail(true);
    await messageAPI.replyViaGmail(messageId, content);
    toast.success('Email sent via Gmail!');
    setMessageInput(''); // Clear input after successful send
  } catch (error: any) {
    const errorMessage = error.response?.data?.error || '';

    // Handle specific error cases from backend
    if (errorMessage === 'gmail not connected') {
      toast.error('Gmail not connected. Connect in your profile menu.');
    } else if (errorMessage === 'message not found') {
      toast.error('Message not found');
    } else if (errorMessage === 'message was not received via email') {
      toast.error('This message was not received via email and cannot be replied to');
    } else {
      toast.error('Failed to send email via Gmail');
    }
  } finally {
    setSendingViaGmail(false);
  }
};
```

#### B. Add "Send via Gmail" Button

**Important**: The "Send via Gmail" button should only appear when:
1. User has Gmail connected (`user?.gmailConnected`)
2. Currently viewing an inbox message OR a thread with email context
3. The message was received via email (has `externalMessageId`)

Add button next to existing "Send" button in message input area (lines 564-588):

```typescript
<div className="flex gap-2">
  <Button
    onClick={handleSendMessage}
    disabled={sendingMessage || sendingViaGmail || !messageInput.trim()}
  >
    {sendingMessage ? 'Sending...' : 'Send'}
  </Button>

  {/* Show Gmail button only if connected and message has email context */}
  {user?.gmailConnected && currentThread && (
    <Button
      onClick={() => handleSendViaGmail(getReplyableMessageId(), messageInput)}
      disabled={sendingMessage || sendingViaGmail || !messageInput.trim()}
      variant="outline"
    >
      {sendingViaGmail ? 'Sending via Gmail...' : 'Send via Gmail'}
    </Button>
  )}
</div>
```

Note: You'll need to add logic to find the most recent message with `externalMessageId` in the thread to use as the reply target.

#### C. Add State Variables
```typescript
const [sendingViaGmail, setSendingViaGmail] = useState(false);
```

#### D. Show Gmail Badge on Sent Messages
Add visual indicator on messages sent via Gmail (optional enhancement):
- Add `sentViaGmail?: boolean` to Message type
- Show small Gmail icon badge on message bubble
- Backend would need to track this in database

---

## Phase 3: User Experience Enhancements

### 3.1 Gmail Connection Status Indicator

**Option A**: Show in NavUser dropdown (already covered in 1.3)

**Option B**: Create a dedicated Settings page
**New file**: `frontend/app/settings/page.tsx`

Features:
- Gmail connection section
- Show connected email or "Not connected" status
- Connect/Disconnect button
- Display OAuth scopes granted
- Link to Google account permissions

### 3.2 Onboarding Flow Update (Optional)
**File**: [frontend/app/onboarding/page.tsx](frontend/app/onboarding/page.tsx)

Consider adding Gmail connection as optional step:
- After user sets preferences
- Show value proposition: "Connect Gmail to send replies"
- Allow skip with "I'll do this later"

### 3.3 Empty State Messaging
Update empty states to mention Gmail capability:
- When no messages in thread
- When viewing inbox message
- "Connect Gmail to send replies as emails"

---

## Phase 4: Error Handling & Edge Cases

### 4.1 Gmail Not Connected State
When user tries to send via Gmail but hasn't connected:

**In ChatPane.tsx**:
- Check `user?.gmailConnected` before showing button
- Show tooltip: "Connect Gmail in your profile to send emails"
- Or auto-hide button if not connected

### 4.2 Token Expiration
Backend auto-refreshes tokens, but if refresh fails:
- Backend returns error
- Frontend shows: "Gmail connection expired. Please reconnect."
- Provide link to reconnect (opens OAuth flow)

### 4.3 OAuth Callback Errors
**In** `frontend/app/gmail-error/page.tsx`:

Handle these URL params from Google:
- `?error=access_denied` - User denied permission
- `?error=...` - Other OAuth errors

Show friendly messages:
- "You need to grant permission to send emails via Gmail"
- "Something went wrong. Please try again."
- Provide "Try Again" button that redirects to dashboard

### 4.4 Network Errors
Follow existing error handling pattern:
```typescript
catch (error: any) {
  const message = error.response?.data?.error || 'Failed to send email';
  toast.error(message);
  console.error('Gmail send error:', error);
}
```

---

## Phase 5: Testing Plan

### 5.1 Manual Testing Checklist
- [ ] OAuth flow: Click "Connect Gmail" → redirects to Google
- [ ] OAuth consent: Grant permissions → redirects back
- [ ] Callback handling: Successfully processes callback
- [ ] Status update: NavUser shows "Gmail Connected ✓"
- [ ] Auth context: `user.gmailConnected` is true
- [ ] Send button: "Send via Gmail" button appears
- [ ] Send email: Successfully sends via Gmail API
- [ ] Toast notification: Success toast appears
- [ ] Gmail sent folder: Email appears in Gmail
- [ ] Email threading: Reply properly threaded
- [ ] Disconnect: Can disconnect Gmail
- [ ] Status update: Button returns to "Connect Gmail"
- [ ] Error handling: Graceful error messages
- [ ] Token refresh: Works after token expires

### 5.2 Edge Cases to Test
- [ ] Try to send without connecting Gmail first
- [ ] OAuth error (deny permissions)
- [ ] Expired/invalid tokens
- [ ] Network errors during send
- [ ] Send to message without `externalMessageID`
- [ ] Concurrent requests (send multiple emails quickly)

---

## Implementation Order

### Step 1: OAuth Infrastructure (20 min)
1. Update `frontend/lib/api.ts` with `gmailAPI`
2. Update User interface (add gmailConnected, gmailEmail)
3. Create `frontend/app/gmail-error/page.tsx`

### Step 2: Connection UI (30 min)
4. Update `frontend/components/dashboard/NavUser.tsx`
5. Add Gmail status fetching with useEffect
6. Add connect/disconnect handlers

### Step 3: Toast Notifications (10 min)
6. Install Sonner: `npm install sonner`
7. Add Toaster to `frontend/app/layout.tsx`

### Step 4: Email Sending (30 min)
8. Add `replyViaGmail` to messageAPI
9. Update `frontend/components/dashboard/ChatPane.tsx`
10. Add "Send via Gmail" button
11. Add send handler with error handling

### Step 5: Testing & Polish (30 min)
12. Manual testing of full flow
13. Error handling improvements
14. Loading states and UX polish

**Total estimated time: ~2 hours**

---

## Critical Files Summary

### Files to Create
1. `frontend/app/gmail-error/page.tsx` - OAuth error handler

### Files to Modify
1. `frontend/lib/api.ts` - Add gmailAPI namespace, update User type
2. `frontend/components/dashboard/NavUser.tsx` - Add Gmail connect/disconnect
3. `frontend/components/dashboard/ChatPane.tsx` - Add send via Gmail button
4. `frontend/app/layout.tsx` - Add Toaster component

### Dependencies to Install
1. `sonner` - Toast notifications

---

## Backend Integration Points

The backend already provides these endpoints (from GMAIL_SETUP.md):

**OAuth Flow**:
- `GET /api/v1/gmail/connect` - Get OAuth URL (returns `authUrl` and `state`)
- `GET /oauth/callback` - Handle OAuth callback (redirects to `/dashboard` or `/gmail-error`)
- `GET /api/v1/gmail/status` - Check connection status (returns `connected` and optionally `gmailEmail`)
- `POST /api/v1/gmail/disconnect` - Disconnect Gmail

**Email Sending**:
- `POST /api/v1/messages/{messageId}/reply-via-gmail` - Send reply

All endpoints are protected except `/oauth/callback` (public).

---

## Security Considerations

1. **Token Storage**: Tokens stored server-side only (encrypted)
2. **Frontend Security**: Only displays connection status, never tokens
3. **CORS**: Ensure OAuth callback URL matches backend `GOOGLE_REDIRECT_URL`
4. **HTTPS**: In production, ensure all OAuth redirects use HTTPS
5. **State Parameter**: Backend handles CSRF protection via state parameter

---

## Future Enhancements (Post-MVP)

1. **Email Templates**: Pre-fill common responses
2. **Scheduling**: Schedule emails to send later
3. **CC/BCC**: Support for additional recipients
4. **Attachments**: Support file attachments
5. **Email Preview**: Preview email before sending
6. **Signature**: Add custom email signature
7. **Analytics**: Track email open rates (if using tracking pixels)
8. **Multi-account**: Support multiple Gmail accounts

---

## Success Criteria

The implementation is complete when:
- [ ] Users can connect Gmail via OAuth flow
- [ ] Connection status is visible in UI
- [ ] Users can send replies via Gmail
- [ ] Success/error notifications work properly
- [ ] Emails appear in user's Gmail sent folder
- [ ] Replies are properly threaded
- [ ] Users can disconnect Gmail
- [ ] All error cases handled gracefully
- [ ] No console errors or warnings
