# RFC: [TICKET-ID] Título de la funcionalidad

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  PLANTILLA MEDIANA — Funcionalidades completas                          ║
  ║                                                                          ║
  ║  Usar para: una feature completa (bitácora de cotejo, exportación PDF,   ║
  ║  flujo de estados de un documento, sistema de notificaciones, etc.)      ║
  ║                                                                          ║
  ║  Si el cambio toca más de 3 épicas, requiere migración de datos          ║
  ║  compleja, o involucra un rediseño de arquitectura, usá la plantilla     ║
  ║  GRANDE.                                                                 ║
  ║                                                                          ║
  ║  INSTRUCCIONES:                                                          ║
  ║  1. Copiá esta plantilla y renombrá el archivo.                          ║
  ║  2. Completá las secciones. Las marcadas [OPCIONAL] se pueden borrar.    ║
  ║  3. Borrá los comentarios HTML cuando esté listo para revisión.          ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->

| Campo            | Valor                                      |
|------------------|--------------------------------------------|
| **Ticket**       | [TICKET-ID](link al ticket)                |
| **Autor(es)**    | @nombre1, @nombre2                         |
| **Estado**       | 🟡 Borrador / 🔵 En revisión / 🟢 Aprobado / 🔴 Rechazado |
| **Creado**       | YYYY-MM-DD                                 |
| **Última edición** | YYYY-MM-DD                               |
| **Revisores**    | @revisor1, @revisor2, @revisor3            |

---

## Historial de versiones

| Versión | Fecha      | Autor   | Cambios                        |
|---------|------------|---------|--------------------------------|
| 0.1     | YYYY-MM-DD | @nombre | Borrador inicial               |

---

## Contexto y motivación 🚀

<!--
  ¿De dónde sale esta necesidad? ¿Qué problema tiene el usuario hoy?
  Incluí links a tickets, conversaciones o documentos relacionados.
  3-5 oraciones, no más.
-->

**Problema:** Describir el dolor del usuario o la limitación actual.

**Contexto:** Describir la situación actual y por qué ahora es el momento de resolverlo.

**Documentos relacionados:**
- [Link a épica o RFC previo](url)
- [Link a diseño en Figma](url)

---

## Objetivo 🎯

<!--
  ¿Qué se logra si esto se implementa bien?
  ¿Qué explícitamente NO es objetivo de este RFC?
-->

### Objetivos

- ✅ Objetivo 1
- ✅ Objetivo 2
- ✅ Objetivo 3

### No-objetivos

- ❌ No-objetivo 1 — razón breve
- ❌ No-objetivo 2 — razón breve

---

## Alcance 📋

### Incluye

- Funcionalidad A
- Funcionalidad B
- Funcionalidad C

### No incluye

- Funcionalidad X — NTH / horizonte
- Funcionalidad Y — depende de [otra épica]

---

## Diseño de producto 🧠

<!--
  Describí el flujo completo desde la perspectiva del usuario.
  Esto es lo que alinea a front, back, UX y QA antes de implementar.
  Usá escenarios numerados.
-->

### Flujo principal

1. El usuario [acción inicial]
2. El sistema [respuesta]
3. El usuario [siguiente acción]
4. El sistema [resultado final]

### Flujos alternativos

#### Si [condición A]

1. ...
2. ...

#### Si [condición B]

1. ...
2. ...

### Reglas de negocio

<!--
  Las reglas que todo el equipo necesita conocer, independientemente
  de si trabajan en front o back.
-->

| # | Regla | Ejemplo |
|---|-------|---------|
| 1 | Descripción de la regla | "Si el docente tiene 2 materias, ve ambas en la misma vista" |
| 2 | Descripción de la regla | "El coordinador puede sobreescribir el cálculo de clases" |

### Decisiones por provincia

<!--
  Si hay comportamiento que varía por cliente/provincia, documentalo acá.
  Esto impacta feature flags y configuración.
-->

| Decisión | Quién decide | Default | Notas |
|----------|-------------|---------|-------|
| Decisión 1 | Provincia | valor | Notas |
| Decisión 2 | Provincia | valor | Notas |

---

## Frontend 🖥️

<!--
  Describí qué se construye o modifica en el front.
  Pensá en: componentes, stores, queries, rutas, estados.
-->

### Vista general

<!--
  Describí la pantalla o sección nueva/modificada a alto nivel.
  Si hay diseño en Figma, linkealo acá.
-->

**Figma:** [Link al diseño](url) (si existe)

### Componentes

| Componente | Acción | Descripción | Archivo |
|------------|--------|-------------|---------|
| NombreComponente | Crear | Qué hace | `src/components/...` |
| OtroComponente | Modificar | Qué cambia | `src/components/...` |

### Estado (Zustand / React Query)

<!--
  ¿Se crea o modifica un store?
  ¿Hay nueva query/mutation con TanStack Query?
  ¿Hay cache invalidation a considerar?
-->

