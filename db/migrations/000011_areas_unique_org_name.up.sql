-- Enforce uniqueness of area name within an organization (HU-3.2 acceptance criteria).
ALTER TABLE areas
    ADD CONSTRAINT areas_organization_name_unique UNIQUE (organization_id, name);
