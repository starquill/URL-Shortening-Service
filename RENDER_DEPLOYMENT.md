# Render.com Deployment Guide (100% FREE - No Credit Card!)

Deploy your URL Shortening Service to Render.com with ZERO cost and NO credit card.

## ✅ Why Render.com?

- ✅ **NO credit card required**
- ✅ Free web service (750 hours/month)
- ✅ Free PostgreSQL database
- ✅ Auto-deploy from GitHub
- ✅ Free SSL certificate
- ✅ Simple setup (no CLI needed)

## 🎯 Two Deployment Options

### Option A: With Redis (Free External)
- Render.com for app + PostgreSQL (free)
- Redis Labs for Redis cache (30MB free)

### Option B: Without Redis (Simpler)
- Render.com for app + PostgreSQL (free)
- No Redis (app still works, just slightly slower)

---

## 📦 Step 1: Prepare Your Code

First, let's make Redis optional in your app:

### Update config.go to handle missing Redis

Your app already gracefully handles Redis errors, so it will work without Redis!

### Create `render.yaml` (optional)

```yaml
services:
  - type: web
    name: url-shortener
    env: docker
    dockerfilePath: ./Dockerfile
    envVars:
      - key: SERVER_PORT
        value: 8080
      - key: DATABASE_URL
        fromDatabase:
          name: url-shortener-db
          property: connectionString
      - key: BASE_URL
        sync: false
      - key: REDIS_URL
        value: localhost:6379
      - key: REDIS_TTL
        value: 24h

databases:
  - name: url-shortener-db
    databaseName: url_shortening_service
    user: postgres
```

---

## 🚀 Step 2: Deploy via Render Dashboard (No CLI!)

### 2.1 Sign Up (No Credit Card)

1. Go to: https://render.com
2. Click "Get Started for Free"
3. Sign up with:
   - GitHub (recommended)
   - GitLab
   - Google
   - Or Email
4. **No credit card asked!** ✅

### 2.2 Push Code to GitHub

```bash
# If not already on GitHub
git remote add origin https://github.com/yourusername/url-shortener.git
git push -u origin main
```

### 2.3 Create PostgreSQL Database

1. In Render dashboard, click **"New +"**
2. Select **"PostgreSQL"**
3. Settings:
   - **Name**: `url-shortener-db`
   - **Database**: `url_shortening_service`
   - **User**: `postgres`
   - **Region**: Oregon (or closest to you)
   - **Instance Type**: **Free**
4. Click **"Create Database"**
5. ✅ Database created! (takes ~2 minutes)

### 2.4 Create Web Service

1. Click **"New +"** again
2. Select **"Web Service"**
3. Connect your GitHub repository
4. Settings:
   - **Name**: `url-shortener`
   - **Region**: Same as database
   - **Branch**: `main`
   - **Environment**: **Docker**
   - **Instance Type**: **Free**

### 2.5 Configure Environment Variables

In the "Environment Variables" section, add:

```
DATABASE_URL = (auto-filled from database)
BASE_URL = https://url-shortener.onrender.com
FRONTEND_URL = https://url-shortener.onrender.com
SERVER_PORT = 8080
REDIS_URL = localhost:6379
REDIS_TTL = 24h
```

**Note:** We use `localhost:6379` as placeholder. Redis won't be available but app will work.

### 2.6 Deploy!

1. Click **"Create Web Service"**
2. Render will:
   - Clone your repo
   - Build Docker image
   - Run database migrations
   - Deploy your app
3. Takes ~5-10 minutes
4. ✅ Your app is live!

---

## 🔴 Step 3: Add Redis (Optional)

If you want caching, use Redis Labs free tier:

### 3.1 Sign Up for Redis Cloud

1. Go to: https://redis.com/try-free/
2. Sign up (no credit card)
3. Create free database:
   - **Name**: url-shortener-cache
   - **Cloud**: AWS
   - **Region**: Same as Render
   - **Size**: 30MB (free)

### 3.2 Get Redis URL

1. In Redis Cloud dashboard, click your database
2. Copy the **Endpoint**: `redis-xxxxx.cloud.redislabs.com:12345`
3. Copy the **Password**

