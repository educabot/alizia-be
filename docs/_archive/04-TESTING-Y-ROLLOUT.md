# Testing, Rollout, Riesgos y Dependencias — Alizia v2

## QA — Estrategia de testing

### Precondiciones

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

### Matriz de testing por fase

#### Fase 1: Setup

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 1.1 | GET /health | 200 `{"status": "ok"}` | Alta |
| 1.2 | Request sin Authorization header a ruta protegida | 401 `missing_token` | Alta |
| 1.3 | Request con JWT inválido | 401 `invalid_token` | Alta |
| 1.4 | Request con JWT válido | 200 + claims en context | Alta |
| 1.5 | Request con JWT de otra org | Datos filtrados por org_id del token | Alta |

#### Fase 2: Admin/Integration

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

#### Fase 3: Coordination Documents

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

#### Fase 4: AI Generation

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

#### Fase 5: Teaching

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

#### Fase 6: Resources

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

### Testing de permisos y roles

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

### Testing de multi-tenancy

| # | Escenario | Resultado |
|---|-----------|-----------|
| T1 | Crear área en org A, listar desde org B | Org B no ve el área |
| T2 | Crear topic en org A, buscar desde org B | No encontrado |
| T3 | Crear documento en org A, GET desde org B | 404 |
| T4 | Dos orgs con mismo nombre de área | Ambas coexisten, IDs distintos |

### Coverage target: 80%

Reportado automáticamente en PRs via GitHub Actions.

---

## Rollout 📈

| Fase | Alcance | Criterio para avanzar | Duración estimada |
|------|---------|----------------------|-------------------|
| 1 | Staging — equipo interno | CI verde, /health responde, auth funciona, CRUD básico | 1 semana |
| 2 | Org piloto (1 provincia) | Coordinador crea doc + genera con IA + publica. Docente planifica 1 clase | 2 semanas |
| 3 | 2-3 orgs adicionales | Feedback positivo, sin bugs bloqueantes, config por org funciona | 2 semanas |
| 4 | Todas las orgs | Métricas de éxito alcanzadas, 80% coverage, docs API actualizados | Continuo |

### Plan de rollback

1. **Railway**: revertir al deploy anterior (1 click en dashboard)
2. **Migraciones**: ejecutar `.down.sql` correspondiente
3. **Config**: revertir JSON de organización a versión anterior
4. **Comunicar** al equipo y documentar causa del rollback

### Monitoreo post-deploy

| Qué monitorear | Herramienta | Umbral de alerta |
|---|---|---|
| Errores 5xx | Bugsnag | > 5 en 5 minutos |
| Latencia de endpoints | Railway logs (slog) | p95 > 2 segundos |
| Healthcheck | Railway built-in | /health no responde en 30s |
| Errores de IA | Bugsnag + slog | Azure OpenAI timeout > 30s |
| Login failures | slog (login_failed events) | > 20 en 5 minutos desde misma IP |

---

## Dependencias 👥

### Internas

| Dependencia | Tipo | Bloqueante | Estado | Notas |
|-------------|------|------------|--------|-------|
| team-ai-toolkit | Librería Go | Sí (Fase 1) | ✅ Creado | Repo con tests, compila limpio |
| auth-service | Microservicio | No (futuro) | ⬜ Futuro (no bloqueante) | Planificado para reemplazar Auth0. Alizia v2 arranca con Auth0 |
| Auth0 tenant config | Infra | Sí (Fase 1) | ⬜ Configurar | Domain + audience + API (mismo sistema que tich-cronos) |
| Railway account | Infra | Sí (Fase 1) | ⬜ Configurar | Cuenta + proyecto + PostgreSQL |
| PostgreSQL en Railway | Infra | Sí (Fase 1) | ⬜ Provisionar | O DB externa |
| Azure OpenAI access | Servicio | Sí (Fase 4) | ✅ Ya disponible | Mismo acceso que el POC |
| Diseño UX/UI | Entregable | No (backend first) | ⬜ En progreso | Frontend es RFC separado |

