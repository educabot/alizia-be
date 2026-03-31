# HU-5.5: Estados de planificación

> Como docente, necesito marcar mi planificación como lista para que sea visible al coordinador, y como coordinador necesito ver el progreso de planificación de mi equipo.

**Fase:** 5 — Planificación docente
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Estados: `pending` → `in_progress` → `published`
- [ ] El docente marca como published cuando termina de editar
- [ ] Endpoint `PATCH /api/v1/lesson-plans/:id/status` cambia el estado
- [ ] El coordinador puede ver el estado de todas las planificaciones de su área
- [ ] Solo el docente asignado puede cambiar el estado de su plan
- [ ] Para marcar como `published` se requiere: actividades configuradas + proposal generada

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 5.5.1 | [Usecase: cambiar estado](./tareas/T-5.5.1-usecase-cambiar-estado.md) | src/core/usecases/ | ⬜ |
| 5.5.2 | [Endpoints de estado y vista coordinador](./tareas/T-5.5.2-endpoints.md) | src/entrypoints/ | ⬜ |
| 5.5.3 | [Tests](./tareas/T-5.5.3-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-5.1: Modelo de datos](../HU-5.1-modelo-datos-planificacion/HU-5.1-modelo-datos-planificacion.md) — Enum lesson_plan_status
- [HU-5.4: Generación con IA](../HU-5.4-generacion-ia/HU-5.4-generacion-ia.md) — Proposal debe existir para marcar published

## Diseño técnico

### Máquina de estados

```
[pending] ──→ [in_progress] ──→ [published]
```

El flujo es unidireccional: el docente avanza de pending a in_progress cuando comienza a editar, y a published cuando finaliza.

### Validaciones al marcar published

1. **Actividades configuradas:** Al menos 1 apertura, 1+ desarrollo, 1 cierre
2. **Proposal generada:** El campo `proposal` no debe ser null/vacío
3. **Solo el docente asignado** puede cambiar estado

### Vista coordinador

```
GET /api/v1/coordination-documents/:id/planning-progress
```

Response:
```json
{
  "coordination_document_id": 5,
  "subjects": [
    {
      "subject_id": 1,
      "subject_name": "Matemáticas",
      "teacher": "Doc. García",
      "total_classes": 20,
      "published_classes": 12,
      "pending_classes": 5,
      "not_started": 3,
      "progress_percentage": 60
    }
  ],
  "overall_progress": 55
}
```

## Test cases

- 5.22: Marcar published con proposal → status updated
- 5.23: Marcar published sin proposal → 422
- 5.24: Transición pending → in_progress → ok
- 5.25: Otro docente intenta cambiar estado → 403
- 5.26: GET planning-progress como coordinador → stats correctas
