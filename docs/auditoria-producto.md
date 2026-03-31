# Auditoría Producto vs Técnico — Material para reunión
**Fecha:** 2026-03-30
**Fuente producto:** `Aliciav3/Docs` (Mintlify, 13 épicas)
**Fuente técnica:** `alizia-be/docs/rfc-alizia/` (11 épicas, HUs, modelo de datos, arquitectura)
---
## A. Producto ya actualizó — Técnico debe alinear
Estos cambios los hizo producto. El equipo técnico debe actualizar su documentación para coincidir.
| # | Hallazgo | Acción para técnico |
|---|----------|---------------------|
| A1 | **Estados doc. coordinación** — Producto define: pendiente → en progreso → publicado. Técnico tiene: draft → published → archived | Cambiar el enum a `pending → in_progress → published`. Eliminar `archived` o dejarlo como horizonte. Actualizar `04-documento-coordinacion.md`, HU-4.5, y migraciones |
| A2 | **Estados planificación** — Producto define 3 estados (pendiente, en proceso, publicada). Técnico tiene 2 (pending, planned) | Agregar estado intermedio `in_progress`. Actualizar `05-planificacion-docente.md`, HU-5.5, y migraciones |
| A3 | **Auth MVP** — Producto ahora dice solo email+password. Técnico ya está alineado (JWT) | Sin acción técnica requerida. Solo confirmar que la abstracción multi-proveedor no es MVP |
| A4 | **Onboarding** — Producto dice MVP (product tour + datos de perfil, alta manual). Técnico dice Post-MVP | **Decisión requerida en reunión:** ¿Se mantiene onboarding como MVP con alcance reducido (solo product tour + perfil), o se posterga todo? |
| A5 | **Terminología "profesor"** — Debe ser "docente" en todos lados | Buscar y reemplazar "profesor" / "Prof." en la doc técnica. Archivos: `01-roles-accesos.md:55`, seeds con "Prof. García" |
| A6 | **Terminología "plan de clase"** — Debe ser "planificación docente" | Revisar uso en prosa castellana del RFC. Aceptable en código (`lesson_plan`), no en documentación funcional |
| A7 | **Terminología "la IA"** — Debe ser "Alizia" o "el Asistente IA" | Buscar y reemplazar en: `T-4.3.2`, `T-4.3.4`, `T-5.6.4`, `T-8.3.3` |
| A8 | **Typo "assistente"** — Directorio `06-assistente-ia/` tiene doble 's' | Renombrar a `06-asistente-ia/` |
---
## B. Técnico definió bien — Producto adoptó
Estos son features que el técnico desarrolló y que producto incorporó. No requieren acción adicional, solo awareness.
| # | Feature | Origen | Adoptado en producto |
|---|---------|--------|----------------------|
| B1 | Wizard de creación de 3 pasos | HU-4.2 | `documento-de-coordinacion/overview.mdx` |
| B2 | Asistencia de navegación con triggers proactivos | HU-6.6 + tareas | `asistente-ia/overview.mdx` |
| B3 | Widget "Requiere atención" con heurísticas | HU-7.1, T-7.1.1 | `dashboard/overview.mdx` |
| B4 | Progreso de planificación por materia | HU-7.1, T-7.1.1 | `dashboard/overview.mdx` |
| B5 | Actividades didácticas por momento | HU-3.6 + tareas | `integracion/overview.mdx` |
---
## C. Puntos abiertos — Requieren decisión conjunta en reunión
### C1. Dashboard: ¿MVP o Post-MVP?
| Producto dice | Técnico dice |
|---------------|--------------|
| MVP completo (documentos, cursos, notificaciones, requiere atención) | **Toda la épica es Post-MVP** (HU-7.1, 7.2, 7.3) |
**Pregunta:** Sin dashboard, ¿qué ve el usuario al loguearse? ¿Tiene sentido un MVP sin landing page?
**Propuesta:** Dashboard mínimo en MVP (vista de documentos + cursos), notificaciones y "Requiere atención" como post-MVP.
### C2. Bitácora y repropuesta: ¿MVP o Post-MVP?
| Producto dice | Técnico dice |
|---------------|--------------|
| MVP: bitácora de cotejo, recolección de datos, repropuesta | Post-MVP: HU-5.6 y HU-5.7 |
**Pregunta:** La bitácora es el diferenciador clave (IA que aprende del aula). ¿Se puede hacer un MVP sin ella? ¿O al menos incluir la bitácora básica sin la repropuesta?
**Propuesta:** Bitácora básica en MVP (el docente reporta cómo fue la clase). Repropuesta como post-MVP.
### C3. Customización por organización (IA): ¿MVP o Post-MVP?
| Producto dice | Técnico dice |
|---------------|--------------|
| Sub-épica activa: prompts configurables por org | Post-MVP completa (HU-6.5) |
**Pregunta:** Los prompts configurables (HU-6.2) ya están en MVP. ¿Es suficiente o se necesita la personalización completa (tono, límites, vocabulario)?
**Propuesta:** Prompts configurables por org en MVP (ya cubierto por HU-6.2). Tono, feature flags de IA y vocabulario custom como post-MVP.
### C4. Edición asistida por IA de planificación docente
| Producto dice | Técnico dice |
|---------------|--------------|
| MVP: "Edición directa y asistida por IA de la propuesta" | Solo generación (HU-5.4). No hay chat/edición IA para planificación |
**Pregunta:** Doc. coordinación tiene chat con Alizia (HU-4.6). ¿La planificación docente también necesita chat, o alcanza con edición manual + regeneración?
### C5. Exportación PDF de planificación
| Producto dice | Técnico dice |
|---------------|--------------|
| MVP: exportación como PDF | Mencionada en decisiones técnicas, pero no tiene HU ni tarea |
**Pregunta:** ¿Necesita HU propia o se resuelve como feature liviana al final del MVP?
### C6. Doc publicado no se puede editar — ¿Cómo se maneja?
El técnico define que un documento publicado no se puede editar y debe volver a draft. Producto no lo contemplaba.
**Pregunta:** ¿El coordinador puede "despublicar" para editar? ¿O se crea una nueva versión? ¿Qué pasa con las planificaciones ya hechas sobre un doc que se despublica?
### C7. Rol "admin": ¿Quién es?
El técnico usa extensivamente un rol `admin` con endpoints protegidos por `RequireRole("admin")`. Producto solo define coordinador y docente.
**Pregunta:** ¿El admin es un rol interno del equipo de implementación (no visible al usuario final)? ¿O es un rol de producto que necesita definición funcional?
### C8. Áreas: ¿opcionales o siempre presentes?
| Producto dice | Técnico dice |
|---------------|--------------|
| "No se fuerza la existencia de áreas como concepto obligatorio" | "Si la provincia no organiza por áreas, se crea un área genérica" |
**Pregunta:** ¿El coordinador ve "Área genérica" en la UI, o las áreas desaparecen por completo cuando la provincia no las usa?
### C9. Edición en documentos publicados
| Restricción: doc publicado no se puede editar | HU-4.5 | `documento-de-coordinacion/overview.mdx` |
---
## D. Contradicciones internas del técnico (para que dev-team resuelva)
Estos son problemas **dentro** de la doc técnica que el equipo de desarrollo debe resolver internamente:
| # | Contradicción | Archivos involucrados | Acción sugerida |
|---|---------------|----------------------|-----------------|
| D1 | **GORM vs sqlx** — La comparativa de arquitectura dice `~~GORM~~ → sqlx` pero todos los demás docs dicen GORM | `comparativa-arquitectura.md:472` vs `arquitectura.md:20`, `rfc-alizia.md:359` | Definir cuál es la decisión real y actualizar el doc incorrecto |
| D2 | **Cloud Functions vs Railway** — Comparativa de arquitectura y auth-service mencionan Cloud Functions; RFC y demás dicen Railway | `comparativa-arquitectura.md` (6+ lugares), `auth-service-futuro.md:85` vs `rfc-alizia.md:360` | Actualizar docs obsoletos que aún dicen Cloud Functions |
| D3 | **JWT HS256 vs JWKS** — team-ai-toolkit documenta HMAC-HS256; RFC y tareas mencionan JWKS | `team-ai-toolkit.md:63-111` vs `rfc-alizia.md:354`, `tareas.md:11` | Definir cuál mecanismo se usa y unificar |
| D4 | **Tabla `course_subjects` doble definición** — HU-1.4 con `organization_id` + UNIQUE; HU-3.4 con `start_date/end_date`, sin UNIQUE | `HU-1.4` vs `HU-3.4` | Unificar en un solo schema antes de implementar |
| D5 | **"disciplinas" vs "materias"** — Producto usa "disciplinas"; técnico usa "materias" y `subjects` | Múltiples archivos | Agregar al glosario. Propuesta: "materias" en prosa, `subjects` en código |
---
## E. Gaps pendientes (para planificar)
### E1. Features de producto sin cobertura técnica
| Feature | Épica | Estado |
|---------|-------|--------|
| Procesamiento de bitácora para recomendaciones IA | Asistente IA | Sin context builder ni tools. Depende de decisión C2 |
| "Clases coordinadas" — sincronización real entre docentes | Doc. coordinación | Solo hay marcador visual `is_shared`, sin sincronización |
| Identidad visual por org (logo, colores, nombre) | Cosmos | No aparece en ningún config JSONB |
| Feature flags genéricos por módulo | Cosmos | Solo flags de IA, no para módulos completos |
| Vinculación recurso → planificación | Contenido | En criterio de aceptación sin tarea dedicada |
### E2. Épicas sin presencia técnica
| Épica | Estado producto | Nota |
|-------|----------------|------|
| **Aprendizaje Adaptativo** | Overview definido, MVP pendiente | Horizonte. Agregar al índice técnico como épica futura |
| **Inclusión** | Spec completa (~480 líneas) | Horizonte. Agregar al índice para no olvidar |
### E3. Cosmos fragmentado
Producto define Cosmos como componente core. El técnico tiene un placeholder vacío, pero la funcionalidad **ya está implementada parcialmente** en HU-3.1 (organizations.config), HU-6.5 (personalización IA), T-3.1.5 (seed).
**Acción sugerida:** Documentar que `organizations.config` ES la implementación MVP de Cosmos.
### E4. Biblioteca compartida: discrepancia de alcance
Producto la excluye del MVP ("horizonte"). Técnico la implementa de facto en HU-8.4: "El docente puede ver todos los recursos de su organización".
**Pregunta:** ¿La dejamos como está o la excluimos del MVP?
### E5. "Fuentes curadas, no internet abierto" sin enforcement
El tipo "Creación libre" no tiene fuente anclada y acepta instrucciones arbitrarias. ¿Se necesita restricción técnica o alcanza con el diseño del prompt?
---
## F. Lo que está bien alineado
- Tabla auto-referencial de tópicos (jerarquía dinámica con profundidad pre-computada)
- Prompt + JSON Schema por sección (configuración sin código)
- Migraciones incrementales de BD
- JSON por organización (patrón Cosmos)
- Railway como deploy (monolito modular)
- WhatsApp: ambos coinciden en que está pendiente de definición
- Setup Infraestructura: épica puramente técnica, correctamente ausente de producto
- Contenido: alcance general alineado (MVP Fase 6, tipos configurables, generación por IA)