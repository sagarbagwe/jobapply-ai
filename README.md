# JobApply AI

Human-in-the-loop AI job discovery and application assistant.

> **Safety:** JobApply AI never bypasses CAPTCHA, authentication, rate limits, or anti-bot controls, and never submits an application without an explicit user click.

## MVP

- Chrome Manifest V3 extension (React + TypeScript)
- Editable profile and job preferences stored in `chrome.storage.local`
- Job-page detection and structured extraction
- Deterministic local match scoring
- Application-form field detection and review-first autofill
- Go REST API with pluggable `AIProvider`
- Gemini provider; API keys stay on the server

## Monorepo

```text
extension/  Chrome MV3 extension
backend/    Go API
infra/      Docker Compose (PostgreSQL + Redis)
docs/       Architecture and roadmap
```

## Quick start

### Extension

```bash
cd extension
npm install
npm run build
```

Load `extension/dist` via `chrome://extensions` → **Developer mode** → **Load unpacked**.

### Backend

```bash
cp .env.example .env
# Add GEMINI_API_KEY. Never commit .env.
docker compose up --build
```

Health check: `GET http://localhost:8080/healthz`.

## Privacy

Resume-derived data remains local in the MVP unless the user explicitly invokes an AI feature. Raw resumes are not retained by the backend. See [SECURITY.md](SECURITY.md).

## Status

This is the production-oriented foundation for Phase 1. See [docs/ROADMAP.md](docs/ROADMAP.md).
