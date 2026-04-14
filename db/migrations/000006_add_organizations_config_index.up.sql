-- GIN index on organizations.config for efficient JSONB queries
CREATE INDEX IF NOT EXISTS idx_organizations_config ON organizations USING GIN (config);
