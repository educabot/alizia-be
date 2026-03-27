# team-ai-toolkit — Libreria compartida de infraestructura Go

## Que es

Modulo Go reutilizable (`github.com/educabot/team-ai-toolkit`) que contiene toda la infraestructura comun entre los proyectos backend de Educabot. Cualquier proyecto nuevo importa `team-ai-toolkit` y arranca con: servidor HTTP (Gin), autenticacion JWT con HMAC-HS256, conexion a DB via GORM, logging con Bugsnag, paginacion, manejo de errores estandarizado, transacciones y abstraccion de framework.

> **NOTA:** Alizia y tich-cronos usan **JWT con HMAC-HS256** via `team-ai-toolkit/tokens`. El paquete `tokens` **crea y valida** JWT usando un secreto compartido (`JWT_SECRET`). El auth-service propio es un **plan futuro**. La arquitectura actual usa `tokens.Toker` para emitir y validar tokens directamente en cada proyecto.

**No contiene logica de dominio.** Solo infraestructura que no depende de ningun proyecto especifico.

---

## Contexto: Ecosistema Educabot

```
                    ┌──────────────────────┐
                    │  auth-service         │  (FUTURO)
                    │  Centralizara la      │
                    │  emision de JWT       │
                    │                       │
                    │  Hoy: cada proyecto   │
                    │  usa tokens.Toker     │
                    │  con JWT_SECRET       │
                    └──────────┬────────────┘
                               │
                               │ JWT firmado (HMAC-HS256)
                               │
          ┌────────────────────┼─────────────────────┐
          │                    │                      │
          ▼                    ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│    alizia-api   │  │  tich-cronos    │  │  Futuro proyecto │
│   (monolito)    │  │  (Cloud Fns)    │  │                  │
│                 │  │                 │  │                  │
│  boot.NewRouter │  │  boot.Functions │  │                  │
│  boot.NewServer │  │                 │  │                  │
│  DB: GORM       │  │  DB: GORM       │  │  DB: GORM        │
└────────┬────────┘  └────────┬────────┘  └────────┬─────────┘
         │                    │                     │
         └────────────────────┼─────────────────────┘
                              │
                              ▼
                    ┌──────────────────────┐
                    │  team-ai-toolkit     │  ← Esta libreria
                    │                      │
                    │  Importada por todos  │
                    │  los proyectos como   │
                    │  dependencia Go       │
                    └──────────────────────┘
```

**3 repos separados:**

| Repo | Tipo | Proposito |
|------|------|-----------|
| `educabot/auth-service` | Microservicio (futuro, no en uso) | Planificado para centralizar emision de JWT en el futuro |
| `educabot/team-ai-toolkit` | Libreria Go (no se deploya) | Infraestructura compartida. Se importa en `go.mod` |
| `educabot/alizia-api` | Monolito (deploy propio) | Plataforma Alizia. Usa `boot.NewRouter` + `boot.NewServer` |
| `educabot/tich-cronos` | Cloud Functions (deploy propio) | Plataforma TiCh. Usa `boot.Functions` |

---

## Autenticacion — JWT con HMAC-HS256 via tokens.Toker

### Como funciona

El paquete `tokens` contiene `Toker`, una struct que **crea y valida** JWT usando HMAC-HS256 con un secreto compartido (`JWT_SECRET`).

```go
// Crear un Toker
toker := tokens.New(cfg.JWTSecret) // HMAC-HS256 shared secret

// Crear un JWT
token, err := toker.Create(id, name, email, roles, duration)

// Validar un JWT
claims, err := toker.Validate(tokenString)
```

**Flujo:**
```
Login
  → toker.Create(id, name, email, roles, duration) → JWT firmado con HMAC-HS256
    → Frontend guarda JWT
      → Cada request envia JWT en Authorization: Bearer header
        → ValidateTokenMiddleware valida JWT con el mismo secret
          → Extrae Claims: ID, Name, Email, Roles
```

