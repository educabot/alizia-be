# HU-6.2: Prompts configurables

> Como desarrollador, necesito un sistema de templates de prompts que cargue los prompts desde la configuracion de la organizacion, reemplace placeholders con valores en runtime, y valide las salidas contra JSON Schemas.

**Fase:** 4 — Asistente IA
**Prioridad:** Alta (bloqueante para generacion de secciones y chat)
**Estimacion:** —

---

## Criterios de aceptacion

- [ ] Los templates de prompts se definen en `organizations.config` por tipo de seccion
- [ ] Sistema de placeholders: `{variable_name}` se reemplaza en runtime con valores reales
- [ ] Cada seccion generable tiene: `key`, `ai_prompt` (template), `output_schema` (JSON Schema opcional)
- [ ] El system prompt builder compone: definicion de rol + contexto + instruccion especifica
- [ ] Si un placeholder no tiene valor en runtime → error (no se omite silenciosamente)
- [ ] Los prompts soportan templates multi-linea con formato markdown
- [ ] Validacion de config: los campos requeridos en los templates de prompts se verifican al inicio

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.2.1 | [Prompt template engine](./tareas/T-6.2.1-prompt-template-engine.md) | src/core/usecases/ai/prompt_engine.go | ⬜ |
| 6.2.2 | [System prompt builder](./tareas/T-6.2.2-system-prompt-builder.md) | src/core/usecases/ai/system_prompt.go | ⬜ |
| 6.2.3 | [Tests](./tareas/T-6.2.3-tests.md) | tests/ai/prompt_test.go | ⬜ |

## Dependencias

- HU-3.1 completada (Organizaciones — config de org donde se almacenan los prompts)
- HU-6.1 completada (AI Provider — consume los prompts renderizados)

## Test cases

- 6.6: Template con placeholders → todos reemplazados correctamente
- 6.7: Template con placeholder sin valor → error listando las keys faltantes
- 6.8: System prompt construido con secciones de rol + contexto + instruccion
- 6.9: Validacion de JSON Schema sobre salida de IA → valida
- 6.10: Validacion de JSON Schema sobre salida de IA → invalida → error con detalles
