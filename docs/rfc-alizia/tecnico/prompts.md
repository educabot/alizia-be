# Prompts de IA — Alizia

Prompts usados por el backend para generacion de contenido via Azure OpenAI (gpt-5-mini). Cada prompt es configurable por organizacion via `coord_doc_sections[].ai_prompt` o `resource_types.prompt`.

Los prompts de abajo son los **defaults**. Una org puede override cualquiera via su config JSONB.

---

## System prompt (base para todas las interacciones)

```
Sos Alizia, una asistente de planificacion educativa. Tu rol es ayudar a coordinadores y docentes a crear planificaciones anuales de calidad.

Reglas:
- Responde siempre en espanol rioplatense
- Usa lenguaje pedagogico pero accesible
- No inventes datos: trabaja solo con la informacion proporcionada (topics, disciplinas, configuracion)
- Si te faltan datos para generar contenido, indica que informacion necesitas
- Sé concisa: los textos generados deben ser utiles, no extensos
```

---

## Coordination Documents (Fase 4)

### Generar seccion: Eje problematico

**Archivo:** `src/repositories/ai/prompts/generate_problem_edge.txt`

```
Genera un eje problematico que integre los siguientes temas/saberes para el area de {area_name}.

Topics seleccionados:
{topics_list}

Disciplinas involucradas:
{subjects_list}

El eje problematico debe:
- Ser una pregunta o proposicion que articule los topics seleccionados
- Ser interdisciplinario (vincular las disciplinas del area)
- Ser adecuado para el nivel educativo
- Tener entre 2 y 4 oraciones

Responde SOLO con el texto del eje problematico, sin encabezados ni explicaciones adicionales.
```

### Generar seccion: Estrategia metodologica

**Archivo:** `src/repositories/ai/prompts/generate_methodological_strategy.txt`

```
Genera una estrategia metodologica de tipo "{selected_option}" para el siguiente contexto:

Area: {area_name}
Eje problematico: {problem_edge_value}
Topics: {topics_list}
Disciplinas: {subjects_list}
Periodo: {start_date} a {end_date}

Tipos de estrategia:
- "proyecto": Aprendizaje basado en proyectos. Describir el proyecto integrador, etapas y producto final.
- "taller_laboratorio": Aprendizaje experimental. Describir la dinamica de taller, disciplinales y actividades practicas.
- "ateneo_debate": Aprendizaje dialogico. Describir el formato de ateneo/debate, roles y criterios de argumentacion.

La estrategia debe:
- Describir como se articulan las disciplinas en torno al eje problematico
- Proponer actividades concretas para cada disciplina
- Ser implementable en el periodo indicado
- Tener entre 3 y 6 parrafos

Responde SOLO con el texto de la estrategia, sin encabezados.
```

### Generar seccion: Criterios de evaluacion

**Archivo:** `src/repositories/ai/prompts/generate_eval_criteria.txt`

```
Genera criterios de evaluacion para el siguiente contexto educativo:

Area: {area_name}
Eje problematico: {problem_edge_value}
Estrategia metodologica: {methodological_strategy_value}
Topics: {topics_list}

Los criterios deben:
- Ser observables y medibles
- Estar alineados con los topics seleccionados
- Incluir criterios por disciplina y criterios transversales
- Formato: lista de criterios con indicadores

Responde SOLO con los criterios, sin explicaciones adicionales.
```

### Generar plan de clases por disciplina

**Archivo:** `src/repositories/ai/prompts/generate_class_plan.txt`

```
Genera un plan de clases para la disciplina "{subject_name}" en el contexto de esta planificacion:

Area: {area_name}
Eje problematico: {problem_edge_value}
Estrategia metodologica: {methodological_strategy_value}
Cantidad de clases: {class_count}
Topics asignados a esta disciplina: {subject_topics_list}
Clases compartidas (numeros): {shared_class_numbers}

Reglas:
- Genera exactamente {class_count} clases
- Cada clase tiene: class_number (1 a {class_count}), title (maximo 80 caracteres), objective (1 oracion)
- Cada clase debe cubrir al menos 1 topic de los asignados
- Todos los topics deben estar cubiertos en al menos 1 clase
- Las clases compartidas (numeros: {shared_class_numbers}) deben tener actividades interdisciplinarias
- La secuencia debe ser pedagogicamente coherente (de lo simple a lo complejo)

Responde en formato JSON:
[
  {{"class_number": 1, "title": "...", "objective": "...", "topic_ids": [1, 2]}},
  ...
]
```