**Claims struct:**
```go
type Claims struct {
    ID     string   `json:"id"`
    Name   string   `json:"name"`
    Email  string   `json:"email"`
    Avatar string   `json:"avatar"`
    Roles  []string `json:"roles"`
    jwt.RegisteredClaims
}
```

**Metodos de Claims:** `HasRole(role)`, `HasAnyRole(roles...)`

**Middlewares disponibles:**
- `ValidateTokenMiddleware(toker, env)` — Valida Bearer JWT; en env `Test` auto-mockea claims
- `OptionalTokenMiddleware(toker, env)` — Igual pero permite requests sin autenticacion
- `RequireRole(roles...)` — Verifica roles, retorna 403 si no los tiene
- `ValidateServiceCredentials(allowedServices, serviceSecret)` — Basic auth para comunicacion service-to-service
- `GetClaims(req)` — Extrae `*Claims` del contexto del request

**Constantes:** `ClaimsKey`, `IDKey`, `EmailKey`, `TokenKey`, `TestUserID`

### Auth Service propio (FUTURO)

> El auth-service esta planificado para centralizar la emision de tokens a futuro. No es necesario para el lanzamiento de Alizia.

---

## Modulo y dependencias

```
module github.com/educabot/team-ai-toolkit

go 1.25.0
```

```
require (
    github.com/gin-gonic/gin              v1.12.0
    github.com/gin-contrib/cors            v1.7.6
    github.com/bugsnag/bugsnag-go/v2       v2.6.3
    github.com/bugsnag/bugsnag-go-gin      v1.0.0
    github.com/golang-jwt/jwt/v5           v5.3.1
    github.com/google/uuid                 v1.6.0
    github.com/stretchr/testify            v1.11.1
    gorm.io/gorm                           v1.31.1
    gorm.io/driver/postgres                v1.6.0
)
```

---

## Estructura de la libreria

