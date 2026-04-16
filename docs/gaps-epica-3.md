# Gaps de implementación — Épica 3 (Integración)

> Endpoints documentados en el RFC de Épica 3 (y onboarding relacionado) que **todavía no están implementados** en el backend. Surgió del análisis FE↔BE al intentar levantar el front en local y ver pantallas vacías.

**Fecha del análisis:** 2026-04-16
**Branch analizada:** `master`
**Fuente de verdad:** `docs/rfc-alizia/tecnico/endpoints.md` + `docs/rfc-alizia/epicas/03-integracion/*` + `docs/rfc-alizia/epicas/02-onboarding/*`

---

## TL;DR — qué bloquea al frontend

1. **`GET /subjects`** con `?area_id?` — el FE lo llama en bootstrap para poblar selectores de disciplinas; hoy tira 404 y deja todo vacío.
2. **`GET /course-subjects`** con `?course_id?&teacher_id?` — necesario para listar asignaciones docente/curso/disciplina en múltiples pantallas.
3. **`GET /topics?parent_id=`** — documentado en RFC, actualmente sólo parsea `level`.
4. **`PUT /areas/:id`** — criterio de aceptación de HU-3.2.
5. **`POST /courses` contrato** — alinear con RFC (acepta sólo `{ name }`; hoy el entity lleva `year`).

Todo lo demás está implementado.

---

## Estado actual del BE

Extraído de `src/app/web/mapping.go`:

### Epica 3 — rutas existentes

| Método | Ruta | Handler | Roles |
|---|---|---|---|
| GET | `/organizations/me` | Admin.HandleGetOrganization | any |
| PATCH | `/organizations/me/config` | Admin.HandleUpdateOrgConfig | admin |
| GET | `/areas` | Admin.HandleListAreas | any |
| POST | `/areas` | Admin.HandleCreateArea | coord, admin |
| GET | `/areas/:id/subjects` | Admin.HandleListSubjects | any |
| POST | `/areas/:id/coordinators` | Admin.HandleAssignCoordinator | admin |
| DELETE | `/areas/:id/coordinators/:user_id` | Admin.HandleRemoveCoordinator | admin |
| POST | `/subjects` | Admin.HandleCreateSubject | coord, admin |
| GET | `/topics` | Admin.HandleGetTopics | any |
| POST | `/topics` | Admin.HandleCreateTopic | coord, admin |
| PATCH | `/topics/:id` | Admin.HandleUpdateTopic | coord, admin |
| GET | `/activities` | Admin.HandleListActivities | any |
| POST | `/activities` | Admin.HandleCreateActivity | admin |
| GET | `/courses` | Courses.HandleListCourses | any |
| GET | `/courses/:id` | Courses.HandleGetCourse | any |
| POST | `/courses` | Courses.HandleCreateCourse | admin |
| POST | `/courses/:id/students` | Courses.HandleAddStudent | admin |
| POST | `/courses/:id/time-slots` | Courses.HandleCreateTimeSlot | admin |
| GET | `/courses/:id/schedule` | Courses.HandleGetSchedule | any |
| POST | `/course-subjects` | Courses.HandleAssignCourseSubject | admin |
| GET | `/course-subjects/:id/shared-class-numbers` | Courses.HandleGetSharedClassNumbers | any |

### Onboarding (Epica 2) — rutas existentes

| Método | Ruta | Handler |
|---|---|---|
| GET | `/users/me/onboarding-status` | Onboarding.HandleGetStatus |
| POST | `/users/me/onboarding/complete` | Onboarding.HandleComplete |
| GET | `/users/me/profile` | Onboarding.HandleGetProfile |
| PUT | `/users/me/profile` | Onboarding.HandleSaveProfile |
| GET | `/users/me/onboarding/tour-steps` | Onboarding.HandleGetTourSteps |
| GET | `/onboarding-config` | Onboarding.HandleGetConfig |

Onboarding está completo.

---

## Gaps (Épica 3)

