# back-config — Librería compartida de infraestructura Go

## Qué es

Módulo Go reutilizable (`github.com/educabot/back-config`) que contiene toda la infraestructura común entre los proyectos backend de Educabot. Cualquier proyecto nuevo importa `back-config` y arranca con: servidor HTTP, auth JWT, conexión a DB, logging, paginación, errores estandarizados y abstracción de framework.

**No contiene lógica de dominio.** Solo infraestructura que no depende de ningún proyecto específico.

---

## Contexto: Ecosistema Educabot

```
                    ┌──────────────────────┐
                    │    Auth Service       │  ← Microservicio propio (repo separado)
                    │                      │
                    │  - Login / Register   │
                    │  - JWT RS256 (firma)  │
                    │  - Refresh tokens     │
                    │  - Password reset     │
                    │  - organizations      │
                    │  - users + roles      │
                    │  - sessions           │
                    │                      │
                    │  DB: auth_db          │
                    └──────────┬───────────┘
                               │
                               │ JWT firmado con private key
                               │
          ┌────────────────────┼─────────────────────┐
          │                    │                      │
          ▼                    ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Alizia v2     │  │  tich-cronos    │  │  Futuro proyecto │
│   (monolito)    │  │  (refactorizado)│  │                 │
│                 │  │                 │  │                 │
│  DB: alizia_db  │  │  DB: cronos_db  │  │  DB: propia     │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                     │
         └────────────────────┼─────────────────────┘
                              │
                              ▼
                    ┌──────────────────────┐
                    │    back-config        │  ← Esta librería
                    │                      │
                    │  Importada por todos  │
                    │  los proyectos como   │
                    │  dependencia Go       │
                    └──────────────────────┘
```

**3 repos separados:**

| Repo | Tipo | Propósito |
|------|------|-----------|
| `educabot/auth-service` | Microservicio (deploy propio) | Emite y gestiona JWT. Base de datos propia con users, orgs, roles, sessions |
| `educabot/back-config` | Librería Go (no se deploya) | Infraestructura compartida. Se importa en `go.mod` |
| `educabot/alizia-api` | Monolito (deploy propio) | Plataforma Alizia. Importa back-config |
| `educabot/tich-cronos` | Monolito (deploy propio) | Plataforma TiCh. Importa back-config |

---

## Auth Service — El microservicio

### Qué hace

Es el **único servicio que maneja credenciales y emite tokens**. Los demás proyectos nunca tocan passwords, nunca crean tokens, solo los validan.

### Endpoints

```
POST /auth/login              → Email + password → JWT access + refresh token
POST /auth/register           → Crea usuario + org (o asigna a org existente)
POST /auth/refresh            → Refresh token → nuevo JWT access
POST /auth/password-reset     → Envía email con link de reset
POST /auth/password-reset/confirm → Nuevo password con token del email
GET  /auth/me                 → Info del usuario autenticado
```

### JWT emitido (RS256)

```json
{
  "sub": 42,
  "org_id": 1,
  "roles": ["teacher", "coordinator"],
  "email": "carlos@school.edu",
  "name": "Carlos Coordinador",
  "iat": 1742486400,
  "exp": 1742490000
}
```

Firmado con **private key** (solo el auth service la tiene).

### Base de datos (auth_db)

```sql
organizations (id, name, slug, config, created_at)
users (id, organization_id, email, password_hash, name, avatar_url, created_at)
user_roles (id, user_id, role, UNIQUE(user_id, role))
refresh_tokens (id, user_id, token_hash, expires_at, created_at)
```

### Comunicación con los otros servicios

**Zero HTTP calls en runtime.** Los proyectos que consumen el auth service solo necesitan la **public key** para validar tokens. Si el auth service se cae, los usuarios ya logueados siguen trabajando.

```
Login (HTTP call al auth service, solo 1 vez)
  → Auth service devuelve JWT
    → Frontend guarda JWT
      → Cada request envía JWT en Authorization header
        → Backend valida JWT con PUBLIC KEY (local, sin HTTP)
          → Extrae user_id, org_id, roles del token
```

### El auth service importa back-config también

```go
// auth-service/go.mod
require github.com/educabot/back-config v1.x.x
```

Usa: `boot/`, `web/`, `dbconn/`, `applog/`, `config/`, `errors/`. No usa `tokens/` (él es el que CREA los tokens, no los valida).

---

## back-config — Estructura de la librería