```
team-ai-toolkit/
│
├── config/                              # Configuracion y variables de entorno
│   ├── environment.go                   # Environment type: Local, Develop, Staging, Production, Test
│   ├── base.go                          # BaseConfig struct + LoadBase()
│   └── env.go                           # Helpers: EnvOr(), MustEnv(), EnvSplit(), GetEnvUint()
│
├── web/                                 # Abstraccion HTTP framework-agnostic
│   ├── request.go                       # Request interface
│   │                                    #   Param(), Query(), Header()
│   │                                    #   Bind(), BindJSON(), BindQuery()
│   │                                    #   Set(), Get(), Context(), Next()
│   ├── response.go                      # Response struct { Status int, Body any }
│   │                                    #   JSON(), OK(), Created(), NoContent()
│   │                                    #   Err(), ErrResponse()
│   ├── handler.go                       # Handler = func(req Request) Response
│   ├── interceptor.go                   # Interceptor = func(req Request) Response
│   ├── mock_request.go                  # MockRequest para testing
│   │
│   └── gin/                             # Adaptador Gin
│       ├── handler.go                   # Adapt(h web.Handler) gin.HandlerFunc
│       ├── middleware.go                # AdaptMiddleware(m web.Interceptor) gin.HandlerFunc
│       └── request.go                   # NewRequest(c *gin.Context) *Request
│
├── boot/                                # Bootstrap de servidor HTTP
│   ├── router.go                        # NewRouter(env, allowedOrigins, middleware...) *gin.Engine
│   │                                    #   Configura Recovery, Logger, CORS
│   ├── server.go                        # NewServer(port, engine) *Server
│   │                                    #   Timeouts: Read=10s, Write=30s, Idle=60s
│   │                                    #   Run() — ListenAndServe
│   │                                    #   Shutdown() — graceful con 10s timeout
│   └── functions.go                     # Functions type para Cloud Functions (usado por tich-cronos)
│                                        #   NewFunctions(), Add(), POST/GET/PUT/DELETE/PATCH()
│                                        #   Path(), WithAuth()
│                                        #   Public(), Private(), Optional(), Internal()
│
├── dbconn/                              # Conexion a PostgreSQL via GORM
│   └── postgres.go                      # PostgresConfig struct
│                                        #   URL, Charset, Timeout
│                                        #   MaxOpenConnections, MaxIdleConnections
│                                        # NewPostgresConnector(config) PostgresConnector
│                                        # PostgresConnector.Connect() (*gorm.DB, error)
│
├── tokens/                              # Creacion y validacion de JWT (HMAC-HS256)
│   ├── toker.go                         # Toker struct
│   │                                    #   New(secret string) Toker
│   │                                    #   Toker.Create(id, name, email, roles, duration) (string, error)
│   │                                    #   Toker.Validate(tokenString) (*Claims, error)
│   ├── claims.go                        # Claims struct
│   │                                    #   ID string, Name string, Email string
│   │                                    #   Avatar string, Roles []string
│   │                                    #   Embeds jwt.RegisteredClaims
│   │                                    #   HasRole(role), HasAnyRole(roles...)
│   ├── middleware.go                    # ValidateTokenMiddleware(toker, env) web.Interceptor
│   │                                    #   Extrae Bearer token, valida JWT, inyecta Claims
│   │                                    #   En env Test: auto-mockea claims
│   │                                    # OptionalTokenMiddleware(toker, env) web.Interceptor
│   │                                    #   Igual pero permite requests sin auth
│   ├── roles.go                         # RequireRole(roles...) web.Interceptor
│   │                                    #   Chequea Claims.Roles, retorna 403 si no tiene
│   ├── service.go                       # ValidateServiceCredentials(allowedServices, serviceSecret)
│   │                                    #   Basic auth para service-to-service
│   └── context.go                       # GetClaims(req) *Claims
│                                        # Constantes: ClaimsKey, IDKey, EmailKey, TokenKey, TestUserID
│
├── errors/                              # Errores compartidos + mapeo HTTP
│   ├── errors.go                        # Sentinel errors:
│   │                                    #   ErrNotFound, ErrValidation, ErrUnauthorized
│   │                                    #   ErrForbidden, ErrDuplicate, ErrConflict
│   │                                    # Re-exports: Is, As, New (de stdlib)
│   └── handler.go                       # HandleError(err) web.Response
│                                        #   NotFound → 404, Validation → 400
│                                        #   Unauthorized → 401, Forbidden → 403
│                                        #   Duplicate → 409, Conflict → 409
│                                        #   default → 500
│
├── pagination/                          # Paginacion con patron "has more"
│   ├── pagination.go                    # Pagination struct { Limit uint, Offset uint }
│   │                                    #   ApplyDefaults(defaultLimit, maxLimit)
│   │                                    #   FetchLimit() — retorna Limit+1 para check "has more"
│   │                                    #   Scope(p) func(db *gorm.DB) *gorm.DB — GORM scope
│   └── response.go                      # PaginatedResponse[T] { Items []T, More bool }
│                                        #   NewPaginatedResponse[S,T](records, limit, transform)
│                                        #   HasMore[T](records, limit) — trim + bool
│
├── transactions/                        # Transacciones con GORM y patron Persistor
│   ├── transactor.go                    # Transactor interface: Run(task, dep, deps...) error
│   │                                    # New(db *gorm.DB) Transactor → TransactorImpl
│   │                                    # TransactorImpl.Run() — begin tx, wrap, execute, commit/rollback
│   ├── persistor.go                     # Persistor interface: With(tx *gorm.DB) Persistor
│   │                                    # Dependencies = []Persistor
│   │                                    # Task = func(dependencies Dependencies) error
│   ├── helpers.go                       # OfType[P Persistor](deps) P — generic helper
│   └── mock.go                          # TransactorMock — ejecuta task sin transaccion real
│
├── applog/                              # Logging
│   ├── logger.go                        # Logger interface: Error(msg, fields...), ConsoleError(msg, fields...)
│   │                                    # Field struct { Key, Value }
│   │                                    # NewField(key, value) Field
│   ├── test_logger.go                   # NewTestLogger() *TestLogger — imprime errores a stdout
│   ├── mock_logger.go                   # NewMockLogger() *MockLogger — no-op para tests silenciosos
│   │
│   └── bugsnag/                         # Integracion con Bugsnag
│       └── bugsnag.go                   # BugsnagLogger struct — implementa Logger, envia a Bugsnag
│                                        # NewBugsnagLogger() *BugsnagLogger
│                                        # Configure(apiKey, env) — configura Bugsnag globalmente
│                                        # GinMiddleware(apiKey, env, projectPackages) gin.HandlerFunc
│
├── wrappers/                            # Tipos custom para PostgreSQL
│   ├── uuid_slice.go                    # UUIDSlice = []uuid.UUID con Scan/Value para arrays UUID
│   ├── sql_time.go                      # SQLTime struct { Hour, Minute, Second }
│   │                                    #   Scan/Value para columnas TIME de PostgreSQL
│   │
│   └── date/                            # Tipo Date sin componente de hora
│       └── date.go                      # Date struct { Year, Month, Day }
│                                        # New(), FromTime(), FromStr(), Parse()
│                                        # IsZero(), ToTime(), String()
│                                        # JSON marshal/unmarshal como "YYYY-MM-DD"
│                                        # SQL Scan/Value driver interface
│
├── rest/                                # (legacy) Helpers de response
│   └── response.go                      # ErrFormat, OkFormat[T] — funcionalidad ahora en web/response.go
│
├── go.mod
└── go.sum
```

