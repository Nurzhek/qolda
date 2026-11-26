# qolda

Campus platform for renting, swapping and selling items between students of one university
(Senior Project, Team 4, Fall 2026). Monolith: one Go binary serves the REST API and the Vue app.

**Fall 2026 MVP scope** (see `docs/Requirements_and_Design_Team4`): university sign-in with a
domain allowlist (non-university accounts are rejected), roles, JWT; user profile; rent / exchange /
sale listings with photos; catalogue search and filters; PostgreSQL schema; OpenAPI; Docker Compose;
CI with tests for auth and listings. Orders, deposits and Kaspi, chat, ratings
and moderation are Spring 2027.

```
backend/   Go: gotemplate (Clean Architecture + DDD), gRPC + grpc-gateway, pgx/mobone   → backend/CLAUDE.md
web/       Vue 3 + Vite + Nuxt UI + Pinia, Feature-Sliced Design                     → web/CLAUDE.md
docs/      kickoff and requirements/design documents
```

## Run everything (Docker)

```bash
cp .env.example .env    # set JWT_SECRET (openssl rand -hex 32)
docker compose up --build
```

Open http://localhost:8090. With `AUTH_DEV_LOGIN=true` you can sign in with any `@nu.edu.kz`
address without Google. API docs (Swagger UI): http://localhost:8090/docs/; "Try it out" uses the
session of the web app signed in on the same browser.

Demo data (12 students, 38 listings; safe to re-run, optional photos in `demo/photos/`):

```bash
python3 demo/seed.py
```

## Develop

```bash
docker compose up -d db                       # Postgres on localhost:5435
cd backend && cp .env.example .env && go run ./cmd   # API on :8090
cd web && npx -y pnpm@10 install && npx -y pnpm@10 dev   # http://localhost:5173, proxies /api and /files
```

## Google sign-in

Create an OAuth client (Web application) in Google Cloud Console, add the redirect URI
`<PUBLIC_URL>/api/auth/google/callback`, put `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` into `.env`.
The e-mail domain is checked on our server after Google returns the profile, so a personal Gmail
account is rejected even if Google accepts it.

## Config (env)

| var | default | |
|---|---|---|
| `ALLOWED_DOMAINS` | `nu.edu.kz` | comma-separated; subdomains allowed |
| `AUTH_DEV_LOGIN` | `false` | e-mail sign-in without IdP (dev / pilot fallback) |
| `CAMPUSES` | `academic_block,dormitory,family_housing,off_campus` | campus codes |
| `JWT_SECRET` | — | required unless `DEBUG=true` |
| `ACCESS_TTL` / `REFRESH_TTL` | `15m` / `720h` | |
| `PUBLIC_URL` | `http://localhost:8090` | builds the Google redirect URI |
| `UPLOAD_DIR` | `./uploads` (`/data/uploads` in Docker) | listing photos |
