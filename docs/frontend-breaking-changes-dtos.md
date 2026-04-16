# Breaking Changes para el Frontend — Migración a Response DTOs

**Branch:** `feature/sl/epica-3-integracion`
**Fecha:** 2026-04-16
**Scope:** Todos los endpoints bajo `/api/v1` que devolvían entidades GORM crudas pasaron a devolver DTOs explícitos. Esto cambia la forma del JSON.

---

## TL;DR — Qué cambia globalmente

En **todas** las responses afectadas:

| Campo | Antes | Ahora |
|---|---|---|
| `created_at` | presente | **removido** |
| `updated_at` | presente | **removido** |
| `organization_id` | presente | **removido** (ya viene en el JWT) |

En cualquier `user` embebido (coordinator, teacher, etc.):

| Campo | Antes | Ahora |
|---|---|---|
| `password_hash` | nunca apareció (ya filtrado) | sigue oculto |
| `profile_data` | presente como `{}` | **removido** |
| `onboarding_completed_at` | presente | **removido** |
| `roles` | presente | **removido** |

> **Recomendación:** si el front guarda estos campos en alguna interfaz TypeScript o store, hay que sacarlos. Si los muestra en UI, hay que pedirlos por endpoint específico (no vienen de "free" en cada listado).

El envelope `{"description": ...}` **NO cambia** — sigue igual que antes (es convención del toolkit).

---

## Endpoints afectados — uno por uno

### 1. Areas

#### `GET /api/v1/areas` — listar áreas
#### `POST /api/v1/areas` — crear área

**Antes (entidad cruda):**
```json
{
  "description": {
    "items": [{
      "id": 1,
      "organization_id": "a0eebc99-...",
      "name": "Ciencias",
      "description": "...",
      "subjects": [{
        "id": 5,
        "organization_id": "a0eebc99-...",
        "area_id": 1,
        "name": "Matemática",
        "description": "...",
        "created_at": "2026-04-16T12:00:00Z",
        "updated_at": "2026-04-16T12:00:00Z"
      }],
      "coordinators": [{
        "id": 1,
        "area_id": 1,
        "user_id": 2,
        "created_at": "2026-04-16T12:00:00Z",
        "user": {
          "id": 2,
          "organization_id": "a0eebc99-...",
          "email": "coord@neuquen.edu.ar",
          "first_name": "Carlos",
          "last_name": "Coordinador",
          "avatar_url": null,
          "onboarding_completed_at": "2026-04-16T17:55:11Z",
          "profile_data": {},
          "roles": [{"id": 2, "user_id": 2, "role": "coordinator"}],
          "created_at": "2026-04-16T12:00:00Z",
          "updated_at": "2026-04-16T12:00:00Z"
        }
      }],
      "created_at": "2026-04-16T12:00:00Z",
      "updated_at": "2026-04-16T12:00:00Z"
    }],
    "more": false
  }
}
```

**Ahora (DTO):**
```json
{
  "description": {
    "items": [{
      "id": 1,
      "name": "Ciencias",
      "description": "...",
      "subjects": [{
        "id": 5,
        "area_id": 1,
        "name": "Matemática",
        "description": "..."
      }],
      "coordinators": [{
        "id": 1,
        "area_id": 1,
        "user": {
          "id": 2,
          "email": "coord@neuquen.edu.ar",
          "first_name": "Carlos",
          "last_name": "Coordinador",
          "avatar_url": null
        }
      }]
    }],
    "more": false
  }
}
```

**Cambios clave:**
- `subjects` y `coordinators` siempre son arrays (nunca `null`). Vacíos llegan como `[]`.
- `coordinator.user_id` **removido** (ahora viene dentro de `user.id`).
- `coordinator.user` puede ser `null` (sucede en la response de `POST /areas/:id/coordinators` porque el usecase no popula User).
- `coordinator.created_at` **removido**.

---

### 2. Subjects

#### `GET /api/v1/areas/:id/subjects`
#### `POST /api/v1/subjects`

**Ahora:**
```json
{ "id": 5, "area_id": 1, "name": "Matemática", "description": "..." }
```

Removido: `organization_id`, `created_at`, `updated_at`.

---

### 3. Topics

#### `GET /api/v1/topics?level=N`
#### `POST /api/v1/topics`
#### `PATCH /api/v1/topics/:id`

**Ahora:**
```json
{
  "id": 10,
  "parent_id": null,
  "name": "Álgebra",
  "description": "...",
  "level": 1,
  "children": [
    { "id": 11, "parent_id": 10, "name": "Ecuaciones", "level": 2 }
  ]
}
```

**Cambios:**
- `children` se omite si está vacío (`omitempty`). Antes podía venir como `null` o `[]`.
- Removido: `organization_id`, `created_at`, `updated_at`.

---

### 4. Activities (Templates)

#### `GET /api/v1/activities?moment=apertura|desarrollo|cierre`
#### `POST /api/v1/activities`

> **Nota:** son templates de actividad por momento de clase, no instancias dentro de un lesson plan.

**Ahora:**
```json
{
  "id": 1,
  "moment": "apertura",
  "name": "Lluvia de ideas",
  "description": "...",
  "duration_minutes": 15
}
```

**Cambios:**
- Renombrado: `moment_type` → `moment`, `title` → `name`, `duration_min` → `duration_minutes`.
- Removidos: `lesson_plan_id`, `sort_order` (no aplican a templates), `organization_id`, `created_at`, `updated_at`.
- `description` y `duration_minutes` ahora son opcionales (omitempty).