---

## Variables de entorno (BaseConfig)

`config.LoadBase()` lee las siguientes variables de entorno:

| Variable | Descripcion | Requerida |
|----------|-------------|-----------|
| `PORT` | Puerto del servidor HTTP | No (default segun proyecto) |
| `ENV` | Entorno: `local`, `develop`, `staging`, `production`, `test` | No (default: `local`) |
| `DATABASE_URL` | Connection string de PostgreSQL | Si |
| `JWT_SECRET` | Secreto compartido para HMAC-HS256 JWT | Si |
| `ALLOWED_ORIGINS` | Origenes CORS separados por coma | No |
| `API_KEY_BUGSNAG` | API key de Bugsnag para error reporting | No |

```go
type BaseConfig struct {
    Port           string
    Env            Environment
    DatabaseURL    string
    JWTSecret      string
    AllowedOrigins []string
    BugsnagAPIKey  string
}
```

**Helpers de config:**
- `EnvOr(key, fallback)` — Lee variable de entorno o retorna fallback
- `MustEnv(key)` — Lee variable de entorno o panic si no existe
- `EnvSplit(key, sep)` — Lee y splitea por separador
- `GetEnvUint(key, fallback)` — Lee como uint
- `IsProduction()` — Retorna `true` para `Production` o `Staging`

---

## Patrones clave

### Patron de paginacion "has more"

En lugar de contar el total de registros (query extra), el toolkit usa el patron "fetch N+1":

```go
// 1. El request trae Limit y Offset
var p pagination.Pagination
req.BindQuery(&p)
p.ApplyDefaults(20, 100) // default 20, max 100

// 2. Se buscan Limit+1 registros
var records []Entity
db.Scopes(pagination.Scope(p)).Limit(p.FetchLimit()).Find(&records)

// 3. Si vinieron Limit+1, hay mas paginas
items, hasMore := pagination.HasMore(records, p.Limit)

// 4. Response
return pagination.PaginatedResponse[EntityDTO]{
    Items: items,  // maximo Limit elementos (se recorta el +1)
    More:  hasMore, // true si habia mas de Limit registros
}
```

### Patron Transactor/Persistor

Para transacciones que involucran multiples repositorios:

```go
// Cada repositorio implementa Persistor
type UserRepo struct { db *gorm.DB }
func (r *UserRepo) With(tx *gorm.DB) transactions.Persistor {
    return &UserRepo{db: tx}
}

// El usecase recibe un Transactor
type CreateUserUsecase struct {
    transactor transactions.Transactor
    userRepo   *UserRepo
    auditRepo  *AuditRepo
}

// Ejecutar con transaccion
err := uc.transactor.Run(
    func(deps transactions.Dependencies) error {
        userRepo := transactions.OfType[*UserRepo](deps)
        auditRepo := transactions.OfType[*AuditRepo](deps)
        // Ambos repos usan la misma transaccion
        userRepo.Create(user)
        auditRepo.Log(action)
        return nil
    },
    uc.userRepo,  // persistors que seran wrapped con tx
    uc.auditRepo,
)
// Si error → rollback automatico, si ok → commit
```

