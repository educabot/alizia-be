# HU-1.4: Asignación organizacional

> Como admin, necesito asignar usuarios a áreas y cursos para que cada persona vea solo los datos relevantes a su rol dentro de la organización.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `area_coordinators` permite asignar coordinadores a áreas (M2M)
- [ ] Tabla `course_subjects` vincula curso + materia + docente
- [ ] Endpoint admin para asignar coordinador a un área
- [ ] Endpoint admin para asignar docente a curso-materia
- [ ] Un coordinador solo ve documentos de sus áreas asignadas
- [ ] Un docente solo ve los cursos/materias donde está asignado
- [ ] Multi-tenancy: todas las queries filtran por `organization_id`
- [ ] Seed incluye asignaciones de ejemplo

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.4.1 | [Migración: area_coordinators + course_subjects](./tareas/T-1.4.1-migracion.md) | db/migrations/ | ⬜ |
| 1.4.2 | [Entities y providers para asignaciones](./tareas/T-1.4.2-entities-providers.md) | internal/admin/ | ⬜ |
| 1.4.3 | [Endpoints admin de asignación](./tareas/T-1.4.3-endpoints-admin.md) | internal/admin/entrypoints/ | ⬜ |
| 1.4.4 | [Filtrado por asignación en queries](./tareas/T-1.4.4-filtrado-asignacion.md) | internal/*/repositories/ | ⬜ |
| 1.4.5 | [Tests de asignación y filtrado](./tareas/T-1.4.5-tests-asignacion.md) | *_test.go | ⬜ |

## Modelo de datos

```
area_coordinators (id, area_id FK, user_id FK, UNIQUE(area_id, user_id))
course_subjects (id, course_id FK, subject_id FK, teacher_id FK, organization_id FK, school_year, UNIQUE(course_id, subject_id, school_year))
```

## Dependencias

- HU-1.2 completada (users y organizations en DB)
- HU-1.3 completada (RequireRole middleware para proteger endpoints admin)
- Épica 3 (integración): areas, subjects, courses ya deben existir en DB

## Test cases

- Admin asigna coordinador a área → 201, registro creado
- Admin asigna coordinador duplicado → 409 `already_assigned`
- Coordinador consulta documentos → solo ve los de sus áreas
- Docente consulta cursos → solo ve donde está asignado
- Request desde otra org → no ve asignaciones ajenas
