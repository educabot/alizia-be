# HU-4.6: Chat con Alizia

> Como coordinador, necesito chatear con Alizia para editar el documento de coordinación por lenguaje natural, sin tener que buscar y editar cada campo manualmente.

**Fase:** 3 — Coordination Documents
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `POST /api/v1/coordination-documents/:id/chat` acepta mensaje del usuario
- [ ] Alizia puede modificar el documento via function calling (tools)
- [ ] Tools disponibles (set mínimo): `update_section`, `append_to_section`, `update_class` (unificado) (P16)
- [ ] **Historial de chat persistido en backend** — tabla `coord_doc_chat_messages` (P15)
- [ ] Endpoint `GET /api/v1/coordination-documents/:id/chat/history` para retomar conversación
- [ ] **Auto-compactación** cuando la conversación es muy larga (límite configurable) (P15)
- [ ] Alizia tiene contexto del documento completo (secciones, disciplinas, plan de clases, topics)
- [ ] Las modificaciones via chat se aplican inmediatamente al documento
- [ ] Funciona en documentos `in_progress` y `published` — pero en `published`, **Alizia no ofrece tools de edición de clases** y lo comunica si el coordinador pide editar clases (P12/P16)
- [ ] **Tool calls best-effort**: si un tool falla, los demás se ejecutan igual. Errores se reportan individualmente. Se reintenta 1 vez los fallos (P18)
- [ ] **Manejo de errores IA**: 1 reintento automático con timeout reducido. Si falla 2 veces, error amigable (P17)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.6.1 | [Definición de tools (function calling)](./tareas/T-4.6.1-tools-definition.md) | src/core/usecases/ | ⬜ |
| 4.6.2 | [Usecase: chat con contexto del documento](./tareas/T-4.6.2-usecase-chat.md) | src/core/usecases/ | ⬜ |
| 4.6.3 | [Endpoint POST chat + GET history](./tareas/T-4.6.3-endpoint-chat.md) | src/entrypoints/ | ⬜ |
| 4.6.4 | [Tests](./tareas/T-4.6.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.3: Secciones dinámicas](../HU-4.3-secciones-dinamicas/HU-4.3-secciones-dinamicas.md) — update_section usa la misma lógica
- [HU-4.4: Plan de clases](../HU-4.4-plan-clases-por-materia/HU-4.4-plan-clases-por-materia.md) — update_class usa la misma lógica
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Azure OpenAI con function calling

## Diseño técnico

### Tools (function calling) — Set mínimo (P16 — Decisión)

| Tool | Descripción | Parámetros | Disponible en published |
|------|-------------|------------|------------------------|
| `update_section` | Reescribe contenido de una sección | `section_key: string, content: string` | Sí |
| `append_to_section` | Agrega texto al final de una sección | `section_key: string, content: string` | Sí |
| `update_class` | Edita título, objetivo y/o topics de una clase | `class_id: int, title?: string, objective?: string, topic_ids?: int[]` | **No** |

3 tools unificados. `update_class_title` y `update_class_topics` del RFC original se fusionan en `update_class`.

**Cuando el documento está published**, Alizia recibe en su system prompt: "El documento está publicado. Las clases son inmutables. Solo podés editar secciones. Si el coordinador pide editar clases, explicale que las clases no se pueden modificar en un documento publicado."

### Historial persistido en backend (P15 — Decisión)

**Tabla:** `coord_doc_chat_messages`

