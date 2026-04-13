ALTER TABLE users ADD COLUMN onboarding_completed_at TIMESTAMP NULL;

COMMENT ON COLUMN users.onboarding_completed_at IS 'Timestamp when user completed onboarding flow. NULL = not completed.';
