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
| 2.9 | Crear course_subject (curso + materia + docente) | 201 | Alta |
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
| 3.2 | Asignar materias + class_count al doc | coordination_document_subjects creados | Alta |
| 3.3 | Asignar topics a cada materia | coord_doc_subject_topics creados | Alta |
| 3.4 | Validar que todos los topics del doc estén distribuidos | Error si algún topic no asignado a ninguna materia | Alta |
| 3.5 | PATCH sections (actualizar sección dinámica) | sections JSONB actualizado | Alta |
| 3.6 | PATCH sections con key inexistente en config | 400 `invalid section key` | Media |
| 3.7 | Publicar documento (draft → published) | Status actualizado | Alta |
| 3.8 | Publicar documento ya publicado | 400 o idempotente | Media |
| 3.9 | Archivar documento publicado | Status → archived | Media |
| 3.10 | DELETE documento en draft | 200 + eliminado | Media |
| 3.11 | DELETE documento publicado | 400 `cannot delete published document` | Alta |
| 3.12 | GET documento completo (con todas las junction tables) | Todas las relaciones cargadas | Alta |
| 3.13 | Docente ve doc publicado | 200 (lectura) | Alta |
| 3.14 | Docente intenta editar doc (si config restringe) | 403 | Media |
| 3.15 | Listar documentos filtrados por area_id | Solo docs del área solicitada | Media |

### Fase 4: AI Generation

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 4.1 | POST /coordination-documents/:id/generate | Secciones populadas según ai_prompt de config | Alta |
| 4.2 | Generación crea plan de clases por materia | coord_doc_classes creadas con class_number, title, objective | Alta |
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
| 5.10 | Status cambia a planned | Status actualizado | Media |
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

## Coverage target: 80%

Reportado automáticamente en PRs via GitHub Actions.
