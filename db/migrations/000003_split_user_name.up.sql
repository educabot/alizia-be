-- 000003_split_user_name.up.sql
-- Splits users.name into first_name + last_name.
--
-- Idempotent by design: migration 000001 was later consolidated to create
-- first_name/last_name directly, so a fresh install already has the target
-- shape and this migration must be a no-op. Existing environments that
-- ran 000001's original form (with a single `name` column) still get the
-- rename on their next migrate run.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'name'
    ) THEN
        ALTER TABLE users RENAME COLUMN name TO first_name;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'last_name'
    ) THEN
        ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NOT NULL DEFAULT '';
    END IF;
END $$;
