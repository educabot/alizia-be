# Catalogo de errores — Alizia API

Todos los errores siguen el formato estandar de `team-ai-toolkit/errors`. La respuesta HTTP contiene:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Descripcion legible del error"
  }
}
```

- `code`: constante en UPPER_SNAKE_CASE, estable para programmatic matching
- `message`: texto descriptivo, puede cambiar entre versiones

---

## Errores genericos (team-ai-toolkit/errors)

Estos errores vienen de la libreria compartida y aplican a todos los endpoints.

| HTTP | Code | Descripcion | Cuando ocurre |
|------|------|-------------|---------------|
| 400 | `VALIDATION_ERROR` | Request body invalido o campos requeridos faltantes | Body mal formado, campos vacios, tipos incorrectos |
| 401 | `UNAUTHORIZED` | Token JWT ausente, expirado o invalido | Header Authorization faltante o JWT invalido |
| 403 | `FORBIDDEN` | Usuario no tiene el rol requerido | Ej: teacher intenta crear area (requiere admin) |
| 404 | `NOT_FOUND` | Recurso no encontrado | ID inexistente o recurso de otra org |
| 409 | `DUPLICATE` | Registro duplicado viola constraint UNIQUE | Ej: area con mismo nombre en la org |
| 409 | `CONFLICT` | Conflicto de estado o concurrencia | Ej: operacion concurrente sobre el mismo recurso |
| 500 | `INTERNAL_ERROR` | Error interno del servidor | Error no manejado, falla de DB, etc. |

> **Nota sobre multi-tenancy:** Un `NOT_FOUND` se retorna tanto si el recurso no existe como si pertenece a otra organizacion. Nunca se expone informacion cross-tenant.

---

## Errores de dominio — Admin (Fase 2)

### Areas

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 409 | `DUPLICATE_AREA` | Ya existe un area con ese nombre en la org | POST /areas, PUT /areas/:id |
| 404 | `FORBIDDEN_AREA_NOT_IN_ORG` | El area no pertenece a la org del usuario | PUT /areas/:id |

### Coordinators

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `USER_NOT_COORDINATOR` | El usuario no tiene rol coordinator | POST /areas/:id/coordinators |
| 409 | `DUPLICATE_COORDINATOR` | El usuario ya es coordinador de esa area | POST /areas/:id/coordinators |

### Subjects

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 409 | `DUPLICATE_SUBJECT` | Ya existe una disciplina con ese nombre en la misma area | POST /subjects |

### Courses

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 409 | `DUPLICATE_COURSE` | Ya existe un curso con ese nombre en la org | POST /courses |

### Course Subjects

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 409 | `DUPLICATE_COURSE_SUBJECT` | Ya existe esa combinacion curso + disciplina + school_year | POST /course-subjects |

### Topics

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `TOPIC_EXCEEDS_MAX_LEVELS` | El topic excede `config.topic_max_levels` de la org | POST /topics |
| 400 | `TOPIC_WRONG_LEVEL` | El topic no esta al nivel `config.topic_selection_level` | POST /coordination-documents |
| 400 | `TOPIC_NOT_IN_ORG` | El topic no pertenece a la org | POST /coordination-documents |

### Time Slots

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `SHARED_CLASSES_DISABLED` | Se enviaron 2 course_subjects pero `shared_classes_enabled = false` | POST /courses/:id/time-slots |
| 400 | `COURSE_SUBJECT_WRONG_COURSE` | El course_subject no pertenece al curso del time_slot | POST /courses/:id/time-slots |
| 409 | `TIME_SLOT_OVERLAP` | Ya existe un time_slot en ese dia/hora para el curso | POST /courses/:id/time-slots |

### Activities

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `INVALID_MOMENT` | El momento no es apertura, desarrollo o cierre | POST /activities |

---

## Errores de dominio — Coordination Documents (Fase 3)

### Wizard de creacion

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `SUBJECT_NOT_IN_AREA` | Una disciplina no pertenece al area del documento | POST /coordination-documents |
| 400 | `TOPICS_NOT_FULLY_DISTRIBUTED` | No todos los topics del doc estan asignados a alguna disciplina | POST /coordination-documents |

### Edicion y estados

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `DOCUMENT_NOT_DRAFT` | El documento no esta en estado draft (no se puede editar/eliminar) | PATCH, DELETE, POST /generate |
| 400 | `INVALID_SECTION_KEY` | La key de seccion no existe en `config.coord_doc_sections` | PATCH /coordination-documents/:id |
| 400 | `DOCUMENT_NO_TOPICS` | El documento no tiene topics asignados (no se puede generar) | POST /generate |

### Publicacion

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `DOCUMENT_MISSING_REQUIRED_SECTIONS` | Faltan secciones marcadas como `required` en la config | PATCH (status → published) |
| 400 | `DOCUMENT_NO_CLASS_PLANS` | No hay plan de clases generado para alguna disciplina | PATCH (status → published) |
| 400 | `TOPICS_NOT_FULLY_COVERED` | Hay topics de disciplina no asignados a ninguna clase | PATCH (status → published) |
| 400 | `INVALID_STATUS_TRANSITION` | Transicion de estado no permitida (ej: archived → draft) | PATCH (status change) |

### Chat

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `INVALID_TOOL_CALL` | Alizia intento ejecutar un tool con parametros invalidos | POST /chat |

---

## Errores de dominio — Teaching (Fase 5)

### Lesson Plans

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `NO_PUBLISHED_DOCUMENT` | No hay documento de coordinacion publicado para la disciplina | GET /course-subjects/:id/lesson-plans, POST /lesson-plans |
| 409 | `LESSON_PLAN_ALREADY_EXISTS` | Ya existe un lesson plan para esa clase y course_subject | POST /lesson-plans |
| 400 | `INVALID_MOMENT_ACTIVITIES` | Cantidad de actividades invalida por momento | POST, PATCH /lesson-plans |
| 400 | `TOPIC_NOT_IN_DOCUMENT` | Un topic no esta asignado a la disciplina en el doc de coordinacion | POST /lesson-plans |
| 400 | `ACTIVITY_NOT_IN_MOMENT` | La actividad no corresponde al momento indicado | POST /lesson-plans/:id/generate-activity |
| 400 | `LESSON_PLAN_INCOMPLETE` | El lesson plan no tiene todos los momentos completos | PATCH /lesson-plans/:id/status |

#### Detalle: `INVALID_MOMENT_ACTIVITIES`

Se retorna cuando:
- `apertura` tiene != 1 actividad
- `desarrollo` tiene < 1 o > `config.desarrollo_max_activities` actividades
- `cierre` tiene != 1 actividad
- Algun momento falta en el request
- Un activity_id no existe o no corresponde al momento (ej: actividad de apertura usada en cierre)

El mensaje incluye detalle:
```json
{
  "error": {
    "code": "INVALID_MOMENT_ACTIVITIES",
    "message": "Momento 'desarrollo' tiene 4 actividades, maximo permitido: 3",
    "details": {
      "moment": "desarrollo",
      "count": 4,
      "max": 3
    }
  }
}
```

---

## Errores de dominio — Resources (Fase 6)

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 400 | `RESOURCE_TYPE_REQUIRES_FONT` | El tipo de recurso requiere una fuente pero no se envio `font_id` | POST /resources |
| 400 | `RESOURCE_TYPE_DISABLED` | El tipo de recurso esta deshabilitado para la org | POST /resources |
| 400 | `RESOURCE_NOT_DRAFT` | El recurso no esta en estado draft (no se puede regenerar) | POST /resources/:id/generate |
| 400 | `RESOURCE_NO_FONT` | El recurso requiere font pero no tiene una asignada | POST /resources/:id/generate |
| 400 | `CONTENT_SCHEMA_MISMATCH` | El content enviado no matchea el output_schema del tipo | PATCH /resources/:id |

---

## Errores de dominio — AI (Fase 4)

| HTTP | Code | Descripcion | Endpoint |
|------|------|-------------|----------|
| 502 | `AI_GENERATION_ERROR` | Error al comunicarse con Azure OpenAI | Todos los /generate, /chat |
| 504 | `AI_GENERATION_TIMEOUT` | Azure OpenAI no respondio a tiempo (> 30s) | Todos los /generate, /chat |
| 429 | `AI_RATE_LIMITED` | Rate limit de Azure OpenAI alcanzado | Todos los /generate, /chat |

### Detalle: `AI_GENERATION_ERROR`

Errores de la API de Azure OpenAI se mapean asi:

| Azure OpenAI error | Alizia code | HTTP |
|---------------------|-------------|------|
| 400 (prompt invalido) | `AI_GENERATION_ERROR` | 502 |
| 401/403 (auth) | `AI_GENERATION_ERROR` | 502 |
| 429 (rate limit) | `AI_RATE_LIMITED` | 429 |
| 500+ (server error) | `AI_GENERATION_ERROR` | 502 |
| Timeout | `AI_GENERATION_TIMEOUT` | 504 |

> El backend nunca expone detalles internos de Azure OpenAI al cliente. El mensaje generico es "Error al generar contenido con IA. Intenta nuevamente."

---

## Transiciones de estado validas

### Coordination Document

```
pending → in_progress    (comenzar edición)
in_progress → published  (publicar)
```

Cualquier otra transicion retorna `INVALID_STATUS_TRANSITION`.

### Lesson Plan

```
pending → in_progress → published    (marcar como en progreso, luego publicar)
```

### Resource

```
draft → active       (activar recurso)
```

---

## Implementacion en Go

Los errores de dominio de Alizia extienden `team-ai-toolkit/errors`:

```go
// src/core/providers/errors.go

