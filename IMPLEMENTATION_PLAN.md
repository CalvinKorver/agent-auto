# Chat Integration with Model Context - Implementation Plan

## Executive Summary

**Good News:** Your backend is complete and well-architected. The Claude Sonnet 4.5 integration is fully implemented with excellent context management. The frontend message display works perfectly. The **only missing piece** is wiring up the message input UI (textarea + send button).

**Current Context Assessment:** Already excellent and sufficient
- ✅ Last 10 messages for conversation history
- ✅ User preferences (year, make, model)
- ✅ Thread metadata (seller name, type)
- ✅ System prompt with negotiation guidelines

## Implementation Phases

### Phase 1: Wire Up Message Input (CRITICAL - Primary Work)

**Goal:** Enable users to send messages and receive AI-enhanced responses

**File:** [frontend/components/dashboard/ChatPane.tsx](frontend/components/dashboard/ChatPane.tsx)

**Current Issue:**
- Lines 258-276: Message input UI exists but has NO handlers
- No `onChange` on textarea
- No `onClick` on send button
- No state management for input value

**Changes Required:**

#### 1. Add Imports
```typescript
import { useRef } from 'react'; // Add to existing React import
import { Button } from '@/components/ui/button';
```

#### 2. Add State Variables (after line 17)
```typescript
const [messageInput, setMessageInput] = useState('');
const [sendingMessage, setSendingMessage] = useState(false);
const messagesEndRef = useRef<HTMLDivElement>(null);
```

#### 3. Create Message Send Handler
```typescript
const handleSendMessage = async () => {
  if (!selectedThreadId || !messageInput.trim() || sendingMessage) {
    return;
  }

  const content = messageInput.trim();
  setSendingMessage(true);

  try {
    const response = await messageAPI.createMessage(selectedThreadId, {
      content,
      sender: 'user'
    });

    // Add both user and agent messages to local state
    setMessages(prev => [
      ...prev,
      response.userMessage,
      ...(response.agentMessage ? [response.agentMessage] : [])
    ]);

    // Clear input
    setMessageInput('');
  } catch (error) {
    console.error('Failed to send message:', error);
    alert('Failed to send message. Please try again.');
  } finally {
    setSendingMessage(false);
  }
};
```

#### 4. Add Keyboard Handler
```typescript
const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault();
    handleSendMessage();
  }
};
```

#### 5. Add Auto-Scroll Effect
```typescript
useEffect(() => {
  messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
}, [messages]);
```

#### 6. Update Textarea (replace lines 262-266)
```typescript
<textarea
  value={messageInput}
  onChange={(e) => setMessageInput(e.target.value)}
  onKeyDown={handleKeyDown}
  placeholder="Type message... AI will assist (Ctrl+Enter to send)"
  rows={3}
  disabled={sendingMessage}
  className="w-full px-4 py-3 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed"
/>
```

#### 7. Update Send Button (replace lines 268-273)
```typescript
<Button
  onClick={handleSendMessage}
  disabled={sendingMessage || !messageInput.trim()}
  size="lg"
  className="px-6"
>
  {sendingMessage ? 'Sending...' : 'Send'}
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
  </svg>
</Button>
```

#### 8. Add Scroll Target
After messages area, before closing div (around line 255):
```typescript
<div ref={messagesEndRef} />
```

---

### Phase 2: Testing (CRITICAL - Validation)

**Manual Test Checklist:**
- [ ] Login and select a thread
- [ ] Type message in textarea and click send
- [ ] Verify user message appears
- [ ] Verify loading indicator shows
- [ ] Verify AI agent response appears (enhanced version)
- [ ] Verify both messages persist after page refresh
- [ ] Test Ctrl+Enter keyboard shortcut
- [ ] Test empty message prevention (button should be disabled)
- [ ] Test error handling (stop backend, try to send)
- [ ] Switch between threads and verify context isolation

**API Test:**
```bash
# Test message creation endpoint
curl -X POST http://localhost:8080/api/v1/threads/{thread-id}/messages \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"content": "What is your best price?", "sender": "user"}'
```

