# Épica 4: Documento de coordinación

> Creación asistida, edición colaborativa y gestión de estados de documentos de coordinación areal.

**Estado:** MVP
**Fase de implementación:** Fase 3

---

## Problema

Los equipos de docentes necesitan documentos complejos (como el itinerario del área) que deben alinear múltiples personas en el trabajo a realizarse en un plazo futuro dado. Este proceso es manual, lento y difícil de articular entre roles. Sin herramienta:

- El coordinador arma planillas sueltas sin estructura ni trazabilidad
- Los docentes no tienen visibilidad de lo que se planificó para su materia
- No hay forma de generar contenido pedagógico alineado al diseño curricular
- Las clases compartidas entre materias no se coordinan

## Objetivos

- Crear documentos de coordinación via wizard guiado (3 pasos)
- Soportar secciones dinámicas configurables por org (eje problemático, estrategia, criterios)
- Generar plan de clases por materia con IA (título, objetivo, topics por clase)
- Permitir edición directa y asistida por IA (chat con Alizia + function calling)
- Gestionar estados del documento (draft → published → archived)
- Soportar clases compartidas como diferenciador clave

## Alcance MVP

**Incluye:**

- Wizard de creación en 3 pasos (topics → período + class_count → asignar topics a materias)
- CRUD completo del documento con 6 tablas normalizadas
- Secciones dinámicas según `config.coord_doc_sections` (JSONB)
- Generación de secciones con IA (eje problemático, estrategia metodológica, etc.)
- Plan de clases por materia generado por IA
- Edición directa de secciones y clases
- Chat con Alizia (function calling: update_section, update_class, etc.)
- Publicación del documento → visible para docentes

**No incluye:**

- Planificación del clase a clase detallado → ver [Épica 5](../05-planificacion-docente/05-planificacion-docente.md)
- Creación de recursos didácticos → ver [Épica 8](../08-contenido-recursos/08-contenido-recursos.md)
- Motor de IA (prompts, Azure OpenAI) → ver [Épica 6](../06-assistente-ia/06-asistente-ia.md)

## Principios de diseño

- **Propuesta primero:** Alizia genera una primera versión; el equipo edita y valida.
- **Alineación vertical:** Todo lo que se planifica debe poder trazarse hasta los lineamientos provinciales (topics).
- **Colaboración entre roles:** Coordinadores definen el marco, docentes validan y ajustan.
- **Configurable, no hardcoded:** Las secciones del documento se definen en la config de la org.

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-4.1 | [Modelo de datos del documento](./HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) | 6 tablas normalizadas, entities, providers, repository | Fase 3 | 5 |
| HU-4.2 | [Wizard de creación](./HU-4.2-wizard-creacion/HU-4.2-wizard-creacion.md) | 3 pasos: topics, período + class_count, asignar topics a materias | Fase 3 | 5 |
| HU-4.3 | [Secciones dinámicas](./HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) | JSONB sections según config org, edición, generación IA | Fase 3 | 4 |
| HU-4.4 | [Plan de clases por materia](./HU-4.4-plan-clases-por-materia/HU-4.4-plan-clases-por-materia.md) | coord_doc_classes, generación IA, class_topics | Fase 3 | 5 |
| HU-4.5 | [Publicación y estados](./HU-4.5-publicacion-estados/HU-4.5-publicacion-estados.md) | draft → published → archived, validaciones al publicar | Fase 3 | 3 |
| HU-4.6 | [Chat con Alizia](./HU-4.6-chat-alizia/HU-4.6-chat-alizia.md) | Function calling para editar secciones y clases via chat | Fase 3 | 4 |

---

## Decisiones técnicas

- El documento usa **6 tablas normalizadas** en vez de JSONB anidados (lección del POC). Esto permite JOINs, FKs reales y validación por BD.
- Las **secciones son dinámicas** — definidas en `config.coord_doc_sections`. Cada org elige qué secciones tiene, sus labels, tipos de input y prompts de IA.
- El **período** es un nombre libre con fechas custom — no se fuerza semestre/cuatrimestre.
- El `class_count` se calcula automáticamente desde la grilla horaria pero el **coordinador puede override** (± feriados).
- Los **topics se seleccionan al nivel** definido por `config.topic_selection_level`. En el wizard paso 3, se distribuyen entre materias.
- Al **publicar**, se valida que todos los topics del documento estén distribuidos en al menos una materia.
- **Clases compartidas** se detectan automáticamente desde time_slots y se muestran en el plan de clases.
- El chat con Alizia usa **function calling** — tools genéricos como `update_section(key, content)` y `update_class(class_id, title, topics)`.

## Decisiones de cada cliente

- Qué secciones tiene el documento (eje problemático, estrategia, criterios de evaluación, etc.)
- Las estrategias metodológicas disponibles (proyecto, taller, ateneo, laboratorio)
- El nivel de edición que tiene el docente sobre el documento del coordinador
- Si los topics pueden repetirse entre materias o son exclusivos

## Épicas relacionadas

- **[Épica 3: Integración](../03-integracion/03-integracion.md)** — Provee áreas, materias, topics, time_slots que alimentan el wizard
- **[Épica 5: Planificación docente](../05-planificacion-docente/05-planificacion-docente.md)** — Consume el documento publicado como base para lesson plans
- **[Épica 6: Asistente IA](../06-assistente-ia/06-asistente-ia.md)** — Motor de generación de secciones, plan de clases y chat
- **[Épica 1: Roles y accesos](../01-roles-accesos/01-roles-accesos.md)** — RequireRole(coordinator) para crear/editar, teacher para leer

## Test cases asociados

- Wizard: crear documento con 3 pasos → documento en draft con topics y subjects
- Secciones: PATCH con section_key inválida → 422
- Generación: POST generate → secciones + plan de clases generados
- Publicación: publicar sin todos los topics distribuidos → 422
- Publicación: publicar ok → estado published, visible para teachers
- Chat: function calling update_section → sección actualizada
- Delete: eliminar documento published → 403 (solo draft)
