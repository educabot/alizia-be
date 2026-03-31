# HU-3.5: Grilla horaria y clases compartidas

> Como admin, necesito definir la grilla horaria semanal de cada curso y marcar qué clases son compartidas entre dos disciplinas.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `time_slots` con: id, course_id, day_of_week, start_time, end_time, created_at
- [ ] Tabla `time_slot_subjects` con: id, time_slot_id, course_subject_id, UNIQUE(time_slot_id, course_subject_id)
- [ ] Endpoint `POST /api/v1/courses/:id/time-slots` crea slot con sus subjects
- [ ] Clase normal: 1 time_slot → 1 time_slot_subject
- [ ] Clase compartida: 1 time_slot → 2 time_slot_subjects
- [ ] Si `config.shared_classes_enabled = false` → máximo 1 subject por slot
- [ ] Si `config.shared_classes_enabled = true` → máximo 2 subjects por slot
- [ ] Trigger valida que ambos course_subjects pertenezcan al mismo course_id que el time_slot
- [ ] Endpoint de detalle de curso incluye la grilla horaria con subjects resueltos
- [ ] Query para detectar clases compartidas de un curso funciona correctamente
- [ ] Query para calcular qué class_numbers son compartidos para un course_subject

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 3.5.1 | [Migración: time_slots, time_slot_subjects](./tareas/T-3.5.1-migracion.md) | db/migrations/ | ⬜ |
| 3.5.2 | [Trigger: validate_time_slot_subject](./tareas/T-3.5.2-trigger-validacion.md) | db/migrations/ | ⬜ |
| 3.5.3 | [Entities y providers](./tareas/T-3.5.3-entities-providers.md) | src/core/ | ⬜ |
| 3.5.4 | [Repository + queries de shared classes](./tareas/T-3.5.4-repository-shared-classes.md) | src/repositories/ | ⬜ |
| 3.5.5 | [Endpoints y handlers](./tareas/T-3.5.5-endpoints-handlers.md) | src/entrypoints/ | ⬜ |
| 3.5.6 | [Seed y tests](./tareas/T-3.5.6-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-3.4: Cursos y asignaciones](../HU-3.4-cursos-alumnos-asignaciones/HU-3.4-cursos-alumnos-asignaciones.md) — courses y course_subjects deben existir
- [HU-3.1: Organizaciones](../HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — config con `shared_classes_enabled`

## Diseño técnico

### Modelo normalizado (reemplaza JSONB schedule)

```
time_slots           time_slot_subjects
┌─────────────┐     ┌──────────────────┐
│ id          │     │ id               │
│ course_id   │──┐  │ time_slot_id     │──→ time_slots
│ day_of_week │  │  │ course_subject_id│──→ course_subjects
│ start_time  │  │  └──────────────────┘
│ end_time    │  │
└─────────────┘  │  Clase normal: 1 registro
                 │  Clase compartida: 2 registros
```

### Trigger de validación

```sql
CREATE OR REPLACE FUNCTION validate_time_slot_subject() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM course_subjects cs
        JOIN time_slots ts ON ts.course_id = cs.course_id
        WHERE cs.id = NEW.course_subject_id AND ts.id = NEW.time_slot_id
    ) THEN
        RAISE EXCEPTION 'course_subject does not belong to the same course as the time_slot';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### Query: detectar clases compartidas

```sql
SELECT ts.day_of_week, ts.start_time, ts.end_time,
       array_agg(cs.id) AS course_subject_ids
FROM time_slots ts
JOIN time_slot_subjects tss ON tss.time_slot_id = ts.id
JOIN course_subjects cs ON cs.id = tss.course_subject_id
WHERE ts.course_id = $1
GROUP BY ts.id
HAVING count(*) > 1;
```

### Query: calcular class_numbers compartidos

```sql
WITH subject_slots AS (
    SELECT ts.day_of_week, ts.start_time, ts.id AS slot_id,
           ROW_NUMBER() OVER (ORDER BY ts.day_of_week, ts.start_time) AS weekly_position
    FROM time_slots ts
    JOIN time_slot_subjects tss ON tss.time_slot_id = ts.id
    WHERE tss.course_subject_id = $1
),
shared_positions AS (
    SELECT ss.weekly_position
    FROM subject_slots ss
    JOIN time_slot_subjects tss ON tss.time_slot_id = ss.slot_id
    GROUP BY ss.weekly_position, ss.slot_id
    HAVING count(*) > 1
)
SELECT weekly_position + (week * classes_per_week) AS class_number
FROM shared_positions, generate_series(0, total_weeks - 1) AS week;
```

## Test cases

- 3.20: POST time-slot con 1 subject → clase normal
- 3.21: POST time-slot con 2 subjects → clase compartida
- 3.22: POST time-slot con 2 subjects y `shared_classes_enabled = false` → 422
- 3.23: POST time-slot con course_subject de otro curso → error (trigger)
- 3.24: GET course/:id → grilla con slots y subjects resueltos
- 3.25: Query shared classes → retorna solo slots con 2+ subjects
