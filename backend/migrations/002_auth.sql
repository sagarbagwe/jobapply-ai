ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash text;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
CREATE INDEX IF NOT EXISTS users_email_live_idx ON users(lower(email)) WHERE deleted_at IS NULL;
