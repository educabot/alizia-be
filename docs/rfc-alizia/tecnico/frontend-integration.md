# Guia de integracion frontend — Alizia API

Guia para el equipo de frontend (React + TypeScript) sobre como consumir la API y manejar los patrones dinamicos del backend.

---

## Autenticacion

### Login

El frontend obtiene un JWT via `team-ai-toolkit/tokens`. El flujo es:

1. Usuario ingresa email + password
2. Frontend envia `POST /api/v1/auth/login` con `{ email, password }`
3. Backend valida credenciales y retorna JWT
4. Frontend guarda JWT en memoria (no localStorage por seguridad)
5. Todas las requests llevan header `Authorization: Bearer <JWT>`

### Claims del JWT

```typescript
interface JWTClaims {
  id: string;
  name: string;
  email: string;
  avatar: string;
  roles: ("teacher" | "coordinator" | "admin")[];
  // organization_id se extrae en el backend, no esta en el JWT visible
}
```

### Manejo de expiracion

- JWT tiene expiracion (configurable, default 24h)
- Si el backend retorna `401 UNAUTHORIZED`, el frontend redirige a login
- No hay refresh token en el MVP

---

## Paginacion

Todos los endpoints GET de listado usan paginacion con patron "has more".

### Request

```
GET /api/v1/areas?limit=20&offset=0
```

| Param | Default | Max | Descripcion |
|-------|---------|-----|-------------|
| `limit` | 20 | 100 | Items por pagina |
| `offset` | 0 | - | Items a saltar |

### Response

```typescript
interface PaginatedResponse<T> {
  items: T[];
  more: boolean; // true = hay mas items despues de estos
}
```

### Implementacion en frontend

```typescript
// Cargar siguiente pagina
const loadMore = async (currentOffset: number) => {
  const res = await api.get(`/areas?limit=20&offset=${currentOffset}`);
  setItems(prev => [...prev, ...res.items]);
  setHasMore(res.more);
  setOffset(currentOffset + res.items.length);
};
```

> El backend internamente hace `LIMIT N+1` y retorna N items. Si trajo N+1, `more = true`.

---

## Multi-tenancy (transparente)

El frontend **no necesita enviar** `organization_id` en ningun request. El backend lo extrae automaticamente del JWT.

- Listados: filtrados por la org del usuario
- Creacion: el `organization_id` se asigna automaticamente
- Acceso a recurso de otra org: retorna `404 NOT_FOUND` (no `403`)

---

## Secciones dinamicas del documento de coordinacion

Este es el patron mas complejo de la API. Las secciones del documento son **configurables por organizacion**.

### Paso 1: Obtener la configuracion

Al cargar un documento (`GET /coordination-documents/:id`), la response incluye `org_config`:

```json
{
  "org_config": {
    "coord_doc_sections": [
      {
        "key": "problem_edge",
        "label": "Eje problematico",
        "type": "text",
        "ai_prompt": "...",
        "required": true
      },
      {
        "key": "methodological_strategy",
        "label": "Estrategia metodologica",
        "type": "select_text",
        "options": ["proyecto", "taller_laboratorio", "ateneo_debate"],
        "ai_prompt": "...",
        "required": true
      },
      {
        "key": "eval_criteria",
        "label": "Criterios de evaluacion",
        "type": "text",
        "ai_prompt": "...",
        "required": false
      }
    ]
  }
}
```

### Paso 2: Renderizar segun type

```typescript
type SectionType = "text" | "select_text" | "markdown";

interface SectionConfig {
  key: string;
  label: string;
  type: SectionType;
  options?: string[];  // Solo para select_text
  ai_prompt: string;
  required: boolean;
}
```

| Type | UI Component | Descripcion |
|------|-------------|-------------|
| `text` | Textarea | Texto libre editable |
| `select_text` | Select + Textarea | Primero elige opcion del select, luego edita el texto generado/escrito |
| `markdown` | Editor markdown | Texto con formato (futuro) |

### Paso 3: Leer valores actuales

Los valores estan en `document.sections`:

```json
{
  "sections": {
    "problem_edge": {
      "value": "Texto del eje problematico..."
    },
    "methodological_strategy": {
      "selected_option": "proyecto",
      "value": "Texto de la estrategia..."
    }
  }
}
```

