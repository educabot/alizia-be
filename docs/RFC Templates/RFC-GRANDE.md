# RFC: [TICKET-ID] Título del RFC

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  PLANTILLA GRANDE — Épicas, APIs completas, refactors cross-equipo      ║
  ║                                                                          ║
  ║  Usar para: una épica completa, una API con múltiples flujos,            ║
  ║  un refactor de arquitectura, un cambio que toca 3+ disciplinas,         ║
  ║  o cualquier cosa que requiera más de 2 semanas de desarrollo.           ║
  ║                                                                          ║
  ║  INSTRUCCIONES:                                                          ║
  ║  1. Copiá esta plantilla y renombrá el archivo.                          ║
  ║  2. Completá las secciones relevantes. Las marcadas [OPCIONAL]           ║
  ║     se pueden borrar si no aplican.                                      ║
  ║  3. Las secciones sin marca son OBLIGATORIAS.                            ║
  ║  4. Borrá los comentarios HTML cuando esté listo para revisión.          ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->

| Campo              | Valor                                      |
|--------------------|--------------------------------------------|
| **Ticket**         | [TICKET-ID](link al ticket)                |
| **Autor(es)**      | @nombre1, @nombre2                         |
| **Estado**         | 🟡 Borrador / 🔵 En revisión / 🟢 Aprobado / 🔴 Rechazado / ⚪ Deprecado |
| **Tipo**           | Épica / API / Refactor                     |
| **Creado**         | YYYY-MM-DD                                 |
| **Última edición** | YYYY-MM-DD                                 |
| **Revisores**      | @revisor1, @revisor2, @revisor3            |
| **Decisión**       | Pendiente / Aprobada con opción X          |

---

## Historial de versiones

| Versión | Fecha      | Autor   | Cambios                          |
|---------|------------|---------|----------------------------------|
| 0.1     | YYYY-MM-DD | @nombre | Borrador inicial                 |
| 0.2     | YYYY-MM-DD | @nombre | Incorpora feedback de revisión   |
| 1.0     | YYYY-MM-DD | @nombre | Versión aprobada                 |

---

## Índice

<!--
  En un RFC grande el índice es obligatorio.
  Actualizalo si agregás o sacás secciones.
-->

