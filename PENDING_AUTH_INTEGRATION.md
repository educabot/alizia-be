# Pending: Auth Login/Logout Integration

**Branch:** `feature/sl/auth-login-logout`
**Estado:** Esperando release `v1.8.0` de `team-ai-toolkit` en GitHub.

## Contexto

Se agregaron los endpoints `POST /api/v1/auth/login` y `POST /api/v1/auth/logout` al backend, consumiendo el nuevo paquete `auth` que se agregó al `team-ai-toolkit`.

- **team-ai-toolkit**: branch `feature/auth-package`, commit `d5c4386`, pusheado a origin. Listo para taggear `v1.8.0`.
- **alizia-be**: cambios en esta branch, **sin commitear**, usando `replace` local al toolkit para desarrollo.

## Validación end-to-end ya realizada (con replace local)

- `POST /auth/login` con `admin@neuquen.edu.ar` / `admin123` → 200 + JWT (24h)
- `POST /auth/login` password incorrecta → 401 `{"code":"unauthorized","description":"invalid credentials"}`
- `POST /auth/login` body `{}` → 400 `{"code":"bad_request","description":"email and password are required"}`
- `GET /api/v1/users/me/onboarding-status` con JWT fresco → 200 (tenant middleware extrae `org_id` del claim `aud`)
- `POST /auth/logout` → 200 `{"status":"logged_out"}`

## Archivos modificados / nuevos en esta branch

**Modificados:**
- `cmd/app.go` — `NewHandlers` ahora recibe `*Repositories`
- `cmd/handlers.go` — wire `ttauth.NewLoginHandler` (duración 24h) y `ttauth.NewLogoutHandler`
- `cmd/repositories.go` — agrega `AuthCredentials ttauth.CredentialsProvider`
- `db/seeds/seed.sql` — seeds con `password_hash` bcrypt de "admin123"
- `go.mod`, `go.sum` — `replace github.com/educabot/team-ai-toolkit => ../team-ai-toolkit` (temporal)
- `src/app/web/mapping.go` — grupo público `/api/v1/auth/{login,logout}` sin middleware de auth/tenant
- `src/entrypoints/containers.go` — agrega `Login` y `Logout` como `web.Handler`

**Nuevos:**
- `scripts/hash_password/main.go` — CLI helper para generar hashes bcrypt (cost 12)
- `src/repositories/auth/credentials_provider.go` — implementación GORM de `CredentialsProvider`
- `src/repositories/auth/credentials_provider_test.go` — 7 tests con testify/mock

## Pasos pendientes (ejecutar cuando la release v1.8.0 esté publicada)

```bash
cd C:/Users/Educabot/Desktop/Educabot/Repositorios/alizia-be
git checkout feature/sl/auth-login-logout

# 1. Bump a la versión publicada
go get github.com/educabot/team-ai-toolkit@v1.8.0

# 2. Quitar el replace local
go mod edit -dropreplace github.com/educabot/team-ai-toolkit
go mod tidy

# 3. Validar build y tests
go build ./...
go vet ./...
go test ./...

# 4. Rebuild y relanzar backend
# (matar proceso actual si sigue corriendo con binary antiguo)

# 5. Re-correr los 5 curl de validación end-to-end
#    (login OK, login wrong pass, login empty body, GET protegido con JWT, logout)

# 6. Commit y push
git add -A
git commit -m "feat(auth): add login/logout endpoints using team-ai-toolkit v1.8.0"
git push -u origin feature/sl/auth-login-logout
```

## Notas

- El backend quedó corriendo en background (task `beqnckszl`) con el binary que usa el `replace` local. Sigue válido para testing manual mientras tanto.
- Los seeds de usuarios todos tienen password `admin123` (bcrypt cost 12).
- Para regenerar el hash: `go run ./scripts/hash_password admin123`