Expected response includes both userMessage and agentMessage.

---

### Phase 3: Optional Enhancements (POST-MVP)

#### Enhancement 1: Cross-Thread Offer Sharing (Medium Priority)

**Purpose:** Enable AI to reference offers from other threads for competitive negotiation

**How It Works:**
- Users manually "track" specific offers from any seller thread (e.g., "Dealer A offered $25k OTD")
- When chatting with ANY seller, Claude sees ALL tracked offers across all threads
- Claude can say things like: "I have another dealer offering $25k - can you match that?"
- Creates competitive pressure across negotiations

**Implementation Steps:**

**Step 1: Create Offer Tracking API**

Create [backend/internal/api/handlers/offer.go](backend/internal/api/handlers/offer.go):
```go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"carbuyer/internal/db/models"
	"carbuyer/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferHandler struct {
	db *gorm.DB
}

func NewOfferHandler(db *gorm.DB) *OfferHandler {
	return &OfferHandler{db: db}
}

// TrackOffer creates a new tracked offer
func (h *OfferHandler) TrackOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid thread ID"})
		return
	}

	var req struct {
		OfferText string     `json:"offerText"`
		MessageID *uuid.UUID `json:"messageId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	// Verify thread belongs to user
	var thread models.Thread
	if err := h.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "thread not found"})
		return
	}

	offer := models.TrackedOffer{
		ThreadID:  threadID,
		MessageID: req.MessageID,
		OfferText: req.OfferText,
		TrackedAt: time.Now(),
	}

	if err := h.db.Create(&offer).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to track offer"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(offer)
}

// GetTrackedOffers gets all tracked offers for the user
func (h *OfferHandler) GetTrackedOffers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	var offers []models.TrackedOffer
	err := h.db.Joins("JOIN threads ON tracked_offers.thread_id = threads.id").
		Where("threads.user_id = ?", userID).
		Preload("Thread").
		Order("tracked_offers.tracked_at DESC").
		Find(&offers).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to get offers"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"offers": offers})
}

// DeleteTrackedOffer removes a tracked offer
func (h *OfferHandler) DeleteTrackedOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	offerIDStr := chi.URLParam(r, "offerId")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid offer ID"})
		return
	}

	// Verify offer belongs to user's thread
	result := h.db.Joins("JOIN threads ON tracked_offers.thread_id = threads.id").
		Where("tracked_offers.id = ? AND threads.user_id = ?", offerID, userID).
		Delete(&models.TrackedOffer{})

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to delete offer"})
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "offer not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted successfully"})
}
```

**Step 2: Update Routes**

Update [backend/internal/api/routes.go](backend/internal/api/routes.go) to add offer routes:
```go
// Add to RegisterRoutes function:
offerHandler := handlers.NewOfferHandler(db)

r.Route("/api/v1", func(r chi.Router) {
	r.Use(middleware.AuthMiddleware(db))

	// ... existing routes ...

	// Offer routes
	r.Get("/offers", offerHandler.GetTrackedOffers)
	r.Post("/threads/{id}/offers", offerHandler.TrackOffer)
	r.Delete("/offers/{offerId}", offerHandler.DeleteTrackedOffer)
})
```

**Step 3: Update MessageService**

Modify [backend/internal/services/message.go](backend/internal/services/message.go#L59-L106):

After line 87, add:
```go
// Get tracked offers across all user's threads for competitive context
var trackedOffers []models.TrackedOffer
s.db.Joins("JOIN threads ON tracked_offers.thread_id = threads.id").
	Where("threads.user_id = ?", userID).
	Preload("Thread").
	Order("tracked_offers.tracked_at DESC").
	Limit(10). // Last 10 offers
	Find(&trackedOffers)
