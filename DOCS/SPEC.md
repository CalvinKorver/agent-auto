# Lolo AI Application - Technical Specification

## Overview
An AI-powered car buying assistant that helps users negotiate and communicate with multiple car sellers simultaneously through dedicated agent-based chat threads.

## Architecture

### Frontend
- **Framework**: Next.js
- **UI Layout**:
  - Left sidebar: List of active agent chats/threads
  - Main area: Selected chat conversation view

### Backend
- **Framework**: Go (Golang)
- **LLM Provider**: Anthropic Claude (initial implementation)

## Core Concepts

### Master Agent Architecture
- **Single master agent** handles all seller communications across all threads
- Agent has awareness of ALL active negotiations simultaneously
- Can leverage information from one thread to inform negotiations in others
- Enables cross-thread strategies (e.g., "Seller A offered $X, can you match?")
- Maintains consistent negotiation approach across all conversations
- utilizes and references the log of different potential offers from sellers

### Thread-Based Chat System
- Each thread represents a communication channel with a specific seller entity:
  - Individual private seller
  - Car dealership
  - Other seller types
- Users can select/switch between different threads from the sidebar
- Each thread maintains its own message history
- Master agent routes responses to the appropriate thread based on context
- Maintains a single page front end that keeps the chats on the side and in the main bar, shows the current chat
- Looks like this:
+----------------------+------------------------------------------------+
| SIDEBAR              | MAIN CONTENT AREA                              |
+----------------------+------------------------------------------------+
| **Lolo AI** |                                                |
|                      |                                                |
| TARGET VEHICLE:      |                                                |
| Mazda CX-90          |                                                |
|                      |                                                |
|                      |         +----------------------------+         |
| THREADS              |         |                            |         |
| (No active negotiations)|      |    No Active Negotiations  |         |
|                      |         |                            |         |
|                      |         |  Add a seller in the       |         |
|                      |         |  sidebar to begin chatting.|         |
|                      |         |                            |         |
|                      |         +----------------------------+         |
|                      |                                                |
|                      |                                                |
| [ + New Seller Thread ]|                                                |
|                      |                                                |
+----------------------+------------------------------------------------+

### User Preferences (Car Selection)
**MVP Scope:**
- Make, model, and year selection (e.g., "2024 Mazda CX-90")
- Defined once at the start of the car buying process
- All agent chats inherit and respect these preferences

**Post-MVP:**
- Price threshold (low/high range)
- Additional filters and preferences (mileage, condition, features, etc.)

### AI Negotiation Prompt
- **Hardcoded system prompt** (for MVP)
- Instructs the LLM that it is a "master car negotiator"
- Enforces adherence to user's defined preferences (make, model, price constraints)
- Ensures consistent negotiation strategy across all agent threads

## UI Mockups

Visual mockups of key screens are available in the `ui-mockup/` directory:

### Setup Screen
**File**: [ui-mockup/setupscreen.png](ui-mockup/setupscreen.png)

Shows the initial onboarding screen where users define their target vehicle preferences. Features:
- Clean, centered card layout
- Car icon and app branding
- Clear messaging about locked preferences
- Input fields for Make and Model (Note: Year field needs to be added to mockup)
- Prominent "START NEGOTIATING" CTA button

### Main Dashboard (Empty State)
**File**: [ui-mockup/main1.png](ui-mockup/main1.png)

Shows the main dashboard after preferences are set but before any seller threads are created. Features:
- Dark sidebar with app branding
- Target vehicle display (shows selected make/model)
- "THREADS" section showing "(No active negotiations)"
- "+ New Seller Thread" button at bottom of sidebar
- Empty state in main content area with helpful prompt
- Clean, professional design with good use of whitespace

### Recommended Additional Mockups

The following screens should be mocked up to complete the UI design:

1. **Login/Registration Screen** (`login.png`)
   - Email and password input fields
   - "Login" and "Register" options
   - Consistent branding with setup screen

2. **Active Chat View** (`active-chat.png`)
   - Sidebar with multiple seller threads listed
   - Active thread highlighted
   - Main area showing chat conversation between user/agent/seller
   - Message input box at bottom
   - Display of tracked offers in sidebar

