# Épicas — Alizia

## Índice

### MVP (9 épicas)

| # | Épica | Fase | Estado |
|---|-------|------|--------|
| 0 | [Setup e infraestructura](./00-setup-infraestructura/00-setup-infraestructura.md) | Fase 1 | MVP (9 tareas) |
| 1 | [Roles y accesos](./01-roles-accesos/01-roles-accesos.md) | Fase 1–2 | MVP |
| 2 | [Onboarding](./02-onboarding/02-onboarding.md) | Fase 2 | MVP |
| 3 | [Integración](./03-integracion/03-integracion.md) | Fase 2 | MVP |
| 4 | [Documento de coordinación](./04-documento-coordinacion/04-documento-coordinacion.md) | Fase 3 | MVP |
| 5 | [Planificación docente](./05-planificacion-docente/05-planificacion-docente.md) | Fase 5 | MVP |
| 6 | [Asistente IA](./06-asistente-ia/06-asistente-ia.md) | Fase 4 | MVP |
| 8 | [Contenido y recursos](./08-contenido-recursos/08-contenido-recursos.md) | Fase 6 | MVP |
| 7 | [Dashboard](./07-dashboard/07-dashboard.md) | Fase 7 | MVP |

### Post-MVP / Futuro

| # | Épica | Estado |
|---|-------|--------|
| 9 | [WhatsApp](./09-whatsapp/09-whatsapp.md) | Pendiente definición |
| 10 | [Cosmos](./10-cosmos/10-cosmos.md) | MVP (transversal — config por org) |
| 11 | [Aprendizaje Adaptativo](./11-aprendizaje-adaptativo/11-aprendizaje-adaptativo.md) | Horizonte |
| 12 | [Inclusión](./12-inclusion/12-inclusion.md) | Horizonte |

---

## Mapa de dependencias

```
Setup e infra ──→ Roles y accesos ──→ Onboarding ──→ (usuario listo para usar la plataforma)
        │
        ▼
   Integración ──→ Documento de coordinación ──→ Planificación docente
        │                    │                          │
        │                    ▼                          ▼
        │              Asistente IA ◄───────────── Bitácora
        │                    │
        ▼                    ▼
    Contenido ◄──────── Asistente IA

   Dashboard ◄── (consume estado de Docs + Planificaciones)

   WhatsApp ◄── Asistente IA (canal alternativo)

   Cosmos ── (transversal: organizations.config JSONB)
```

---

## Patrones transversales

Estos patrones aparecen en múltiples épicas:

| Patrón | Épicas que lo usan | Descripción |
|---|---|---|
| JSON de configuración por org | 1, 2, 3, 4, 5, 8 | Configuración provincial centralizada (feature flags, nombres de niveles, tipos habilitados) |
| Prompt + JSON Schema por sección | 4, 6, 8 | Cada output generado por IA tiene su prompt y schema configurable por provincia |
| Feature flags por organización | 2, 4, 5, 8 | Funcionalidades que se activan/desactivan por cliente |
| Clases coordinadas | 4, 5 | Diferenciador clave: sincronización entre docentes que comparten horario |
| Decisiones por provincia | Todas | Cada cliente (provincia) puede customizar comportamiento sin cambios de código |
