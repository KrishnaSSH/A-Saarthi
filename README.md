# A-Saarthi

An accessible front end for Ayushman Bharat: multilingual, minimal navigation, automatic location detection, and plain-language rejection reasons. Svelte frontend, Go backend, talks to the official ABDM gateway ([docs](https://docs.abdm.gov.in/)).

## Run

```sh
go run ./cmd/server                  # :8080
cd frontend && npm i && npm run dev     # :5173, proxies /api
```

`npm run build` in `frontend` makes the Go server serve the app from `:8080`.

## ABDM credentials

```sh
export ABDM_CLIENT_ID=... ABDM_CLIENT_SECRET=...
export ABDM_BASE=https://dev.abdm.gov.in   # sandbox (default)
```

Without them, `/api/abha/search` returns 503; everything else works.

## API

| Route | Purpose |
|---|---|
| `GET /api/reasons?lang=` | Rejection code → why / how to fix |
| `GET /api/location?lat=&lon=` | Coordinates → structured Indian address |
| `POST /api/check` | Validate an application before submit |
| `POST /api/abha/search` | ABHA lookup via official gateway |

Add reasons in `backend/errors.go`; add languages in `frontend/src/lib/i18n.js`.

## Deploy (Vercel)

The frontend is built to static files; `api/index.go` serves the Go API as a serverless function (`vercel.json` rewrites `/api/*` to it). `.github/workflows/deploy.yml` tests, then deploys previews for PRs and production for `main`.

One-time setup:

1. `npx vercel link` in this repo, then read `.vercel/project.json` for `orgId` and `projectId`.
2. Add GitHub repo secrets: `VERCEL_TOKEN` (vercel.com/account/tokens), `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`.
3. Add `ABDM_CLIENT_ID`, `ABDM_CLIENT_SECRET` (and optionally `ABDM_BASE`) in the Vercel project's environment variables.
