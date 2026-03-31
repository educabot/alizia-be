# HU-8.6: Creación libre y vínculo recurso–planificación

> Como docente, necesito crear recursos con lineamientos propios (no solo tipos predefinidos) y vincularlos a clases de mi planificación para usarlos antes, durante o después de la clase.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `POST /api/v1/resources` acepta `type: "free"` con campo `guidelines` (texto libre del usuario que describe qué quiere generar)
- [ ] La IA genera el recurso usando los guidelines del usuario como prompt principal + contexto de la disciplina/curso
- [ ] Tabla `resource_lesson_plan_links` vincula un recurso con una clase de un lesson plan (resource_id, lesson_plan_id, class_number, moment: before|during|after)
- [ ] Endpoint `POST /api/v1/resources/:id/link` asocia recurso a clase
- [ ] Endpoint `GET /api/v1/lesson-plans/:id/resources` retorna recursos vinculados por clase
- [ ] Endpoint `DELETE /api/v1/resources/:id/link/:linkId` desvincula
- [ ] Un recurso puede estar vinculado a múltiples clases
- [ ] Una clase puede tener múltiples recursos

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.3.1 | Migración: tabla resource_lesson_plan_links | migrations/ | ⬜ |
| 8.3.2 | Usecase: creación libre con guidelines | src/core/usecases/ | ⬜ |
| 8.3.3 | Endpoints: link/unlink recurso-plan | src/entrypoints/ | ⬜ |
| 8.3.4 | Endpoint: recursos por lesson plan | src/entrypoints/ | ⬜ |
| 8.3.5 | Tests | tests/ | ⬜ |

## Dependencias

- [HU-8.1: Modelo de datos recursos](../HU-8.1-modelo-datos-recursos/HU-8.1-modelo-datos-recursos.md) — tabla resources, generación por tipo
- [HU-5.3: Lesson plans con class_plan](../../05-planificacion-docente/HU-5.3-lesson-plans/HU-5.3-lesson-plans.md) — lesson plans con class_plan

## Diseño técnico

### Tabla resource_lesson_plan_links

```sql
CREATE TABLE resource_lesson_plan_links (
  id SERIAL PRIMARY KEY,
  resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
  lesson_plan_id INT NOT NULL REFERENCES teacher_lesson_plans(id) ON DELETE CASCADE,
  class_number INT NOT NULL,
  moment VARCHAR(10) NOT NULL CHECK (moment IN ('before', 'during', 'after')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(resource_id, lesson_plan_id, class_number, moment)
);
```

## Test cases

- Crear recurso type "free" con guidelines → IA genera usando guidelines como prompt
- Vincular recurso a clase 3 momento "before" → link creado
- GET recursos de lesson plan → retorna agrupados por clase
- Vincular mismo recurso a 2 clases → OK
- Desvincular → link eliminado, recurso permanece
