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
- [ ] Secciones de tipo `select_text` almacenan `{ selected_option: "...", value: "..." }`
- [ ] Endpoint `POST /api/v1/coordination-documents/:id/generate` genera contenido IA para todas las secciones
- [ ] La generación usa el `ai_prompt` de cada sección en la config
- [ ] El placeholder `{selected_option}` se reemplaza con la opción elegida

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

### JSONB sections del documento

```json
{
  "problem_edge": {
    "value": "¿Cómo las lógicas de poder y saber configuran..."
  },
  "methodological_strategy": {
    "selected_option": "proyecto",
    "value": "Implementaremos un proyecto interdisciplinario..."
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

## Test cases

- 4.10: PATCH section válida → contenido actualizado
- 4.11: PATCH section_key inválida (no existe en config) → 422
- 4.12: POST generate → todas las secciones generadas con IA
- 4.13: Sección select_text sin selected_option → 422
