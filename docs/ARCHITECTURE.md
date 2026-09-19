# Architecture

```text
Chrome MV3 extension
├─ Popup: profile, preferences, match, review
├─ Content script: detect job/forms; fill only after approval
├─ Service worker: orchestration, dedupe, notifications
└─ Local storage: profile, preferences, seen jobs, applications
          │ HTTPS
Go API
├─ AIProvider interface → GeminiProvider
├─ validation and safety policy
├─ PostgreSQL (planned persistence)
└─ Redis (planned dedupe/jobs)
```

## Trust boundaries

The extension never contains provider secrets. Page content is treated as untrusted input. The content script extracts text and form metadata but does not execute page-provided code. AI output is structured, validated, and displayed before use.

## Match score

The MVP score is deterministic and explainable: required/preferences gates plus skill, title, location, and experience signals. It is an internal relevance score—not a probability of interview or offer.
