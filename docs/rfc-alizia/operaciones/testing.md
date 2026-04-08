# QA — Estrategia de testing

## Precondiciones

- PostgreSQL corriendo con schema migrado (4 migraciones: init, coordination, teaching, resources)
- Organización seed con config de ejemplo:
  - `topic_max_levels: 3`
  - `topic_level_names: ["Núcleos", "Áreas de conocimiento", "Categorías"]`
  - `topic_selection_level: 3`
  - `shared_classes_enabled: true`
  - `desarrollo_max_activities: 3`
  - `coord_doc_sections` con 3 secciones (problem_edge, methodological_strategy, eval_criteria)
- Usuarios seed: 1 admin, 1 coordinator, 2 teachers
- Topics seed: jerarquía de 3 niveles (2 núcleos → 4 áreas → 8 categorías)
- Areas seed: 1 área con 2 subjects
- Courses seed: 1 curso con schedule (incluyendo 1 clase compartida)
- Activities seed: 2 por momento (6 total)

---

## Matriz de testing por fase

### Fase 1: Setup

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 1.1 | GET /health | 200 `{"status": "ok"}` | Alta |
| 1.2 | Request sin Authorization header a ruta protegida | 401 `missing_token` | Alta |
| 1.3 | Request con JWT inválido | 401 `invalid_token` | Alta |
| 1.4 | Request con JWT válido | 200 + claims en context | Alta |
| 1.5 | Request con JWT de otra org | Datos filtrados por org_id del token | Alta |

### Fase 2: Admin/Integration

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 2.1 | Crear área | 201 + área creada con organization_id del JWT | Alta |
| 2.2 | Crear área sin role coordinator/admin | 403 `forbidden` | Alta |
| 2.3 | Listar áreas (filtra por org) | Solo áreas de la org del token | Alta |
| 2.4 | Crear subject en área | 201 + subject vinculado al área | Alta |
| 2.5 | Crear topic nivel 1 (parent_id=NULL) | 201 + level=1 | Alta |
| 2.6 | Crear topic nivel 2 (parent_id=topic_level_1) | 201 + level=2 | Alta |
| 2.7 | Crear topic nivel 4 cuando max_levels=3 | 400 `topic exceeds max level` | Alta |
| 2.8 | Crear curso + students | 201 | Media |
| 2.9 | Crear course_subject (curso + disciplina + docente) | 201 | Alta |
| 2.10 | Crear time_slot | 201 | Media |
| 2.11 | Crear time_slot_subject | 201 | Alta |
| 2.12 | Crear 2 time_slot_subjects en mismo slot (clase compartida) | 201 si shared_classes_enabled | Alta |
| 2.13 | Crear 2 time_slot_subjects cuando shared_classes_enabled=false | 400 | Alta |
| 2.14 | Trigger: course_subject de otro curso en time_slot | Error del trigger | Alta |
| 2.15 | Crear activities por momento | 201 | Media |

### Fase 3: Coordination Documents

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 3.1 | Crear documento (wizard paso 1: topics) | 201 + doc en draft + coord_doc_topics creados | Alta |
| 3.2 | Asignar disciplinas + class_count al doc | coordination_document_subjects creados | Alta |
| 3.3 | Asignar topics a cada disciplina | coord_doc_subject_topics creados | Alta |
| 3.4 | Validar que todos los topics del doc estén distribuidos | Error si algún topic no asignado a ninguna disciplina | Alta |
| 3.5 | PATCH sections (actualizar sección dinámica) | sections JSONB actualizado | Alta |
| 3.6 | PATCH sections con key inexistente en config | 400 `invalid section key` | Media |
| 3.7 | Publicar documento (in_progress → published) | Status actualizado | Alta |
| 3.8 | Publicar documento ya publicado | 400 o idempotente | Media |
| 3.9 | Comenzar edición (pending → in_progress) | Status actualizado | Media |
| 3.10 | DELETE documento en pending | 200 + eliminado | Media |
| 3.11 | DELETE documento publicado | 400 `cannot delete published document` | Alta |
| 3.12 | GET documento completo (con todas las junction tables) | Todas las relaciones cargadas | Alta |
| 3.13 | Docente ve doc publicado | 200 (lectura) | Alta |
| 3.14 | Docente intenta editar doc (si config restringe) | 403 | Media |
| 3.15 | Listar documentos filtrados por area_id | Solo docs del área solicitada | Media |

