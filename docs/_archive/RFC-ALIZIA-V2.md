# RFC: Alizia v2 — Plataforma de planificación educativa anual

| Campo              | Valor                                      |
|--------------------|--------------------------------------------|
| **Autor(es)**      | Equipo Backend + Equipo de Producto        |
| **Estado**         | 🟡 Borrador                                |
| **Tipo**           | Épica / API / Refactor                     |
| **Creado**         | 2026-03-25                                 |
| **Última edición** | 2026-03-25                                 |
| **Revisores**      | Pendiente                                  |
| **Decisión**       | Pendiente                                  |

---

## Historial de versiones

| Versión | Fecha      | Autor   | Cambios |
|---------|------------|---------|---------|
| 0.1     | 2026-03-25 | Equipo Backend | Borrador inicial |
| 0.2     | 2026-03-25 | Equipo Backend | RFC completo — fusión de notion.md (producto) + proposal-der-v2 (datos) + arquitectura Go |

---

## Índice

- [Contexto y motivación](#contexto-y-motivación-)
- [Objetivo](#objetivo-)
- [Alcance](#alcance-)
- [Diseño de producto](#diseño-de-producto-)
- [Épicas](#épicas)
- [Patrones transversales](#patrones-transversales)
- [Arquitectura general](#arquitectura-general-️)
- [Backend — Endpoints](#backend--endpoints-️)
- [Backend — Modelo de datos](#backend--modelo-de-datos-)
- [Backend — Lógica y configuración](#backend--lógica-y-configuración-)
- [QA — Estrategia de testing](#qa--estrategia-de-testing-)
- [Alternativas evaluadas](#alternativas-evaluadas-)
- [Rollout](#rollout-)
- [Dependencias](#dependencias-)
- [Riesgos](#riesgos-)
- [Preguntas abiertas](#preguntas-abiertas-)
- [Glosario](#glosario-)
- [Tareas](#tareas-)

---

## Contexto y motivación 🚀

### Problema

Los coordinadores y docentes de instituciones educativas planifican anualmente de forma fragmentada: documentos sueltos, WhatsApp, planillas, sin alineación entre áreas. No existe una herramienta que:

1. Permita al coordinador definir un plan de área que los docentes hereden
2. Genere contenido pedagógico con IA alineado al diseño curricular provincial
3. Recolecte feedback de lo que pasó en clase para mejorar planificaciones futuras
4. Se adapte a la estructura curricular de cada provincia sin cambios de código

### Contexto

- **POC actual**: Backend en FastAPI (Python), frontend vanilla HTML/JS, 10 tablas, sin multi-tenancy
- **POC validó**: El flujo coordinador → documento de coordinación → docente → lesson plan funciona
- **v2 es un rewrite completo**: Backend en Go, frontend en React, 26+ tablas, multi-tenant, IA integrada
- **Dos productos comparten infraestructura**: Alizia + TiCh/Tuni usan Auth0 para autenticación y team-ai-toolkit como librería compartida
- **Arquitectura ya definida**: Ver ARQUITECTURA-GO-ALIZIA-V2.md

### Documentos relacionados

- [ARQUITECTURA-GO-ALIZIA-V2.md](./ARQUITECTURA-GO-ALIZIA-V2.md) — Arquitectura técnica del backend
- [proposal-der-v2.md](./proposal-der-v2.md) — Modelo de datos completo (26 tablas)
- [architecture-auth-service.md](./architecture-auth-service.md) — Auth service propio (plan futuro, no bloqueante)
- [BACK-CONFIG-LIBRERIA.md](./BACK-CONFIG-LIBRERIA.md) — Librería compartida (team-ai-toolkit)
- [notion.md](./notion.md) — Épicas de producto (10 épicas)

---

## Objetivo 🎯

### Objetivos

- ✅ Rewrite completo del backend en Go con Clean Architecture
- ✅ Multi-tenancy: cada provincia/institución es un tenant con configuración propia
- ✅ Documentos de coordinación con generación IA y edición colaborativa
- ✅ Planificación docente clase a clase heredada del documento de coordinación
- ✅ Sistema de recursos (guías, fichas) con generación IA configurable por org
- ✅ Clases compartidas (coordinadas) como diferenciador de producto
- ✅ Auth via Auth0 (mismo sistema que tich-cronos), con posibilidad de migrar a auth-service propio en el futuro

### No-objetivos

- ❌ Frontend (este RFC es solo backend)
- ❌ WhatsApp integration (Épica 9, pendiente de definición)
- ❌ Cosmos (Épica 10, sin definición)
- ❌ Social login (Google/Microsoft) — futuro
- ❌ Admin panel UI — administración por API
- ❌ Migración de datos del POC — se arranca de cero
- ❌ Roles de directivos o supervisores — horizonte
- ❌ Gestión de múltiples instituciones por usuario — por definir
- ❌ Informe de proceso de alumnos — horizonte
- ❌ Trayectorias de refuerzo personalizadas — horizonte

### Métricas de éxito

| Métrica | Valor esperado | Cómo se mide |
|---------|---------------|--------------|
| Documentos de coordinación creados | ≥1 por área piloto | Query DB |
| Lesson plans generados | ≥1 por docente piloto | Query DB |
| Tiempo de generación IA | < 30s por sección | Logs |
| Cobertura de tests | ≥ 80% | CI/CD coverage report |
| Uptime | > 99.5% | Railway healthcheck |

---

## Alcance 📋

### Incluye (MVP)

- Épica 1: Roles y accesos (login, roles, multi-org)
- Épica 3: Integración (diseño curricular, topics, horarios)
- Épica 4: Documento de coordinación (wizard, secciones, IA, publicación)
- Épica 5: Planificación docente (lesson plans, momentos, IA) — sin bitácora ni repropuesta
- Épica 6: Asistente IA (generación, chat, function calling)
- Épica 8: Contenido (recursos, tipos, generación IA, library)

### No incluye (futuro)

- Épica 2: Onboarding — NTH post-MVP
- Épica 7: Dashboard — NTH, depende de Épica 4 y 5
- Épica 9: WhatsApp — pendiente definición
- Épica 10: Cosmos — pendiente definición
- Bitácora de cotejo (audio) — parte de Épica 5, post-MVP
- Repropuesta automática — parte de Épica 5, post-MVP
- Export PDF — NTH

### Fases de implementación

| Fase | Qué incluye | Dependencia |
|------|-------------|-------------|
| 1 — Setup | Repo, CI/CD, Railway, DB, auth integration (Auth0), /health | team-ai-toolkit |
| 2 — Admin/Integration (Épica 1 + 3) | Orgs, areas, subjects, topics, courses, time slots | Fase 1 |
| 3 — Coordination Documents (Épica 4) | Wizard, secciones dinámicas, CRUD, publicación | Fase 2 |
| 4 — AI Generation (Épica 6) | Generación de secciones, plan de clases, chat | Fase 3 |
| 5 — Teaching (Épica 5) | Lesson plans, momentos, actividades, generación | Fase 3 |
| 6 — Resources (Épica 8) | Tipos de recurso, fonts, generación IA, library | Fase 2 |

---

## Diseño de producto 🧠

### Resumen del sistema

Sistema multi-tenant de planificación educativa anual. Cada organización (colegio, provincia) es un tenant con configuración propia via `organizations.config` (JSONB).

**Roles:**
- **Coordinador**: crea documentos de coordinación que definen qué se enseña por materia en un período (topics, plan de clases, secciones configurables por org)
- **Docente**: a partir del documento de coordinación, crea lesson plans clase a clase (momentos con actividades, fuentes) y recursos (guías, fichas, etc.)
- **Admin**: gestión de la organización, usuarios, configuración

**Flujo principal:**
1. Coordinador crea **coordination_document** para un área → selecciona topics, asigna a materias, genera plan de clases con IA
2. Docente ve el plan y crea **teacher_lesson_plans** por clase → selecciona actividades por momento (apertura/desarrollo/cierre), IA genera contenido por actividad
3. Docente crea **resources** (guías de lectura, fichas de curso, etc.) usando tipos configurables con generación IA

### Principios de diseño

1. **Provincial-first** — Cada implementación respeta la estructura y terminología de la provincia
2. **Proponer, no imponer** — La IA propone; el usuario decide
3. **Configurable, no customizable** — JSON config por org, no código custom por cliente
4. **Simple sobre abstracto** — Si un patrón emerge en 3+ clientes, entonces genericizar
5. **Rol define el flujo** — La experiencia del usuario cambia según su rol desde el inicio
6. **Del área al aula** — La planificación individual nace del acuerdo colectivo
7. **Fuentes curadas** — Los recursos se generan desde fuentes oficiales, no desde internet abierto
8. **IA que aprende del aula** — Las propuestas mejoran con el feedback real del docente
9. **Voz del docente** — La bitácora acepta audio libre, sin formato rígido

### Flujos de usuario

#### Flujo 1: Setup de la organización (Admin)

**Actor:** Admin / Equipo de implementación
**Precondición:** Organización creada en Auth0

1. Crear organización con config JSONB (niveles de topics, secciones, feature flags)
2. Crear usuarios y asignar roles (teacher, coordinator, admin)
3. Crear áreas y asignar coordinadores
4. Crear materias en cada área
5. Crear cursos y alumnos
6. Crear course_subjects (curso + materia + docente + período)
7. Definir time_slots (grilla horaria semanal)
8. Si clases compartidas habilitadas → 2 time_slot_subjects por slot
9. Cargar topics (jerarquía según topic_max_levels de la config)
10. Cargar activities (por momento: apertura/desarrollo/cierre)

**Resultado:** Organización lista para que coordinadores creen documentos.

#### Flujo 2: Coordinador crea documento de coordinación

**Actor:** Coordinador
**Precondición:** Área con materias y topics cargados

1. Coordinador selecciona área
2. **Wizard paso 1**: Selecciona topics al nivel configurado por `topic_selection_level`
3. **Wizard paso 2**: Define período (fechas custom) + cantidad de clases por materia
4. **Wizard paso 3**: Asigna topics a cada materia
5. Sistema crea documento en estado `draft`
6. Coordinador completa secciones dinámicas (según `config.coord_doc_sections`)
7. Opcionalmente genera secciones con IA ("Generar con Alizia")
8. IA genera: eje problemático, estrategia metodológica, plan de clases por materia
9. Coordinador revisa, edita directo o via chat con Alizia (function calling: `update_section`, `update_class`, etc.)
10. Coordinador publica → estado `published` → visible para docentes

**Resultado:** Documento publicado con secciones completas y plan de clases por materia.

#### Flujo 3: Docente planifica clase a clase

**Actor:** Docente
**Precondición:** Documento de coordinación publicado para su materia

1. Docente ve el plan de clases heredado del documento de coordinación
2. Selecciona una clase (class_number)
3. Crea teacher_lesson_plan: título, objetivo, topics
4. Selecciona actividades por momento:
   - **Apertura**: exactamente 1 actividad
   - **Desarrollo**: 1 a `config.desarrollo_max_activities` (default 3) actividades
   - **Cierre**: exactamente 1 actividad
5. Selecciona fuentes educativas (global o por momento)
6. Opcionalmente genera contenido por actividad con IA (`activityContent`)
7. Edita directo o via chat con Alizia
8. Estado cambia a `planned`

**Resultado:** Lesson plan con momentos, actividades y contenido generado.

#### Flujo 4: Docente crea recurso

**Actor:** Docente
**Precondición:** Tipos de recurso habilitados en la org

1. Docente elige tipo de recurso (guía de lectura, ficha de curso, etc.)
2. Si el tipo `requires_font` → selecciona font (fuente educativa)
3. Se resuelve el prompt: `organization_resource_types.custom_prompt` ?? `resource_types.prompt`
4. Se resuelve el output_schema: `custom_output_schema` ?? `output_schema`
5. Se envía al LLM con contexto (font, course_subject, etc.)
6. La respuesta se guarda en `resources.content` (JSONB) según el schema
7. Docente edita directo o via chat con Alizia
8. Recurso queda en library, reutilizable por otros docentes de la org

**Resultado:** Recurso generado, editable, exportable, reutilizable.

### Reglas de negocio

| # | Regla | Ejemplo | Aplica a |
|---|-------|---------|----------|
| 1 | Cada org define niveles de topics via config | `topic_max_levels: 3`, nombres: "Núcleos", "Áreas", "Categorías" | Back |
| 2 | Topics se seleccionan al nivel `topic_selection_level` | Si level=3, se eligen categorías, no núcleos | Back + Front |
| 3 | Clases compartidas solo si `shared_classes_enabled` | 2 materias en mismo time_slot, ambas del mismo área | Back |
| 4 | Secciones del doc son dinámicas según `coord_doc_sections` | Cada sección tiene key, label, type, ai_prompt, required | Back + Front |
| 5 | Momentos didácticos son enum fijo: apertura, desarrollo, cierre | desarrollo permite 1 a `desarrollo_max_activities` actividades | Back |
| 6 | Resource types pueden ser públicos (todas las orgs) o privados (1 org) | `organization_id IS NULL` = público | Back |
| 7 | Un usuario puede tener múltiples roles | teacher + coordinator en la misma org | Auth0 + Back |
| 8 | Mismo email puede existir en orgs distintas | `UNIQUE(email, organization_id)` | Auth0 + Back |
| 9 | El período del documento es texto libre con fechas custom | No se fuerza semestre/cuatrimestre | Back |
| 10 | Un docente por materia por curso (first-come-first-serve si hay conflicto) | El primero en planificar escribe | Back |
| 11 | Coordinador puede override manual de class_count (± feriados) | Ajuste manual sobre el cálculo automático | Back |
| 12 | Todos los topics del documento deben estar distribuidos entre materias | Validación al publicar | Back |
| 13 | Filter por materia en library es soft (UX), no permisos | Un docente de matemáticas puede ver recursos de ciencias | Front |
| 14 | Permisos sobre el doc de coordinación son configurables por org | Quién edita, quién solo visualiza | Back |

### Decisiones por provincia

| Decisión | Quién decide | Default | Impacto en config |
|----------|-------------|---------|-------------------|
| Niveles de topics (profundidad) | Provincia | 3 | `topic_max_levels`, `topic_level_names` |
| Nivel de selección de topics | Provincia | 3 | `topic_selection_level` |
| Clases compartidas habilitadas | Provincia | true | `shared_classes_enabled` |
| Secciones del documento | Provincia | problem_edge + methodological_strategy + eval_criteria | `coord_doc_sections` |
| Max actividades en desarrollo | Provincia | 3 | `desarrollo_max_activities` |
| Tipos de recurso habilitados | Provincia | Todos los públicos | `organization_resource_types` |
| Permisos del docente sobre el doc | Provincia | Solo lectura | Config de permisos |
| Datos requeridos en onboarding | Provincia | Perfil básico | Config de onboarding |
| Estrategias metodológicas disponibles | Provincia | proyecto, taller, ateneo | `coord_doc_sections[].options` |
| Tipos de actividad por momento | Provincia | Definidos con equipo pedagógico | Tabla `activities` |

### Estados y ciclo de vida

#### Documento de coordinación
```
[draft] ──(publicar / coordinador)──▶ [published] ──(archivar / coordinador)──▶ [archived]
```

#### Lesson plan
```
[pending] ──(planificar / docente)──▶ [planned]
```

#### Recurso
```
[draft] ──(activar / docente)──▶ [active]
```

---

## Épicas

### Épica 1: Roles y accesos

> Autenticación, roles, permisos y asignación organizacional de usuarios.

**Problema:** La plataforma opera con múltiples roles (coordinador, docente, y potencialmente directivos) que tienen permisos distintos sobre los mismos documentos y cursos. Se necesita un sistema que controle quién puede crear, editar y visualizar cada recurso.

**Objetivos:**
- Autenticar usuarios de forma segura
- Definir roles con permisos diferenciados (coordinador crea documentos, docente planifica clases)
- Asignar usuarios a instituciones, áreas y cursos

**Alcance MVP:**
- Autenticación de usuarios (email + password via Auth0)
- Roles de coordinator, teacher y admin con permisos diferenciados
- Asignación de usuarios a organizaciones y cursos

**No incluye:**
- Roles de directivos o supervisores → horizonte
- Gestión de múltiples instituciones por usuario → por definir

**Decisiones técnicas:**
- Un usuario puede tener **múltiples roles dentro de una misma organización**. Un docente puede ser profesor de dos materias y coordinador de un área — no hay restricción. La experiencia no necesita "Escoger un rol".
- Si un usuario trabaja en **dos instituciones distintas**, tiene dos cuentas separadas (un usuario por organización).
- El mecanismo de autenticación puede variar por provincia. Para el MVP: email + password via Auth0 (mismo sistema que tich-cronos).
- Los permisos sobre el documento de coordinación (quién edita, quién solo visualiza) son **configurables por organización**.

**Épicas relacionadas:** Onboarding (flujo post-auth), Documento de coordinación (permisos), Planificación docente (acceso a cursos)

---

### Épica 2: Onboarding (post-MVP)

> Carga de datos iniciales y product tour para nuevos usuarios.

**Problema:** Un usuario nuevo necesita completar su perfil y entender la plataforma antes de ser productivo. Sin un flujo guiado, el tiempo hasta el primer uso real es alto.

**Objetivos:**
- Capturar los datos necesarios del usuario al primer ingreso
- Guiar al usuario por las funcionalidades clave según su rol

**Decisiones técnicas:**
- El onboarding se dispara **post-autenticación al primer ingreso**. Los datos de la institución ya están cargados vía Integración.
- Los datos requeridos se definen como **configuración por organización** (mismo JSON de config).
- El product tour se adapta al **rol y a los feature flags activos** de la organización.

---

### Épica 3: Integración

> Importación de datos y configuración de la estructura curricular de cada provincia.

**Problema:** Cada provincia tiene su propio diseño curricular, terminología y estructura de conocimientos. El sistema necesita incorporar estos datos como base para que todo lo que Alizia genere esté alineado con la realidad de cada jurisdicción.

**Objetivos:**
- Importar diseños curriculares, NAPs y fuentes oficiales de cada provincia
- Modelar la estructura curricular provincial (áreas, disciplinas, tópicos, sub-tópicos, actividades)
- Que estos datos sean la fuente de verdad para la generación de documentos y las respuestas del asistente IA

**Alcance MVP:**
- Importación de diseño curricular y fuentes oficiales
- Configuración de la estructura curricular (jerarquía de conocimientos)
- Customizaciones por institución: grillas horarias, alta de docentes

**Decisiones técnicas:**
- La configuración de cada organización se almacena en un **JSON de configuración** a nivel organización — mismo patrón que TUNI.
- Los tópicos se modelan en una **tabla única auto-referencial** (foreign key a sí misma). La profundidad se **pre-computa y almacena** para evitar recursiones costosas.
- Cada organización define la **profundidad máxima** y los **nombres por nivel**.
- El setup inicial se realiza **manualmente por el equipo de implementación** — no hay backoffice self-service en el MVP.
- Las migraciones se hacen **incrementalmente**. Lección aprendida de TUNI: schema gigante upfront genera redefiniciones constantes.
- Un área puede contener **una o más disciplinas**. No se fuerza la existencia de áreas como concepto obligatorio.

---

### Épica 4: Documento de coordinación

> Creación asistida, edición colaborativa y gestión de estados de documentos de coordinación areal.

**Problema:** Los equipos de docentes necesitan documentos complejos que deben alinear múltiples personas en el trabajo a realizarse en un plazo futuro dado. Este proceso es manual, lento y difícil de articular entre roles.

**Objetivos:**
- Generar una primera versión del documento con IA, alineada a lineamientos provinciales
- Permitir la edición colaborativa entre roles
- Asegurar que el documento sea la base sobre la cual se planifica el clase a clase

**Alcance MVP:**
- Generación asistida del documento basado en topics seleccionados
- Cálculo automático de cantidad de clases por disciplina según grilla horaria
- Selección de sub-topics (pueden repetirse entre disciplinas, configurable)
- Generación de cronograma tentativo con objetivo por clase
- Edición directa y asistida por IA
- Publicación del documento

**Sub-épicas:**

| Componente | Descripción |
|---|---|
| Asistente de creación | Genera elementos del documento: eje problemático, estrategia metodológica y cronograma |
| Editor de documento | Edición directa y asistida por IA de cada sección |
| Flujo de estados | Gestión del ciclo de vida: draft → published → archived |

**Decisiones técnicas:**
- El **período** es nombre libre con rango de fechas custom — no se asume duración fija.
- Cada sección se define mediante **JSON Schema** (estructura output, prompt, opciones). Customizable por provincia sin cambios de código.
- **Clases coordinadas** es feature flag. Si dos materias coinciden en horario, las modificaciones se sincronizan. Diferenciador clave — ninguna plataforma del mercado lo ofrece.
- El cálculo de clases se basa en la grilla horaria. Coordinador puede override manual (± feriados).
- La selección de sub-topics respeta la **profundidad configurada por organización**. UI se adapta dinámicamente.

**Decisiones por provincia:**
- Estrategias metodológicas disponibles (proyecto, taller, ateneo, laboratorio) — requieren validación con equipo pedagógico
- Nivel de edición del docente sobre el documento del coordinador

---

### Épica 5: Planificación docente

> Planificación del clase a clase por momentos, con asistencia de IA y feedback post-dictado.

**Problema:** Los docentes planifican sus clases sin conexión con lo acordado a nivel área, ni con lo que están dictando otros docentes. Re-planificar en base a lo que ocurrió depende 100% del docente, su energía y su memoria.

**Objetivos:**
- Cada clase alineada con el documento de coordinación del área
- Cada clase coherente con lo que se dicta en otras disciplinas
- Generar propuesta inicial de actividades por momento
- Incorporar feedback de clases anteriores (post-MVP)
- Personalización del docente sin perder alineación curricular

**Alcance MVP:**
- Visualización del cronograma heredado del documento de coordinación
- Edición del objetivo de clase
- Selección de actividades por momento con recomendaciones de IA
- Personalización: anclar clase a recurso o comentario
- Generación de propuesta detallada
- Gestión de estados (pending → planned)

**Post-MVP:**
- Bitácora de cotejo (reportar cómo fue la clase, soporta audio)
- Recolección activa de datos (Alizia pregunta por info faltante)
- Repropuesta: cambios sugeridos a clases futuras ya planificadas basados en bitácora

**Sub-épicas:**

| Componente | Descripción |
|---|---|
| Plan de clase | Selección de actividades por momento y generación de propuesta |
| Momentos didácticos | Configuración de tipos de actividad por momento (configurable por org) |
| Incorporación de fuentes | Anclaje de la clase a un recurso o fuente específica |
| Edición del documento | Edición directa o asistida, gestión de estado |
| Bitácora (post-MVP) | Recolección del resultado de una clase |
| Repropuesta (post-MVP) | Sugerencias de cambios para clases futuras ya planificadas |

**Decisiones técnicas:**
- Un docente por materia por curso (first-come-first-serve si hay conflicto).
- Planificación exportable como **PDF** (template configurable) — NTH.
- **Momentos didácticos** (apertura, desarrollo, cierre) son fijos. Tipos de actividad dentro de cada momento son **configurables por organización**.
- Con feature flag de **clases coordinadas**, la planificación muestra indicador de clase compartida y sincroniza.
- Bitácora: el docente graba **audio libre** (no formato rígido). Adopción depende de que sea natural.

**Decisiones por provincia:**
- Tipos de actividad por momento se definen con cada equipo pedagógico
- Formato y profundidad de la bitácora requiere validación

---

### Épica 6: Asistente IA

> Asistente de inteligencia artificial que genera, edita y aprende del uso de la plataforma.

**Problema:** Producir documentos y planificaciones alineadas curricularmente es lento y requiere expertise. Un asistente genérico no conoce el contexto provincial ni el historial del aula.

**Objetivos:**
- Generar primeras versiones de documentos basadas en diseño curricular
- Ediciones masivas o puntuales por instrucción natural
- Incorporar historial de clases para mejorar recomendaciones
- Mantener alineación con fuentes oficiales de la provincia

**Alcance MVP:**
- Generación de documentos de coordinación (secciones + plan de clases)
- Recomendación de actividades por momento
- Edición asistida: instrucciones en lenguaje natural → Alizia modifica
- Procesamiento de feedback post-clase (post-MVP)

**Sub-épicas:**

| Componente | Descripción |
|---|---|
| Modificación de contenido | Edición asistida por instrucción natural |
| Asistencia de uso y navegación | Ayuda contextual |
| Customización por cliente | Tono, límites, fuentes por provincia |

**Decisiones técnicas:**
- Opera como **LLM con tools** (function calling). Puede leer y modificar el documento activo, recomendar actividades, consultar estructura curricular.
- Cada sección generada tiene **prompt + JSON Schema** por organización. Variables contextuales (tópicos, disciplina, grilla) inyectadas en el prompt.
- El **planificador de clases** recibe todos los temas, disciplinas y clases disponibles, y distribuye contenido inteligentemente.
- Generación de secciones arranca con **prompt simple**. Sofisticación viene con uso real, no con sobre-ingeniería anticipada.
- Chat permite instrucciones naturales ("cambiá la actividad del cierre por algo más dinámico"). Internamente: chat con tools sobre documento activo.

**Function calling tools:**

| Tool | Descripción |
|---|---|
| `update_section(section_key, content)` | Actualiza una sección del documento. Valida que `section_key` exista en schema de la org |
| `update_class(class_number, title, objective)` | Modifica una clase del plan |
| `update_class_topics(class_number, topic_ids)` | Cambia topics de una clase |

---

### Épica 7: Dashboard (post-MVP)

> Vista consolidada del estado de documentos, cursos y notificaciones.

**Problema:** Coordinadores y docentes no tienen un lugar único donde ver el estado de sus documentos, planificaciones y cursos.

**Objetivos:**
- Visibilidad rápida del estado de documentos y planificaciones
- Centralizar acceso a cursos asignados
- Notificar cambios relevantes

**Decisiones técnicas:**
- Lo que ve cada usuario depende de **rol y configuración de la org**. Coordinador ve documentos + cursos del área. Docente ve planificaciones + clases próximas.
- Notificaciones cubren: publicación de documento (docente ya puede planificar), modificaciones en clases coordinadas, plazos próximos.

---

### Épica 8: Contenido

> Biblioteca de recursos didácticos y herramientas de creación basadas en fuentes oficiales.

**Problema:** Los docentes necesitan recursos didácticos adaptados a su contexto curricular. Crearlos desde cero es lento y recurrir a fuentes no curadas genera inconsistencias.

**Objetivos:**
- Creación de recursos a partir de fuentes oficiales validadas
- Tipos predefinidos (ficha de cátedra, guía de lectura, etc.) + creación libre
- Recursos utilizables en planificaciones de clase

**Alcance MVP:**
- Creación desde fuentes oficiales
- Tipos predefinidos (más a definir con ministerio)
- Edición directa y asistida por IA
- Exportación para impresión
- Disponibilización para uso en clase

**Sub-épicas:**

| Componente | Descripción |
|---|---|
| Biblioteca | Repositorio de recursos creados, reutilizables por org |
| Creación de contenido | Generación asistida desde fuentes + tipo seleccionado |

**Decisiones técnicas:**
- Cada org **habilita qué tipos** tiene disponibles. Se activan según acuerdo con equipo pedagógico provincial.
- **Biblioteca central** por org. Todos los docentes pueden explorar y reutilizar antes de generar nuevo. Reduce costos y promueve consistencia.
- Filtro por materia es **restricción soft** (UX, no permisos).
- Cada tipo tiene **prompt + JSON Schema** configurable por provincia.
- MVP: si segundo cliente necesita variante, **duplicar y adaptar**. Genericidad cuando aparezca el patrón real.

**Tipos públicos vs privados:**

| `organization_id` | Visibilidad | Ejemplo |
|---|---|---|
| `NULL` | Público: visible para todas las orgs | `lecture_guide`, `course_sheet` |
| Set | Privado: solo para esa org | Tipos custom |

**Flujo de generación IA:**

1. Docente elige tipo → (si `requires_font`) elige fuente → crea `resources` con `content = {}`
2. Resuelve prompt: `organization_resource_types.custom_prompt` ?? `resource_types.prompt`
3. Resuelve output_schema: `custom_output_schema` ?? `output_schema`
4. Envía al LLM con contexto (font, course_subject, etc.)
5. Respuesta se guarda en `resources.content` (JSONB) según schema
6. Frontend renderiza `content` dinámicamente según `output_schema`
7. Chat con Alizia puede editar secciones del `content`

---

### Épica 9: WhatsApp (pendiente definición)

> Canal de WhatsApp para interacción con Alizia fuera de la plataforma.

**Decisiones técnicas:**
- Canal **cross-producto** (aplica a Alizia + TUNI + futuros productos).
- Comportamiento inicial: **efecto wow + casos básicos** (pregunta rápida, resumen, disparar proceso en plataforma). No replicar funcionalidad completa.
- Requiere **autenticación secundaria**: vincular teléfono con usuario de plataforma.

---

### Épica 10: Cosmos (pendiente definición)

> Sin contenido definido.

---

## Patrones transversales

| Patrón | Épicas | Descripción |
|---|---|---|
| JSON de configuración por org | 1, 2, 3, 4, 5, 8 | Configuración provincial centralizada (feature flags, nombres de niveles, tipos habilitados) |
| Prompt + JSON Schema por sección | 4, 6, 8 | Cada output generado por IA tiene prompt y schema configurable por provincia |
| Feature flags por organización | 2, 4, 5, 8 | Funcionalidades que se activan/desactivan por cliente |
| Clases coordinadas | 4, 5 | Diferenciador clave: sincronización entre docentes que comparten horario |
| Decisiones por provincia | Todas | Cada cliente customiza comportamiento sin cambios de código |

### Mapa de dependencias entre épicas

```
Roles y accesos ──→ Onboarding ──→ (usuario listo para usar la plataforma)
        │
        ▼
   Integración ──→ Documento de coordinación ──→ Planificación docente
        │                    │                          │
        │                    ▼                          ▼
        │              Asistente IA ◄───────────── Bitácora
        │                    │
        ▼                    ▼
    Contenido ◄──────── Asistente IA

   Dashboard ◄── (consume estado de Docs + Planificaciones)

   WhatsApp ◄── Asistente IA (canal alternativo)

   Cosmos ── (pendiente)
```

---

## Arquitectura general 🏗️

```
┌─────────────────┐                    ┌──────────────────────┐
│  Alizia Frontend │                    │    Auth0              │
│  React + TS      │                    │    (SaaS)             │
└────────┬────────┘                    │    JWT + JWKS         │
         │                              └──────────┬───────────┘
         │ HTTPS                                   │ JWT firmado
         ▼                                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      alizia-api (este RFC)                       │
│                                                                  │
│  Go 1.26 + Gin (via team-ai-toolkit/web)                        │
│  Clean Architecture: entities → providers → usecases → handlers  │
│  GORM + PostgreSQL                                               │
│  Deploy: Railway (Docker container)                              │
│                                                                  │
│  Módulos:                                                        │
│  ├── admin        → orgs, areas, subjects, topics, courses       │
│  ├── coordination → documentos, wizard, secciones, publicación   │
│  ├── teaching     → lesson plans, momentos, actividades          │
│  ├── resources    → fonts, tipos de recurso, recursos generados  │
│  └── ai           → Azure OpenAI, function calling, chat         │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                    ┌──────────┴──────────┐
                    │                     │
             ┌──────┴──────┐    ┌────────┴────────┐
             │  PostgreSQL  │    │  Azure OpenAI    │
             │  (Railway)   │    │  (gpt-5-mini)    │
             └─────────────┘    └─────────────────┘
```

### Stack técnico

| Componente | Tecnología |
|---|---|
| Lenguaje | Go 1.26 |
| Framework | Gin (abstraído via team-ai-toolkit/web) |
| ORM | GORM (estándar empresa) |
| DB | PostgreSQL |
| Auth | Auth0 JWT + Bearer tokens (team-ai-toolkit/tokens valida via JWKS) |
| AI | Azure OpenAI SDK (gpt-5-mini) |
| Logging | slog (team-ai-toolkit/applog) |
| Error tracking | Bugsnag (team-ai-toolkit/applog/bugsnag) |
| Testing | testify + GORM, target 80% |
| Linting | golangci-lint (15+ linters) |
| Deploy | Railway (Docker, auto-deploy desde GitHub) |

Ver [ARQUITECTURA-GO-ALIZIA-V2.md](./ARQUITECTURA-GO-ALIZIA-V2.md) para estructura de directorios completa, patrones de código, y decisiones técnicas detalladas.

---

## Backend — Endpoints ⚙️

### Admin (Fase 2)

| Método | Ruta | Descripción | Auth | Roles |
|--------|------|-------------|------|-------|
| POST | `/api/v1/areas` | Crear área | Sí | coordinator, admin |
| GET | `/api/v1/areas` | Listar áreas de la org | Sí | Todos |
| PUT | `/api/v1/areas/:id` | Actualizar área | Sí | coordinator, admin |
| POST | `/api/v1/areas/:id/coordinators` | Asignar coordinador | Sí | admin |
| POST | `/api/v1/subjects` | Crear materia | Sí | coordinator, admin |
| GET | `/api/v1/subjects` | Listar materias | Sí | Todos |
| POST | `/api/v1/courses` | Crear curso | Sí | admin |
| GET | `/api/v1/courses` | Listar cursos | Sí | Todos |
| GET | `/api/v1/courses/:id` | Detalle con students + schedule | Sí | Todos |
| POST | `/api/v1/courses/:id/time-slots` | Crear time slot | Sí | admin |
| POST | `/api/v1/topics` | Crear topic | Sí | admin |
| GET | `/api/v1/topics` | Listar topics (tree) | Sí | Todos |

### Coordination Documents (Fase 3)

| Método | Ruta | Descripción | Auth | Roles |
|--------|------|-------------|------|-------|
| POST | `/api/v1/coordination-documents` | Crear documento (wizard) | Sí | coordinator |
| GET | `/api/v1/coordination-documents` | Listar docs (filtro por area) | Sí | coordinator, teacher |
| GET | `/api/v1/coordination-documents/:id` | Detalle completo | Sí | coordinator, teacher |
| PATCH | `/api/v1/coordination-documents/:id` | Actualizar (sections, status) | Sí | coordinator |
| DELETE | `/api/v1/coordination-documents/:id` | Eliminar (solo draft) | Sí | coordinator |
| POST | `/api/v1/coordination-documents/:id/subjects` | Asignar materias + class_count | Sí | coordinator |
| POST | `/api/v1/coordination-documents/:id/generate` | Generar secciones + plan con IA | Sí | coordinator |
| POST | `/api/v1/coordination-documents/:id/chat` | Chat con Alizia (function calling) | Sí | coordinator |

### Teaching (Fase 5)

| Método | Ruta | Descripción | Auth | Roles |
|--------|------|-------------|------|-------|
| GET | `/api/v1/course-subjects/:id/lesson-plans` | Lesson plans del docente | Sí | teacher |
| POST | `/api/v1/lesson-plans` | Crear lesson plan | Sí | teacher |
| PATCH | `/api/v1/lesson-plans/:id` | Actualizar | Sí | teacher |
| POST | `/api/v1/lesson-plans/:id/generate-activity` | Generar contenido IA por actividad | Sí | teacher |

### Resources (Fase 6)

| Método | Ruta | Descripción | Auth | Roles |
|--------|------|-------------|------|-------|
| GET | `/api/v1/resource-types` | Tipos disponibles para la org | Sí | teacher |
| GET | `/api/v1/fonts` | Fuentes educativas del area | Sí | Todos |
| POST | `/api/v1/resources` | Crear recurso | Sí | teacher |
| PATCH | `/api/v1/resources/:id` | Actualizar | Sí | teacher |
| POST | `/api/v1/resources/:id/generate` | Generar con IA | Sí | teacher |

### AI (Fase 4)

| Método | Ruta | Descripción | Auth | Roles |
|--------|------|-------------|------|-------|
| POST | `/api/v1/chat` | Chat general con Alizia | Sí | Todos |

---

## Backend — Modelo de datos ⛁

### Resumen: 26 tablas

El modelo completo con diagrama ER, SQL, queries de ejemplo, triggers, constraints, y cambios vs POC está en [proposal-der-v2.md](./proposal-der-v2.md).

### Enums

```sql
CREATE TYPE coord_doc_status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE lesson_plan_status AS ENUM ('pending', 'planned');
CREATE TYPE resources_mode AS ENUM ('global', 'per_moment');
CREATE TYPE class_moment AS ENUM ('apertura', 'desarrollo', 'cierre');
CREATE TYPE member_role AS ENUM ('teacher', 'coordinator', 'admin');
CREATE TYPE resource_status AS ENUM ('draft', 'active');
```

### Tablas y su propósito

| Tabla | Descripción |
|-------|-------------|
| `organizations` | Tenant. Cada cliente es una org con config JSONB custom |
| `users` | Docentes, coordinadores, admins. Pertenecen a una única org (en auth_db) |
| `user_roles` | Roles del usuario (teacher, coordinator, admin). Puede tener varios |
| `areas` | Agrupación de materias (ej: "Ciencias"). Opcional según config |
| `area_coordinators` | Qué usuarios coordinan qué áreas (M2M) |
| `subjects` | Materias individuales (ej: "Matemáticas"). Pertenecen a un área |
| `topics` | Jerarquía dinámica de temas/saberes. Self-referential con niveles configurables |
| `courses` | Grupos de alumnos (ej: "2do 1era") |
| `students` | Alumnos de un curso. Solo lectura informativa |
| `course_subjects` | Instancia: curso + materia + docente + período lectivo |
| `time_slots` | Slots horarios semanales de un curso (día + hora) |
| `time_slot_subjects` | course_subjects por slot. 2 registros = clase compartida |
| `coordination_documents` | **Output principal**. Planificación anual con `sections` JSONB dinámico |
| `coord_doc_topics` | Junction: topics seleccionados para el doc |
| `coordination_document_subjects` | Materias en doc con `class_count` |
| `coord_doc_subject_topics` | Junction: topics asignados a materia dentro del doc |
| `coord_doc_classes` | Plan de clases por materia (class_number, title, objective) |
| `coord_doc_class_topics` | Junction: topics cubiertos en cada clase |
| `activities` | Actividades didácticas predefinidas por org y momento |
| `teacher_lesson_plans` | Plan de clase del docente. Moments JSONB con actividades + contenido IA |
| `lesson_plan_topics` | Junction: topics cubiertos en un lesson plan |
| `lesson_plan_moment_fonts` | Junction: font por momento o global |
| `fonts` | Fuentes educativas (PDFs, videos, docs). `is_validated` = aprobado por coordinadores |
| `resource_types` | Tipos con prompt IA y output_schema. `org_id NULL` = público |
| `organization_resource_types` | Override por org: enable/disable + custom prompt/schema |
| `resources` | Recurso generado. `content` JSONB según output_schema del tipo |

### Relaciones clave (del diagrama ER)

```
organizations ──1:N──▶ areas ──1:N──▶ subjects
organizations ──1:N──▶ topics (self-referential via parent_id)
organizations ──1:N──▶ courses ──1:N──▶ students
                                courses ──1:N──▶ course_subjects
                                courses ──1:N──▶ time_slots ──1:N──▶ time_slot_subjects

areas ──1:N──▶ coordination_documents ──1:N──▶ coord_doc_topics
                                       ──1:N──▶ coordination_document_subjects
                                                  ──1:N──▶ coord_doc_subject_topics
                                                  ──1:N──▶ coord_doc_classes
                                                             ──1:N──▶ coord_doc_class_topics

course_subjects ──1:N──▶ teacher_lesson_plans ──1:N──▶ lesson_plan_topics
                                               ──1:N──▶ lesson_plan_moment_fonts

resource_types ──1:N──▶ resources
resource_types ──1:N──▶ organization_resource_types
fonts ──0:N──▶ resources (opcional)
```

### Normalización: coordination_documents

**Antes (POC — JSONB + arrays):**
```
coordination_documents.subjects_data → JSONB con topic_ids[], class_plan[], category_ids[]
```

**Ahora (tablas normalizadas):**
```
coordination_documents
  └── coord_doc_topics (doc ↔ topic)
  └── coordination_document_subjects (doc ↔ subject + class_count)
        └── coord_doc_subject_topics (subject en doc ↔ topic)
        └── coord_doc_classes (class_number, title, objective)
              └── coord_doc_class_topics (clase ↔ topic)
```

### Constraints UNIQUE en junction tables

| Tabla | Constraint |
|-------|-----------|
| `user_roles` | `UNIQUE(user_id, role)` |
| `area_coordinators` | `UNIQUE(area_id, user_id)` |
| `time_slot_subjects` | `UNIQUE(time_slot_id, course_subject_id)` |
| `coord_doc_topics` | `UNIQUE(coordination_document_id, topic_id)` |
| `coord_doc_subject_topics` | `UNIQUE(coord_doc_subject_id, topic_id)` |
| `coord_doc_class_topics` | `UNIQUE(coord_doc_class_id, topic_id)` |
| `lesson_plan_topics` | `UNIQUE(lesson_plan_id, topic_id)` |
| `lesson_plan_moment_fonts` | `UNIQUE(lesson_plan_id, moment, font_id)` |
| `organization_resource_types` | `UNIQUE(organization_id, resource_type_id)` |

### Topics: jerarquía dinámica

Reemplaza las 3 tablas fijas del POC (`problematic_nuclei`, `knowledge_areas`, `categories`) con una sola tabla self-referential.

**Reglas de nivel:**
- `parent_id IS NULL` → `level = 1` (raíz)
- `parent_id IS NOT NULL` → `level = parent.level + 1`
- No puede exceder `config.topic_max_levels`

**Recálculo automático** al mover un topic:

```sql
WITH RECURSIVE tree AS (
    SELECT id, parent_id,
           COALESCE((SELECT level FROM topics WHERE id = NEW.parent_id), 0) + 1 AS level
    FROM topics WHERE id = NEW.id
    UNION ALL
    SELECT t.id, t.parent_id, tree.level + 1
    FROM topics t JOIN tree ON t.parent_id = tree.id
)
UPDATE topics SET level = tree.level
FROM tree WHERE topics.id = tree.id;
```

### Horarios: clases compartidas

**Clase normal**: 1 `time_slot` → 1 `time_slot_subject`
**Clase compartida**: 1 `time_slot` → 2 `time_slot_subjects`

**Trigger de validación:**

```sql
CREATE OR REPLACE FUNCTION validate_time_slot_subject() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM course_subjects cs
        JOIN time_slots ts ON ts.course_id = cs.course_id
        WHERE cs.id = NEW.course_subject_id AND ts.id = NEW.time_slot_id
    ) THEN
        RAISE EXCEPTION 'course_subject does not belong to the same course as the time_slot';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### Migraciones

| Orden | Migración | Tablas |
|-------|-----------|--------|
| 1 | init | Enums + organizations + areas + subjects + topics + courses + students + course_subjects + time_slots + time_slot_subjects + activities |
| 2 | coordination | coordination_documents + coord_doc_topics + coordination_document_subjects + coord_doc_subject_topics + coord_doc_classes + coord_doc_class_topics |
| 3 | teaching | teacher_lesson_plans + lesson_plan_topics + lesson_plan_moment_fonts |
| 4 | resources | fonts + resource_types + organization_resource_types + resources |

---

## Backend — Lógica y configuración 🔧

### Configuración completa por organización

```jsonc
{
  // --- Taxonomía de temas ---
  "topic_max_levels": 3,
  "topic_level_names": [
    "Núcleos problemáticos",
    "Áreas de conocimiento",
    "Categorías"
  ],
  "topic_selection_level": 3,

  // --- Clases compartidas ---
  "shared_classes_enabled": true,

  // --- Secciones del documento de coordinación ---
  "coord_doc_sections": [
    {
      "key": "problem_edge",
      "label": "Eje problemático",
      "type": "text",
      "ai_prompt": "Generá un eje problemático que integre las categorías seleccionadas...",
      "required": true
    },
    {
      "key": "methodological_strategy",
      "label": "Estrategia metodológica",
      "type": "select_text",
      "options": ["proyecto", "taller_laboratorio", "ateneo_debate"],
      "ai_prompt": "Generá una estrategia metodológica de tipo {selected_option}...",
      "required": true
    },
    {
      "key": "eval_criteria",
      "label": "Criterios de evaluación",
      "type": "text",
      "ai_prompt": "Generá criterios de evaluación para las categorías seleccionadas...",
      "required": false
    }
  ],

  // --- Lesson plans ---
  "desarrollo_max_activities": 3
}
```

### Estructura JSONB: sections del coordination document

```json
{
  "problem_edge": {
    "value": "¿Cómo las lógicas de poder y saber configuran..."
  },
  "methodological_strategy": {
    "selected_option": "proyecto",
    "value": "Implementaremos un ateneo-debate interdisciplinario..."
  },
  "eval_criteria": {
    "value": "Los criterios de evaluación serán..."
  }
}
```

### Estructura JSONB: moments del teacher_lesson_plan

```json
{
  "apertura": {
    "activities": [1],
    "activityContent": { "1": "Texto generado por IA para actividad 1..." }
  },
  "desarrollo": {
    "activities": [3, 5],
    "activityContent": { "3": "...", "5": "..." }
  },
  "cierre": {
    "activities": [8],
    "activityContent": { "8": "..." }
  }
}
```

### Query: tipos de recurso disponibles para una org

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

### Query: detectar clases compartidas

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

### Query: clases con topics de un documento

```sql
SELECT cdc.class_number, cdc.title, array_agg(t.name) AS topics
FROM coord_doc_classes cdc
JOIN coordination_document_subjects cds ON cds.id = cdc.coord_doc_subject_id
LEFT JOIN coord_doc_class_topics cdct ON cdct.coord_doc_class_id = cdc.id
LEFT JOIN topics t ON t.id = cdct.topic_id
WHERE cds.coordination_document_id = $1 AND cds.subject_id = $2
GROUP BY cdc.id ORDER BY cdc.class_number;
```

---

## QA — Estrategia de testing 🧪

### Precondiciones

- PostgreSQL corriendo con schema migrado
- Organización seed con config de ejemplo (3 niveles de topics, clases compartidas habilitadas)
- Usuarios seed: admin, coordinator, teacher
- Topics seed: jerarquía de 3 niveles
- Areas, subjects, courses, time_slots seed

### Matriz por fase

#### Fase 2: Admin/Integration

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 2.1 | Crear área | 201 | Alta |
| 2.2 | Crear topic respetando max_levels | 201 si OK, 400 si excede | Alta |
| 2.3 | Crear time_slot con clase compartida | 201 si enabled, 400 si disabled | Alta |
| 2.4 | Trigger: course_subject debe pertenecer al mismo curso | Error si no pertenece | Alta |
| 2.5 | Listar topics como árbol | Jerarquía correcta | Media |

#### Fase 3: Coordination Documents

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 3.1 | Crear doc vía wizard | 201 + doc en draft | Alta |
| 3.2 | Asignar materias + class_count | Subjects vinculados | Alta |
| 3.3 | Asignar topics a materias | Validar que todos los topics del doc estén cubiertos | Alta |
| 3.4 | Publicar doc | Status → published | Alta |
| 3.5 | Docente no puede editar doc publicado (si config lo restringe) | 403 | Alta |
| 3.6 | DELETE solo funciona en draft | 400 si published | Media |

#### Fase 4: AI Generation

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 4.1 | Generar secciones con IA | Secciones populadas según config | Alta |
| 4.2 | Generar plan de clases | Classes creadas con topics distribuidos | Alta |
| 4.3 | Chat update_section | Sección actualizada, key validada contra config | Alta |
| 4.4 | Chat update_class | Clase modificada | Media |

#### Fase 5: Teaching

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 5.1 | Crear lesson plan heredando de doc | 201 + datos del doc | Alta |
| 5.2 | Seleccionar actividades respetando limits | 1 apertura, 1-3 desarrollo, 1 cierre | Alta |
| 5.3 | Generar contenido por actividad | activityContent populado | Media |
| 5.4 | Fonts global vs por momento | Comportamiento correcto según resources_mode | Media |

#### Fase 6: Resources

| # | Caso | Resultado esperado | Prioridad |
|---|------|--------------------|-----------|
| 6.1 | Listar tipos disponibles para org | Solo públicos + privados de la org | Alta |
| 6.2 | Tipo con custom_prompt usa override | Prompt correcto en generación | Alta |
| 6.3 | Tipo con requires_font sin font | 400 | Media |
| 6.4 | Generar recurso con IA | content populado según output_schema | Alta |

### Coverage target: 80%

---

## Alternativas evaluadas 🔀

### Arquitectura

| Criterio | POC (FastAPI) | Go + GORM + Railway |
|----------|:---:|:---:|
| Performance | ❌ Python | ✅ Go |
| Multi-tenancy | ❌ | ✅ org_id en JWT |
| Equipo conoce | ✅ Python | ✅ Go (tich-cronos) |
| Escalabilidad | ❌ | ✅ |
| Infra compartida | ❌ | ✅ team-ai-toolkit |

> **Elegido: Go + GORM + Railway**

### ORM

| Criterio | GORM | sqlx |
|----------|:---:|:---:|
| Equipo lo conoce | ✅ | ❌ |
| CRUD rápido | ✅ | ❌ |
| Queries complejas | ❌ (Raw SQL) | ✅ |
| Performance | Buena | Óptima |

> **Elegido: GORM** — estándar empresa. sqlx documentado como alternativa futura en ARQUITECTURA-GO-ALIZIA-V2.md.

---

## Rollout 📈

| Fase | Alcance | Criterio para avanzar |
|------|---------|----------------------|
| 1 | Staging — equipo interno | CI verde, /health, auth funciona |
| 2 | Org piloto (1 provincia) | Coordinador crea doc + docente planifica |
| 3 | 2-3 orgs adicionales | Feedback positivo, sin bugs bloqueantes |
| 4 | Todas las orgs | Métricas de éxito alcanzadas |

### Plan de rollback

1. Railway: revertir al deploy anterior (1 click)
2. Migraciones: ejecutar `.down.sql`
3. Comunicar al equipo

---

## Dependencias 👥

| Dependencia | Tipo | Bloqueante | Estado |
|-------------|------|------------|--------|
| team-ai-toolkit | Librería Go | Sí | ✅ Creado |
| auth-service | Microservicio | No (futuro) | ⬜ Futuro (no bloqueante — se arranca con Auth0) |
| Railway account | Infra | Sí (Fase 1) | ⬜ Configurar |
| PostgreSQL en Railway | Infra | Sí (Fase 1) | ⬜ Provisionar |
| Azure OpenAI access | Servicio | Sí (Fase 4) | ✅ Ya disponible |
| Diseño UX/UI | Entregable | No (backend first) | ⬜ En progreso |
| Auth0 tenant config | Infra | Sí (Fase 1) | ⬜ Configurar (domain + audience + API) |

---

## Riesgos ⚠️

| # | Riesgo | Probabilidad | Impacto | Mitigación |
|---|--------|-------------|---------|------------|
| 1 | GORM genera queries N+1 | Media | Medio | Preload explícito + Raw SQL para JOINs complejos |
| 2 | Config JSONB inmanejable | Baja | Alto | Validación en backend, schema documentation |
| 3 | IA genera contenido de baja calidad | Media | Medio | Prompts iterativos, review humano antes de publicar |
| 4 | Multi-tenancy data leak | Baja | Crítico | Middleware de tenant en TODAS las rutas, tests de isolation |
| 5 | Railway downtime | Baja | Medio | Dockerfile portable, migrar a otro hosting en horas |
| 6 | CTE recursivo de topics lento con muchos niveles | Baja | Medio | Level precalculado, solo recalcular rama afectada |

---

## Preguntas abiertas ❓

| # | Pregunta | Estado |
|---|----------|--------|
| 1 | ¿Cómo se cargan los datos iniciales de una provincia? (manual, CSV, API) | 🟡 Pendiente |
| 2 | ¿Quién crea las organizaciones? (super admin o self-service) | 🟡 Pendiente |
| 3 | ¿Los docentes pueden ver docs de otras áreas? | 🟡 Pendiente |
| 4 | ¿Se mantiene historial de versiones de documents? | 🟡 Pendiente |
| 5 | ¿Concurrent users esperados por org? | 🟡 Pendiente |
| 6 | ¿Qué pasa si el docente no tiene internet al grabar bitácora? (escuelas rurales) | 🟡 Pendiente |
| 7 | ¿Se permite subida de fuentes propias del docente? (decisión por provincia) | 🟡 Pendiente |

---

## Glosario 📖

| Término | Definición |
|---------|-----------|
| Coordination Document | Documento de planificación anual de un área, creado por el coordinador |
| Lesson Plan | Plan de clase individual creado por el docente, hereda del coordination document |
| Topic | Tema/saber en la jerarquía curricular (self-referential, niveles configurables) |
| Class Moment | Momento didáctico: apertura, desarrollo, cierre |
| Shared Class / Clase compartida | Dos materias enseñadas simultáneamente por dos docentes en el mismo horario. Diferenciador clave del producto |
| Font | Fuente educativa (PDF, video, documento) — del español "fuente", NO tipografía |
| Resource Type | Tipo de recurso generado por IA (guía de lectura, ficha de curso, etc.) |
| Organization / org | Tenant: una provincia, escuela, o universidad con configuración propia |
| team-ai-toolkit | Librería Go compartida con infra reutilizable (web, boot, tokens, etc.) |
| Railway | Plataforma de hosting para containers Docker |
| Bitácora | Registro post-clase del docente sobre cómo fue la clase (soporta audio). Post-MVP |
| Repropuesta | Sugerencia automática de cambios a clases futuras basada en bitácora. Post-MVP |
| NAP | Núcleos de Aprendizajes Prioritarios — lineamientos curriculares nacionales argentinos |

---

## Tareas 📝

### Fase 1 — Setup

| # | Tarea | Estado |
|---|-------|--------|
| 1.1 | Crear repo alizia-api con estructura de directorios | ⬜ |
| 1.2 | Configurar go.mod con team-ai-toolkit | ⬜ |
| 1.3 | Configurar CI (GitHub Actions: test + lint) | ⬜ |
| 1.4 | Provisionar Railway + PostgreSQL | ⬜ |
| 1.5 | Configurar Auth0 tenant (domain, audience, API) para staging + prod | ⬜ |
| 1.6 | Deploy inicial (/health responde) | ⬜ |
| 1.7 | Integrar auth middleware | ⬜ |

### Fase 2 — Admin/Integration

| # | Tarea | Estado |
|---|-------|--------|
| 2.1 | Migración init (enums + tablas base + triggers) | ⬜ |
| 2.2 | CRUD areas + area_coordinators | ⬜ |
| 2.3 | CRUD subjects | ⬜ |
| 2.4 | CRUD courses + students + course_subjects | ⬜ |
| 2.5 | CRUD topics (jerarquía con validación de niveles) | ⬜ |
| 2.6 | CRUD time_slots + time_slot_subjects (trigger same-course) | ⬜ |
| 2.7 | CRUD activities (por momento) | ⬜ |
| 2.8 | Tests de integración | ⬜ |

### Fase 3 — Coordination Documents

| # | Tarea | Estado |
|---|-------|--------|
| 3.1 | Migración coordination (6 tablas) | ⬜ |
| 3.2 | Crear documento (wizard 3 pasos) | ⬜ |
| 3.3 | Asignar materias + topics a materias | ⬜ |
| 3.4 | CRUD secciones dinámicas (sections JSONB) | ⬜ |
| 3.5 | Status workflow (draft → published → archived) | ⬜ |
| 3.6 | GET documento completo (con JOINs/Preloads) | ⬜ |
| 3.7 | Tests | ⬜ |

### Fase 4 — AI Generation

| # | Tarea | Estado |
|---|-------|--------|
| 4.1 | Azure OpenAI client wrapper | ⬜ |
| 4.2 | Generar secciones del doc (prompt por sección desde config) | ⬜ |
| 4.3 | Generar plan de clases (distribuir topics) | ⬜ |
| 4.4 | Chat con function calling (update_section, update_class, update_class_topics) | ⬜ |
| 4.5 | Tests con mock de AI client | ⬜ |

### Fase 5 — Teaching

| # | Tarea | Estado |
|---|-------|--------|
| 5.1 | Migración teaching (3 tablas + activities ya en Fase 2) | ⬜ |
| 5.2 | Crear lesson plan (hereda de doc) | ⬜ |
| 5.3 | Seleccionar actividades por momento (validar limits) | ⬜ |
| 5.4 | Seleccionar fonts (global o por momento) | ⬜ |
| 5.5 | Generar contenido por actividad (IA) | ⬜ |
| 5.6 | Tests | ⬜ |

### Fase 6 — Resources

| # | Tarea | Estado |
|---|-------|--------|
| 6.1 | Migración resources (4 tablas) | ⬜ |
| 6.2 | CRUD resource types + org overrides | ⬜ |
| 6.3 | CRUD fonts | ⬜ |
| 6.4 | Crear recurso + generar con IA (prompt resolution) | ⬜ |
| 6.5 | Query tipos disponibles por org | ⬜ |
| 6.6 | Tests | ⬜ |
