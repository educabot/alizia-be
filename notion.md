# **RFC: Épicas de Alizia — Visión completa del producto**

**Fecha:** 2026-03-13 **Autor:** Equipo de producto **Estado:** En revisión — esperando comentarios del equipo

---

> Este documento consolida los overviews de las 10 épicas que componen Alizia. El objetivo es que el equipo revise, comente y valide el alcance, las decisiones técnicas y las dependencias antes de avanzar con las especificaciones detalladas.
>
>
> **Cómo comentar:** Dejá tus comentarios directamente en este documento o en el canal correspondiente, indicando la épica y sección.
>

---

## **1. Roles y accesos**

> Autenticación, roles, permisos y asignación organizacional de usuarios.
>

### **Problema**

La plataforma opera con múltiples roles (coordinador, docente, y potencialmente directivos) que tienen permisos distintos sobre los mismos documentos y cursos. Se necesita un sistema que controle quién puede crear, editar y visualizar cada recurso o informes.

### **Objetivos**

- Autenticar usuarios de forma segura
- Definir roles con permisos diferenciados (coordinador crea documentos, docente planifica clases)
- Asignar usuarios a instituciones, áreas y cursos

### **Alcance MVP**

**Incluye:**

- Autenticación de usuarios
- Roles de coordinador y docente con permisos diferenciados
- Asignación de usuarios a instituciones y cursos

**No incluye:**

- Roles de directivos o supervisores → horizonte
- Gestión de múltiples instituciones por usuario → por definir

### **Principios de diseño**

- **Rol define el flujo:** La experiencia del usuario cambia según su rol desde el inicio.
- **Asignación clara:** Cada usuario sabe a qué cursos e instituciones tiene acceso.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Autenticación | Login y gestión de sesión |
| Roles y permisos | Definición de roles y qué puede hacer cada uno |
| Asignación organizacional | Vinculación de usuarios a instituciones, áreas y cursos |

### **Decisiones de cada cliente**

- Los roles adicionales a coordinador y docente dependen de cada provincia
- El modelo de permisos sobre el documento de coordinación (quién edita, quién solo visualiza) es decisión de cada cliente

### **Decisiones técnicas**

- Un usuario puede tener **múltiples roles dentro de una misma organización**. Un docente puede ser profesor de dos materias y coordinador de un área — no hay restricción. Idealmente la experiencia no necesita "Escoger un rol" para el usuario.
- En el MVP, si un usuario trabaja en **dos instituciones distintas**, tiene dos cuentas separadas (un usuario por organización).
- El mecanismo de autenticación puede variar por provincia: mail + contraseña, cuentas institucionales (ej: Google Workspace del ministerio), u otros proveedores. Para el MVP, limitamos a mail + password
- Los permisos sobre el documento de coordinación (quién edita, quién solo visualiza) son **configurables por organización**. No se hardcodea que "el coordinador edita y el docente solo ve" porque hay provincias donde el docente también interviene.

### **Épicas relacionadas**

- **Onboarding** — Flujo post-autenticación para nuevos usuarios
- **Documento de coordinación** — Permisos de edición y visualización del documento
- **Planificación docente** — Acceso del docente a sus cursos y clases

---

## **2. Onboarding**

> Carga de datos iniciales y product tour para nuevos usuarios.
>

### **Problema**

Un usuario nuevo necesita completar su perfil y entender la plataforma antes de ser productivo. Sin un flujo guiado, el tiempo hasta el primer uso real es alto.

### **Objetivos**

- Capturar los datos necesarios del usuario al primer ingreso
- Guiar al usuario por las funcionalidades clave según su rol

### **Alcance MVP**

**Incluye:**

- Carga de datos del usuario (perfil, disciplinas, etc.)
- Product tour contextual según rol (coordinador o docente)

**No incluye:**

- Configuración de institución o estructura curricular → ver Integración
- Alta de usuarios → ver Roles y accesos

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Datos | Formulario de carga de información del usuario al primer ingreso |
| Product Tour | Recorrido guiado por la plataforma adaptado al rol del usuario |