### Fase 4: AI Generation

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 4.1 | POST /coordination-documents/:id/generate | Secciones populadas según ai_prompt de config | Alta |
| 4.2 | Generación crea plan de clases por disciplina | coord_doc_classes creadas con class_number, title, objective | Alta |
| 4.3 | Generación asigna topics a cada clase | coord_doc_class_topics creados | Alta |
| 4.4 | Chat: "cambiá el eje problemático por algo más corto" | update_section ejecutado, sección actualizada | Alta |
| 4.5 | Chat: update_section con key inválida | Error de validación | Media |
| 4.6 | Chat: update_class modifica título y objetivo | Clase actualizada | Media |
| 4.7 | Chat: update_class_topics cambia topics de una clase | coord_doc_class_topics actualizados | Media |
| 4.8 | Generación con sección tipo select_text | selected_option respetado en el prompt | Media |

### Fase 5: Teaching

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 5.1 | Crear lesson plan heredando de doc publicado | 201 + class_number, title, objective del doc | Alta |
| 5.2 | Crear lesson plan sin doc publicado | 400 `no published coordination document` | Alta |
| 5.3 | Seleccionar 1 actividad de apertura | moments.apertura.activities = [id] | Alta |
| 5.4 | Seleccionar 0 actividades de apertura | 400 `apertura requires exactly 1 activity` | Alta |
| 5.5 | Seleccionar 4 actividades de desarrollo (max=3) | 400 `max 3 activities in desarrollo` | Alta |
| 5.6 | Seleccionar topics (subconjunto de los del doc) | lesson_plan_topics creados | Alta |
| 5.7 | Seleccionar fonts modo global | lesson_plan_moment_fonts con moment=NULL | Media |
| 5.8 | Seleccionar fonts por momento | lesson_plan_moment_fonts con moment=apertura/desarrollo/cierre | Media |
| 5.9 | Generar contenido por actividad | activityContent populado en moments JSONB | Alta |
| 5.10 | Status cambia a published | Status actualizado | Media |
| 5.11 | Clase compartida muestra indicador | Respuesta incluye flag de shared class | Media |

### Fase 6: Resources

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 6.1 | GET /resource-types (org con overrides) | Tipos públicos habilitados + privados de la org | Alta |
| 6.2 | Tipo público deshabilitado por org | No aparece en la lista | Alta |
| 6.3 | Tipo con custom_prompt | Prompt override usado en generación | Alta |
| 6.4 | Tipo con custom_output_schema | Schema override usado | Media |
| 6.5 | Crear recurso con tipo que requires_font sin font | 400 `font required` | Alta |
| 6.6 | Crear recurso con font | 201 + font_id guardado | Alta |
| 6.7 | Generar recurso con IA | content JSONB populado según output_schema | Alta |
| 6.8 | GET /fonts filtrado por area | Solo fonts del área | Media |
| 6.9 | Font con is_validated=false | No visible para docentes en la API | Media |

---

## Testing de permisos y roles

| # | Escenario | Rol | Acción | Resultado |
|---|-----------|-----|--------|-----------|
| P1 | Coordinator crea documento | coordinator | POST /coordination-documents | 201 |
| P2 | Teacher no puede crear documento | teacher | POST /coordination-documents | 403 |
| P3 | Teacher crea lesson plan | teacher | POST /lesson-plans | 201 |
| P4 | Coordinator no crea lesson plan | coordinator (sin role teacher) | POST /lesson-plans | 403 |
| P5 | Admin crea área | admin | POST /areas | 201 |
| P6 | Teacher no crea área | teacher | POST /areas | 403 |
| P7 | Multi-rol: coordinator+teacher crea doc Y plan | coordinator+teacher | POST ambos | 201 ambos |
| P8 | User de org A no ve datos de org B | user org A | GET /areas (org B) | [] vacío |

---

## Testing de multi-tenancy

| # | Escenario | Resultado |
|---|-----------|-----------|
| T1 | Crear área en org A, listar desde org B | Org B no ve el área |
| T2 | Crear topic en org A, buscar desde org B | No encontrado |
| T3 | Crear documento en org A, GET desde org B | 404 |
| T4 | Dos orgs con mismo nombre de área | Ambas coexisten, IDs distintos |

---

## Fixtures y datos de test

### Seed SQL para tests de integracion

Los tests de integracion corren contra PostgreSQL real. Se usa un seed SQL que crea un dataset minimo para cubrir todos los flujos.

**Archivo:** `db/seeds/test_seed.sql`

