# Preguntas Abiertas — Epica 4: Documento de Coordinacion

> **Objetivo:** Definir decisiones de diseno pendientes antes de implementar.
> Cada pregunta tiene entre 1 y 4 alternativas segun corresponda.
>
> **Fecha:** 2026-04-22
> **Participantes:** Juan, Jose, Sebastian

---

## Contexto rapido

El **documento de coordinacion** es el plan maestro de area que define:

1. **QUE** ensenar (seleccion de topics del curriculum provincial)
2. **DONDE** (distribucion de topics entre materias del area)
3. **CUANDO** (plan clase por clase por materia, generado con IA)

Mas secciones narrativas configurables por org (eje problematico, estrategia metodologica, criterios de evaluacion).

**Flujo del coordinador:**

```
CREAR (wizard) → GENERAR (IA) → EDITAR (humano + IA) → PUBLICAR (validar) → AJUSTAR?
   pending?        → in_progress      in_progress        → published          published?
```

---

## Resumen de preguntas

| #       | Pregunta corta                           | Alternativas | HU    |
| ------- | ---------------------------------------- | ------------ | ----- |
| P1      | Estado pending necesario?                | 3            | 4.1   |
| P2      | Multiples docs por area?                 | 3            | 4.1   |
| P3      | Sin grilla horaria?                      | 1 (obvia)    | 4.2   |
| P4      | Topics compartidos entre materias?       | 1 (obvia)    | 4.2   |
| P5      | Class count = 0?                         | 1 (obvia)    | 4.2   |
| P6      | select_text: cuando elige la opcion?     | 3            | 4.3   |
| P7      | Granularidad de re-generacion?           | 4            | 4.3/4 |
| P8      | Tipos de seccion para MVP?               | 1 (obvia)    | 4.3   |
| P9      | IA genera count incorrecto de clases?    | 2            | 4.4   |
| P10     | CRUD manual de clases?                   | 2            | 4.4   |
| P11     | Clases compartidas en generacion IA?     | 2            | 4.4   |
| **P12** | **Publicado es editable?**               | **4**        | **4.5** |
| P13     | Que se valida al publicar?               | 3            | 4.5   |
| P14     | Docentes ven no-publicados?              | 1 (obvia)    | 4.5   |
| P15     | Persistir historial de chat?             | 2            | 4.6   |
| P16     | Que tools tiene Alizia en el chat?       | 3            | 4.6   |
| P17     | Falla de IA: como se maneja?             | 2            | 4.6   |
| P18     | Tool calls atomicos?                     | 1 (obvia)    | 4.6   |
| P19     | Edicion concurrente?                     | 2            | Cross |
| P20     | Undo/versionado?                         | 3            | Cross |

> **P12 es la pregunta central.** La respuesta condiciona P7, P10, P15, P16, P19 y P20.

---

## HU-4.1 — Modelo de Datos

### P1: El estado `pending` es necesario?

El flujo dice `pending -> in_progress -> published`. Pero no queda claro **quien** ni **cuando** transiciona de pending a in_progress.

**A) Eliminar `pending` — solo 2 estados**
- El documento se crea directamente en `in_progress`
- Simplifica todo: menos logica de transicion, menos estados que testear
- El wizard YA recopila datos sustanciales (topics, materias, class_count) — el documento nunca esta realmente "vacio"

**B) Mantener `pending` — transicion manual**
- `pending` = "cree el documento pero todavia no empece a trabajarlo"
- El coordinador aprieta "Comenzar edicion" para pasar a `in_progress`
- Util si el coordinador quiere dejar documentos creados para despues
- Pregunta: si el wizard ya recopila todo, cuando estaria un doc en `pending` y por que?

**C) Mantener `pending` — transicion automatica**
- `pending` = recien creado por el wizard
- Pasa a `in_progress` automaticamente al primer generate o primera edicion
- `pending` es solo un marcador de "todavia no se genero contenido IA"
- Permite distinguir "creado pero sin contenido IA" vs "ya tiene contenido"

**Decision:** `[ Mantener el pending pero en el momento que hacer click en el tema o arrancas a escribir algo pasas a in progrres de manera automatica. Solo esta en pending cuando est recien creado y nunca edito nada]`

---

