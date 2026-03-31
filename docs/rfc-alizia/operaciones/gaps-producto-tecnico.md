# Gaps: Producto → Tecnico

> Features documentadas en producto (alizia-docs) sin cobertura tecnica en el RFC.
> Generado: 2026-03-30

## Resumen

17 gaps identificados en 8 epicas.

Las epicas con mejor cobertura son Setup (solo tecnica, sin contraparte producto), Integracion, Documento de coordinacion y Planificacion docente. Las epicas con mayor deuda son Cosmos (3 sub-features MVP sin HU), Dashboard (widget "Requiere atencion" sin HU dedicada) y las epicas horizonte (Aprendizaje adaptativo, Inclusion) que tienen especificaciones funcionales ricas en producto pero solo placeholders en el RFC.

---

## Gaps por epica

### Epica 1: Roles y accesos

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Permisos configurables por org sobre el documento de coordinacion (quien edita, quien solo visualiza) | Parcial | Media | Producto dice explicitamente que es "decision de cada cliente" y que no se hardcodea. El RFC menciona esto en decisiones tecnicas de la epica pero no hay HU ni tarea que implemente la configuracion de permisos por org. HU-1.3 define RequireRole estatico (coordinator/teacher), no permisos dinamicos. |

### Epica 2: Onboarding

Cobertura completa. El RFC tiene HU-2.1 a HU-2.4 que cubren: deteccion de primer ingreso, formulario de perfil, product tour y configuracion por org.

*(Sin gaps)*

### Epica 3: Integracion

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Importacion de diseno curricular y fuentes oficiales (documentos, NAPs) | Sin cobertura | Media | Producto menciona "Importacion de diseno curricular, NAPs, fuentes oficiales" como sub-epica. El RFC modela topics y la jerarquia, pero no hay HU que cubra la carga/importacion de las fuentes documentales que alimentan a la IA (PDFs, documentos curriculares). El seed cubre datos de estructura pero no fuentes textuales. Esto impacta directamente la calidad de generacion IA. |

### Epica 4: Documento de coordinacion

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Colaboracion entre roles (docente valida y ajusta el cronograma) | Parcial | Media | Producto indica "Colaboracion entre roles: Coordinadores definen el marco, docentes validan y ajustan el cronograma". El RFC cubre permisos de lectura para teacher pero no hay flujo tecnico de "el docente sugiere cambios al coordinador" ni mecanismo de feedback/aprobacion colaborativo. |

### Epica 5: Planificacion docente

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Recoleccion activa de datos por Alizia (preguntas sobre informacion faltante, alumnos con seguimiento) | Parcial | Media | Producto dice "Alizia pregunta activamente por informacion faltante o alumnos con un caso que merezca seguimiento". HU-5.6 cubre preguntas de seguimiento en la bitacora, pero la recoleccion activa como feature autonoma (Alizia que proactivamente pide datos fuera de bitacora) no tiene HU dedicada. |
| Personalizacion: anclar clase a recurso (cancion, lectura, etc.) o comentario personalizado | Parcial | Baja | Producto lo menciona como feature MVP. HU-5.3 cubre "fuentes" y "notas" en el modelo pero el flujo especifico de "anclar a un recurso externo" (URL, titulo de cancion, etc.) no esta detallado como campo/flujo tecnico explicito. |

### Epica 6: Asistente IA

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Procesamiento de feedback post-clase (bitacora) para ajustar futuras propuestas a nivel del motor IA | Parcial | Media | Producto lo lista como alcance MVP del asistente IA: "Procesamiento de feedback post-clase (bitacora de cotejo) para ajustar futuras propuestas". El RFC tiene HU-5.7 (Repropuesta) en Epica 5, pero dentro de Epica 6 no hay HU ni context builder que integre datos de bitacora como contexto para generaciones futuras. La conexion bidireccional bitacora→motor IA esta en las dependencias pero no materializada como tarea en Epica 6. |
| Recomendacion de actividades para cada momento de una clase | Parcial | Media | Producto lo lista como alcance MVP: "Recomendacion de actividades para cada momento de una clase". HU-5.4 genera propuesta de clase pero el flujo especifico de "recomendar actividades del catalogo de HU-3.6 segun contexto" no tiene usecase dedicado. La IA genera texto libre, no selecciona del catalogo. |

### Epica 7: Dashboard

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Widget "Requiere atencion" con heuristicas configurables por org | Parcial | Media | Producto define heuristicas especificas: "documentos sin publicar hace mas de N dias (default 7), docentes sin planificar a menos de N dias del inicio (default 14), documentos completos pero no publicados". HU-7.1 menciona el widget pero no hay tarea dedicada para el job/query de heuristicas configurables ni para almacenar los thresholds en organizations.config. |
| Vista de cursos asignados al usuario (como seccion del dashboard) | Parcial | Baja | Producto lo lista como feature MVP separada: "Vista de cursos asignados al usuario". HU-7.1 y HU-7.2 lo incluyen como widget pero no hay endpoint especifico de "mis cursos" — depende de logica existente. |

