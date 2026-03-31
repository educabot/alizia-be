# HU-7.2: Dashboard docente

> Como docente, necesito una vista donde ver mis clases próximas, el estado de mis planificaciones y los documentos de coordinación publicados para mis disciplinas.

**Fase:** Fase 7
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] El docente ve al ingresar un resumen de sus próximas clases (basado en cronograma)
- [ ] Se muestra el estado de planificación por disciplina: clases planificadas vs pendientes
- [ ] Se listan los documentos de coordinación publicados que aplican a sus disciplinas
- [ ] Se indica si hay un documento nuevo que aún no revisó
- [ ] Las clases compartidas se marcan visualmente

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 7.2.1 | [Endpoints de agregación docente](./tareas/T-7.2.1-endpoints-docente.md) | handlers/, usecases/dashboard/ | ⬜ |
| 7.2.2 | [Tracking de vistas de documentos](./tareas/T-7.2.2-tracking-vistas.md) | entities, migracion | ⬜ |
| 7.2.3 | [Frontend dashboard docente](./tareas/T-7.2.3-frontend-docente.md) | frontend/ | ⬜ |
| 7.2.4 | [Tests](./tareas/T-7.2.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-5.2: Cronograma docente](../../05-planificacion-docente/HU-5.2-cronograma-docente/HU-5.2-cronograma-docente.md) — Cronograma de clases
- [HU-5.5: Estados de planificación](../../05-planificacion-docente/HU-5.5-estados-planificacion/HU-5.5-estados-planificacion.md) — Estado de planificaciones
- [HU-4.5: Publicación y estados](../../04-documento-coordinacion/HU-4.5-publicacion-estados/HU-4.5-publicacion-estados.md) — Documentos publicados

## Diseño de producto

### Widgets del dashboard

| Widget | Datos | Acción |
|--------|-------|--------|
| **Próximas clases** | Próximas 3-5 clases con fecha, disciplina, título y estado de planificación | Click → ir a planificar o ver planificación |
| **Mis planificaciones** | Por disciplina: barra de progreso (planificadas/total), estado general | Click → ir al cronograma |
| **Documentos publicados** | Documentos de coordinación de mi área, con badge "nuevo" si no lo abrí | Click → ver documento |
| **Clases compartidas** | Próximas clases compartidas con nombre del otro docente y disciplina | Click → ver planificación |

### Ejemplo visual

```
┌─────────────────────────────────────────────────┐
│  Hola, [Docente]              Mis disciplinas: 3   │
├──────────────────────┬──────────────────────────┤
│ Próximas clases      │ Mis planificaciones      │
│                      │                          │
│ Lun 31 - Matemáticas │ Matemáticas (3a)         │
│ "Ecuaciones" ✅      │ ██████████░░ 8/10 clases │
│                      │                          │
│ Mar 1 - Física       │ Física (3a)              │
│ "Cinemática" ⬜      │ ████░░░░░░░░ 3/10 clases │
│                      │                          │
│ Mié 2 - Matemáticas  │ Matemáticas (5b)         │
│ "Funciones" 🤝 ⬜    │ ░░░░░░░░░░░░ 0/8 clases  │
├──────────────────────┼──────────────────────────┤
│ Documentos           │ Clases compartidas       │
│                      │                          │
│ 🆕 Doc Mar-Jul (3a)  │ Mié 2 - Matemáticas 🤝  │
│    Publicado 28/03   │ con Doc. López (Física) │
│ Doc Mar-Jul (5b)     │                          │
│    Publicado 15/03   │ Mié 9 - Matemáticas 🤝  │
│                      │ con Doc. López (Física) │
└──────────────────────┴──────────────────────────┘
```

### Consideraciones

- "Próximas clases" se calcula a partir del cronograma heredado del documento de coordinación y la fecha actual
- El badge "nuevo" requiere trackear la última vez que el docente abrió cada documento (tabla `user_document_views` o similar)
- Las clases compartidas merecen visibilidad especial porque requieren coordinación entre docentes
- El docente puede tener múltiples disciplinas en múltiples cursos — el dashboard debe agrupar de forma clara
- Considerar un estado "urgente" para clases próximas sin planificar (ej: dentro de los próximos 3 días hábiles)
