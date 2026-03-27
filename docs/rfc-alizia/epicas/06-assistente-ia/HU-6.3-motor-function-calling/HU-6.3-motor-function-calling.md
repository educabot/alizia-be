# HU-6.3: Motor de function calling

> Como desarrollador, necesito un framework de function calling para que las features de chat con IA puedan definir tools, ejecutarlas y manejar conversaciones multi-turn con tool calls de forma consistente.

**Fase:** 4 — Asistente IA
**Prioridad:** Alta (bloqueante para HU-4.6 y cualquier chat con tools)
**Estimacion:** —

---

## Criterios de aceptacion

- [ ] Tool registry: registrar tools con nombre, descripcion, parametros (JSON Schema) y handler function
- [ ] Tool execution loop: enviar mensajes → si hay tool_calls → ejecutar handlers → agregar resultados → repetir
- [ ] Maximo 5 iteraciones por turno de conversacion (safety cap)
- [ ] Si un tool handler retorna error → enviar el error como tool result al LLM (que informe al usuario)
- [ ] Los resultados de tools se serializan como strings JSON
- [ ] Soporte para multiples tool calls en una sola respuesta (parallel tool calling)
- [ ] Reutilizable en distintos contextos (coord doc chat, lesson plan chat, etc.)
- [ ] Cada contexto registra su propio set de tools

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.3.1 | [Tool registry](./tareas/T-6.3.1-tool-registry.md) | src/core/usecases/ai/tool_registry.go | ⬜ |
| 6.3.2 | [Tool execution loop](./tareas/T-6.3.2-tool-execution-loop.md) | src/core/usecases/ai/tool_loop.go | ⬜ |
| 6.3.3 | [Context-specific tool sets](./tareas/T-6.3.3-context-tool-sets.md) | src/core/usecases/ai/tool_sets.go | ⬜ |
| 6.3.4 | [Tests](./tareas/T-6.3.4-tests.md) | tests/ai/tool_test.go | ⬜ |

## Dependencias

- HU-6.1 completada (AI Provider con metodo `ChatWithTools` y tipos `ChatMessage`, `ToolCall`, `ToolDefinition`)

## Test cases

- 6.11: Registrar tool → disponible en el registry via `GetDefinitions`
- 6.12: Single tool call → handler ejecutado → resultado enviado al LLM → respuesta final
- 6.13: Multiples tool calls paralelos → todos los handlers ejecutados → resultados enviados
- 6.14: Tool handler con error → error enviado como resultado al LLM → el LLM informa al usuario gracefully
- 6.15: Loop alcanza max iteraciones → retorna ultimo mensaje con warning