### **Decisiones de cada cliente**

- Los datos requeridos en el onboarding pueden variar por provincia
- El contenido del product tour depende de las funcionalidades habilitadas para cada cliente

### **Decisiones técnicas**

- El onboarding se dispara **post-autenticación al primer ingreso**. Los datos de la institución y la estructura curricular ya están cargados vía Integración — el onboarding solo captura datos del usuario, no de la organización.
- Los datos requeridos se definen como **configuración por organización** (mismo JSON de config). Una provincia puede pedir disciplinas y experiencia docente; otra puede no pedir nada adicional al perfil básico.
- El product tour se adapta al **rol y a los feature flags activos** de la organización. Si una org no tiene habilitado contenido, el tour no muestra esa sección.

### **Épicas relacionadas**

- **Roles y accesos** — Define el rol que determina el flujo de onboarding
- **Integración** — Los datos de institución ya están cargados cuando el usuario hace onboarding

---

## **3. Integración**

> Importación de datos y configuración de la estructura curricular de cada provincia.
>

### **Problema**

Cada provincia tiene su propio diseño curricular, terminología y estructura de conocimientos. El sistema necesita incorporar estos datos como base para que todo lo que Alizia genere esté alineado con la realidad de cada jurisdicción.

### **Objetivos**

- Importar diseños curriculares, NAPs y fuentes oficiales de cada provincia
- Modelar la estructura curricular provincial (áreas, disciplinas, tópicos, sub-tópicos, actividades, tipos de contenidos a crear)
- Que estos datos sean la fuente de verdad para la generación de documentos y las respuestas del asistente IA

### **Alcance MVP**

**Incluye:**

- Importación de diseño curricular y fuentes oficiales
- Configuración de la estructura curricular (jerarquía de conocimientos, tópicos, etc.)
- Customizaciones por institución: Grillas horarias por institución, alta de docentes, etc.

**No incluye:**

- Gestión de usuarios e instituciones → ver Roles y accesos
- Carga de datos del docente o alumnos → ver Onboarding

*Aclaración: Sí se incluye el impacto y presencia de esos datos por institución, no la carga.*

### **Principios de diseño**

- **Provincial first:** Cada implementación respeta la estructura y terminología de la provincia.
- **Fuente de verdad única:** Los datos importados alimentan toda la plataforma.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Importación de datos | Carga de diseño curricular, NAPs, fuentes oficiales y grillas horarias |
| Estructura curricular | Modelado de la jerarquía de conocimientos según la provincia |

### **Decisiones del cliente**

- La granularidad de la estructura curricular varía (Neuquén usa conocimientos y saberes + categorías; otras provincias pueden diferir)

### **Decisiones técnicas**

- La configuración de cada organización se almacena en un **JSON de configuración** a nivel organización — mismo patrón que usamos en TUNI. Este JSON contiene: nombres de niveles de tópicos, profundidad máxima permitida, feature flags (clases coordinadas, subida de fuentes propias, etc.) y cualquier parametrización provincial.
- Los tópicos se modelan en una **tabla única auto-referencial** (foreign key que apunta a sí misma). La profundidad de cada nodo se **pre-computa y almacena** para evitar recursiones costosas en queries. Si se mueve un tópico, se re-computa solo la rama afectada.
- Cada organización define la **profundidad máxima** de su jerarquía de tópicos y los **nombres por nivel** (ej: nivel 1 = "Conocimientos y saberes", nivel 2 = "Núcleos problemáticos", nivel 3 = "Categorías"). El frontend renderiza dinámicamente según esta configuración.
- El setup inicial de cada provincia (estructura curricular, grillas, docentes) se realiza **manualmente por el equipo de implementación** — no hay backoffice de auto-servicio en el MVP. Esto es deliberado: lo que vendemos es que hacemos el setup del diseño curricular.
- Las migraciones de base de datos se hacen **incrementalmente** a medida que se desarrollan las funcionalidades — no se crea el DER completo upfront. Lección aprendida de TUNI: un schema gigante genera redefiniciones constantes y deuda de mantenimiento.
- Un área puede contener **una o más disciplinas**. Si la provincia no organiza por áreas, el coordinador crea un documento con las materias que correspondan. No se fuerza la existencia de áreas como concepto obligatorio.

