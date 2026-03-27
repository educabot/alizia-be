# HU-6.4: Contexto y enriquecimiento

> Como desarrollador, necesito funciones de context builder que transformen entidades del dominio en texto enriquecido para prompts de LLM, para que las generaciones de IA tengan toda la informacion relevante disponible.

**Fase:** 4 — Asistente IA
**Prioridad:** Alta
**Estimacion:** —

---

## Criterios de aceptacion

- [ ] Context builder para documentos de coordinacion: incluye area, periodo, topics, materias con topics, secciones existentes
- [ ] Context builder para planificaciones: incluye info de la clase, actividades por momento, fuente, topics, estrategia del doc de coordinacion, contexto de clase compartida
- [ ] Context builder para chat: construye un resumen del documento completo para el system prompt (secciones, materias, plan de clases)
- [ ] Todos los builders retornan strings (no structs) — optimizados para consumo del LLM
- [ ] El contexto se estructura con headers markdown para claridad
- [ ] Cuando una entidad no tiene datos para un campo, se omite (no se pone "N/A")

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.4.1 | [Context builder coordinacion](./tareas/T-6.4.1-context-builder-coordinacion.md) | src/core/usecases/ai/context_coord_doc.go | ⬜ |
| 6.4.2 | [Context builder planificacion](./tareas/T-6.4.2-context-builder-planificacion.md) | src/core/usecases/ai/context_lesson_plan.go | ⬜ |
| 6.4.3 | [Tests](./tareas/T-6.4.3-tests.md) | tests/ai/context_test.go | ⬜ |

## Dependencias

- HU-4.1 completada (Modelo de datos documento — entidades usadas como input)
- HU-5.1 completada (Modelo de datos planificacion — entidades usadas como input)
- HU-6.2 completada (System prompt builder — el contexto se agrega via AddContext)

## Test cases

- 6.19: BuildCoordDocContext con documento completo → incluye area, topics, materias, secciones
- 6.20: BuildCoordDocContext con secciones vacias → omite bloque de secciones
- 6.21: BuildLessonPlanContext con clase compartida → incluye info de la otra materia
- 6.22: BuildLessonPlanContext sin fuente → omite bloque de fuente
- 6.23: BuildChatContext → resumen del documento con todas las secciones y planes de clase
