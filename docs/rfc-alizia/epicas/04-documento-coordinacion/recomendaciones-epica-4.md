# Recomendaciones — Epica 4: Documento de Coordinacion

> **Objetivo:** Postura recomendada para cada pregunta abierta, con razonamiento.
> Documento historico — las **decisiones finales** estan en las preguntas abiertas y reflejadas en los HU/tareas del RFC.
>
> **Fecha:** 2026-04-22
> **Autor:** Sebastian (con analisis de Claude)
> **Reunion de decisiones:** 2026-04-22 (Juan, Jose, Sebastian)

---

## Resumen: recomendacion vs decision final

| # | Pregunta | Recomendacion | Decision final | Coincide |
|---|----------|---------------|----------------|----------|
| P1 | Estado pending? | A) Eliminar | C) Mantener con auto-transicion | No |
| P2 | Multiples docs por area? | B) Un published por area | A) Sin restriccion + archived | No |
| P3 | Sin grilla horaria? | Permitir manual | Confirmado | Si |
| P4 | Topics compartidos? | Si, siempre | Confirmado | Si |
| P5 | Class count = 0? | No, minimo 1 | Confirmado | Si |
| P6 | select_text: cuando elige? | A) Antes de generar | C) Generar N variantes + manual | No |
| P7 | Granularidad de re-generacion? | D) MVP simple, v2 granular | E) Sin boton, todo por chat | No |
| P8 | Tipos de seccion? | Solo text/select_text | Confirmado | Si |
| P9 | IA genera count incorrecto? | A) Reintentar 1 vez | A) Confirmado | Si |
| P10 | CRUD manual de clases? | A) No en MVP | C) Inmutables post-publish, teacher edita | No |
| P11 | Clases compartidas en IA? | A) Solo marcar | B) IA coordina | No |
| P12 | Publicado es editable? | B) Editable con warning | C) Secciones si, clases no | No |
| P13 | Que se valida al publicar? | B) Medio | Sub-topics: warning confirmable | Parcial |
| P14 | Docentes ven no-publicados? | No | Confirmado | Si |
| P15 | Persistir chat? | A) No en MVP | B) Persistir + auto-compact | No |
| P16 | Tools del chat? | A) Set minimo | A) Confirmado (sin tools clase en published) | Si |
| P17 | Falla de IA? | A) 1 reintento | A) Confirmado + timeout reducido | Si |
| P18 | Tool calls atomicos? | Best-effort | Confirmado + retry | Si |
| P19 | Edicion concurrente? | A) Last write wins | Auto-guardado (TBD) | No |
| P20 | Undo/versionado? | B) 1 snapshot | B+C) 3 snapshots | Parcial |

**Resultado:** 9 coincidencias, 9 divergencias, 2 parciales. Las divergencias van todas en direccion de mas funcionalidad para MVP, priorizando UX del coordinador.

---

## Detalle original de recomendaciones

> Las recomendaciones a continuacion son el input pre-reunion. Para la decision final y el razonamiento, ver `preguntas-abiertas-epica-4.md` y los documentos de cada HU.

### P1: Eliminar `pending` — solo `in_progress` y `published`

El wizard ya recopila datos sustanciales: nombre, periodo, topics, materias, class_count. El documento nunca esta realmente "vacio" despues de la creacion. No hay un escenario claro donde `pending` aporte valor que `in_progress` no cubra.

**Decision final diferente:** Se mantiene `pending` con transicion automatica. El equipo valora poder distinguir "creado pero nunca tocado" de "en edicion activa".

---

### P2: Un `published` por area, multiples borradores

Los docentes necesitan saber exactamente cual es EL documento de su area. Si hay 3 publicados, no saben cual mirar.

**Decision final diferente:** Sin restriccion. Multiples documentos en cualquier estado. El frontend maneja la visualizacion. Se agrega estado `archived`.

---

### P6: Antes de generar — el coordinador elige primero

La decision pedagogica ("proyecto" vs "taller" vs "ateneo") es del coordinador, no de la IA.

**Decision final diferente:** Generar N variantes (una por opcion) + opcion manual. El equipo prioriza la experiencia exploratoria del coordinador.

---

### P7: MVP simple (todo o nada), v2 granular — Opcion D

Un solo boton "Generar" que regenera todo.

**Decision final diferente:** Sin boton "Regenerar". Toda la re-generacion va por el chat con Alizia. Simplifica la UI y da control granular via lenguaje natural.

---

### P11: Solo marcar clases compartidas, no coordinar con IA

La coordinacion de contenido entre materias por IA es un problema complejo.

**Decision final diferente:** IA coordina generando secuencialmente con contexto. El equipo acepta la complejidad adicional por el valor pedagogico.

---

### P12: Publicado = editable con warning

Razonamiento original: la realidad educativa exige poder editar post-publicacion.

**Decision final diferente:** Parcialmente editable. Secciones si, clases no. Las clases son la base de la planificacion docente — inmutables protege al teacher.

---

### P15: No persistir historial de chat en MVP

Lo que importa son los cambios en el documento, no la conversacion.

**Decision final diferente:** Persistir con auto-compactacion. El equipo valora la continuidad de la conversacion entre sesiones.

---

### P19: Last write wins para MVP

La probabilidad de edicion concurrente es baja en MVP.

**Decision final diferente:** Auto-guardado (detalles de implementacion TBD). El equipo quiere explorar auto-save antes de decidir estrategia de concurrencia.

---

### P20: Snapshot automatico antes de generar

Un solo snapshot es suficiente para el caso principal.

**Decision final diferente:** 3 snapshots. El equipo quiere mas margen de undo sin llegar al historial completo.
