# Épica 10: Cosmos

> Personalización de la plataforma por cliente: identidad visual, feature flags, nomenclatura configurable y configuración de IA por organización.

**Estado:** MVP (transversal — no tiene fase propia, se implementa dentro de las épicas que lo consumen)
**Fase de implementación:** Transversal (se materializa en HU-3.1, HU-6.5, HU-2.4, HU-8.2)

---

## Problema

Cada ministerio provincial tiene su propia identidad, nomenclatura, estructura curricular y prioridades pedagógicas. Un sistema genérico genera fricción y rechazo. Se necesita que cada organización vea la plataforma como propia, con sus colores, términos y comportamiento de IA adaptado.

## Objetivos

- Permitir que cada organización tenga su identidad visual (colores, logo, nombre de plataforma)
- Habilitar/deshabilitar módulos y funcionalidades por cliente (feature flags)
- Adaptar la terminología del sistema a los términos de cada provincia
- Configurar el comportamiento del asistente IA por contexto del cliente

## Alcance MVP

**Incluye:**

- Feature flags binarios por módulo (on/off por organización)
- Nomenclatura configurable (nombres de niveles de topics, secciones de documento)
- Configuración de tipos de recurso, actividades y secciones de documento por organización
- Prompts y JSON Schemas por organización para generaciones de IA
- Identidad visual básica (nombre de plataforma, logo, colores primarios)

**No incluye:**

- Panel self-service para que el cliente configure (el equipo de Alizia configura manualmente)
- Marketplace de módulos de terceros → horizonte
- Niveles intermedios de feature flags (solo binario en MVP)

## Principios de diseño

- **Configuración, no código:** Los cambios por cliente se resuelven con configuración (JSON, feature flags), no con branches ni deploys separados
- **Defaults sensatos:** Una organización nueva arranca con una configuración base funcional
- **El cliente se ve reflejado:** Colores, logos y terminología hacen que la plataforma se sienta como propia del ministerio

---

## Implementación técnica

### Cosmos = `organizations.config` (JSONB)

Cosmos no es un servicio separado — es el campo `config` JSONB de la tabla `organizations`. Toda la personalización por cliente vive ahí. Esta decisión sigue el patrón de TUNI.

### Schema completo de `organizations.config`

```json
{
  // --- Estructura curricular (HU-3.1, HU-3.3) ---
  "topic_max_levels": 3,
  "topic_level_names": ["Núcleos problemáticos", "Áreas de conocimiento", "Categorías"],
  "topic_selection_level": 3,

  // --- Clases y enseñanza (HU-3.5) ---
  "shared_classes_enabled": true,
  "desarrollo_max_activities": 3,

  // --- Secciones del documento de coordinación (HU-4.3) ---
  "coord_doc_sections": [
    {
      "key": "problem_edge",
      "label": "Eje problemático",
      "type": "text",
      "ai_prompt": "Generá un eje problemático que articule las categorías seleccionadas...",
      "required": true
    },
    {
      "key": "methodological_strategy",
      "label": "Estrategia metodológica",
      "type": "select_text",
      "options": ["proyecto", "taller_laboratorio", "ateneo_debate"],
      "ai_prompt": "Generá una estrategia metodológica detallada...",
      "required": true
    },
    {
      "key": "eval_criteria",
      "label": "Criterios de evaluación",
      "type": "text",
      "ai_prompt": "Generá criterios de evaluación...",
      "required": false
    }
  ],

  // --- Personalización de IA (HU-6.5) ---
  "ai_settings": {
    "tone": "neutral_formal",
    "max_generation_length": 2000,
    "max_chat_interactions": 20,
    "enabled_features": ["generation", "chat"],
    "curriculum_sources": [],
    "vocabulary_overrides": {},
    "prompt_overrides": {}
  },

  // --- Feature flags por módulo ---
  "modules": {
    "planificacion": true,
    "contenido": true,
    "aprendizaje_adaptativo": false,
    "inclusion": false
  },

  // --- Identidad visual ---
  "visual_identity": {
    "platform_name": "Alizia",
    "logo_url": null,
    "primary_color": "#6366F1",
    "secondary_color": "#1E1B4B"
  },

  // --- Nomenclatura configurable ---
  "nomenclature": {
    "coordination_document": "Documento de coordinación",
    "lesson_plan": "Planificación docente",
    "area": "Área"
  },

  // --- Onboarding (HU-2.4) ---
  "onboarding": {
    "allow_skip": true,
    "profile_fields": [],
    "tour_steps": []
  }
}
```

### Valores de tono de IA

