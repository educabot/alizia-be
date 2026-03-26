# HU-4.5: Publicación y estados

> Como coordinador, necesito publicar el documento para que los docentes lo vean, y archivar documentos viejos.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Estados: `draft` → `published` → `archived`
- [ ] Solo coordinadores pueden cambiar estado
- [ ] Al publicar, se valida que todos los topics del documento estén distribuidos entre las materias
- [ ] Al publicar, se valida que todas las secciones requeridas tengan contenido
- [ ] Documento publicado es visible para docentes (GET listar y detalle)
- [ ] Solo documentos en `draft` se pueden eliminar (DELETE)
- [ ] Documento publicado no se puede editar (secciones ni clases) → debe volverse a draft o archivarse
- [ ] DELETE en documento published → 403

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.5.1 | [Usecase: publicar documento](./tareas/T-4.5.1-usecase-publicar.md) | src/core/usecases/ | ⬜ |
| 4.5.2 | [Usecase: archivar y eliminar](./tareas/T-4.5.2-usecase-archivar-eliminar.md) | src/core/usecases/ | ⬜ |
| 4.5.3 | [Tests de estados](./tareas/T-4.5.3-tests-estados.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.1: Modelo de datos](../HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) — Campo status enum
- [HU-4.3: Secciones](../HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) — Validar secciones requeridas

## Diseño técnico

### Máquina de estados

```
[draft] ──(publicar)──→ [published] ──(archivar)──→ [archived]
```

### Validaciones al publicar

1. **Topics distribuidos:** Cada topic en `coord_doc_topics` debe aparecer en al menos un `coord_doc_subject_topics`. Si hay topics sin asignar → 422 con lista de topics faltantes.

2. **Secciones requeridas:** Cada sección con `required: true` en `config.coord_doc_sections` debe tener `value` no vacío en `sections` JSONB. Si faltan → 422 con lista de secciones faltantes.

### Query: topics no distribuidos

```sql
SELECT t.id, t.name
FROM coord_doc_topics cdt
JOIN topics t ON t.id = cdt.topic_id
WHERE cdt.coordination_document_id = $1
  AND cdt.topic_id NOT IN (
    SELECT cdst.topic_id
    FROM coord_doc_subject_topics cdst
    JOIN coordination_document_subjects cds ON cds.id = cdst.coord_doc_subject_id
    WHERE cds.coordination_document_id = $1
  );
```

## Test cases

- 4.19: Publicar con todos los topics distribuidos → published
- 4.20: Publicar con topics sin distribuir → 422 con detalle
- 4.21: Publicar con sección requerida vacía → 422
- 4.22: DELETE en draft → ok
- 4.23: DELETE en published → 403
- 4.24: Archivar documento published → archived
- 4.25: Docente puede ver documento published → 200
