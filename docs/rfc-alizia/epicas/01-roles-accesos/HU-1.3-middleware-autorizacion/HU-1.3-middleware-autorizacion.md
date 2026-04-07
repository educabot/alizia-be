# HU-1.3: Middleware de autorización

> Como coordinador, necesito que solo los usuarios con el rol adecuado puedan acceder a ciertos endpoints, para garantizar la seguridad de los datos.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta (bloqueante para rutas protegidas por rol)
**Estimación:** —

---

## Criterios de aceptación

- [x] Middleware `RequireRole(roles...)` rechaza requests si el usuario no tiene al menos uno de los roles requeridos
- [x] Interceptor chain funcional: Auth → Tenant → RequireRole → Handler
- [x] Request con rol insuficiente → 403 `forbidden`
- [x] Request sin claims (token válido pero sin roles) → 403 `forbidden`
- [x] Roles se leen de `tokens.GetClaims(ctx).Roles`
- [x] Error responses unificadas con el formato estándar (`{"error": "..."}`)
- [x] Tests cubren combinaciones de roles permitidos y denegados

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.3.1 | [Implementar RequireRole middleware](./tareas/T-1.3.1-require-role.md) | src/entrypoints/middleware/require_role.go | ✅ |
| 1.3.2 | [Interceptor chain: Auth → Tenant → Role](./tareas/T-1.3.2-interceptor-chain.md) | src/app/web/mapping.go | ✅ |
| 1.3.3 | [Error handling unificado para auth/authz](./tareas/T-1.3.3-error-handling.md) | team-ai-toolkit (web.Err) | ✅ |
| 1.3.4 | [Tests de autorización por rol](./tareas/T-1.3.4-tests-authz.md) | src/entrypoints/middleware/*_test.go | ✅ |

## Dependencias

- HU-1.1 completada (JWT middleware inyectando claims en context)
- HU-1.2 completada (roles definidos como enum `member_role`)

## Test cases

- Request con rol `coordinator` a endpoint que requiere `coordinator` → 200
- Request con rol `teacher` a endpoint que requiere `coordinator` → 403 `forbidden`
- Request con roles `[teacher, coordinator]` a endpoint que requiere `coordinator` → 200
- Request sin claims en context → 403 `forbidden`
- Request a endpoint sin RequireRole → pasa (solo auth + tenant)

## Notas

> **Nota:** En MVP los permisos son estaticos: coordinador edita documentos, docente solo lee. Si una organizacion necesita permisos distintos, se puede extender el middleware para leer una configuracion de permisos desde `organizations.config`. Esto no requiere HU dedicada — se implementa como extension del middleware existente cuando surja la necesidad.
