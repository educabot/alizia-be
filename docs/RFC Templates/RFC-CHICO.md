# RFC: [TICKET-ID] Título del cambio

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  PLANTILLA CHICA — Cambios acotados                                     ║
  ║                                                                          ║
  ║  Usar para: agregar un botón, un filtro, un campo, un estado,            ║
  ║  un ajuste de UX, un endpoint simple, un fix de lógica.                  ║
  ║                                                                          ║
  ║  Si el cambio involucra más de 2 flujos o más de 3 archivos por          ║
  ║  disciplina, considerá usar la plantilla MEDIANA.                        ║
  ║                                                                          ║
  ║  INSTRUCCIONES:                                                          ║
  ║  1. Copiá esta plantilla y renombrá el archivo.                          ║
  ║  2. Completá todas las secciones — son pocas, todas son obligatorias.    ║
  ║  3. Borrá los comentarios HTML cuando esté listo para revisión.          ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->

| Campo            | Valor                                      |
|------------------|--------------------------------------------|
| **Ticket**       | [TICKET-ID](link al ticket)                |
| **Autor(es)**    | @nombre                                    |
| **Estado**       | 🟡 Borrador / 🔵 En revisión / 🟢 Aprobado |
| **Creado**       | YYYY-MM-DD                                 |
| **Revisores**    | @revisor1, @revisor2                       |

---

## Qué y por qué 🎯

<!--
  En 2-4 oraciones: qué se quiere hacer y por qué.
  No es un documento de producto — es contexto suficiente para que
  cualquier dev del equipo entienda el cambio sin abrir Jira.
-->

**Qué:** Describir el cambio concreto (ej: "Agregar botón de exportar PDF en la vista de planificación").

**Por qué:** Describir la motivación (ej: "Los docentes necesitan imprimir la planificación para entregarla en la institución").

**Referencia:** [Link al diseño en Figma / ticket / conversación](url)

---

## Alcance 📋

<!--
  Sé explícito sobre lo que entra y lo que no.
  Esto evita que el reviewer pregunte "¿y esto también lo hacemos?"
-->

| Incluye | No incluye |
|---------|------------|
| Cosa 1  | Cosa A — razón breve |
| Cosa 2  | Cosa B — razón breve |

---

## Frontend 🖥️

<!--
  Describí qué cambia en el front. Pensá en:
  - ¿Qué componente se crea o modifica?
  - ¿Qué acción dispara? (click, submit, navegación)
  - ¿Qué estado se afecta? (Zustand store, React Query cache)
  - ¿Hay un loading state, error state, empty state?
-->

### Componente(s) afectado(s)

| Componente | Acción | Archivo aproximado |
|------------|--------|--------------------|
| Nombre del componente | Crear / Modificar | `src/components/...` |

### Comportamiento esperado

<!--
  Describí el flujo desde la perspectiva del usuario:
  1. El usuario hace X
  2. Se muestra Y
  3. Si hay error, se muestra Z
-->

1. El usuario hace clic en [elemento]
2. Se dispara [acción]
3. El sistema muestra [resultado]

### Estados de UI

| Estado | Qué se muestra |
|--------|----------------|
| Default | Descripción |
| Loading | Descripción |
| Error | Descripción |
| Vacío | Descripción (si aplica) |

---

## Backend ⚙️

<!--
  Describí qué cambia en el back. Pensá en:
  - ¿Se crea o modifica un endpoint?
  - ¿Cambia el modelo de datos? (nueva columna, nuevo campo en response)
  - ¿Hay lógica de negocio nueva?
  Si el cambio es solo de front, poné "Sin cambios en backend" y listo.
-->

### Endpoint(s)

<!--
  Si no hay cambio de endpoint, borrá esta sub-sección.
-->

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST/PUT/DELETE` | `/api/v1/ruta` | Qué hace |

**Request:**

```json
{
  "campo": "valor"
}
```

**Response:**

| Status | Descripción | Payload |
|--------|-------------|---------|
| 200 | OK | `{ "id": "..." }` |
| 400 | Validación | `{ "error": "..." }` |

### Modelo de datos

<!--
  Si no hay cambio de schema, borrá esta sub-sección.
-->

```sql
-- Describí el cambio: nueva columna, nuevo índice, etc.
ALTER TABLE tabla ADD COLUMN nuevo_campo TYPE DEFAULT valor;
```

### Lógica de negocio

<!--
  Describí brevemente la lógica nueva o modificada.
  No es necesario poner código — alcanza con describir las reglas.
-->

- Regla 1: Si X entonces Y
- Regla 2: Validar que Z antes de guardar

---

## Recomendaciones para UX 🎨

<!--
  Esta sección NO la completa UX — la completa el owner del RFC
  para pedir lo que necesita del equipo de diseño.

  Si el cambio es puramente técnico sin impacto visual, poné
  "Sin impacto en UX" y borrá el resto.
-->

### ¿Se necesita diseño?

<!-- Sí / No / Ya existe (link a Figma) -->

### Qué pedirle a UX

<!--
  Sé específico. En vez de "diseñar el botón", pedí:
  - "Ubicación del botón en la vista de planificación"
  - "Estado del botón cuando la exportación está en progreso"
  - "Qué feedback visual recibe el usuario al completar"
-->

- [ ] Pedido 1: Descripción concreta de lo que necesitás
- [ ] Pedido 2: Descripción concreta
- [ ] Pedido 3: Descripción concreta

### Contexto para UX

<!--
  Info que le sirve al diseñador para tomar decisiones:
  - ¿Dónde vive este componente? (qué pantalla, qué sección)
  - ¿Quién lo usa? (coordinador, docente, ambos)
  - ¿Hay restricciones técnicas? (ej: "no podemos hacer drag & drop acá")
-->

- Pantalla: [nombre de la vista]
- Rol(es): [coordinador / docente / ambos]
- Restricciones: [si las hay]

---

## QA — Flujos a testear 🧪

<!--
  Listá los flujos que QA debe verificar manualmente.
  Esto se convierte en la planilla de testing del cambio.
  Cada fila es un caso de prueba.
-->

| # | Flujo | Pasos | Resultado esperado | Estado |
|---|-------|-------|--------------------|--------|
| 1 | Happy path | 1. Hacer X → 2. Verificar Y | Se muestra Z correctamente | ⬜ Pendiente |
| 2 | Error case | 1. Hacer X sin completar Y | Se muestra mensaje de error | ⬜ Pendiente |
| 3 | Edge case | 1. Hacer X con datos límite | Se comporta según regla Z | ⬜ Pendiente |

### Precondiciones

<!--
  ¿Qué necesita QA para poder testear esto?
  Ej: "Tener un usuario docente con al menos 1 curso asignado"
-->

- Precondición 1
- Precondición 2

---

## Preguntas abiertas ❓

<!--
  Si no hay preguntas, borrá la sección.
-->

| # | Pregunta | Responsable | Estado |
|---|----------|-------------|--------|
| 1 | ¿Pregunta? | @nombre | 🟡 Pendiente |

---

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  Checklist antes de enviar a revisión:                                  ║
  ║  □ ¿Completé "Qué y por qué"?                                          ║
  ║  □ ¿Definí alcance (incluye / no incluye)?                              ║
  ║  □ ¿Describí el cambio de front O puse "Sin cambios"?                   ║
  ║  □ ¿Describí el cambio de back O puse "Sin cambios"?                    ║
  ║  □ ¿Pedí lo que necesito de UX O puse "Sin impacto"?                    ║
  ║  □ ¿Listé al menos el happy path en QA?                                 ║
  ║  □ ¿Borré los comentarios HTML?                                         ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->