```sql
CREATE TABLE coord_doc_chat_messages (
    id SERIAL PRIMARY KEY,
    coordination_document_id INTEGER NOT NULL REFERENCES coordination_documents(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,  -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    tool_calls JSONB,           -- tool calls del assistant si aplica
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Endpoints:**
- `POST /coordination-documents/:id/chat` — envía mensaje (ya no recibe history, se lee del backend)
- `GET /coordination-documents/:id/chat/history` — retorna historial completo (paginado)

**Ventajas:**
- El coordinador puede retomar la conversación en otro dispositivo
- No se pierde contexto al cerrar el browser
- El backend tiene control total del historial

### Auto-compactación (P15 — Decisión)

Cuando la conversación es muy larga (configurable, por defecto **50 mensajes** o **~30K tokens estimados**):

1. Se toma el historial completo
2. Se envía a la IA con prompt: "Resumí esta conversación en un párrafo, destacando las acciones realizadas y el estado actual"
3. Se reemplaza el historial antiguo por un mensaje `system` con el resumen
4. Los mensajes nuevos se agregan después del resumen

Esto evita exceder el context window del LLM y reduce costos.

### Request (simplificado — sin history)

```json
{
  "message": "Cambiá el eje problemático para que se enfoque más en la sustentabilidad ambiental"
}
```

El backend carga el historial desde `coord_doc_chat_messages` automáticamente.

### Response

```json
{
  "message": "Listo, actualicé el eje problemático con foco en sustentabilidad ambiental.",
  "actions": [
    {
      "tool": "update_section",
      "result": "Sección 'problem_edge' actualizada",
      "success": true
    }
  ]
}
```

### Contexto del sistema

El system prompt incluye:
- Nombre y período del documento
- **Estado del documento** (para que Alizia sepa qué tools ofrecer)
- Secciones actuales con su contenido
- Disciplinas con sus topics y class_count
- Plan de clases resumido (con IDs para que los tools funcionen)
- Tools disponibles con sus schemas
- Instrucciones de comportamiento según estado (published → no editar clases)

### Flujo

```
Usuario escribe mensaje
  → Backend carga historial desde coord_doc_chat_messages
  → Backend arma context (documento completo + estado)
  → Envía a Azure OpenAI con tools (filtrados por estado)
  → Si respuesta tiene tool_calls → ejecuta cada uno (best-effort)
  → Persiste mensajes (user + assistant) en coord_doc_chat_messages
  → Retorna respuesta de Alizia + resultado de las acciones
```

### Tool calls best-effort + retry (P18 — Decisión)

Si Alizia ejecuta 3 tools y alguno falla:
- Tool 1 (update_section): OK → se persiste
- Tool 2 (update_class con class_id inválido): falla → se reintenta 1 vez → si sigue fallando, se reporta
- Tool 3 (update_section): OK → se persiste

La respuesta incluye qué acciones se ejecutaron y cuáles fallaron:

```json
{
  "actions": [
    {"tool": "update_section", "result": "OK", "success": true},
    {"tool": "update_class", "result": "Error: clase no encontrada", "success": false},
    {"tool": "update_section", "result": "OK", "success": true}
  ]
}
```

### Manejo de errores de IA (P17 — Decisión)

- Timeout reducido (configurable, por defecto 45s)
- 1 reintento automático para errores retriables (timeout, 5xx, 429)
- Errores no retriables (4xx): error directo
- Si falla después del reintento: 503 con mensaje amigable
- **No reintentar** tool calls que ya modificaron datos — riesgo de duplicación

### Chat como canal de re-generación (P7 — Decisión)

No hay botón "Regenerar" en la UI. Si el coordinador quiere regenerar contenido, usa el chat:

> Coordinador: "Reescribí toda la estrategia metodológica enfocándote en aprendizaje basado en problemas"
> Alizia: [ejecuta update_section("methodological_strategy", nuevo_contenido)] "Listo, reescribí la estrategia metodológica..."

Esto da al coordinador control granular (puede pedir cambios específicos) sin riesgo de sobreescribir todo el documento.

### Nota: tools genéricos por org

Las secciones del documento varían por organización (JSON Schema en `config.coord_doc_sections`). Los tools (`update_section`, etc.) son genéricos y funcionan con cualquier schema. El LLM recibe las section_keys disponibles en el system prompt.

## Test cases

- 4.26: Chat "cambiá el eje problemático" → update_section ejecutado
- 4.27: Chat "poné más horas en matemáticas" → Alizia responde que no puede
- 4.28: Chat en documento in_progress → todos los tools disponibles
- 4.29: Chat en documento published → solo tools de secciones, clases bloqueadas
- 4.30: Chat "cambiá el título de la clase 3" en published → Alizia explica que no puede
- 4.31: Tool update_section con key inválida → error manejado, Alizia informa
- 4.32: Tool call falla → retry 1 vez, si sigue fallando reportar
- 4.33: Historial persistido → cerrar browser, reabrir, historial disponible
- 4.34: Auto-compactación → conversación larga se resume automáticamente
- 4.35: GET history → historial paginado del chat