### **Épicas relacionadas**

- **Documento de coordinación** — Consume la estructura curricular para generar documentos
- **Planificación docente** — Usa los datos importados como contexto para la planificación

---

## **4. Documento de coordinación**

> Creación asistida, edición colaborativa y gestión de estados de documentos de coordinación areal.
>

### **Problema**

Los equipos de docentes necesitan documentos complejos (como el itinerario del área) que deben alinear múltiples personas en el trabajo a realizarse en un plazo futuro dado. Este proceso es manual, lento y difícil de articular entre roles.

### **Objetivos**

- Generar una primera versión del documento, que siga los lineamientos correspondientes y las mejores prácticas
- Permitir la edición colaborativa entre roles (por ej: coordinadores y docentes)
- Asegurar que el documento resultante sea la base sobre la cual se planifica el clase a clase

### **Alcance MVP**

**Incluye:**

- Generación asistida del documento basado en un "topic", que incluya los elementos que la provincia quiera, por ej: Eje problemático y estrategia metodológica
- Cálculo automático de cantidad de clases por disciplina según grilla horaria
- Selección de "sub-topic", donde habrá múltiples disponibles para cada disciplina, que pueden repetirse o no en otras disciplinas (configurable según provincia)
- Generación de un cronograma de clases tentativo, con objetivo por clase para cada disciplina
- Edición directa y asistida por IA del documento
- Publicación del documento para que todos los roles puedan acceder

**No incluye:**

- Planificación del clase a clase detallado → ver Planificación docente
- Creación de recursos didácticos → ver Contenido

### **Principios de diseño**

- **Propuesta primero:** Alicia genera una primera versión; el equipo edita y valida.
- **Alineación vertical:** Todo lo que se planifica debe poder trazarse hasta los lineamientos provinciales.
- **Colaboración entre roles:** Coordinadores definen el marco, docentes validan y ajustan el cronograma.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Asistente de creación | Genera elementos del documento: eje problemático, estrategia metodológica y cronograma a partir del tema y documentación base (diseño curricular, etc.) |
| Editor de documento | Edición directa y asistida por IA de cada sección del documento |
| Flujo de estados | Gestión del ciclo de vida: borrador, en revisión, publicado |

### **Decisiones de cada cliente**

- Las estrategias metodológicas disponibles (proyecto, taller, ateneo, laboratorio). Además requieren validación con cada equipo pedagógico provincial
- El nivel de edición que tiene el docente sobre el documento del coordinador es decisión de cada provincia

### **Decisiones técnicas**

- El **período** del documento (cuatrimestre, bimestre, semestre) es un nombre libre con rango de fechas custom — no se asume una duración fija ni un calendario predeterminado. El coordinador define inicio y fin.
- Cada sección del documento se define mediante un **JSON Schema** que especifica: estructura del output, prompt de generación y opciones disponibles. Esto permite que cada provincia customice qué secciones tiene y cómo se generan, sin cambios de código.
- El documento tiene estados: **pendiente → borrador → publicado**.
- **Clases coordinadas entre docentes** es un feature flag por organización. Si dos materias coinciden en horario, las modificaciones de un docente son visibles para el otro. Esto se destaca como diferenciador clave del producto — ninguna plataforma del mercado lo ofrece.
- El cálculo de clases por disciplina se basa en la grilla horaria importada. El coordinador puede sobreescribir manualmente (±) para contemplar feriados, días institucionales o situaciones no previstas por el sistema.
- La selección de sub-topics para el documento respeta la **profundidad configurada por organización**. Si la org define que se seleccionan hasta nivel 2, solo se muestran esos niveles — se diseñarán UI/UX para cada caso y el sistema se adaptará dinámicamente.

