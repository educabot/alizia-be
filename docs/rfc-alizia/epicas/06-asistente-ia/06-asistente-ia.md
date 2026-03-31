# Épica 6: Asistente IA

> Infraestructura de inteligencia artificial compartida: provider, prompts, function calling y contexto.

**Estado:** MVP
**Fase de implementación:** Fase 4

---

## Problema

Múltiples épicas necesitan generar texto, estructurar respuestas en JSON, ejecutar tools via function calling, y armar contexto rico para los prompts. Sin una infraestructura compartida, cada feature reimplementa la integración con Azure OpenAI, los prompts quedan hardcodeados y no son configurables por provincia, y no hay un patrón consistente para el function calling.

## Objetivos

- Proveer un AI Provider reutilizable con interfaz clara (texto, JSON, chat con tools)
- Centralizar la gestión de prompts configurables por organización
- Estandarizar el patrón de function calling con un motor de ejecución de tools
- Construir context builders que enriquezcan los prompts con datos relevantes del dominio
- Permitir customización del comportamiento IA por organización (tono, límites, fuentes)

## Alcance MVP

**Incluye:**

- AI Provider con Azure OpenAI (GenerateText, GenerateJSON, ChatWithTools)
- Sistema de prompts configurables por org con placeholders y JSON Schema
- Motor de function calling: definición de tools, ejecución en loop, multi-turn
- Context builders para documentos de coordinación y planificaciones

- Customización por organización (tono, límites, fuentes permitidas)
- Asistencia de navegación contextual dentro de la plataforma

**No incluye:**

- Entrenamiento de modelos custom → fuera de alcance
- RAG con documentos curriculares → horizonte
- Voice-to-text para bitácora → ver Épica 5 post-MVP

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-6.1 | [AI Provider](./HU-6.1-ai-provider/HU-6.1-ai-provider.md) | Cliente Azure OpenAI con interfaz reutilizable | Fase 4 | 4 |
| HU-6.2 | [Prompts configurables](./HU-6.2-prompts-configurables/HU-6.2-prompts-configurables.md) | Templates de prompts por org con placeholders y JSON Schema | Fase 4 | 3 |
| HU-6.3 | [Motor de function calling](./HU-6.3-motor-function-calling/HU-6.3-motor-function-calling.md) | Framework para definir y ejecutar tools en conversations | Fase 4 | 4 |
| HU-6.4 | [Contexto y enriquecimiento](./HU-6.4-contexto-enriquecimiento/HU-6.4-contexto-enriquecimiento.md) | Builders de contexto para diferentes entidades del dominio | Fase 4 | 3 |
| HU-6.5 | [Customización por organización](./HU-6.5-customizacion-organizacion/HU-6.5-customizacion-organizacion.md) | Ajustar tono, límites y fuentes por provincia | Fase 4 | 5 |
| HU-6.6 | [Asistencia de navegación](./HU-6.6-asistencia-navegacion/HU-6.6-asistencia-navegacion.md) | Ayuda contextual dentro de la plataforma | Fase 4 | 4 |

---

## Decisiones técnicas

- Se usa **Azure OpenAI** (modelo `gpt-5-mini`) con `max_completion_tokens` en vez de `max_tokens`. No se customiza `temperature`.
- El AI Provider expone tres métodos: `GenerateText` (texto libre), `GenerateJSON` (structured output con schema), `ChatWithTools` (conversation con function calling).
- Los prompts se almacenan en `organizations.config` como templates con placeholders (`{variable}`). Cada sección generada por IA tiene su propio prompt.
- El motor de function calling sigue un loop: enviar mensaje → si hay tool_calls → ejecutar → enviar resultados → repetir hasta respuesta final. Máximo 5 iteraciones por seguridad.
- Los context builders son funciones puras que transforman entidades del dominio en strings optimizados para el prompt.
- Se usa `team-ai-toolkit` como shared Go library cuando aplique.

## Decisiones de cada cliente

- Los prompts por sección son configurables por organización
- El tono del asistente (formal, informal, rioplatense) se define por provincia
- Las fuentes curriculares permitidas varían por jurisdicción

## Épicas relacionadas

- **[Épica 4: Documento de coordinación](../04-documento-coordinacion/04-documento-coordinacion.md)** — Consume GenerateText (HU-4.3), GenerateJSON (HU-4.4), ChatWithTools (HU-4.6)
- **[Épica 5: Planificación docente](../05-planificacion-docente/05-planificacion-docente.md)** — Consume GenerateText para propuestas de clase (HU-5.4)
- **[Épica 8: Contenido y recursos](../08-contenido-recursos/08-contenido-recursos.md)** — Consumirá GenerateText para generación de recursos didácticos

## Test cases asociados

- Fase 4: Tests 6.1–6.18 (AI provider, prompts, function calling, contexto, customización por org, asistencia de navegación)