```sql
-- Organizacion de test
INSERT INTO organizations (id, name, slug, config) VALUES (1, 'Test University', 'test-uni', '{
  "topic_max_levels": 3,
  "topic_level_names": ["Nucleos", "Areas de conocimiento", "Categorias"],
  "topic_selection_level": 3,
  "shared_classes_enabled": true,
  "desarrollo_max_activities": 3,
  "coord_doc_sections": [
    {"key": "problem_edge", "label": "Eje problematico", "type": "text", "ai_prompt": "...", "required": true},
    {"key": "methodological_strategy", "label": "Estrategia metodologica", "type": "select_text", "options": ["proyecto", "taller_laboratorio", "ateneo_debate"], "ai_prompt": "...", "required": true},
    {"key": "eval_criteria", "label": "Criterios de evaluacion", "type": "text", "ai_prompt": "...", "required": false}
  ]
}');

-- Segunda org para tests de multi-tenancy
INSERT INTO organizations (id, name, slug, config) VALUES (2, 'Other School', 'other-school', '{
  "topic_max_levels": 2,
  "topic_level_names": ["Unidades", "Temas"],
  "topic_selection_level": 2,
  "shared_classes_enabled": false,
  "desarrollo_max_activities": 2,
  "coord_doc_sections": [
    {"key": "problem_edge", "label": "Eje", "type": "text", "ai_prompt": "...", "required": true}
  ]
}');

-- Usuarios
INSERT INTO users (id, organization_id, email, name, password_hash) VALUES
  (1, 1, 'admin@test.com', 'Admin User', '$2a$12$...'),       -- admin
  (2, 1, 'coord@test.com', 'Coordinator User', '$2a$12$...'),  -- coordinator
  (3, 1, 'teacher1@test.com', 'Teacher One', '$2a$12$...'),    -- teacher
  (4, 1, 'teacher2@test.com', 'Teacher Two', '$2a$12$...'),    -- teacher
  (5, 2, 'other@other.com', 'Other Org User', '$2a$12$...');   -- admin org 2

INSERT INTO user_roles (user_id, role) VALUES
  (1, 'admin'), (2, 'coordinator'), (3, 'teacher'), (4, 'teacher'), (5, 'admin');

-- Areas
INSERT INTO areas (id, organization_id, name, description) VALUES
  (1, 1, 'Ciencias Exactas', 'Area de ciencias exactas y naturales'),
  (2, 2, 'Ciencias', 'Area unica');

INSERT INTO area_coordinators (area_id, user_id) VALUES (1, 2);

-- Subjects
INSERT INTO subjects (id, organization_id, area_id, name) VALUES
  (1, 1, 1, 'Matematicas'),
  (2, 1, 1, 'Fisica');

-- Topics (3 niveles)
INSERT INTO topics (id, organization_id, parent_id, name, level) VALUES
  (1, 1, NULL, 'Pensamiento logico-matematico', 1),     -- Nucleo
  (2, 1, 1,    'Aritmetica', 2),                         -- Area de conocimiento
  (3, 1, 1,    'Algebra', 2),
  (4, 1, 2,    'Suma y resta', 3),                       -- Categoria
  (5, 1, 2,    'Multiplicacion', 3),
  (6, 1, 3,    'Ecuaciones lineales', 3),
  (7, 1, NULL, 'Pensamiento cientifico', 1),
  (8, 1, 7,    'Mecanica', 2),
  (9, 1, 8,    'Cinematica', 3),
  (10, 1, 8,   'Dinamica', 3);

-- Courses
INSERT INTO courses (id, organization_id, name) VALUES (1, 1, '3ro 1era');

-- Students
INSERT INTO students (id, course_id, name) VALUES
  (1, 1, 'Juan Perez'), (2, 1, 'Maria Garcia'), (3, 1, 'Pedro Lopez');

-- Course Subjects
INSERT INTO course_subjects (id, course_id, subject_id, teacher_id, start_date, end_date, school_year) VALUES
  (1, 1, 1, 3, '2026-03-01', '2026-11-30', 2026),  -- Matematicas con Teacher One
  (2, 1, 2, 4, '2026-03-01', '2026-11-30', 2026);  -- Fisica con Teacher Two

-- Time Slots (incluyendo clase compartida)
INSERT INTO time_slots (id, course_id, day_of_week, start_time, end_time) VALUES
  (1, 1, 1, '08:00', '09:30'),  -- Lunes 08:00 (Matematicas)
  (2, 1, 2, '09:45', '11:15'),  -- Martes 09:45 (Fisica)
  (3, 1, 3, '08:00', '09:30'),  -- Miercoles 08:00 (compartida)
  (4, 1, 4, '08:00', '09:30'),  -- Jueves 08:00 (Matematicas)
  (5, 1, 5, '09:45', '11:15');  -- Viernes 09:45 (Matematicas)

INSERT INTO time_slot_subjects (time_slot_id, course_subject_id) VALUES
  (1, 1),    -- Lunes: Matematicas
  (2, 2),    -- Martes: Fisica
  (3, 1),    -- Miercoles: Matematicas (compartida)
  (3, 2),    -- Miercoles: Fisica (compartida)
  (4, 1),    -- Jueves: Matematicas
  (5, 1);    -- Viernes: Matematicas

-- Activities (2 por momento)
INSERT INTO activities (id, organization_id, moment, name, description, duration_minutes) VALUES
  (1, 1, 'apertura', 'Lluvia de ideas', 'Recoleccion de ideas previas', 15),
  (2, 1, 'apertura', 'Pregunta disparadora', 'Pregunta para activar conocimiento previo', 10),
  (3, 1, 'desarrollo', 'Trabajo en grupo', 'Actividad colaborativa', 30),
  (4, 1, 'desarrollo', 'Resolucion de problemas', 'Ejercicios individuales', 25),
  (5, 1, 'desarrollo', 'Exposicion dialogada', 'Explicacion con participacion', 20),
  (6, 1, 'cierre', 'Puesta en comun', 'Socializacion de resultados', 15),
  (7, 1, 'cierre', 'Autoevaluacion', 'Reflexion individual', 10);

-- Fonts
INSERT INTO fonts (id, organization_id, name, description, file_url, file_type, area_id, is_validated) VALUES
  (1, 1, 'Guia de aritmetica', 'PDF con ejercicios de suma y resta', 'https://example.com/aritmetica.pdf', 'pdf', 1, true),
  (2, 1, 'Video cinematica', 'Video tutorial de movimiento', 'https://example.com/cinematica.mp4', 'video', 1, true),
  (3, 1, 'Borrador no validado', 'Font sin validar', 'https://example.com/draft.pdf', 'pdf', 1, false);

-- Resource Types (1 publico, 1 privado)
INSERT INTO resource_types (id, key, name, description, prompt, output_schema, organization_id, requires_font) VALUES
  (1, 'lecture_guide', 'Guia de lectura', 'Guia para acompanar lectura', 'Genera una guia...', '{"type":"object","properties":{"title":{"type":"string"},"sections":{"type":"array"}}}', NULL, true),
  (2, 'course_sheet', 'Ficha de curso', 'Ficha resumen del curso', 'Genera una ficha...', '{"type":"object","properties":{"title":{"type":"string"},"objectives":{"type":"array"}}}', NULL, false);

-- Reset sequences
SELECT setval('organizations_id_seq', 10);
SELECT setval('users_id_seq', 10);
SELECT setval('areas_id_seq', 10);
SELECT setval('subjects_id_seq', 10);
SELECT setval('topics_id_seq', 20);
SELECT setval('courses_id_seq', 10);
SELECT setval('students_id_seq', 10);
SELECT setval('course_subjects_id_seq', 10);
SELECT setval('time_slots_id_seq', 10);
SELECT setval('activities_id_seq', 10);
SELECT setval('fonts_id_seq', 10);
SELECT setval('resource_types_id_seq', 10);
```

