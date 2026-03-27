# Épica 8: Contenido y recursos

> Biblioteca de recursos didácticos y herramientas de creación basadas en fuentes oficiales.

**Estado:** MVP
**Fase de implementación:** Fase 6

---

## Problema

Los docentes necesitan recursos didácticos (fichas de cátedra, guías de lectura, imágenes, videos, etc.) adaptados a su contexto curricular. Crearlos desde cero es lento y recurrir a fuentes no curadas genera inconsistencias.

## Objetivos

- Permitir la creación de recursos a partir de fuentes oficiales validadas por el ministerio
- Ofrecer tipos de recurso predefinidos (ficha de cátedra, guía de lectura, entre otros) y también "Creación libre" donde el usuario define los lineamientos de lo que quiere
- Que los recursos creados sean utilizables dentro de las planificaciones de clase, para ser entregados antes, durante o pos clase

## Alcance MVP

**Incluye:**

- Creación de recursos a partir de fuentes oficiales provistas por el cliente
- Tipos de recurso predefinidos (ficha de cátedra, guía de lectura, más a definir con el ministerio)
- Edición directa y asistida por IA del recurso generado
- Permitir la exportación para impresión
- Disponibilizar para uso del recurso en clase (Planificación)

**No incluye:**

- Subida de fuentes propias del docente → decisión pendiente por provincia
- Biblioteca compartida entre docentes → horizonte

## Principios de diseño

- **Fuentes curadas:** Los recursos se generan desde fuentes oficiales, no desde internet abierto.
- **Listo para el aula:** El recurso generado debe ser usable directamente con los alumnos, incluyendo la información del aula.

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-8.1 | [Modelo de datos recursos](./HU-8.1-modelo-datos-recursos/HU-8.1-modelo-datos-recursos.md) | Tablas para fonts, tipos de recurso, recursos y config por org | Fase 6 | 4 |
| HU-8.2 | [Tipos de recurso y configuración](./HU-8.2-tipos-recurso-configuracion/HU-8.2-tipos-recurso-configuracion.md) | Gestión de tipos por org, resolución de prompt/schema | Fase 6 | 3 |
| HU-8.3 | [Generación con IA](./HU-8.3-generacion-ia-recursos/HU-8.3-generacion-ia-recursos.md) | Generar recursos a partir de tipo, fuente y contexto | Fase 6 | 4 |
| HU-8.4 | [Biblioteca y exploración](./HU-8.4-biblioteca-exploracion/HU-8.4-biblioteca-exploracion.md) | Listado, filtros y búsqueda de recursos de la org | Fase 6 | 3 |
| HU-8.5 | [Edición y exportación](./HU-8.5-edicion-exportacion/HU-8.5-edicion-exportacion.md) | Edición manual/asistida, publicación y exportación | Fase 6 | 4 |

---

## Sub-épicas

| Componente | Descripción |
|---|---|
| Biblioteca | Repositorio de recursos creados por el docente |
| Creación de contenido | Generación asistida de recursos a partir de fuentes y tipo seleccionado |

## Decisiones de cada cliente

- Los tipos de recurso disponibles se definen con cada equipo pedagógico provincial
- Si se permite o no que el docente suba fuentes propias es decisión de cada provincia

## Decisiones técnicas

- Cada organización **habilita qué tipos de recurso** tiene disponibles. No todos los clientes ven los mismos tipos — se activan según acuerdo con el equipo pedagógico provincial. Un tipo de recurso puede existir en el sistema y no estar habilitado para una organización.
- El concepto de **biblioteca** es central: los recursos creados se almacenan a nivel organización. Todos los docentes de la misma organización pueden explorar y reutilizar recursos existentes antes de generar uno nuevo. Esto reduce costos de generación y promueve consistencia.
- El filtro por materia opera como **restricción soft** (UX, no permisos). Un docente de matemáticas no ve recursos de ciencias naturales por default, pero a nivel permiso el acceso es por organización.
- Cada tipo de recurso tiene un **prompt y un JSON Schema** que define la estructura del output. Esto permite que una misma funcionalidad (ej: guía de lectura) genere resultados con formatos distintos según la provincia, modificando solo la configuración.
- Arrancar simple: en el MVP, si un segundo cliente necesita una variante de un tipo existente, se **duplica y adapta** en vez de sobre-ingenierizar un sistema de templates parametrizables. La genericidad se construye cuando aparezca el patrón real.

### Tipos públicos vs privados

| `organization_id` | Visibilidad | Ejemplo |
|---|---|---|
| `NULL` | Público: visible para todas las orgs | `lecture_guide`, `course_sheet` |
| Set | Privado: solo para esa org | Tipos custom |

### Flujo de generación IA

1. Docente elige tipo → (si `requires_font`) elige fuente → crea `resources` con `content = {}`
2. Resuelve prompt: `organization_resource_types.custom_prompt` ?? `resource_types.prompt`
3. Resuelve output_schema: `custom_output_schema` ?? `output_schema`
4. Envía al LLM con contexto (font, course_subject, etc.)
5. Respuesta se guarda en `resources.content` (JSONB) según schema
6. Frontend renderiza `content` dinámicamente según `output_schema`
7. Chat con Alizia puede editar secciones del `content`

## Épicas relacionadas

- **Planificación docente** — Los recursos se pueden incorporar en las clases
- **Asistente IA** — Motor de generación y edición de contenido
- **Integración** — Provee las fuentes oficiales del diseño curricular
