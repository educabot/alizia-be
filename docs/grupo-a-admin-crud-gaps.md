# Grupo A — Gaps de Admin CRUD pendientes

**Fecha:** 2026-04-17
**Origen:** Análisis FE ↔ BE contra `rfc-alizia/epicas`. Son endpoints que el FE necesita para completar la administración de taxonomía y asignaciones, pero que **ni existen hoy ni están planificados en el RFC** (o están solo como modelo de datos, sin handlers).
**Scope:** solo endpoints admin/coordinator que completan CRUD existente. No incluye los "contenedores vacíos" (coordinación, planificación, recursos) — esos sí están cubiertos por las Épicas 4, 5, 8.

---

## TL;DR — Qué falta

| Recurso | GET list | GET id | POST | PATCH | DELETE | Estado actual |
|---|---|---|---|---|---|---|
| **areas** | ✅ | ✅ | ✅ | ✅ (`PUT`) | ✅ | Completo (referencia) |
| **subjects** | ✅ | ⚠️ (solo vía `/areas/:id/subjects`) | ✅ | ❌ **falta** | ❌ **falta** | Incompleto |
| **topics** | ✅ (tree) | ⚠️ (no hay GET individual) | ✅ | ✅ | ❌ **falta** | Incompleto |
| **courses** | ✅ | ✅ | ✅ | ❌ **falta** | ❌ **falta** | Incompleto |
| **course-subjects** | ✅ | ✅ | ✅ | ❌ **falta** | ❌ **falta** | Incompleto |
| **activities** | ✅ | ❌ **falta** | ✅ | ❌ **falta** | ❌ **falta** | Muy incompleto |
| **users** | ❌ **falta** | ❌ **falta** (solo `/users/me/*`) | — (fuera de scope) | — | — | No listable |

**Resumen:**
- **10 endpoints nuevos** a crear (7 PATCH/DELETE + 1 GET activity + 1 GET subject + 1 GET /users).
- Areas queda como referencia de cómo implementar el patrón completo (ver `src/core/usecases/admin/delete_area.go` con `CountDependencies` para devolver 409 en lugar de cascadear).
- GET /users es el más urgente para el FE: sin él, asignación de coordinators y de `teacher_id` en course-subjects se hacen con inputs manuales numéricos.

---

## Patrones de referencia (no reinventar)

Todos los endpoints nuevos deben seguir la arquitectura existente:

1. **Usecase en `src/core/usecases/admin/`** con struct `XxxRequest { OrgID uuid.UUID; ID int64; ... }`, método `Validate() error` y `Execute(ctx, req) error`.
2. **Handler en `src/entrypoints/admin.go`** (o `courses.go` según dominio) que parsea param, invoca usecase, devuelve `web.OK(mapXxx(result))` o `rest.HandleError(err)`.
3. **Errores estándar** vía `providers.ErrValidation`, `providers.ErrNotFound`, `providers.ErrConflict` — `rest.HandleError` los traduce a HTTP (400, 404, 409).
4. **Multi-tenant:** siempre leer `middleware.OrgID(req)` del JWT, nunca confiar en body/query.
5. **DTOs explícitos** (no entidades GORM crudas) — seguir la misma convención de `frontend-breaking-changes-dtos.md`: sin `created_at`, `updated_at`, `organization_id`.
6. **DELETE con dependencias:** antes de borrar, contar referencias y devolver `ErrConflict` con mensaje accionable si hay dependencias. **No cascadear** — es destructivo e irrecuperable. Patrón de referencia: `delete_area.go:44-64`.
7. **Rutas en `src/app/web/mapping.go`:** `coordOnly` para POST/PATCH, `adminOnly` para DELETE (consistente con areas).
8. **Tests:** cada usecase con `xxx_test.go` cubriendo validation, not found, conflict, happy path.

---

## 1. Subjects — PATCH y DELETE

### `PATCH /api/v1/subjects/:id`

**Auth:** `coordOnly` (coordinator o admin)
**Propósito:** renombrar/editar descripción o mover a otra área.

**Request body:**
```json
{
  "name": "Matemática III",       // opcional
  "description": "...",            // opcional
  "area_id": 5                     // opcional — mover a otra área
}
```

**Reglas:**
- Al menos un campo presente (si no, `400`).
- Si `area_id` viene, validar que el área pertenezca a la org del caller.
- Si `name` + `area_id` ya existe en otra subject → `409 Conflict`.