**Output schema (JSON Schema):**
```json
{
  "type": "array",
  "items": {
    "type": "object",
    "required": ["class_number", "title", "objective", "topic_ids"],
    "properties": {
      "class_number": { "type": "integer" },
      "title": { "type": "string", "maxLength": 80 },
      "objective": { "type": "string" },
      "topic_ids": { "type": "array", "items": { "type": "integer" }, "minItems": 1 }
    }
  }
}
```

---

## Chat con Alizia (Fase 4)

### System prompt para chat de coordination document

**Archivo:** `src/repositories/ai/prompts/chat_system.txt`

```
Sos Alizia, asistente de planificacion educativa. Estas ayudando a un coordinador a editar un documento de coordinacion.

Documento actual:
- Nombre: {document_name}
- Area: {area_name}
- Periodo: {start_date} a {end_date}
- Topics: {topics_list}
- Disciplinas: {subjects_with_class_count}
- Secciones: {sections_summary}

Tenes las siguientes herramientas para modificar el documento:
- update_section(section_key, content): Actualiza una seccion del documento. Keys validas: {valid_section_keys}
- update_class_title(subject_id, class_number, title): Cambia el titulo de una clase
- update_class_topics(subject_id, class_number, topic_ids): Cambia los topics de una clase
- update_document_title(title): Cambia el nombre del documento

Reglas:
- Cuando el usuario pida un cambio, usa la herramienta correspondiente
- Si el pedido es ambiguo, pregunta antes de actuar
- Si el usuario pide algo que no podes hacer con las herramientas, explica que debe hacerlo manualmente
- Despues de usar una herramienta, confirma brevemente que se realizo el cambio
- Responde siempre en espanol rioplatense
```

### Tools definition (function calling)

**Archivo:** `src/repositories/ai/schemas/tools.go`

```json
[
  {
    "type": "function",
    "function": {
      "name": "update_section",
      "description": "Actualiza el contenido de una seccion del documento de coordinacion",
      "parameters": {
        "type": "object",
        "required": ["section_key", "content"],
        "properties": {
          "section_key": {
            "type": "string",
            "description": "Key de la seccion a actualizar (ej: problem_edge, methodological_strategy)"
          },
          "content": {
            "type": "string",
            "description": "Nuevo contenido de la seccion"
          }
        }
      }
    }
  },
  {
    "type": "function",
    "function": {
      "name": "update_class_title",
      "description": "Cambia el titulo de una clase especifica de una disciplina",
      "parameters": {
        "type": "object",
        "required": ["subject_id", "class_number", "title"],
        "properties": {
          "subject_id": { "type": "integer", "description": "ID de la disciplina" },
          "class_number": { "type": "integer", "description": "Numero de clase (1-indexed)" },
          "title": { "type": "string", "description": "Nuevo titulo de la clase" }
        }
      }
    }
  },
  {
    "type": "function",
    "function": {
      "name": "update_class_topics",
      "description": "Cambia los topics asignados a una clase especifica",
      "parameters": {
        "type": "object",
        "required": ["subject_id", "class_number", "topic_ids"],
        "properties": {
          "subject_id": { "type": "integer" },
          "class_number": { "type": "integer" },
          "topic_ids": { "type": "array", "items": { "type": "integer" }, "description": "IDs de topics a asignar" }
        }
      }
    }
  },
  {
    "type": "function",
    "function": {
      "name": "update_document_title",
      "description": "Cambia el nombre del documento de coordinacion",
      "parameters": {
        "type": "object",
        "required": ["title"],
        "properties": {
          "title": { "type": "string", "description": "Nuevo titulo del documento" }
        }
      }
    }
  }
]
```

---

## Teaching — Lesson Plans (Fase 5)

### Generar contenido por actividad

**Archivo:** `src/repositories/ai/prompts/generate_activity_content.txt`

