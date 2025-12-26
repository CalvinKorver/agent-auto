# Railway Deployment Setup

Quick reference for deploying the backend to Railway.

## Step 1: Create Project via Dashboard

1. Go to [railway.app](https://railway.app)
2. Sign in with GitHub
3. Click "New Project"
4. Select "Deploy from GitHub repo"
5. Choose your `car-buyer` repository
6. Railway will automatically detect the `Dockerfile`

## Step 2: Configure Environment Variables

After deployment starts, add these variables in the Railway dashboard:

Go to your service → **Variables** tab → Click "New Variable"

### Required Variables

```bash
PORT=8080
ENVIRONMENT=production
DATABASE_URL=<your-neon-production-connection-string>
JWT_SECRET=<click-generate-or-use-openssl-rand-base64-32>
JWT_EXPIRATION_HOURS=24
ANTHROPIC_API_KEY=<your-anthropic-api-key>
ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT_AUTH=5
RATE_LIMIT_API=100
```

### Optional Variables (for email features)

```bash
MAILGUN_API_KEY=
MAILGUN_DOMAIN=
MAILGUN_WEBHOOK_SIGNING_KEY=
```

## Step 3: Generate Public Domain

1. Go to **Settings** tab
2. Navigate to **Networking** section
3. Click "Generate Domain"
4. Copy the domain (e.g., `car-buyer-backend-production.up.railway.app`)
5. Save this URL - you'll need it for Vercel frontend configuration

## Step 4: Verify Deployment

Test the health endpoint:

```bash
curl https://your-domain.railway.app/health
```

Expected response:
```json
{"status":"healthy","database":"connected"}
```

## Step 5: Update CORS After Vercel Deploy

Once you deploy the frontend to Vercel and get the URL:

1. Go back to Railway → **Variables** tab
2. Update `ALLOWED_ORIGINS`:
   ```
   https://your-app.vercel.app
   ```
3. For multiple origins (including preview):
   ```
   https://your-app.vercel.app,https://car-buyer-git-*.vercel.app
   ```

## Troubleshooting

### Build Fails

**Check logs**:
1. Go to Deployments tab
2. Click the failed deployment
3. View build logs

**Common issues**:
- Missing Dockerfile → Ensure `Dockerfile` exists in root
- Wrong path → `railway.json` should reference `Dockerfile` (not `backend/Dockerfile`)
- Missing dependencies → Ensure `go.mod` and `go.sum` are committed

### Deployment Starts but Crashes

**Check runtime logs**:
1. Go to Deployments tab
2. Click the deployment
3. View runtime logs

**Common issues**:
- Missing environment variables → Check all required vars are set
- Database connection failed → Verify `DATABASE_URL` format with `?sslmode=require`
- Wrong port → Ensure `PORT=8080` is set

### Health Check Fails

```bash
# Check if service is running
curl https://your-domain.railway.app/health

# Check API root
curl https://your-domain.railway.app/api/v1
```

**If database connection fails**:
- Verify Neon database is running
- Check connection string has `?sslmode=require`
- Ensure database allows connections from Railway IPs

## Using Railway CLI (Optional)

### Install
```bash
npm i -g @railway/cli
railway login
```

### Link Project
```bash
cd /Users/calvinkorver/car-buyer
railway link
# Select your project
```

### View Logs
```bash
railway logs
```

### Set Variables via CLI
```bash
railway variables set KEY=value
```

### Deploy from CLI
```bash
railway up
```

## Creating Preview Environment

For PR deployments with separate database:

1. **Create Preview Environment**:
   - Railway Dashboard → Environments → "New Environment"
   - Name: `preview`

2. **Configure Preview Variables**:
   - Same as production, but:
   - `ENVIRONMENT=preview`
   - `DATABASE_URL=<neon-preview-branch-url>`
   - `ALLOWED_ORIGINS=https://car-buyer-git-*.vercel.app`

3. **Auto-Deploy PRs**:
   - Railway will automatically deploy PRs to preview environment
   - Each PR gets its own deployment URL

## Monitoring

### View Metrics
- Railway Dashboard → Your Service → Metrics
- CPU, Memory, Network usage

### View Logs
- Railway Dashboard → Your Service → Deployments → Latest → Logs
- Or use CLI: `railway logs --tail`

### Set Up Alerts
- Railway Dashboard → Settings → Notifications
- Configure webhook or email alerts

## Cost Management

### Free Tier
- $5 credit per month
- Enough for hobby projects

### Paid Plans
- Pay for what you use
- ~$5-10/month for small apps
- View usage in Dashboard → Billing

### Optimize Costs
- Use smallest container size that works
- Enable auto-sleep for dev environments
- Use Neon's free tier for database

## Next Steps

1. ✅ Deploy backend to Railway
2. ✅ Get backend URL
3. → Deploy frontend to Vercel with backend URL
4. → Update CORS in Railway with Vercel URL
5. → Test end-to-end flow
6. → Set up preview environments

See [QUICKSTART.md](QUICKSTART.md) for complete deployment guide.