3. **New Thread Modal/Dialog** (`new-thread-modal.png`)
   - Form to add new seller
   - Fields: Seller Name, Seller Type (dropdown: Private/Dealership/Other)
   - Create and Cancel buttons

4. **Offer Tracking UI** (`offer-tracking.png`)
   - How the "Track this offer?" prompt appears in chat
   - Global offers section in sidebar showing all tracked offers
   - Visual indicator when AI detects potential offer

5. **Chat with Tracked Offers** (`chat-with-offers.png`)
   - Full dashboard showing active negotiation
   - Multiple threads in sidebar
   - Active chat with message history
   - Tracked offers displayed prominently
   - Shows how all pieces work together

## User Flow

### Initial Onboarding
When a user first opens the dashboard:
1. **Check for existing preferences**: If no user preferences are set, redirect to preference input form (see `setupscreen.png`)
2. **Preference Input Form**: Collect required information:
   - Year (dropdown or number input, e.g., 2020-2025)
   - Make (text input or dropdown, e.g., "Mazda")
   - Model (text input, e.g., "CX-90")
3. **Save and Continue**: Once preferences are submitted, user is directed to the main dashboard
4. **Locked Preferences**: For the current buying session, these preferences are locked and apply to all agent threads

### Main Dashboard Flow
After preferences are set, users can:
- View the empty state if no threads exist yet (see `main1.png`)
- Create their first agent thread for a seller
- Begin negotiating through the AI assistant

## Key Features

### 1. Multi-Thread Chat Management
- Create new agent threads for different sellers
- View all active threads in sidebar
- Switch between threads seamlessly
- Each thread is isolated but follows the same user preferences

### 2. User Preference Definition
- Initial setup: User specifies desired make/model/year
- Preferences are locked for the buying session
- All agents operate within these constraints

### 3. AI-Powered Negotiation
- Each message sent by user can be augmented/refined by AI
- Master agent maintains negotiation strategy across all threads
- Agent has full context of all negotiations when generating responses
- Respects user's constraints and preferences
- Uses Anthropic Claude API for message generation

### 4. Offer Tracking System
**Simple Text-Based Tracking (MVP):**
- Agent identifies potential offers in conversation
- User prompted to "Track" offers they want to remember
- Tracked offers stored as simple text entries per thread
- Master agent receives all tracked offers across all threads in its context
+----------------------+------------------------------------------------+
| SIDEBAR              | MAIN CONTENT AREA - ACTIVE THREAD              |
+----------------------+------------------------------------------------+
| **Lolo AI** |                                                |
|                      | Header: **Downtown Toyota Dealership** [⚙️]    |
| TARGET VEHICLE:      |------------------------------------------------|
| Mazda CX-90          | (Chat History Scroll Area)                     |
|                      |                                                |
|                      | [USER]: Do you have any room on the sticker    |
|                      | price?                                         |
| ==================== |                                                |
| **GLOBAL OFFERS** | [AGENT]: (AI refines and sends negotiation)    |
| * $39.5k (Toyota)    | "We are serious buyers ready to move today if  |
| * $41k (CarMax)      | the price is right. What is your best out-the- |
| ==================== | door price?"                                   |
|                      |                                                |
| THREADS              | [SELLER - Downtown Toyota]: We can do $39,500  |
| > Downtown Toyota    | plus tax and title fees if you buy today.      |
|   Private Seller Bob |                                                |
|   CarMax             |    +---------------------------------------+   |
|                      |    | 🤖 AI detected potential offer:         |   |
|                      |    | "$39,500 plus tax/title"                |   |
| [ + New Seller Thread ]|  | [ TRACK THIS OFFER ]  [ Ignore ]        |   |
|                      |    +---------------------------------------+   |
|                      |                                                |
|                      | [USER]: That sounds interesting. Let me think. |
|                      |                                                |
|                      |                                                |
|                      |------------------------------------------------|
|                      | [ Type message... AI will assist ]  [ SEND > ] |
+----------------------+------------------------------------------------+

**User Flow:**
1. Seller makes an offer in conversation (e.g., "$39k plus $1200 in fees")
2. Agent detects potential offer and prompts user: "Track this offer?"
3. User clicks "Track" button or ignores
4. If tracked, offer is added to thread's offer log
5. All tracked offers from all threads are injected into agent context for future messages