### 3.3 Update Render Environment

1. Go to your Render web service
2. Update environment variable:
   ```
   REDIS_URL = redis-xxxxx.cloud.redislabs.com:12345
   ```
3. Add new variable:
   ```
   REDIS_PASSWORD = your-password-here
   ```
4. Click "Save Changes"
5. Service will auto-redeploy

---

## ✅ Step 4: Test Your Deployment

### Get your URL

Your app will be at: `https://url-shortener.onrender.com`

### Test endpoints

```bash
# Health check
curl https://url-shortener.onrender.com/health

# Create short URL
curl -X POST https://url-shortener.onrender.com/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

# Test redirect
# Visit the short URL in browser
```

---

## 📊 Monitoring

### View Logs

1. Go to your service dashboard
2. Click "Logs" tab
3. See real-time logs

### View Metrics

Click "Metrics" tab to see:
- CPU usage
- Memory usage
- Request count

---

## 🔄 Auto-Deploy

Render automatically deploys when you push to GitHub!

```bash
# Make changes
git add .
git commit -m "update feature"
git push

# Render automatically:
# 1. Detects push
# 2. Builds new image
# 3. Deploys
# 4. Zero downtime!
```

---

## 💰 Cost Breakdown (FREE!)

| Resource | Render Free Tier | Your Usage | Cost |
|----------|------------------|------------|------|
| Web Service | 750 hours/month | ~720 hours | ₹0 |
| PostgreSQL | 1GB storage, 90 days free then pauses when inactive | <100MB | ₹0 |
| SSL Certificate | Included | Yes | ₹0 |
| Redis | Not included | Optional | ₹0 (use Redis Labs) |
| **TOTAL** | | | **₹0** |

**Note:** After 90 days, free databases pause when inactive. Just click "Resume" when needed.

---

## ⚠️ Limitations of Free Tier

1. **Sleeps after 15 min inactivity** - First request takes ~30 seconds to wake up
2. **750 hours/month** - Roughly 31 days, so keep it running 24/7
3. **Database pauses after 90 days** - Can resume anytime
4. **No Redis included** - Use external Redis Labs or skip it

---

## 🎓 Update Your Resume

```
"Deployed production URL shortening service to Render.com using 
Docker containerization, PostgreSQL database, automated CI/CD from 
GitHub, and external Redis caching. Application serves users globally 
with auto-scaling and zero-downtime deployments."
```

---

## 🐛 Troubleshooting

### Build fails?
- Check Render logs
- Verify Dockerfile is correct
- Ensure all dependencies in go.mod

### Database connection fails?
- Check DATABASE_URL environment variable
- Ensure database is in same region
- Check database logs in Render dashboard

### App sleeps after inactivity?
- This is normal on free tier
- Use UptimeRobot.com (free) to ping your app every 5 min
- Keeps it awake 24/7

### Need to run migrations manually?
- Render runs them automatically on startup
- Or connect to shell: Click "Shell" in dashboard

---

## 🚀 Next Steps

1. ✅ Deploy to Render
2. ✅ Test all endpoints
3. ✅ (Optional) Add Redis Labs
4. ✅ Set up UptimeRobot to keep it awake
5. ✅ Add to portfolio

**Your live URL:**
`https://your-app-name.onrender.com`

**Total cost: ₹0** 🎉

---

## Alternative: Deploy via Blueprint (Even Easier!)

Create `render.yaml` in your repo root:

```yaml
services:
  - type: web
    name: url-shortener
    env: docker
    dockerfilePath: ./Dockerfile
    plan: free
    envVars:
      - key: SERVER_PORT
        value: 8080
      - key: DATABASE_URL
        fromDatabase:
          name: url-shortener-db
          property: connectionString
      - key: REDIS_URL
        value: localhost:6379

databases:
  - name: url-shortener-db
    databaseName: url_shortening_service
    plan: free
```

Then:
1. Push to GitHub
2. In Render, click "New" → "Blueprint"
3. Connect repo
4. Click "Apply"
5. Done!