### Externas

| Dependencia | Tipo | Bloqueante | Contacto |
|-------------|------|------------|----------|
| Azure OpenAI | API LLM | Sí (Fase 4) | Ya configurado |
| SendGrid | Email (auth-service) | No (noop funciona) | Cuenta por configurar |
| Equipo pedagógico provincial | Contenido | No (usamos defaults) | Reuniones pendientes |

---

## Riesgos ⚠️

| # | Riesgo | Probabilidad | Impacto | Mitigación |
|---|--------|-------------|---------|------------|
| 1 | GORM genera queries N+1 en documents con 8+ JOINs | Media | Medio | Preload explícito en GORM + `db.Raw()` para queries complejas. Documentado en arquitectura como patrón. Si >50% de repos usan Raw, evaluar migración a sqlx |
| 2 | Config JSONB por org se vuelve inmanejable con muchos campos | Baja | Alto | Validación estricta en backend, schema documentado, defaults sensatos. No permitir campos arbitrarios |
| 3 | IA genera contenido de baja calidad o desalineado | Media | Medio | Prompts iterativos (empezar simple, mejorar con uso real). Review humano obligatorio antes de publicar (status draft). Prompts configurables por provincia |
| 4 | Multi-tenancy data leak (org A ve datos de org B) | Baja | Crítico | Middleware de tenant en TODAS las rutas. Tests específicos de isolation (ver matriz T1-T4). org_id viene del JWT, no del request |
| 5 | Railway downtime afecta servicio | Baja | Medio | Dockerfile portable. Si Railway cae, migrar a Render/Fly.io/VPS en horas. Zero vendor lock-in |
| 6 | CTE recursivo de topics lento con muchos niveles | Baja | Medio | Level precalculado en tabla. Solo recalcular rama afectada al mover topic. Max 5 niveles en la práctica |
| 7 | Equipo pedagógico no define tipos de recurso a tiempo | Media | Bajo | Arrancar con 2 tipos genéricos (lecture_guide, course_sheet). Agregar más post-MVP |
| 8 | Frontend y backend desalineados en formato de JSONB | Media | Medio | Swagger/OpenAPI como contrato. Validar schemas en CI |
| 9 | Clases compartidas generan edge cases no contemplados | Media | Medio | Feature flag deshabilitado por defecto. Activar solo en orgs que lo necesiten y testear exhaustivamente |

---

## Preguntas abiertas ❓

| # | Pregunta | Área | Estado |
|---|----------|------|--------|
| 1 | ¿Cómo se cargan los datos iniciales de una provincia? (manual, CSV, API) | Producto/Ops | 🟡 Pendiente |
| 2 | ¿Quién crea las organizaciones? (super admin de Educabot o self-service) | Producto | 🟡 Pendiente |
| 3 | ¿Los docentes pueden ver docs de otras áreas? (permiso configurable por org) | Producto | 🟡 Pendiente |
| 4 | ¿Se mantiene historial de versiones de coordination documents? | Producto | 🟡 Pendiente |
| 5 | ¿Cuántos concurrent users se esperan por org? (para dimensionar Railway) | Infra | 🟡 Pendiente |
| 6 | ¿Qué pasa si el docente no tiene internet al grabar bitácora? (escuelas rurales) | Producto | 🟡 Pendiente |
| 7 | ¿Se permite subida de fuentes propias del docente? (decisión por provincia) | Producto | 🟡 Pendiente |
| 8 | ¿Cómo se exporta el lesson plan a PDF? (template configurable, formato) | Producto/Front | 🟡 Pendiente |
| 9 | ¿El dashboard necesita notificaciones push o solo polling? | Frontend | 🟡 Pendiente |
| 10 | ¿Rate limiting en endpoints de generación IA? (costo por request) | Backend | 🟡 Pendiente |