**Response:**
```json
{
  "description": {
    "id": 12,
    "name": "Matemática III",
    "description": "...",
    "area_id": 5
  }
}
```

### `DELETE /api/v1/subjects/:id`

**Auth:** `adminOnly`
**Propósito:** eliminar una subject que no está en uso.

**Dependencias a chequear antes de borrar:**
- `course_subjects` que la referencien → si > 0, `409` con mensaje: `"subject has dependencies (N course-subjects); remove them before deleting"`.
- `topics` vinculados a la subject (si el modelo lo permite).

**Response 204** (sin body).

---

## 2. Topics — DELETE (con cascada controlada del subárbol)

### `DELETE /api/v1/topics/:id`

**Auth:** `adminOnly`
**Propósito:** eliminar un tema de la jerarquía curricular.

**Comportamiento crítico — decidir y documentar:**

| Opción | Pros | Contras |
|---|---|---|
| **A) Rechazar si tiene hijos** (consistente con areas) | Simple, seguro, predecible | Obliga al usuario a borrar hoja por hoja |
| **B) Cascadear todo el subárbol** (declarado en HU-3.3 del RFC) | UX más ágil | Destructivo; si hay referencias en lesson-plans → conflicto |

**Recomendación:** Opción A por consistencia con `delete_area.go`. Si producto pide B, que sea opt-in vía `?cascade=true` y aun así rechace si hay `lesson_plans` o `coordination_documents` referenciando nodos del subárbol.

**Dependencias a chequear:**
- `topics WHERE parent_id = :id` (hijos directos).
- Referencias desde `coord_doc_class_topics`, `lesson_plan_topics`, etc. (cuando existan las tablas).

**Response 204** en éxito, `409` si hay dependencias.

---

## 3. Courses — PATCH y DELETE

### `PATCH /api/v1/courses/:id`

**Auth:** `adminOnly`
**Propósito:** corregir nombre, cambiar ciclo lectivo o nivel.

**Request body:**
```json
{
  "name": "3ro A",              // opcional
  "level": "secundaria",         // opcional
  "school_year": 2026            // opcional
}
```

**Reglas:** al menos un campo. Unicidad (`name` + `school_year` por org) → `409` si duplica.

### `DELETE /api/v1/courses/:id`

**Auth:** `adminOnly`

**Dependencias a chequear:**
- `course_subjects WHERE course_id = :id` → `409`.
- `students` del curso → `409` (o decisión de producto: ¿se archivan alumnos?).
- `time_slots WHERE course_id = :id`.

**Response 204.**

---

## 4. Course-Subjects — PATCH y DELETE (alta prioridad FE)

Usados para reasignar docente y dar de baja una asignación.

### `PATCH /api/v1/course-subjects/:id`

**Auth:** `adminOnly`
**Propósito:** cambiar docente asignado, fechas vigentes o ciclo.

**Request body:**
```json
{
  "teacher_id": 42,              // opcional — cambio de docente
  "start_date": "2026-03-01",    // opcional
  "end_date": "2026-12-15",      // opcional
  "school_year": 2026            // opcional
}
```

**Reglas:**
- Si `teacher_id` viene, validar que el user exista en la org y tenga rol `teacher`.
- `start_date <= end_date` si ambas presentes.
- No permitir overlapping con otra course-subject activa del mismo curso+subject (si es regla de negocio).

**Response:** el course-subject actualizado (mismo shape que `POST`).

### `DELETE /api/v1/course-subjects/:id`

**Auth:** `adminOnly`

**Dependencias a chequear:**
- `lesson_plans WHERE course_subject_id = :id` → `409` (cuando exista la tabla).
- `time_slots` que incluyan esta course-subject → desasociar del time_slot o `409`.

**Response 204.**

**Nota:** este es el CRUD más pedido por el FE. Sin PATCH no se puede reasignar docente desde la UI de admin, que es el caso más común durante el año lectivo.

---

## 5. Activities — GET detail, PATCH y DELETE

### `GET /api/v1/activities/:id`

**Auth:** any authenticated.
**Motivo:** Hoy solo hay list. El FE necesita detail para el editor.

### `PATCH /api/v1/activities/:id`

**Auth:** `adminOnly`
**Request body:** cualquier subset de los campos que acepta `POST /activities` (name, description, moment, tags, etc.).