### Epica 8: Contenido y recursos

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Creacion libre (usuario define lineamientos de lo que quiere) | Sin cobertura | Media | Producto incluye "Creacion libre donde el usuario define los lineamientos de lo que quiere" como tipo de recurso. El RFC solo cubre tipos predefinidos con prompt+schema fijos. No hay HU ni flujo para que el usuario escriba un prompt/lineamiento libre y genere un recurso ad-hoc. |
| Disponibilizar recurso para uso en clase (vinculo recurso→planificacion) | Sin cobertura | Media | Producto dice "Que los recursos creados sean utilizables dentro de las planificaciones de clase, para ser entregados antes, durante o pos clase". No hay HU tecnica que vincule resources con lesson_plans ni endpoint para asociar un recurso a una clase planificada. |

### Epica 9: WhatsApp

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Epica completa | Sin cobertura | Baja | Tanto producto como RFC marcan esta epica como "pendiente de definicion". No hay gap real porque ambos lados estan alineados en que falta definicion. Se incluye por completitud. |

### Epica 10: Cosmos

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Feature flags por modulo (modules.planificacion, modules.contenido, etc.) | Sin cobertura | Alta | Producto lo lista como MVP. El RFC de Cosmos lo documenta en el schema de config pero lo marca como "Por agregar" — no hay HU ni tarea que implemente la logica de verificar feature flags en los endpoints/handlers. Sin esto, no se puede habilitar/deshabilitar modulos por cliente. |
| Identidad visual por organizacion (logo, colores, nombre de plataforma) | Sin cobertura | Alta | Producto lo lista como MVP: "Identidad visual por organizacion (logo, colores, nombre)". El RFC de Cosmos lo documenta en el schema (`visual_identity`) pero lo marca como "Por agregar". No hay HU, endpoint ni tarea para servir estos datos al frontend. |
| Nomenclatura configurable (terminos por provincia) | Sin cobertura | Alta | Producto lo lista como MVP: "Nomenclatura configurable (ej: itinerario del area vs documento de coordinacion)". El RFC de Cosmos lo documenta en el schema (`nomenclature`) pero lo marca como "Por agregar". Sin HU ni implementacion tecnica. |

### Epica 11: Aprendizaje adaptativo

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Epica completa (perfil de alumno, actividades, indicadores de progreso, adaptacion de propuestas) | Sin cobertura | Baja | Producto tiene un overview con sub-epicas definidas (actividades para alumnos, perfil de alumno, indicadores, adaptacion de propuestas, integracion con terceros). El RFC solo tiene un placeholder de 1 parrafo: "Pendiente de definicion tecnica". Severidad baja porque ambos lados marcan la epica como "Horizonte". |

### Epica 12: Inclusion

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Epica completa (valija, planificador, QR, fichas, feedback, 4 flujos detallados, modelo de datos, 9 bloques de RFs) | Sin cobertura | Baja | Producto tiene una especificacion funcional completa (spec.mdx) con modelo de datos, flujos de usuario, requerimientos funcionales detallados (RF01-RF09) y roadmap por fases. El RFC solo tiene un placeholder de 1 parrafo. Severidad baja porque es "Horizonte", pero la brecha de documentacion es la mas grande de todo el proyecto. |

---

## Gaps transversales (producto/modularizacion.mdx)

| Feature en producto | Estado en RFC | Severidad | Nota |
|---|---|---|---|
| Ecosistema de terceros (modulos 5, 6, 7... con requisitos de integracion) | Sin cobertura | Baja | Producto describe un ecosistema donde terceros pueden construir modulos que se integran con Alizia. No hay cobertura tecnica, pero es horizonte. |
| Composicion entre modulos (interacciones cruzadas tipo Planificacion+Contenido, Adaptativo+Inclusion) | Sin cobertura | Baja | Producto documenta una tabla de composiciones entre modulos. No hay patron tecnico definido para estas interacciones. |

---

## Resumen por severidad

| Severidad | Cantidad | Detalle |
|---|---|---|
| Alta | 3 | Cosmos: feature flags, identidad visual, nomenclatura — todos MVP sin HU tecnica |
| Media | 9 | Gaps parciales en Roles, Integracion, Doc. coordinacion, Planificacion, Asistente IA, Dashboard, Contenido |
| Baja | 5 | WhatsApp (ambos pendientes), Aprendizaje adaptativo, Inclusion (horizonte), Ecosistema terceros, Composicion modulos |

## Recomendacion de accion

1. **Prioridad inmediata:** Crear HUs o tareas en Cosmos (Epica 10) para feature flags, identidad visual y nomenclatura. Estos 3 items son MVP y estan documentados en el schema pero no tienen implementacion tecnica asignada.
2. **Prioridad alta:** Definir el flujo tecnico de "Creacion libre" de recursos (Epica 8) y el vinculo recurso→planificacion. Sin estos, dos features MVP de producto quedan sin backend.
3. **Prioridad media:** Cerrar los gaps parciales: permisos configurables por org (Epica 1), importacion de fuentes documentales (Epica 3), colaboracion docente-coordinador (Epica 4), heuristicas de "Requiere atencion" (Epica 7).
4. **Sin accion inmediata:** Aprendizaje adaptativo, Inclusion y WhatsApp estan alineados entre producto y RFC como epicas no priorizadas.
