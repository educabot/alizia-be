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
- [ ] Enum `coord_doc_status` creado: pending, in_progress, published
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
  └── coordination_document_subjects (doc ↔ subject + class_count)
        ├── coord_doc_subject_topics (subject en doc ↔ topic)
        └── coord_doc_classes (class_number, title, objective)
              └── coord_doc_class_topics (clase ↔ topic)
```

## Test cases

- 4.1: Crear documento → todas las tablas relacionadas se crean correctamente
- 4.2: GET detalle → retorna documento con topics, subjects, classes, todo preloaded
- 4.3: Eliminar documento → cascade elimina todas las junction tables
