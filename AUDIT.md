# Auditoría de Proyecto — Alizia BE — 2026-04-08

**Stack:** Go 1.26.1 | GORM | PostgreSQL 16 | Gin | JWT | Railway
**Estado general:** SALUDABLE

## RESUMEN EJECUTIVO

El proyecto sigue Clean Architecture de forma ejemplar: los usecases nunca importan infraestructura, todas las dependencias fluyen en una sola dirección, y el DI manual en `cmd/` es claro y trazable. La seguridad es sólida con aislamiento multi-tenant verificado en 3 capas (middleware → usecase → repository). Los principales gaps son: tests faltantes en handlers/repositories, código stub sin implementar (coordination/teaching/resources), y archivos utilitarios vacíos. No hay vulnerabilidades críticas.

## SCORECARD

| Dimensión | Estado | Nota |
|-----------|--------|------|
| Arquitectura | **OK** | Clean Architecture impecable, 0 violaciones |
| Calidad de código | **WARN** | Stubs vacíos, archivos sin usar, error silenciado en app.go |
| Testing | **WARN** | Usecases 80% cubiertos, handlers y repos 0% |
| Seguridad | **OK** | SQL injection safe, tenant isolation excelente, JWT en todas las rutas |
| API/Contratos | **OK** | Endpoints consistentes, validación en todos los requests |
| Dependencias | **OK** | Todas actualizadas, go.sum commiteado |
| Documentación | **OK** | CLAUDE.md completo, RFC docs, README útil |

## HALLAZGOS CRITICOS

Ninguno. No hay vulnerabilidades de seguridad, data loss, o crashes inminentes.

## HALLAZGOS IMPORTANTES

### 1. Error silenciado en `cmd/app.go:62-66`

```go
func (a *App) Close() {
    sqlDB, _ := a.db.DB()  // error ignorado
    err := sqlDB.Close()
    if err != nil {
        return  // no loguea el error
    }
}
```

**Fix:** Loguear ambos errores:
```go
func (a *App) Close() {
    sqlDB, err := a.db.DB()
    if err != nil {
        log.Printf("error getting sql.DB: %v", err)
        return
    }
    if err := sqlDB.Close(); err != nil {
        log.Printf("error closing database: %v", err)
    }
}
```

### 2. Coordination usecases sin tests (2/10 usecases sin cubrir)

- `src/core/usecases/coordination/create_document.go` — 0 tests
- `src/core/usecases/coordination/get_document.go` — 0 tests
- Falta `MockCoordinationProvider` en `src/mocks/providers/`

### 3. Handlers sin tests (0% coverage)

Ningún handler en `src/entrypoints/` tiene tests HTTP. Los 8 handlers (admin + onboarding) no se testean end-to-end.

### 4. Repositories sin tests (0% coverage)

Los 4 repos implementados (`user.go`, `area.go`, `area_coordinator.go`, `organization.go`) no tienen tests.

### 5. Interfaces definidas pero nunca implementadas

En `src/core/providers/`:
- `SubjectProvider` (admin.go:41)
- `TopicProvider` (admin.go:47)
- `CourseProvider` (admin.go:52)
- `TimeSlotProvider` (admin.go:57)
- `TeachingProvider` (teaching.go:11)
- `ResourceProvider` (resources.go:11)
- `FontProvider` (resources.go:18)
- `ResourceTypeProvider` (resources.go:22)
- `AIClient` (ai.go:27)

### 6. Archivos y structs vacíos (stubs)

| Archivo | Tipo |
|---------|------|
| `src/entrypoints/coordination.go` | Container struct vacío |
| `src/entrypoints/teaching.go` | Container struct vacío |
| `src/entrypoints/resources.go` | Container struct vacío |
| `src/repositories/coordination/repository.go` | Solo constructor, 0 métodos |
| `src/repositories/teaching/repository.go` | Solo constructor, 0 métodos |
| `src/repositories/resources/repository.go` | Solo constructor, 0 métodos |
| `src/utils/json.go` | Archivo vacío |
| `src/utils/slices.go` | Archivo vacío |

