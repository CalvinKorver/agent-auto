# Deployment Guide

This guide covers deploying the car-buyer application to production.

## Architecture

- **Frontend**: Vercel (Next.js)
- **Backend**: Railway or Render (containerized Go)
- **Database**: Neon PostgreSQL (serverless)

## Database Setup (Neon)

### 1. Create Neon Account
1. Go to [neon.tech](https://neon.tech) and sign up
2. Create a new project called "car-buyer"

### 2. Create Database Branches

**Production Branch** (main):
- Automatically created with project
- Use for production environment

**Preview Branch**:
```bash
# Create via Neon dashboard or CLI
neon branches create --name preview --parent main
```

**Development Branch** (optional):
```bash
# For team development, or use local PostgreSQL
neon branches create --name development --parent main
```

### 3. Get Connection Strings

For each branch, copy the connection string from Neon dashboard:
- Format: `postgres://[user]:[password]@[host]/[dbname]?sslmode=require`
- Store securely - never commit to git

## Environment Variables

### Local Development

**Backend** (`/backend/.env`):
```bash
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgresql://localhost:5432/carbuyer  # or Neon dev branch
JWT_SECRET=<generate-with-openssl-rand-base64-32>
JWT_EXPIRATION_HOURS=24
ANTHROPIC_API_KEY=sk-ant-xxx
ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT_AUTH=5
RATE_LIMIT_API=100
```

**Frontend** (`/frontend/.env.local`):
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_ENVIRONMENT=development
```

### Preview Environment

Set in Railway/Render and Vercel dashboards:

**Backend (Railway/Render)**:
- `ENVIRONMENT=preview`
- `DATABASE_URL=<neon-preview-branch-url>`
- `JWT_SECRET=<unique-preview-secret>`
- `ANTHROPIC_API_KEY=<your-key>`
- `ALLOWED_ORIGINS=https://car-buyer-git-*.vercel.app`
- All other vars from .env.example

**Frontend (Vercel Preview)**:
- `NEXT_PUBLIC_API_URL=https://your-backend-preview.railway.app/api/v1`
- `NEXT_PUBLIC_ENVIRONMENT=preview`

### Production Environment

**Backend (Railway/Render)**:
- `ENVIRONMENT=production`
- `DATABASE_URL=<neon-production-url>`
- `JWT_SECRET=<unique-production-secret>`
- `ANTHROPIC_API_KEY=<your-key>`
- `ALLOWED_ORIGINS=https://your-domain.com`
- All other vars from .env.example

**Frontend (Vercel Production)**:
- `NEXT_PUBLIC_API_URL=https://your-backend.railway.app/api/v1`
- `NEXT_PUBLIC_ENVIRONMENT=production`

## Backend Deployment

### Option A: Railway (Recommended)

1. **Install Railway CLI**:
```bash
npm i -g @railway/cli
railway login
```

2. **Create Project**:
```bash
cd backend
railway init
railway link
```

3. **Add Environment Variables**:
```bash
railway variables set DATABASE_URL="<neon-url>"
railway variables set JWT_SECRET="<secret>"
railway variables set ANTHROPIC_API_KEY="<key>"
railway variables set ALLOWED_ORIGINS="<vercel-url>"
# ... add all other variables
```

4. **Deploy**:
```bash
railway up
```

5. **Get URL**:
```bash
railway status
# Note the deployment URL
```

### Option B: Render

1. Go to [render.com](https://render.com)
2. Create new "Web Service"
3. Connect your GitHub repository
4. Configure:
   - **Root Directory**: `backend`
   - **Build Command**: `go build -o bin/server cmd/server/main.go`
   - **Start Command**: `./bin/server`
   - **Environment**: Add all variables from .env.example
5. Deploy

## Frontend Deployment (Vercel)

### 1. Connect Repository

1. Go to [vercel.com](https://vercel.com)
2. Import your GitHub repository
3. Configure:
   - **Framework Preset**: Next.js
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `.next`

### 2. Environment Variables

**Production**:
- Go to Project Settings → Environment Variables
- Add variables for "Production" environment:
  - `NEXT_PUBLIC_API_URL`: Your Railway/Render backend URL
  - `NEXT_PUBLIC_ENVIRONMENT`: `production`

**Preview**:
- Add same variables for "Preview" environment
- Use preview backend URL if available

### 3. Deploy

- Automatic on push to main (production)
- Automatic on PR (preview)

## Database Migrations

The app currently uses GORM AutoMigrate. For production:

### Recommended: Add Migration Versioning

1. Install migrate tool:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Create migrations:
```bash
migrate create -ext sql -dir backend/migrations -seq initial_schema
```

3. Update deployment to run migrations before starting server

## Post-Deployment Checklist

- [ ] All environment variables set correctly
- [ ] Database connections tested
- [ ] CORS configured for production domains
- [ ] JWT secrets are unique per environment
- [ ] Health check endpoint working
- [ ] Vercel deployment successful
- [ ] Backend deployment successful
- [ ] Test user registration flow
- [ ] Test authentication end-to-end
- [ ] Monitor logs for errors

## Monitoring & Logging

### Vercel
- Logs: Vercel Dashboard → Deployments → Logs
- Analytics: Built-in Vercel Analytics

### Railway
- Logs: `railway logs`
- Metrics: Railway Dashboard

### Neon
- Database metrics: Neon Dashboard → Monitoring
- Connection pooling: Enabled by default

## Security Notes

1. **Never commit secrets** - Use environment variables
2. **Rotate secrets regularly** - Especially JWT_SECRET
3. **Use different secrets** per environment
4. **Enable SSL** - Neon requires `?sslmode=require`
5. **Review CORS** - Only allow necessary origins
6. **Rate limiting** - Configure appropriately for production

## Scaling Considerations

- Neon: Auto-scales, consider upgrading plan if needed
- Railway: Auto-scales based on load
- Vercel: Automatically scales edge functions

## Troubleshooting

### "Database connection failed"
- Check DATABASE_URL format
- Ensure `?sslmode=require` for Neon
- Verify network access in Neon dashboard

### "CORS error"
- Check ALLOWED_ORIGINS matches frontend URL exactly
- Include protocol (https://)
- No trailing slash

### "JWT invalid"
- Ensure JWT_SECRET matches between deployments
- Check token expiration settings
