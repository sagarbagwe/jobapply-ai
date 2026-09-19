# Roadmap and release status

## Implemented
- [x] MV3 extension, editable profile/preferences, PDF resume parsing, multi-resume metadata
- [x] Structured job detection, explainable matching, duplicate detection, notifications
- [x] Review-first form filling and manual final submission boundary
- [x] Local tracker plus authenticated PostgreSQL tracker synchronization
- [x] User registration/login with short-lived signed JWT access tokens
- [x] Tenant-scoped profile and application APIs
- [x] Server-side data export and confirmed account deletion
- [x] Gemini provider isolation, strict structured outputs, no raw-PDF persistence
- [x] Rate/body limits, exact-origin CORS, security/no-store headers
- [x] PostgreSQL schema, Docker deployment, CI race tests, audit, and release packaging
- [x] Privacy design, threat model, security policy, and deployment checklist

## Operational release gates
These require the production operator or external reviewers rather than additional repository code:
- [ ] Provision the TLS API hostname and secrets manager
- [ ] Set the final Chrome extension ID in `ALLOWED_ORIGINS` and store manifest
- [ ] Configure encrypted PostgreSQL, backups, restore drills, monitoring, and spending alerts
- [ ] Run permitted ATS compatibility testing and independent security/privacy review
- [ ] Complete Chrome Web Store disclosures and review

No release may bypass CAPTCHA, authentication, rate limits, anti-bot controls, or the user's final Submit action.
