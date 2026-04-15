# Auditoría de Proyecto — Alizia-BE — 2026-04-15

**Fecha:** 2026-04-15
**Stack:** Go 1.26, Gin, GORM, PostgreSQL 16, JWT, Railway
**Estado general:** NECESITA ATENCIÓN

## RESUMEN EJECUTIVO

El proyecto tiene una arquitectura Clean Architecture bien implementada con patrones consistentes en usecases, buena cobertura de tests (92.6%), y separación de capas correcta sin dependencias circulares. Sin embargo, se identifican **bugs de multi-tenancy en el repositorio de usuarios**, **falta de paginación en todos los endpoints list**, **validaciones incompletas en varios usecases**, y **código muerto del módulo coordination** que nunca se implementó. La seguridad general es buena (JWT, RBAC, SQL parametrizado) pero necesita rate limiting en login y limpieza de tokens de test del `.env.example`.

## SCORECARD

| Dimensión | Estado | Nota |
|-----------|--------|------|
| Arquitectura | OK | Clean Architecture correcta, sin circular deps |
| Calidad de código | WARN | Bugs multi-tenancy, validaciones faltantes, código muerto |
| Testing | OK | 92.6% usecases con tests, buena calidad |
| Seguridad | WARN | Sin rate limiting, tokens de test en .env.example |
| API/Contratos | WARN | Sin paginación, contenedores vacíos |
| Dependencias | OK | Deps actualizadas, CI con lint+tests |
| Documentación | OK | CLAUDE.md y RFC existen |

---

## HALLAZGOS CRÍTICOS

### C1. Multi-tenancy roto en UpdateProfileData y CompleteOnboarding

**Archivo:** `src/repositories/admin/user.go:93-108`

Las queries `CompleteOnboarding` y `UpdateProfileData` filtran solo por `userID` sin incluir `organization_id`. Un usuario autenticado en org A podría modificar datos de un usuario en org B si conoce su ID.

```go
// ACTUAL (vulnerable):
Where("id = ?", userID)

// CORRECTO:
Where("id = ? AND organization_id = ?", userID, orgID)
```

**Impacto:** Cross-tenant data modification.
**Fix:** Agregar `orgID` como parámetro a ambos métodos y actualizar la interface `UserProvider` + usecases que los llaman.

### C2. Sin paginación en endpoints list

**Archivos afectados:** Todos los repositorios con `List*` — `area.go`, `course.go`, `subject.go`, `activity_template.go`, `topic.go`

Ningún endpoint list tiene LIMIT. Con datos suficientes, un request puede retornar miles de registros causando OOM o timeouts.

```go
// ACTUAL:
r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&areas).Error

// CORRECTO:
r.db.WithContext(ctx).Where("organization_id = ?", orgID).Limit(100).Offset(offset).Find(&areas).Error
```

**Impacto:** Degradación de performance, potencial DoS.

### C3. Sin rate limiting en POST /auth/login

Se removió rate limiting en commit `9a06a45` por riesgo de memory leak. No se reimplementó.

**Impacto:** Vulnerable a ataques de fuerza bruta.
**Fix:** Implementar rate limiting con Redis o usar middleware como `gin-contrib/ratelimit`.

---

## HALLAZGOS IMPORTANTES

### I1. Falta validación de DayOfWeek en CreateTimeSlot

**Archivo:** `src/core/usecases/admin/create_time_slot.go:23-42`

El campo `DayOfWeek` acepta cualquier int (incluidos negativos y >6). Debería validar rango 0-6.

```go
if r.DayOfWeek < 0 || r.DayOfWeek > 6 {
    return fmt.Errorf("%w: day_of_week must be between 0 and 6", providers.ErrValidation)
}
```

### I2. strconv.ParseInt retorna 500 en vez de 400

**Archivos:** `src/entrypoints/onboarding.go`, `src/entrypoints/admin.go`, `src/entrypoints/courses.go`

Cuando `middleware.UserID(req)` retorna "" o un path param `:id` no es numérico, `strconv.ParseInt` falla y `rest.HandleError(err)` lo mapea a 500 (Internal Server Error) cuando debería ser 400 (Bad Request).

**Fix:** Wrappear el error con `providers.ErrValidation`:
```go
userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
if err != nil {
    return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
}
```

### I3. Falta validación de formato StartTime/EndTime

**Archivo:** `src/core/usecases/admin/create_time_slot.go:30-34`

Solo valida que no sea vacío, no que sea formato HH:MM válido ni que StartTime < EndTime.

### I4. JSON unmarshal silencioso en onboarding

**Archivos:** `src/core/usecases/onboarding/save_profile.go:76`, `get_tour_steps.go:121,131`, `get_config.go:54`

Si `org.Config` tiene JSON malformado, `json.Unmarshal` falla silenciosamente y retorna valores por defecto, enmascarando problemas de integridad.

### I5. Módulo coordination: código muerto

**Archivos afectados:**
- `src/core/entities/coordination.go` — 6 entities sin tablas en BD
- `src/core/providers/coordination.go` — Interface de 16 métodos
- `src/repositories/coordination/repository.go` — Stub vacío (solo constructor)
- `src/core/usecases/coordination/create_document.go` — Usecase completo pero nunca instanciado
- `src/core/usecases/coordination/get_document.go` — Usecase completo pero nunca instanciado
- `src/entrypoints/coordination.go` — Container vacío `struct{}`

