# HU-1.1: Autenticación JWT

> Como usuario, necesito autenticarme con email y contraseña para acceder a la plataforma con mis permisos correspondientes.

**Fase:** 1 — Setup
**Prioridad:** Alta (bloqueante para rutas protegidas)
**Estimación:** —

---

## Criterios de aceptación

- [x] JWT auth configurado con HS256 shared secret para staging
- [x] JWT middleware valida tokens via HS256 shared secret (team-ai-toolkit/tokens)
- [x] Claims extraídos del JWT: user_id, org_id (via Audience), roles, email, name
- [x] Tenant middleware inyecta org_id en el contexto (via JWT Audience claim)
- [x] Request sin token → 401 `unauthorized`
- [x] Request con token inválido → 401 `unauthorized`
- [ ] Request con token de otra org → datos filtrados por org_id
- [ ] Refresh token rotation configurada (30 días abs, 7 días inactividad)
- [x] Endpoint POST /auth/logout invalida sesión
- [x] CORS configurado para orígenes permitidos via env var (boot.NewRouter)
- [x] Rate limiting en endpoints de IA (generate, chat)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.1.1 | [Configurar JWT auth](./tareas/T-1.1.1-configurar-jwt.md) | — | ✅ |
| 1.1.2 | [Integrar JWT middleware (HS256)](./tareas/T-1.1.2-jwt-middleware.md) | src/app/web/mapping.go | ✅ |
| 1.1.3 | [Integrar tenant middleware](./tareas/T-1.1.3-tenant-middleware.md) | src/entrypoints/middleware/tenant.go | ✅ |
| 1.1.4 | [Config: JWT auth env vars](./tareas/T-1.1.4-config-jwt.md) | config/config.go | ✅ |
| 1.1.5 | [Tests de autenticación](./tareas/T-1.1.5-tests-auth.md) | src/entrypoints/middleware/*_test.go | ✅ |
| 1.1.6 | [Refresh, logout y CORS](./tareas/T-1.1.6-refresh-logout-cors.md) | src/entrypoints/auth.go, middleware/ratelimit.go | ✅ |

## Dependencias

- Épica 0 completada (/health respondiendo)
- JWT auth configurado (HS256 shared secret)
- team-ai-toolkit/tokens funcional

## Test cases

- 1.1: Request sin Authorization header → 401 `missing_token`
- 1.3: Request con JWT inválido → 401 `invalid_token`
- 1.4: Request con JWT válido → 200 + claims en context
- 1.5: Request con JWT de otra org → datos filtrados por org_id
- 1.6: POST /auth/logout → 200 con logout_url válida
- 1.7: Request desde origen no permitido → CORS bloqueado
- 1.8: Exceder rate limit → 429 rate_limit_exceeded
