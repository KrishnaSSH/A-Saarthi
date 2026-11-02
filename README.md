# A-Saarthi

An accessible front end for Ayushman Bharat: multilingual, minimal navigation, automatic location detection, and plain-language rejection reasons. Svelte frontend, Go backend, talks to the official ABDM gateway ([docs](https://docs.abdm.gov.in/)).

## Run

```sh
cd backend && go run .                  # :8080
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
