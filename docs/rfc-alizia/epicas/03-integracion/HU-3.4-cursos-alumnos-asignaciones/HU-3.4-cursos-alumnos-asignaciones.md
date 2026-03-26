# HU-3.4: Cursos, alumnos y asignaciones

> Como admin, necesito crear cursos, cargar alumnos y asignar docentes a materias para estructurar la vida académica de la institución.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `courses` con: id, organization_id, name, created_at
- [ ] Tabla `students` con: id, course_id, name, created_at
- [ ] Tabla `course_subjects` con: id, course_id, subject_id, teacher_id, start_date, end_date, school_year, created_at
- [ ] Endpoints para courses: POST, GET (listar), GET /:id (detalle con students)
- [ ] Endpoints para students: POST (bulk o individual dentro de un curso)
- [ ] Endpoint para course_subjects: POST (asignar docente a materia en curso)
- [ ] GET /courses/:id incluye students y course_subjects con teacher y subject preloaded
- [ ] `course_subjects` vincula curso + materia + docente + período lectivo
- [ ] Un docente puede estar asignado a múltiples course_subjects
- [ ] Todos los queries filtran por organization_id del JWT

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 3.4.1 | [Migración: courses, students, course_subjects](./tareas/T-3.4.1-migracion.md) | db/migrations/ | ⬜ |
| 3.4.2 | [Entities](./tareas/T-3.4.2-entities.md) | src/core/entities/ | ⬜ |
| 3.4.3 | [Providers y repository](./tareas/T-3.4.3-providers-repository.md) | src/core/providers/, src/repositories/ | ⬜ |
| 3.4.4 | [Endpoints y handlers](./tareas/T-3.4.4-endpoints-handlers.md) | src/entrypoints/ | ⬜ |
| 3.4.5 | [Seed y tests](./tareas/T-3.4.5-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — FK organization_id
- [HU-3.2: Áreas y materias](../HU-3.2-areas-materias/HU-3.2-areas-materias.md) — subjects debe existir para course_subjects
- [HU-1.2: Modelo de usuarios](../../01-roles-accesos/HU-1.2-modelo-usuarios-roles/HU-1.2-modelo-usuarios-roles.md) — users debe existir para teacher_id

## Diseño técnico

### course_subjects

```sql
CREATE TABLE course_subjects (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id),
    subject_id INTEGER NOT NULL REFERENCES subjects(id),
    teacher_id INTEGER NOT NULL REFERENCES users(id),
    start_date DATE,
    end_date DATE,
    school_year INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Es la **instancia** que dice: "en el curso 3A, la materia Matemáticas la da el docente Juan durante 2026". Es la tabla central que conecta el mundo admin con el mundo docente.

## Test cases

- 3.15: POST course → 201
- 3.16: POST student en curso → alumno asociado
- 3.17: POST course_subject → asignación docente-materia-curso
- 3.18: GET course/:id → incluye students, course_subjects con teacher y subject
- 3.19: Asignar docente de otra org → error (FK + multi-tenancy)