**Benefits:**
- Cross-thread price awareness enables competitive negotiation
- User controls what gets tracked (quality over quantity)
- Agent can reference better offers from other sellers
- Simple implementation without complex parsing

## Authentication & User Management

### Authentication Strategy (MVP)
**Simple Email/Password Authentication**
- User registration with email and password
- Login with email/password
- Session-based authentication using JWT tokens
- No social login (OAuth) in MVP

### Auth Flow
1. **Registration**:
   - User provides email and password
   - Backend validates email format and password strength
   - Password hashed using bcrypt
   - User record created in database
   - JWT token issued upon successful registration

2. **Login**:
   - User provides email and password
   - Backend validates credentials
   - JWT token issued upon successful login
   - Token stored in httpOnly cookie or localStorage

3. **Protected Routes**:
   - All API endpoints (except auth endpoints) require valid JWT
   - Frontend checks for valid token before rendering dashboard
   - Redirect to login if token is missing or expired

4. **Logout**:
   - Clear JWT token from client
   - Invalidate session on backend (optional for MVP)

### User Model
```
User:
- id: uuid
- email: string (unique, indexed)
- passwordHash: string
- createdAt: timestamp
- updatedAt: timestamp
```

### Auth Endpoints
- `POST /auth/register` - Create new user account
- `POST /auth/login` - Authenticate user and return JWT
- `POST /auth/logout` - Invalidate session
- `GET /auth/me` - Get current authenticated user info

### Security Considerations (MVP)
- Password hashing with bcrypt (cost factor 10+)
- JWT tokens with expiration (24 hours)
- HTTPS only in production
- Basic rate limiting on auth endpoints
- Email validation and password strength requirements

### Post-MVP Auth Features
- Password reset flow
- Email verification
- OAuth providers (Google, Apple)
- Multi-factor authentication (MFA)
- Session management across devices

## API Design

### Base URL
```
Development: http://localhost:8080/api/v1
Production: https://api.carbuyer.com/v1
```

### Authentication
All API requests (except auth endpoints) must include JWT token:
```
Authorization: Bearer <jwt_token>
```

### API Endpoints Specification

#### Authentication Endpoints

**POST /auth/register**
```json
Request:
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response (201 Created):
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}

Errors:
400 - Invalid email format or weak password
409 - Email already exists
```

**POST /auth/login**
```json
Request:
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response (200 OK):
{
  "user": {
    "id": "uuid",
    "email": "user@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}

Errors:
401 - Invalid credentials
```

**GET /auth/me**
```json
Response (200 OK):
{
  "id": "uuid",
  "email": "user@example.com",
  "preferences": {
    "year": 2024,
    "make": "Mazda",
    "model": "CX-90"
  }
}

Errors:
401 - Unauthorized (invalid/missing token)
```

#### User Preferences Endpoints

**GET /preferences**
```json
Response (200 OK):
{
  "year": 2024,
  "make": "Mazda",
  "model": "CX-90",
  "createdAt": "2024-01-15T10:30:00Z"
}

Response (404 Not Found):
{
  "error": "No preferences set"
}
```

**POST /preferences**
```json
Request:
{
  "year": 2024,
  "make": "Mazda",
  "model": "CX-90"
}

Response (201 Created):
{
  "year": 2024,
  "make": "Mazda",
  "model": "CX-90",
  "createdAt": "2024-01-15T10:30:00Z"
}

Errors:
400 - Invalid year/make/model
401 - Unauthorized
```

#### Thread Endpoints

**POST /threads**
```json
Request:
{
  "sellerName": "Downtown Toyota",
  "sellerType": "dealership"
}

Response (201 Created):
{
  "id": "thread-uuid",
  "sellerName": "Downtown Toyota",
  "sellerType": "dealership",
  "createdAt": "2024-01-15T10:30:00Z",
  "messageCount": 0
}

Errors:
400 - Missing required fields
401 - Unauthorized
404 - No preferences set (must set preferences first)
```

