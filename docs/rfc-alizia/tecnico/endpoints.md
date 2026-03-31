# Backend — Endpoints API

Todos los endpoints requieren autenticacion via Bearer token (JWT via team-ai-toolkit/tokens) salvo `/health`.

**Base URL:** `/api/v1`

**Convenciones:**
- Autenticacion: `Authorization: Bearer <JWT>`
- Content-Type: `application/json`
- Paginacion: query params `limit` (default 20, max 100) y `offset` (default 0)
- Respuesta paginada: `{ "items": [...], "more": true|false }`
- Filtro por org: automatico via JWT claims (el backend extrae `organization_id` del token)
- Errores: ver [errores.md](./errores.md) para catalogo completo

---

## Health (Fase 1)

### `GET /health`

Health check. No requiere autenticacion.

**Response `200`:**
```json
{ "status": "ok" }
```

---

## Admin (Fase 2)

### `POST /api/v1/areas`

Crear area.

**Roles:** coordinator, admin

**Request:**
```json
{
  "name": "Ciencias Exactas",
  "description": "Area de ciencias exactas y naturales"
}
```

**Response `201`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "name": "Ciencias Exactas",
  "description": "Area de ciencias exactas y naturales",
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `DUPLICATE_AREA`

---

### `GET /api/v1/areas`

Listar areas de la org.

**Roles:** Todos

**Query params:** `limit`, `offset`

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "organization_id": 1,
      "name": "Ciencias Exactas",
      "description": "Area de ciencias exactas y naturales",
      "created_at": "2026-03-25T10:00:00Z"
    }
  ],
  "more": false
}
```

---

### `PUT /api/v1/areas/:id`

Actualizar area.

**Roles:** coordinator, admin

**Request:**
```json
{
  "name": "Ciencias Naturales",
  "description": "Actualizada"
}
```

**Response `200`:** Area actualizada (mismo schema que POST).

**Errores:** `NOT_FOUND`, `FORBIDDEN_AREA_NOT_IN_ORG`, `DUPLICATE_AREA`

---

### `POST /api/v1/areas/:id/coordinators`

Asignar coordinador a un area.

**Roles:** admin

**Request:**
```json
{
  "user_id": 5
}
```

**Response `201`:**
```json
{
  "id": 1,
  "area_id": 1,
  "user_id": 5
}
```

**Errores:** `NOT_FOUND`, `USER_NOT_COORDINATOR`, `DUPLICATE_COORDINATOR`

---

### `POST /api/v1/subjects`

Crear disciplina.

**Roles:** coordinator, admin

**Request:**
```json
{
  "area_id": 1,
  "name": "Matematicas",
  "description": "Matematica aplicada"
}
```

**Response `201`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "area_id": 1,
  "name": "Matematicas",
  "description": "Matematica aplicada",
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (area), `DUPLICATE_SUBJECT`

---

### `GET /api/v1/subjects`

Listar disciplinas de la org.

**Roles:** Todos

**Query params:** `limit`, `offset`, `area_id` (opcional, filtra por area)

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "organization_id": 1,
      "area_id": 1,
      "name": "Matematicas",
      "description": "Matematica aplicada",
      "created_at": "2026-03-25T10:00:00Z"
    }
  ],
  "more": false
}
```

---

### `POST /api/v1/courses`

Crear curso.

**Roles:** admin

**Request:**
```json
{
  "name": "3ro 1era"
}
```

**Response `201`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "name": "3ro 1era",
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `DUPLICATE_COURSE`

---

### `GET /api/v1/courses`

Listar cursos de la org.

**Roles:** Todos

**Query params:** `limit`, `offset`

**Response `200`:** Paginado con items de curso (mismo schema que POST response).

---

### `GET /api/v1/courses/:id`

Detalle de curso con students y schedule.

**Roles:** Todos

