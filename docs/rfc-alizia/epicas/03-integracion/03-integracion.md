# Épica 3: Integración

> Importación de datos y configuración de la estructura curricular de cada provincia.

**Estado:** MVP
**Fase de implementación:** Fase 2

---

## Problema

Cada provincia tiene su propio diseño curricular, terminología y estructura de conocimientos. El sistema necesita incorporar estos datos como base para que todo lo que Alizia genere esté alineado con la realidad de cada jurisdicción. Sin estos datos no se pueden crear documentos de coordinación ni planificaciones.

## Objetivos

- Modelar organizaciones como tenants con configuración propia (JSONB)
- Crear y gestionar áreas, materias y su relación con coordinadores
- Modelar la jerarquía curricular (topics) de forma dinámica y configurable por org
- Gestionar cursos, alumnos y asignaciones docente-materia-curso
- Definir grillas horarias semanales con soporte para clases compartidas
- Cargar actividades didácticas por momento de clase

## Alcance MVP

**Incluye:**

- CRUD de organizaciones con config JSONB
- CRUD de áreas, materias, asignación de coordinadores
- Topics como tabla auto-referencial con niveles dinámicos
- CRUD de cursos, students, course_subjects
- Time slots + time_slot_subjects (clases compartidas)
- Activities por momento (apertura, desarrollo, cierre)
- Todos los endpoints admin de Fase 2

**No incluye:**

- UI de administración → admin opera por API/scripts
- Importación masiva desde archivos (CSV, Excel) → horizonte
- Onboarding de usuarios → ver [Épica 2: Onboarding](../02-onboarding/02-onboarding.md)
- Autenticación y roles → ver [Épica 1: Roles y accesos](../01-roles-accesos/01-roles-accesos.md)

## Principios de diseño

- **Provincial first:** Cada implementación respeta la estructura y terminología de la provincia.
- **Fuente de verdad única:** Los datos importados alimentan toda la plataforma.
- **Setup manual deliberado:** El equipo de implementación hace el setup inicial por provincia — no hay backoffice self-service en el MVP. Lo que vendemos es que hacemos el setup del diseño curricular.
- **Migraciones incrementales:** No se crea el DER completo upfront. Cada HU agrega sus tablas.

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-3.1 | [Organizaciones y configuración](./HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) | Tabla organizations, config JSONB, CRUD, seed | Fase 2 | 5 |
| HU-3.2 | [Áreas y materias](./HU-3.2-areas-materias/HU-3.2-areas-materias.md) | Tablas areas, subjects, area_coordinators, CRUD, endpoints | Fase 2 | 6 |
| HU-3.3 | [Topics — jerarquía curricular](./HU-3.3-topics-jerarquia-curricular/HU-3.3-topics-jerarquia-curricular.md) | Tabla self-referential, niveles dinámicos, tree endpoint | Fase 2 | 6 |
| HU-3.4 | [Cursos, alumnos y asignaciones](./HU-3.4-cursos-alumnos-asignaciones/HU-3.4-cursos-alumnos-asignaciones.md) | courses, students, course_subjects, CRUD | Fase 2 | 5 |
| HU-3.5 | [Grilla horaria y clases compartidas](./HU-3.5-grilla-horaria-clases-compartidas/HU-3.5-grilla-horaria-clases-compartidas.md) | time_slots, time_slot_subjects, triggers, shared classes | Fase 2 | 6 |
| HU-3.6 | [Actividades didácticas](./HU-3.6-actividades-didacticas/HU-3.6-actividades-didacticas.md) | activities por momento, CRUD, validación por config | Fase 2 | 4 |

---

## Decisiones técnicas

- La configuración de cada organización se almacena en `organizations.config` (JSONB) — mismo patrón que TUNI. Controla: niveles de topics, secciones del doc, clases compartidas, feature flags, etc.
- Los topics se modelan en una **tabla única auto-referencial**. La profundidad se **pre-computa** en el campo `level`. Si se mueve un topic, se re-computa solo la rama afectada.
- Cada org define la **profundidad máxima** de su jerarquía (`topic_max_levels`) y los **nombres por nivel** (`topic_level_names`).
- Las clases compartidas se modelan con **2 `time_slot_subjects`** por slot. Trigger valida que ambos `course_subjects` pertenezcan al mismo `course_id`.
- Un área puede contener una o más materias. Si la provincia no organiza por áreas, se crea un área genérica.
- El setup inicial de cada provincia lo hace el equipo de implementación — no hay backoffice self-service.
- Las migraciones son **incrementales**: cada HU agrega sus tablas, no se crea todo upfront.

## Decisiones de cada cliente

- Profundidad de la jerarquía de topics y nombres por nivel
- Si usan áreas o no
- Si habilitan clases compartidas
- Actividades didácticas disponibles por momento
- Config general de la org (feature flags, secciones del doc, etc.)

## Épicas relacionadas

- **[Épica 1: Roles y accesos](../01-roles-accesos/01-roles-accesos.md)** — Provee autenticación, roles y multi-tenancy. Organizations + users ya existen de Épica 1.
- **[Épica 4: Documento de coordinación](../04-documento-coordinacion/04-documento-coordinacion.md)** — Consume áreas, materias, topics, time_slots para crear documentos.
- **[Épica 5: Planificación docente](../05-planificacion-docente/05-planificacion-docente.md)** — Usa course_subjects, activities, topics para lesson plans.
- **[Épica 6: Asistente IA](../06-assistente-ia/06-asistente-ia.md)** — Usa topics y config de org como contexto para generación.

## Test cases asociados

- Fase 2: CRUD de áreas, materias, topics, cursos
- Topics: crear topic excediendo `topic_max_levels` → error
- Topics: mover topic y verificar re-cómputo de levels
- Time slots: crear clase compartida con `shared_classes_enabled = false` → error
- Time slots: crear clase compartida con subjects de distinto curso → error (trigger)
- Activities: crear actividad con momento inválido → error
- Multi-tenancy: listar datos solo filtra por org del JWT
