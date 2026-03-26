# HU-4.4: Plan de clases por materia

> Como coordinador, necesito un plan de clases generado por IA para cada materia del documento, con título, objetivo y topics por clase.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] `coord_doc_classes` se generan por IA al llamar `POST /generate` (junto con secciones)
- [ ] Cada clase tiene: class_number, title, objective
- [ ] Cada clase tiene topics asignados (coord_doc_class_topics)
- [ ] Se generan tantas clases como `class_count` de la materia
- [ ] Los topics asignados a la materia (coord_doc_subject_topics) se distribuyen entre las clases
- [ ] El coordinador puede editar título, objetivo y topics de cada clase individualmente
- [ ] Las clases compartidas (de HU-3.5) se marcan en la respuesta del detalle
- [ ] Endpoint para editar una clase individual: `PATCH /api/v1/coord-doc-classes/:id`

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.4.1 | [Usecase: generar plan de clases](./tareas/T-4.4.1-usecase-generar-plan.md) | src/core/usecases/ | ⬜ |
| 4.4.2 | [Prompt y schema para generación](./tareas/T-4.4.2-prompt-schema.md) | src/core/usecases/ | ⬜ |
| 4.4.3 | [Endpoint PATCH clase individual](./tareas/T-4.4.3-endpoint-editar-clase.md) | src/entrypoints/ | ⬜ |
| 4.4.4 | [Integrar shared classes en respuesta](./tareas/T-4.4.4-shared-classes.md) | src/core/usecases/ | ⬜ |
| 4.4.5 | [Tests](./tareas/T-4.4.5-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.2: Wizard](../HU-4.2-wizard-creacion/HU-4.2-wizard-creacion.md) — Documento con subjects y topics creados
- [HU-3.5: Grilla horaria](../../03-integracion/HU-3.5-grilla-horaria-clases-compartidas/HU-3.5-grilla-horaria-clases-compartidas.md) — Para marcar clases compartidas
- [Épica 6: Asistente IA](../../06-assistente-ia/06-asistente-ia.md) — Azure OpenAI para generar el plan

## Diseño técnico

### Generación IA del plan de clases

Para cada materia del documento, se envía al LLM:

**Input:**
- Materia (nombre)
- Topics asignados a la materia
- class_count
- Secciones del documento ya generadas (eje problemático, estrategia)

**Output esperado:**
```json
[
  {
    "class_number": 1,
    "title": "Introducción al pensamiento algebraico",
    "objective": "Que los estudiantes identifiquen patrones numéricos...",
    "topic_ids": [5]
  },
  {
    "class_number": 2,
    "title": "Ecuaciones de primer grado",
    "objective": "Que los estudiantes resuelvan ecuaciones simples...",
    "topic_ids": [5, 8]
  }
]
```

### Clases compartidas en la respuesta

El GET detalle del documento incluye para cada clase:

```json
{
  "class_number": 3,
  "title": "Interdisciplina: Matemáticas y Física",
  "is_shared": true,
  "shared_with_subject": "Física"
}
```

`is_shared` se calcula cruzando con los shared class numbers de HU-3.5.

## Test cases

- 4.14: POST generate → plan de clases creado para cada materia
- 4.15: Cantidad de clases generadas == class_count
- 4.16: Todos los topics de la materia distribuidos en al menos una clase
- 4.17: PATCH clase → título y topics actualizados
- 4.18: Clases compartidas marcadas correctamente
