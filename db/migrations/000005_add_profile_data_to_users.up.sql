ALTER TABLE users ADD COLUMN profile_data JSONB NOT NULL DEFAULT '{}';

COMMENT ON COLUMN users.profile_data IS 'Dynamic profile fields captured during onboarding. Schema defined by organization config.';