```

Update line 99 call to include offers:
```go
agentContent, err := s.claudeService.GenerateNegotiationResponse(
	content,
	prefs.Year,
	prefs.Make,
	prefs.Model,
	thread.SellerName,
	recentMessages,
	trackedOffers, // Add this parameter
)
```

**Step 4: Update ClaudeService**

Modify [backend/internal/services/claude.go](backend/internal/services/claude.go#L25-L54):

Update function signature (line 25):
```go
func (s *ClaudeService) GenerateNegotiationResponse(
	userMessage string,
	year int,
	make string,
	model string,
	sellerName string,
	messageHistory []models.Message,
	trackedOffers []models.TrackedOffer, // Add this parameter
) (string, error) {
```

After line 54 (after main system prompt), add:
```go
// Add competitive offers to system prompt
if len(trackedOffers) > 0 {
	systemPrompt += "\n\nCompetitive Offers (from other sellers):\n"
	for _, offer := range trackedOffers {
		sellerInfo := "Unknown Seller"
		if offer.Thread != nil {
			sellerInfo = offer.Thread.SellerName
		}
		systemPrompt += fmt.Sprintf("- %s (from %s)\n", offer.OfferText, sellerInfo)
	}
	systemPrompt += "\nUse these offers as leverage when negotiating. You can mention 'another dealer' without naming them specifically."
}
```

**Step 5: Frontend API Client**

Add to [frontend/lib/api.ts](frontend/lib/api.ts):
```typescript
export interface TrackedOffer {
  id: string;
  threadId: string;
  messageId?: string;
  offerText: string;
  trackedAt: string;
  thread?: {
    sellerName: string;
    sellerType: string;
  };
}

export const offerAPI = {
  trackOffer: async (threadId: string, offerText: string, messageId?: string) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/threads/${threadId}/offers`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ offerText, messageId }),
    });
    if (!response.ok) throw new Error('Failed to track offer');
    return response.json();
  },

  getTrackedOffers: async (): Promise<{ offers: TrackedOffer[] }> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/offers`, {
      method: 'GET',
      headers: getHeaders(),
    });
    if (!response.ok) throw new Error('Failed to get tracked offers');
    return response.json();
  },

  deleteTrackedOffer: async (offerId: string) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/offers/${offerId}`, {
      method: 'DELETE',
      headers: getHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete offer');
    return response.json();
  },
};
```

**Step 6: UI for Tracking Offers**

Add to [ChatPane.tsx](frontend/components/dashboard/ChatPane.tsx):

Add state for tracking:
```typescript
const [trackingOffer, setTrackingOffer] = useState<string | null>(null);
```

Add handler:
```typescript
const handleTrackOffer = async (message: Message) => {
  if (!selectedThreadId) return;

  const offerText = prompt('Enter a brief description of this offer:', message.content.substring(0, 100));
  if (!offerText) return;

  setTrackingOffer(message.id);
  try {
    await offerAPI.trackOffer(selectedThreadId, offerText, message.id);
    alert('Offer tracked! Claude will reference this in other negotiations.');
  } catch (error) {
    console.error('Failed to track offer:', error);
    alert('Failed to track offer');
  } finally {
    setTrackingOffer(null);
  }
};
```

In message rendering (for seller messages), add button:
```typescript
{isSeller && (
  <Button
    variant="ghost"
    size="sm"
    onClick={() => handleTrackOffer(message)}
    disabled={trackingOffer === message.id}
    className="mt-2 text-xs"
  >
    {trackingOffer === message.id ? 'Tracking...' : 'Track Offer'}
  </Button>
)}
```

#### Enhancement 2: Display User Preferences in Sidebar (Low Priority)

**File:** [frontend/components/dashboard/ThreadPane.tsx](frontend/components/dashboard/ThreadPane.tsx)

Add after line 110:
```typescript
<div className="border-b border-slate-700 px-4 py-3">
  <div className="text-xs text-slate-400 uppercase mb-1">Target Vehicle</div>
  <div className="text-sm text-white font-medium">
    {user?.preferences?.year} {user?.preferences?.make} {user?.preferences?.model}
  </div>
</div>
```

#### Enhancement 3: Increase Message History Limit (Optional)

**Current:** 10 messages for Claude context
**Recommendation:** Keep at 10, or increase to maximum 20

