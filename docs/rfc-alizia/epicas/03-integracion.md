# Épica 3: Integración

> Importación de datos y configuración de la estructura curricular de cada provincia.

**Estado:** MVP
**Fase de implementación:** Fase 2

---

## Problema

Cada provincia tiene su propio diseño curricular, terminología y estructura de conocimientos. El sistema necesita incorporar estos datos como base para que todo lo que Alizia genere esté alineado con la realidad de cada jurisdicción.

## Objetivos

- Importar diseños curriculares, NAPs y fuentes oficiales de cada provincia
- Modelar la estructura curricular provincial (áreas, disciplinas, tópicos, sub-tópicos, actividades, tipos de contenidos a crear)
- Que estos datos sean la fuente de verdad para la generación de documentos y las respuestas del asistente IA

## Alcance MVP

**Incluye:**

- Importación de diseño curricular y fuentes oficiales
- Configuración de la estructura curricular (jerarquía de conocimientos, tópicos, etc.)
- Customizaciones por institución: Grillas horarias por institución, alta de docentes, etc.

**No incluye:**

- Gestión de usuarios e instituciones → ver Roles y accesos
- Carga de datos del docente o alumnos → ver Onboarding

*Aclaración: Sí se incluye el impacto y presencia de esos datos por institución, no la carga.*

## Principios de diseño

- **Provincial first:** Cada implementación respeta la estructura y terminología de la provincia.
- **Fuente de verdad única:** Los datos importados alimentan toda la plataforma.

## Sub-épicas

| Componente | Descripción |
|---|---|
| Importación de datos | Carga de diseño curricular, NAPs, fuentes oficiales y grillas horarias |
| Estructura curricular | Modelado de la jerarquía de conocimientos según la provincia |

## Decisiones del cliente

- La granularidad de la estructura curricular varía (Neuquén usa conocimientos y saberes + categorías; otras provincias pueden diferir)

## Decisiones técnicas

- La configuración de cada organización se almacena en un **JSON de configuración** a nivel organización — mismo patrón que usamos en TUNI. Este JSON contiene: nombres de niveles de tópicos, profundidad máxima permitida, feature flags (clases coordinadas, subida de fuentes propias, etc.) y cualquier parametrización provincial.
- Los tópicos se modelan en una **tabla única auto-referencial** (foreign key que apunta a sí misma). La profundidad de cada nodo se **pre-computa y almacena** para evitar recursiones costosas en queries. Si se mueve un tópico, se re-computa solo la rama afectada.
- Cada organización define la **profundidad máxima** de su jerarquía de tópicos y los **nombres por nivel** (ej: nivel 1 = "Conocimientos y saberes", nivel 2 = "Núcleos problemáticos", nivel 3 = "Categorías"). El frontend renderiza dinámicamente según esta configuración.
- El setup inicial de cada provincia (estructura curricular, grillas, docentes) se realiza **manualmente por el equipo de implementación** — no hay backoffice de auto-servicio en el MVP. Esto es deliberado: lo que vendemos es que hacemos el setup del diseño curricular.
- Las migraciones de base de datos se hacen **incrementalmente** a medida que se desarrollan las funcionalidades — no se crea el DER completo upfront. Lección aprendida de TUNI: un schema gigante genera redefiniciones constantes y deuda de mantenimiento.
- Un área puede contener **una o más disciplinas**. Si la provincia no organiza por áreas, el coordinador crea un documento con las materias que correspondan. No se fuerza la existencia de áreas como concepto obligatorio.

## Épicas relacionadas

- **Documento de coordinación** — Consume la estructura curricular para generar documentos
- **Planificación docente** — Usa los datos importados como contexto para la planificación