**GET /threads**
```json
Response (200 OK):
{
  "threads": [
    {
      "id": "thread-uuid-1",
      "sellerName": "Downtown Toyota",
      "sellerType": "dealership",
      "createdAt": "2024-01-15T10:30:00Z",
      "lastMessageAt": "2024-01-15T11:45:00Z",
      "messageCount": 12
    },
    {
      "id": "thread-uuid-2",
      "sellerName": "Private Seller Bob",
      "sellerType": "private",
      "createdAt": "2024-01-14T09:00:00Z",
      "lastMessageAt": "2024-01-14T15:30:00Z",
      "messageCount": 8
    }
  ]
}

Errors:
401 - Unauthorized
```

**GET /threads/:id**
```json
Response (200 OK):
{
  "id": "thread-uuid",
  "sellerName": "Downtown Toyota",
  "sellerType": "dealership",
  "createdAt": "2024-01-15T10:30:00Z",
  "messageCount": 12
}

Errors:
401 - Unauthorized
404 - Thread not found
403 - Thread belongs to different user
```

#### Message Endpoints

**GET /threads/:id/messages**
```json
Query Params:
- limit: number (default 50, max 100)
- offset: number (default 0)

Response (200 OK):
{
  "messages": [
    {
      "id": "msg-uuid-1",
      "threadId": "thread-uuid",
      "sender": "user",
      "content": "Do you have any room on the sticker price?",
      "timestamp": "2024-01-15T10:35:00Z"
    },
    {
      "id": "msg-uuid-2",
      "threadId": "thread-uuid",
      "sender": "agent",
      "content": "We are serious buyers ready to move today if the price is right. What is your best out-the-door price?",
      "timestamp": "2024-01-15T10:35:30Z"
    },
    {
      "id": "msg-uuid-3",
      "threadId": "thread-uuid",
      "sender": "seller",
      "content": "We can do $39,500 plus tax and title fees if you buy today.",
      "timestamp": "2024-01-15T10:40:00Z"
    }
  ],
  "total": 12,
  "hasMore": false
}

Errors:
401 - Unauthorized
404 - Thread not found
403 - Thread belongs to different user
```

**POST /threads/:id/messages**
```json
Request:
{
  "content": "What's the mileage on this vehicle?",
  "sender": "user"
}

Response (201 Created):
{
  "userMessage": {
    "id": "msg-uuid-4",
    "threadId": "thread-uuid",
    "sender": "user",
    "content": "What's the mileage on this vehicle?",
    "timestamp": "2024-01-15T11:00:00Z"
  },
  "agentMessage": {
    "id": "msg-uuid-5",
    "threadId": "thread-uuid",
    "sender": "agent",
    "content": "Could you please provide the current mileage and service history for this 2024 Mazda CX-90?",
    "timestamp": "2024-01-15T11:00:15Z"
  }
}

Note: Backend calls Claude API to generate agent's enhanced message
Errors:
400 - Empty message content
401 - Unauthorized
404 - Thread not found
403 - Thread belongs to different user
```

#### Offer Tracking Endpoints

**POST /threads/:id/offers/track**
```json
Request:
{
  "messageId": "msg-uuid-3",
  "offerText": "$39,500 plus tax and title fees"
}

Response (201 Created):
{
  "id": "offer-uuid",
  "threadId": "thread-uuid",
  "messageId": "msg-uuid-3",
  "offerText": "$39,500 plus tax and title fees",
  "trackedAt": "2024-01-15T11:05:00Z"
}

Errors:
400 - Missing required fields
401 - Unauthorized
404 - Thread or message not found
```

**GET /offers**
```json
Response (200 OK):
{
  "offers": [
    {
      "id": "offer-uuid-1",
      "threadId": "thread-uuid-1",
      "sellerName": "Downtown Toyota",
      "offerText": "$39,500 plus tax and title fees",
      "trackedAt": "2024-01-15T11:05:00Z"
    },
    {
      "id": "offer-uuid-2",
      "threadId": "thread-uuid-2",
      "sellerName": "CarMax",
      "offerText": "$41,000 out the door",
      "trackedAt": "2024-01-14T16:00:00Z"
    }
  ]
}

Note: Returns all tracked offers across all threads for current user
Errors:
401 - Unauthorized
```

