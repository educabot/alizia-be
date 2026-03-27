# HU-7.3: Sistema de notificaciones

> Como usuario, necesito recibir notificaciones cuando ocurran eventos relevantes (publicación de documentos, plazos, cambios) para no perderme información importante.

**Fase:** Post-MVP
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Notificaciones in-app visibles desde un ícono en el header (campana con badge de conteo)
- [ ] Cada notificación tiene: tipo, mensaje, timestamp, link a la entidad relacionada, leída/no leída
- [ ] El usuario puede marcar notificaciones como leídas (individual y "marcar todas como leídas")
- [ ] Las notificaciones se generan automáticamente al ocurrir eventos del sistema
- [ ] Se soportan al menos los eventos definidos en la lista de eventos

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 7.3.1 | [Modelo de datos notificaciones](./tareas/T-7.3.1-modelo-datos-notificaciones.md) | entities, migracion | ⬜ |
| 7.3.2 | [Generación por eventos](./tareas/T-7.3.2-generacion-eventos.md) | usecases/notifications/ | ⬜ |
| 7.3.3 | [Job programado de plazos](./tareas/T-7.3.3-job-plazos.md) | usecases/notifications/ | ⬜ |
| 7.3.4 | [Endpoints de notificaciones](./tareas/T-7.3.4-endpoints-notificaciones.md) | handlers/ | ⬜ |
| 7.3.5 | [Frontend de notificaciones](./tareas/T-7.3.5-frontend-notificaciones.md) | frontend/ | ⬜ |
| 7.3.6 | [Tests](./tareas/T-7.3.6-tests.md) | tests/ | ⬜ |

## Dependencias

- [Épica 1: Roles y accesos](../../01-roles-accesos/01-roles-accesos.md) — Determina quién recibe cada notificación
- [Épica 4: Documento de coordinación](../../04-documento-coordinacion/04-documento-coordinacion.md) — Eventos de publicación y edición
- [Épica 5: Planificación docente](../../05-planificacion-docente/05-planificacion-docente.md) — Eventos de planificación

## Diseño de producto

### Eventos que generan notificaciones

| Evento | Destinatario | Mensaje ejemplo |
|--------|-------------|-----------------|
| Documento publicado | Docentes del área | "Se publicó el documento de coordinación Mar-Jul para Ciencias" |
| Documento vuelto a draft | Docentes del área | "El documento Mar-Jul fue revertido a borrador por el coordinador" |
| Clase compartida modificada | Docente de la otra materia | "Se modificó la planificación de Clase 3 (Física) que compartís con Matemáticas" |
| Plazo próximo | Docentes sin planificar | "Tenés 5 clases sin planificar y la primera es en 3 días" |
| Planificación completa | Coordinador del área | "Prof. García completó la planificación de Matemáticas (3a)" |
| Documento archivado | Docentes del área | "El documento Mar-Jul fue archivado" |

### Modelo conceptual

```
notifications
├── id
├── user_id (destinatario)
├── organization_id
├── type (enum: doc_published, doc_drafted, shared_class_modified, deadline_approaching, planning_complete, doc_archived)
├── title
├── message
├── entity_type (coordination_document, lesson_plan, course_subject)
├── entity_id
├── read_at (null = no leída)
├── created_at
```

### Generación de notificaciones

Dos estrategias complementarias:

1. **Event-driven:** Al ejecutar acciones (publicar documento, completar planificación), el usecase genera las notificaciones para los destinatarios correspondientes. Es síncrono dentro de la transacción o asíncrono via cola.

2. **Scheduled:** Un job periódico (ej: diario a las 8am) genera notificaciones de plazos próximos. Revisa clases sin planificar cuya fecha está dentro de los próximos N días configurables.

### Consideraciones

- Las notificaciones son **in-app solamente** para MVP. Push notifications y email son horizonte.
- El badge de conteo en el header muestra notificaciones no leídas. Se actualiza al cargar la página (no websocket en MVP).
- Las notificaciones de plazos deben ser configurables por org: cuántos días antes alertar, frecuencia de recordatorio.
- Evitar notificaciones excesivas: si un coordinador publica y despublica varias veces, no generar notificación cada vez (debounce o solo notificar en la transición a published).
- Considerar un límite de retención (ej: 90 días) para no acumular notificaciones indefinidamente.

### Preguntas abiertas

- ¿El coordinador quiere ver un resumen diario o notificación por cada docente que completa?
- ¿Las notificaciones de plazo son por materia, por curso, o globales?
- ¿Se necesita un canal de notificación adicional (WhatsApp, email) desde el inicio?
- ¿El usuario puede configurar qué notificaciones quiere recibir?