### **Épicas relacionadas**

- **Integración** — Provee diseño curricular y estructura que alimentan la generación
- **Planificación docente** — Consume el documento como insumo para el clase a clase
- **Asistente IA** — Motor de generación y edición del documento

---

## **5. Planificación docente**

> Planificación del clase a clase por momentos, con asistencia de IA y feedback post-dictado.
>

### **Problema**

Los docentes planifican sus clases sin una conexión clara con lo acordado a nivel área, ni con lo que están dictando en clases los otros docentes. Re-planificar en base a lo que ocurrió en clases anteriores es algo que depende 100% del docente, de su energía y su memoria. Falta un sistema que alinee, proponga, recuerde y aprenda del día a día del aula.

### **Objetivos**

- Que cada clase planificada esté alineada con el documento de coordinación del área
- Que cada clase sea coherente con lo que se está dictando en otras disciplinas
- Generar una propuesta inicial de actividades por momento (apertura, desarrollo, cierre)
- Incorporar el feedback de clases anteriores para mejorar las propuestas futuras
- Permitir personalización del docente sin perder la alineación curricular

### **Alcance MVP**

**Incluye:**

- Visualización del cronograma de clases heredado del documento de coordinación
- Edición del objetivo de clase
- Selección de actividades por momento (apertura, desarrollo, cierre) con recomendaciones de IA
- Personalización: el docente puede anclar la clase a un recurso (canción, lectura, etc.) o a un comentario personalizado
- Generación de la propuesta detallada de clase
- Edición directa y asistida por IA de la propuesta
- Gestión de estados de las planificaciones (Pendiente/planificar, en proceso, publicada)
- Bitácora de cotejo: el docente reporta cómo fue la clase (soporta audio)
- Recolección de datos: Alicia pregunta activamente por información faltante o alumnos con un caso que merezca seguimiento
- Repropuesta: En base a la bitácora proponemos cambios a las siguientes clases ya planificadas (en proceso o publicadas) y tendremos eso en cuenta para la generación de las pendientes

**No incluye:**

¿Qué pasa si el docente no tiene internet al momento de grabar la bitácora? Escuelas
rurales argentinas pueden tener conectividad limitada. ¿Se sube después? ¿Funciona offline? Esto aplica también a la planificación en general.

- Informe de proceso (resumen de progreso del alumno por área) → horizonte
- Trayectorias de refuerzo personalizadas → horizonte
- Creación de recursos didácticos → ver Contenido

### **Principios de diseño**

- **Del área al aula:** La planificación individual nace del acuerdo colectivo.
- **IA que aprende del aula:** Las propuestas mejoran con el feedback real del docente.
- **Voz del docente:** La bitácora acepta audio libre, sin formato rígido.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Plan de clase | Selección de actividades por momento y generación de la propuesta detallada |
| Momentos didácticos | Configuración de los tipos de actividad por momento (apertura, desarrollo, cierre) |
| Incorporación de fuentes | Anclaje de la clase a un recurso o fuente específica |
| Edición del documento | Permitir que el docente edite directamente o a través de asistente el doc, gestión de su estado |
| Bitácora | Recolección del resultado de una clase, guiado y simple para el docente |
| Repropuesta | Sugerencias de cambios para clases futuras ya planificadas (Proactividad del sistema) |

### **Decisiones de cada cliente**

- Los tipos de actividad disponibles por momento se definen con cada equipo pedagógico provincial
- El formato y profundidad de la bitácora de cotejo requiere validación

### **Decisiones técnicas**