```
back-config/
│
├── web/                                 # Abstracción HTTP framework-agnostic
│   ├── request.go                       # Request interface
│   │                                    #   Param(), Query(), Header(), Body()
│   │                                    #   Bind(), Set(), Get(), Next()
│   ├── response.go                      # Response struct
│   │                                    #   JSON(status, body), Err(status, code, msg)
│   ├── handler.go                       # Handler func(Request) Response
│   ├── interceptor.go                   # Interceptor func(Request) Response (middleware)
│   ├── error.go                         # Error response helpers
│   │
│   └── gin/                             # Adaptador Gin (reemplazable por chi/, echo/, etc.)
│       ├── handler.go                   # Adapt(web.Handler) → gin.HandlerFunc
│       ├── request.go                   # GinRequest implementa web.Request
│       ├── middleware.go                # AdaptMiddleware(web.Interceptor) → gin.HandlerFunc
│       └── engine.go                    # Helpers
│
├── boot/                                # Bootstrap de servidor HTTP
│   ├── server.go                        # Server struct
│   │                                    #   NewServer(port, engine) → *Server
│   │                                    #   Run() — ListenAndServe con timeouts
│   │                                    #   Shutdown() — graceful con context timeout
│   └── gin.go                           # NewEngine(env, allowedOrigins) → *gin.Engine
│                                        #   Recovery, CORS, slog middleware, /health
│
├── dbconn/                              # Conexión a PostgreSQL
│   └── postgres.go                      # MustConnect(dsn) → *sqlx.DB
│                                        #   MaxOpenConns: 25, MaxIdleConns: 10
│                                        #   ConnMaxLifetime: 5min
│
├── tokens/                              # Cliente JWT del Auth Service
│   ├── jwt.go                           # ValidateJWT(token, publicKey) → (*Claims, error)
│   │                                    #   Parsea y valida RS256
│   │                                    #   Verifica exp, iat
│   ├── claims.go                        # Claims struct
│   │                                    #   UserID int64, OrgID int64
│   │                                    #   Roles []string, Email string, Name string
│   ├── middleware.go                    # NewAuthInterceptor(publicKey) → web.Interceptor
│   │                                    #   Extrae Bearer token de Authorization header
│   │                                    #   Valida JWT, inyecta Claims en context
│   │                                    #   Retorna 401 si falta o es inválido
│   ├── tenant.go                        # NewTenantInterceptor() → web.Interceptor
│   │                                    #   Lee OrgID del Claims ya validado
│   │                                    #   Lo inyecta en context para multi-tenancy
│   ├── roles.go                         # RequireRole(roles...) → web.Interceptor
│   │                                    #   Chequea que Claims.Roles contenga al menos 1
│   │                                    #   Retorna 403 si no tiene permisos
│   └── context.go                       # MustClaimsFromContext(req) → Claims
│                                        #   UserIDFromContext(req) → int64
│                                        #   OrgIDFromContext(req) → int64
│
├── applog/                              # Logging con slog
│   ├── logger.go                        # Setup(env string)
│   │                                    #   prod/staging → JSON a stdout
│   │                                    #   local/develop → text con colores
│   └── test_logger.go                   # SetupTest() — logger silencioso para tests
│
├── pagination/                          # Paginación
│   ├── pagination.go                    # Pagination struct {Page, PerPage, Offset()}
│   │                                    #   ParseFromQuery(req web.Request) → Pagination
│   │                                    #   Defaults: page=1, per_page=20, max=100
│   └── response.go                      # PaginatedResponse[T] {Data, Total, Page, PerPage}
│
├── transactions/                        # Transacciones con sqlx
│   ├── transactor.go                    # RunInTx(ctx, db, func(tx) error) error
│   │                                    #   Begin → fn(tx) → Commit o Rollback
│   ├── dbtx.go                          # DBTX interface
│   │                                    #   GetContext, SelectContext, ExecContext
│   │                                    #   Implementado por *sqlx.DB y *sqlx.Tx
│   └── mock.go                          # MockDBTX para tests
│
├── errors/                              # Errores compartidos + mapeo HTTP
│   ├── errors.go                        # Sentinel errors comunes
│   │                                    #   ErrNotFound, ErrValidation, ErrUnauthorized
│   │                                    #   ErrForbidden, ErrDuplicate, ErrConflict
│   └── handler.go                       # HandleError(err) → web.Response
│                                        #   errors.Is() → mapea a HTTP status + code
│                                        #   not_found → 404
│                                        #   validation_error → 400
│                                        #   unauthorized → 401
│                                        #   forbidden → 403
│                                        #   duplicate → 409
│                                        #   default → 500 + log
│
├── config/                              # Helpers de configuración
│   ├── env.go                           # EnvOr(key, fallback) string
│   │                                    #   MustEnv(key) string — panic si falta
│   └── base.go                          # BaseConfig struct
│                                        #   Port, Env, DatabaseURL, AuthPublicKey
│                                        #   AllowedOrigins []string
│                                        #   LoadBase() → BaseConfig
│
├── go.mod                               # module github.com/educabot/back-config
└── go.sum
```

---

## Dependencias de back-config

