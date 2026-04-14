CREATE TABLE topics (
    id         BIGSERIAL    PRIMARY KEY,
    organization_id UUID    NOT NULL REFERENCES organizations(id),
    parent_id  BIGINT       REFERENCES topics(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    description TEXT,
    level      INTEGER      NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_topics_organization_id ON topics(organization_id);
CREATE INDEX idx_topics_parent_id ON topics(parent_id);
CREATE INDEX idx_topics_level ON topics(organization_id, level);