### P2: Puede existir mas de un documento activo por area?

El RFC no define unicidad. Carlos podria crear 3 documentos para "Ciencias Exactas".

**A) Sin restriccion**
- Multiples documentos en cualquier estado por area
- Maximo flexibility: borradores, versiones alternativas, documentos por periodo
- Riesgo: los docentes no saben cual mirar

**B) Un `published` por area, multiples borradores**
- Solo un documento `published` por area a la vez
- Si publicas uno nuevo, el anterior se archiva automaticamente
- Multiples `in_progress`/`pending` permitidos (borradores)
- Los docentes siempre ven exactamente un documento vigente

**C) Un activo por area/periodo**
- UNIQUE(area_id, period) para docs no-archivados
- Fuerza que cada periodo tenga un solo documento
- Mas restrictivo pero mas claro para todos

**Decision:** `[ Ser flexibles , sin restricciones mostrar todo desde el backend y desde el front dar reglas o cosas visuales para que esten todos conviviendo a la vez. Ademas sumar la posibilidad del estado de archived ]`

---

## HU-4.2 — Wizard de Creacion

### P3: Que pasa si no hay grilla horaria cargada? — OBVIA

**Recomendacion: permitir carga manual con sugerencia cuando hay grilla.**

Si hay grilla → pre-llena `class_count` con el calculo automatico.
Si no hay grilla → campo vacio, coordinador ingresa a mano.

Bloquear el wizard porque el admin no cargo la grilla es mala UX y no tiene beneficio real. El coordinador sabe cuantas clases tiene — la grilla es una conveniencia, no un requisito.

**Confirmar o corregir:** `[ A priori si , eviar casos ]`

---

### P4: Los topics pueden repetirse entre materias? — OBVIA

**Recomendacion: si, permitido siempre.**

La interdisciplinariedad es un pilar del modelo pedagogico. El mismo topic "Sustentabilidad" puede (y deberia) aparecer en Ciencias y en Humanidades.

El schema ya lo soporta (no hay UNIQUE que lo impida). Agregar un flag configurable no aporta valor para MVP y agrega complejidad.

**Confirmar o corregir:** `[ SI]`

---

### P5: El class_count puede ser 0? — OBVIA

**Recomendacion: no. Minimo 1.**

Si una materia no tiene clases, no tiene sentido incluirla en el documento. El coordinador simplemente no la agrega en el paso 2 del wizard.

class_count = 0 crea edge cases raros: que muestra el plan de clases? como distribuye topics? Complejidad innecesaria.

**Confirmar o corregir:** `[ SI ]`

---

## HU-4.3 — Secciones Dinamicas

### P6: Para `select_text`, cuando elige Carlos la opcion?

Ejemplo: "Estrategia metodologica" tiene opciones: proyecto, taller, ateneo. La IA necesita saber cual para generar contenido relevante.

**A) Antes de generar (pre-requisito)**
- Carlos selecciona "proyecto" ANTES de apretar "Generar"
- Si hay secciones `select_text` sin opcion seleccionada, el boton "Generar" muestra: "Selecciona las opciones pendientes"
- Mas barato en tokens (1 generacion por seccion)
- UX clara: el coordinador toma la decision, la IA ejecuta

**B) Alizia sugiere la opcion y genera**
- El coordinador NO elige previamente
- Alizia analiza los topics seleccionados y sugiere "proyecto" como la mejor opcion, generando el contenido
- El coordinador puede cambiar la opcion despues y regenerar esa seccion
- Mas "magico" pero el coordinador pierde agencia en la decision inicial

**C) Generar las N variantes y dejar elegir**
- Alizia genera una version para cada opcion (proyecto, taller, ateneo)
- Carlos ve las 3, elige la que mas le gusta
- UX exploratoria: el coordinador descubre que estrategia le sirve
- 3x mas caro en tokens, 3x mas latencia
- Muy buena UX pero costosa

**Decision:** `[ Opcion C con posibilidad de una opcion mas para agregar una opcion escrita a mano por el coordinador ]`

---

### P7: Como funciona la re-generacion?

Esta pregunta combina las anteriores P7 (sobreescribe?) y P12 (por materia o todo?). Son la misma decision de diseno.

