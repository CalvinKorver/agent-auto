## Otto Repo Notes

- Purpose: AI-assisted car-buying app. Backend Go API, frontend Next.js 14 (app router + Tailwind + shadcn). Production: backend on Railway, frontend on https://www.OttoAI.com.

### Key Docs
- `README.md` for stack/overview, `DOCS/QUICKSTART.md` for prod deploy, `DOCS/IMPLEMENTATION_PLAN.md` for chat feature status, `DOCS/SPEC.md` for API/architecture, `backend/MIGRATIONS.md` for migration approach, `backend/TESTING.md` for Gmail OAuth tests.

### Local Development
- One-shot: `./dev.sh` (installs/runs Air for backend hot reload, then `npm run dev` in frontend).
- Backend only: `cd backend && go run cmd/server/main.go` (loads `.env`, auto-migrates via GORM).
- Frontend only: `cd frontend && npm run dev`.
- Health: `curl http://localhost:8080/health` → `{"status":"healthy","database":"connected"}`.

### Environment Variables (backend, required unless default noted)
- `PORT` (default 8080), `ENVIRONMENT`, `DATABASE_URL`, `JWT_SECRET`, `JWT_EXPIRATION_HOURS` (default 24), `ANTHROPIC_API_KEY`.
- CORS: `ALLOWED_ORIGINS` comma list (first used for Gmail redirect).
- Email/Gmail: `MAILGUN_API_KEY`, `MAILGUN_DOMAIN`, `MAILGUN_WEBHOOK_SIGNING_KEY`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` (default `http://localhost:3000/oauth/callback`), `TOKEN_ENCRYPTION_KEY`.
- Rate limits: `RATE_LIMIT_AUTH`, `RATE_LIMIT_API`.

### Backend Map
- Entry: `backend/cmd/server/main.go` sets router, CORS, health; wires services/handlers; routes under `/api/v1` plus `/oauth/callback`.
- Config: `backend/internal/config/config.go` reads env + defaults.
- DB/models: `backend/internal/db` (AutoMigrate on start).
- Services: `backend/internal/services/` (auth, preferences, threads, messages + Claude, email via Mailgun + Gmail, Gmail OAuth tokens).
- HTTP handlers: `backend/internal/api/handlers/`.
- Notable routes (all JWT except health/webhooks): auth register/login/me/logout; preferences get/post; threads CRUD + messages; offers (create under thread, list all); inbox assign/archive; Gmail connect/status/disconnect; message reply via Gmail; webhooks for inbound/test email.

### Frontend Map
- Next.js app router in `frontend/app/`; main screens under `app/dashboard`, auth under `login`/`register`, onboarding under `onboarding`.
- Shared UI in `frontend/components/` (shadcn-based), API client `frontend/lib/api.ts`, auth context `frontend/contexts/AuthContext.tsx`.

### Testing
- Gmail OAuth manual flow in `backend/TESTING.md`.
- Quick login for tests: `curl -X POST http://localhost:8080/api/v1/auth/login -d '{"email":"test@example.com","password":"testpass123"}'` (from testing doc).

### Migrations
- Currently GORM AutoMigrate on startup; `backend/MIGRATIONS.md` documents plan to move to `golang-migrate` with commands/examples.

### Deployment Cheats
- Backend: Dockerfile at repo root (Railway auto-detect). Render alt: build `go build -o bin/server cmd/server/main.go`, start `./bin/server`.
- Frontend: deploy `frontend/` on Vercel with `npm run build`.
- Vercel env: `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_ENVIRONMENT`. Backend CORS must include deployed frontend URL.

### Logging
- **Backend (Go)**: Uses standard library `log` package. Import with `import "log"`. Use `log.Printf()` for formatted logging, `log.Println()` for simple messages. Example: `log.Printf("Login error for email %s: %v", email, err)`. All authentication errors (login/register failures) are logged with email and error details.
- **Frontend (TypeScript/Next.js)**: Uses browser `console.error()` for error logging. Example: `console.error('Login error:', errorMessage, err)`. All authentication errors (login/register failures) are logged to browser console with error message and full error object.

### Helpful References
- UI/feature status: `DOCS/IMPLEMENTATION_PLAN.md` (Phase 1 chat complete; offer tracking planned).
- Design mocks: `ui-mockup/`.
- Gmail-specific helper scripts: `backend/get-oauth-url.sh`, `backend/test-gmail-oauth.sh`.

