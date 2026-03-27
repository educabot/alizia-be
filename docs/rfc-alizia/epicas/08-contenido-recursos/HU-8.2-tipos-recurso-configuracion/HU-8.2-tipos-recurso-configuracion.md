# HU-8.2: Tipos de recurso y configuración por organización

> Como administrador, necesito definir qué tipos de recurso están disponibles para cada organización y poder customizar prompts y schemas por provincia.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Existen tipos de recurso públicos disponibles para todas las organizaciones (guía de lectura, ficha de cátedra, creación libre)
- [ ] Cada organización puede habilitar/deshabilitar tipos de recurso
- [ ] Cada organización puede sobreescribir el prompt y el output_schema de un tipo
- [ ] Los tipos privados solo son visibles para la organización que los creó
- [ ] Endpoint para listar tipos disponibles de una organización (públicos habilitados + privados)
- [ ] Endpoint para configurar tipos por organización (habilitar, customizar)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.2.1 | [Endpoints de tipos de recurso](./tareas/T-8.2.1-endpoints-tipos.md) | handlers/ | ⬜ |
| 8.2.2 | [Resolución de prompt y schema](./tareas/T-8.2.2-resolucion-prompt-schema.md) | usecases/ | ⬜ |
| 8.2.3 | [Tests](./tareas/T-8.2.3-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-8.1: Modelo de datos](../HU-8.1-modelo-datos-recursos/HU-8.1-modelo-datos-recursos.md) — Tablas base
- [HU-3.1: Organizaciones](../../03-integracion/HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — Config por org

## Test cases

- 8.5: Listar tipos para org → retorna públicos habilitados + privados de la org
- 8.6: Tipo deshabilitado → no aparece en la lista
- 8.7: Prompt custom → se resuelve correctamente sobre el default
