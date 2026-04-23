# Épica 4: Documento de coordinación

> Creación asistida, edición colaborativa y gestión de estados de documentos de coordinación areal.

**Estado:** MVP
**Fase de implementación:** Fase 3

---

## Problema

Los equipos de docentes necesitan documentos complejos (como el itinerario del área) que deben alinear múltiples personas en el trabajo a realizarse en un plazo futuro dado. Este proceso es manual, lento y difícil de articular entre roles. Sin herramienta:

- El coordinador arma planillas sueltas sin estructura ni trazabilidad
- Los docentes no tienen visibilidad de lo que se planificó para su disciplina
- No hay forma de generar contenido pedagógico alineado al diseño curricular
- Las clases compartidas entre disciplinas no se coordinan

## Objetivos

- Crear documentos de coordinación via wizard guiado (3 pasos)
- Soportar secciones dinámicas configurables por org (eje problemático, estrategia, criterios)
- Generar plan de clases por disciplina con IA (título, objetivo, topics por clase)
- Permitir edición directa y asistida por IA (chat con Alizia + function calling)
- Gestionar estados del documento (pending → in_progress → published / archived)
- Soportar clases compartidas como diferenciador clave

## Alcance MVP

**Incluye:**

- Wizard de creación en 3 pasos (topics → período + class_count → asignar topics a disciplinas)
- CRUD completo del documento con 6 tablas normalizadas + tablas de chat y snapshots
- Secciones dinámicas según `config.coord_doc_sections` (JSONB)
- Generación inicial de secciones con IA, con variantes para `select_text` (el coordinador elige entre opciones generadas o escribe una propia)
- Plan de clases por disciplina generado por IA, con coordinación de clases compartidas
- Edición directa de secciones; clases editables solo en estado `in_progress`
- Chat con Alizia (function calling: update_section, update_class, etc.) — sin botón "Regenerar" separado, toda la re-generación va por chat
- Historial de chat persistido en backend con auto-compactación
- Publicación del documento con validación (warning no bloqueante para sub-topics sin asignar)
- Documentos publicados: secciones editables con warning, clases inmutables
- Sistema de snapshots (hasta 3 versiones) para restaurar estados anteriores
- Auto-guardado (detalles de implementación TBD)

**No incluye:**

- Propagación automática de cambios a lesson plans existentes → post-MVP
- Templates de documentos (admin pre-arma esqueleto, coordinador instancia) → horizonte
- Planificación del clase a clase detallado → ver [Épica 5](../05-planificacion-docente/05-planificacion-docente.md)
- Creación de recursos didácticos → ver [Épica 8](../08-contenido-recursos/08-contenido-recursos.md)
- Motor de IA (prompts, Azure OpenAI) → ver [Épica 6](../06-asistente-ia/06-asistente-ia.md)

## Principios de diseño

- **Propuesta primero:** Alizia genera una primera versión; el equipo edita y valida.
- **Alineación vertical:** Todo lo que se planifica debe poder trazarse hasta los lineamientos provinciales (topics).
- **Colaboración entre roles:** Coordinadores definen el marco, docentes validan y ajustan.
- **Configurable, no hardcoded:** Las secciones del documento se definen en la config de la org.
- **Chat como canal principal de edición IA:** No hay botón "Regenerar" — toda la interacción con IA post-generación inicial pasa por el chat con Alizia.

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-4.1 | [Modelo de datos del documento](./HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) | 6 tablas normalizadas + chat_messages + snapshots, entities, providers, repository | Fase 3 | 5 |
| HU-4.2 | [Wizard de creación](./HU-4.2-wizard-creacion/HU-4.2-wizard-creacion.md) | 3 pasos: topics, período + class_count, asignar topics a disciplinas | Fase 3 | 5 |
| HU-4.3 | [Secciones dinámicas](./HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) | JSONB sections según config org, edición, generación IA con variantes | Fase 3 | 4 |
| HU-4.4 | [Plan de clases por disciplina](./HU-4.4-plan-clases-por-materia/HU-4.4-plan-clases-por-materia.md) | coord_doc_classes, generación IA coordinada, class_topics | Fase 3 | 5 |
| HU-4.5 | [Publicación y estados](./HU-4.5-publicacion-estados/HU-4.5-publicacion-estados.md) | pending → in_progress → published / archived, validaciones, secciones editables post-publish | Fase 3 | 3 |
| HU-4.6 | [Chat con Alizia](./HU-4.6-chat-alizia/HU-4.6-chat-alizia.md) | Function calling para editar secciones y clases, historial persistido, auto-compact | Fase 3 | 4 |

