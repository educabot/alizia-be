# Alizia API

Backend Go para la plataforma educativa Alizia. Planificacion anual, coordinacion de areas, planificacion docente, generacion de recursos con IA.

## Features

- **Coordinacion**: Documentos de planificacion anual por area con generacion IA
- **Planificacion docente**: Lesson plans con momentos didacticos y actividades
- **Recursos**: Biblioteca de recursos generados por IA a partir de fuentes oficiales
- **Chat IA**: Asistente Alizia con function calling para editar documentos
- **Multi-tenant**: Aislamiento por organizacion via JWT claims

## Quick Start

```bash
docker compose up -d                # Levanta PostgreSQL
cp .env.example .env                # Configurar variables
make migrate                        # Correr migraciones
make run                            # Arranca con Air (hot reload)
```

## Stack

| Componente | Tecnologia |
|---|---|
| Language | Go 1.26+ |
| ORM | GORM |
| Database | PostgreSQL 16 |
| AI | Azure OpenAI |
| Auth | JWT |
| Deploy | Railway (Docker) |

## Architecture

```
cmd/            <- Entry point + DI manual (1 archivo por responsabilidad)
src/core/       <- Dominio puro: entities, providers (interfaces), usecases
src/entrypoints/<- HTTP handlers (REST)
src/repositories/<- Implementaciones GORM + raw SQL
src/mocks/      <- Mocks de todas las capas
src/app/        <- Route mapping
config/         <- Config con env vars
db/migrations/  <- SQL migrations (up/down)
```

Detalle completo en `docs/rfc-alizia/`.

## Development

```bash
make build        # CGO_ENABLED=0 go build -o alizia-api ./cmd
make test         # go test -race ./...
make test-cover   # test + coverage report
make vet          # go vet ./...
make lint         # golangci-lint run
make docker       # docker compose up -d
make migrate      # Correr migraciones
```

## Testing

Target: 90% coverage. Ver `TESTING.md` para convenciones y guia completa.

```bash
go test ./...
```

## API

- Health: `GET /health`
- Base: `/api/v1`

Todos los endpoints autenticados via Bearer token (JWT).

## Environment Variables

**Requeridas:**
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret for JWT validation
- `AZURE_OPENAI_API_KEY` - Azure OpenAI key
- `AZURE_OPENAI_ENDPOINT` - Azure OpenAI endpoint

**Opcionales:**
- `PORT` - Puerto del servidor (default: 8080)
- `ENV` - Entorno: local, staging, production (default: local)
- `AZURE_OPENAI_MODEL` - Modelo (default: gpt-5-mini)
- `BUGSNAG_API_KEY` - Error tracking
- `ALLOWED_ORIGINS` - CORS origins separados por coma

## Deploy

Push a main -> GitHub Actions (test + lint) -> Railway auto-deploy.

## License

Private.