## OBSERVACIONES

### Formato de timestamp hardcodeado
- `src/core/usecases/onboarding/get_status.go:58` usa `"2006-01-02T15:04:05Z07:00"` inline
- Extraer a constante: `const RFC3339Format = time.RFC3339`

### Duplicación de parseo de userID en handlers
- `strconv.ParseInt(middleware.UserID(req), 10, 64)` se repite 5 veces en `src/entrypoints/onboarding.go`
- Se podría extraer a helper en middleware, pero 3 líneas repetidas no justifican abstracción por ahora

### Tokens de test en `.env.example`
- Lines 14-21 tienen JWTs válidos por 10 años. Son útiles para desarrollo local pero técnicamente son secretos committed. Aceptable para dev, pero podrían moverse a un archivo `TESTING.md`.

## LO QUE ESTA BIEN

1. **Clean Architecture perfecta** — 0 violaciones. Usecases solo importan entities + providers (interfaces). Nunca infraestructura.
2. **Patrón de usecase 100% consistente** — Los 10 usecases siguen: Interface + impl + NewXxx + Execute + Request.Validate()
3. **Error handling ejemplar** — Sentinel errors (`ErrValidation`, `ErrNotFound`, `ErrConflict`) con wrapping contextual via `fmt.Errorf("%w: ...")`
4. **Multi-tenancy en 3 capas** — Middleware extrae org_id del JWT, usecases validan, repositories filtran por organization_id
5. **SQL injection safe** — 100% queries parametrizadas via GORM, 0 concatenación de strings
6. **DI manual claro** — Bootstrap en `cmd/` es explícito: repos → usecases → handlers → routes
7. **Middleware composition** — Auth → Tenant → Role, bien ordenado y testeado
8. **Repos con error mapping** — Traducen errores de GORM/PostgreSQL a errores de dominio
9. **Mocks bien estructurados** — 4/5 interfaces mockeadas, tests usan testify/mock correctamente
10. **Config segura** — `MustEnv()` para vars requeridas (fail-fast), 0 hardcoded values

## PLAN DE ACCION RECOMENDADO

### Prioridad 1 — Fix rápido (1 hora)
- [ ] Fix error handling en `cmd/app.go:62-66`
- [ ] Eliminar archivos vacíos (`src/utils/json.go`, `src/utils/slices.go`)
- [ ] Reemplazar timestamp hardcodeado con `time.RFC3339` en `get_status.go:58`

### Prioridad 2 — Tests faltantes (medio día)
- [ ] Crear `MockCoordinationProvider` en `src/mocks/providers/`
- [ ] Tests para `create_document.go` y `get_document.go`
- [ ] Tests para handlers HTTP (admin + onboarding) con httptest

### Prioridad 3 — Limpieza (cuando se trabaje en esas épicas)
- [ ] Implementar o eliminar interfaces no usadas (Subject, Topic, Course, etc.)
- [ ] Implementar o eliminar container/repo stubs (coordination, teaching, resources)

### No hacer (over-engineering)
- No agregar tests de repository con DB real (el mock testing en usecases es suficiente para esta etapa)
- No abstraer el parseo de userID (3 líneas repetidas no justifican abstracción)
- No agregar rate limiting (no hay endpoints públicos sin auth)

## METRICAS

| Métrica | Valor |
|---------|-------|
| Archivos Go | 72 |
| LOC (aprox) | 2,305 |
| Tests | 12 archivos, ~1,464 LOC |
| Usecases testeados | 8/10 (80%) |
| Handlers testeados | 0/8 (0%) |
| Repos testeados | 0/7 (0%) |
| Middleware testeado | 4/4 (100%) |
| Dependencias | Todas actualizadas |
| Vulnerabilidades | 1 low (Dependabot #61) |