Todo este código es no-funcional: no hay migraciones, el repositorio no implementa la interface, los usecases no están en DI, y no hay endpoints HTTP. Igualmente con `teaching.go` y `resources.go`.

### I6. Test faltante: list_subjects

**Archivo faltante:** `src/core/usecases/admin/list_subjects_test.go`

Es el único usecase de admin sin test.

### I7. Tokens JWT de larga duración en .env.example

**Archivo:** `.env.example:14-21`

Tokens de test con expiración en año 2036. Si alguien los usa en un ambiente real con el mismo JWT_SECRET, son válidos.

**Fix:** Reemplazar con tokens de expiración corta o eliminarlos y documentar cómo generarlos.

---

## OBSERVACIONES

### O1. Org config keys hardcodeadas como strings

Las keys de configuración de organización (`shared_classes_enabled`, `topic_max_levels`, `features`, `coord_doc_sections`) están dispersas como strings literales en los usecases. Moverlas a constantes mejoraría mantenibilidad.

### O2. RemoveCoordinator sin test de validación

`src/core/usecases/admin/remove_coordinator_test.go` no tiene subtests para validaciones (patrón que sí siguen todos los demás tests).

### O3. Dockerfile sin non-root user

El container corre como root. Agregar un usuario no-root mejoraría seguridad en producción.

### O4. CI sin security scanning

No hay `govulncheck`, `gosec`, ni dependency scanning en el pipeline. Agregar al menos `govulncheck` detectaría vulnerabilidades conocidas.

### O5. Linters faltantes en .golangci.yml

No tiene `gosec` (seguridad) ni `errorlint` (mejor error handling). Son los dos más útiles para agregar.

### O6. ListByCourse queries sin org_id

`src/repositories/admin/student.go:27`, `course_subject.go:35`, `timeslot.go:27` — métodos `ListByCourse` aceptan `courseID` directo sin validar que el curso pertenezca a la org. Mitigado porque los usecases validan antes, pero riesgo si se reutilizan repos directamente.

---

## LO QUE ESTÁ BIEN

1. **Clean Architecture impecable** — 0 violaciones de capas, 0 circular dependencies. Usecases nunca importan infraestructura.
2. **Patrón de usecases 100% consistente** — Los 27 usecases siguen el mismo patrón `Request.Validate()` + `Execute()` sin excepciones.
3. **Tests de alta calidad** — 25/27 usecases testeados. Tests verifican comportamiento real (idempotencia, deduplicación, feature flags), no implementación.
4. **Mocks bien diseñados** — 11/11 providers activos mockeados correctamente con testify.
5. **Seguridad de auth sólida** — JWT con validación completa, roles RBAC en 3 niveles (any, coordinator, admin), extracción segura de OrgID desde Audience claim, contraseñas con argon2id.
6. **SQL 100% parametrizado** — Todas las queries GORM usan placeholders. 0 riesgo de SQL injection.
7. **Error handling centralizado** — `rest.HandleError()` mapea errores de dominio a HTTP sin exponer internals. No revela si usuario existe o password incorrecto en login.
8. **Multi-tenancy correcto en 95% de queries** — Casi todas las queries filtran por `organization_id`.
9. **DI manual claro** — Bootstrap en `cmd/` es explícito y trazable.
10. **CI funcional** — Tests con race detector + golangci-lint + coverage en cada PR.

---

## PLAN DE ACCIÓN RECOMENDADO

### Prioridad 1 — Bugs (esta semana)
1. Fix multi-tenancy en `user.go`: agregar `organization_id` a `CompleteOnboarding` y `UpdateProfileData`
2. Fix strconv.ParseInt: wrappear errores de parsing con `ErrValidation` para retornar 400
3. Agregar validación `DayOfWeek` rango 0-6 en `create_time_slot.go`
4. Agregar validación formato `StartTime`/`EndTime` (HH:MM) y que start < end

### Prioridad 2 — Seguridad (próximo sprint)
5. Implementar rate limiting en POST /auth/login
6. Reemplazar tokens de test en .env.example
7. Agregar `gosec` a `.golangci.yml` y CI
8. Agregar non-root user al Dockerfile

### Prioridad 3 — API robustez (backlog)
9. Implementar paginación en todos los endpoints list (limit/offset, max 100)
10. Agregar test para `list_subjects`
11. Mejorar test de `remove_coordinator` con subtests de validación

### Prioridad 4 — Limpieza (backlog)
12. Decidir sobre módulo coordination: limpiar código muerto o documentar como WIP
13. Extraer org config keys a constantes
14. Agregar logging a JSON unmarshal silenciosos en onboarding

---

## MÉTRICAS

| Métrica | Valor |
|---------|-------|
| Archivos Go (producción) | ~85 |
| Usecases | 27 |
| Usecases testeados | 25/27 (92.6%) |
| Endpoints HTTP | 27 |
| Handlers testeados | 0/27 (0%) |
| Middleware testeado | 4/4 (100%) |
| Repos testeados | 1/12 (auth) |
| Hallazgos críticos | 3 |
| Hallazgos importantes | 7 |
| Observaciones | 6 |
