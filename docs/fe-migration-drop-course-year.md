# FE Migration — Drop `courses.year`

**Branch:** `feature/sl/epica-3-integracion`
**Fecha:** 2026-04-16
**Gap:** G-5 (alineación contrato `POST /courses` con RFC)
**Migración DB:** `000013_drop_courses_year.up.sql`

---

## TL;DR

El campo `year` desaparece de **courses**. El año académico ahora vive **únicamente** en `course_subjects.school_year`. Un mismo `course` (ej. "2do 1era") puede reutilizarse a través de múltiples años académicos sin duplicarse, con un `course_subject` por año.

| Recurso | Antes | Ahora |
|---|---|---|
| `Course` (entity) | tenía `year: int` | sin `year` |
| `course_subjects[].school_year` | ya existía | **única fuente de verdad** del año |

---

## Cambios por endpoint

### 1. `POST /api/v1/courses` — crear curso

**Request body — antes:**
```json
{
  "name": "2do 1era",
  "year": 2026
}
```

**Request body — ahora:**
```json
{
  "name": "2do 1era"
}
```

> Si el FE sigue mandando `year`, el BE lo ignora silenciosamente (no hay validación que rechace campos extra). No es error pero tampoco hace nada — sacarlo del payload.

**Response — antes:**
```json
{
  "id": 1,
  "name": "2do 1era",
  "year": 2026,
  "students": [],
  "course_subjects": []
}
```

**Response — ahora:**
```json
{
  "id": 1,
  "name": "2do 1era",
  "students": [],
  "course_subjects": []
}
```

---

### 2. `GET /api/v1/courses` — listar cursos

Cada item del array deja de incluir `year`.

```json
{
  "items": [
    { "id": 1, "name": "2do 1era", "students": [], "course_subjects": [] },
    { "id": 2, "name": "3ro 2da", "students": [], "course_subjects": [] }
  ],
  "has_more": false
}
```

---

### 3. `GET /api/v1/courses/:id` — detalle de curso

Idem: el campo `year` ya no aparece en el objeto raíz. El año académico se lee desde `course_subjects[].school_year`.

```json
{
  "id": 1,
  "name": "2do 1era",
  "students": [...],
  "course_subjects": [
    {
      "id": 10,
      "course_id": 1,
      "subject_id": 5,
      "teacher_id": 12,
      "school_year": 2026,
      "start_date": "2026-03-01",
      "end_date": "2026-12-15",
      "subject": { "id": 5, "name": "Matemática" },
      "teacher": { "id": 12, "first_name": "Ana", "last_name": "Pérez" }
    }
  ]
}
```

---

## Acciones requeridas en el FE

### TypeScript / interfaces
- Sacar `year: number` del type `Course` (o equivalente).
- Sacar `year` del payload del form de creación de curso.
- Si hay un store/cache de cursos, limpiarlo o resetearlo en el primer load post-deploy para evitar leer la propiedad inexistente.

### UI
- Si la pantalla "Listado de cursos" mostraba el año al lado del nombre del curso, derivarlo desde `course_subjects[0].school_year` (o desde el filtro de año activo en el contexto, si existe).
- Si el form de "Crear curso" pedía el año al usuario:
  - **Opción A**: removerlo del form (curso es atemporal; el año lo define la asignación de materia/profesor).
  - **Opción B (recomendada UX)**: mantener el input de año en el form, pero usarlo en el siguiente paso (asignar materias) como `school_year` del primer `course_subject`. NO mandarlo en el `POST /courses`.

### Filtro "cursos del año X"
- **Antes:** `GET /courses` y filtrar por `course.year`.
- **Ahora:** `GET /course-subjects?school_year=X` (cuando se agregue el filtro — ver nota abajo) y agrupar por `course_id`. O bien `GET /courses` + para cada uno mirar sus `course_subjects[].school_year`.

> **Nota:** el endpoint `GET /course-subjects` ya soporta filtros opcionales (`course_id`, `subject_id`, `teacher_id`). Si el FE necesita filtro por `school_year`, abrir un ticket — es trivial agregarlo en BE.

---

## Compatibilidad

- **No hay versionado de API**: el cambio es breaking inmediato.
- **No hay rollback automático**: la migración `000013_drop_courses_year.down.sql` recupera la columna pero los valores originales se pierden (default = año actual).
- **Coordinar deploy**: mergear FE y BE en la misma ventana, o el FE va a recibir `undefined` en `course.year` y romper si no defendió el acceso.

---

## Checklist FE

- [ ] Type `Course` sin `year`
- [ ] Payload `POST /courses` sin `year`
- [ ] UI listado de cursos: derivar año desde `course_subjects` o quitar la columna
- [ ] UI form crear curso: ajustado según opción A o B
- [ ] Tests/mocks de FE actualizados
- [ ] Coordinar release con BE (mismo deploy)