**Response `200`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "name": "3ro 1era",
  "created_at": "2026-03-25T10:00:00Z",
  "students": [
    { "id": 1, "name": "Juan Perez", "created_at": "2026-03-25T10:00:00Z" }
  ],
  "time_slots": [
    {
      "id": 1,
      "day_of_week": 1,
      "start_time": "08:00",
      "end_time": "09:30",
      "subjects": [
        {
          "course_subject_id": 1,
          "subject_id": 1,
          "subject_name": "Matematicas",
          "teacher_id": 5,
          "teacher_name": "Doc. García"
        }
      ]
    },
    {
      "id": 3,
      "day_of_week": 3,
      "start_time": "08:00",
      "end_time": "09:30",
      "subjects": [
        {
          "course_subject_id": 1,
          "subject_id": 1,
          "subject_name": "Matematicas",
          "teacher_id": 5,
          "teacher_name": "Doc. García"
        },
        {
          "course_subject_id": 2,
          "subject_id": 2,
          "subject_name": "Fisica",
          "teacher_id": 6,
          "teacher_name": "Doc. López"
        }
      ]
    }
  ]
}
```

> Cuando `subjects` tiene 2 elementos, es una clase compartida.

**Errores:** `NOT_FOUND`

---

### `POST /api/v1/courses/:id/time-slots`

Crear time slot para un curso.

**Roles:** admin

**Request:**
```json
{
  "day_of_week": 1,
  "start_time": "08:00",
  "end_time": "09:30",
  "course_subject_ids": [1]
}
```

Para clase compartida, enviar 2 IDs:
```json
{
  "day_of_week": 3,
  "start_time": "08:00",
  "end_time": "09:30",
  "course_subject_ids": [1, 2]
}
```

**Response `201`:**
```json
{
  "id": 1,
  "course_id": 1,
  "day_of_week": 1,
  "start_time": "08:00",
  "end_time": "09:30",
  "subjects": [
    { "course_subject_id": 1, "subject_name": "Matematicas" }
  ]
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (course o course_subject), `SHARED_CLASSES_DISABLED`, `COURSE_SUBJECT_WRONG_COURSE`, `TIME_SLOT_OVERLAP`

---

### `POST /api/v1/topics`

Crear topic en la jerarquia curricular.

**Roles:** admin

**Request:**
```json
{
  "name": "Aritmetica basica",
  "description": "Operaciones fundamentales",
  "parent_id": 1
}
```

`parent_id` es null para topics raiz (level 1).

**Response `201`:**
```json
{
  "id": 5,
  "organization_id": 1,
  "parent_id": 1,
  "name": "Aritmetica basica",
  "description": "Operaciones fundamentales",
  "level": 2,
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `TOPIC_EXCEEDS_MAX_LEVELS`, `NOT_FOUND` (parent)

---

### `GET /api/v1/topics`

Listar topics como arbol.

**Roles:** Todos

**Query params:** `level` (opcional, filtrar por nivel), `parent_id` (opcional, hijos de un topic)

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "name": "Pensamiento logico-matematico",
      "level": 1,
      "children": [
        {
          "id": 2,
          "name": "Aritmetica",
          "level": 2,
          "children": [
            { "id": 5, "name": "Suma y resta", "level": 3, "children": [] }
          ]
        }
      ]
    }
  ]
}
```

> Sin query params devuelve el arbol completo. Con `level=N` devuelve solo ese nivel (flat). Con `parent_id=N` devuelve hijos directos.

---

### `POST /api/v1/course-subjects`

Crear asignacion curso + disciplina + docente.

**Roles:** admin

**Request:**
```json
{
  "course_id": 1,
  "subject_id": 1,
  "teacher_id": 5,
  "start_date": "2026-03-01",
  "end_date": "2026-11-30",
  "school_year": 2026
}
```

**Response `201`:**
```json
{
  "id": 1,
  "course_id": 1,
  "subject_id": 1,
  "teacher_id": 5,
  "start_date": "2026-03-01",
  "end_date": "2026-11-30",
  "school_year": 2026,
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (course, subject, teacher), `DUPLICATE_COURSE_SUBJECT`

---

### `GET /api/v1/course-subjects`

Listar asignaciones de la org.

**Roles:** Todos

**Query params:** `limit`, `offset`, `course_id`, `subject_id`, `teacher_id` (todos opcionales)

**Response `200`:** Paginado con items de course_subject (mismo schema que POST response).

---

### `POST /api/v1/activities`

Crear actividad didactica.

**Roles:** admin, coordinator

**Request:**
```json
{
  "moment": "desarrollo",
  "name": "Trabajo en grupo",
  "description": "Los alumnos trabajan colaborativamente",
  "duration_minutes": 30
}
```

**Response `201`:**
```json
{
  "id": 3,
  "organization_id": 1,
  "moment": "desarrollo",
  "name": "Trabajo en grupo",
  "description": "Los alumnos trabajan colaborativamente",
  "duration_minutes": 30,
  "created_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `INVALID_MOMENT`

---

### `GET /api/v1/activities`

Listar actividades de la org.

**Roles:** Todos

**Query params:** `moment` (opcional: `apertura`, `desarrollo`, `cierre`)

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "moment": "apertura",
      "name": "Lluvia de ideas",
      "description": "...",
      "duration_minutes": 15
    }
  ],
  "more": false
}
```

---

## Coordination Documents (Fase 3)

### `POST /api/v1/coordination-documents`

Crear documento de coordinacion (wizard completo).

**Roles:** coordinator

**Request:**
```json
{
  "name": "Planificacion anual Ciencias 2026",
  "area_id": 1,
  "start_date": "2026-03-01",
  "end_date": "2026-11-30",
  "topic_ids": [5, 8, 12],
  "subjects": [
    {
      "subject_id": 1,
      "class_count": 20,
      "topic_ids": [5, 8]
    },
    {
      "subject_id": 2,
      "class_count": 15,
      "topic_ids": [8, 12]
    }
  ]
}
```

- `topic_ids` (root): topics seleccionados al nivel `topic_selection_level`
- `subjects[].topic_ids`: subconjunto asignado a cada disciplina

**Response `201`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "name": "Planificacion anual Ciencias 2026",
  "area_id": 1,
  "start_date": "2026-03-01",
  "end_date": "2026-11-30",
  "status": "draft",
  "sections": {},
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z",
  "topics": [
    { "id": 5, "name": "Suma y resta", "level": 3 },
    { "id": 8, "name": "Ecuaciones", "level": 3 },
    { "id": 12, "name": "Cinematica", "level": 3 }
  ],
  "subjects": [
    {
      "id": 1,
      "subject_id": 1,
      "subject_name": "Matematicas",
      "class_count": 20,
      "topics": [
        { "id": 5, "name": "Suma y resta" },
        { "id": 8, "name": "Ecuaciones" }
      ],
      "classes": []
    }
  ]
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (area, topics, subjects), `TOPIC_WRONG_LEVEL`, `TOPIC_NOT_IN_ORG`, `SUBJECT_NOT_IN_AREA`, `TOPICS_NOT_FULLY_DISTRIBUTED`

---

### `GET /api/v1/coordination-documents`

Listar documentos de coordinacion.

**Roles:** coordinator, teacher

**Query params:** `limit`, `offset`, `area_id` (opcional), `status` (opcional: `pending`, `in_progress`, `published`)

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "name": "Planificacion anual Ciencias 2026",
      "area_id": 1,
      "area_name": "Ciencias Exactas",
      "start_date": "2026-03-01",
      "end_date": "2026-11-30",
      "status": "draft",
      "created_at": "2026-03-25T10:00:00Z",
      "updated_at": "2026-03-25T10:00:00Z"
    }
  ],
  "more": false
}
```