- Se asume **un docente por materia por curso**. Si excepcionalmente hay dos, opera first-come-first-serve: el primero en planificar escribe, el segundo ve los cambios. El documento de coordinación es la referencia compartida.
- La planificación debe poder **exportarse como PDF** para impresión y para integración con plataformas provinciales existentes donde los docentes reportan sus planificaciones. El formato de exportación será un template configurable.
- Los **momentos didácticos** (apertura, desarrollo, cierre) son fijos como estructura, pero los tipos de actividad dentro de cada momento son configurables por organización. Cada provincia define su catálogo de actividades.
- Cuando existe el feature flag de **clases coordinadas**, la planificación muestra un indicador de clase compartida y las modificaciones se sincronizan entre los docentes involucrados.
- ;cente graba lo que pasó y el sistema lo procesa. No se fuerza un formato rígido porque la adopción depende de que sea natural y rápido.

### **Épicas relacionadas**

- **Documento de coordinación** — Provee el cronograma y objetivos de clase
- **Contenido** — Recursos disponibles para incorporar en la planificación
- **Asistente IA** — Genera propuestas y procesa feedback

### **Notas**

- El informe de proceso y las trayectorias de refuerzo están en horizonte, pendientes de priorización

---

## **6. Asistente IA**

> Asistente de inteligencia artificial que genera, edita y aprende del uso de la plataforma.
>

### **Problema**

Los docentes y coordinadores necesitan producir documentos y planificaciones alineadas curricularmente, pero el proceso es lento y requiere expertise. Un asistente genérico no conoce el contexto provincial ni el historial del aula.

### **Objetivos**

- Generar primeras versiones de documentos y planificaciones basadas en el diseño curricular
- Permitir ediciones masivas o puntuales por instrucción natural del usuario
- Incorporar el historial de clases (bitácora) para mejorar las recomendaciones futuras
    - WhatsApp requiere una **capa de autenticación secundaria**: vincular número de teléfono con el usuario de la plataforma, y definir permisos de qué puede hacer desde el canal vs. qué requiere entrar a la plataforma.
- Mantener todas las generaciones alineadas con las fuentes oficiales de la provincia

### **/Alcance MVP’**

**Incluye:**

- Generación de documentos de coordinación
- Recomendación de actividades para cada momento de una clase
- Edición asistida: el usuario da instrucciones en lenguaje natural y Alicia modifica el documento
- Procesamiento de feedback post-clase (bitácora de cotejo) para ajustar futuras propuestas

### **Principios de diseño**

- **Curada y validada:** Las generaciones se basan en fuentes oficiales y criterio pedagógico colectivo.
- **Propuesta, no imposición:** Alicia propone; el docente decide.
- **Memoria del aula:** El feedback de clases anteriores informa las recomendaciones.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Modificación de contenido | Edición asistida de documentos y planificaciones por instrucción natural |
| Asistencia de uso y navegación | Ayuda contextual dentro de la plataforma |
| Customización por cliente | Adaptación del comportamiento del asistente según la provincia o institución |

### **Decisiones del cliente**

- La customización por cliente (tono, límites, fuentes permitidas) requiere definición por provincia

### **Decisiones técnicas**

- El asistente opera como una **LLM con tools** (function calling). Puede leer y modificar el documento que el usuario está editando, recomendar actividades, consultar la estructura curricular, y acceder a la bitácora de clases anteriores.
- Cada sección generada por IA se define con un **prompt + JSON Schema** por organización. El prompt incorpora variables contextuales (tópicos seleccionados, disciplina, grilla horaria) y el schema fuerza el formato del output. Esto permite customización por provincia sin cambios de código.
- El **planificador de clases** recibe todos los temas asignados al documento, las disciplinas involucradas y la cantidad de clases disponibles, y distribuye el contenido en el tiempo. Si hay mucho tiempo, lo esparce; si hay poco, lo compacta.
- La generación de secciones del documento de coordinación (eje problemático, estrategia metodológica, criterios de evaluación) arranca con un prompt simple. Se prevé que habrá **mucho prueba y error** — la sofisticación vendrá con el uso real, no con sobre-ingeniería anticipada.
- El asistente de chat permite al usuario pedir modificaciones en lenguaje natural ("cambiá la actividad del cierre por algo más dinámico") y el sistema modifica la sección correspondiente. Internamente es un chat con tools que opera sobre el documento activo.

### **Épicas relacionadas**

