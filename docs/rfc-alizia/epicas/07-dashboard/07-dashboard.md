# Épica 7: Dashboard

> Vista consolidada del estado de documentos, planificaciones y notificaciones por rol.

**Estado:** Post-MVP
**Fase de implementación:** Post-MVP

---

## Problema

Coordinadores y docentes no tienen un lugar único donde ver el estado de sus documentos, planificaciones y cursos. La información está dispersa y no hay visibilidad del progreso general. Un coordinador no sabe rápidamente cuántos docentes ya planificaron; un docente no sabe si hay un documento nuevo publicado.

## Objetivos

- Dar visibilidad rápida del estado de documentos de coordinación y planificaciones
- Centralizar el acceso a cursos y materias asignadas
- Notificar cambios relevantes (publicaciones, actualizaciones, plazos)
- Diferenciar la vista según el rol del usuario (coordinador vs docente)

## Alcance MVP

**Incluye:**

- Dashboard de coordinador: estado de documentos, progreso de planificación por docente
- Dashboard de docente: clases próximas, planificaciones pendientes, documentos publicados
- Sistema de notificaciones in-app para eventos clave

**No incluye:**

- Métricas de uso o analytics del docente → horizonte
- Reportes de progreso de alumnos → horizonte
- Notificaciones push / email → horizonte
- Dashboard de admin / superadmin → horizonte

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-7.1 | [Dashboard coordinador](./HU-7.1-dashboard-coordinador/HU-7.1-dashboard-coordinador.md) | Vista consolidada del coordinador: documentos, progreso, cursos | Post-MVP | 3 |
| HU-7.2 | [Dashboard docente](./HU-7.2-dashboard-docente/HU-7.2-dashboard-docente.md) | Vista del docente: clases próximas, planificaciones, documentos | Post-MVP | 4 |
| HU-7.3 | [Sistema de notificaciones](./HU-7.3-sistema-notificaciones/HU-7.3-sistema-notificaciones.md) | Notificaciones in-app para eventos relevantes | Post-MVP | 6 |

---

## Decisiones técnicas

- Lo que ve cada usuario depende de su **rol y la configuración de la organización**. Un coordinador ve el estado de sus documentos y los cursos del área; un docente ve sus planificaciones y las clases próximas.
- Las notificaciones cubren eventos clave: publicación de un documento, modificaciones en clases coordinadas, plazos próximos. El alcance exacto se define con el primer cliente.
- El dashboard NO requiere datos en tiempo real — un refresh manual o polling cada N minutos es suficiente para MVP.
- Los widgets del dashboard podrían ser configurables por org en una iteración posterior.

## Decisiones de cada cliente

- Qué información se muestra en el dashboard puede variar según el rol y la provincia
- Frecuencia y tipos de notificaciones requieren validación con usuarios reales
- El nivel de detalle del progreso de planificación (por docente, por materia, por clase) se define con cada equipo

## Épicas relacionadas

- **[Épica 1: Roles y accesos](../01-roles-accesos/01-roles-accesos.md)** — El rol define qué ve cada usuario
- **[Épica 4: Documento de coordinación](../04-documento-coordinacion/04-documento-coordinacion.md)** — Estado de documentos en el dashboard del coordinador
- **[Épica 5: Planificación docente](../05-planificacion-docente/05-planificacion-docente.md)** — Progreso de planificación y clases próximas