> El listado es liviano (sin sections, topics ni subjects). Usar GET /:id para detalle completo.

---

### `GET /api/v1/coordination-documents/:id`

Detalle completo del documento con todas las relaciones.

**Roles:** coordinator, teacher

**Response `200`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "name": "Planificacion anual Ciencias 2026",
  "area_id": 1,
  "area_name": "Ciencias Exactas",
  "start_date": "2026-03-01",
  "end_date": "2026-11-30",
  "status": "draft",
  "sections": {
    "problem_edge": {
      "value": "Como las logicas de poder y saber configuran..."
    },
    "methodological_strategy": {
      "selected_option": "proyecto",
      "value": "Implementaremos un ateneo-debate interdisciplinario..."
    }
  },
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z",
  "topics": [
    { "id": 5, "name": "Suma y resta", "level": 3 }
  ],
  "subjects": [
    {
      "id": 1,
      "coord_doc_subject_id": 1,
      "subject_id": 1,
      "subject_name": "Matematicas",
      "class_count": 20,
      "topics": [
        { "id": 5, "name": "Suma y resta" }
      ],
      "classes": [
        {
          "id": 1,
          "class_number": 1,
          "title": "Introduccion a la suma",
          "objective": "El alumno comprende la operacion de suma",
          "topics": [
            { "id": 5, "name": "Suma y resta" }
          ],
          "is_shared": false
        },
        {
          "id": 3,
          "class_number": 3,
          "title": "Practica compartida",
          "objective": "Practica interdisciplinaria",
          "topics": [
            { "id": 5, "name": "Suma y resta" }
          ],
          "is_shared": true
        }
      ]
    }
  ],
  "org_config": {
    "coord_doc_sections": [
      { "key": "problem_edge", "label": "Eje problematico", "type": "text", "required": true },
      { "key": "methodological_strategy", "label": "Estrategia metodologica", "type": "select_text", "options": ["proyecto", "taller_laboratorio", "ateneo_debate"], "required": true }
    ]
  }
}
```

> `org_config` incluye las secciones disponibles para que el frontend las renderice dinamicamente. `is_shared` indica si esa clase cae en un time slot compartido.

**Errores:** `NOT_FOUND`

---

### `PATCH /api/v1/coordination-documents/:id`

Actualizar documento (sections, status, nombre, fechas).

**Roles:** coordinator

**Request — actualizar secciones:**
```json
{
  "sections": {
    "problem_edge": {
      "value": "Nuevo texto del eje problematico..."
    },
    "methodological_strategy": {
      "selected_option": "taller_laboratorio",
      "value": "Implementaremos talleres de laboratorio..."
    }
  }
}
```

**Request — actualizar nombre y fechas:**
```json
{
  "name": "Nuevo nombre",
  "start_date": "2026-04-01",
  "end_date": "2026-12-15"
}
```

Todos los campos son opcionales. Solo se actualizan los enviados.

**Response `200`:** Documento actualizado (mismo schema que GET /:id).

**Errores:** `NOT_FOUND`, `DOCUMENT_NOT_DRAFT`, `INVALID_SECTION_KEY`, `VALIDATION_ERROR`

---

### `DELETE /api/v1/coordination-documents/:id`

Eliminar documento. Solo documentos en estado `draft`.

**Roles:** coordinator

**Response `204`:** Sin body.

**Errores:** `NOT_FOUND`, `DOCUMENT_NOT_DRAFT`

---

### `POST /api/v1/coordination-documents/:id/generate`

Generar contenido con IA: secciones + plan de clases por disciplina.

**Roles:** coordinator

**Request (opcional — regenerar solo algunas secciones):**
```json
{
  "section_keys": ["problem_edge", "methodological_strategy"],
  "regenerate_class_plans": true
}
```

Si no se envia body, genera todo.

**Response `200`:**
```json
{
  "sections_generated": ["problem_edge", "methodological_strategy"],
  "class_plans_generated": [
    { "subject_id": 1, "subject_name": "Matematicas", "classes_count": 20 },
    { "subject_id": 2, "subject_name": "Fisica", "classes_count": 15 }
  ]
}
```

> Operacion asincrona en el futuro. Por ahora sincrona (< 30s esperado).

**Errores:** `NOT_FOUND`, `DOCUMENT_NOT_DRAFT`, `AI_GENERATION_ERROR`, `DOCUMENT_NO_TOPICS`

---

### `POST /api/v1/coordination-documents/:id/chat`

Chat con Alizia. Soporta function calling para modificar el documento.

**Roles:** coordinator

**Request:**
```json
{
  "message": "Cambia el titulo de la clase 3 de Matematicas a 'Practica de ecuaciones'",
  "history": [
    { "role": "user", "content": "Genera un eje problematico" },
    { "role": "assistant", "content": "He generado el siguiente eje..." }
  ]
}
```

- `history`: historial completo de la conversacion (el frontend lo mantiene)
- `message`: mensaje nuevo del usuario

**Response `200`:**
```json
{
  "role": "assistant",
  "content": "Listo! He actualizado el titulo de la clase 3 de Matematicas a 'Practica de ecuaciones'.",
  "tool_calls": [
    {
      "function": "update_class_title",
      "arguments": {
        "subject_id": 1,
        "class_number": 3,
        "title": "Practica de ecuaciones"
      }
    }
  ],
  "document_updated": true
}
```

**Tools disponibles (function calling):**

| Function | Parametros | Descripcion |
|----------|-----------|-------------|
| `update_section` | `section_key`, `content` | Actualiza una seccion del documento |
| `update_class_title` | `subject_id`, `class_number`, `title` | Cambia titulo de una clase |
| `update_class_topics` | `subject_id`, `class_number`, `topic_ids` | Cambia topics de una clase |
| `update_document_title` | `title` | Cambia nombre del documento |

**Errores:** `NOT_FOUND`, `DOCUMENT_NOT_DRAFT`, `AI_GENERATION_ERROR`, `INVALID_TOOL_CALL`

---

## Teaching (Fase 5)

### `GET /api/v1/course-subjects/:id/lesson-plans`

Listar lesson plans del docente para un course_subject.

**Roles:** teacher

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "course_subject_id": 1,
      "coordination_document_id": 1,
      "class_number": 1,
      "title": "Introduccion a la suma",
      "status": "pending",
      "is_shared": false,
      "coord_class": {
        "title": "Introduccion a la suma",
        "objective": "El alumno comprende la operacion de suma",
        "topics": [
          { "id": 5, "name": "Suma y resta" }
        ]
      }
    },
    {
      "id": null,
      "course_subject_id": 1,
      "coordination_document_id": 1,
      "class_number": 3,
      "title": null,
      "status": "pending",
      "is_shared": true,
      "coord_class": {
        "title": "Practica compartida",
        "objective": "Practica interdisciplinaria",
        "topics": [
          { "id": 5, "name": "Suma y resta" }
        ]
      }
    }
  ],
  "more": false
}
```

