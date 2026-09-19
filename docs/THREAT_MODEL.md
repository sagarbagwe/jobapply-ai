# Threat model

| Risk | Control |
|---|---|
| Provider key extraction | Keys exist only in backend environment |
| Prompt injection in job pages | Page text is untrusted; structured outputs, fixed instructions, user review |
| Fabricated applicant facts | Grounding-only prompts, missing fields, review gate |
| Cross-site form mistakes | Semantic mapping, confidence fallback, no submit action |
| Abuse and cost spikes | Request/body limits, per-IP rate limit, production bearer token |
| Sensitive retention | No raw PDF persistence; deletion controls; no-store headers |
| Duplicate applications | Canonical URL and company/title checks |
| Automated-site policy violations | User-initiated detection/fill only; no CAPTCHA or anti-bot bypass |
