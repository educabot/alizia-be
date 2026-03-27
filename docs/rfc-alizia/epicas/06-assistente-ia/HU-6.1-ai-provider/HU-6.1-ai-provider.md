# HU-6.1: AI Provider

> Como desarrollador, necesito un AI Provider reutilizable que abstraiga la integración con Azure OpenAI, para que todas las features puedan generar texto, JSON estructurado y chatear con tools a través de una interfaz consistente.

**Fase:** 4 — Asistente IA
**Prioridad:** Alta (bloqueante para todas las features que usan IA)
**Estimacion:** —

---

## Criterios de aceptacion

- [ ] Interfaz `AIProvider` con 3 metodos: `GenerateText`, `GenerateJSON`, `ChatWithTools`
- [ ] Implementacion con Azure OpenAI usando el SDK oficial de Go (`azopenai`)
- [ ] Configuracion via variables de entorno (`AZURE_OPENAI_ENDPOINT`, `AZURE_OPENAI_API_KEY`)
- [ ] Modelo: `gpt-5-mini` con `max_completion_tokens` (no `max_tokens`)
- [ ] Retry con backoff exponencial (3 intentos) para errores de rate limit (429)
- [ ] Errores tipados: `ErrAIRateLimit`, `ErrAIContextLength`, `ErrAIUnavailable`
- [ ] Logging de uso de tokens por request (prompt_tokens, completion_tokens)
- [ ] Timeout configurable (default 30s)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.1.1 | [Definicion de interfaz del provider](./tareas/T-6.1.1-provider-interface.md) | src/core/providers/ai_provider.go | ⬜ |
| 6.1.2 | [Implementacion Azure OpenAI](./tareas/T-6.1.2-azure-openai-impl.md) | src/infra/ai/azure_openai.go | ⬜ |
| 6.1.3 | [Configuracion e inicializacion](./tareas/T-6.1.3-config-init.md) | src/infra/ai/config.go, src/core/providers/errors.go | ⬜ |
| 6.1.4 | [Tests](./tareas/T-6.1.4-tests.md) | tests/ai/*_test.go | ⬜ |

## Dependencias

- Libreria compartida `team-ai-toolkit` (config helpers: `EnvOr`, `MustEnv`)
- SDK de Azure OpenAI para Go (`github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai`)

## Test cases

- 6.1: `GenerateText` con prompt valido → retorna string no vacio
- 6.2: `GenerateJSON` con schema → retorna JSON valido que matchea el schema
- 6.3: `ChatWithTools` con tool calls → ejecuta handler, retorna respuesta final
- 6.4: Rate limit → retry con backoff → exito en el reintento
- 6.5: Context length excedido → retorna `ErrAIContextLength`