### Error Response Format
All error responses follow this structure:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {} // Optional additional context
  }
}
```

### Rate Limiting
- Auth endpoints: 5 requests per minute per IP
- All other endpoints: 100 requests per minute per user
- Response headers include:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`

## Technical Components

### Frontend (Next.js)
- **Preference Input Form** (initial onboarding)
  - Year selector
  - Make input field
  - Model input field
  - Submit button

    A[Simple Centered Card/Modal] --> B(Title: "Let's define your target vehicle")
    B --> C{Input: Vehicle Make e.g., Mazda}
    C --> D{Input: Vehicle Model e.g., CX-90}
    D --> E[Button: Start Negotiating]
    E --> F(Redirects to Main Dashboard)

    style A fill:#f9f9f9,stroke:#333,stroke-width:2px
    style E fill:#4CAF50,stroke:#333,stroke-width:1px,color:white

- **Chat Sidebar Component**
  - List of active threads
  - Thread creation button
  - Display seller names
- **Message Thread View**
  - Message history display
  - Sender identification (user/agent/seller)
  - Timestamps
- **Message Composition Interface**
  - Text input area
  - Send button
  - (Future) AI suggestion toggle

### Backend (Go)
- **API endpoints for:**
  - `GET /preferences` - Retrieve current user preferences
  - `POST /preferences` - Create/update user preferences
  - `POST /threads` - Create new agent thread
  - `GET /threads` - Fetch list of all threads
  - `GET /threads/:id/messages` - Get messages for a specific thread
  - `POST /threads/:id/messages` - Send message within a thread
  - `POST /threads/:id/offers/track` - Track an offer in a thread
- **Integration with Anthropic Claude API**
  - System prompt injection with user preferences
  - Message generation with full context across threads
- **Data Storage**
  - User preference persistence
  - Message history storage
  - Tracked offers storage

### Data Models (Detailed)

```
User:
- id: uuid (primary key)
- email: string (unique, indexed)
- passwordHash: string
- createdAt: timestamp
- updatedAt: timestamp
- Relationships:
  - has one UserPreferences
  - has many Threads

UserPreferences:
- id: uuid (primary key)
- userId: uuid (foreign key -> User.id, unique)
- make: string
- model: string
- year: integer
- createdAt: timestamp
- updatedAt: timestamp
- (future) priceMin: decimal
- (future) priceMax: decimal

Thread:
- id: uuid (primary key)
- userId: uuid (foreign key -> User.id, indexed)
- sellerName: string
- sellerType: enum (private, dealership, other)
- createdAt: timestamp
- updatedAt: timestamp
- Relationships:
  - belongs to User
  - has many Messages
  - has many TrackedOffers

Message:
- id: uuid (primary key)
- threadId: uuid (foreign key -> Thread.id, indexed)
- sender: enum (user, agent, seller)
- content: text
- timestamp: timestamp
- Relationships:
  - belongs to Thread

TrackedOffer:
- id: uuid (primary key)
- threadId: uuid (foreign key -> Thread.id, indexed)
- messageId: uuid (foreign key -> Message.id, nullable)
- offerText: text
- trackedAt: timestamp
- Relationships:
  - belongs to Thread
  - optionally references Message
```

## Database

### Database Choice
**PostgreSQL** (recommended for MVP)
- Mature, reliable, open-source
- Excellent support for UUIDs and JSON
- Good Go library support (pgx, GORM)
- Easy to deploy (Railway, Render, Supabase)

### Alternative: SQLite
- For local development or small-scale deployment
- Simpler setup, single file database
- Good for prototyping

### Schema Indexing Strategy
```sql
-- Users table
CREATE INDEX idx_users_email ON users(email);

-- Threads table
CREATE INDEX idx_threads_user_id ON threads(user_id);
CREATE INDEX idx_threads_created_at ON threads(created_at DESC);

-- Messages table
CREATE INDEX idx_messages_thread_id ON messages(thread_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp DESC);

-- TrackedOffers table
CREATE INDEX idx_tracked_offers_thread_id ON tracked_offers(thread_id);
```

### Migrations
- Use migration tool (e.g., golang-migrate, Goose)
- Version-controlled migration files
- Up/down migrations for reversibility

## Environment Variables

