-- seed.sql
-- Initial data for development and testing

-- Organization: Provincia de Neuquén
INSERT INTO organizations (id, name, slug, config) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Provincia de Neuquén', 'neuquen', '{
        "topic_levels": 3,
        "shared_classes": true,
        "max_activities_per_class": 5,
        "document_sections": ["objectives", "content", "methodology", "evaluation", "resources"],
        "features": {
            "ai_chat": true,
            "shared_classes": true,
            "resource_library": false
        },
        "onboarding": {
            "skip_allowed": false,
            "profile_fields": [
                {"key": "disciplines", "label": "Disciplinas que enseña", "type": "multiselect", "options": ["Matemática", "Física", "Historia", "Lengua", "Biología", "Química", "Geografía"], "required": true},
                {"key": "experience_years", "label": "Años de experiencia", "type": "number", "required": true},
                {"key": "institution", "label": "Institución", "type": "text", "required": false},
                {"key": "education_level", "label": "Nivel educativo", "type": "select", "options": ["Inicial", "Primario", "Secundario", "Superior"], "required": true}
            ],
            "tour_steps": [
                {"key": "welcome", "title": "Bienvenido a Alizia", "description": "Alizia te ayuda a planificar el año escolar de forma colaborativa.", "order": 1},
                {"key": "explore", "title": "Explorá la plataforma", "description": "Navegá las secciones para descubrir las herramientas disponibles.", "order": 2},
                {"key": "coordination", "title": "Documento de coordinación", "description": "Creá y gestioná documentos de coordinación para tu área.", "order": 3, "roles": ["coordinator", "admin"]},
                {"key": "planning", "title": "Planificación docente", "description": "Armá tu planificación de clases a partir de la coordinación.", "order": 4, "roles": ["teacher"]},
                {"key": "ai_assistant", "title": "Asistente IA", "description": "Usá la IA para generar contenido y obtener sugerencias.", "order": 5, "requires_feature": "ai_chat"}
            ]
        }
    }')
ON CONFLICT (id) DO NOTHING;

-- Users
-- password_hash for all seeded users = argon2id('admin123', OWASP 2024 params)
-- Regenerate with: go run ./scripts/hash_password admin123
-- NB: argon2id is salted — each run returns a different hash. All 4 users
-- share the same hash string here only because the seed is deterministic for
-- local dev; in real usage each user must have its own salt/hash pair.
INSERT INTO users (id, organization_id, email, first_name, last_name, password_hash) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@neuquen.edu.ar',    'Ana',    'Admin',        '$argon2id$v=19$m=19456,t=2,p=1$KVZrrFCXd0/xP3whwMoErQ$sM8triBmp3RFIIpm0j6JPEMXuuoCa/JWUet61LyRw7c'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'coord@neuquen.edu.ar',    'Carlos', 'Coordinador',  '$argon2id$v=19$m=19456,t=2,p=1$KVZrrFCXd0/xP3whwMoErQ$sM8triBmp3RFIIpm0j6JPEMXuuoCa/JWUet61LyRw7c'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher1@neuquen.edu.ar', 'María',  'Docente',      '$argon2id$v=19$m=19456,t=2,p=1$KVZrrFCXd0/xP3whwMoErQ$sM8triBmp3RFIIpm0j6JPEMXuuoCa/JWUet61LyRw7c'),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher2@neuquen.edu.ar', 'Pedro',  'Multirol',     '$argon2id$v=19$m=19456,t=2,p=1$KVZrrFCXd0/xP3whwMoErQ$sM8triBmp3RFIIpm0j6JPEMXuuoCa/JWUet61LyRw7c')
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