> Devuelve todas las clases del plan de coordinacion. Si el docente ya creo un lesson plan, tiene `id` y datos propios. Si no, `id` es null y solo tiene la info del coord_class. `is_shared` marca clases compartidas.

**Errores:** `NOT_FOUND`, `NO_PUBLISHED_DOCUMENT`

---

### `POST /api/v1/lesson-plans`

Crear lesson plan para una clase.

**Roles:** teacher

**Request:**
```json
{
  "course_subject_id": 1,
  "coordination_document_id": 1,
  "class_number": 1,
  "title": "Mi plan para la clase 1",
  "objective": "Que los alumnos dominen la suma basica",
  "knowledge_content": "Contenido conceptual...",
  "didactic_strategies": "Estrategias a usar...",
  "class_format": "presencial",
  "topic_ids": [5],
  "resources_mode": "per_moment",
  "moments": {
    "apertura": {
      "activities": [1]
    },
    "desarrollo": {
      "activities": [3, 5]
    },
    "cierre": {
      "activities": [8]
    }
  },
  "fonts": {
    "global": [],
    "apertura": [10],
    "desarrollo": [10, 11],
    "cierre": []
  }
}
```

- `moments`: IDs de actividades por momento. Apertura y cierre deben tener exactamente 1. Desarrollo 1 a `config.desarrollo_max_activities`.
- `fonts`: si `resources_mode: "global"`, usar `{ "global": [10, 11] }`. Si `"per_moment"`, distribuir por momento.