- **Documento de coordinación** — Alicia genera y edita el documento
- **Planificación docente** — Alicia propone actividades y procesa la bitácora
- **Contenido** — Alicia genera recursos didácticos

---

## **7. Dashboard**

> Vista consolidada del estado de documentos, cursos y notificaciones.
>

### **Problema**

Coordinadores y docentes no tienen un lugar único donde ver el estado de sus documentos, planificaciones y cursos. La información está dispersa y no hay visibilidad del progreso general.

### **Objetivos**

- Dar visibilidad rápida del estado de documentos de coordinación y planificaciones
- Centralizar el acceso a cursos asignados
- Notificar cambios relevantes (publicaciones, actualizaciones, plazos)

### **Alcance MVP**

**Incluye:**

- Vista de estado de documentos de coordinación (borrador, publicado, etc.)
- Vista de cursos asignados al usuario
- Sistema de notificaciones

**No incluye:**

- Métricas de uso o analytics del docente → horizonte
- Reportes de progreso de alumnos → horizonte

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Estado de documentos | Visualización del estado de documentos y planificaciones |
| Cursos | Listado y acceso a cursos asignados |
| Notificaciones | Alertas sobre publicaciones, cambios y plazos |

### **Decisiones de cada cliente**

- Qué información se muestra en el dashboard puede variar según el rol y la provincia

### **Decisiones técnicas**

- Lo que ve cada usuario en el dashboard depende de su **rol y la configuración de la organización**. Un coordinador ve el estado de sus documentos y los cursos del área; un docente ve sus planificaciones y las clases próximas.
- Las notificaciones cubren eventos clave: publicación de un documento de coordinación (el docente ya puede planificar), modificaciones en clases coordinadas, y plazos próximos. El alcance exacto de notificaciones se define con el primer cliente.

### **Épicas relacionadas**

- **Documento de coordinación** — Los documentos se visualizan en el dashboard
- **Planificación docente** — Las planificaciones se visualizan en el dashboard
- **Roles y accesos** — El rol define qué ve cada usuario

---

## **8. Contenido**

> Biblioteca de recursos didácticos y herramientas de creación basadas en fuentes oficiales.
>

### **Problema**

Los docentes necesitan recursos didácticos (fichas de cátedra, guías de lectura, imágenes, videos, etc.) adaptados a su contexto curricular. Crearlos desde cero es lento y recurrir a fuentes no curadas genera inconsistencias.

### **Objetivos**

- Permitir la creación de recursos a partir de fuentes oficiales validadas por el ministerio
- Ofrecer tipos de recurso predefinidos (ficha de cátedra, guía de lectura, entre otros) y también "Creación libre" donde el usuario define los lineamientos de lo que quiere
- Que los recursos creados sean utilizables dentro de las planificaciones de clase, para ser entregados antes, durante o pos clase

### **Alcance MVP**

**Incluye:**

- Creación de recursos a partir de fuentes oficiales provistas por el cliente
- Tipos de recurso predefinidos (ficha de cátedra, guía de lectura, más a definir con el ministerio)
- Edición directa y asistida por IA del recurso generado
- Permitir la exportación para impresión
- Disponibilizar para uso del recurso en clase (Planificación)

**No incluye:**

- Subida de fuentes propias del docente → decisión pendiente por provincia
- Biblioteca compartida entre docentes → horizonte

### **Principios de diseño**

- **Fuentes curadas:** Los recursos se generan desde fuentes oficiales, no desde internet abierto.
- **Listo para el aula:** El recurso generado debe ser usable directamente con los alumnos, incluyendo la información del aula.

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Biblioteca | Repositorio de recursos creados por el docente |
| Creación de contenido | Generación asistida de recursos a partir de fuentes y tipo seleccionado |

### **Decisiones de cada cliente**

- Los tipos de recurso disponibles se definen con cada equipo pedagógico provincial
- Si se permite o no que el docente suba fuentes propias es decisión de cada provincia

### **Decisiones técnicas**

