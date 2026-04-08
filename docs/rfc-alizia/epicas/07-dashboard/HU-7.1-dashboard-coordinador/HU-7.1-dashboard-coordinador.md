# HU-7.1: Dashboard coordinador

> Como coordinador, necesito una vista consolidada donde ver el estado de mis documentos de coordinación, el progreso de planificación de los docentes y los cursos de mi área.

**Fase:** Fase 7
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] El coordinador ve al ingresar un resumen de sus documentos de coordinación con estado (pendiente, en progreso, publicado)
- [ ] Se muestra el progreso de planificación: cuántos docentes ya planificaron sus clases por documento publicado
- [ ] Se listan los cursos del área con acceso rápido al detalle
- [ ] Se muestra un indicador de documentos que requieren acción (sin publicar, con ediciones pendientes)
- [ ] Los datos se cargan al entrar y se refrescan con acción del usuario (no real-time)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 7.1.1 | [Endpoints de agregación](./tareas/T-7.1.1-endpoints-agregacion.md) | handlers/, usecases/dashboard/ | ⬜ |
| 7.1.2 | [Frontend dashboard coordinador](./tareas/T-7.1.2-frontend-coordinador.md) | frontend/ | ⬜ |
| 7.1.3 | [Tests](./tareas/T-7.1.3-tests.md) | tests/ | ⬜ |
| 7.1.4 | Heuristicas configurables para widget "Requiere atencion": query que detecta documentos sin publicar hace >N dias (default 7), docentes sin planificar a <N dias del inicio (default 14), documentos completos sin publicar. Thresholds configurables en organizations.config | usecases/dashboard/ | ⬜ |

## Dependencias

- [HU-1.2: Modelo de usuarios y roles](../../01-roles-accesos/HU-1.2-modelo-usuarios-roles/HU-1.2-modelo-usuarios-roles.md) — Rol coordinador
- [HU-4.5: Publicación y estados](../../04-documento-coordinacion/HU-4.5-publicacion-estados/HU-4.5-publicacion-estados.md) — Estado de documentos
- [HU-5.5: Estados de planificación](../../05-planificacion-docente/HU-5.5-estados-planificacion/HU-5.5-estados-planificacion.md) — Progreso de planificación

## Diseño de producto

### Widgets del dashboard

| Widget | Datos | Acción |
|--------|-------|--------|
| **Mis documentos** | Lista de documentos con badge de estado, fecha de creación | Click → ir al documento |
| **Progreso de planificación** | Por documento publicado: N/M docentes planificaron, barra de progreso | Click → ver detalle por disciplina |
| **Cursos del área** | Lista de cursos con cantidad de disciplinas y docentes asignados | Click → ir al curso |
| **Requiere atención** | Documentos en borrador hace más de N días, docentes sin planificar | Click → ir al item |

### Ejemplo visual

```
┌─────────────────────────────────────────────────┐
│  Hola, [Coordinador]         Área: Ciencias     │
├──────────────────────┬──────────────────────────┤
│ Mis documentos       │ Progreso planificación   │
│                      │                          │
│ ● Doc Mar-Jul  🟢   │ Doc Mar-Jul (publicado)  │
│ ● Doc Ago-Nov  🟡   │ ████████░░ 4/5 docentes  │
│                      │                          │
├──────────────────────┼──────────────────────────┤
│ Cursos del área      │ Requiere atención        │
│                      │                          │
│ 3a - 4 disciplinas      │ ⚠ Doc Ago-Nov en borrador│
│ 5b - 3 disciplinas      │   hace 15 días           │
│ 6a - 4 disciplinas      │ ⚠ Doc. García sin       │
│                      │   planificar (Matemáticas)│
└──────────────────────┴──────────────────────────┘
```

### Consideraciones

- El progreso de planificación ya existe como endpoint en HU-5.5 (`GET /coordination-documents/:id/planning-progress`) — el dashboard lo consume
- "Requiere atención" es una heurística configurable (ej: borrador > 7 días, docente sin planificar a 2 semanas del inicio)
- El dashboard debe cargar rápido — las queries de agregación deben ser eficientes (considerar vistas disciplinalizadas si hay volumen)