### G-1 · `GET /subjects` — listado flat con filtro opcional por área

**Referencia RFC:** `tecnico/endpoints.md` §164 + HU-3.2 criterio "CRUD endpoints para disciplinas: POST, GET (listar)".

**Contrato esperado:**
```
GET /api/v1/subjects?area_id={id}&limit={n}&offset={n}
```
- `area_id` opcional — sin él, devuelve todas las disciplinas de la org.
- Response paginado `{ items: Subject[], more: bool }`.

**Por qué importa:** hoy el FE hace `subjectsApi.list()` en bootstrap (sin area_id) para poblar un cache global. Sin este endpoint el FE no puede listar disciplinas transversalmente.

**Archivos a tocar:**
- `src/core/usecases/admin/list_subjects.go` — relajar `Validate()` para permitir `AreaID == 0` (hoy lo rechaza como validation error). Considerar separar en dos casos de uso (`ListSubjects` para nested, `ListAllSubjects` para flat) o usar un sentinel.
- `src/entrypoints/admin.go` — nuevo handler `HandleListAllSubjects` que lee `area_id` de query en vez de path param.
- `src/app/web/mapping.go` — registrar `api.GET("/subjects", ...)`.
- `src/repositories/subject.go` — método que filtre por `organization_id` y opcionalmente `area_id`.

**Tests:**
- Listar sin filtro retorna todas las disciplinas de la org.
- Listar con `area_id` filtra correctamente.
- Multi-tenancy: org A no ve subjects de org B.

---

### G-2 · `GET /course-subjects` — listado con filtros

**Referencia RFC:** `tecnico/endpoints.md` §443.

**Contrato esperado:**
```
GET /api/v1/course-subjects?course_id={id}&subject_id={id}&teacher_id={id}&limit={n}&offset={n}
```
- Todos los filtros opcionales.
- Response paginado con el mismo schema que `POST /course-subjects`.

**Por qué importa:** el FE lo usa en pantallas de docentes (para listar sus asignaciones) y en el wizard de creación de planes de clase (para elegir un curso-disciplina).

**Archivos a tocar:**
- `src/core/usecases/admin/list_course_subjects.go` — nuevo caso de uso.
- `src/entrypoints/courses.go` — nuevo `HandleListCourseSubjects` que parsee query params.
- `src/app/web/mapping.go` — `api.GET("/course-subjects", ...)`.
- `src/repositories/course_subject.go` — método con WHERE dinámico.

**Tests:**
- Filtrar por `teacher_id` devuelve sólo asignaciones de ese docente.
- Filtrar por `course_id` devuelve todas las disciplinas de ese curso.
- Sin filtros devuelve todas las asignaciones de la org.

---

### G-3 · `GET /topics?parent_id=` — filtro por padre

**Referencia RFC:** `tecnico/endpoints.md` §378 — "Con `parent_id=N` devuelve hijos directos".

**Estado actual:** `HandleGetTopics` en `src/entrypoints/admin.go:150` sólo parsea `level`. El query `parent_id` se ignora silenciosamente.

**Contrato esperado:**
```
GET /api/v1/topics?parent_id={id}
```
Devuelve los topics cuyo `parent_id` coincide (flat, no árbol).

**Archivos a tocar:**
- `src/entrypoints/admin.go:150` — parsear `parent_id` de query y setearlo en `GetTopicsRequest`.
- `src/core/usecases/admin/get_topics.go` — agregar campo `ParentID *int64` al request y lógica de filtrado.
- `src/repositories/topic.go` — query con WHERE por parent_id cuando se pasa.

**Tests:**
- `GET /topics?parent_id=1` retorna sólo hijos directos.
- `GET /topics?level=2&parent_id=1` combina ambos filtros.

---

### G-4 · `PUT /areas/:id` — actualizar área

**Referencia RFC:** HU-3.2 criterio "CRUD endpoints para áreas: POST, GET (listar), PUT".

**Contrato esperado:**
```
PUT /api/v1/areas/:id
Body: { name?: string, description?: string }
```
Response `200` con el área actualizada.

