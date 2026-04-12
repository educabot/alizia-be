# Auth Login/Logout Integration — DONE

**Branch:** `feature/sl/auth-login-logout`
**Estado:** Integrado contra `team-ai-toolkit@v1.8.0` (tag publicado).

Este archivo se mantiene como registro del cierre de la integración. Ya no
hay nada bloqueante. Se puede borrar en cualquier commit futuro.

## Resumen final

- `team-ai-toolkit` v1.8.0 publicado con el paquete `auth/` (primitives-only).
- `alizia-be` consume v1.8.0 sin `replace` local.
- Password hashing: **argon2id** (OWASP 2024 params: 19 MiB, t=2, p=1). Se
  descartó bcrypt durante la review del PR del toolkit.
- JWT issuer: `alizia-be` (constante en `cmd/handlers.go`).
- Login/logout se wirean en alizia-be como handlers propios
  (`src/entrypoints/auth.go`) porque el toolkit intencionalmente **no** expone
  un `NewLoginHandler` genérico (ver `auth/TODO_V1.9.md` en el toolkit).

## Cambios vs el plan original

| Área | Plan original (v1.7.x) | Realidad (v1.8.0) |
|---|---|---|
| Hashing | bcrypt cost 12 | argon2id OWASP |
| Login handler | `ttauth.NewLoginHandler` del toolkit | handler propio en `src/entrypoints/auth.go` |
| Logout handler | `ttauth.NewLogoutHandler` del toolkit | handler propio en `src/entrypoints/auth.go` |
| `tokens.New` | `New(secret)` | `New(secret, issuer)` |
| `ComparePassword` | devuelve `bool` | devuelve `(bool, error)` |
| JWT con Audience | implícito | `Toker.CreateWithClaims` con `Audience=[orgID]` |

## Archivos finales

**Modificados:**
- `cmd/handlers.go` — wire handlers propios, issuer en `tokens.New`
- `cmd/repositories.go` — `AuthCredentials ttauth.CredentialsProvider`
- `db/seeds/seed.sql` — hashes argon2id de "admin123"
- `go.mod`, `go.sum` — `team-ai-toolkit v1.8.0`, sin `replace`
- `src/app/web/mapping.go` — `/api/v1/auth/{login,logout}` públicos
- `src/entrypoints/containers.go` — `Login` y `Logout` como `web.Handler`
- `src/repositories/auth/credentials_provider.go` — ajuste a firma
  `ComparePassword (bool, error)`
- `src/entrypoints/middleware/auth_integration_test.go`,
  `chain_integration_test.go` — issuer en `tokens.New`
- `scripts/hash_password/main.go` — usa `ttauth.HashPassword` (argon2id)

**Nuevos:**
- `src/entrypoints/auth.go` — `NewLoginHandler` + `NewLogoutHandler` propios
- `src/repositories/auth/credentials_provider_test.go` — 7 tests con
  testify/mock

## Credenciales de dev

Todos los users del seed comparten el hash de `admin123`.
Regenerar: `go run ./scripts/hash_password admin123` (cada corrida devuelve
un hash distinto por el salt random — pegar el resultante en los 4 rows).