**File:** [backend/internal/services/message.go](backend/internal/services/message.go:82)
```go
// Change Limit(10) to Limit(20) if needed
s.db.Where("thread_id = ?", threadID).Order("timestamp DESC").Limit(20).Find(&recentMessages)
```

**Note:** Only increase if you find 10 messages insufficient. More context = higher API costs and slower responses.

---

## Critical Files Reference

### Files to Modify (Phase 1 - Required):
1. [frontend/components/dashboard/ChatPane.tsx](frontend/components/dashboard/ChatPane.tsx) - PRIMARY WORK HERE

### Files to Modify (Phase 3 - Optional):
1. [backend/internal/services/message.go](backend/internal/services/message.go)
2. [backend/internal/services/claude.go](backend/internal/services/claude.go)
3. [frontend/components/dashboard/ThreadPane.tsx](frontend/components/dashboard/ThreadPane.tsx)

### Files That Work Perfectly (NO CHANGES NEEDED):
1. [backend/internal/api/handlers/message.go](backend/internal/api/handlers/message.go) - API endpoints complete
2. [frontend/lib/api.ts](frontend/lib/api.ts) - API client ready
3. [backend/internal/db/models/message.go](backend/internal/db/models/message.go) - Database models solid

---

## Context Management Analysis

### Current Context (Excellent and Sufficient)

**What Claude receives per message:**
1. **System Prompt:**
   - User's car requirements (year, make, model)
   - Seller name
   - Negotiation guidelines

2. **Message History:**
   - Last 10 messages from thread
   - Formatted with sender labels (user/agent/seller)

3. **Current Message:**
   - User's draft message to enhance

**Assessment:** This is well-designed and provides sufficient context for effective negotiation assistance. The 10-message limit balances context quality with API efficiency.

### Why Current Context Is Good Enough

- ✅ **10 messages = ~5 conversation turns** - Sufficient for coherent conversation
- ✅ **User preferences included** - AI knows what car they want
- ✅ **Seller context present** - Knows who they're negotiating with
- ✅ **Clear AI role** - Enhance messages for negotiation
- ✅ **Cost-efficient** - Doesn't waste tokens on excessive history
- ✅ **Fast responses** - Less context = faster Claude API calls

---

## Implementation Strategy

### Recommended Path

**Step 1 (Today - 2 hours):**
- Implement Phase 1: Wire up message input in ChatPane
- Test basic send/receive flow
- Validate end-to-end with real messages

**Step 2 (Testing - 30 minutes):**
- Run through manual test checklist
- Verify AI responses are contextually appropriate
- Test error scenarios

**Step 3 (Optional - Later):**
- Add tracked offers context if competitive negotiation is needed
- Display user preferences in sidebar for visibility
- Consider message history adjustments based on usage

### What NOT to Do

- ❌ Don't redesign backend Claude integration (it's excellent)
- ❌ Don't add WebSockets/real-time updates (not needed for MVP)
- ❌ Don't implement message editing/deletion (post-MVP)
- ❌ Don't add file attachments (out of scope)
- ❌ Don't over-engineer context management (current is optimal)

---

## Expected Results

**After Phase 1:**
- Users can type and send messages
- AI agent receives user input and enhances it for negotiation
- Both user and agent messages display in chat
- Conversation flows naturally
- Context is maintained across the conversation
- Loading states provide clear feedback

**Example Flow:**
1. User types: "price?"
2. User clicks Send
3. User message appears: "price?"
4. Loading indicator shows: "Crafting negotiation message..."
5. Agent message appears: "We are serious buyers for this 2024 Mazda CX-90. Could you please provide your best out-the-door pricing, including all fees and taxes?"
6. Both messages saved to database and persist

---

## Success Criteria

✅ Users can send messages via textarea and send button
✅ AI agent responds with enhanced negotiation message
✅ Loading states show during AI generation
✅ Errors are handled gracefully
✅ Messages persist across page refreshes
✅ Context is maintained within threads
✅ Keyboard shortcuts work (Ctrl+Enter)
✅ Auto-scroll to new messages
✅ Thread isolation works correctly
