# LeBlanc

Leblanc is split into a Go API (`server/`) and a Vue 3 front-end (`website/`). The API uses MongoDB for persistence.

## Prerequisites

- Go 1.24+
- Node 20+ (for the Vue app)
- MongoDB running locally (default connection string `mongodb://127.0.0.1:27017`)

## Backend setup

```bash
cd server
cp .env.example .env # or edit .env with your credentials
go run ./cmd/seed    # creates the database, indexes, and sample drinks
go run .
```

The seed command uses the `.env` variables `MONGO_URI` and `MONGO_DB` to create:

- unique indexes on `users.nameLower` and `users.emailLower`
- a curated set of drinks that power the `/drinks` and recommendation endpoints

If you re-run the seeder it will skip inserting drinks when the collection is not empty.

## Front-end

```bash
cd website/LeBlanc\ web
npm install   
npm run dev         ## for local
vercel --prod       ## for deployed 
```

Configure `VITE_API_BASE` in `website/LeBlanc web/.env` if the backend is not on `http://localhost:4000`.

## Deployments

- Frontend (Vercel): https://le-blanc-web.vercel.app  
  - SPA rewrite lives in `website/LeBlanc web/vercel.json`.
  - Deploy from `website/LeBlanc web` with `vercel --prod` after setting `VITE_API_BASE` (point it to the API), `VITE_EMAILJS_*`, and `VITE_ADMIN_EMAIL`.
- Backend (Fly.io): https://server-wandering-tree-4946.fly.dev  
  - Uses `server/Dockerfile` and `server/fly.toml` (internal port `8080`).
  - Deploy with `flyctl deploy --config fly.toml --dockerfile Dockerfile` from `server/` and set secrets such as `MONGO_URI`, `MONGO_DB`, `ADMIN_NAME`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `FRONTEND_VERIFY_URL=https://le-blanc-web.vercel.app/verify`, and `EMAIL_REQUIRE_MX`.
