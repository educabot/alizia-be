# Épica 7: Dashboard

> Vista consolidada del estado de documentos, cursos y notificaciones.

**Estado:** POST-MVP / FUTURO

---

## Problema

Coordinadores y docentes no tienen un lugar único donde ver el estado de sus documentos, planificaciones y cursos. La información está dispersa y no hay visibilidad del progreso general.

## Objetivos

- Dar visibilidad rápida del estado de documentos de coordinación y planificaciones
- Centralizar el acceso a cursos asignados
- Notificar cambios relevantes (publicaciones, actualizaciones, plazos)

## Alcance MVP

**Incluye:**

- Vista de estado de documentos de coordinación (borrador, publicado, etc.)
- Vista de cursos asignados al usuario
- Sistema de notificaciones

**No incluye:**

- Métricas de uso o analytics del docente → horizonte
- Reportes de progreso de alumnos → horizonte

## Sub-épicas

| Componente | Descripción |
|---|---|
| Estado de documentos | Visualización del estado de documentos y planificaciones |
| Cursos | Listado y acceso a cursos asignados |
| Notificaciones | Alertas sobre publicaciones, cambios y plazos |

## Decisiones de cada cliente

- Qué información se muestra en el dashboard puede variar según el rol y la provincia

## Decisiones técnicas

- Lo que ve cada usuario en el dashboard depende de su **rol y la configuración de la organización**. Un coordinador ve el estado de sus documentos y los cursos del área; un docente ve sus planificaciones y las clases próximas.
- Las notificaciones cubren eventos clave: publicación de un documento de coordinación (el docente ya puede planificar), modificaciones en clases coordinadas, y plazos próximos. El alcance exacto de notificaciones se define con el primer cliente.

## Épicas relacionadas

- **Documento de coordinación** — Los documentos se visualizan en el dashboard
- **Planificación docente** — Las planificaciones se visualizan en el dashboard
- **Roles y accesos** — El rol define qué ve cada usuario
