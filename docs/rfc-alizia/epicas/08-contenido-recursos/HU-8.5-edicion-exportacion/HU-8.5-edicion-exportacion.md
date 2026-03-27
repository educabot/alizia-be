# HU-8.5: Edición y exportación

> Como docente, necesito editar un recurso generado (manualmente o con ayuda de Alizia) y exportarlo para impresión, para tenerlo listo para el aula.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] El docente puede editar manualmente cualquier sección del contenido generado
- [ ] El docente puede chatear con Alizia para pedir modificaciones asistidas al recurso
- [ ] Los cambios se guardan en `resources.content`
- [ ] El recurso puede pasar de `generated` a `published` cuando el docente lo da por listo
- [ ] El recurso publicado se puede exportar a formato imprimible (vista de impresión HTML)
- [ ] El recurso se puede vincular a un lesson plan desde la planificación

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.5.1 | [Endpoints de edición](./tareas/T-8.5.1-endpoints-edicion.md) | handlers/ | ⬜ |
| 8.5.2 | [Chat con Alizia para recursos](./tareas/T-8.5.2-chat-alizia-recursos.md) | usecases/ | ⬜ |
| 8.5.3 | [Frontend de edición y exportación](./tareas/T-8.5.3-frontend-edicion.md) | frontend/ | ⬜ |
| 8.5.4 | [Tests](./tareas/T-8.5.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-8.3: Generación con IA](../HU-8.3-generacion-ia-recursos/HU-8.3-generacion-ia-recursos.md) — Recurso generado
- [HU-6.3: Motor de function calling](../../06-assistente-ia/HU-6.3-motor-function-calling/HU-6.3-motor-function-calling.md) — ChatWithTools para edición asistida
- [HU-5.3: Planificación por momentos](../../05-planificacion-docente/HU-5.3-planificacion-por-momentos/HU-5.3-planificacion-por-momentos.md) — Vínculo de recurso a lesson plan

## Test cases

- 8.16: Editar sección del content → content actualizado
- 8.17: Chat con Alizia → modifica content via function calling
- 8.18: Publicar recurso → status cambia a published
- 8.19: Exportar recurso → vista imprimible generada