| Store / Query | Acción | Descripción |
|---------------|--------|-------------|
| `useNombreStore` | Crear / Modificar | Qué estado maneja |
| `useNombreQuery` | Crear | Qué datos fetchea, cuándo se invalida |
| `useNombreMutation` | Crear | Qué endpoint llama, qué invalida en cache |

### Rutas

<!--
  Si se crea una ruta nueva o se modifica una existente.
  Si no hay cambio de rutas, borrá esta sub-sección.
-->

| Ruta | Componente | Descripción |
|------|-----------|-------------|
| `/nueva-ruta` | PaginaNueva | Qué muestra |

### Estados de UI por pantalla

| Pantalla / Componente | Default | Loading | Error | Vacío |
|-----------------------|---------|---------|-------|-------|
| NombreComponente | Descripción | Descripción | Descripción | Descripción |

### Consideraciones de front

<!--
  Cosas que el dev de front necesita saber:
  - ¿Hay responsive/mobile a considerar?
  - ¿Hay accesibilidad especial? (aria labels, keyboard nav)
  - ¿Se usa un componente de Radix UI existente?
  - ¿Hay animaciones o transiciones?
-->

- Consideración 1
- Consideración 2

---

## Backend ⚙️

<!--
  Describí qué se construye o modifica en el back.
  Pensá en: endpoints, modelos, lógica, queries, jobs.
-->

### Endpoints

#### `METHOD /api/v1/ruta`

**Descripción:** Qué hace este endpoint.

**Request:**

```json
{
  "campo_1": "valor",
  "campo_2": 123
}
```

**Parámetros:**

| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| campo_1 | `string` | Sí | Descripción |
| campo_2 | `int` | No | Descripción |

**Response:**

| Status | Causa | Payload |
|--------|-------|---------|
| 200 | OK | `{ "data": [...] }` |
| 400 | Validación fallida | `{ "error": "..." }` |
| 404 | No encontrado | `{ "error": "..." }` |

#### `METHOD /api/v1/otra-ruta`

<!-- Repetir estructura por cada endpoint -->

---

### Modelo de datos

<!--
  Nuevas tablas, columnas, índices, o cambios al schema.
-->

```sql
-- Describí la migración
CREATE TABLE nueva_tabla (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    campo_1     VARCHAR(255) NOT NULL,
    campo_2     INTEGER     DEFAULT 0,
    org_id      UUID        NOT NULL REFERENCES organizations(id),
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nueva_tabla_org ON nueva_tabla (org_id);
```

**Migraciones:**

| Orden | Migración | Reversible | Requiere downtime |
|-------|-----------|------------|-------------------|
| 1 | Crear tabla X | Sí | No |
| 2 | Agregar columna Y | Sí | No |

### Lógica de negocio

<!--
  Describí las reglas que el backend tiene que implementar.
  No es necesario poner código completo — alcanza con pseudo-lógica.
-->

- **Regla 1:** Si [condición] entonces [acción]. Motivo: [por qué].
- **Regla 2:** Validar [qué] antes de [operación]. Motivo: [por qué].

### Configuración por organización [OPCIONAL]

<!--
  Si el feature usa el JSON de configuración por org.
-->

```json
{
  "feature_nombre": {
    "enabled": true,
    "config": {
      "param_1": "valor",
      "param_2": 100
    }
  }
}
```

### Consideraciones de back

<!--
  Cosas que el dev de back necesita saber:
  - ¿Hay permisos/roles involucrados?
  - ¿Hay performance concerns? (N+1, queries pesadas)
  - ¿Se necesita cache (Redis)?
  - ¿Hay jobs asincrónicos? (Cloud Tasks)
-->

- Consideración 1
- Consideración 2

---

## Recomendaciones para UX 🎨

<!--
  Qué necesitás del equipo de diseño para esta feature.
  Si ya existe el diseño completo en Figma, linkeá y describí
  solo lo que falta o lo que necesita ajuste.
-->

### Estado del diseño

<!-- Elegí uno:
  - ✅ Diseño completo en Figma: [link](url)
  - 🟡 Diseño parcial — falta: [lista]
  - ❌ Sin diseño — se necesita antes de implementar
  - ⚪ No requiere diseño (cambio invisible para el usuario)
-->

### Qué pedirle a UX

- [ ] Pedido 1: Descripción concreta
- [ ] Pedido 2: Descripción concreta
- [ ] Pedido 3: Descripción concreta

### Contexto para UX

| Aspecto | Detalle |
|---------|---------|
| Pantalla(s) | Nombre de la vista donde vive |
| Rol(es) | Coordinador / Docente / Ambos |
| Dispositivo | Desktop / Mobile / Ambos |
| Componentes existentes | ¿Hay algo en el design system que se puede reutilizar? |
| Restricciones técnicas | Limitaciones que UX debe conocer |
| Referencia | Links a features similares en el producto o en competidores |

---

## QA — Plan de testing 🧪

