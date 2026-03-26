# HU-1.1: Autenticación con Auth0

> Como usuario, necesito autenticarme con email y contraseña para acceder a la plataforma con mis permisos correspondientes.

**Fase:** 1 — Setup
**Prioridad:** Alta (bloqueante para rutas protegidas)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Auth0 tenant configurado con domain + audience para staging
- [ ] JWT middleware valida tokens via JWKS (team-ai-toolkit/tokens)
- [ ] Claims extraídos del JWT: user_id, org_id, roles, email, name
- [ ] Tenant middleware inyecta org_id en el contexto
- [ ] Request sin token → 401 `missing_token`
- [ ] Request con token inválido → 401 `invalid_token`
- [ ] Request con token de otra org → datos filtrados por org_id

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.1.1 | [Configurar Auth0 tenant](./tareas/T-1.1.1-configurar-auth0.md) | — | ⬜ |
| 1.1.2 | [Integrar JWT middleware (JWKS)](./tareas/T-1.1.2-jwt-middleware.md) | cmd/main.go | ⬜ |
| 1.1.3 | [Integrar tenant middleware](./tareas/T-1.1.3-tenant-middleware.md) | cmd/main.go | ⬜ |
| 1.1.4 | [Config: Auth0 env vars](./tareas/T-1.1.4-config-auth0.md) | config/config.go | ⬜ |
| 1.1.5 | [Tests de autenticación](./tareas/T-1.1.5-tests-auth.md) | — | ⬜ |

## Dependencias

- Épica 0 completada (/health respondiendo)
- Auth0 tenant creado (mismo sistema que tich-cronos)
- team-ai-toolkit/tokens funcional

## Test cases

- 1.1: Request sin Authorization header → 401 `missing_token`
- 1.3: Request con JWT inválido → 401 `invalid_token`
- 1.4: Request con JWT válido → 200 + claims en context
- 1.5: Request con JWT de otra org → datos filtrados por org_id