### Backend (Go)
```bash
# Server
PORT=8080
ENVIRONMENT=development # or production

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/carbuyer

# Authentication
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24

# Anthropic API
ANTHROPIC_API_KEY=sk-ant-...

# CORS
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Rate Limiting
RATE_LIMIT_AUTH=5  # requests per minute
RATE_LIMIT_API=100 # requests per minute
```

### Frontend (Next.js)
```bash
# API
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Environment
NEXT_PUBLIC_ENVIRONMENT=development
```

## MVP Scope

### In Scope:
- **Authentication & User Management**
  - Email/password registration and login
  - JWT-based authentication
  - Protected routes and API endpoints
- **User Preferences**
  - Initial onboarding form for year/make/model
  - Preference persistence per user
  - Redirect logic if preferences not set
- **UI Components**
  - Login/registration pages
  - Preference input form
  - Chat sidebar with thread list
  - Main chat view with message history
  - Message composition interface
- **Thread Management**
  - Create new seller threads
  - List all threads for authenticated user
  - Switch between threads
  - Display thread metadata (seller name, type, message count)
- **Message System**
  - Send user messages
  - AI agent message generation via Claude API
  - Display message history with sender identification
  - Manual seller message entry (for MVP testing)
- **Offer Tracking**
  - Track offers from conversations
  - View all tracked offers across threads
  - Display tracked offers in sidebar
- **Backend API**
  - RESTful API with all specified endpoints
  - PostgreSQL database integration
  - Anthropic Claude API integration
  - JWT middleware for authentication
  - Basic error handling and validation
  - Rate limiting

### Out of Scope (Post-MVP):
- Price threshold filtering
- Advanced negotiation strategies
- Multiple car preference profiles
- Real-time notifications (WebSockets)
- Seller response simulation/integration
- Payment or transaction features
- OAuth social login
- Password reset flow
- Email verification
- Multi-factor authentication
- Mobile responsive design (desktop-first MVP)
- Thread archiving/deletion
- Message editing/deletion
- File attachments in messages
- Search functionality
- Analytics and reporting

## System Prompt (Initial Draft)

```
You are an expert car negotiation assistant helping a buyer communicate with car sellers.
Your goal is to secure the best possible deal while maintaining professional and respectful communication.

User's Requirements:
- Year: {year}
- Make: {make}
- Model: {model}

Guidelines:
- Always negotiate within the user's specified requirements
- Be firm but polite in negotiations
- Ask relevant questions about vehicle condition, history, and pricing
- Highlight any concerns or red flags
- Work towards the best price and terms for the buyer
- Never deviate from the specified year, make, and model
```

## Technology Stack Summary

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context API (or Zustand for simplicity)
- **HTTP Client**: Fetch API or Axios
- **Form Handling**: React Hook Form
- **Validation**: Zod

### Backend
- **Language**: Go 1.21+
- **Web Framework**: Chi or Gin (lightweight, fast routers)
- **Database ORM**: GORM or sqlx
- **Database**: PostgreSQL 15+
- **Authentication**: JWT (golang-jwt/jwt library)
- **Password Hashing**: bcrypt
- **API Client**: Anthropic Go SDK
- **Migrations**: golang-migrate or Goose

### Infrastructure & Deployment

#### Development
- **Frontend**: `npm run dev` on localhost:3000
- **Backend**: `go run main.go` on localhost:8080
- **Database**: Local PostgreSQL or Docker container

#### Production (Recommended)
- **Frontend Hosting**: Vercel (optimized for Next.js)
- **Backend Hosting**: Railway, Render, or Fly.io
- **Database**: Railway PostgreSQL, Supabase, or Neon
- **Environment Secrets**: Platform-specific secret management

