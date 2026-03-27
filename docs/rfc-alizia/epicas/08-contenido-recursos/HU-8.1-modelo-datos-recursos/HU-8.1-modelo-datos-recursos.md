# HU-8.1: Modelo de datos de recursos

> Como equipo de desarrollo, necesito las tablas normalizadas para los tipos de recurso, recursos generados y la configuración por organización.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Alta (bloqueante para todo lo demás de esta épica)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `fonts` con: id, organization_id, name, description, file_url, file_type, thumbnail_url, area_id, is_validated, created_at
- [ ] Tabla `resource_types` con: id, slug, name, description, prompt, output_schema (JSONB), requires_font (bool), is_active (bool), organization_id (nullable), created_at
- [ ] Tabla `organization_resource_types` (junction: org ↔ resource_type + custom_prompt + custom_output_schema)
- [ ] Tabla `resources` con: id, organization_id, resource_type_id, course_subject_id, title, content (JSONB), font_id (FK fonts, nullable), status, created_by, created_at, updated_at
- [ ] Enum `resource_status`: draft, generated, published
- [ ] Tipos públicos (organization_id NULL) visibles para todas las orgs
- [ ] Tipos privados (organization_id set) visibles solo para esa org
- [ ] Entities Go con GORM tags y relaciones
- [ ] Provider interfaces para CRUD
- [ ] Repository GORM con preloads

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.1.1 | [Migración: tablas de recursos](./tareas/T-8.1.1-migracion.md) | db/migrations/ | ⬜ |
| 8.1.2 | [Entities y providers](./tareas/T-8.1.2-entities-providers.md) | src/core/ | ⬜ |
| 8.1.3 | [Repository GORM](./tareas/T-8.1.3-repository.md) | src/repositories/ | ⬜ |
| 8.1.4 | [Tests](./tareas/T-8.1.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../../03-integracion/HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — FK organization_id
- [HU-3.4: Cursos y asignaciones](../../03-integracion/HU-3.4-cursos-alumnos-asignaciones/HU-3.4-cursos-alumnos-asignaciones.md) — FK course_subject_id
- [HU-1.2: Modelo de usuarios y roles](../../01-roles-accesos/HU-1.2-modelo-usuarios-roles/HU-1.2-modelo-usuarios-roles.md) — FK created_by
- [HU-3.2: Áreas y materias](../../03-integracion/HU-3.2-areas-materias/HU-3.2-areas-materias.md) — FK area_id en fonts

## Diseño técnico

### Modelo normalizado

```
fonts (fuentes educativas oficiales por org)

resource_types (tipos globales + privados por org)
  └── organization_resource_types (activación + customización por org)

resources (recursos generados)
  ├── resource_type_id → resource_types
  ├── font_id → fonts (fuente usada, si aplica)
  ├── course_subject_id → course_subjects (contexto de materia/curso)
  └── created_by → users
```

## Test cases

- 8.1: Crear resource_type público → visible para todas las orgs
- 8.2: Crear resource_type privado → solo visible para la org dueña
- 8.3: Crear resource con content JSONB → almacena y retorna correctamente
- 8.4: GET resource con preloads → incluye type, course_subject, creator
