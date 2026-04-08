# HU-4.6: Chat con Alizia

> Como coordinador, necesito chatear con Alizia para editar el documento de coordinación por lenguaje natural, sin tener que buscar y editar cada campo manualmente.

**Fase:** 3 — Coordination Documents
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `POST /api/v1/coordination-documents/:id/chat` acepta mensaje + historial
- [ ] Alizia puede modificar el documento via function calling (tools)
- [ ] Tools disponibles: `update_section`, `update_class_title`, `update_class_topics`
- [ ] El historial de chat se envía desde el frontend en cada request (no se persiste en backend)
- [ ] Alizia tiene contexto del documento completo (secciones, disciplinas, plan de clases, topics)
- [ ] Las modificaciones via chat se aplican inmediatamente al documento
- [ ] Funciona en documentos `in_progress` y `published` (documento vivo, editable post-publicación)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.6.1 | [Definición de tools (function calling)](./tareas/T-4.6.1-tools-definition.md) | src/core/usecases/ | ⬜ |
| 4.6.2 | [Usecase: chat con contexto del documento](./tareas/T-4.6.2-usecase-chat.md) | src/core/usecases/ | ⬜ |
| 4.6.3 | [Endpoint POST chat](./tareas/T-4.6.3-endpoint-chat.md) | src/entrypoints/ | ⬜ |
| 4.6.4 | [Tests](./tareas/T-4.6.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.3: Secciones dinámicas](../HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) — update_section usa la misma lógica
- [HU-4.4: Plan de clases](../HU-4.4-plan-clases-por-materia/HU-4.4-plan-clases-por-materia.md) — update_class usa la misma lógica
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Azure OpenAI con function calling

## Diseño técnico

### Tools (function calling)

| Tool | Descripción | Parámetros |
|------|-------------|------------|
| `update_section` | Actualiza contenido de una sección | `section_key: string, content: string` |
| `update_class_title` | Cambia título de una clase | `class_id: int, title: string` |
| `update_class_topics` | Actualiza topics de una clase | `class_id: int, topic_ids: int[]` |
| `append_to_section` | Agrega texto al final de una sección | `section_key: string, content: string` |

### Request

```json
{
  "message": "Cambiá el eje problemático para que se enfoque más en la sustentabilidad ambiental",
  "history": [
    {"role": "user", "content": "Generá el documento"},
    {"role": "assistant", "content": "Listo, generé las secciones..."}
  ]
}
```

### Contexto del sistema

El system prompt incluye:
- Nombre y período del documento
- Secciones actuales con su contenido
- Disciplinas con sus topics y class_count
- Plan de clases resumido
- Tools disponibles con sus schemas

### Flujo

```
Usuario escribe mensaje
  → Backend arma context (documento completo)
  → Envía a Azure OpenAI con tools
  → Si respuesta tiene tool_calls → ejecuta cada uno
  → Retorna respuesta de Alizia + resultado de las acciones
```

### Nota: tools genéricos por org

Las secciones del documento varían por organización (JSON Schema en `config.coord_doc_sections`). Los tools listados arriba (`update_section`, etc.) son ejemplos concretos — la implementación real debe ser genérica para funcionar con cualquier schema. Ver nota detallada en HU-6.3 sobre JSON Path como solución.

## Test cases

- 4.26: Chat "cambiá el eje problemático" → update_section ejecutado
- 4.27: Chat "poné más horas en matemáticas" → Alizia responde que no puede (no hay tool para eso)
- 4.28: Chat en documento published → permitido (documento vivo, editable post-publicación)
- 4.29: Tool update_section con key inválida → error manejado, Alizia informa
- 4.30: Historial se envía y se mantiene coherente