Carlos genera contenido, edita manualmente, y quiere regenerar algo.

**A) Todo o nada + warning**
- POST /generate regenera TODO: todas las secciones + plan de clases de todas las materias
- Warning previo: "Esto va a reemplazar todo el contenido generado y tus ediciones. Continuar?"
- Simple de implementar. Un boton, una accion
- Riesgo: el coordinador pierde trabajo. Pero es claro y predecible

**B) Separado: secciones vs plan de clases**
- Dos botones: "Regenerar secciones" y "Regenerar plan de clases"
- El coordinador puede regenerar secciones sin tocar las clases y viceversa
- Complejidad media. Cubre el caso comun: "las secciones estan bien, pero el plan no"

**C) Granular por seccion y por materia**
- Boton "Regenerar" por cada seccion individual
- Boton "Regenerar plan" por cada materia individual
- Maximo control: regenerar solo Fisica sin tocar Matematica
- Mas complejo pero protege las ediciones manuales de lo que ya esta bien
- Endpoint: POST /generate con body `{ "sections": ["problem_edge"], "subject_ids": [2] }`

**D) Mixto: todo junto para MVP, granular despues**
- MVP: opcion A (todo o nada + warning)
- v2: migrar a opcion C (granular)
- Permite avanzar rapido sin cerrar puertas

**Decision:** `[ Opcion E como las clases no se pueden tocar si la edicion es con IA que la tool te pemita editar las secciones editables sin ningun boton adicional necesario, no hay un boton regenerar, el chat ALizia te lo regenera hablando con ella]`

---

### P8: Tipos de seccion para MVP? — OBVIA

**Recomendacion: solo `text` y `select_text`.**

Cubren los 3 casos del seed de Neuquen (eje problematico, estrategia metodologica, criterios de evaluacion). Si alguna org necesita otro tipo (checklist, multi_select, etc.), se agrega cuando aparezca el caso real.

No disenar para lo hipotetico.

**Confirmar o corregir:** `[ SI ]`

---

## HU-4.4 — Plan de Clases

### P9: Que pasa si la IA genera un numero incorrecto de clases?

Carlos pidio 48 clases pero la IA devuelve 45 o 51.

**A) Reintentar 1 vez con prompt reforzado**
- Si el count no matchea, reintentar agregando al prompt: "IMPORTANTE: genera EXACTAMENTE 48 clases"
- Si falla 2 veces, guardar lo que vino + warning al usuario con la diferencia
- Balance: intenta ser correcto pero no se rompe si la IA falla

**B) Aceptar siempre + ajustar class_count**
- Guardar las clases que vengan
- Actualizar `class_count` al numero real generado
- Warning: "Se generaron 45 clases en vez de 48. Podes agregar las faltantes manualmente (si P10=B) o regenerar"
- Mas resiliente, menos latencia
- Riesgo: el class_count pierde significado como "lo que el coordinador definio"

**Decision:** `[ Opcion A ]`

---

### P10: El coordinador puede agregar/quitar clases manualmente?

Hoy solo puede editar titulo, objetivo y topics de clases existentes.

**A) No en MVP — solo editar existentes**
- PATCH titulo, objetivo, topics. Nada mas
- Para cambiar cantidad: regenerar (ver P7)
- Simple de implementar
- Limitante si el coordinador necesita 1 clase mas o menos

**B) CRUD completo en MVP**
- POST /coord-doc-classes — agregar clase (al final o en posicion especifica)
- DELETE /coord-doc-classes/:id — eliminar clase
- Auto-renumera class_number al agregar/eliminar
- Mas control, mejor UX
- Mas endpoints y validaciones
- Habilita que P9 sea mas flexible (si faltan clases, el coordinador las agrega)

**Decision:** `[ C no puede editar las clases una vez publicadas, las puede editar el teacher solo ]`

---

### P11: Clases compartidas — como se manejan en la generacion IA?

Si Matematica y Fisica comparten un slot horario, ciertas clases son "compartidas".

**A) Solo marcar, no coordinar (MVP)**
- La IA genera cada materia independientemente
- En el GET response, las clases compartidas se marcan: `is_shared: true, shared_with: "Fisica"`
- La coordinacion de contenido queda en el coordinador humano
- Post-MVP: la IA podria coordinar
- Simple, predecible, no agrega riesgo a la generacion

