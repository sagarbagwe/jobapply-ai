# Security policy

- Keep Gemini and other provider keys only on the backend.
- Never commit `.env`, resumes, exports, or personal application data.
- Final application submission always requires explicit user action.
- Do not bypass CAPTCHA, login controls, rate limits, robots directives, or anti-bot measures.
- Minimize retention: store structured resume data only; delete raw uploads after parsing unless the user opts in.
- Encrypt sensitive fields at rest in production and use TLS in transit.
- Restrict CORS to published extension IDs before production deployment.

Report vulnerabilities privately to the repository owner rather than opening a public issue.
