# HU-1.3: Middleware de autorización

> Como coordinador, necesito que solo los usuarios con el rol adecuado puedan acceder a ciertos endpoints, para garantizar la seguridad de los datos.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta (bloqueante para rutas protegidas por rol)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Middleware `RequireRole(roles...)` rechaza requests si el usuario no tiene al menos uno de los roles requeridos
- [ ] Interceptor chain funcional: Auth → Tenant → RequireRole → Handler
- [ ] Request con rol insuficiente → 403 `forbidden`
- [ ] Request sin claims (token válido pero sin roles) → 403 `forbidden`
- [ ] Roles se leen de `tokens.GetClaims(ctx).Roles`
- [ ] Error responses unificadas con el formato estándar (`{"error": "..."}`)
- [ ] Tests cubren combinaciones de roles permitidos y denegados

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.3.1 | [Implementar RequireRole middleware](./tareas/T-1.3.1-require-role.md) | internal/middleware/ | ⬜ |
| 1.3.2 | [Interceptor chain: Auth → Tenant → Role](./tareas/T-1.3.2-interceptor-chain.md) | cmd/routes.go | ⬜ |
| 1.3.3 | [Error handling unificado para auth/authz](./tareas/T-1.3.3-error-handling.md) | internal/middleware/ | ⬜ |
| 1.3.4 | [Tests de autorización por rol](./tareas/T-1.3.4-tests-authz.md) | internal/middleware/*_test.go | ⬜ |

## Dependencias

- HU-1.1 completada (JWT middleware inyectando claims en context)
- HU-1.2 completada (roles definidos como enum `member_role`)

## Test cases

- Request con rol `coordinator` a endpoint que requiere `coordinator` → 200
- Request con rol `teacher` a endpoint que requiere `coordinator` → 403 `forbidden`
- Request con roles `[teacher, coordinator]` a endpoint que requiere `coordinator` → 200
- Request sin claims en context → 403 `forbidden`
- Request a endpoint sin RequireRole → pasa (solo auth + tenant)
