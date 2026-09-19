# Production deployment checklist

## Required configuration
- Set `APP_ENV=production`.
- Use a randomly generated `JWT_SECRET` of at least 32 characters from a secrets manager.
- Set `GEMINI_API_KEY` only on the backend.
- Set `ALLOWED_ORIGINS` to the exact published Chrome extension origin.
- Build the extension with `VITE_API_BASE_URL` pointing to the TLS API endpoint.
- Replace development database credentials and require TLS.

## Release checks
- CI extension build, dependency audit, Go vet, race tests, and Docker build pass.
- Apply database migrations before switching traffic.
- Configure backups, restore drills, log retention, alerts, and provider spending limits.
- Complete Chrome Web Store privacy disclosures.
- Test only against supported sites and document their permitted interaction model.
- Commission an external security and privacy review.

## Non-negotiable behavior
No CAPTCHA bypass, credential bypass, rate-limit evasion, restricted scraping, fabricated applicant facts, or automatic final submission.
