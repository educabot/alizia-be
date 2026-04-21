CREATE TYPE class_moment AS ENUM ('apertura', 'desarrollo', 'cierre');

CREATE TABLE activities (
    id               BIGSERIAL    PRIMARY KEY,
    organization_id  UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    moment           class_moment NOT NULL,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    duration_minutes INTEGER,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activities_organization ON activities(organization_id);
CREATE INDEX idx_activities_moment ON activities(organization_id, moment);