```
require (
    github.com/gin-gonic/gin        v1.10.x
    github.com/jmoiron/sqlx          v1.4.x
    github.com/lib/pq                v1.10.x
    github.com/golang-jwt/jwt/v5     v5.x.x
)
```

4 dependencias. Minimalista.

---

## Cómo lo importa cada proyecto

### Alizia v2

```go
// go.mod
module github.com/educabot/alizia-api

require github.com/educabot/back-config v1.x.x
```

```go
// config/config.go
package config

import bcfg "github.com/educabot/back-config/config"

type Config struct {
    bcfg.BaseConfig                     // Port, Env, DatabaseURL, AuthPublicKey, AllowedOrigins
    AzureOpenAIKey      string
    AzureOpenAIEndpoint string
    AzureOpenAIModel    string
    BugsnagAPIKey       string
}

func Load() *Config {
    base := bcfg.LoadBase()
    return &Config{
        BaseConfig:          base,
        AzureOpenAIKey:      bcfg.MustEnv("AZURE_OPENAI_API_KEY"),
        AzureOpenAIEndpoint: bcfg.MustEnv("AZURE_OPENAI_ENDPOINT"),
        AzureOpenAIModel:    bcfg.EnvOr("AZURE_OPENAI_MODEL", "gpt-5-mini"),
        BugsnagAPIKey:       os.Getenv("API_KEY_BUGSNAG"),
    }
}
```

```go
// cmd/app.go
import (
    "github.com/educabot/back-config/boot"
    "github.com/educabot/back-config/dbconn"
    "github.com/educabot/back-config/tokens"
    "github.com/educabot/back-config/applog"
)

func NewApp(cfg *config.Config) *App {
    applog.Setup(cfg.Env)
    db := dbconn.MustConnect(cfg.DatabaseURL)
    engine := boot.NewEngine(cfg.Env, cfg.AllowedOrigins)

    // Auth middleware del auth service compartido
    authMw := tokens.NewAuthInterceptor(cfg.AuthPublicKey)
    tenantMw := tokens.NewTenantInterceptor()

    // ... wiring de repos, usecases, handlers
    server := boot.NewServer(cfg.Port, engine)
    return &App{db: db, server: server}
}
```

### tich-cronos (refactorizado)

```go
// go.mod
module tichacademy.com/tich-cronos

require github.com/educabot/back-config v1.x.x
```

```go
// config/config.go
package config

import bcfg "github.com/educabot/back-config/config"

type Config struct {
    bcfg.BaseConfig
    CanvasSigloClientID     string
    CanvasSigloClientSecret string
    CanvasSigloRedirectURI  string
    LLMOpenAIKey            string
    LLMOpenAIURL            string
    ContentGeneratorURL     string
    BugsnagAPIKey           string
}

func Load() *Config {
    base := bcfg.LoadBase()
    return &Config{
        BaseConfig:              base,
        CanvasSigloClientID:     bcfg.MustEnv("CANVAS_SIGLO_CLIENT_ID"),
        LLMOpenAIKey:            bcfg.MustEnv("LLM_OPENAI_API_KEY"),
        ContentGeneratorURL:     bcfg.MustEnv("CONTENT_GENERATOR_URL"),
        BugsnagAPIKey:           os.Getenv("API_KEY_BUGSNAG"),
        // ...
    }
}
```

```go
// cmd/app.go — mismo patrón que Alizia
import (
    "github.com/educabot/back-config/boot"
    "github.com/educabot/back-config/dbconn"
    "github.com/educabot/back-config/tokens"     // MISMO auth service
    "github.com/educabot/back-config/applog"
)
```

### Auth Service

```go
// go.mod
module github.com/educabot/auth-service

require github.com/educabot/back-config v1.x.x
```

```go
// Usa de back-config:
import (
    "github.com/educabot/back-config/boot"       // Server bootstrap
    "github.com/educabot/back-config/web"         // Handler abstraction
    "github.com/educabot/back-config/dbconn"      // PostgreSQL connection
    "github.com/educabot/back-config/applog"      // Logging
    "github.com/educabot/back-config/config"      // EnvOr(), MustEnv()
    "github.com/educabot/back-config/errors"      // ErrNotFound, HandleError()
    "github.com/educabot/back-config/pagination"  // Si tiene listados
)

// NO usa tokens/ (él CREA los tokens, no los valida)
// En su lugar tiene su propio paquete interno:
//   internal/jwt/issuer.go → SignJWT(claims, privateKey) → token string
```

---

## Qué va y qué NO va en back-config

### Va (infraestructura genérica)

