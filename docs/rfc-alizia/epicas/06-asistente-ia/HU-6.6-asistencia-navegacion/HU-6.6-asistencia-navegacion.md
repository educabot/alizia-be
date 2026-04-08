# HU-6.6: Asistencia en navegacion

> Como usuario (docente o coordinador), necesito ayuda contextual de Alizia mientras navego la plataforma, para entender funcionalidades, recibir sugerencias y resolver dudas sin salir de la aplicacion.

**Fase:** Fase 4
**Prioridad:** Alta
**Estimacion:** --

---

## Contexto

Actualmente Alizia solo esta disponible dentro de la vista de edicion de documentos de coordinacion. Esta HU extiende su presencia a toda la plataforma como asistente contextual, capaz de guiar al usuario segun la pantalla en la que se encuentra y los datos que tiene visibles.

## Escenarios de uso

### Tour guiado
El usuario abre un documento de coordinacion por primera vez. Alizia ofrece un recorrido por las secciones principales: titulo, estrategia metodologica, plan de clases, y como interactuar con cada una.

### Sugerencia de proximos pasos
El usuario esta en una pantalla de planificacion vacia. Alizia sugiere por donde empezar: "Empeza seleccionando un nucleo problematico y las categorias que queres trabajar."

### Explicacion de conceptos
El usuario pregunta "que es un nucleo problematico?". Alizia responde con una explicacion basada en el curriculo de la provincia del usuario, no una definicion generica.

### Deteccion de inactividad
El usuario lleva tiempo en el wizard sin avanzar. Alizia detecta la inactividad y ofrece ayuda: "Necesitas una mano con este paso? Puedo explicarte que significa cada opcion."

### Resumen de documento publicado
El usuario abre un documento ya publicado. Alizia ofrece un resumen de las decisiones clave: cantidad de clases por disciplina, categorias cubiertas, y si hay categorias sin asignar.

## Criterios de aceptacion

- [ ] Existe un widget de chat liviano accesible desde todas las pantallas de la plataforma
- [ ] Alizia conoce la pantalla actual del usuario y los datos visibles en ella
- [ ] Las respuestas de navegacion son cortas y accionables (no explicaciones extensas)
- [ ] El comportamiento proactivo se limita a situaciones de trigger especificas (primera vez, inactividad, pantalla vacia)
- [ ] El usuario puede cerrar/minimizar el widget y su preferencia se persiste
- [ ] Las explicaciones de conceptos pedagogicos se adaptan a la provincia del usuario
- [ ] El widget no interfiere con la interaccion normal de la plataforma

## Comportamiento proactivo vs reactivo

### Proactivo (Alizia inicia)
- Primera visita a una pantalla clave
- Deteccion de inactividad prolongada en un wizard
- Pantalla vacia donde se espera contenido
- Documento recien publicado (ofrece resumen)

### Reactivo (usuario inicia)
- El usuario hace clic en el widget y escribe una pregunta
- El usuario pide explicacion de un termino
- El usuario solicita ayuda con una accion especifica

**Principio:** lo proactivo debe ser infrecuente y dismissable. Lo reactivo debe ser siempre disponible.

## Consideraciones tecnicas

- El widget de chat debe ser un componente global, no atado a una vista especifica
- Context-aware: Alizia recibe metadata de la pantalla actual (ruta, entidad visible, estado)
- Las respuestas de navegacion podrian usar un modelo mas simple/barato que el de generacion de documentos
- El historial de interacciones de navegacion es independiente del historial de chat de documentos
- Considerar cache de respuestas frecuentes para reducir llamadas a la API

## Consideraciones de producto

- El nivel de proactividad aceptable requiere user testing — empezar conservador
- El contenido de ayuda deberia ser hibrido: textos pre-escritos para lo comun, IA para lo especifico
- Usuarios que quieran deshabilitar el asistente deben poder hacerlo desde configuracion
- La integracion con un eventual flujo de onboarding (Epica 2) es deseable pero no bloqueante
- Las metricas clave son: tasa de uso del widget, preguntas mas frecuentes, tasa de dismiss de proactivos

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 6.6.1 | [Navigation context provider](./tareas/T-6.6.1-navigation-context-provider.md) | usecases/ai/navigation_context_builder.go | ⬜ |
| 6.6.2 | [Sistema de triggers proactivos](./tareas/T-6.6.2-proactive-triggers.md) | frontend + backend | ⬜ |
| 6.6.3 | [Widget de chat global](./tareas/T-6.6.3-chat-widget-global.md) | frontend | ⬜ |
| 6.6.4 | [Tests](./tareas/T-6.6.4-tests.md) | tests | ⬜ |

## Dependencias

- HU-6.1 completada (AI Provider)
- HU-6.2 completada (Prompts configurables)
- HU-6.5 deseable (Customizacion por organizacion — para adaptar explicaciones a la provincia)

## Preguntas abiertas

- Cuanta proactividad es aceptable? (requiere user testing)
- El contenido de ayuda debe ser pre-escrito o generado por IA? (enfoque hibrido probable)
- Como se maneja a usuarios que quieren deshabilitar el asistente?
- Se integra con el eventual flujo de onboarding (Epica 2)?
- Que modelo usar para las respuestas de navegacion? (costo vs calidad)
