# HU-3.2: Áreas y materias

> Como admin, necesito crear áreas y materias para estructurar la oferta curricular de la institución.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `areas` con: id, organization_id, name, description, created_at
- [ ] Tabla `subjects` con: id, organization_id, area_id (FK), name, description, created_at
- [ ] Tabla `area_coordinators` con: id, area_id, user_id, UNIQUE(area_id, user_id)
- [ ] CRUD endpoints para áreas: POST, GET (listar), PUT
- [ ] CRUD endpoints para materias: POST, GET (listar)
- [ ] Endpoint para asignar coordinador a área: `POST /api/v1/areas/:id/coordinators`
- [ ] Listar áreas incluye las materias asociadas (preload)
- [ ] Un coordinador puede coordinar múltiples áreas (M2M via area_coordinators)
- [ ] Múltiples coordinadores pueden coordinar la misma área
- [ ] Todos los queries filtran por organization_id del JWT

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 3.2.1 | [Migración: areas, subjects, area_coordinators](./tareas/T-3.2.1-migracion.md) | db/migrations/ | ⬜ |
| 3.2.2 | [Entities](./tareas/T-3.2.2-entities.md) | src/core/entities/ | ⬜ |
| 3.2.3 | [Providers](./tareas/T-3.2.3-providers.md) | src/core/providers/ | ⬜ |
| 3.2.4 | [Repository GORM](./tareas/T-3.2.4-repository.md) | src/repositories/ | ⬜ |
| 3.2.5 | [Endpoints y handlers](./tareas/T-3.2.5-endpoints-handlers.md) | src/entrypoints/ | ⬜ |
| 3.2.6 | [Seed y tests](./tareas/T-3.2.6-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — organizations debe existir
- [HU-1.2: Modelo de usuarios](../../01-roles-accesos/HU-1.2-modelo-usuarios-roles/HU-1.2-modelo-usuarios-roles.md) — users debe existir para area_coordinators
- [HU-1.3: Middleware de autorización](../../01-roles-accesos/HU-1.3-middleware-autorizacion/HU-1.3-middleware-autorizacion.md) — RequireRole(admin, coordinator) para endpoints

## Test cases

- 3.4: POST area → 201 con area creada
- 3.5: POST subject con area_id → materia asociada al área
- 3.6: POST area coordinator → asignación exitosa
- 3.7: GET areas → incluye subjects y coordinators preloaded
- 3.8: Duplicar coordinador en misma área → 409 conflict
