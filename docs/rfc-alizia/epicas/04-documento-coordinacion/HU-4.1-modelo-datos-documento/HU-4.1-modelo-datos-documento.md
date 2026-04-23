# HU-4.1: Modelo de datos del documento

> Como coordinador, necesito que el documento de coordinación esté modelado con tablas normalizadas para que los datos sean confiables, trazables y consultables.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta (bloqueante para todo lo demás de esta épica)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `coordination_documents` con: id, organization_id, name, area_id, start_date, end_date, status (enum), sections (JSONB), created_at, updated_at
- [ ] Tabla `coord_doc_topics` (junction doc ↔ topic)
- [ ] Tabla `coordination_document_subjects` (doc ↔ subject + class_count)
- [ ] Tabla `coord_doc_subject_topics` (subject en doc ↔ topic)
- [ ] Tabla `coord_doc_classes` (class_number, title, objective por disciplina)
- [ ] Tabla `coord_doc_class_topics` (clase ↔ topic)
- [ ] Tabla `coord_doc_chat_messages` (historial de chat persistido) — P15
- [ ] Tabla `coord_doc_snapshots` (hasta 3 snapshots para undo) — P20
- [ ] Enum `coord_doc_status` creado: pending, in_progress, published, **archived** — P2
- [ ] Entities Go con GORM tags y relaciones (preloads)
- [ ] Provider interfaces para CRUD + operaciones complejas
- [ ] Repository GORM con queries de detalle (múltiples preloads)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.1.1 | [Migración: tablas del documento](./tareas/T-4.1.1-migracion.md) | db/migrations/ | ⬜ |
| 4.1.2 | [Entities](./tareas/T-4.1.2-entities.md) | src/core/entities/ | ⬜ |
| 4.1.3 | [Providers](./tareas/T-4.1.3-providers.md) | src/core/providers/ | ⬜ |
| 4.1.4 | [Repository GORM](./tareas/T-4.1.4-repository.md) | src/repositories/ | ⬜ |
| 4.1.5 | [Tests](./tareas/T-4.1.5-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../../03-integracion/HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — FK organization_id
- [HU-3.2: Áreas y disciplinas](../../03-integracion/HU-3.2-areas-materias/HU-3.2-areas-materias.md) — FK area_id, subject_id
- [HU-3.3: Topics](../../03-integracion/HU-3.3-topics-jerarquia-curricular/HU-3.3-topics-jerarquia-curricular.md) — FK topic_id en junction tables

## Diseño técnico

### Modelo normalizado

```
coordination_documents
  ├── coord_doc_topics (doc ↔ topic)
  ├── coordination_document_subjects (doc ↔ subject + class_count)
  │     ├── coord_doc_subject_topics (subject en doc ↔ topic)
  │     └── coord_doc_classes (class_number, title, objective)
  │           └── coord_doc_class_topics (clase ↔ topic)
  ├── coord_doc_chat_messages (historial de chat)
  └── coord_doc_snapshots (hasta 3 snapshots para undo)
```

### Tablas adicionales (decisiones P15 y P20)

**Chat messages (P15):** Historial de chat persistido en backend para que el coordinador pueda retomar conversaciones. Incluye role, content y tool_calls.

**Snapshots (P20):** Antes de cada generación IA, se guarda un snapshot del estado actual (secciones + clases). Se mantienen hasta **3 snapshots** por documento (no historial completo, pero más que un solo undo).

### Relación con topics

El documento de coordinación está relacionado a un **topic padre** (nivel seleccionado en `config.topic_selection_level`). Las materias dentro del documento se relacionan a **topics hijos** (sub-topics). Esta relación es clave para la validación al publicar (P13).

## Test cases

- 4.1: Crear documento → todas las tablas relacionadas se crean correctamente
- 4.2: GET detalle → retorna documento con topics, subjects, classes, todo preloaded
- 4.3: Eliminar documento → cascade elimina todas las junction tables, chat messages y snapshots
- 4.4: Chat message se persiste con role, content y tool_calls
- 4.5: Snapshot se crea y se limita a 3 por documento (el más viejo se elimina)
