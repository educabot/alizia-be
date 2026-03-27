---
name: rfc-docs
description: "Genera y organiza documentación RFC del proyecto: épicas, historias de usuario (HU) y tareas (T). Usa este skill para crear nuevas épicas, desglosar historias, agregar tareas, o reorganizar la documentación existente."
---

# RFC Documentation Generator

Generás y organizás la documentación RFC de Alizia v2 siguiendo la estructura y convenciones establecidas.

## Ubicación

Toda la documentación vive en `docs/rfc-alizia/epicas/`.

## Estructura de directorios

```
docs/rfc-alizia/epicas/
├── epicas.md                              ← Índice general de todas las épicas
├── XX-nombre-epica.md                     ← Summary suelto (épicas sin desglosar)
├── XX-nombre-epica/                       ← Épica desglosada
│   ├── XX-nombre-epica.md                 ← Archivo principal de la épica
│   ├── HU-X.1-nombre-historia/            ← Historia de usuario
│   │   ├── HU-X.1-nombre-historia.md      ← Archivo principal de la HU
│   │   └── tareas/                        ← Subcarpeta de tareas
│   │       ├── T-X.1.1-nombre-tarea.md
│   │       ├── T-X.1.2-nombre-tarea.md
│   │       └── ...
│   ├── HU-X.2-nombre-historia/
│   │   ├── HU-X.2-nombre-historia.md
│   │   └── tareas/
│   └── ...
```

## Convenciones de naming

### CRÍTICO: Nunca usar README.md
Los archivos se nombran igual que su carpeta contenedora. Esto permite identificar el archivo cuando hay muchos tabs abiertos.

- Carpeta `01-roles-accesos/` → archivo `01-roles-accesos.md`
- Carpeta `HU-1.1-autenticacion-auth0/` → archivo `HU-1.1-autenticacion-auth0.md`

### Prefijos
| Prefijo | Significado | Ejemplo |
|---------|-------------|---------|
| `XX-` | Número de épica (00, 01, 02...) | `01-roles-accesos` |
| `HU-X.Y-` | Historia de Usuario (épica.historia) | `HU-1.2-modelo-usuarios-roles` |
| `T-X.Y.Z-` | Tarea (épica.historia.tarea) | `T-1.2.1-migracion` |

### Idioma
- Nombres de carpetas/archivos: **español** (kebab-case)
- Contenido de documentos: **español**
- Código, SQL, nombres técnicos dentro de los docs: **inglés**

### Links internos
Siempre apuntar al archivo explícito, nunca a la carpeta:
- ✅ `./HU-1.1-autenticacion-auth0/HU-1.1-autenticacion-auth0.md`
- ❌ `./HU-1.1-autenticacion-auth0/`

---

## Templates

### Template: Archivo principal de Épica (`XX-nombre.md`)

```markdown
# Épica X: Nombre

> Descripción de una línea.

**Estado:** MVP | Post-MVP | Pendiente definición
**Fase de implementación:** Fase N

---

## Problema

Qué problema resuelve esta épica.

## Objetivos

- Objetivo 1
- Objetivo 2

## Alcance MVP

**Incluye:**

- Feature 1
- Feature 2

**No incluye:**

- Feature futura → horizonte
- Feature por definir → por definir

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-X.1 | [Nombre](./HU-X.1-nombre/HU-X.1-nombre.md) | Descripción corta | Fase N | N |
| HU-X.2 | [Nombre](./HU-X.2-nombre/HU-X.2-nombre.md) | Descripción corta | Fase N | N |

---

## Decisiones técnicas

- Decisión 1
- Decisión 2

## Decisiones de cada cliente

- Decisión configurable 1

## Épicas relacionadas

- **Nombre épica** — Relación

## Test cases asociados

- Fase N: Tests X.X–X.X (descripción)

Ver [testing.md](../../operaciones/testing.md) para la matriz completa.
```

### Template: Historia de Usuario (`HU-X.Y-nombre.md`)