### Dos modos de deploy: Monolito vs Cloud Functions

**Alizia** usa el modo monolito clasico:
```go
router := boot.NewRouter(cfg.Env, cfg.AllowedOrigins, authMiddleware)
// Registrar rutas en el router Gin
server := boot.NewServer(cfg.Port, router)
server.Run()    // ListenAndServe
server.Shutdown() // graceful shutdown
```

**tich-cronos** usa Cloud Functions:
```go
fns := boot.NewFunctions()
fns.Add("createUser",
    fns.POST("/users").
        WithAuth().
        Private(createUserHandler),
)
// Cada funcion se deploya individualmente
```

---

## Como lo importa cada proyecto

### Alizia

```go
// go.mod
module github.com/educabot/alizia-api

require github.com/educabot/team-ai-toolkit v1.x.x
```

```go
// config/config.go
package config

import bcfg "github.com/educabot/team-ai-toolkit/config"

type Config struct {
    bcfg.BaseConfig                     // Port, Env, DatabaseURL, JWTSecret, AllowedOrigins, BugsnagAPIKey
    AzureOpenAIKey      string
    AzureOpenAIEndpoint string
    AzureOpenAIModel    string
}

func Load() *Config {
    base := bcfg.LoadBase()
    return &Config{
        BaseConfig:          base,
        AzureOpenAIKey:      bcfg.MustEnv("AZURE_OPENAI_API_KEY"),
        AzureOpenAIEndpoint: bcfg.MustEnv("AZURE_OPENAI_ENDPOINT"),
        AzureOpenAIModel:    bcfg.EnvOr("AZURE_OPENAI_MODEL", "gpt-5-mini"),
    }
}
```

```go
// cmd/app.go
import (
    "github.com/educabot/team-ai-toolkit/boot"
    "github.com/educabot/team-ai-toolkit/dbconn"
    "github.com/educabot/team-ai-toolkit/tokens"
    "github.com/educabot/team-ai-toolkit/applog/bugsnag"
)

func NewApp(cfg *config.Config) *App {
    // Logging + error reporting
    bugsnag.Configure(cfg.BugsnagAPIKey, string(cfg.Env))
    logger := bugsnag.NewBugsnagLogger()

    // Database via GORM
    connector := dbconn.NewPostgresConnector(dbconn.PostgresConfig{URL: cfg.DatabaseURL})
    db, _ := connector.Connect()

    // JWT — crea y valida tokens con secreto compartido
    toker := tokens.New(cfg.JWTSecret)
    authMw := tokens.ValidateTokenMiddleware(toker, cfg.Env)

    // Router + Server (monolito)
    router := boot.NewRouter(string(cfg.Env), cfg.AllowedOrigins, authMw)
    // ... registrar rutas ...
    server := boot.NewServer(cfg.Port, router)
    return &App{db: db, server: server}
}
```

### tich-cronos

```go
// go.mod
module tichacademy.com/tich-cronos

require github.com/educabot/team-ai-toolkit v1.x.x
```

```go
// Usa boot.Functions para Cloud Functions en lugar de boot.NewRouter + boot.NewServer
import "github.com/educabot/team-ai-toolkit/boot"

fns := boot.NewFunctions()
fns.Add("getStudents",
    fns.GET("/students").
        WithAuth().
        Private(getStudentsHandler),
)
```

---

## Que va y que NO va en team-ai-toolkit

### Va (infraestructura generica)