**Response `201`:**
```json
{
  "id": 1,
  "course_subject_id": 1,
  "coordination_document_id": 1,
  "class_number": 1,
  "title": "Mi plan para la clase 1",
  "objective": "Que los alumnos dominen la suma basica",
  "knowledge_content": "Contenido conceptual...",
  "didactic_strategies": "Estrategias a usar...",
  "class_format": "presencial",
  "status": "pending",
  "resources_mode": "per_moment",
  "topics": [
    { "id": 5, "name": "Suma y resta" }
  ],
  "moments": {
    "apertura": {
      "activities": [1],
      "activityContent": {}
    },
    "desarrollo": {
      "activities": [3, 5],
      "activityContent": {}
    },
    "cierre": {
      "activities": [8],
      "activityContent": {}
    }
  },
  "fonts": [
    { "moment": "apertura", "font_id": 10, "font_name": "Guia de aritmetica" },
    { "moment": "desarrollo", "font_id": 10, "font_name": "Guia de aritmetica" },
    { "moment": "desarrollo", "font_id": 11, "font_name": "Video tutorial" }
  ],
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (course_subject, document, activities, fonts), `NO_PUBLISHED_DOCUMENT`, `INVALID_MOMENT_ACTIVITIES`, `TOPIC_NOT_IN_DOCUMENT`, `LESSON_PLAN_ALREADY_EXISTS`

---

### `PATCH /api/v1/lesson-plans/:id`

Actualizar lesson plan. Mismos campos que POST, todos opcionales.

**Roles:** teacher

**Request (ejemplo parcial):**
```json
{
  "title": "Titulo actualizado",
  "moments": {
    "apertura": { "activities": [2] },
    "desarrollo": { "activities": [3] },
    "cierre": { "activities": [8] }
  }
}
```

**Response `200`:** Lesson plan actualizado (mismo schema que POST response).

**Errores:** `NOT_FOUND`, `VALIDATION_ERROR`, `INVALID_MOMENT_ACTIVITIES`

---

### `POST /api/v1/lesson-plans/:id/generate-activity`

Generar contenido con IA para una actividad especifica.

**Roles:** teacher

**Request:**
```json
{
  "moment": "desarrollo",
  "activity_id": 3
}
```

**Response `200`:**
```json
{
  "moment": "desarrollo",
  "activity_id": 3,
  "content": "Texto generado por IA describiendo como implementar la actividad 'Trabajo en grupo' en el contexto de esta clase..."
}
```

> El contenido generado se guarda automaticamente en `moments.desarrollo.activityContent["3"]` del lesson plan.

**Errores:** `NOT_FOUND`, `ACTIVITY_NOT_IN_MOMENT`, `AI_GENERATION_ERROR`

---

### `PATCH /api/v1/lesson-plans/:id/status`

Cambiar estado del lesson plan.

**Roles:** teacher

**Request:**
```json
{
  "status": "in_progress"
}
```

**Response `200`:**
```json
{
  "id": 1,
  "status": "in_progress"
}
```

**Errores:** `NOT_FOUND`, `INVALID_STATUS_TRANSITION`, `LESSON_PLAN_INCOMPLETE`

---

## Resources (Fase 6)

### `GET /api/v1/resource-types`

Tipos de recurso disponibles para la org (publicos habilitados + privados de la org).

**Roles:** teacher

**Response `200`:**
```json
{
  "items": [
    {
      "id": 1,
      "key": "lecture_guide",
      "name": "Guia de lectura",
      "description": "Guia para acompanar la lectura de un texto",
      "requires_font": true,
      "prompt": "Genera una guia de lectura...",
      "output_schema": {
        "type": "object",
        "properties": {
          "title": { "type": "string" },
          "sections": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "heading": { "type": "string" },
                "content": { "type": "string" }
              }
            }
          }
        }
      },
      "is_custom": false
    }
  ],
  "more": false
}
```

> `is_custom` indica si la org tiene override de prompt/schema. `prompt` y `output_schema` ya resuelven el override (custom ?? default).

---

### `GET /api/v1/fonts`

Fuentes educativas del area.

**Roles:** Todos

**Query params:** `area_id` (requerido), `validated_only` (default `true`)

**Response `200`:**
```json
{
  "items": [
    {
      "id": 10,
      "name": "Guia de aritmetica basica",
      "description": "PDF con ejercicios de suma y resta",
      "file_url": "https://storage.example.com/fonts/aritmetica.pdf",
      "file_type": "pdf",
      "thumbnail_url": "https://storage.example.com/fonts/aritmetica_thumb.png",
      "area_id": 1,
      "is_validated": true,
      "created_at": "2026-03-25T10:00:00Z"
    }
  ],
  "more": false
}
```

---

### `POST /api/v1/resources`

Crear recurso.

**Roles:** teacher

**Request:**
```json
{
  "resource_type_id": 1,
  "title": "Guia de lectura: Numeros naturales",
  "font_id": 10,
  "course_subject_id": 1
}
```

- `font_id`: requerido si el resource_type tiene `requires_font: true`
- `course_subject_id`: opcional, contexto de creacion

**Response `201`:**
```json
{
  "id": 1,
  "organization_id": 1,
  "resource_type_id": 1,
  "resource_type_name": "Guia de lectura",
  "title": "Guia de lectura: Numeros naturales",
  "content": {},
  "user_id": 5,
  "font_id": 10,
  "course_subject_id": 1,
  "status": "draft",
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z"
}
```

**Errores:** `VALIDATION_ERROR`, `NOT_FOUND` (resource_type, font, course_subject), `RESOURCE_TYPE_REQUIRES_FONT`, `RESOURCE_TYPE_DISABLED`

---

### `PATCH /api/v1/resources/:id`

Actualizar recurso (titulo, content, status).

**Roles:** teacher

**Request:**
```json
{
  "title": "Titulo actualizado",
  "content": {
    "title": "Guia de lectura: Numeros naturales",
    "sections": [
      { "heading": "Introduccion", "content": "Texto editado manualmente..." }
    ]
  }
}
```

**Response `200`:** Recurso actualizado (mismo schema que POST response).

**Errores:** `NOT_FOUND`, `VALIDATION_ERROR`, `CONTENT_SCHEMA_MISMATCH`

---

### `POST /api/v1/resources/:id/generate`

Generar contenido del recurso con IA.

**Roles:** teacher

**Request (opcional):**
```json
{
  "custom_instruction": "Enfocate en ejercicios practicos, nivel basico"
}
```

**Response `200`:**
```json
{
  "content": {
    "title": "Guia de lectura: Numeros naturales",
    "sections": [
      {
        "heading": "Antes de leer",
        "content": "Texto generado por IA..."
      },
      {
        "heading": "Durante la lectura",
        "content": "Texto generado por IA..."
      },
      {
        "heading": "Despues de leer",
        "content": "Texto generado por IA..."
      }
    ]
  }
}
```

> El contenido generado se guarda automaticamente en el recurso. La estructura sigue el `output_schema` del resource_type.

**Errores:** `NOT_FOUND`, `RESOURCE_NOT_DRAFT`, `AI_GENERATION_ERROR`, `RESOURCE_NO_FONT` (si requires_font y no tiene)

---

## AI — Chat general (Fase 4)

### `POST /api/v1/chat`

Chat general con Alizia (sin contexto de documento especifico).

**Roles:** Todos

**Request:**
```json
{
  "message": "Que estrategias metodologicas recomiendas para ensenar fracciones?",
  "history": [
    { "role": "user", "content": "Hola" },
    { "role": "assistant", "content": "Hola! Soy Alizia..." }
  ]
}
```

**Response `200`:**
```json
{
  "role": "assistant",
  "content": "Para ensenar fracciones te recomiendo las siguientes estrategias..."
}
```

**Errores:** `AI_GENERATION_ERROR`

---

## Queries SQL de referencia

### Tipos de recurso disponibles para una org

```sql
SELECT rt.*, ort.custom_prompt, ort.custom_output_schema
FROM resource_types rt
LEFT JOIN organization_resource_types ort
    ON ort.resource_type_id = rt.id AND ort.organization_id = $1
WHERE rt.is_active = true
  AND (
    (rt.organization_id IS NULL AND COALESCE(ort.enabled, true) = true)
    OR rt.organization_id = $1
  );
```

### Detectar clases compartidas

```sql
SELECT ts.day_of_week, ts.start_time, ts.end_time,
       array_agg(cs.id) AS course_subject_ids
FROM time_slots ts
JOIN time_slot_subjects tss ON tss.time_slot_id = ts.id
JOIN course_subjects cs ON cs.id = tss.course_subject_id
WHERE ts.course_id = $1
GROUP BY ts.id
HAVING count(*) > 1;
```

### Clases con topics de un documento

```sql
SELECT cdc.class_number, cdc.title, array_agg(t.name) AS topics
FROM coord_doc_classes cdc
JOIN coordination_document_subjects cds ON cds.id = cdc.coord_doc_subject_id
LEFT JOIN coord_doc_class_topics cdct ON cdct.coord_doc_class_id = cdc.id
LEFT JOIN topics t ON t.id = cdct.topic_id
WHERE cds.coordination_document_id = $1 AND cds.subject_id = $2
GROUP BY cdc.id ORDER BY cdc.class_number;
```
