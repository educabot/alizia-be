# HU-4.3: Secciones dinámicas

> Como coordinador, necesito editar y generar con IA las secciones del documento (eje problemático, estrategia, criterios) según lo que mi provincia requiera.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Las secciones del documento se definen en `config.coord_doc_sections` de la org
- [ ] El JSONB `sections` del documento almacena el contenido de cada sección
- [ ] Endpoint `PATCH /api/v1/coordination-documents/:id` permite editar secciones individuales
- [ ] Editar una section_key que no existe en la config → 422
- [ ] Secciones de tipo `text` almacenan `{ value: "..." }`
- [ ] Secciones de tipo `select_text` almacenan `{ selected_option: "...", value: "...", variants: {...} }`
- [ ] Endpoint `POST /api/v1/coordination-documents/:id/generate` genera contenido IA para todas las secciones
- [ ] Para secciones `select_text`: la IA genera **una variante por cada opción** disponible. El coordinador también puede escribir una opción a mano
- [ ] La generación usa el `ai_prompt` de cada sección en la config
- [ ] Solo tipos `text` y `select_text` en MVP
- [ ] Secciones editables en estado `in_progress` y `published` (con warning en published)
- [ ] **No hay botón "Regenerar" separado** — la re-generación de secciones individuales va por el chat con Alizia

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.3.1 | [Usecase: actualizar secciones](./tareas/T-4.3.1-usecase-actualizar-secciones.md) | src/core/usecases/ | ⬜ |
| 4.3.2 | [Usecase: generar secciones con IA](./tareas/T-4.3.2-usecase-generar-secciones.md) | src/core/usecases/ | ⬜ |
| 4.3.3 | [Endpoints PATCH y generate](./tareas/T-4.3.3-endpoints.md) | src/entrypoints/ | ⬜ |
| 4.3.4 | [Tests](./tareas/T-4.3.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.1: Modelo de datos](../HU-4.1-modelo-datos-documento/HU-4.1-modelo-datos-documento.md) — Tabla con campo sections JSONB
- [HU-3.1: Organizaciones](../../03-integracion/HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — Config con coord_doc_sections
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Azure OpenAI para generación

## Diseño técnico

### Config de secciones (organizations.config)

```jsonc
{
  "coord_doc_sections": [
    {
      "key": "problem_edge",
      "label": "Eje problemático",
      "type": "text",
      "ai_prompt": "Generá un eje problemático que integre las categorías seleccionadas...",
      "required": true
    },
    {
      "key": "methodological_strategy",
      "label": "Estrategia metodológica",
      "type": "select_text",
      "options": ["proyecto", "taller_laboratorio", "ateneo_debate"],
      "ai_prompt": "Generá una estrategia metodológica de tipo {selected_option}...",
      "required": true
    }
  ]
}
```

### Generación de variantes para `select_text` (P6 — Decisión)

Para secciones de tipo `select_text`, la IA genera **una variante por cada opción** disponible en `options`. Además, el coordinador puede **escribir una opción a mano** que no esté en la lista predefinida.

**Flujo:**
1. POST /generate dispara la generación de secciones
2. Para cada sección `text`: se genera 1 contenido
3. Para cada sección `select_text`: se generan N contenidos (uno por opción)
4. Las variantes se almacenan todas en el JSONB
5. El coordinador ve las N variantes, elige la que más le gusta, o escribe una propia
6. Al seleccionar, `selected_option` y `value` se actualizan

**Costo:** N llamadas a la IA por sección select_text (ej: 3 opciones = 3 generaciones). Aceptable para MVP dado que las secciones son pocas.

### JSONB sections del documento

```json
{
  "problem_edge": {
    "value": "¿Cómo las lógicas de poder y saber configuran..."
  },
  "methodological_strategy": {
    "selected_option": "proyecto",
    "value": "Implementaremos un proyecto interdisciplinario...",
    "variants": {
      "proyecto": "Implementaremos un proyecto interdisciplinario...",
      "taller_laboratorio": "Se desarrollará un taller de laboratorio...",
      "ateneo_debate": "Se organizará un ateneo de debate..."
    }
  }
}
```

### Opción manual del coordinador

Si el coordinador quiere escribir su propia opción en vez de elegir una variante generada:

```json
{
  "methodological_strategy": {
    "selected_option": "custom",
    "custom_option_label": "Aprendizaje basado en problemas",
    "value": "Texto escrito por el coordinador...",
    "variants": { ... }
  }
}
```

### PATCH request

```json
{
  "sections": {
    "problem_edge": {
      "value": "Nuevo contenido del eje problemático..."
    }
  }
}
```

El PATCH hace merge: solo actualiza las keys enviadas, no sobreescribe todo.

### Re-generación via chat (P7 — Decisión)

**No hay botón "Regenerar"** en la UI. Si el coordinador quiere regenerar una sección, le pide a Alizia en el chat:

> "Alizia, reescribí el eje problemático enfocándote más en sustentabilidad"

Alizia usa `update_section(section_key, content)` para actualizar la sección. Esto protege las ediciones manuales de otras secciones y da control granular al coordinador.

La generación inicial (POST /generate) sigue existiendo como primer paso después del wizard.

### Manejo de errores de IA (P17 — Decisión)

- Timeout: reducido (configurable, por defecto 45s por llamada)
- Si falla (timeout, 5xx, rate limit): **1 reintento automático**
- Si falla 2 veces: error al usuario con mensaje amigable
- Para generación multi-sección: si algunas secciones se generaron OK pero otras fallan, **guardar las exitosas** y reportar error parcial

## Test cases

- 4.10: PATCH section válida → contenido actualizado
- 4.11: PATCH section_key inválida (no existe en config) → 422
- 4.12: POST generate con sección text → contenido generado
- 4.13: POST generate con sección select_text → N variantes generadas (una por opción)
- 4.14: Elegir variante → selected_option y value actualizados
- 4.15: Escribir opción manual → custom_option_label guardado
- 4.16: Sección select_text sin selected_option al publicar → validación en HU-4.5
- 4.17: PATCH sección en documento published → 200 con warning
- 4.18: Falla IA parcial → secciones exitosas guardadas, error reportado