- Cada organización **habilita qué tipos de recurso** tiene disponibles. No todos los clientes ven los mismos tipos — se activan según acuerdo con el equipo pedagógico provincial. Un tipo de recurso puede existir en el sistema y no estar habilitado para una organización.
- El concepto de **biblioteca** es central: los recursos creados se almacenan a nivel organización. Todos los docentes de la misma organización pueden explorar y reutilizar recursos existentes antes de generar uno nuevo. Esto reduce costos de generación y promueve consistencia.
- El filtro por materia opera como **restricción soft** (UX, no permisos). Un docente de matemáticas no ve recursos de ciencias naturales por default, pero a nivel permiso el acceso es por organización.
- Cada tipo de recurso tiene un **prompt y un JSON Schema** que define la estructura del output. Esto permite que una misma funcionalidad (ej: guía de lectura) genere resultados con formatos distintos según la provincia, modificando solo la configuración.
- Arrancar simple: en el MVP, si un segundo cliente necesita una variante de un tipo existente, se **duplica y adapta** en vez de sobre-ingenierizar un sistema de templates parametrizables. La genericidad se construye cuando aparezca el patrón real.

### **Épicas relacionadas**

- **Planificación docente** — Los recursos se pueden incorporar en las clases
- **Asistente IA** — Motor de generación y edición de contenido
- **Integración** — Provee las fuentes oficiales del diseño curricular

---

## **9. WhatsApp**

> Canal de WhatsApp para interacción con Alizia fuera de la plataforma.
>

### **Problema**

*Pendiente de definición con más contexto.*

### **Objetivos**

*Pendiente.*

### **Alcance MVP**

*Pendiente.*

### **Sub-épicas**

| Componente | Descripción |
| --- | --- |
| Comportamiento | Definición de qué puede hacer Alizia vía WhatsApp y sus límites |
| Conexión | Integración técnica con la API de WhatsApp |

### **Decisiones de cada cliente**

- Pendiente definir el alcance funcional del canal

### **Decisiones técnicas**

- El canal se piensa como **cross-producto** — la experiencia de integración con WhatsApp aplica tanto a Alizia como a TUNI y futuros productos. Aprender a operar bien el canal acá nos da un impulso para toda la empresa.
- El comportamiento inicial apunta a **efecto wow + casos básicos**: el docente pregunta algo rápido, recibe un resumen, o dispara un proceso que queda listo en la plataforma. No se busca replicar toda la funcionalidad del producto en WhatsApp.

### **Épicas relacionadas**

- **Asistente IA** — El comportamiento del canal se basa en el asistente IA

---

## **10. Cosmos**

> Pendiente de definición.
>

*Esta épica aún no tiene contenido definido.*

---

## **Mapa de dependencias entre épicas**

```
Roles y accesos ──→ Onboarding ──→ (usuario listo para usar la plataforma)
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

   Cosmos ── (pendiente)
```

---

## **Patrones transversales**

Estos patrones aparecen en múltiples épicas y merecen atención como decisiones de arquitectura:

| Patrón | Épicas que lo usan | Descripción |
| --- | --- | --- |
| JSON de configuración por org | Integración, Onboarding, Doc. coord., Contenido | Configuración provincial centralizada (feature flags, nombres de niveles, tipos habilitados) |
| Prompt + JSON Schema por sección | Doc. coord., Asistente IA, Contenido | Cada output generado por IA tiene su prompt y schema configurable por provincia |
| Feature flags por organización | Onboarding, Doc. coord., Planificación, Contenido | Funcionalidades que se activan/desactivan por cliente |
| Clases coordinadas | Doc. coord., Planificación | Diferenciador clave: sincronización entre docentes que comparten horario |
| Decisiones por provincia | Todas | Cada cliente (provincia) puede customizar comportamiento sin cambios de código |

---

> **Siguiente paso:** Revisá cada épica, dejá tus comentarios y dudas. Una vez validado este RFC, avanzamos con las especificaciones funcionales detalladas (spec.mdx) de cada épica.
>