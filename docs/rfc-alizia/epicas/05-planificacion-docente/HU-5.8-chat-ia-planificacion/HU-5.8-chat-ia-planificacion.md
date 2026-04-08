# HU-5.8: Chat con Alizia para planificación docente

> Como docente, necesito chatear con Alizia para editar mi planificación por lenguaje natural, sin tener que buscar y editar cada campo manualmente.

**Fase:** Fase 5
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `POST /api/v1/lesson-plans/:id/chat` acepta mensaje + historial
- [ ] Alizia puede modificar la planificación via function calling (tools)
- [ ] Tools disponibles: `update_activity_content`, `swap_activity`, `update_objective`, `update_font`
- [ ] El historial de chat se envía desde el frontend en cada request (no se persiste en backend)
- [ ] Alizia tiene contexto de la planificación completa (momentos, actividades, contenido, fuentes, documento de coordinación)
- [ ] Las modificaciones via chat se aplican inmediatamente a la planificación
- [ ] Solo funciona en planificaciones `in_progress` (no pending ni published)

## Concepto

Análogo a HU-4.6 (Chat con Alizia para documento de coordinación) pero aplicado a la planificación docente. El docente puede pedir cambios como "cambiá la actividad del cierre por algo más dinámico" y Alizia modifica la sección correspondiente.

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 5.8.1 | Definición de tools (function calling) | src/core/usecases/ | ⬜ |
| 5.8.2 | Usecase: chat con contexto de planificación | src/core/usecases/ | ⬜ |
| 5.8.3 | Endpoint POST chat | src/entrypoints/ | ⬜ |
| 5.8.4 | Tests | tests/ | ⬜ |

## Dependencias

- [HU-5.3: Planificación por momentos](../HU-5.3-planificacion-por-momentos/HU-5.3-planificacion-por-momentos.md) — Estructura de momentos y actividades
- [HU-5.4: Generación con IA](../HU-5.4-generacion-ia/HU-5.4-generacion-ia.md) — Contexto de generación
- [HU-4.6: Chat con Alizia](../../04-documento-coordinacion/HU-4.6-chat-alizia/HU-4.6-chat-alizia.md) — Patrón de chat con tools
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Azure OpenAI con function calling