```markdown
# HU-X.Y: Nombre de la historia

> Como [rol], necesito [acción] para [beneficio].

**Fase:** N — Nombre de fase
**Prioridad:** Alta | Media | Baja
**Estimación:** —

---

## Criterios de aceptación

- [ ] Criterio 1
- [ ] Criterio 2
- [ ] Criterio 3

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| X.Y.1 | [Nombre tarea](./tareas/T-X.Y.1-nombre.md) | path/al/archivo | ⬜ |
| X.Y.2 | [Nombre tarea](./tareas/T-X.Y.2-nombre.md) | path/al/archivo | ⬜ |

## Dependencias

- HU-X.Y completada (qué se necesita)
- Épica N completada (qué se necesita)

## Test cases

- X.Y: Descripción del test → resultado esperado
```

### Template: Tarea (`T-X.Y.Z-nombre.md`)

```markdown
# T-X.Y.Z: Nombre de la tarea

**Historia:** HU-X.Y — Nombre de la historia
**Tipo:** Backend | Frontend | Infra | CI/CD | Config | Testing | DB
**Estado:** ⬜ Pendiente
**Fase:** Post-MVP  ← (incluir solo si la tarea es Post-MVP)

---

## Descripción

Qué hay que hacer y por qué.

## Implementación

(Código de referencia, SQL, configuración, pasos — lo que aplique)

## Notas

- Nota técnica relevante
- Decisión de diseño
```

### Tareas Post-MVP

Las HUs marcadas como Post-MVP **también deben tener tareas desglosadas** si tienen suficiente detalle técnico (entidades, tablas, endpoints, decisiones de diseño definidas). Cada tarea Post-MVP debe incluir:

- `**Fase:** Post-MVP` en el header del archivo
- Mismo nivel de detalle que las tareas MVP (migración, entities, endpoints, tests)
- Código de referencia en Go siguiendo los patrones existentes del proyecto

---

## Numeración

La numeración es jerárquica y consistente:

```
Épica 1
├── HU-1.1 (primera historia de la épica 1)
│   ├── T-1.1.1 (primera tarea de HU-1.1)
│   ├── T-1.1.2
│   └── T-1.1.3
├── HU-1.2 (segunda historia de la épica 1)
│   ├── T-1.2.1
│   └── T-1.2.2
```

Si se agrega una HU entre existentes, renumerar las siguientes. Si se mueve una HU a otra épica, renumerar ambas épicas.

## Índice general (epicas.md)

Cuando se crea una nueva épica, actualizar `docs/rfc-alizia/epicas/epicas.md` con la nueva entrada en la tabla correspondiente (MVP o Post-MVP).

## Estados de tareas

| Emoji | Estado |
|-------|--------|
| ⬜ | Pendiente |
| 🔄 | En progreso |
| ✅ | Completada |

## Proceso para crear una épica nueva

1. Crear el `.md` suelto si solo es definición de producto (sin desglose técnico)
2. Cuando se desglose, crear la carpeta `XX-nombre/` con su archivo principal
3. Crear subcarpetas `HU-X.Y-nombre/` por cada historia
4. Dentro de cada HU, crear `tareas/` con los archivos `T-X.Y.Z-nombre.md`
5. Actualizar `epicas.md` (índice general)
6. Actualizar links en archivos que referencien la épica

## Proceso para agregar una HU a épica existente

1. Crear carpeta `HU-X.Y-nombre/` dentro de la épica
2. Crear `HU-X.Y-nombre.md` y `tareas/` con sus tareas
3. Actualizar la tabla de historias en el archivo principal de la épica
4. Verificar que la numeración sea consistente

## Referencia rápida de archivos existentes

Antes de crear, verificar qué épicas ya están desglosadas:
```bash
# Ver estructura actual
find docs/rfc-alizia/epicas -type d | sort

# Ver épicas sueltas (sin desglosar)
ls docs/rfc-alizia/epicas/*.md

# Ver épicas desglosadas
ls -d docs/rfc-alizia/epicas/*/
```
