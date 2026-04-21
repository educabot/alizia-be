-- seed.sql
-- Initial data for development and testing

-- Organization: Provincia de Neuquén
INSERT INTO organizations (id, name, slug, config) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Provincia de Neuquén', 'neuquen', '{
        "topic_max_levels": 3,
        "topic_level_names": ["Núcleos problemáticos", "Áreas de conocimiento", "Categorías"],
        "topic_selection_level": 3,
        "shared_classes_enabled": true,
        "desarrollo_max_activities": 3,
        "coord_doc_sections": [
            {
                "key": "problem_edge",
                "label": "Eje problemático",
                "type": "text",
                "ai_prompt": "Generá un eje problemático que articule las categorías seleccionadas para las disciplinas del área, considerando el contexto educativo de nivel secundario.",
                "required": true
            },
            {
                "key": "methodological_strategy",
                "label": "Estrategia metodológica",
                "type": "select_text",
                "options": ["proyecto", "taller_laboratorio", "ateneo_debate"],
                "ai_prompt": "Generá una estrategia metodológica detallada para abordar el eje problemático, explicando cómo se articulan las disciplinas del área y qué actividades concretas se proponen.",
                "required": true
            },
            {
                "key": "eval_criteria",
                "label": "Criterios de evaluación",
                "type": "text",
                "ai_prompt": "Generá criterios de evaluación que permitan valorar el aprendizaje de las categorías seleccionadas, considerando evaluación formativa y sumativa.",
                "required": false
            }
        ],
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
-- password_hash for all seeded users = argon2id('admin123', OWASP 2024 params).
-- Regenerate with: go run ./scripts/hash_password admin123
-- Each user has its own salt/hash so the seed mirrors production behavior and
-- avoids the anti-pattern of reusing argon2id hashes across accounts.
INSERT INTO users (id, organization_id, email, first_name, last_name, password_hash) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@neuquen.edu.ar',    'Ana',    'Admin',        '$argon2id$v=19$m=19456,t=2,p=1$VaKh6LQ9ytjvW3XEJ5WqCQ$X7V3mReWhC3CWAt7DB7gbxnjevO+PFa1GbSD1q77MHA'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'coord@neuquen.edu.ar',    'Carlos', 'Coordinador',  '$argon2id$v=19$m=19456,t=2,p=1$El4xWACc0BOEvclmu81kLQ$hHzBQxlsJfAteoGS4ttimQSYT16txqYjopSpKecPtN4'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher1@neuquen.edu.ar', 'María',  'Docente',      '$argon2id$v=19$m=19456,t=2,p=1$e+Pr+aRlNh3FyXDe5VQ5Zg$FqhB48JWwnYjCV4OHk8HeNjhjs9Rb0gIiC8I1oPI0eY'),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'teacher2@neuquen.edu.ar', 'Pedro',  'Multirol',     '$argon2id$v=19$m=19456,t=2,p=1$l/Dgy/ijeQtapaWaq9N8PQ$9qtRAuWVcrOiaqfvZwpnhfJZJdzyg2p9/b1v8qYpbRc')
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

-- Topics: 3-level hierarchy (Núcleos → Áreas de conocimiento → Categorías)
-- Level 1: Núcleos problemáticos
INSERT INTO topics (id, organization_id, parent_id, name, description, level) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'Pensamiento Lógico-Matemático', 'Núcleo orientado al desarrollo del razonamiento lógico', 1),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'Comunicación y Lenguaje', 'Núcleo orientado a competencias comunicativas', 1)
ON CONFLICT (id) DO NOTHING;

