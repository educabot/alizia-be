# HU-4.4: Plan de clases por disciplina

> Como coordinador, necesito un plan de clases generado por IA para cada disciplina del documento, con título, objetivo y topics por clase.

**Fase:** 3 — Coordination Documents
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] `coord_doc_classes` se generan por IA al llamar `POST /generate` (junto con secciones)
- [ ] Cada clase tiene: class_number, title, objective
- [ ] Cada clase tiene topics asignados (coord_doc_class_topics)
- [ ] Se generan tantas clases como `class_count` de la disciplina
- [ ] Si la IA genera un count incorrecto: **reintentar 1 vez** con prompt reforzado. Si falla 2 veces: guardar lo que vino + warning al usuario (P9)
- [ ] Los topics asignados a la disciplina (coord_doc_subject_topics) se distribuyen entre las clases
- [ ] El coordinador puede editar título, objetivo y topics de cada clase **solo en estado in_progress** (P10)
- [ ] **Clases inmutables post-publicación** — solo el teacher puede editar clases en su planificación (P10/P12)
- [ ] Las clases compartidas se **coordinan por IA**: generación secuencial con contexto entre materias (P11)
- [ ] Endpoint para editar una clase individual: `PATCH /api/v1/coord-doc-classes/:id` (solo in_progress)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 4.4.1 | [Usecase: generar plan de clases](./tareas/T-4.4.1-usecase-generar-plan.md) | src/core/usecases/ | ⬜ |
| 4.4.2 | [Prompt y schema para generación](./tareas/T-4.4.2-prompt-schema.md) | src/core/usecases/ | ⬜ |
| 4.4.3 | [Endpoint PATCH clase individual](./tareas/T-4.4.3-endpoint-editar-clase.md) | src/entrypoints/ | ⬜ |
| 4.4.4 | [Integrar shared classes en generación IA](./tareas/T-4.4.4-shared-classes.md) | src/core/usecases/ | ⬜ |
| 4.4.5 | [Tests](./tareas/T-4.4.5-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-4.2: Wizard](../HU-4.2-wizard-creacion/HU-4.2-wizard-creacion.md) — Documento con subjects y topics creados
- [HU-3.5: Grilla horaria](../../03-integracion/HU-3.5-grilla-horaria-clases-compartidas/HU-3.5-grilla-horaria-clases-compartidas.md) — Para identificar clases compartidas
- [Épica 6: Asistente IA](../../06-asistente-ia/06-asistente-ia.md) — Azure OpenAI para generar el plan

## Diseño técnico

### Generación IA del plan de clases

Para cada disciplina del documento, se envía al LLM:

**Input:**
- Disciplina (nombre)
- Topics asignados a la disciplina
- class_count
- Secciones del documento ya generadas (eje problemático, estrategia)
- Clases compartidas identificadas (P11)
- Contexto de materias ya generadas si hay clases compartidas

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

### Reintento por count incorrecto (P9 — Decisión)

Si la IA genera un número diferente al `class_count` esperado:

1. **Primer intento**: prompt normal
2. **Si count != esperado**: reintentar con prompt reforzado: "IMPORTANTE: genera EXACTAMENTE {n} clases, ni más ni menos"
3. **Si el segundo intento también falla**: guardar lo que vino + warning al usuario: "Se generaron {actual} clases en vez de {expected}. Podés ajustar manualmente"

### Clases inmutables post-publicación (P10/P12 — Decisión)

- En estado `in_progress`: el coordinador puede editar título, objetivo y topics de cualquier clase via `PATCH /coord-doc-classes/:id`
- En estado `published`: las clases son **inmutables**. El endpoint retorna `403`
- Solo el **teacher** puede editar clases en su propia planificación (Épica 5)
- Alizia en el chat **no ofrece tools de edición de clases** cuando el documento está publicado

### Coordinación IA de clases compartidas (P11 — Decisión)

Si dos materias comparten un slot horario (detectado desde time_slots de HU-3.5), la IA coordina el contenido:

**Flujo de generación coordinada:**
1. Identificar pares de materias con clases compartidas
2. Generar la primera materia normalmente
3. Generar la segunda materia con contexto: "Las clases {3, 7, 15} son compartidas con {Matemática}. El plan de Matemática para esas clases es: {contexto}. Asegurá coherencia temática y complementariedad"
4. Si hay más de 2 materias con slots compartidos, generar en secuencia pasando contexto acumulado

**Marcado en respuesta:**
```json
{
  "class_number": 3,
  "title": "Interdisciplina: Matemáticas y Física",
  "is_shared": true,
  "shared_with_subject": "Física"
}
```

**Riesgos y mitigación:**
- Si la primera generación es mala, puede contaminar la segunda → el coordinador puede pedir a Alizia que regenere via chat
- Prompts más largos = más tokens → aceptable dado que es una operación infrecuente

### Manejo de errores de IA (P17 — Decisión)

- Timeout reducido (configurable)
- 1 reintento automático para errores retriables (timeout, 5xx, 429)
- Para generación multi-materia: si algunas materias se generaron OK pero otras fallan, **guardar las exitosas** y reportar error parcial
- Las secciones son pre-requisito: si secciones fallan, no intentar generar clases

## Test cases

- 4.14: POST generate → plan de clases creado para cada disciplina
- 4.15: Cantidad de clases generadas == class_count
- 4.16: Count incorrecto → reintento con prompt reforzado
- 4.17: Segundo intento también incorrecto → guardar + warning
- 4.18: Todos los topics de la disciplina distribuidos en al menos una clase
- 4.19: PATCH clase en in_progress → título y topics actualizados
- 4.20: PATCH clase en published → 403 (clases inmutables)
- 4.21: Clases compartidas: generación con contexto de otra materia
- 4.22: Clases compartidas marcadas correctamente en respuesta
- 4.23: Falla IA parcial (1 materia OK, 1 falla) → guardar exitosa, reportar error