<!--
  Esta sección es el insumo para que QA arme su planilla de testing.
  Cubrir: happy path, errores, edge cases, y permisos/roles.
-->

### Precondiciones

- Precondición 1 (ej: "Tener usuario coordinador con al menos 1 documento publicado")
- Precondición 2

### Casos de prueba

| # | Categoría | Flujo | Pasos | Resultado esperado | Rol | Estado |
|---|-----------|-------|-------|--------------------|-----|--------|
| 1 | Happy path | Flujo principal completo | 1. X → 2. Y → 3. Z | Se muestra resultado correcto | Docente | ⬜ |
| 2 | Happy path | Variante con dato A | 1. X → 2. Y | Se comporta según regla A | Coordinador | ⬜ |
| 3 | Error | Campo inválido | 1. Dejar campo vacío → 2. Submit | Mensaje de error claro | Docente | ⬜ |
| 4 | Error | Sin permisos | 1. Intentar acción sin rol | Se bloquea con mensaje | Docente | ⬜ |
| 5 | Edge case | Sin datos previos | 1. Entrar sin datos cargados | Empty state correcto | Ambos | ⬜ |
| 6 | Edge case | Datos límite | 1. Cargar máximo permitido | Se comporta correctamente | Docente | ⬜ |
| 7 | Permisos | Coordinador ve X | 1. Entrar como coordinador | Ve opciones de coordinador | Coordinador | ⬜ |
| 8 | Permisos | Docente no ve Y | 1. Entrar como docente | No ve opciones de coordinador | Docente | ⬜ |

### Regresión

<!--
  ¿Qué funcionalidad existente podría romperse con este cambio?
-->

| # | Funcionalidad existente | Qué verificar | Estado |
|---|------------------------|---------------|--------|
| 1 | Feature X | Que siga funcionando Y | ⬜ |

---

## Alternativas evaluadas 🔀

<!--
  Al menos 2 alternativas. Si no hay alternativas reales,
  explicá por qué la solución propuesta es la única viable.
-->

| Alternativa | Descripción | Ventajas | Desventajas |
|-------------|-------------|----------|-------------|
| A: [nombre] | Qué propone | ✅ Pro 1, ✅ Pro 2 | ❌ Con 1 |
| B: [nombre] | Qué propone | ✅ Pro 1 | ❌ Con 1, ❌ Con 2 |
| C: [nombre] (elegida) | Qué propone | ✅ Pro 1, ✅ Pro 2 | ❌ Con 1 |

**Opción elegida: C** — Justificación breve del por qué.

---

## Dependencias 👥

| Dependencia | Tipo | Bloqueante | Estado | Responsable |
|-------------|------|------------|--------|-------------|
| Épica / RFC X | Interna | Sí / No | ✅ / 🟡 / 🔴 | @nombre |
| API externa Y | Externa | Sí / No | ✅ / 🟡 / 🔴 | @nombre |

---

## Preguntas abiertas ❓

| # | Pregunta | Área | Responsable | Estado |
|---|----------|------|-------------|--------|
| 1 | ¿Pregunta? | Producto / Front / Back / UX | @nombre | 🟡 Pendiente |
| 2 | ¿Pregunta? | Producto / Front / Back / UX | @nombre | 💬 Respondida: respuesta |

---

## Tareas 📝

| # | Tarea | Disciplina | Sección | Asignado | Estado | Ticket |
|---|-------|-----------|---------|----------|--------|--------|
| 1 | Diseño en Figma | UX | Recomendaciones UX | @nombre | ⬜ | — |
| 2 | Crear migración SQL | Back | Modelo de datos | @nombre | ⬜ | — |
| 3 | Implementar endpoint X | Back | Endpoints | @nombre | ⬜ | — |
| 4 | Crear componente Y | Front | Componentes | @nombre | ⬜ | — |
| 5 | Integrar con endpoint | Front | Estado | @nombre | ⬜ | — |
| 6 | Testing manual | QA | Plan de testing | @nombre | ⬜ | — |

---

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  Checklist antes de enviar a revisión:                                  ║
  ║  □ ¿Contexto y motivación claros?                                       ║
  ║  □ ¿Objetivos y no-objetivos definidos?                                 ║
  ║  □ ¿Alcance con incluye / no incluye?                                   ║
  ║  □ ¿Flujo de producto descrito paso a paso?                             ║
  ║  □ ¿Sección de Front completa o marcada "Sin cambios"?                  ║
  ║  □ ¿Sección de Back completa o marcada "Sin cambios"?                   ║
  ║  □ ¿Pedidos a UX concretos o marcado "No requiere diseño"?              ║
  ║  □ ¿Plan de QA con al menos happy path + error + edge case?             ║
  ║  □ ¿Al menos 2 alternativas evaluadas?                                  ║
  ║  □ ¿Preguntas abiertas listadas?                                        ║
  ║  □ ¿Tareas con disciplina asignada?                                     ║
  ║  □ ¿Borré los comentarios HTML?                                         ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->