---

## Decisiones técnicas

- El documento usa **6 tablas normalizadas** en vez de JSONB anidados (lección del POC). Esto permite JOINs, FKs reales y validación por BD.
- Las **secciones son dinámicas** — definidas en `config.coord_doc_sections`. Cada org elige qué secciones tiene, sus labels, tipos de input y prompts de IA.
- El **período** es un nombre libre con fechas custom — no se fuerza semestre/cuatrimestre.
- El `class_count` se calcula automáticamente desde la grilla horaria pero el **coordinador puede override** (± feriados). Si no hay grilla cargada, el coordinador ingresa manualmente. **Mínimo: 1.**
- Los **topics se seleccionan al nivel** definido por `config.topic_selection_level`. En el wizard paso 3, se distribuyen entre disciplinas. **Los topics pueden repetirse entre disciplinas** — la interdisciplinariedad es un pilar del modelo.
- Al **publicar**, se valida que cada sub-topic esté asignado a al menos una materia. Si hay sub-topics sin asignar, se muestra un **warning confirmable** (no bloqueante): "Te quedó tal sub-topic sin asignar, ¿querés continuar?"
- **Clases compartidas**: la IA coordina el contenido de clases compartidas entre disciplinas, generando secuencialmente y pasando contexto de una materia a la siguiente.
- El chat con Alizia usa **function calling** — tools genéricos como `update_section(key, content)` y `update_class(class_id, title, topics)`. Las secciones varían por org (JSON Schema), los tools son genéricos y usan **JSON Path** para indicar qué parte del documento modificar.
- **No hay botón "Regenerar"** separado. Toda la edición asistida por IA post-generación inicial va por el chat con Alizia. El coordinador le pide en lenguaje natural qué cambiar.
- El **historial de chat se persiste en backend** (tabla `coord_doc_chat_messages`). Cuando la conversación es muy larga, se aplica auto-compactación con un límite configurable.
- **Sistema de snapshots**: antes de cada generación IA, se guarda un snapshot del estado actual (secciones + clases). Se mantienen hasta **3 snapshots** por documento para poder restaurar.
- **Auto-guardado**: el documento se guarda automáticamente (detalles de implementación TBD: frecuencia, mecanismo de polling/push).

### Documento publicado: secciones editables, clases inmutables (P12)

Un documento publicado **es parcialmente editable**:
- **Secciones narrativas** (eje problemático, estrategia, criterios): editables post-publicación con warning "Los cambios no se propagan a planificaciones docentes ya creadas"
- **Plan de clases** (títulos, objetivos, topics por clase): **inmutables** post-publicación. Solo el teacher puede editar clases en su planificación
- Alizia en el chat **no ofrece tools de edición de clases** cuando el documento está publicado

### Estados del documento

```
[pending] ──(primera edición/interacción)──→ [in_progress] ──(publicar)──→ [published]
                                                  │                              │
                                                  └──(archivar)──→ [archived] ←─┘
```

- **pending**: recién creado por el wizard, nunca editado. Transición **automática** a `in_progress` al primer click/edición/interacción
- **in_progress**: el coordinador está trabajando activamente
- **published**: visible para docentes. Secciones editables, clases inmutables
- **archived**: archivado manualmente. Sin restricción de cantidad por área — múltiples documentos pueden coexistir en cualquier estado

### Templates vs instancias (horizonte)

