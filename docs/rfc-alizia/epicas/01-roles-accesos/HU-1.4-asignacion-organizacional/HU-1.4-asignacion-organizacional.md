# HU-1.4: Asignación organizacional

> Como admin, necesito asignar coordinadores a áreas para que cada persona vea solo los datos relevantes a su rol dentro de la organización.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

> **Nota:** La tabla `course_subjects` se define en [HU-3.4](../../03-integracion/HU-3.4-cursos-alumnos-asignaciones/HU-3.4-cursos-alumnos-asignaciones.md). Esta HU solo gestiona la asignación de coordinadores a áreas.

## Criterios de aceptación

- [ ] Tabla `area_coordinators` permite asignar coordinadores a áreas (M2M)
- [ ] Endpoint admin para asignar coordinador a un área
- [ ] Un coordinador solo ve documentos de sus áreas asignadas
- [ ] Multi-tenancy: todas las queries filtran por `organization_id`
- [ ] Seed incluye asignaciones de ejemplo

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.4.1 | [Migración: area_coordinators](./tareas/T-1.4.1-migracion.md) | db/migrations/ | ⬜ |
| 1.4.2 | [Entities y providers para area_coordinators](./tareas/T-1.4.2-entities-providers.md) | internal/admin/ | ⬜ |
| 1.4.3 | [Endpoints admin de asignación de coordinadores](./tareas/T-1.4.3-endpoints-admin.md) | internal/admin/entrypoints/ | ⬜ |
| 1.4.4 | [Filtrado por asignación en queries](./tareas/T-1.4.4-filtrado-asignacion.md) | internal/*/repositories/ | ⬜ |
| 1.4.5 | [Tests de asignación y filtrado](./tareas/T-1.4.5-tests-asignacion.md) | *_test.go | ⬜ |

## Modelo de datos

```
area_coordinators (id, area_id FK, user_id FK, UNIQUE(area_id, user_id))
```

## Dependencias

- HU-1.2 completada (users y organizations en DB)
- HU-1.3 completada (RequireRole middleware para proteger endpoints admin)
- Épica 3 (integración): areas ya deben existir en DB

## Test cases

- Admin asigna coordinador a área → 201, registro creado
- Admin asigna coordinador duplicado → 409 `already_assigned`
- Coordinador consulta documentos → solo ve los de sus áreas
- Request desde otra org → no ve asignaciones ajenas