#### Alternative: Docker Compose
```yaml
version: '3.8'
services:
  frontend:
    build: ./frontend
    ports: ["3000:3000"]
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

  backend:
    build: ./backend
    ports: ["8080:8080"]
    environment:
      - DATABASE_URL=postgresql://postgres:password@db:5432/carbuyer
      - JWT_SECRET=${JWT_SECRET}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=carbuyer
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Project Structure

### Frontend (Next.js)
```
frontend/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── (dashboard)/
│   │   ├── layout.tsx           # Dashboard layout with sidebar
│   │   ├── page.tsx              # Main dashboard (redirects to onboarding if no prefs)
│   │   └── onboarding/
│   │       └── page.tsx          # Preference input form
│   ├── layout.tsx
│   └── page.tsx                  # Landing page (redirects to login or dashboard)
├── components/
│   ├── auth/
│   │   ├── LoginForm.tsx
│   │   └── RegisterForm.tsx
│   ├── dashboard/
│   │   ├── Sidebar.tsx
│   │   ├── ThreadList.tsx
│   │   ├── ChatView.tsx
│   │   ├── MessageList.tsx
│   │   ├── MessageInput.tsx
│   │   └── OfferTracker.tsx
│   └── onboarding/
│       └── PreferenceForm.tsx
├── lib/
│   ├── api.ts                    # API client functions
│   ├── auth.ts                   # Auth utilities
│   └── types.ts                  # TypeScript types
├── contexts/
│   └── AuthContext.tsx
└── package.json
```

### Backend (Go)
```
backend/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go          # Auth handlers
│   │   │   ├── preferences.go
│   │   │   ├── threads.go
│   │   │   ├── messages.go
│   │   │   └── offers.go
│   │   ├── middleware/
│   │   │   ├── auth.go          # JWT middleware
│   │   │   ├── cors.go
│   │   │   └── ratelimit.go
│   │   └── router.go             # Route definitions
│   ├── db/
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── preferences.go
│   │   │   ├── thread.go
│   │   │   ├── message.go
│   │   │   └── offer.go
│   │   └── db.go                 # Database connection
│   ├── services/
│   │   ├── auth.go               # Auth business logic
│   │   ├── claude.go             # Claude API integration
│   │   ├── negotiation.go        # Negotiation logic
│   │   └── offers.go
│   └── config/
│       └── config.go             # Configuration management
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_preferences.up.sql
│   └── ...
├── go.mod
└── go.sum
```

## Development Workflow

### Initial Setup
1. **Clone repository** and create project structure
2. **Backend setup**:
   ```bash
   cd backend
   go mod init carbuyer
   go get github.com/go-chi/chi/v5
   go get gorm.io/gorm
   go get gorm.io/driver/postgres
   go get github.com/golang-jwt/jwt/v5
   go get github.com/anthropics/anthropic-sdk-go
   ```
3. **Frontend setup**:
   ```bash
   cd frontend
   npx create-next-app@latest . --typescript --tailwind --app
   npm install axios react-hook-form zod
   ```
4. **Database setup**:
   ```bash
   createdb carbuyer
   # Run migrations
   ```
5. **Environment variables**: Copy `.env.example` to `.env` and fill in values

### Development Flow
1. Start database (if not already running)
2. Start backend: `cd backend && go run cmd/server/main.go`
3. Start frontend: `cd frontend && npm run dev`
4. Access app at http://localhost:3000

### Testing Strategy (Post-MVP)
- **Frontend**: React Testing Library, Jest
- **Backend**: Go standard testing package
- **E2E**: Playwright or Cypress
- **API Testing**: Postman collections

## Next Steps

### Phase 1: Foundation (Week 1)
1. Set up project repositories and structure
2. Initialize Next.js frontend with basic routing
3. Initialize Go backend with Chi/Gin router
4. Set up PostgreSQL database and initial migrations
5. Configure environment variables and secrets

### Phase 2: Authentication (Week 2)
1. Implement user registration and login endpoints
2. Build JWT authentication middleware
3. Create login/register UI components
4. Implement protected route logic on frontend
5. Test auth flow end-to-end

### Phase 3: Core Features (Week 3-4)
1. Build user preferences system (backend + frontend)
2. Implement onboarding flow
3. Create thread management (CRUD operations)
4. Build message system with Claude integration
5. Design and implement chat UI components

### Phase 4: Advanced Features (Week 5)
1. Implement offer tracking system
2. Build global offers display in sidebar
3. Refine Claude system prompt for negotiations
4. Add rate limiting and error handling
5. Polish UI/UX

### Phase 5: Testing & Deployment (Week 6)
1. End-to-end testing of all features
2. Security audit (SQL injection, XSS, CSRF)
3. Performance optimization
4. Deploy to production environment
5. Monitor and iterate based on feedback