### `DELETE /api/v1/activities/:id`

**Auth:** `adminOnly`

**Dependencias a chequear:**
- Referencias desde `lesson_plans` o `coordination_documents` (cuando existan).

---

## 6. Users — listado (bloqueador FE)

### `GET /api/v1/users`

**Auth:** `adminOnly` (y tal vez `coordOnly` si coordinators necesitan listar docentes de su área).
**Propósito:** poblar dropdowns de:
- Asignación de coordinator a un área (`POST /areas/:id/coordinators` necesita `user_id`).
- Asignación de `teacher_id` al crear/editar course-subject.

**Query params:**
```
?role=teacher           # opcional — filtrar por rol
?role=coordinator
?area_id=5              # opcional — solo coords/teachers de un área (si aplica)
?search=juan            # opcional — LIKE en name/email
?limit=50&offset=0
```

**Response:**
```json
{
  "description": {
    "items": [
      {
        "id": 42,
        "name": "Juan Pérez",
        "email": "juan@org.edu.ar",
        "roles": ["teacher"]
      }
    ],
    "total": 128,
    "has_more": true
  }
}
```

**Campos excluidos (consistente con DTO convention):** `password_hash`, `profile_data`, `onboarding_completed_at`, `created_at`, `updated_at`, `organization_id`.

**Nota sobre `auth-service-futuro.md`:** el documento posterga la gestión completa de users a un servicio externo. Pero **listado read-only filtrado por org** es necesario ya y no implica mover la autenticación — puede vivir en Alizia-be como view y migrarse cuando el auth-service esté listo.

---

## Priorización sugerida

| Prioridad | Endpoint | Razón |
|---|---|---|
| **P0** | `GET /users` (con filtros por rol) | Desbloquea 2 flujos admin que hoy usan inputs manuales numéricos |
| **P0** | `PATCH /course-subjects/:id` | Reasignación de docente es caso de uso frecuente |
| **P1** | `DELETE /course-subjects/:id` | Baja de asignación |
| **P1** | `PATCH + DELETE /courses/:id` | Corrección y limpieza de taxonomía |
| **P2** | `PATCH + DELETE /subjects/:id` | Menor frecuencia; errores de alta son raros |
| **P2** | `DELETE /topics/:id` | El RFC lo declara en HU-3.3 pero falta implementación |
| **P2** | `GET + PATCH + DELETE /activities/:id` | Catálogo didáctico — baja rotación |

---

## Checklist de implementación por endpoint

Para cada uno de los 10 endpoints nuevos:

- [ ] Usecase en `src/core/usecases/admin/` con `Validate()` + `Execute()`.
- [ ] Test del usecase (`_test.go`) — validation, not found, conflict, happy path.
- [ ] Repository method en `src/repositories/admin/` (o `courses/`) si no existe.
- [ ] Handler en `src/entrypoints/admin.go` o `courses.go`.
- [ ] Mapping en `mapping.go` bajo `coordOnly` (PATCH) o `adminOnly` (DELETE).
- [ ] DTO response en `mapping.go` del paquete web (sin campos excluidos).
- [ ] Para DELETE: `CountDependencies` + `ErrConflict` con mensaje accionable.
- [ ] Entrada en `frontend-breaking-changes-dtos.md` si cambia shape de response existente.

---

## Relación con el RFC

Abrir HUs nuevas (o extender existentes) en:

- **HU-3.2 (Areas y Materias)** → agregar tareas `PATCH /subjects/:id`, `DELETE /subjects/:id`.
- **HU-3.3 (Topics)** → T-3.3.5 ya menciona DELETE; crear la tarea de implementación.
- **HU-3.4 (Courses, Course-Subjects)** → agregar `PATCH + DELETE` para ambos.
- **HU-3.6 (Activities)** → agregar `GET /:id`, `PATCH`, `DELETE`.
- **HU-1.2 (Modelo users/roles)** → agregar tarea `GET /users` con filtros (read-only, no migra a auth-service).

---

## Fuera de scope de este documento

- **Contenedores vacíos** (coordinación, planificación, recursos, chat, dashboard, notificaciones): ya cubiertos por Épicas 4, 5, 6, 7, 8. No repetir acá.
- **Gestión completa de users** (create/update/delete de users, password reset, invitations): delegado al auth-service futuro.
- **`/fonts` CRUD**: gap identificado en HU-8.1 pero corresponde a Épica 8.
