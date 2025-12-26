# Project Instructions (Lolo AI)

## Deployment & Infrastructure
- Backend: Railway (`https://agent-auto-production.up.railway.app/`, port 8080). Dockerfile in repo root; Render alternative build `go build -o bin/server cmd/server/main.go`, start `./bin/server`.
- Frontend: https://www.agentauto.com on Vercel (root `frontend/`, build `npm run build`).

## Development Quickstart
- Run both with hot reload: `./dev.sh` (installs/runs Air in `backend/`, then `npm run dev` in `frontend/`).
- Backend only: `cd backend && go run cmd/server/main.go` (loads `.env`, AutoMigrate).
- Frontend only: `cd frontend && npm run dev`.
- Health: `curl http://localhost:8080/health`.
- Gmail test creds (from `backend/TESTING.md`): `test@example.com` / `testpass123`.

## Repo Map
- Backend entry/router: `backend/cmd/server/main.go` (chi, CORS, routes, OAuth callback).
- Config/env: `backend/internal/config/config.go`.
- Services: `backend/internal/services/` (auth, preferences, threads, messages+Claude, email/Mailgun, Gmail OAuth, offers).
- Handlers: `backend/internal/api/handlers/`; middleware in `backend/internal/api/middleware/`.
- DB: `backend/internal/db` (AutoMigrate) and models in `backend/internal/db/models/`.
- Frontend: Next.js app router under `frontend/app/`; UI components `frontend/components/`; API client `frontend/lib/api.ts`; auth context `frontend/contexts/AuthContext.tsx`.
- Docs: `instructions.md` (cheatsheet), `README.md`, `DOCS/QUICKSTART.md`, `DOCS/IMPLEMENTATION_PLAN.md`, `backend/MIGRATIONS.md`, `backend/TESTING.md`.

## API Surface (JWT unless noted)
- `/health`.
- `/api/v1/auth` register | login | me | logout.
- `/api/v1/preferences` get/post.
- `/api/v1/threads` list/create/get/delete + `/threads/{id}/messages` get/post + `/threads/{id}/offers` post.
- `/api/v1/offers` get.
- `/api/v1/inbox/messages` get | assign | delete.
- `/api/v1/gmail` connect | status | disconnect.
- `/api/v1/messages/{messageId}/reply-via-gmail`.
- Webhooks (public): `/api/v1/webhooks/email/inbound|test`; OAuth callback `/oauth/callback`.

## Environment Variables (backend)
- Required: `DATABASE_URL`, `JWT_SECRET`, `ANTHROPIC_API_KEY`.
- Also: `PORT` (default 8080), `ENVIRONMENT`, `ALLOWED_ORIGINS` comma list (first used for Gmail redirect), `RATE_LIMIT_AUTH`, `RATE_LIMIT_API`.
- Email/Gmail: `MAILGUN_API_KEY`, `MAILGUN_DOMAIN`, `MAILGUN_WEBHOOK_SIGNING_KEY`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` (default `http://localhost:3000/oauth/callback`), `TOKEN_ENCRYPTION_KEY`.

## Frontend Development Notes
- Prefer component extraction over large files. Keep components modular/reusable.
- Default to shadcn components. Available: Accordion, Alert Dialog, Alert, Aspect Ratio, Avatar, Badge, Breadcrumb, Button Group, Button, Calendar, Card, Carousel, Chart, Checkbox, Collapsible, Combobox, Command, Context Menu, Data Table, Date Picker, Dialog, Drawer, Dropdown Menu, Empty, Field, Form, Hover Card, Input Group, Input OTP, Input, Item, Kbd, Label, Menubar, Native Select, Navigation Menu, Pagination, Popover, Progress, Radio Group, Resizable, Scroll Area, Select, Separator, Sheet, Sidebar, Skeleton, Slider, Sonner, Spinner, Switch, Table, Tabs, Textarea, Toast, Toggle Group, Toggle, Tooltip, Typography.

## Migrations & Data
- AutoMigrate runs on start; `backend/MIGRATIONS.md` outlines migration strategy to move to `golang-migrate` with commands/examples.

## Testing
- Gmail OAuth + reply flow in `backend/TESTING.md` (login, get auth URL, complete browser OAuth, verify status, send reply).