- [Contexto y motivación](#contexto-y-motivación-)
- [Objetivo](#objetivo-)
- [Alcance](#alcance-)
- [Diseño de producto](#diseño-de-producto-)
- [Arquitectura general](#arquitectura-general-)
- [Frontend](#frontend-️)
- [Backend — Endpoints](#backend--endpoints-️)
- [Backend — Modelo de datos](#backend--modelo-de-datos-)
- [Backend — Lógica y configuración](#backend--lógica-y-configuración-)
- [Recomendaciones para UX](#recomendaciones-para-ux-)
- [QA — Estrategia de testing](#qa--estrategia-de-testing-)
- [Alternativas evaluadas](#alternativas-evaluadas-)
- [Rollout](#rollout-)
- [Dependencias](#dependencias-)
- [Riesgos](#riesgos-)
- [Preguntas abiertas](#preguntas-abiertas-)
- [Glosario](#glosario-)
- [Tareas](#tareas-)

---

## Contexto y motivación 🚀

<!--
  ¿Por qué existe este RFC? ¿Qué problema resuelve y para quién?
  Incluí el contexto de negocio que un nuevo integrante del equipo
  necesitaría para entender la motivación.
-->

### Problema

Describir el problema desde la perspectiva del usuario...

### Contexto

Describir la situación actual, intentos previos, y por qué ahora...

### Documentos relacionados

- [Épica / RFC previo](url)
- [Diseño en Figma](url)
- [Ticket principal](url)
- [Conversación de referencia](url)

---

## Objetivo 🎯

### Objetivos

- ✅ Objetivo 1
- ✅ Objetivo 2
- ✅ Objetivo 3

### No-objetivos

- ❌ No-objetivo 1 — razón breve
- ❌ No-objetivo 2 — razón breve

### Métricas de éxito

<!--
  ¿Cómo sabemos que esto funcionó?
  Definir métricas concretas, aunque sean cualitativas en MVP.
-->

| Métrica | Valor esperado | Cómo se mide |
|---------|---------------|--------------|
| Métrica 1 (ej: % de docentes que planifican) | X% | Dashboard / query / manual |
| Métrica 2 (ej: documentos publicados por mes) | N | Dashboard / query |
| Métrica 3 (ej: tiempo promedio de creación) | < X min | Analytics |

---

## Alcance 📋

### Incluye

- Funcionalidad A
- Funcionalidad B
- Funcionalidad C
- Funcionalidad D

### No incluye

- Funcionalidad X — NTH / horizonte
- Funcionalidad Y — depende de [otro RFC]
- Funcionalidad Z — fuera de scope, se aborda en [épica]

### Fases de implementación

<!--
  Si la feature es grande, dividirla en fases ayuda a entregar
  valor incremental y a que QA pueda testear progresivamente.
-->

| Fase | Qué incluye | Entrega estimada | Dependencia |
|------|-------------|------------------|-------------|
| 1 | Funcionalidad core (A + B) | Sprint X | Ninguna |
| 2 | Funcionalidad complementaria (C) | Sprint Y | Fase 1 |
| 3 | Polish y edge cases (D) | Sprint Z | Fase 2 |

---

## Diseño de producto 🧠

<!--
  Esta sección alinea a TODO el equipo en qué se construye.
  Es la fuente de verdad de producto para este RFC.
-->

### Principios de diseño

<!--
  Reglas que guían las decisiones. Útil cuando hay trade-offs
  y el equipo necesita criterio para decidir sin escalar.
-->

1. **Principio 1** — Descripción breve
2. **Principio 2** — Descripción breve
3. **Principio 3** — Descripción breve

### Flujos de usuario

#### Flujo 1: [Nombre del flujo]

**Actor:** Coordinador / Docente / Ambos
**Precondición:** Qué debe existir antes de que el flujo sea posible.

1. El usuario [acción]
2. El sistema [respuesta]
3. El usuario [acción]
4. El sistema [resultado]

**Resultado exitoso:** Descripción del estado final esperado.

#### Flujo 2: [Nombre del flujo]

**Actor:** ...
**Precondición:** ...

1. ...

#### Flujo alternativo: [Nombre]

1. ...

#### Flujo de error: [Nombre]

1. ...

### Reglas de negocio

| # | Regla | Ejemplo | Aplica a |
|---|-------|---------|----------|
| 1 | Descripción | Ejemplo concreto | Front / Back / Ambos |
| 2 | Descripción | Ejemplo concreto | Back |
| 3 | Descripción | Ejemplo concreto | Front |

### Decisiones por provincia

<!--
  Comportamiento que varía por cliente. Esto impacta directamente
  el JSON de configuración por org y los feature flags.
-->

| Decisión | Quién decide | Default | Impacto en config |
|----------|-------------|---------|-------------------|
| Decisión 1 | Provincia | Valor | Feature flag / JSON config |
| Decisión 2 | Provincia | Valor | Feature flag / JSON config |
| Decisión 3 | Institución | Valor | JSON config |

### Estados y ciclo de vida

<!--
  Si la entidad principal tiene estados (borrador, publicado, etc.).
  Describir transiciones y quién puede ejecutar cada una.
-->

```
[Estado A] ──(acción / rol)──▶ [Estado B] ──(acción / rol)──▶ [Estado C]
                                    │
                                    ▼
                               [Estado D]
```

| Transición | De | A | Quién puede | Condiciones |
|------------|-----|---|-------------|-------------|
| Publicar | Borrador | Publicado | Coordinador | Todas las secciones completas |
| Revertir | Publicado | Borrador | Coordinador | — |

---

## Arquitectura general 🏗️

<!--
  Vista de alto nivel de cómo interactúan los componentes.
  Útil para que todo el equipo entienda el big picture.
-->

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend   │────▶│   API (Go)   │────▶│  PostgreSQL  │
│  React + TS  │◀────│  Gin + GORM  │◀────│              │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │  Cloud Tasks │
                    │  (async)     │
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │  LLM / IA    │
                    │  (Azure)     │
                    └─────────────┘
```

<!--
  Adaptá el diagrama al caso real. Si no hay async o IA, sacá esos boxes.
  Si hay Redis, agregalo. El punto es que sea un mapa rápido.
-->

### Componentes involucrados

| Componente | Responsabilidad | Cambio |
|------------|----------------|--------|
| `tuni-ai-webapp` | UI de la feature | Crear / Modificar |
| `tich-cronos` | API + lógica de negocio | Crear / Modificar |
| Cloud Tasks | Jobs asíncronos | Crear (si aplica) |
| PostgreSQL | Persistencia | Nueva tabla / migración |

---

## Frontend 🖥️

### Vista general

<!--
  Describí las pantallas nuevas o modificadas a alto nivel.
  Si hay diseño en Figma, este es el link principal.
-->

**Figma:** [Link al diseño completo](url)

### Mapa de pantallas

<!--
  Listá las pantallas/vistas involucradas y su relación.
-->

| Pantalla | Ruta | Descripción | Estado |
|----------|------|-------------|--------|
| Vista A | `/ruta-a` | Qué muestra | Nueva |
| Vista B | `/ruta-b` | Qué muestra | Nueva |
| Vista C | `/ruta-c` | Qué muestra | Modificada |

### Componentes

| Componente | Tipo | Descripción | Reutiliza |
|------------|------|-------------|-----------|
| ComponenteA | Page | Página principal del flujo | — |
| ComponenteB | Feature | Sección específica dentro de la página | — |
| ComponenteC | UI | Componente reutilizable | Radix Dialog |

### Estado y data fetching

<!--
  Detallar stores de Zustand y queries/mutations de TanStack Query.
  Esto es crítico para que el dev de front entienda el flujo de datos.
-->

#### Stores (Zustand)

| Store | Propósito | Estado que maneja |
|-------|-----------|-------------------|
| `useNombreStore` | Descripción | `{ campo1, campo2, acciones }` |

#### Queries (TanStack Query)

| Hook | Endpoint | Query key | Stale time | Descripción |
|------|----------|-----------|------------|-------------|
| `useNombreQuery` | `GET /api/v1/...` | `['nombre', id]` | 5min | Qué datos trae |

#### Mutations (TanStack Query)

| Hook | Endpoint | Invalidates | Descripción |
|------|----------|-------------|-------------|
| `useNombreMutation` | `POST /api/v1/...` | `['nombre']` | Qué hace, qué invalida |

### Estados de UI

| Pantalla | Default | Loading | Error | Vacío |
|----------|---------|---------|-------|-------|
| Vista A | Descripción | Skeleton / Spinner | Toast + retry | Empty state con CTA |
| Vista B | Descripción | Descripción | Descripción | Descripción |

### Consideraciones de front

- **Responsive:** ¿Hay comportamiento mobile?
- **Accesibilidad:** ¿Keyboard navigation, aria labels, screen reader?
- **Performance:** ¿Listas largas que requieren virtualización? ¿Lazy loading?
- **Optimistic updates:** ¿Alguna mutation se puede aplicar optimisticamente?
- **Feature flags:** ¿Qué se esconde detrás de un flag de LaunchDarkly?

---

## Backend — Endpoints ⚙️

<!--
  Un bloque por cada endpoint. Repetir la estructura.
-->

### `METHOD /api/v1/ruta-1`

| Campo | Valor |
|-------|-------|
| **Descripción** | Qué hace |
| **Auth** | Sí — Roles: coordinador / docente |
| **Rate limit** | Si aplica |

**Request:**

```json
{
  "campo_1": "valor",
  "campo_2": 123
}
```

**Parámetros:**

| Parámetro | Tipo | Requerido | Validación | Descripción |
|-----------|------|-----------|------------|-------------|
| campo_1 | `string` | Sí | max 255 chars | Descripción |
| campo_2 | `int` | No | > 0 | Descripción |

**Response 200:**

```json
{
  "data": {
    "id": "uuid",
    "campo_1": "valor",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Errores:**

| Status | Código | Causa | Payload |
|--------|--------|-------|---------|
| 400 | `VALIDATION_ERROR` | Campo inválido | `{ "error": "...", "field": "..." }` |
| 403 | `FORBIDDEN` | Sin permisos | `{ "error": "..." }` |
| 404 | `NOT_FOUND` | Recurso no existe | `{ "error": "..." }` |
| 409 | `CONFLICT` | Estado inválido para la operación | `{ "error": "..." }` |

**SQL principal:**

```sql
-- Query principal del endpoint
SELECT campo_1, campo_2
FROM tabla
WHERE org_id = $1
  AND activo = true
ORDER BY created_at DESC;
```

---

### `METHOD /api/v1/ruta-2`

<!-- Repetir estructura -->

---

## Backend — Modelo de datos ⛁

### Diagrama de entidades

<!--
  Relaciones entre las tablas principales de este RFC.
-->

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  organizations   │────▶│  entidad_nueva   │────▶│  entidad_hija    │
│  (existente)     │     │  (nueva)         │     │  (nueva)         │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

### Nuevas tablas

```sql
CREATE TABLE entidad_nueva (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID        NOT NULL REFERENCES organizations(id),
    titulo          VARCHAR(255) NOT NULL,
    estado          VARCHAR(50) NOT NULL DEFAULT 'borrador',
    config          JSONB       DEFAULT '{}',
    created_by      UUID        NOT NULL REFERENCES users(id),
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_estado CHECK (estado IN ('borrador', 'en_revision', 'publicado'))
);

-- Índices
CREATE INDEX idx_entidad_nueva_org ON entidad_nueva (org_id);
CREATE INDEX idx_entidad_nueva_estado ON entidad_nueva (org_id, estado);
```

### Modificaciones a tablas existentes

```sql
ALTER TABLE tabla_existente
    ADD COLUMN nuevo_campo VARCHAR(100) DEFAULT NULL;
```

### Migraciones

| Orden | Migración | Reversible | Downtime | Notas |
|-------|-----------|------------|----------|-------|
| 1 | Crear tabla entidad_nueva | Sí | No | — |
| 2 | Crear tabla entidad_hija | Sí | No | FK a entidad_nueva |
| 3 | Agregar columna a tabla X | Sí | No | Nullable, sin default costoso |

---

## Backend — Lógica y configuración 🔧

### Lógica de negocio

<!--
  Reglas que el backend implementa. Agrupá por dominio.
-->

#### [Dominio 1: ej. Gestión de estados]

- **Regla 1:** Descripción. **Motivo:** por qué.
- **Regla 2:** Descripción. **Motivo:** por qué.

#### [Dominio 2: ej. Cálculos]

- **Regla 3:** Descripción. **Motivo:** por qué.

### Configuración por organización

<!--
  Si esta feature usa el JSON de configuración por org.
-->

```json
{
  "nombre_feature": {
    "enabled": true,
    "nombres_niveles": {
      "1": "Conocimientos y saberes",
      "2": "Núcleos problemáticos"
    },
    "profundidad_maxima": 3,
    "opciones": ["opcion_a", "opcion_b"]
  }
}
```

| Campo | Tipo | Default | Descripción |
|-------|------|---------|-------------|
| `enabled` | `bool` | `false` | Habilita la feature |
| `nombres_niveles` | `map[int]string` | `{}` | Nombres por nivel de jerarquía |
| `profundidad_maxima` | `int` | `2` | Máxima profundidad del árbol |

### Feature flags (LaunchDarkly)

| Flag | Tipo | Default | Descripción | Quién controla |
|------|------|---------|-------------|----------------|
| `ff_nombre_feature` | `bool` | `false` | Habilita la feature completa | Producto |
| `ff_nombre_sub_feature` | `bool` | `false` | Habilita sub-funcionalidad | Producto |

### Jobs asíncronos [OPCIONAL]

<!--
  Si hay Cloud Tasks, schedulers, o procesos asincrónicos.
-->

| Job | Trigger | Frecuencia | Descripción | Timeout |
|-----|---------|------------|-------------|---------|
| `job_nombre` | Cloud Tasks / Scheduler | On demand / Cron | Qué hace | 60s |

### Consideraciones de back

- **Permisos:** ¿Qué roles acceden a qué endpoints?
- **Performance:** ¿Queries que podrían ser costosas? ¿N+1?
- **Cache:** ¿Se usa Redis para algo?
- **Concurrencia:** ¿Puede haber escrituras simultáneas?
- **Idempotencia:** ¿Los endpoints son idempotentes?

---

## Recomendaciones para UX 🎨

### Estado del diseño

<!--
  ✅ Diseño completo en Figma: [link](url)
  🟡 Diseño parcial — falta: [lista]
  ❌ Sin diseño — se necesita antes de implementar
-->

### Qué pedirle a UX — por pantalla

<!--
  Organizá los pedidos por pantalla/flujo para que UX
  pueda trabajar de forma estructurada.
-->

#### Pantalla 1: [nombre]

- [ ] Pedido concreto
- [ ] Pedido concreto

#### Pantalla 2: [nombre]

- [ ] Pedido concreto
- [ ] Pedido concreto

#### Transversales

- [ ] Empty states para todas las vistas nuevas
- [ ] Estados de loading (skeleton vs spinner)
- [ ] Mensajes de error y confirmación
- [ ] Responsive / mobile (si aplica)

### Contexto para UX

| Aspecto | Detalle |
|---------|---------|
| Pantalla(s) | Lista de vistas involucradas |
| Rol(es) | Coordinador / Docente / Ambos |
| Dispositivo | Desktop / Mobile / Ambos |
| Componentes reutilizables | ¿Qué del design system se puede usar? |
| Restricciones técnicas | Limitaciones que UX debe conocer |
| Flujos de referencia | Features similares ya diseñadas o de competidores |
| Volumen de datos | ¿Cuántos items se esperan en listas/tablas? |

---

## QA — Estrategia de testing 🧪

### Precondiciones generales

<!--
  Qué necesita QA configurado antes de empezar a testear.
-->

- Precondición 1 (ej: "Tener organización con configuración de provincia X cargada")
- Precondición 2 (ej: "Tener usuario coordinador y docente en la misma org")
- Precondición 3 (ej: "Feature flag habilitado en staging")

### Matriz de testing por flujo

#### Flujo 1: [nombre]

| # | Caso | Pasos | Resultado esperado | Rol | Prioridad | Estado |
|---|------|-------|--------------------|-----|-----------|--------|
| 1.1 | Happy path | 1. X → 2. Y → 3. Z | Resultado correcto | Coordinador | Alta | ⬜ |
| 1.2 | Variante A | 1. X con dato A | Se comporta según regla A | Coordinador | Alta | ⬜ |
| 1.3 | Error: campo inválido | 1. Dejar campo vacío → 2. Submit | Error claro | Coordinador | Media | ⬜ |
| 1.4 | Edge: sin datos previos | 1. Entrar sin datos | Empty state | Coordinador | Media | ⬜ |

#### Flujo 2: [nombre]

| # | Caso | Pasos | Resultado esperado | Rol | Prioridad | Estado |
|---|------|-------|--------------------|-----|-----------|--------|
| 2.1 | Happy path | ... | ... | Docente | Alta | ⬜ |
| 2.2 | ... | ... | ... | ... | ... | ⬜ |

### Testing de permisos y roles

| # | Escenario | Rol | Acción | Resultado esperado | Estado |
|---|-----------|-----|--------|--------------------|--------|
| P1 | Coordinador accede a X | Coordinador | Navegar a /ruta | Ve la pantalla completa | ⬜ |
| P2 | Docente no accede a Y | Docente | Navegar a /ruta | Redirect o 403 | ⬜ |
| P3 | Multi-rol | Coordinador + Docente | Navegar | Ve ambas funcionalidades | ⬜ |

### Testing cross-browser / dispositivo [OPCIONAL]

| Browser / Device | Resolución | Flujos a verificar | Estado |
|-----------------|------------|-------------------|--------|
| Chrome Desktop | 1920x1080 | Todos | ⬜ |
| Chrome Mobile | 375x812 | Flujos principales | ⬜ |
| Safari Desktop | 1440x900 | Flujos principales | ⬜ |

### Regresión

| # | Funcionalidad existente | Qué verificar | Prioridad | Estado |
|---|------------------------|---------------|-----------|--------|
| R1 | Feature X existente | Que siga funcionando Y | Alta | ⬜ |
| R2 | Feature Z existente | Que no se rompa W | Media | ⬜ |

### Criterios de aceptación

<!--
  ¿Cuándo se considera que la feature está lista para producción?
-->

- [ ] Todos los casos de prioridad Alta pasan ✅
- [ ] Todos los casos de permisos pasan ✅
- [ ] Sin bugs bloqueantes ni críticos abiertos
- [ ] Regresión verificada en flujos impactados
- [ ] Feature flag funciona correctamente (on/off)

---

## Alternativas evaluadas 🔀

### Alternativa A: [nombre]

**Descripción:** Qué propone esta alternativa.

**Ventajas:**
- ✅ Ventaja 1
- ✅ Ventaja 2

**Desventajas:**
- ❌ Desventaja 1
- ❌ Desventaja 2

### Alternativa B: [nombre]

**Descripción:** Qué propone esta alternativa.

**Ventajas:**
- ✅ Ventaja 1

**Desventajas:**
- ❌ Desventaja 1
- ❌ Desventaja 2

### Alternativa C: [nombre] (elegida)

**Descripción:** Qué propone esta alternativa.

**Ventajas:**
- ✅ Ventaja 1
- ✅ Ventaja 2

**Desventajas:**
- ❌ Desventaja 1

### Tabla comparativa

| Criterio | Alt. A | Alt. B | Alt. C (elegida) |
|----------|:------:|:------:|:----------------:|
| Complejidad | Baja | Media | Media |
| Tiempo estimado | X sprints | Y sprints | Z sprints |
| Mantenibilidad | ✅ | ❌ | ✅ |
| Escalabilidad | ❌ | ✅ | ✅ |
| Requiere migración | No | Sí | Sí |

> **Opción elegida: Alternativa C** — Justificación del por qué.

---

## Rollout 📈

<!--
  Plan de despliegue progresivo.
  En Alizia el rollout es por organización (provincia).
-->

| Fase | Alcance | Criterio para avanzar | Responsable |
|------|---------|----------------------|-------------|
| 1 | Staging — equipo interno | Sin bugs bloqueantes | @nombre |
| 2 | Organización piloto (provincia X) | Feedback positivo, sin errores críticos | @nombre |
| 3 | Todas las organizaciones | Métricas de éxito alcanzadas | @nombre |

### Plan de rollback

1. Desactivar feature flag `ff_nombre_feature`
2. Verificar que los usuarios ven el comportamiento anterior
3. Comunicar al equipo y documentar la causa del rollback

---

## Dependencias 👥

### Internas

| Dependencia | Tipo | Bloqueante | Estado | Responsable |
|-------------|------|------------|--------|-------------|
| RFC / Épica X | Precondición | Sí | ✅ / 🟡 / 🔴 | @nombre |
| Migración Y | Técnica | Sí | ⬜ | @nombre |

### Externas

| Dependencia | Tipo | Bloqueante | Estado | Contacto |
|-------------|------|------------|--------|----------|
| API tercero | Integración | Sí / No | ✅ / 🟡 | @contacto |
| Diseño UX | Entregable | Sí | ✅ / 🟡 | @diseñador |

---

## Riesgos ⚠️

<!--
  Riesgos identificados y plan de mitigación.
  Obligatorio en RFCs grandes.
-->

| # | Riesgo | Probabilidad | Impacto | Mitigación |
|---|--------|-------------|---------|------------|
| 1 | Descripción del riesgo | Alta / Media / Baja | Alto / Medio / Bajo | Qué se hace para mitigar |
| 2 | Descripción | Media | Alto | Mitigación |

---

## Preguntas abiertas ❓

| # | Pregunta | Área | Responsable | Estado |
|---|----------|------|-------------|--------|
| 1 | ¿Pregunta? | Producto | @nombre | 🟡 Pendiente |
| 2 | ¿Pregunta? | Frontend | @nombre | 🟡 Pendiente |
| 3 | ¿Pregunta? | Backend | @nombre | 💬 Respondida: respuesta |
| 4 | ¿Pregunta? | UX | @nombre | 🔴 Demorado |

---

## Glosario 📖

<!--
  Términos de dominio que no todo el equipo conoce.
  Si todos los términos son obvios, borrá esta sección.
-->

| Término | Definición |
|---------|-----------|
| Término 1 | Definición |
| Término 2 | Definición |

---

## Tareas 📝

<!--
  Agrupadas por disciplina. Cada tarea referencia la sección del RFC.
  Esto se convierte en el backlog de implementación.
-->

### Backend

| # | Tarea | Sección | Asignado | Estado | Ticket |
|---|-------|---------|----------|--------|--------|
| B1 | Crear migración | Modelo de datos | @nombre | ⬜ | — |
| B2 | Implementar endpoint X | Endpoints | @nombre | ⬜ | — |
| B3 | Implementar endpoint Y | Endpoints | @nombre | ⬜ | — |
| B4 | Lógica de negocio Z | Lógica | @nombre | ⬜ | — |
| B5 | Configurar feature flag | Feature flags | @nombre | ⬜ | — |

### Frontend

| # | Tarea | Sección | Asignado | Estado | Ticket |
|---|-------|---------|----------|--------|--------|
| F1 | Crear página A | Mapa de pantallas | @nombre | ⬜ | — |
| F2 | Crear componente B | Componentes | @nombre | ⬜ | — |
| F3 | Integrar con endpoint X | Estado y data fetching | @nombre | ⬜ | — |
| F4 | Integrar con endpoint Y | Estado y data fetching | @nombre | ⬜ | — |

### UX

| # | Tarea | Sección | Asignado | Estado | Ticket |
|---|-------|---------|----------|--------|--------|
| U1 | Diseño pantalla A | Recomendaciones UX | @nombre | ⬜ | — |
| U2 | Diseño pantalla B | Recomendaciones UX | @nombre | ⬜ | — |
| U3 | Empty states | Recomendaciones UX | @nombre | ⬜ | — |

### QA

| # | Tarea | Sección | Asignado | Estado | Ticket |
|---|-------|---------|----------|--------|--------|
| Q1 | Testing flujo 1 | Estrategia de testing | @nombre | ⬜ | — |
| Q2 | Testing flujo 2 | Estrategia de testing | @nombre | ⬜ | — |
| Q3 | Testing permisos | Estrategia de testing | @nombre | ⬜ | — |
| Q4 | Regresión | Estrategia de testing | @nombre | ⬜ | — |

---

<!--
  ╔══════════════════════════════════════════════════════════════════════════╗
  ║  Checklist antes de enviar a revisión:                                  ║
  ║                                                                          ║
  ║  PRODUCTO                                                                ║
  ║  □ ¿Contexto y motivación claros?                                       ║
  ║  □ ¿Objetivos, no-objetivos y métricas de éxito?                        ║
  ║  □ ¿Alcance con incluye / no incluye / fases?                           ║
  ║  □ ¿Flujos de usuario paso a paso?                                      ║
  ║  □ ¿Reglas de negocio documentadas?                                     ║
  ║  □ ¿Decisiones por provincia listadas?                                  ║
  ║                                                                          ║
  ║  FRONTEND                                                                ║
  ║  □ ¿Mapa de pantallas completo?                                         ║
  ║  □ ¿Componentes listados con acción (crear/modificar)?                  ║
  ║  □ ¿Stores y queries/mutations definidas?                               ║
  ║  □ ¿Estados de UI (loading, error, vacío) por pantalla?                 ║
  ║                                                                          ║
  ║  BACKEND                                                                 ║
  ║  □ ¿Endpoints con request, response y errores?                          ║
  ║  □ ¿Modelo de datos con SQL y migraciones?                              ║
  ║  □ ¿Lógica de negocio con reglas y motivos?                             ║
  ║  □ ¿Configuración por org documentada?                                  ║
  ║  □ ¿Feature flags listados?                                             ║
  ║                                                                          ║
  ║  UX                                                                      ║
  ║  □ ¿Pedidos concretos por pantalla?                                     ║
  ║  □ ¿Contexto suficiente para que UX trabaje?                            ║
  ║                                                                          ║
  ║  QA                                                                      ║
  ║  □ ¿Precondiciones claras?                                              ║
  ║  □ ¿Casos por flujo con happy path + error + edge case?                 ║
  ║  □ ¿Testing de permisos/roles?                                          ║
  ║  □ ¿Regresión identificada?                                             ║
  ║  □ ¿Criterios de aceptación definidos?                                  ║
  ║                                                                          ║
  ║  GENERAL                                                                 ║
  ║  □ ¿Al menos 2 alternativas evaluadas con tabla comparativa?            ║
  ║  □ ¿Plan de rollout con rollback?                                       ║
  ║  □ ¿Dependencias y riesgos identificados?                               ║
  ║  □ ¿Tareas agrupadas por disciplina?                                    ║
  ║  □ ¿Borré los comentarios HTML?                                         ║
  ╚══════════════════════════════════════════════════════════════════════════╝
-->