**B) IA coordina las clases compartidas**
- El prompt incluye contexto: "Las clases 3, 7 y 15 son compartidas con Fisica. Asegura coherencia"
- Genera las dos materias en secuencia: primero Matematica, despues Fisica con contexto de lo que genero para Mate
- Mejor resultado pedagogico
- Significativamente mas complejo y mas tokens
- Riesgo: si la primera generacion es mala, contamina la segunda

**Decision:** `[ Opcion B ]`

---

## HU-4.5 — Publicacion y Estados

### P12: Un documento publicado es editable o no? ⚠️ CRITICA

Hay una **contradiccion en el RFC actual**:
- El overview dice "documento vivo, editable post-publicacion"
- Los endpoints de edicion validan `status == in_progress` y devuelven 403 si es published

Esta es la decision mas importante de la epica. Condiciona P7, P10, P15, P16, P19 y P20.

**A) Publicado = inmutable**
- Una vez publicado, no se toca. Punto
- Si algo cambio, se crea un documento nuevo para el siguiente periodo
- Los docentes siempre trabajan sobre una version estable
- Simple, predecible, sin riesgo de inconsistencia
- Riesgo: rigido. En la realidad, un coordinador NECESITA ajustar cosas durante el cuatrimestre

**B) Publicado = editable con warning**
- Se puede editar todo (secciones, clases, chat con Alizia) estando publicado
- Cada edicion incluye en el response: `"warning": "Los cambios no se propagan a planificaciones docentes ya creadas"`
- No hay notificacion al docente (MVP). El docente ve la version actualizada si entra al documento
- Flexible, refleja la realidad
- Riesgo: el docente hizo un lesson plan basado en clase 15, el coordinador la cambio, el lesson plan quedo desalineado

**C) Publicado = editable solo secciones, clases inmutables**
- Las secciones narrativas (eje problematico, estrategia, etc.) se pueden editar post-publicacion
- El plan de clases (titulo, objetivo, topics por clase) queda fijo al publicar
- Razon: los docentes planifican sobre las clases, no sobre las secciones
- Compromiso entre flexibilidad y estabilidad
- Las secciones son "marco general", las clases son "lo operativo"

**D) Publicado con "despublicar"**
- Publicado es inmutable PERO se puede "despublicar" (volver a in_progress)
- El coordinador edita lo que necesite y re-publica
- Mientras esta despublicado, los docentes ven la ultima version publicada (o un aviso)
- Mas complejo pero da control total sin ambiguedad
- Requiere nuevo estado o logica de "version publicada" vs "version borrador"

**Decision:** `[ Opcion C las clases son inmutables pero no se va a progragar nada, dejar un warning]`

---

### P13: Que se valida al publicar?

**A) Minimo: topics distribuidos + secciones requeridas**
- Cada topic del documento asignado a al menos 1 materia
- Cada seccion con `required: true` tiene contenido no vacio
- NO valida nada a nivel clase (si hay clases generadas, si los topics aparecen en alguna clase)
- Rapido de publicar, responsabilidad en el coordinador

**B) Medio: A + clases generadas**
- Todo lo de A
- Ademas: cada materia debe tener al menos 1 clase generada (no se puede publicar sin plan de clases)
- Garantiza que el documento tiene contenido real, no solo metadata

**C) Estricto: B + cobertura de topics en clases**
- Todo lo de B
- Ademas: cada topic asignado a una materia debe aparecer en al menos 1 clase de esa materia
- Garantiza plan completo y trazable
- Riesgo: demasiado estricto? bloquea publicacion por un topic que el coordinador dejo intencionalmente para "agregar despues"
- Podria ser **warning** en vez de error bloqueante

**Decision:** `[Cada sub topic debe ser asignado al menos a 1 materia warning u opcional? Te quedo tal sub topic sin asignar , queres contnuar? Recordemos que un documento de coordinaciones esta relacionado a un topic padre y las materias relacionadas a un topic hijo ]`

---

### P14: Los docentes pueden ver documentos no publicados? — OBVIA

**Recomendacion: no. Solo documentos publicados.**