José propone separar `coordination_documents` en dos tablas: **templates** (esqueleto creado por admin/seed, inmutable, reutilizable entre cursos) e **instancias** (creadas por el coordinador a partir de un template, con state machine). **MVP: no se implementa.** Se evalúa post-MVP.

### Tabla de feriados (horizonte)

Para el cálculo automático de `class_count` se necesita una fuente de verdad de feriados. **MVP: se carga manualmente o se descuenta en el override del coordinador.**

## Decisiones de cada cliente

- Qué secciones tiene el documento (eje problemático, estrategia, criterios de evaluación, etc.)
- Las estrategias metodológicas disponibles (proyecto, taller, ateneo, laboratorio)
- El nivel de edición que tiene el docente sobre el documento del coordinador
- Los topics pueden repetirse entre disciplinas (confirmado: siempre permitido)

## Resumen de decisiones (preguntas abiertas)

| # | Pregunta | Decisión |
|---|----------|----------|
| P1 | Estado pending | Mantener con transición automática a in_progress al primera interacción |
| P2 | Múltiples docs por área | Sin restricción, agregar estado archived |
| P3 | Sin grilla horaria | Permitir carga manual, grilla es sugerencia |
| P4 | Topics compartidos | Sí, siempre permitido |
| P5 | Class count = 0 | No, mínimo 1 |
| P6 | select_text: cuándo elige | Generar N variantes + opción manual escrita por coordinador |
| P7 | Re-generación | Sin botón regenerar, todo va por chat con Alizia |
| P8 | Tipos de sección MVP | Solo text y select_text |
| P9 | IA genera count incorrecto | Reintentar 1 vez con prompt reforzado |
| P10 | CRUD manual de clases | Clases inmutables post-publicación, solo teacher edita |
| P11 | Clases compartidas en IA | IA coordina (genera secuencialmente con contexto) |
| P12 | Publicado editable | Secciones sí, clases no + warning |
| P13 | Validación al publicar | Sub-topics asignados a materia: warning confirmable |
| P14 | Docentes ven no-publicados | No |
| P15 | Persistir chat | Sí, en backend + auto-compact |
| P16 | Tools del chat | Set mínimo (4 tools), sin tools de clases cuando published |
| P17 | Falla de IA | 1 reintento + timeout reducido |
| P18 | Tool calls atómicos | Best-effort + retry |
| P19 | Edición concurrente | Auto-guardado (implementación TBD) |
| P20 | Undo/versionado | Snapshots, máximo 3 versiones |

## Épicas relacionadas

- **[Épica 3: Integración](../03-integracion/03-integracion.md)** — Provee áreas, disciplinas, topics, time_slots que alimentan el wizard
- **[Épica 5: Planificación docente](../05-planificacion-docente/05-planificacion-docente.md)** — Consume el documento publicado como base para lesson plans
- **[Épica 6: Asistente IA](../06-asistente-ia/06-asistente-ia.md)** — Motor de generación de secciones, plan de clases y chat
- **[Épica 1: Roles y accesos](../01-roles-accesos/01-roles-accesos.md)** — RequireRole(coordinator) para crear/editar, teacher para leer

## Test cases asociados

- Wizard: crear documento con 3 pasos → documento en pending con topics y subjects
- Auto-transición: primera edición en documento pending → in_progress automático
- Secciones: PATCH con section_key inválida → 422
- Variantes select_text: POST generate → N variantes generadas para cada sección select_text
- Generación: POST generate → secciones + plan de clases generados (con coordinación de compartidas)
- Publicación: publicar con sub-topics sin asignar → warning confirmable (no bloqueante)
- Publicación: publicar ok → estado published, visible para teachers
- Chat: function calling update_section → sección actualizada
- Chat published: tools de clases no disponibles cuando documento está publicado
- Chat historial: cerrar y reabrir → historial recuperado del backend
- Delete: eliminar documento pending → 200
- Delete: eliminar documento published → 403
- Archivar: archivar documento → estado archived
- Edit published secciones: PATCH sección en published → 200 con warning
- Edit published clases: PATCH clase en published → 403
- Snapshots: generar → snapshot creado, restaurar → estado anterior recuperado