---

### 5. Organization

#### `GET /api/v1/organization`
#### `PATCH /api/v1/organization/config`

**Ahora:**
```json
{
  "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "name": "Escuela Neuquén",
  "slug": "neuquen",
  "config": { "...": "..." }
}
```

**Cambios:**
- `config` ahora es un objeto JSON parseado (`map[string]any`) en lugar de un string base64/raw JSONB.
- Removido: `created_at`, `updated_at`.

---

### 6. Courses

#### `GET /api/v1/courses` — listar
#### `POST /api/v1/courses` — crear
#### `GET /api/v1/courses/:id` — detalle

**Ahora:**
```json
{
  "id": 1,
  "name": "1° A",
  "students": [
    { "id": 1, "course_id": 1, "name": "Juan Pérez" }
  ],
  "course_subjects": [{
    "id": 1,
    "course_id": 1,
    "subject_id": 5,
    "teacher_id": 3,
    "school_year": 2026,
    "start_date": "2026-03-01",
    "end_date": "2026-12-15",
    "subject": { "id": 5, "name": "Matemática" },
    "teacher": { "id": 3, "first_name": "María", "last_name": "Docente" }
  }]
}
```

**Cambios clave:**
- `Course` ya no tiene `year` (migración `000013_drop_courses_year.up.sql`). El curso es atemporal; el ciclo lectivo vive en `course_subjects[].school_year`.
- `students` y `course_subjects` siempre son arrays (`[]` si vacíos, nunca `null`).
- `start_date`/`end_date` ahora son **strings `YYYY-MM-DD`** (antes era timestamp ISO 8601 completo). Si están vacíos o no seteados se omiten (`omitempty`).
- `subject` (dentro de `course_subjects`) es un objeto compacto `{id, name}` — no la entidad Subject completa.
- `teacher` es un objeto compacto `{id, first_name, last_name}` — sin email, avatar, roles, etc.
- Removido en student/course/course_subject: `organization_id`, `created_at`, `updated_at`.
- **`POST /course-subjects` preloadea `subject` y `teacher` en la response** (fix G-10): la creación devuelve la entidad recargada desde el repo con ambos campos populados, no el struct en memoria.

---

### 7. CourseSubject (asignación profesor-materia-curso)

#### `POST /api/v1/course-subjects`

Misma estructura que `course_subjects[]` arriba. Sin envolver en course.

---

### 8. Students

#### `POST /api/v1/courses/:id/students`

```json
{ "id": 1, "course_id": 1, "name": "Juan Pérez" }
```

Removido: `created_at`, `updated_at`.

---

## Endpoints que NO cambiaron

Estos ya tenían su DTO o devuelven structs propios:

- `GET /api/v1/courses/:id/schedule` (`timeSlotResponse` ya existía)
- `POST /api/v1/courses/:id/time-slots`
- `GET /api/v1/course-subjects/:id/shared-class-numbers`
- `DELETE /api/v1/areas/:id/coordinators/:user_id` (devuelve 204 No Content)

---

## Checklist para el frontend

- [ ] Sacar `created_at` / `updated_at` / `organization_id` de todas las interfaces TS de Area, Subject, Topic, Course, Student, CourseSubject, Activity, Organization, User (resumido).
- [ ] Renombrar campos en `Activity`: `moment_type → moment`, `title → name`, `duration_min → duration_minutes`. Sacar `lesson_plan_id` y `sort_order`.
- [ ] Cambiar tipo de `start_date` / `end_date` en CourseSubject: ahora es `string` (`YYYY-MM-DD`), no `Date` ISO 8601 completo. Parsear con `new Date(s + 'T00:00:00')` o usar `date-fns`.
- [ ] En `Coordinator`: sacar `user_id` (usar `coordinator.user.id`), sacar `created_at`. Manejar `coordinator.user === null` en el render (puede pasar al asignar).
- [ ] En `User` (embebido en coordinator/teacher): sacar `roles`, `profile_data`, `onboarding_completed_at`. Si se necesita rol del usuario, pedirlo por endpoint específico.
- [ ] En `CourseSubject`: el `subject` embebido ahora es `{id, name}`, el `teacher` es `{id, first_name, last_name}`. Si la UI mostraba más datos del teacher (email, avatar), hay que cambiar el flujo o pedir el user completo aparte.
- [ ] En `Organization.config`: tratar como objeto JSON (`Record<string, any>`), ya viene parseado.
- [ ] Verificar que ningún componente espere `null` donde ahora viene `[]` (arrays vacíos): `subjects`, `coordinators`, `students`, `course_subjects`, `children` de topics.

---

## Por qué este cambio

- **Seguridad:** los DTOs garantizan que campos sensibles (`password_hash`, `profile_data` con datos privados) nunca se filtren accidentalmente al cambiar el schema de la DB.
- **Estabilidad de contrato:** cambios en GORM o en migraciones ya no rompen la API. El contrato vive en `src/entrypoints/`.
- **Performance:** el JSON es más chico (sin timestamps, sin campos zero-value, sin metadata redundante).
- **Convención:** documentado en `CLAUDE.md` y patrón canónico en `src/entrypoints/courses.go` (`timeSlotResponse`).

Cualquier duda, el código fuente de los DTOs está en:
- `src/entrypoints/admin.go` (Areas, Subjects, Topics, Activities, Organization, Coordinator)
- `src/entrypoints/courses.go` (Course, Student, CourseSubject, TimeSlot)