- Una seccion sin valor: `{}` o ausente en el objeto
- `selected_option` solo existe para `type: "select_text"`

### Paso 4: Guardar cambios

```
PATCH /api/v1/coordination-documents/:id
{
  "sections": {
    "problem_edge": {
      "value": "Nuevo texto..."
    }
  }
}
```

> Solo enviar las secciones que cambiaron. El backend hace merge, no reemplaza todo el objeto `sections`.

### Paso 5: Generar con IA

```
POST /api/v1/coordination-documents/:id/generate
{
  "section_keys": ["problem_edge"],  // Opcional: solo estas secciones
  "regenerate_class_plans": false     // Opcional: si regenerar el plan de clases
}
```

> Sin body genera todo. El boton "Generar con Alizia" puede generar todo o solo una seccion.

### Validacion al publicar

Antes de permitir "Publicar", el frontend puede validar:

```typescript
const canPublish = (doc: Document, config: SectionConfig[]): string[] => {
  const errors: string[] = [];

  // Secciones requeridas
  for (const section of config.filter(s => s.required)) {
    const value = doc.sections?.[section.key]?.value;
    if (!value || value.trim() === "") {
      errors.push(`Seccion "${section.label}" es requerida`);
    }
  }

  // Todas las materias deben tener plan de clases
  for (const subject of doc.subjects) {
    if (subject.classes.length === 0) {
      errors.push(`${subject.subject_name} no tiene plan de clases`);
    }
  }

  return errors;
};
```

---

## Momentos didacticos (Lesson Plans)

### Estructura de moments

```typescript
interface Moments {
  apertura: {
    activities: number[];       // Exactamente 1 activity ID
    activityContent: Record<string, string>;  // Key = activity_id as string
  };
  desarrollo: {
    activities: number[];       // 1 a desarrollo_max_activities (default 3)
    activityContent: Record<string, string>;
  };
  cierre: {
    activities: number[];       // Exactamente 1 activity ID
    activityContent: Record<string, string>;
  };
}
```

### Obtener actividades disponibles

```
GET /api/v1/activities?moment=apertura
GET /api/v1/activities?moment=desarrollo
GET /api/v1/activities?moment=cierre
```

El frontend debe filtrar por momento para mostrar solo actividades validas en cada seccion.

### Validacion en frontend

```typescript
const validateMoments = (moments: Moments, maxDesarrollo: number): string[] => {
  const errors: string[] = [];
  if (moments.apertura.activities.length !== 1) errors.push("Apertura debe tener 1 actividad");
  if (moments.desarrollo.activities.length < 1 || moments.desarrollo.activities.length > maxDesarrollo)
    errors.push(`Desarrollo debe tener 1 a ${maxDesarrollo} actividades`);
  if (moments.cierre.activities.length !== 1) errors.push("Cierre debe tener 1 actividad");
  return errors;
};
```

### Generar contenido IA por actividad

```
POST /api/v1/lesson-plans/:id/generate-activity
{ "moment": "desarrollo", "activity_id": 3 }
```

El contenido generado se guarda en `moments.desarrollo.activityContent["3"]`. El frontend debe refrescar el lesson plan despues de generar.

---

## Fuentes educativas (Fonts)

### Modo global vs por momento

El lesson plan tiene `resources_mode`:

| Mode | Comportamiento | UI |
|------|---------------|-----|
| `global` | Fonts aplican a toda la clase | Un unico selector de fonts |
| `per_moment` | Fonts asignadas por momento | Selector de fonts dentro de cada momento |

### Al crear/editar lesson plan

```json
// Modo global
{ "resources_mode": "global", "fonts": { "global": [10, 11] } }

// Modo per_moment
{ "resources_mode": "per_moment", "fonts": { "apertura": [10], "desarrollo": [10, 11], "cierre": [] } }
```

---

## Clases compartidas

### En la grilla horaria (course detail)

Un time_slot con 2 subjects es una clase compartida:

```typescript
const isShared = (slot: TimeSlot) => slot.subjects.length > 1;
```

Renderizar con badge o icono para indicar que dos materias se dictan en el mismo horario.

### En el plan de clases (coordination document)

Cada clase tiene `is_shared: boolean`. Las clases compartidas deben tener actividades interdisciplinarias.

### En lesson plans (teacher view)

