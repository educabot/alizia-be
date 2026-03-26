# Épica 6: Asistente IA

> Asistente de inteligencia artificial que genera, edita y aprende del uso de la plataforma.

**Estado:** MVP
**Fase de implementación:** Fase 4

---

## Problema

Los docentes y coordinadores necesitan producir documentos y planificaciones alineadas curricularmente, pero el proceso es lento y requiere expertise. Un asistente genérico no conoce el contexto provincial ni el historial del aula.

## Objetivos

- Generar primeras versiones de documentos y planificaciones basadas en el diseño curricular
- Permitir ediciones masivas o puntuales por instrucción natural del usuario
- Incorporar el historial de clases (bitácora) para mejorar las recomendaciones futuras
- Mantener todas las generaciones alineadas con las fuentes oficiales de la provincia

## Alcance MVP

**Incluye:**

- Generación de documentos de coordinación (secciones + plan de clases)
- Recomendación de actividades para cada momento de una clase
- Edición asistida: el usuario da instrucciones en lenguaje natural y Alicia modifica el documento
- Procesamiento de feedback post-clase (bitácora de cotejo) para ajustar futuras propuestas (post-MVP)

## Principios de diseño

- **Curada y validada:** Las generaciones se basan en fuentes oficiales y criterio pedagógico colectivo.
- **Propuesta, no imposición:** Alicia propone; el docente decide.
- **Memoria del aula:** El feedback de clases anteriores informa las recomendaciones.

## Sub-épicas

| Componente | Descripción |
|---|---|
| Modificación de contenido | Edición asistida de documentos y planificaciones por instrucción natural |
| Asistencia de uso y navegación | Ayuda contextual dentro de la plataforma |
| Customización por cliente | Adaptación del comportamiento del asistente según la provincia o institución |

## Decisiones del cliente

- La customización por cliente (tono, límites, fuentes permitidas) requiere definición por provincia

## Decisiones técnicas

- El asistente opera como una **LLM con tools** (function calling). Puede leer y modificar el documento que el usuario está editando, recomendar actividades, consultar la estructura curricular, y acceder a la bitácora de clases anteriores.
- Cada sección generada por IA se define con un **prompt + JSON Schema** por organización. El prompt incorpora variables contextuales (tópicos seleccionados, disciplina, grilla horaria) y el schema fuerza el formato del output. Esto permite customización por provincia sin cambios de código.
- El **planificador de clases** recibe todos los temas asignados al documento, las disciplinas involucradas y la cantidad de clases disponibles, y distribuye el contenido en el tiempo. Si hay mucho tiempo, lo esparce; si hay poco, lo compacta.
- La generación de secciones del documento de coordinación (eje problemático, estrategia metodológica, criterios de evaluación) arranca con un prompt simple. Se prevé que habrá **mucho prueba y error** — la sofisticación vendrá con el uso real, no con sobre-ingeniería anticipada.
- El asistente de chat permite al usuario pedir modificaciones en lenguaje natural ("cambiá la actividad del cierre por algo más dinámico") y el sistema modifica la sección correspondiente. Internamente es un chat con tools que opera sobre el documento activo.

## Function calling tools

| Tool | Descripción |
|---|---|
| `update_section(section_key, content)` | Actualiza una sección del documento. Valida que `section_key` exista en schema de la org |
| `update_class(class_number, title, objective)` | Modifica una clase del plan |
| `update_class_topics(class_number, topic_ids)` | Cambia topics de una clase |

## Épicas relacionadas

- **Documento de coordinación** — Alicia genera y edita el documento
- **Planificación docente** — Alicia propone actividades y procesa la bitácora
- **Contenido** — Alicia genera recursos didácticos