-- Level 2: Áreas de conocimiento
INSERT INTO topics (id, organization_id, parent_id, name, description, level) VALUES
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Aritmética Básica', 'Operaciones fundamentales y propiedades de números', 2),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Geometría y Medición', 'Figuras, cuerpos geométricos y unidades de medida', 2),
    (5, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Comprensión Lectora', 'Estrategias de lectura y análisis de textos', 2),
    (6, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Producción Escrita', 'Redacción, coherencia y cohesión textual', 2)
ON CONFLICT (id) DO NOTHING;

-- Level 3: Categorías
INSERT INTO topics (id, organization_id, parent_id, name, description, level) VALUES
    (7,  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Suma y resta', 'Operaciones de adición y sustracción', 3),
    (8,  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Multiplicación y división', 'Operaciones de multiplicación, división y tablas', 3),
    (9,  'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Figuras planas', 'Triángulos, cuadriláteros, círculos y propiedades', 3),
    (10, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Cuerpos geométricos', 'Prismas, pirámides, cilindros y elementos', 3),
    (11, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Lectura literal', 'Identificación de información explícita', 3),
    (12, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Lectura inferencial', 'Deducción de información implícita', 3),
    (13, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 6, 'Texto narrativo', 'Escritura de cuentos, relatos y narraciones', 3),
    (14, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 6, 'Texto informativo', 'Escritura de informes, noticias y descripciones', 3)
ON CONFLICT (id) DO NOTHING;

-- Courses (year removed in migration 000013; academic year now lives in course_subjects.school_year)
INSERT INTO courses (id, organization_id, name) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2do 1era'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2do 2da'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '3ro 1era')
ON CONFLICT (id) DO NOTHING;

-- Students
INSERT INTO students (id, course_id, name) VALUES
    (1, 1, 'Lucía Martinez'),
    (2, 1, 'Tomás Rodriguez'),
    (3, 1, 'Valentina Lopez'),
    (4, 2, 'Santiago Diaz'),
    (5, 2, 'Martina Ruiz'),
    (6, 3, 'Emilia Castro'),
    (7, 3, 'Joaquin Vargas')
ON CONFLICT (id) DO NOTHING;

-- Course-Subject assignments (teacher_id references users: 3=María Docente, 4=Pedro Multirol)
INSERT INTO course_subjects (id, organization_id, course_id, subject_id, teacher_id, school_year, start_date, end_date) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 1, 3, 2026, '2026-03-02', '2026-12-15'),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 3, 4, 2026, '2026-03-02', '2026-12-15'),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 1, 3, 2026, '2026-03-02', '2026-12-15'),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 2, 4, 2026, '2026-03-02', '2026-12-15')
ON CONFLICT (id) DO NOTHING;

-- Activity templates (didactic activities by class moment)
INSERT INTO activities (id, organization_id, moment, name, description, duration_minutes) VALUES
    (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'apertura',   'Lluvia de ideas',          'Técnica grupal de generación de ideas', 15),
    (2, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'apertura',   'Pregunta disparadora',     'Pregunta para activar conocimientos previos', 10),
    (3, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'apertura',   'Repaso clase anterior',    'Revisión breve de contenidos previos', 10),
    (4, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'desarrollo', 'Trabajo en grupo',         'Actividad colaborativa en equipos', 30),
    (5, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'desarrollo', 'Exposición dialogada',     'Presentación con intercambio participativo', 25),
    (6, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'desarrollo', 'Resolución de problemas',  'Aplicación de conceptos a situaciones concretas', 30),
    (7, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'desarrollo', 'Análisis de caso',         'Estudio y discusión de un caso real o simulado', 25),
    (8, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'cierre',     'Puesta en común',          'Socialización de resultados y conclusiones', 15),
    (9, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'cierre',     'Síntesis individual',      'Cada alumno resume lo aprendido', 10),
    (10, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'cierre',    'Autoevaluación',           'Reflexión personal sobre el aprendizaje', 10)
ON CONFLICT (id) DO NOTHING;

-- Time slots: weekly schedule grid. Demo: 1 normal class + 1 shared class on course 1 (2do 1era).
-- Slot 1 (lunes 8:00-9:00) = Matemática only.
-- Slot 2 (martes 10:00-11:00) = Matemática + Historia dictadas en simultáneo (clase compartida, HU-3.5).
INSERT INTO time_slots (id, course_id, day_of_week, start_time, end_time) VALUES
    (1, 1, 1, '08:00', '09:00'),
    (2, 1, 2, '10:00', '11:00')
ON CONFLICT (id) DO NOTHING;

-- 1 row per slot = clase normal. 2 rows con el mismo time_slot_id = clase compartida.
INSERT INTO time_slot_subjects (id, time_slot_id, course_subject_id) VALUES
    (1, 1, 1),
    (2, 2, 1),
    (3, 2, 2)
ON CONFLICT (id) DO NOTHING;

-- Reset sequences to avoid conflicts with future inserts
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 0) FROM users));
SELECT setval('user_roles_id_seq', (SELECT COALESCE(MAX(id), 0) FROM user_roles));
SELECT setval('areas_id_seq', (SELECT COALESCE(MAX(id), 0) FROM areas));
SELECT setval('subjects_id_seq', (SELECT COALESCE(MAX(id), 0) FROM subjects));
SELECT setval('area_coordinators_id_seq', (SELECT COALESCE(MAX(id), 0) FROM area_coordinators));
SELECT setval('topics_id_seq', (SELECT COALESCE(MAX(id), 0) FROM topics));
SELECT setval('courses_id_seq', (SELECT COALESCE(MAX(id), 0) FROM courses));
SELECT setval('students_id_seq', (SELECT COALESCE(MAX(id), 0) FROM students));
SELECT setval('course_subjects_id_seq', (SELECT COALESCE(MAX(id), 0) FROM course_subjects));
SELECT setval('activities_id_seq', (SELECT COALESCE(MAX(id), 0) FROM activities));
SELECT setval('time_slots_id_seq', (SELECT COALESCE(MAX(id), 0) FROM time_slots));
SELECT setval('time_slot_subjects_id_seq', (SELECT COALESCE(MAX(id), 0) FROM time_slot_subjects));
