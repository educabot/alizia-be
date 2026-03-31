# HU-5.9: Exportación PDF de planificación

> Como docente, necesito exportar mi planificación como PDF para impresión y para integración con plataformas provinciales.

**Fase:** Fase 5
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] Endpoint `GET /api/v1/lesson-plans/:id/pdf` retorna un archivo PDF
- [ ] El PDF incluye: título de clase, objetivo, actividades por momento con contenido, fuentes, topics
- [ ] El template de PDF es configurable por organización (logo, colores, estructura)
- [ ] El PDF incluye la información del documento de coordinación de referencia
- [ ] Solo se puede exportar planificaciones en estado `published`

## Concepto

La planificación debe poder exportarse como PDF para impresión y para integración con plataformas provinciales existentes donde los docentes reportan sus planificaciones. El formato de exportación será un template configurable por organización.

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 5.9.1 | Generación de PDF (librería Go) | src/core/usecases/ | ⬜ |
| 5.9.2 | Template configurable por org | src/core/usecases/ | ⬜ |
| 5.9.3 | Endpoint GET PDF | src/entrypoints/ | ⬜ |
| 5.9.4 | Tests | tests/ | ⬜ |

## Dependencias

- [HU-5.5: Estados de planificación](../HU-5.5-estados-planificacion/HU-5.5-estados-planificacion.md) — Solo published se exporta
- [HU-5.3: Planificación por momentos](../HU-5.3-planificacion-por-momentos/HU-5.3-planificacion-por-momentos.md) — Contenido a exportar
