# HU-4.2: Wizard de creación

> Como coordinador, necesito un wizard de 3 pasos para crear un documento de coordinación seleccionando topics, definiendo el período y asignando topics a cada disciplina.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `POST /api/v1/coordination-documents` acepta los datos del wizard completo
- [ ] Paso 1: seleccionar topics al nivel definido por `config.topic_selection_level`
- [ ] Paso 2: definir nombre, período (start_date, end_date) y class_count por disciplina
- [ ] Paso 3: asignar topics seleccionados a cada disciplina del área
- [ ] El class_count se calcula automáticamente desde la grilla horaria (time_slots) **si existe**. Si no hay grilla cargada, el coordinador ingresa manualmente (P3)
- [ ] **class_count mínimo: 1** — si una materia no tiene clases, no se agrega al documento (P5)
- [ ] El coordinador puede override el class_count manualmente (± feriados) con campo opcional `observations`
- [ ] **Topics pueden repetirse entre disciplinas** — la interdisciplinariedad es un pilar del modelo (P4)
- [ ] El documento se crea en estado `pending` — transición automática a `in_progress` al primer edit/interacción (P1)
- [ ] Se crean todas las junction tables: coord_doc_topics, coordination_document_subjects, coord_doc_subject_topics
- [ ] Solo coordinadores del área pueden crear documentos para esa área

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.2.1 | [Usecase: crear documento](./tareas/T-4.2.1-usecase-crear-documento.md) | src/core/usecases/ | ⬜ |
| 4.2.2 | [Cálculo automático de class_count](./tareas/T-4.2.2-calculo-class-count.md) | src/core/usecases/ | ⬜ |
| 4.2.3 | [Endpoint POST y handler](./tareas/T-4.2.3-endpoint-crear.md) | src/entrypoints/ | ⬜ |
| 4.2.4 | [Endpoint GET listar y detalle](./tareas/T-4.2.4-endpoints-listar-detalle.md) | src/entrypoints/ | ⬜ |
| 4.2.5 | [Seed y tests](./tareas/T-4.2.5-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-4.1: Modelo de datos](../HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) — Tablas deben existir
- [HU-3.5: Grilla horaria](../../03-integracion/HU-3.5-grilla-horaria-clases-compartidas/HU-3.5-grilla-horaria-clases-compartidas.md) — Para calcular class_count automático (opcional)
- [HU-1.4: Asignación organizacional](../../01-roles-accesos/HU-1.4-asignacion-organizacional/HU-1.4-asignacion-organizacional.md) — Verificar que el usuario es coordinador del área

## Diseño técnico

### Request POST (wizard completo)

```json
{
  "name": "Itinerario Ciencias Exactas - 1er cuatrimestre",
  "area_id": 1,
  "start_date": "2026-03-15",
  "end_date": "2026-07-15",
  "topic_ids": [5, 8, 12, 15],
  "subjects": [
    {
      "subject_id": 1,
      "class_count": 20,
      "topic_ids": [5, 8],
      "observations": ""
    },
    {
      "subject_id": 2,
      "class_count": 18,
      "topic_ids": [12, 15],
      "observations": "Se descuentan 2 clases por jornadas institucionales"
    }
  ]
}
```

### Validaciones

- `class_count >= 1` para cada subject (P5)
- `topic_ids` no vacío a nivel documento
- Cada `subject_id` pertenece al área indicada
- Cada `topic_id` en subjects debe estar incluido en `topic_ids` del documento
- Topics pueden repetirse entre subjects (P4)

### Cálculo automático de class_count

```
class_count = (slots por semana de la disciplina) × (semanas en el período)
semanas = (end_date - start_date) / 7
```

Si hay grilla horaria → pre-llenar con el cálculo automático.
Si no hay grilla → campo vacío, coordinador ingresa a mano (P3).

El frontend muestra el valor calculado (si existe); el coordinador puede modificarlo.

### Flujo del wizard

```
Paso 1: GET /api/v1/topics?level=3  →  Seleccionar topics
                                          │
Paso 2: GET /api/v1/areas/:id       →  Ver disciplinas del área
         + calcular class_count         Definir período + override class_count
         (si hay grilla)                  │
                                          │
Paso 3: Asignar topics a disciplinas   →  Drag & drop topics → subjects
         (repeticiones permitidas)        │
                                          │
POST /api/v1/coordination-documents  →  Crear documento (pending)
```

## Test cases

- 4.4: POST wizard completo → documento creado en estado pending con topics y subjects
- 4.5: POST con topic que no existe → 422
- 4.6: POST con disciplina que no es del área → 422
- 4.7: POST sin ser coordinador del área → 403
- 4.8: POST con class_count = 0 → 422 (mínimo 1)
- 4.9: POST con topic repetido entre disciplines → ok (interdisciplinariedad)
- 4.10: POST sin grilla horaria → ok, class_count manual
- 4.11: GET listar por área → solo docs de la org y área
- 4.12: GET detalle → documento con todos los preloads