### Helper de test en Go

```go
// db/testhelper/setup.go
func SetupTestDB(t *testing.T) *gorm.DB {
    db := dbconn.NewPostgresConnector(testConfig).Connect()
    // Ejecutar migraciones
    RunMigrations(db)
    // Ejecutar seed
    RunSeedSQL(db, "db/seeds/test_seed.sql")
    // Cleanup al terminar
    t.Cleanup(func() { TruncateAllTables(db) })
    return db
}
```

### Mock de Azure OpenAI

Para tests que involucran generacion IA, usar un mock del provider:

```go
// src/mocks/providers/ai_mock.go
type MockAIProvider struct {
    GenerateSectionFunc  func(ctx context.Context, prompt string) (string, error)
    GenerateClassPlanFunc func(ctx context.Context, params ClassPlanParams) ([]ClassPlan, error)
    ChatFunc             func(ctx context.Context, messages []Message, tools []Tool) (*ChatResponse, error)
}
```

Los tests de integracion de IA usan este mock. Los tests E2E contra Azure OpenAI real se corren manualmente (no en CI) marcados con `//go:build e2e`.

### JWT de test

`team-ai-toolkit/tokens` en environment `Test` auto-genera claims. Para tests manuales:

```go
// Generar JWT de test para cada rol
adminToken := testhelper.TokenForUser(1, "admin")        // Admin org 1
coordToken := testhelper.TokenForUser(2, "coordinator")  // Coordinator org 1
teacherToken := testhelper.TokenForUser(3, "teacher")    // Teacher org 1
otherOrgToken := testhelper.TokenForUser(5, "admin")     // Admin org 2 (multi-tenancy)
```

---

## Coverage target: 80%

Reportado automaticamente en PRs via GitHub Actions.

### Que cubrir por capa

| Capa | Tipo de test | Cobertura esperada |
|------|-------------|-------------------|
| Usecases | Unit tests (mocks de providers) | 90%+ |
| Repositories | Integration tests (PostgreSQL real) | 80%+ |
| Handlers | Integration tests (HTTP + mocks) | 70%+ |
| Triggers | SQL tests (INSERT + expect error) | 100% |

### Que NO testear

- Codigo generado por GORM (AutoMigrate, etc.)
- Libreria team-ai-toolkit (tiene sus propios tests)
- Azure OpenAI responses (mock en CI, E2E manual)
