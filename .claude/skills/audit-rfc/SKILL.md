---
name: audit-rfc
description: "Audita la documentación RFC de Alizia: completitud, consistencia interna, links rotos, alineación entre docs técnicos y épicas/HUs/tareas."
---

# RFC Audit — Alizia

Ejecutás una auditoría completa de la documentación RFC del proyecto, identificando problemas de completitud, consistencia, y alineación.

## Ubicación del RFC

```
docs/rfc-alizia/
├── rfc-alizia.md                  ← Documento central
├── epicas/                        ← Épicas, HUs, Tareas
│   ├── epicas.md                  ← Índice general
│   └── XX-nombre/                 ← Épicas desglosadas
├── tecnico/                       ← Docs técnicos
│   ├── arquitectura.md
│   ├── modelo-de-datos.md
│   ├── endpoints.md
│   ├── errores.md
│   ├── prompts.md
│   ├── frontend-integration.md
│   └── team-ai-toolkit.md
└── operaciones/
    ├── testing.md
    └── deploy.md
```

## Proceso de auditoría

Lanzar **3 agentes en paralelo**, cada uno con un foco distinto:

### Agente 1: Auditoría de Épicas y Flujo de Trabajo

**Foco:** Estructura de épicas, HUs y tareas.

**Checklist:**
- [ ] Todas las épicas listadas en `epicas.md` existen como directorio
- [ ] Cada épica desglosada tiene archivo principal con tabla de HUs
- [ ] Cada HU tiene carpeta con archivo principal y subcarpeta `tareas/`
- [ ] Cada tarea referenciada en la HU existe como archivo
- [ ] Numeración jerárquica consistente (épica → HU → tarea)
- [ ] Links internos apuntan a archivos existentes (no a carpetas, no rotos)
- [ ] Estados de tareas actualizados (sin contradicciones)
- [ ] Naming sigue convenciones: `XX-nombre/XX-nombre.md`, `HU-X.Y-nombre/`, `T-X.Y.Z-nombre.md`
- [ ] No hay archivos `README.md` (convención del proyecto)
- [ ] Cada épica tiene secciones: Problema, Objetivos, Alcance MVP, Historias, Decisiones técnicas
- [ ] Épicas Post-MVP también tienen tareas desglosadas si hay suficiente detalle técnico

**Output:** Lista de issues con severidad (CRITICAL / WARNING / INFO).

### Agente 2: Auditoría Técnica

**Foco:** Consistencia entre documentos técnicos y el RFC central.

**Checklist:**
- [ ] `rfc-alizia.md` referencia todos los docs técnicos en su índice
- [ ] Modelo de datos (`modelo-de-datos.md`):
  - Todas las tablas mencionadas en el RFC existen en el modelo
  - Campos JSONB tienen schema documentado
  - Triggers definidos con SQL completo
  - Índices definidos para queries frecuentes
  - Relaciones FK son consistentes con las entidades
- [ ] Endpoints (`endpoints.md`):
  - Cada endpoint tiene: método, ruta, roles, request body, response body, errores
  - Endpoints alineados con las HUs (no hay endpoints huérfanos ni HUs sin endpoint)
  - Paginación documentada donde aplica
- [ ] Errores (`errores.md`):
  - Cada código de error referenciado en endpoints está catalogado
  - Transiciones de estado documentadas
  - Códigos HTTP correctos para cada tipo de error
- [ ] Prompts (`prompts.md`):
  - Cada prompt tiene placeholders documentados
  - Variables de contexto tienen fuente identificada
  - Output schemas definidos donde se espera JSON
- [ ] Frontend integration (`frontend-integration.md`):
  - Tipos TypeScript alineados con el modelo de datos
  - Patrones de paginación consistentes con el backend
  - Manejo de errores cubre los códigos del catálogo
- [ ] Arquitectura (`arquitectura.md`):
  - Capas documentadas: entities, providers, usecases, handlers, repositories
  - Import paths usan `go-alizia/src/...`
  - Decisiones técnicas clave documentadas con justificación
- [ ] Stack tecnológico del RFC central coincide con los docs técnicos
- [ ] No hay contradicciones entre documentos (ej: un doc dice GORM, otro dice sqlx)

**Output:** Lista de inconsistencias con los archivos involucrados y líneas.

### Agente 3: Inventario y Cobertura

**Foco:** Qué falta por documentar.

**Checklist:**
- [ ] Contar épicas totales vs desglosadas vs pendientes
- [ ] Contar HUs totales vs con tareas vs sin tareas
- [ ] Contar tareas totales por estado (pendiente/progreso/completada)
- [ ] Identificar HUs que mencionan endpoints no documentados en `endpoints.md`
- [ ] Identificar tablas del modelo no cubiertas por ninguna épica
- [ ] Identificar docs técnicos referenciados pero inexistentes
- [ ] Verificar que `testing.md` tiene fixtures/seeds para las entidades principales
- [ ] Verificar que `deploy.md` existe y está actualizado
- [ ] Cross-reference: para cada fase del RFC, verificar que todas sus épicas tienen HUs desglosadas

**Output:** Tabla resumen con métricas + lista de gaps.

## Formato del reporte final

Consolidar los 3 agentes en un reporte único con esta estructura:

```markdown
# Auditoría RFC Alizia — {fecha}

## Resumen ejecutivo

| Métrica | Valor |
|---------|-------|
| Épicas totales | X |
| Épicas desglosadas | X |
| HUs totales | X |
| Tareas totales | X (⬜ Y / 🔄 Z / ✅ W) |
| Docs técnicos | X / Y esperados |
| Issues CRITICAL | X |
| Issues WARNING | X |

## 1. Issues críticos (CRITICAL)

Problemas que bloquean el desarrollo o causan confusión.

- **[CRITICAL-001]** Descripción — archivo(s) afectado(s)

## 2. Warnings (WARNING)

Problemas menores que deberían resolverse.

- **[WARNING-001]** Descripción — archivo(s) afectado(s)

## 3. Informativos (INFO)

Sugerencias de mejora.

- **[INFO-001]** Descripción

## 4. Inconsistencias entre documentos

| Doc A | Doc B | Inconsistencia |
|-------|-------|---------------|
| ... | ... | ... |

## 5. Documentos faltantes

| Documento | Prioridad | Motivo |
|-----------|-----------|--------|
| ... | ... | ... |

## 6. Cobertura por fase

| Fase | Épicas | HUs desglosadas | Tareas | Cobertura |
|------|--------|-----------------|--------|-----------|
| ... | ... | ... | ... | ...% |
```

## Notas importantes

- **No modificar nada** durante la auditoría — solo reportar
- Usar paths relativos desde la raíz del proyecto en el reporte
- Si un issue tiene fix obvio, sugerirlo en la descripción
- Priorizar CRITICAL sobre WARNING sobre INFO
- El reporte se imprime en la conversación, NO se guarda como archivo (a menos que el usuario lo pida)

## Convenciones a validar (del skill rfc-docs)

- Archivos se nombran igual que su carpeta contenedora (nunca README.md)
- Prefijos: `XX-` épica, `HU-X.Y-` historia, `T-X.Y.Z-` tarea
- Nombres en español kebab-case, código/SQL en inglés
- Links internos apuntan a archivos explícitos, no carpetas
- Import paths: `go-alizia/src/...` (no `go-alizia-v2/src/...`)
- Auth: JWT via team-ai-toolkit (no Auth0)
- Deploy: Railway (no Cloud Run/Functions)
- El proyecto se llama "Alizia" (sin V2)
