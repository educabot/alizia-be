# HU-8.3: Generación de recursos con IA

> Como docente, necesito generar un recurso didáctico seleccionando un tipo, una fuente (si aplica) y el contexto de mi materia, para obtener contenido listo para usar en el aula.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] El docente puede seleccionar un tipo de recurso disponible para su organización
- [ ] Si el tipo requiere fuente (`requires_font`), el docente selecciona/ingresa la referencia
- [ ] Para "Creación libre", el docente escribe instrucciones de lo que quiere generar
- [ ] Se genera el recurso usando el AI Provider con el prompt y schema resueltos
- [ ] El contenido generado se almacena en `resources.content` según el output_schema
- [ ] El recurso se crea en estado `generated` tras la generación exitosa
- [ ] El contexto de materia/curso se incluye en el prompt para personalización

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.3.1 | [Usecase de generación](./tareas/T-8.3.1-usecase-generacion.md) | usecases/resources/ | ⬜ |
| 8.3.2 | [Endpoint de creación y generación](./tareas/T-8.3.2-endpoint-generacion.md) | handlers/ | ⬜ |
| 8.3.3 | [Context builder para recursos](./tareas/T-8.3.3-context-builder.md) | usecases/ai/ | ⬜ |
| 8.3.4 | [Tests](./tareas/T-8.3.4-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-8.1: Modelo de datos](../HU-8.1-modelo-datos-recursos/HU-8.1-modelo-datos-recursos.md) — Tablas base
- [HU-8.2: Tipos de recurso](../HU-8.2-tipos-recurso-configuracion/HU-8.2-tipos-recurso-configuracion.md) — Resolución de prompt/schema
- [HU-6.1: AI Provider](../../06-assistente-ia/HU-6.1-ai-provider/HU-6.1-ai-provider.md) — GenerateJSON

## Test cases

- 8.8: Generar guía de lectura con fuente → content con estructura del schema
- 8.9: Generar ficha de cátedra sin fuente → content generado correctamente
- 8.10: Creación libre con instrucciones custom → respeta las instrucciones
- 8.11: Tipo que requiere fuente sin fuente → error de validación