| Paquete | Que resuelve | Quien lo usa |
|---------|-------------|-------------|
| `config/` | Variables de entorno, BaseConfig, helpers | Todos |
| `web/` | Abstraccion HTTP framework-agnostic (Handler, Request, Response) | Todos |
| `web/gin/` | Adaptador Gin (Adapt, AdaptMiddleware, NewRequest) | Todos (hoy) |
| `boot/` | Server lifecycle (NewRouter + NewServer) y Cloud Functions (Functions) | Todos |
| `dbconn/` | Conexion PostgreSQL via GORM con pool | Todos |
| `tokens/` | Creacion y validacion JWT HMAC-HS256, middlewares de auth y roles | Todos |
| `errors/` | Sentinel errors + HandleError() mapeo a HTTP | Todos |
| `pagination/` | Pagination struct, GORM scope, PaginatedResponse[T] con "has more" | Todos |
| `transactions/` | Transactor/Persistor pattern con GORM | Todos |
| `applog/` | Logger interface, TestLogger, MockLogger | Todos |
| `applog/bugsnag/` | BugsnagLogger + GinMiddleware para error reporting | Todos |
| `wrappers/` | UUIDSlice, SQLTime para PostgreSQL | Segun necesidad |
| `wrappers/date/` | Date type (YYYY-MM-DD) con JSON y SQL support | Segun necesidad |
| `rest/` | (legacy) ErrFormat, OkFormat — usar web/ en su lugar | Deprecado |

### NO va (dominio especifico)

| Cosa | Por que NO | Donde vive |
|------|-----------|------------|
| Entities/modelos | Cada proyecto tiene su dominio | `proyecto/src/core/entities/` |
| Providers/interfaces | Especificas del dominio | `proyecto/src/core/providers/` |
| Usecases | Logica de negocio propia | `proyecto/src/core/usecases/` |
| Handlers | Endpoints propios | `proyecto/src/entrypoints/` |
| Repositories | Queries propias | `proyecto/src/repositories/` |
| Migraciones | Schema propio | `proyecto/db/migrations/` |
| Config struct completo | Cada proyecto tiene campos distintos | `proyecto/config/` |
| AI client | Alizia usa Azure OpenAI, cronos puede usar otro | `proyecto/src/repositories/ai/` |

---

## Versionamiento

team-ai-toolkit usa **Go modules + semver**:

```
v1.0.0 → primera version estable
v1.1.0 → se agrega wrappers/date (backward compatible)
v1.2.0 → se agrega boot.Functions (backward compatible)
v2.0.0 → cambio en firma de Toker (breaking change)
```

Los proyectos fijan la version en `go.mod`:
```
require github.com/educabot/team-ai-toolkit v1.2.0
```

Actualizar es un `go get github.com/educabot/team-ai-toolkit@latest` + correr tests.

---

## Resumen

| Pregunta | Respuesta |
|----------|-----------|
| **Que es team-ai-toolkit?** | Libreria Go con infraestructura compartida |
| **Se deploya?** | No. Se importa como dependencia en `go.mod` |
| **Go version?** | 1.25.0 |
| **DB?** | GORM (gorm.io/gorm + gorm.io/driver/postgres) |
| **Auth?** | JWT con HMAC-HS256 via `tokens.Toker` (crea y valida tokens con `JWT_SECRET`) |
| **Que contiene?** | config/, web/, web/gin/, boot/, dbconn/, tokens/, errors/, pagination/, transactions/, applog/, applog/bugsnag/, wrappers/, wrappers/date/, rest/ (legacy) |
| **Quien lo usa?** | alizia-api (monolito), tich-cronos (Cloud Functions), futuros proyectos |
| **Dos modos de deploy?** | `boot.NewRouter` + `boot.NewServer` (Alizia) o `boot.Functions` (tich-cronos) |
| **Paginacion?** | Patron "has more": fetch N+1, trim, retorna `{ Items, More }` |
| **Transacciones?** | Patron Transactor/Persistor: wrap repositorios con tx, commit/rollback automatico |
| **Que NO va?** | Logica de dominio, entities, usecases, handlers, migraciones |
| **Como se versiona?** | Semver via Go modules (v1.0.0, v1.1.0, v2.0.0) |
