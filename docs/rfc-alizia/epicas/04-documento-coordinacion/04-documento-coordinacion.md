# Épica 4: Documento de coordinación

> Creación asistida, edición colaborativa y gestión de estados de documentos de coordinación areal.

**Estado:** MVP
**Fase de implementación:** Fase 3

---

## Problema

Los equipos de docentes necesitan documentos complejos (como el itinerario del área) que deben alinear múltiples personas en el trabajo a realizarse en un plazo futuro dado. Este proceso es manual, lento y difícil de articular entre roles.

## Objetivos

- Generar una primera versión del documento, que siga los lineamientos correspondientes y las mejores prácticas
- Permitir la edición colaborativa entre roles (por ej: coordinadores y docentes)
- Asegurar que el documento resultante sea la base sobre la cual se planifica el clase a clase

## Alcance MVP

**Incluye:**

- Generación asistida del documento basado en un "topic", que incluya los elementos que la provincia quiera, por ej: Eje problemático y estrategia metodológica
- Cálculo automático de cantidad de clases por disciplina según grilla horaria
- Selección de "sub-topic", donde habrá múltiples disponibles para cada disciplina, que pueden repetirse o no en otras disciplinas (configurable según provincia)
- Generación de un cronograma de clases tentativo, con objetivo por clase para cada disciplina
- Edición directa y asistida por IA del documento
- Publicación del documento para que todos los roles puedan acceder

**No incluye:**

- Planificación del clase a clase detallado → ver Planificación docente
- Creación de recursos didácticos → ver Contenido

## Principios de diseño

- **Propuesta primero:** Alicia genera una primera versión; el equipo edita y valida.
- **Alineación vertical:** Todo lo que se planifica debe poder trazarse hasta los lineamientos provinciales.
- **Colaboración entre roles:** Coordinadores definen el marco, docentes validan y ajustan el cronograma.

## Sub-épicas

| Componente | Descripción |
|---|---|
| Asistente de creación | Genera elementos del documento: eje problemático, estrategia metodológica y cronograma a partir del tema y documentación base (diseño curricular, etc.) |
| Editor de documento | Edición directa y asistida por IA de cada sección del documento |
| Flujo de estados | Gestión del ciclo de vida: borrador, en revisión, publicado |

## Decisiones de cada cliente

- Las estrategias metodológicas disponibles (proyecto, taller, ateneo, laboratorio). Además requieren validación con cada equipo pedagógico provincial
- El nivel de edición que tiene el docente sobre el documento del coordinador es decisión de cada provincia

## Decisiones técnicas

- El **período** del documento (cuatrimestre, bimestre, semestre) es un nombre libre con rango de fechas custom — no se asume una duración fija ni un calendario predeterminado. El coordinador define inicio y fin.
- Cada sección del documento se define mediante un **JSON Schema** que especifica: estructura del output, prompt de generación y opciones disponibles. Esto permite que cada provincia customice qué secciones tiene y cómo se generan, sin cambios de código.
- El documento tiene estados: **draft → published → archived**.
- **Clases coordinadas entre docentes** es un feature flag por organización. Si dos materias coinciden en horario, las modificaciones de un docente son visibles para el otro. Esto se destaca como diferenciador clave del producto — ninguna plataforma del mercado lo ofrece.
- El cálculo de clases por disciplina se basa en la grilla horaria importada. El coordinador puede sobreescribir manualmente (±) para contemplar feriados, días institucionales o situaciones no previstas por el sistema.
- La selección de sub-topics para el documento respeta la **profundidad configurada por organización**. Si la org define que se seleccionan hasta nivel 2, solo se muestran esos niveles — se diseñarán UI/UX para cada caso y el sistema se adaptará dinámicamente.

## Épicas relacionadas

- **Integración** — Provee diseño curricular y estructura que alimentan la generación
- **Planificación docente** — Consume el documento como insumo para el clase a clase
- **Asistente IA** — Motor de generación y edición del documento
