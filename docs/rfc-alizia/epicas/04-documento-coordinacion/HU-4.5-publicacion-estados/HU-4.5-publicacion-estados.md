# HU-4.5: Publicación y estados

> Como coordinador, necesito publicar el documento para que los docentes lo vean, y archivar documentos viejos.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Estados: `pending` → `in_progress` → `published` + `archived`
- [ ] Transición `pending → in_progress` es **automática** al primer click/edición/interacción (P1)
- [ ] Solo coordinadores pueden cambiar estado
- [ ] Sin restricción de cantidad de documentos por área — múltiples pueden coexistir en cualquier estado (P2)
- [ ] Al publicar, se valida que cada sub-topic esté asignado a al menos 1 materia. Si hay sub-topics sin asignar → **warning confirmable** (no bloqueante): "Te quedó tal sub-topic sin asignar, ¿querés continuar?" (P13)
- [ ] Al publicar, se valida que todas las secciones requeridas tengan contenido
- [ ] Documento publicado es visible para docentes (GET listar y detalle)
- [ ] Docentes **no ven** documentos no publicados (P14)
- [ ] Solo documentos en `pending` se pueden eliminar (DELETE)
- [ ] Documento publicado: **secciones editables** con warning "Los cambios no se propagan a planificaciones ya creadas" (P12)
- [ ] Documento publicado: **clases inmutables** — solo el teacher edita en su planificación (P10/P12)
- [ ] DELETE en documento published → 403
- [ ] Endpoint para archivar: `PATCH` con `status: "archived"` (desde in_progress o published)
- [ ] La propagación automática de cambios a lesson plans existentes es **post-MVP**

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.5.1 | [Usecase: publicar documento](./tareas/T-4.5.1-usecase-publicar.md) | src/core/usecases/ | ⬜ |
| 4.5.2 | [Usecase: archivar, auto-transición y eliminar](./tareas/T-4.5.2-usecase-archivar-eliminar.md) | src/core/usecases/ | ⬜ |
| 4.5.3 | [Tests de estados](./tareas/T-4.5.3-tests-estados.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.1: Modelo de datos](../HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) — Campo status enum
- [HU-4.3: Secciones](../HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) — Validar secciones requeridas

## Diseño técnico

### Máquina de estados

```
[pending] ──(auto: primera interacción)──→ [in_progress] ──(publicar)──→ [published]
                                                │                              │
                                                └──(archivar)──→ [archived] ←──┘
```

**Transiciones válidas:**
- `pending → in_progress`: automática al primer edit/click/interacción
- `in_progress → published`: manual, con validación
- `in_progress → archived`: manual
- `published → archived`: manual

**Transiciones inválidas:** no hay camino de retorno (no se puede "despublicar")

### Auto-transición pending → in_progress (P1 — Decisión)

El estado `pending` existe para marcar documentos recién creados que nunca fueron tocados. La transición a `in_progress` es **automática** y ocurre cuando el coordinador:
- Edita cualquier sección
- Interactúa con el chat de Alizia
- Hace cualquier modificación al documento

El backend detecta que el documento está en `pending` y lo transiciona antes de aplicar el cambio:

```go
if doc.Status == "pending" {
    doc.Status = "in_progress"
    // actualizar en la misma transacción
}
```

### Sin restricción por área (P2 — Decisión)

Múltiples documentos pueden coexistir en cualquier estado para la misma área. El frontend maneja la visualización (etiquetas, filtros, ordenamiento) para que los usuarios naveguen entre documentos.

### Validaciones al publicar (P13 — Decisión)

#### 1. Sub-topics asignados a materia — WARNING confirmable

Recordar: un documento de coordinación está relacionado a un **topic padre** y las materias a **topics hijos** (sub-topics). Cada sub-topic del documento debe estar asignado a al menos 1 materia.

Si hay sub-topics sin asignar, se retorna un **warning que el usuario puede confirmar**:

```json
{
  "warnings": [
    {
      "type": "unassigned_subtopics",
      "message": "Te quedaron sub-topics sin asignar a ninguna materia",
      "details": [
        {"id": 5, "name": "Ecuaciones lineales"},
        {"id": 8, "name": "Funciones"}
      ]
    }
  ],
  "requires_confirmation": true
}
```

El frontend muestra: "Te quedó tal sub-topic sin asignar, ¿querés continuar?" con botón Confirmar/Cancelar.

Si el coordinador confirma, se re-envía con `"force": true`:
```json
POST /api/v1/coordination-documents/:id/publish
{ "force": true }
```

#### 2. Secciones requeridas con contenido — BLOQUEANTE

Cada sección con `required: true` en `config.coord_doc_sections` debe tener `value` no vacío. Si faltan → 422 con lista de secciones faltantes.

### Query: sub-topics no asignados

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

### Editabilidad post-publicación (P12 — Decisión)

| Elemento | in_progress | published |
|----------|-------------|-----------|
| Secciones narrativas | Editable | Editable + warning |
| Plan de clases | Editable | **Inmutable** (403) |
| Chat con Alizia | Todos los tools | Solo tools de secciones |

El warning en el response de edición de secciones publicadas:
```json
{
  "data": { ... },
  "warning": "Los cambios no se propagan automáticamente a planificaciones docentes ya creadas"
}
```

## Test cases

- 4.19: Auto-transición: editar documento pending → status cambia a in_progress automáticamente
- 4.20: Publicar con todos los sub-topics asignados → published
- 4.21: Publicar con sub-topics sin asignar → warning confirmable con lista
- 4.22: Publicar con force=true después de warning → published
- 4.23: Publicar con sección requerida vacía → 422 bloqueante
- 4.24: DELETE en pending → ok
- 4.25: DELETE en published → 403
- 4.26: Archivar in_progress → archived
- 4.27: Archivar published → archived
- 4.28: Docente puede ver documento published → 200
- 4.29: Docente no ve documento in_progress → filtrado
- 4.30: PATCH sección en published → 200 con warning
- 4.31: PATCH clase en published → 403 (inmutable)
- 4.32: Múltiples documentos published para misma área → todos visibles
