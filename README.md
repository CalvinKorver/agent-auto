# Lolo AI Application

An AI-powered car buying assistant that helps users negotiate and communicate with multiple car sellers simultaneously through dedicated agent-based chat threads.

## Quick Links

- **[Production Deployment Guide](QUICKSTART.md)** - Deploy to production in 30 minutes
- **[Deployment Documentation](DEPLOYMENT.md)** - Detailed deployment information
- **[Technical Specification](SPEC.md)** - Complete API and architecture docs
- **[Database Migrations](backend/MIGRATIONS.md)** - Migration strategy and best practices

## Project Structure

```
car-buyer/
├── frontend/              # Next.js frontend application
├── backend/               # Go backend API server
│   ├── cmd/server/       # Main application entry point
│   ├── internal/         # Internal packages
│   │   ├── api/         # HTTP handlers and middleware
│   │   ├── config/      # Configuration management
│   │   ├── db/          # Database connection
│   │   ├── models/      # Database models
│   │   └── services/    # Business logic
│   └── MIGRATIONS.md    # Migration documentation
├── .github/workflows/    # CI/CD pipelines
├── Dockerfile           # Backend container configuration
├── .dockerignore        # Docker build exclusions
├── railway.json         # Railway deployment config
├── render.yaml          # Render deployment config
├── vercel.json          # Vercel deployment config
├── DEPLOYMENT.md        # Detailed deployment guide
├── QUICKSTART.md        # Quick deployment guide (30 min)
└── SPEC.md              # Technical specification
```

## Technology Stack

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context API
- **HTTP Client**: Axios
- **Form Handling**: React Hook Form
- **Validation**: Zod

### Backend
- **Language**: Go 1.21+
- **Web Framework**: Chi (lightweight, fast router)
- **Database ORM**: GORM
- **Database**: PostgreSQL 15+
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **AI Provider**: Anthropic Claude SDK

## Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.21+
- **PostgreSQL** 15+
- **Anthropic API Key**

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd car-buyer
```

### 2. Database Setup

Create a PostgreSQL database:

```bash
createdb carbuyer
```

Or using psql:

```sql
CREATE DATABASE carbuyer;
```

### 3. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Copy environment variables
cp .env.example .env

# Edit .env and fill in your values:
# - DATABASE_URL: Your PostgreSQL connection string
# - JWT_SECRET: A secure random string
# - ANTHROPIC_API_KEY: Your Anthropic API key

# Run the server (will auto-migrate database)
go run cmd/server/main.go
```

The backend API will start on `http://localhost:8080`

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Environment variables are already configured in .env.local
# Update if needed

# Run the development server
npm run dev
```

The frontend will start on `http://localhost:3000`

## Development Workflow

### Running Both Services

Terminal 1 (Backend):
```bash
cd backend
go run cmd/server/main.go
```

Terminal 2 (Frontend):
```bash
cd frontend
npm run dev
```

### Database Migrations

The backend uses GORM's AutoMigrate feature, which automatically creates and updates tables based on the model definitions. Migrations run automatically when the server starts.

## Project Status - Phase 1 Complete

✅ **Phase 1: Foundation** (Complete)
- Project structure created
- Next.js frontend initialized with TypeScript and Tailwind CSS
- Go backend initialized with Chi router
- PostgreSQL database models and schema defined
- Environment configuration set up
- Project documentation created

### Phase 1 Deliverables

**Backend Structure:**
- Database models (User, UserPreferences, Thread, Message, TrackedOffer)
- Database connection and auto-migration setup
- Configuration management with environment variables
- Project dependencies installed (Chi, GORM, PostgreSQL driver, JWT, Anthropic SDK, bcrypt, CORS)

**Frontend Structure:**
- Next.js 14 with App Router
- TypeScript configuration
- Tailwind CSS styling
- Dependencies installed (Axios, React Hook Form, Zod, @hookform/resolvers)

**Configuration:**
- Environment variable templates for both frontend and backend
- Git ignore files
- Go module configuration

## Next Steps - Phase 2: Authentication

The next phase will implement:
1. User registration and login endpoints
2. JWT authentication middleware
3. Login/register UI components
4. Protected route logic on frontend
5. End-to-end auth flow testing

## API Documentation

Base URL: `http://localhost:8080/api/v1`

Detailed API documentation is available in [SPEC.md](SPEC.md#api-design)

## Environment Variables

### Backend (.env)

```bash
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgresql://user:password@localhost:5432/carbuyer
JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRATION_HOURS=24
ANTHROPIC_API_KEY=sk-ant-your-api-key-here
ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT_AUTH=5
RATE_LIMIT_API=100
```

### Frontend (.env.local)

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_ENVIRONMENT=development
```

## Key Features (Planned)

### Authentication & User Management
- Email/password registration and login
- JWT-based authentication
- Protected routes and API endpoints

### User Preferences
- Initial onboarding form for year/make/model
- Preference persistence per user
- Redirect logic if preferences not set

### Thread Management
- Create new seller threads
- List all threads for authenticated user
- Switch between threads
- Display thread metadata

### Message System
- Send user messages
- AI agent message generation via Claude API
- Display message history with sender identification
- Manual seller message entry (for MVP testing)

### Offer Tracking
- Track offers from conversations
- View all tracked offers across threads
- Display tracked offers in sidebar

## Contributing

See [SPEC.md](SPEC.md) for detailed technical specifications and development guidelines.

## Production Deployment

Ready to deploy? See [QUICKSTART.md](QUICKSTART.md) for step-by-step deployment instructions.

### Recommended Stack
- **Frontend**: Vercel (Next.js native)
- **Backend**: Railway or Render (containerized Go)
- **Database**: Neon (serverless PostgreSQL with branching)
- **Cost**: ~$5-10/month to start

### Quick Deploy Checklist
- [ ] Create Neon PostgreSQL databases (prod, preview, local)
- [ ] Deploy backend to Railway/Render
- [ ] Deploy frontend to Vercel
- [ ] Configure environment variables
- [ ] Set up preview environments
- [ ] Test end-to-end flow

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed information.

## Documentation

- [QUICKSTART.md](QUICKSTART.md) - Production deployment in 30 minutes
- [DEPLOYMENT.md](DEPLOYMENT.md) - Detailed deployment documentation
- [SPEC.md](SPEC.md) - Complete technical specification and API docs
- [backend/MIGRATIONS.md](backend/MIGRATIONS.md) - Database migration strategy
- [.env.template](.env.template) - Environment variable template

## License

[Add your license here]