Antes de publicar es trabajo interno del coordinador. El docente no tiene por que ver borradores incompletos ni opinar sobre un documento a medio armar. El momento de feedback es despues de publicar.

Si en el futuro se necesita un flujo de "revision por docentes antes de publicar", se agrega como feature explicito, no como acceso silencioso.

**Confirmar o corregir:** `[SI efectivamente no pueden verlos ]`

---

## HU-4.6 — Chat con Alizia

### P15: El historial del chat se persiste en el backend?

Hoy el frontend manda el historial en cada request. Si Carlos cierra el browser, pierde la conversacion.

**A) No persistir (MVP)**
- Frontend maneja el historial en memoria o localStorage
- Simple de implementar
- Los cambios de Alizia YA se persisten en el documento — lo que se pierde es solo la conversacion
- Si Carlos cierra y vuelve, puede seguir chateando (sin contexto de lo anterior)
- Considerar: localStorage persiste entre tabs y sesiones del mismo browser

**B) Persistir en backend**
- Tabla `coord_doc_chat_messages` (doc_id, role, content, tool_calls, created_at)
- GET /coordination-documents/:id/chat/history para retomar
- El coordinador puede retomar la conversacion en otro dispositivo
- Agrega: nueva tabla, nuevo endpoint, manejo de contexto largo (truncar/resumir si crece mucho)
- Pregunta: tiene valor real la conversacion o solo importa el resultado en el documento?

**Decision:** `[ Opcion B, dejar anotado que pasa cuando la conversaciones muy larga (auto compact con un limit ) ]`

---

### P16: Que tools tiene Alizia en el chat?

Tools actuales en el RFC: `update_section`, `append_to_section`, `update_class_title`, `update_class_topics`.

**A) Set minimo (4 tools)**
- `update_section(section_key, content)` — reescribir seccion
- `append_to_section(section_key, content)` — agregar al final de seccion
- `update_class(class_id, title?, objective?, topic_ids?)` — editar clase (unificado en 1 tool)
- 3 tools en vez de 4 (merge update_class_title + update_class_topics en uno)
- Si Carlos pide algo fuera de scope, Alizia responde que no puede
- Cubre el 80% de los casos

**B) Set medio (6-7 tools)**
- Todo lo de A
- `regenerate_section(section_key)` — regenerar una seccion individual con IA
- `regenerate_subject_plan(subject_id)` — regenerar plan de una materia
- `swap_classes(class_id_1, class_id_2)` — intercambiar orden de 2 clases
- Cubre el 95% de los casos
- Depende de P7 (si se permite regeneracion granular)

**C) Set extendido (si P10 = B)**
- Todo lo de B
- `add_class(subject_id, after_class_number, title, objective, topic_ids)` — agregar clase
- `remove_class(class_id)` — eliminar clase
- Maximo poder pero mas superficie de error de la IA

**Decision:** `[  Con A deberias poder hacer todo, aclarar cuando el documento este publish, alizia deberia decir que no puede editar clases ]`

---

### P17: Falla de IA — como se maneja?

Aplica a generacion (HU-4.3/4.4) y chat (HU-4.6).

**A) 1 reintento automatico + guardar parcial**
- Timeout: 60s por llamada a Azure OpenAI
- Si falla (timeout, 5xx, rate limit): reintentar 1 vez
- Si falla 2 veces: 503 con mensaje amigable
- Para generacion multi-paso (secciones + clases): si secciones se generaron OK pero clases fallan, guardar las secciones y reportar error parcial
- Balance entre resiliencia y simplicidad

**B) Sin reintento, error directo**
- Si falla, error inmediato al usuario: "El servicio de IA no esta disponible"
- El usuario reintenta manualmente
- Mas simple, menos "magia"
- Riesgo: errores transitorios (rate limit momentaneo) frustran al usuario

**Decision:** `[ Opcion A bajando el time out mostrando antes ]`

---

### P18: Tool calls atomicos? — OBVIA

**Recomendacion: best-effort (ejecutar todas, reportar fallos).**

Si Alizia ejecuta 3 tools y la segunda falla:
- Tool 1 (update_section): OK, se persiste
- Tool 2 (update_class_title con class_id invalido): falla, se reporta
- Tool 3 (update_class_topics): OK, se persiste

