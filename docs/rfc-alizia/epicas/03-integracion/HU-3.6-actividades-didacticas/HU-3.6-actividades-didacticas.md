# HU-3.6: Actividades didácticas

> Como admin, necesito cargar las actividades didácticas disponibles por momento de clase para que los docentes las seleccionen al planificar.

**Fase:** 2 — Admin/Integration
**Prioridad:** Media (necesaria para Épica 5)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `activities` con: id, organization_id, moment (enum: apertura, desarrollo, cierre), name, description, duration_minutes, created_at
- [ ] Enum `class_moment` creado: 'apertura', 'desarrollo', 'cierre'
- [ ] Endpoint POST para crear actividades (admin)
- [ ] Endpoint GET para listar actividades (filtrable por momento)
- [ ] Cada actividad pertenece a un momento fijo (no puede ser de dos momentos)
- [ ] Config de org define `desarrollo_max_activities` (máx actividades en momento desarrollo)
- [ ] Seed con actividades de ejemplo por momento

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 3.6.1 | [Migración: enum class_moment + tabla activities](./tareas/T-3.6.1-migracion.md) | db/migrations/ | ⬜ |
| 3.6.2 | [Entity y provider](./tareas/T-3.6.2-entity-provider.md) | src/core/ | ⬜ |
| 3.6.3 | [Repository y endpoints](./tareas/T-3.6.3-repository-endpoints.md) | src/repositories/, src/entrypoints/ | ⬜ |
| 3.6.4 | [Seed y tests](./tareas/T-3.6.4-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — FK organization_id, config con `desarrollo_max_activities`

## Diseño técnico

### Momentos de clase (enum fijo)

| Momento | Actividades permitidas | Descripción |
|---------|----------------------|-------------|
| `apertura` | Exactamente 1 | Inicio de la clase |
| `desarrollo` | 1 a `config.desarrollo_max_activities` (default 3) | Cuerpo principal |
| `cierre` | Exactamente 1 | Cierre y síntesis |

Los momentos son un enum fijo — no configurables por org. Lo que sí es configurable es cuántas actividades permite el momento de desarrollo.

### Ejemplo de actividades

```
Apertura:
  - Lluvia de ideas (15 min)
  - Pregunta disparadora (10 min)
  - Repaso clase anterior (10 min)

Desarrollo:
  - Trabajo en grupo (30 min)
  - Exposición dialogada (25 min)
  - Resolución de problemas (30 min)
  - Análisis de caso (25 min)

Cierre:
  - Puesta en común (15 min)
  - Síntesis individual (10 min)
  - Autoevaluación (10 min)
```

## Test cases

- 3.26: POST activity con momento válido → 201
- 3.27: GET activities → lista completa
- 3.28: GET activities?moment=desarrollo → solo actividades de desarrollo
- 3.29: Multi-tenancy: cada org ve solo sus actividades
