-- seed.sql
-- Initial data for development and testing

-- Organization: Provincia de Neuquén
INSERT INTO organizations (id, name, slug, config) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Provincia de Neuquén', 'neuquen', '{
        "topic_levels": 3,
        "shared_classes": true,
        "max_activities_per_class": 5,
        "document_sections": ["objectives", "content", "methodology", "evaluation", "resources"]
    }')
ON CONFLICT (id) DO NOTHING;

-- Users
INSERT INTO users (id, organization_id, email, first_name, last_name) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@neuquen.edu.ar', 'Ana', 'Admin'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'coord@neuquen.edu.ar', 'Carlos', 'Coordinador'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher1@neuquen.edu.ar', 'María', 'Docente'),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher2@neuquen.edu.ar', 'Pedro', 'Multirol')
ON CONFLICT (id) DO NOTHING;

-- Roles
INSERT INTO user_roles (user_id, role) VALUES
    (1, 'admin'),
    (2, 'coordinator'),
    (3, 'teacher'),
    (4, 'teacher'),
    (4, 'coordinator')
ON CONFLICT (user_id, role) DO NOTHING;

-- Areas
INSERT INTO areas (id, organization_id, name, description) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Ciencias', 'Área de ciencias exactas y naturales'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Humanidades', 'Área de ciencias sociales y humanidades')
ON CONFLICT (id) DO NOTHING;

-- Subjects
INSERT INTO subjects (id, organization_id, area_id, name) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Matemática'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Física'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Historia'),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Lengua')
ON CONFLICT (id) DO NOTHING;

-- Area coordinators (Carlos coordina Ciencias, Pedro coordina Humanidades)
INSERT INTO area_coordinators (area_id, user_id) VALUES
    (1, 2),
    (2, 4)
ON CONFLICT (area_id, user_id) DO NOTHING;

-- Reset sequences to avoid conflicts with future inserts
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 0) FROM users));
SELECT setval('user_roles_id_seq', (SELECT COALESCE(MAX(id), 0) FROM user_roles));
SELECT setval('areas_id_seq', (SELECT COALESCE(MAX(id), 0) FROM areas));
SELECT setval('subjects_id_seq', (SELECT COALESCE(MAX(id), 0) FROM subjects));
SELECT setval('area_coordinators_id_seq', (SELECT COALESCE(MAX(id), 0) FROM area_coordinators));
