# HU-6.5: Customizacion por organizacion

> Como administrador de la plataforma, necesito personalizar el comportamiento de Alizia por organizacion, para que cada provincia tenga un asistente alineado con su estilo pedagogico y requerimientos curriculares.

**Fase:** Post-MVP
**Prioridad:** Media
**Estimacion:** --

---

## Contexto

Cada jurisdiccion educativa tiene su propia terminologia, documentos curriculares de referencia, y expectativas sobre como debe comunicarse un asistente de planificacion. Una configuracion unica no sirve para todas las provincias. Esta HU define las dimensiones de customizacion que la plataforma debe soportar para que Alizia se adapte a cada organizacion.

## Dimensiones de customizacion

### Tono

Control sobre el estilo comunicacional de Alizia:
- Formal vs informal
- Variaciones regionales (rioplatense, neutro, etc.)
- Nivel de tecnicismo pedagogico (basico vs avanzado)

### Limites

Restricciones sobre la generacion de contenido:
- Longitud maxima del texto generado por seccion
- Que secciones pueden ser generadas por IA vs completadas manualmente
- Cantidad maxima de interacciones de chat por sesion (si aplica)

### Fuentes curriculares

Documentos de referencia que Alizia puede citar y usar como base:
- Varian segun jurisdiccion (NAP nacionales, DCP provinciales, etc.)
- Deben poder actualizarse sin deploy (configuracion, no codigo)

### Prompts custom

Sobreescritura de prompts por tipo de seccion:
- Permite que el equipo pedagogico de cada provincia ajuste las instrucciones
- Aplica a: estrategia metodologica, plan de clases, respuestas de chat
- Se inyectan a traves del SystemPromptBuilder existente (HU-6.2)

### Vocabulario

Mapeo de terminologia provincial:
- Reemplazo de terminos en el texto generado (post-procesamiento)
- Ejemplos: "contenidos" -> "saberes", "planificacion" -> "secuencia didactica", "NAP" -> "contenidos prioritarios"

### Feature flags

Control granular sobre que funcionalidades de IA estan habilitadas:
- Generacion de secciones
- Chat con Alizia
- Recomendaciones proactivas
- Si una feature esta deshabilitada, el endpoint retorna 403

## Modelo de configuracion

La configuracion vive en `organizations.config.ai_settings`:

```json
{
  "ai_settings": {
    "tone": "rioplatense_informal",
    "max_generation_length": 2000,
    "enabled_features": ["generation", "chat", "recommendations"],
    "curriculum_sources": ["nap_2024", "dcp_mendoza_2023"],
    "vocabulary_overrides": {
      "contenidos": "saberes",
      "planificación": "secuencia didáctica"
    },
    "prompt_overrides": {
      "strategy": "template custom...",
      "class_plan": "template custom..."
    }
  }
}
```

## Criterios de aceptacion

- [ ] Cada organizacion tiene un campo `ai_settings` en su configuracion
- [ ] El tono configurado se inyecta en todos los system prompts via SystemPromptBuilder
- [ ] Los feature flags controlan la disponibilidad de endpoints de IA (403 si deshabilitado)
- [ ] Las sobreescrituras de vocabulario se aplican como post-procesamiento sobre el texto generado
- [ ] Los prompts custom reemplazan los prompts por defecto cuando estan definidos
- [ ] Las fuentes curriculares configuradas se incluyen como contexto en los prompts
- [ ] La longitud maxima se respeta en todas las generaciones (truncar o instruir al modelo)
- [ ] Existe una configuracion por defecto razonable para organizaciones sin customizacion

## Consideraciones

- El tono afecta TODOS los system prompts — necesita inyectarse en el SystemPromptBuilder (HU-6.2)
- Las sobreescrituras de vocabulario son post-procesamiento simple (find & replace), no regeneracion
- Requiere colaboracion cercana con cada equipo pedagogico provincial
- Arrancar con 2-3 configuraciones e iterar basandose en feedback real
- Validar que los prompts custom no rompan el formato esperado de las respuestas
- Considerar un modo "preview" para que el equipo pedagogico pruebe cambios antes de publicar

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.5.1 | [Schema de ai_settings](./tareas/T-6.5.1-ai-settings-schema.md) | entities/ai_settings.go, migracion | ⬜ |
| 6.5.2 | [Post-procesador de vocabulario](./tareas/T-6.5.2-vocabulary-postprocessor.md) | usecases/ai/vocabulary_processor.go | ⬜ |
| 6.5.3 | [Middleware de feature flags](./tareas/T-6.5.3-feature-flags-middleware.md) | middleware/ai_feature_gate.go | ⬜ |
| 6.5.4 | [Inyeccion de tono y prompts custom](./tareas/T-6.5.4-tone-prompt-injection.md) | usecases/ai/tone_instructions.go | ⬜ |
| 6.5.5 | [Tests](./tareas/T-6.5.5-tests.md) | tests | ⬜ |

## Dependencias

- HU-6.1 completada (AI Provider)
- HU-6.2 completada (Prompts configurables — SystemPromptBuilder)
- Modelo de organizaciones implementado con campo `config` JSONB

## Preguntas abiertas

- Quien tiene permisos para modificar la configuracion de IA? (admin global vs admin de org)
- Se necesita versionado de configuraciones para rollback?
- Como se testean los cambios de prompts antes de aplicarlos en produccion?
- Que metricas se recolectan para evaluar la efectividad de cada configuracion?