La respuesta incluye que acciones se ejecutaron y cuales fallaron. El coordinador ve el resultado y puede corregir.

Razon: las tools son independientes entre si. Si "actualizar eje problematico" funciono y "cambiar titulo clase 5" fallo por ID invalido, no tiene sentido rollbackear el eje problematico. Seria peor perder un cambio bueno por un error no relacionado.

**Confirmar o corregir:** `[SI mas retry]`

---

## Cross-cutting (toda la epica)

### P19: Edicion concurrente (dos tabs, dos coordinadores)?

Carlos tiene dos tabs abiertas, o Pedro tambien es coordinador del area.

**A) Last write wins (MVP)**
- Sin locking, sin deteccion de conflictos
- El ultimo en guardar gana
- Simple
- Para MVP es aceptable: es raro que dos coordinadores editen el mismo documento al mismo tiempo
- Riesgo real bajo para MVP; alto para produccion con muchos usuarios

**B) Optimistic locking**
- Campo `version` (integer, auto-increment) en el documento
- PATCH envia la version que leyo; si cambio, 409 Conflict
- Frontend muestra: "El documento fue modificado. Recarga para ver los cambios"
- Mas seguro, complejidad moderada (un campo extra + un check en cada PATCH)
- Pregunta: vale la pena para MVP o es over-engineering?

**Decision:** `[ HAcemos auto guardado? SI lo hacemos cuando y como ? Cuando buscamos la nueva version ? En tiempo real estoy haciendo un GET cada N segundos?  ]`

---

### P20: Undo/versionado de cambios?

Carlos regenero y perdio todas las ediciones manuales. O hizo un cambio via chat que rompio algo.

**A) No en MVP — solo warnings preventivos**
- Sin undo, sin versionado, sin snapshots
- Warning claro antes de acciones destructivas: "Esto va a reemplazar el contenido. Continuar?"
- El coordinador asume la responsabilidad
- Simple, no agrega tablas ni logica

**B) Snapshot automatico antes de generar**
- Antes de cada POST /generate, se guarda copia de `sections` + `classes` en tabla `coord_doc_snapshots`
- Boton "Restaurar ultima version" que vuelve al snapshot
- Solo 1 snapshot (el ultimo). No es un historial completo
- Complejidad moderada: 1 tabla, 1 endpoint, logica de restore
- Protege contra el caso mas comun: "regenere y perdi todo"

**C) Historial completo de versiones**
- Cada cambio crea una version (snapshot completo del documento)
- UI para ver diff entre versiones y restaurar cualquiera
- Maximo safety pero maxima complejidad
- Definitivamente post-MVP

**Decision:** `[ Mezcla de opcion B Y C no el historial cmpleto pero si una N cantidad hacia atras de vesiones, de momento 3 ]`

---

## Dependencias entre preguntas

```
P12 (publicado editable?) ──→ condiciona:
  ├── P7  (granularidad de re-generacion)
  ├── P10 (CRUD manual de clases)
  ├── P15 (persistir chat — mas relevante si published es editable)
  ├── P16 (tools — mas tools si published es editable)
  ├── P19 (concurrencia — mas critico si published es editable)
  └── P20 (undo — mas necesario si published es editable)

P10 (CRUD clases?) ──→ condiciona:
  ├── P16 (tools add_class/remove_class)
  └── P9  (count incorrecto — menos grave si se pueden agregar clases)

P7 (re-generacion?) ──→ condiciona:
  └── P20 (snapshots — mas necesario si re-generar sobreescribe todo)
```

**Orden sugerido para la reunion:**
1. P12 primero (es la decision central)
2. P7 (define como funciona la re-generacion)
3. P10 (define cuanto control tiene el coordinador sobre las clases)
4. El resto en orden

---

## Notas para la reunion

- Las preguntas marcadas **"OBVIA"** tienen una recomendacion clara. Revisar rapido, confirmar o corregir, y avanzar
- **P12 es la pregunta que define el caracter del producto.** Tomense tiempo ahi
- Las opciones no son mutuamente excluyentes con "post-MVP": decidir que es MVP y que es v2 es una respuesta valida
- Si una pregunta genera mucha discusion, anotarla para una segunda sesion en vez de trabar la reunion