| Paquete | Qué resuelve | Quién lo usa |
|---------|-------------|-------------|
| `web/` | Abstracción HTTP framework-agnostic | Todos |
| `web/gin/` | Adaptador Gin | Todos (hoy) |
| `boot/` | Server lifecycle, timeouts, shutdown | Todos |
| `dbconn/` | Conexión PostgreSQL con sqlx | Todos |
| `tokens/` | Validación JWT del auth service | Alizia, tich-cronos, futuros |
| `applog/` | Setup de slog | Todos |
| `pagination/` | Parse page/per_page + response wrapper | Todos |
| `transactions/` | RunInTx(), DBTX interface | Todos |
| `errors/` | Sentinel errors + HandleError() | Todos |
| `config/` | EnvOr(), MustEnv(), BaseConfig | Todos |

### NO va (dominio específico)

| Cosa | Por qué NO | Dónde vive |
|------|-----------|------------|
| Entities/modelos | Cada proyecto tiene su dominio | `proyecto/src/core/entities/` |
| Providers/interfaces | Específicas del dominio | `proyecto/src/core/providers/` |
| Usecases | Lógica de negocio propia | `proyecto/src/core/usecases/` |
| Handlers | Endpoints propios | `proyecto/src/entrypoints/` |
| Repositories | Queries propias | `proyecto/src/repositories/` |
| Migraciones | Schema propio | `proyecto/db/migrations/` |
| Config struct completo | Cada proyecto tiene campos distintos | `proyecto/config/` |
| AI client | Alizia usa Azure OpenAI, cronos puede usar otro | `proyecto/src/repositories/ai/` |
| Mocks | Mockean interfaces propias del proyecto | `proyecto/src/mocks/` |
| JWT issuer (private key) | Solo el auth service firma tokens | `auth-service/internal/jwt/` |
| Prompts/schemas AI | Contenido específico del producto | `proyecto/src/repositories/ai/prompts/` |

---

## Versionamiento

back-config usa **Go modules + semver**:

```
v1.0.0 → primera versión estable
v1.1.0 → se agrega pagination/response.go (backward compatible)
v1.2.0 → se agrega web/gin/engine.go (backward compatible)
v2.0.0 → se cambia firma de tokens.ValidateJWT (breaking change)
```

Los proyectos fijan la versión en `go.mod`:
```
require github.com/educabot/back-config v1.2.0
```

Actualizar es un `go get github.com/educabot/back-config@latest` + correr tests.

---

## Estructura final de los 4 repos

```
educabot/
├── back-config/                 # Librería compartida (no se deploya)
│   ├── web/                     #   Abstracción HTTP
│   ├── boot/                    #   Server bootstrap
│   ├── dbconn/                  #   PostgreSQL connection
│   ├── tokens/                  #   JWT validation (client del auth service)
│   ├── applog/                  #   Logging
│   ├── pagination/              #   Paginación
│   ├── transactions/            #   Transacciones DB
│   ├── errors/                  #   Errores compartidos
│   └── config/                  #   Config helpers + BaseConfig
│
├── auth-service/                # Microservicio de autenticación (deploy propio)
│   ├── cmd/                     #   Entry point
│   ├── internal/                #   JWT issuer (private key), bcrypt, sessions
│   ├── db/migrations/           #   organizations, users, user_roles, refresh_tokens
│   └── go.mod                   #   importa back-config
│
├── alizia-api/                  # Monolito Alizia v2 (deploy propio)
│   ├── cmd/                     #   Entry point + DI manual
│   ├── src/                     #   core/, entrypoints/, repositories/, mocks/
│   ├── config/                  #   Config propio (extiende BaseConfig)
│   ├── db/migrations/           #   24 tablas de dominio educativo
│   └── go.mod                   #   importa back-config
│
└── tich-cronos/                 # Monolito TiCh refactorizado (deploy propio)
    ├── cmd/                     #   Entry point + DI manual (sin Wire)
    ├── src/                     #   core/, entrypoints/, repositories/, mocks/
    ├── config/                  #   Config propio (extiende BaseConfig)
    ├── db/migrations/           #   26 tablas de dominio educativo
    └── go.mod                   #   importa back-config
```

---

## Resumen

| Pregunta | Respuesta |
|----------|-----------|
| **¿Qué es back-config?** | Librería Go con infraestructura compartida |
| **¿Se deploya?** | No. Se importa como dependencia en `go.mod` |
| **¿Qué contiene?** | web/, boot/, dbconn/, tokens/, applog/, pagination/, transactions/, errors/, config/ |
| **¿Quién lo usa?** | Alizia v2, tich-cronos, auth-service, futuros proyectos |
| **¿Qué es auth-service?** | Microservicio propio que reemplaza Auth0. Emite JWT RS256 |
| **¿Cómo se relacionan?** | Auth service firma tokens. back-config/tokens/ los valida. Los proyectos importan back-config |
| **¿Qué NO va?** | Lógica de dominio, entities, usecases, handlers, migraciones |
| **¿Cómo se versiona?** | Semver via Go modules (v1.0.0, v1.1.0, v2.0.0) |