El listado de lesson plans incluye `is_shared` por clase. Mostrar indicador visual para que el docente sepa que esa clase se comparte con otro docente.

---

## Chat con Alizia

### Mantener historial

El frontend mantiene el historial completo y lo envia en cada request:

```typescript
interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

// Estado
const [history, setHistory] = useState<ChatMessage[]>([]);

// Enviar mensaje
const sendMessage = async (message: string) => {
  const res = await api.post(`/coordination-documents/${docId}/chat`, {
    message,
    history
  });

  setHistory(prev => [
    ...prev,
    { role: "user", content: message },
    { role: "assistant", content: res.content }
  ]);
};
```

### Detectar cambios en el documento

La response del chat incluye `document_updated: boolean`. Si es `true`, el frontend debe refrescar el documento:

```typescript
if (res.document_updated) {
  const updatedDoc = await api.get(`/coordination-documents/${docId}`);
  setDocument(updatedDoc);
}
```

---

## Resources (Recursos)

### Flujo de creacion

```
1. GET /resource-types → mostrar tipos disponibles
2. Si tipo.requires_font → GET /fonts?area_id=X → mostrar selector de fonts
3. POST /resources → crear recurso (draft, content vacio)
4. POST /resources/:id/generate → generar contenido con IA
5. PATCH /resources/:id → editar content manualmente
```

### Renderizar content dinamico

El `content` de un recurso sigue el `output_schema` del tipo. El frontend debe renderizar dinamicamente:

```typescript
// output_schema del tipo "lecture_guide"
{
  "type": "object",
  "properties": {
    "title": { "type": "string" },
    "sections": {
      "type": "array",
      "items": {
        "properties": {
          "heading": { "type": "string" },
          "content": { "type": "string" }
        }
      }
    }
  }
}

// Render: iterar sections, mostrar heading como H3 y content como textarea editable
```

---

## Manejo de errores

Ver [errores.md](./errores.md) para el catalogo completo.

```typescript
interface APIError {
  error: {
    code: string;    // Usar para logica (switch/case)
    message: string; // Mostrar al usuario
    details?: any;   // Info adicional (ej: INVALID_MOMENT_ACTIVITIES)
  };
}

// Handler generico
const handleError = (err: APIError) => {
  switch (err.error.code) {
    case "UNAUTHORIZED":
      redirectToLogin();
      break;
    case "FORBIDDEN":
      toast.error("No tenes permisos para esta accion");
      break;
    case "AI_GENERATION_ERROR":
    case "AI_GENERATION_TIMEOUT":
    case "AI_RATE_LIMITED":
      toast.error("Error al generar con IA. Intenta nuevamente.");
      break;
    default:
      toast.error(err.error.message);
  }
};
```

---

## Tipos TypeScript de referencia

```typescript
// Coordination Document (response completa)
interface CoordinationDocument {
  id: number;
  organization_id: number;
  name: string;
  area_id: number;
  area_name: string;
  start_date: string;   // ISO date
  end_date: string;
  status: "draft" | "published" | "archived";
  sections: Record<string, SectionValue>;
  topics: Topic[];
  subjects: DocumentSubject[];
  org_config: { coord_doc_sections: SectionConfig[] };
  created_at: string;
  updated_at: string;
}

interface SectionValue {
  value?: string;
  selected_option?: string;
}

interface DocumentSubject {
  id: number;
  coord_doc_subject_id: number;
  subject_id: number;
  subject_name: string;
  class_count: number;
  topics: Topic[];
  classes: DocumentClass[];
}

interface DocumentClass {
  id: number;
  class_number: number;
  title: string;
  objective: string;
  topics: Topic[];
  is_shared: boolean;
}

interface Topic {
  id: number;
  name: string;
  level?: number;
}

// Lesson Plan
interface LessonPlan {
  id: number | null;
  course_subject_id: number;
  coordination_document_id: number;
  class_number: number;
  title: string | null;
  objective?: string;
  status: "pending" | "planned";
  is_shared: boolean;
  moments?: Moments;
  coord_class: {
    title: string;
    objective: string;
    topics: Topic[];
  };
}

// Resource
interface Resource {
  id: number;
  resource_type_id: number;
  resource_type_name: string;
  title: string;
  content: Record<string, any>;  // Dinamico segun output_schema
  status: "draft" | "active";
  font_id?: number;
  course_subject_id?: number;
}
```
