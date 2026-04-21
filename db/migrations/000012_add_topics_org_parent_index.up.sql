-- Composite index for org-scoped subtree lookups (descendantsOf, GetTopicTree).
-- Convention: every foreign key MUST have an index. Composite indexes on
-- (organization_id, <fk>) take precedence when queries are tenant-scoped.
CREATE INDEX IF NOT EXISTS idx_topics_org_parent ON topics(organization_id, parent_id);
