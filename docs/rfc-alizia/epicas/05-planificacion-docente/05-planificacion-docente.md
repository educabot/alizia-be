# Épica 5: Planificación docente

> Planificación del clase a clase por momentos, con asistencia de IA y feedback post-dictado.

**Estado:** MVP (parcial — bitácora y repropuesta son post-MVP)
**Fase de implementación:** Fase 5

---

## Problema

Los docentes planifican sus clases sin una conexión clara con lo acordado a nivel área, ni con lo que están dictando en clases los otros docentes. Re-planificar en base a lo que ocurrió en clases anteriores es algo que depende 100% del docente, de su energía y su memoria. Falta un sistema que alinee, proponga, recuerde y aprenda del día a día del aula.

## Objetivos

- Que cada clase planificada esté alineada con el documento de coordinación del área
- Que cada clase sea coherente con lo que se está dictando en otras disciplinas
- Generar una propuesta inicial de actividades por momento (apertura, desarrollo, cierre)
- Incorporar el feedback de clases anteriores para mejorar las propuestas futuras
- Permitir personalización del docente sin perder la alineación curricular

## Alcance MVP

**Incluye:**

- Visualización del cronograma de clases heredado del documento de coordinación
- Edición del objetivo de clase
- Selección de actividades por momento (apertura, desarrollo, cierre) con recomendaciones de IA
- Personalización: el docente puede anclar la clase a un recurso (canción, lectura, etc.) o a un comentario personalizado
- Generación de la propuesta detallada de clase
- Edición directa y asistida por IA de la propuesta
- Gestión de estados de las planificaciones (pending → planned)

**Post-MVP:**

- Bitácora de cotejo: el docente reporta cómo fue la clase (soporta audio)
- Recolección de datos: Alicia pregunta activamente por información faltante o alumnos con un caso que merezca seguimiento
- Repropuesta: En base a la bitácora proponemos cambios a las siguientes clases ya planificadas (en proceso o publicadas) y tendremos eso en cuenta para la generación de las pendientes

**No incluye:**

- ¿Qué pasa si el docente no tiene internet al momento de grabar la bitácora? Escuelas rurales argentinas pueden tener conectividad limitada. ¿Se sube después? ¿Funciona offline? Esto aplica también a la planificación en general.
- Informe de proceso (resumen de progreso del alumno por área) → horizonte
- Trayectorias de refuerzo personalizadas → horizonte
- Creación de recursos didácticos → ver Contenido

## Principios de diseño

- **Del área al aula:** La planificación individual nace del acuerdo colectivo.
- **IA que aprende del aula:** Las propuestas mejoran con el feedback real del docente.
- **Voz del docente:** La bitácora acepta audio libre, sin formato rígido.

## Sub-épicas

| Componente | Descripción |
|---|---|
| Plan de clase | Selección de actividades por momento y generación de la propuesta detallada |
| Momentos didácticos | Configuración de los tipos de actividad por momento (apertura, desarrollo, cierre) |
| Incorporación de fuentes | Anclaje de la clase a un recurso o fuente específica |
| Edición del documento | Permitir que el docente edite directamente o a través de asistente el doc, gestión de su estado |
| Bitácora (post-MVP) | Recolección del resultado de una clase, guiado y simple para el docente |
| Repropuesta (post-MVP) | Sugerencias de cambios para clases futuras ya planificadas (Proactividad del sistema) |

## Decisiones de cada cliente

- Los tipos de actividad disponibles por momento se definen con cada equipo pedagógico provincial
- El formato y profundidad de la bitácora de cotejo requiere validación

## Decisiones técnicas

- Se asume **un docente por materia por curso**. Si excepcionalmente hay dos, opera first-come-first-serve: el primero en planificar escribe, el segundo ve los cambios. El documento de coordinación es la referencia compartida.
- La planificación debe poder **exportarse como PDF** para impresión y para integración con plataformas provinciales existentes donde los docentes reportan sus planificaciones. El formato de exportación será un template configurable.
- Los **momentos didácticos** (apertura, desarrollo, cierre) son fijos como estructura, pero los tipos de actividad dentro de cada momento son configurables por organización. Cada provincia define su catálogo de actividades.
- Cuando existe el feature flag de **clases coordinadas**, la planificación muestra un indicador de clase compartida y las modificaciones se sincronizan entre los docentes involucrados.
- La **bitácora de cotejo** (post-MVP) funciona por audio libre: el docente graba lo que pasó y el sistema lo procesa. No se fuerza un formato rígido porque la adopción depende de que sea natural y rápido.

## Notas

- El informe de proceso y las trayectorias de refuerzo están en horizonte, pendientes de priorización

## Épicas relacionadas

- **Documento de coordinación** — Provee el cronograma y objetivos de clase
- **Contenido** — Recursos disponibles para incorporar en la planificación
- **Asistente IA** — Genera propuestas y procesa feedback