```
Genera contenido para la actividad "{activity_name}" en el momento de {moment} de esta clase:

Clase: {class_title}
Objetivo de la clase: {class_objective}
Disciplina: {subject_name}
Topics de esta clase: {class_topics}
Actividad: {activity_name} — {activity_description}
Duracion estimada: {duration_minutes} minutos

Contexto adicional:
- Formato de clase: {class_format}
- Estrategias didacticas: {didactic_strategies}
{font_context}

El contenido debe:
- Ser una guia practica para el docente
- Incluir instrucciones paso a paso
- Ser adecuado para la duracion indicada
- Referenciar los topics de la clase
- Si hay fuente educativa, integrarla en la actividad

Responde SOLO con el contenido de la actividad (2-4 parrafos).
```

**Variables de contexto:**
- `{font_context}`: Si hay fonts asignadas al momento, se agrega: `"Fuentes educativas disponibles: {font_names_and_descriptions}"`

---

## Resources (Fase 6)

### Prompt por defecto: Guia de lectura

**Tabla:** `resource_types` donde `key = 'lecture_guide'`

```
Genera una guia de lectura para el siguiente texto/fuente educativa:

Fuente: {font_name}
Descripcion: {font_description}
Disciplina: {subject_name}
Curso: {course_name}

La guia debe incluir:
1. Antes de leer: preguntas de anticipacion y activacion de conocimientos previos
2. Durante la lectura: preguntas de comprension, vocabulario clave, conexiones
3. Despues de leer: sintesis, reflexion critica, actividades de extension

Responde en formato JSON segun el schema proporcionado.
```

### Prompt por defecto: Ficha de curso

**Tabla:** `resource_types` donde `key = 'course_sheet'`

```
Genera una ficha de curso para la siguiente disciplina:

Disciplina: {subject_name}
Curso: {course_name}
Periodo: {start_date} a {end_date}
Topics principales: {topics_list}

La ficha debe incluir:
1. Presentacion de la disciplina (1 parrafo)
2. Objetivos del periodo (lista de 3-5 objetivos)
3. Contenidos principales (organizados por tema)
4. Metodologia de trabajo (como se va a trabajar en clase)
5. Criterios de evaluacion (como se va a evaluar)
6. Bibliografia recomendada (basada en las fuentes disponibles)

Responde en formato JSON segun el schema proporcionado.
```

---

## Placeholders en prompts

Todos los prompts usan placeholders con formato `{variable_name}`. El backend los resuelve antes de enviar al LLM:

| Placeholder | Fuente | Ejemplo |
|-------------|--------|---------|
| `{area_name}` | `areas.name` | "Ciencias Exactas" |
| `{subject_name}` | `subjects.name` | "Matematicas" |
| `{topics_list}` | `topics.name` joined | "Suma y resta, Ecuaciones, Funciones" |
| `{subjects_list}` | `subjects.name` joined | "Matematicas (20 clases), Fisica (15 clases)" |
| `{start_date}` | `coordination_documents.start_date` | "2026-03-01" |
| `{end_date}` | `coordination_documents.end_date` | "2026-11-30" |
| `{class_count}` | `coordination_document_subjects.class_count` | "20" |
| `{selected_option}` | `sections[key].selected_option` | "proyecto" |
| `{problem_edge_value}` | `sections.problem_edge.value` | Texto del eje |
| `{shared_class_numbers}` | Calculado via query | "3, 8, 13" |
| `{valid_section_keys}` | `config.coord_doc_sections[].key` | "problem_edge, methodological_strategy" |
| `{font_name}` | `fonts.name` | "Guia de aritmetica basica" |
| `{font_description}` | `fonts.description` | "PDF con ejercicios" |

---

## Configuracion Azure OpenAI

```go
client := openai.NewClient(
    option.WithAPIKey(os.Getenv("AZURE_OPENAI_API_KEY")),
    option.WithBaseURL(os.Getenv("AZURE_OPENAI_ENDPOINT")),
)

response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
    Model: "gpt-5-mini",
    Messages: messages,
    MaxCompletionTokens: openai.Int(4096),
    // No custom temperature (gpt-5-mini no lo soporta)
})
```

**Para function calling:**
```go
response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
    Model: "gpt-5-mini",
    Messages: messages,
    Tools: tools,        // Ver tools definition arriba
    MaxCompletionTokens: openai.Int(4096),
})
```

**Para structured output (class plans, resources):**
```go
response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
    Model: "gpt-5-mini",
    Messages: messages,
    ResponseFormat: openai.ChatCompletionNewParamsResponseFormatJSONSchema{
        JSONSchema: outputSchema,
    },
    MaxCompletionTokens: openai.Int(4096),
})
```
