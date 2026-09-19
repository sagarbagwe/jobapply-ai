# Roadmap and release gates

## Implemented in the beta foundation
- [x] MV3 extension with editable profile and preferences
- [x] PDF resume parsing through Gemini with no raw-PDF persistence
- [x] Multi-resume metadata
- [x] JSON-LD and heuristic job detection
- [x] Explainable local and Gemini matching
- [x] Review-first semantic form filling
- [x] Local application tracker and duplicate detection
- [x] Saved seen-job IDs and opt-in match notifications
- [x] Gemini provider abstraction endpoints
- [x] PostgreSQL production schema and Docker initialization
- [x] Request limits, rate limiting, API authentication hook, security headers
- [x] Privacy and threat-model documentation

## Required before public production release
- [ ] Replace development bearer token with user OAuth/JWT and tenant-scoped APIs
- [ ] Wire tracker CRUD to PostgreSQL with row-level authorization
- [ ] Add server-side export/deletion jobs and encryption key management
- [ ] Run ATS adapter compatibility tests against permitted targets
- [ ] Add Playwright E2E, accessibility, load, and restore tests
- [ ] Complete Chrome Web Store privacy disclosures and external security review

The repository is a testable beta, not authorization to automate a site that prohibits automation.
