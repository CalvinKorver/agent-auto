# Production Deployment Quickstart

Follow these steps to deploy your car-buyer application to production.

## Prerequisites

- [ ] GitHub repository set up
- [ ] Neon account (for PostgreSQL)
- [ ] Vercel account (for frontend)
- [ ] Railway or Render account (for backend)
- [ ] Anthropic API key

## Step 1: Database Setup (5 minutes)

### Create Neon Project

1. Go to [neon.tech](https://neon.tech) and sign up
2. Create new project: **car-buyer**
3. Note the connection string from the dashboard

### Create Database Branches

**Production** (already created):
- Use the main branch
- Copy connection string

**Preview** (for PR deployments):
```bash
# In Neon dashboard: Create Branch
# Name: preview
# Parent: main
```

**Local Development** (optional):
- Either use local PostgreSQL
- Or create another Neon branch named "development"

### Save Connection Strings
Store these securely - you'll need them next:
- Production: `postgres://user:pass@host/db?sslmode=require`
- Preview: `postgres://user:pass@preview-host/db?sslmode=require`

## Step 2: Backend Deployment (10 minutes)

### Option A: Railway (Recommended)

#### Using Railway Dashboard (Easier for First Deploy)

1. **Create Project via Dashboard**:
   - Go to [railway.app](https://railway.app) and sign in
   - Click "New Project" → "Deploy from GitHub repo"
   - Connect your GitHub repository
   - Select the `car-buyer` repository
   - Railway will detect the Dockerfile automatically

2. **Configure Service**:
   - Railway will auto-detect the `Dockerfile` in the root
   - **Name**: car-buyer-backend
   - **Root Directory**: Leave empty (uses project root)
   - Railway will automatically start building with Docker
   - Wait for initial build to complete (may take 2-3 minutes)

3. **Add Environment Variables**:
   - Go to your service → "Variables" tab
   - Click "New Variable" and add each:
     - `PORT` = `8080`
     - `ENVIRONMENT` = `production`
     - `DATABASE_URL` = `<your-neon-production-url>`
     - `JWT_SECRET` = Click "Generate" or use: `openssl rand -base64 32`
     - `JWT_EXPIRATION_HOURS` = `24`
     - `ANTHROPIC_API_KEY` = `<your-anthropic-key>`
     - `ALLOWED_ORIGINS` = `http://localhost:3000` (will update after Vercel)
     - `RATE_LIMIT_AUTH` = `5`
     - `RATE_LIMIT_API` = `100`

4. **Generate Domain**:
   - Go to "Settings" tab → "Networking"
   - Click "Generate Domain"
   - Copy the domain (e.g., `car-buyer-backend-production.up.railway.app`)

5. **Verify Deployment**:
```bash
curl https://your-domain.railway.app/health
# Should return: {"status":"healthy","database":"connected"}
```

#### Alternative: Using Railway CLI

1. **Install and Login**:
```bash
npm i -g @railway/cli
railway login
```

2. **Link to Existing Project** (after creating via dashboard):
```bash
cd /Users/calvinkorver/car-buyer
railway link
# Select your project from the list
```

3. **Set Variables via CLI** (optional):
```bash
railway variables set PORT=8080
railway variables set ENVIRONMENT=production
# ... etc
```

4. **Deploy via CLI** (optional):
```bash
railway up
```

### Option B: Render

1. Go to [render.com](https://render.com) and sign in
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Configure:
   - **Name**: car-buyer-backend
   - **Root Directory**: backend
   - **Environment**: Go
   - **Build Command**: `go build -o bin/server cmd/server/main.go`
   - **Start Command**: `./bin/server`
5. Add environment variables (same as Railway above)
6. Click "Create Web Service"
7. Copy the deployment URL

## Step 3: Frontend Deployment (5 minutes)

### Deploy to Vercel

1. Go to [vercel.com](https://vercel.com) and sign in
2. Click "Add New..." → "Project"
3. Import your GitHub repository
4. Configure:
   - **Framework Preset**: Next.js (auto-detected)
   - **Root Directory**: frontend
   - **Build Command**: `npm run build` (default)

5. **Add Environment Variables** (Production):
   - `NEXT_PUBLIC_API_URL`: `https://your-backend-url.railway.app/api/v1`
   - `NEXT_PUBLIC_ENVIRONMENT`: `production`

6. Click "Deploy"
7. Copy your Vercel URL (e.g., `https://car-buyer.vercel.app`)

### Update CORS

Go back to Railway/Render and update:
```bash
railway variables set ALLOWED_ORIGINS="https://car-buyer.vercel.app"
# Or add both www and non-www:
# "https://car-buyer.vercel.app,https://www.car-buyer.vercel.app"
```

## Step 4: Preview Environments (10 minutes)

### Backend Preview (Railway)

```bash
# Create preview environment
railway environment create preview

# Link to preview branch
railway link --environment preview

# Set preview variables (same as production but different DB)
railway variables set DATABASE_URL="<neon-preview-branch-url>" --environment preview
railway variables set ENVIRONMENT=preview --environment preview
railway variables set ALLOWED_ORIGINS="https://car-buyer-git-*.vercel.app" --environment preview
# ... add all other variables
```

### Frontend Preview (Vercel)

1. In Vercel project settings → Environments
2. Add variables for **Preview** environment:
   - `NEXT_PUBLIC_API_URL`: `https://your-backend-preview.railway.app/api/v1`
   - `NEXT_PUBLIC_ENVIRONMENT`: `preview`

Vercel automatically creates preview deployments for PRs!

## Step 5: Verify Deployment (5 minutes)

### Health Check
```bash
curl https://your-backend-url.railway.app/health
# Should return: {"status":"healthy","database":"connected"}
```

### Test the App
1. Visit your Vercel URL
2. Register a new account
3. Set preferences
4. Create a thread
5. Send a message

### Monitor Logs

**Vercel**:
- Dashboard → Deployments → Click deployment → Logs

**Railway**:
```bash
railway logs
```

## Step 6: Custom Domain (Optional)

### Add Domain to Vercel
1. Vercel Dashboard → Project Settings → Domains
2. Add your domain
3. Follow DNS configuration instructions

### Update CORS
```bash
railway variables set ALLOWED_ORIGINS="https://yourdomain.com,https://car-buyer.vercel.app"
```

## Environment Variables Checklist

### Backend (Railway/Render)

**Production**:
- [ ] `PORT=8080`
- [ ] `ENVIRONMENT=production`
- [ ] `DATABASE_URL=<neon-production-url>`
- [ ] `JWT_SECRET=<unique-secret>`
- [ ] `JWT_EXPIRATION_HOURS=24`
- [ ] `ANTHROPIC_API_KEY=<your-key>`
- [ ] `ALLOWED_ORIGINS=<vercel-url>`
- [ ] `RATE_LIMIT_AUTH=5`
- [ ] `RATE_LIMIT_API=100`

**Preview** (same variables, different values):
- [ ] `DATABASE_URL=<neon-preview-url>`
- [ ] `ENVIRONMENT=preview`
- [ ] `ALLOWED_ORIGINS=https://car-buyer-git-*.vercel.app`

### Frontend (Vercel)

**Production**:
- [ ] `NEXT_PUBLIC_API_URL=<backend-url>/api/v1`
- [ ] `NEXT_PUBLIC_ENVIRONMENT=production`

**Preview**:
- [ ] `NEXT_PUBLIC_API_URL=<backend-preview-url>/api/v1`
- [ ] `NEXT_PUBLIC_ENVIRONMENT=preview`

## Troubleshooting

### "Failed to connect to database"
- Verify DATABASE_URL has `?sslmode=require` for Neon
- Check Neon dashboard for connection string
- Ensure database branch exists

### "CORS policy error"
- Verify ALLOWED_ORIGINS matches frontend URL exactly
- No trailing slash in URL
- Include protocol (https://)

### "JWT token invalid"
- Ensure JWT_SECRET is the same across deployments
- Check if token expired (JWT_EXPIRATION_HOURS)

### Backend not responding
- Check health endpoint: `curl <backend-url>/health`
- View logs in Railway/Render dashboard
- Verify all required env vars are set

## Next Steps

- [ ] Set up monitoring (Railway metrics, Vercel Analytics)
- [ ] Configure custom domain
- [ ] Set up database backups (Neon automatic)
- [ ] Add error tracking (Sentry, etc.)
- [ ] Set up CI/CD pipeline (.github/workflows)
- [ ] Configure database migrations

## Costs (Free Tier)

- **Neon**: Free tier includes 3 databases
- **Railway**: $5/month after free credits
- **Vercel**: Free for personal projects
- **Total**: ~$5-10/month to start

## Support

If you run into issues:
1. Check deployment logs
2. Verify all environment variables
3. Test health endpoint
4. Review DEPLOYMENT.md for detailed documentation
