-- Areas (agrupación de disciplinas)
CREATE TABLE areas (
    id BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Subjects (disciplinas pertenecientes a un área)
CREATE TABLE subjects (
    id BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Area coordinators (M2M: un coordinador puede coordinar múltiples áreas)
CREATE TABLE area_coordinators (
    id BIGSERIAL PRIMARY KEY,
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(area_id, user_id)
);

-- Indexes
CREATE INDEX idx_areas_organization_id ON areas(organization_id);
CREATE INDEX idx_subjects_organization_id ON subjects(organization_id);
CREATE INDEX idx_subjects_area_id ON subjects(area_id);
CREATE INDEX idx_area_coordinators_area_id ON area_coordinators(area_id);
CREATE INDEX idx_area_coordinators_user_id ON area_coordinators(user_id);