| Valor | Descripción |
|-------|-------------|
| `neutral_formal` | Español neutro, tono profesional (default) |
| `neutral_informal` | Español neutro, tono cercano |
| `rioplatense_formal` | Voseo, tono profesional |
| `rioplatense_informal` | Voseo, tono cercano |

### Endpoints

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `GET /api/v1/organizations/:id` | GET | Retorna organización con config completo |
| `PATCH /api/v1/organizations/:id` | PATCH | Actualiza config con merge parcial de JSONB |

> Los endpoints ya existen en HU-3.1. Cosmos no agrega endpoints propios — usa la misma API de organizaciones.

---

## Historias de usuario

La mayor parte de Cosmos se implementa de forma transversal en otras épicas, con una HU propia para la configuración pública:

| HU | Épica | Qué implementa de Cosmos |
|----|-------|---------------------------|
| [HU-10.1](./HU-10.1-config-publica-org/HU-10.1-config-publica-org.md) | Cosmos | Endpoint public-config, middleware RequireModule, DefaultConfig (modules, visual_identity, nomenclature) |
| [HU-3.1](../03-integracion/HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) | Integración | Tabla organizations con config JSONB, endpoints, seed |
| [HU-6.5](../06-asistente-ia/HU-6.5-customizacion-organizacion/HU-6.5-customizacion-organizacion.md) | Asistente IA | `ai_settings`: tono, límites, fuentes, vocabulario, prompt overrides |
| [HU-2.4](../02-onboarding/HU-2.4-configuracion-onboarding-por-org/HU-2.4-configuracion-onboarding-por-org.md) | Onboarding | `onboarding`: campos de perfil, pasos del tour |
| [HU-8.2](../08-contenido-recursos/HU-8.2-tipos-recurso-configuracion/HU-8.2-tipos-recurso-configuracion.md) | Contenido | Tipos de recurso configurables por org |

### Pendientes de implementación (MVP)

| Funcionalidad | Estado | Nota |
|---------------|--------|------|
| Config JSONB base (topics, secciones, shared classes) | Cubierto por HU-3.1 | Ya definido |
| AI settings (tono, límites, vocabulario) | Cubierto por HU-6.5 | Post-MVP en RFC actual — **producto dice MVP** |
| Feature flags por módulo (`modules`) | Cubierto por HU-10.1 | Middleware RequireModule + DefaultConfig |
| Identidad visual (`visual_identity`) | Cubierto por HU-10.1 | Endpoint public-config + DefaultConfig |
| Nomenclatura configurable (`nomenclature`) | Cubierto por HU-10.1 | Endpoint public-config + DefaultConfig |
| Onboarding config | Cubierto por HU-2.4 | Post-MVP (Épica 2 es post-MVP) |

---

## Decisiones técnicas

- **Cosmos = `organizations.config`**: No es un servicio separado. Toda la personalización vive en el JSONB de la tabla organizations. Patrón idéntico a TUNI.
- **Feature flags binarios en MVP**: Cada módulo se activa o desactiva con un boolean. No hay niveles intermedios ni reglas complejas.
- **Defaults via Go function**: Al crear una organización, se aplica `DefaultConfig()` que tiene valores sensatos para todos los campos.
- **Merge parcial en PATCH**: El endpoint de actualización hace merge a nivel de keys raíz del JSONB, no reemplaza todo el objeto.
- **Sin UI de administración**: El equipo de implementación configura via API directa (curl/Postman) o seeds. No hay panel admin en MVP.

## Decisiones de cada cliente

- Qué módulos están habilitados (planificación, contenido, etc.)
- Identidad visual (colores, logo, nombre)
- Nomenclatura (cómo llaman a documentos, áreas, etc.)
- Configuración de IA (tono, fuentes curriculares, vocabulario)
- Secciones del documento de coordinación (qué secciones, en qué orden, cuáles son obligatorias)

## Épicas relacionadas

- **[Épica 3: Integración](../03-integracion/03-integracion.md)** — HU-3.1 implementa la tabla y endpoints de organizations.config
- **[Épica 6: Asistente IA](../06-asistente-ia/06-asistente-ia.md)** — HU-6.5 implementa ai_settings
- **[Épica 2: Onboarding](../02-onboarding/02-onboarding.md)** — HU-2.4 implementa onboarding config
- **[Épica 4: Documento de coordinación](../04-documento-coordinacion/04-documento-coordinacion.md)** — Consume coord_doc_sections
- **[Épica 8: Contenido y recursos](../08-contenido-recursos/08-contenido-recursos.md)** — Consume resource types config