**Archivos a tocar:**
- `src/core/usecases/admin/update_area.go` — nuevo caso de uso.
- `src/entrypoints/admin.go` — nuevo `HandleUpdateArea`.
- `src/app/web/mapping.go` — `coordOnly.PUT("/areas/:id", ...)`.
- `src/repositories/area.go` — método `Update`.

**Tests:**
- Update parcial sólo modifica los campos enviados.
- Multi-tenancy: no permite editar un área de otra org.

---

### G-5 · `POST /courses` — alinear contrato con RFC

**Referencia RFC:** `tecnico/endpoints.md` §191-215 — request `{ "name": "3ro 1era" }`, response sin `year`.

**Estado actual:**
- `src/entrypoints/courses.go:29` — `createCourseBody { Name, Year }` acepta `year` extra.
- `src/core/entities/course.go` — `Course.Year int` existe en el modelo.
- Response JSON incluye `year`.

**Discrepancia:** el RFC de HU-3.4 tampoco lista `year` en la migración de `courses` (sólo `id, organization_id, name, created_at`). `school_year` vive en `course_subjects`, no en `courses`.

**Decisión pendiente:** ¿año por curso (BE actual) o año por course_subject (RFC)? Impacta:
- Migración de la tabla.
- Schema del body/response.
- Dónde consulta el FE el "año actual" del curso.

**Archivos a tocar (si se alinea con RFC):**
- Migración que dropee `courses.year`.
- Entity `Course` sin `Year`.
- `createCourseBody` sólo con `Name`.

**Alternativa pragmática:** dejar el BE como está y **actualizar el RFC** para que documente el campo `year` en `POST /courses` (superset del RFC).

---

### G-6 · `GET /areas` con preload de subjects (nice-to-have)

**Referencia RFC:** HU-3.2 criterio "Listar áreas incluye las disciplinas asociadas (preload)".

**Estado actual:** `Admin.HandleListAreas` retorna `{ items: Area[], more }` con `Area { id, name, description, created_at }`. No trae subjects embebidos.

**Por qué importa:** el FE hoy hace N+1 queries para armar el dropdown "Área > Disciplina". Con el preload se resuelve en una sola request.

**Impacto:** bajo — el FE puede funcionar sin esto si G-1 (`GET /subjects`) está.

---

## Priorización sugerida

| Orden | Gap | Bloquea FE | Esfuerzo |
|---|---|---|---|
| 1 | G-1 `GET /subjects` | ✅ sí (bootstrap roto) | S |
| 2 | G-2 `GET /course-subjects` | ✅ sí (páginas docente) | M |
| 3 | G-3 `GET /topics?parent_id=` | ⚠️ parcial (FE filtra en memoria) | XS |
| 4 | G-4 `PUT /areas/:id` | ❌ no (feature admin) | S |
| 5 | G-5 alinear `POST /courses` | ❌ no | decisión de producto |
| 6 | G-6 preload subjects en `/areas` | ❌ no | S |

**Recomendación:** hacer G-1 + G-2 + G-3 en el mismo PR — juntos desbloquean el FE y son cambios localizados (mismos módulos `admin`/`courses`). G-4 como PR separado. G-5 y G-6 a discutir.

---

## Notas adicionales

- **Contenedores stubbeados:** `Coordination`, `Teaching`, `Resources` están declarados en `containers.go` pero sin rutas en `mapping.go`. Corresponden a Épicas 4, 5 y 8 — **no** son parte de Épica 3.
- **Onboarding:** completo. El RFC marca `T-2.2.2 (endpoint de perfil)` como Post-MVP pero ya está implementado — es un adelanto, no un gap.
- **Seed (`db/seeds/seed.sql`):** alineado con el RFC (`features`, `skip_allowed`, `tour_steps` con `key`). No tocar.
- **Multi-tenancy:** verificar en todos los nuevos endpoints que el `WHERE organization_id = ?` venga del JWT, no del request.
