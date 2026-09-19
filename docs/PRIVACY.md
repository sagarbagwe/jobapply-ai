# Privacy design

## Data minimization
- The extension stores profile, resume-derived structured data, preferences, seen jobs, and application history locally.
- Raw resume PDFs are sent to the configured backend only when the user chooses a file and are not persisted by the API.
- Job descriptions are sent only when the user requests AI matching or answer generation.

## User controls
Users can edit extracted information and delete all local data from the extension. Production accounts must provide server-side export and deletion before general availability.

## Human control
Generated answers require review. Form filling is initiated by the user. Final submission is never automated.