package providers

import "github.com/educabot/team-ai-toolkit/errors"

// Admin
var ErrDuplicateArea = errors.NewValidation("DUPLICATE_AREA", "Ya existe un area con ese nombre en la organizacion")
var ErrTopicExceedsMaxLevels = errors.NewValidation("TOPIC_EXCEEDS_MAX_LEVELS", "El topic excede la cantidad maxima de niveles configurada")
var ErrSharedClassesDisabled = errors.NewValidation("SHARED_CLASSES_DISABLED", "Las clases compartidas no estan habilitadas para esta organizacion")

// Coordination
var ErrDocumentNotDraft = errors.NewValidation("DOCUMENT_NOT_DRAFT", "El documento debe estar en estado borrador para esta operacion")
var ErrInvalidSectionKey = errors.NewValidation("INVALID_SECTION_KEY", "La seccion no existe en la configuracion de la organizacion")

// Teaching
var ErrInvalidMomentActivities = errors.NewValidation("INVALID_MOMENT_ACTIVITIES", "Cantidad de actividades invalida para el momento")
var ErrLessonPlanAlreadyExists = errors.NewDuplicate("LESSON_PLAN_ALREADY_EXISTS", "Ya existe una planificación docente para esta clase")

// AI
var ErrAIGeneration = errors.New("AI_GENERATION_ERROR", "Error al generar contenido con IA")
var ErrAITimeout = errors.New("AI_GENERATION_TIMEOUT", "Timeout al generar contenido con IA")
var ErrAIRateLimited = errors.New("AI_RATE_LIMITED", "Limite de solicitudes de IA alcanzado")
```

El `HandleError()` de team-ai-toolkit mapea automaticamente:
- `errors.NewValidation()` → HTTP 400
- `errors.NewDuplicate()` → HTTP 409
- `errors.NewNotFound()` → HTTP 404

Para los errores de AI (502, 504, 429) se usa un handler custom en `src/entrypoints/rest/rest.go` que extiende el HandleError base.

---

## Notas para el frontend

1. **Siempre chequear `error.code`**, no el HTTP status ni el mensaje
2. **Errores 401** redirigen a login
3. **Errores 403** muestran "No tenes permisos para esta accion"
4. **Errores 404** cross-tenant son indistinguibles de "no existe" (por diseno)
5. **Errores 429/502/504** de IA: mostrar mensaje amigable + boton "Reintentar"
6. **`INVALID_MOMENT_ACTIVITIES`** incluye `details` con el momento y los limites — usar para mostrar error inline
7. **`DOCUMENT_MISSING_REQUIRED_SECTIONS`** se puede prevenir validando contra `org_config.coord_doc_sections` antes de intentar publicar
