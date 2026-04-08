# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go backend for **Alizia**, a multi-tenant educational planning platform. Coordinators plan how knowledge/skills will be taught across subjects throughout the school year, with AI assistance.

## Rules

- Keep code minimal, easy to navigate and modify
- Code, schema, everything must be in English
- Communication in Spanish
- Follow Clean Architecture by layers (entities, providers, usecases, entrypoints, repositories)
- One file = one responsibility (especially in usecases)
- Usecases NEVER import infrastructure — only providers (interfaces) and entities
- Every usecase Request struct must have a `Validate() error` method returning `providers.ErrValidation` (wrapped with `fmt.Errorf("%w: ...", ...)`), called as the first statement of `Execute`. Always validate tenant scope (`OrgID`) plus all other required fields

## Architecture

```
cmd/              <- Entry point + DI manual (main, app, repositories, usecases, handlers, routes)
src/core/         <- Pure domain: entities, providers (interfaces), usecases
src/entrypoints/  <- HTTP handlers (REST)
src/repositories/ <- GORM + raw SQL implementations
src/mocks/        <- Mocks for all layers
src/app/web/      <- Route mapping
config/           <- Config with env vars
db/migrations/    <- SQL migrations (up/down)
```

## Domain Concepts

### Educational Structure
```
organizations (multi-tenant root)
  └── areas (e.g., "Sciences", "Humanities") — have a coordinator
        └── subjects (e.g., "Mathematics", "Physics")
```

### Knowledge Taxonomy (topics)
```
topics (hierarchical tree, org-scoped)
  └── children topics (recursive)
```

### Coordination Document (main output)
Planning document per area that:
1. Selects topics to cover
2. Distributes topics across subjects
3. Generates class-by-class plan per subject with AI
4. Has sections (dynamic per org, JSONB)
5. States: pending → in_progress → published (documento vivo, editable post-publication)

### Lesson Plans
Teacher's detailed class plan based on coordination document classes.

## Commands

```bash
# Start PostgreSQL
docker compose up -d

# Run the server (hot reload)
make run

# Build
make build

# Run tests
make test

# Lint
make lint

# Run migrations
make migrate
```

## Database

PostgreSQL local: `postgresql://postgres:postgres@localhost:5480/alizia?sslmode=disable`

## API

- Health: `GET /health` → 200 `{"status":"ok"}`
- Base: `/api/v1` (all endpoints authenticated via JWT)

## Environment Variables

See `.env.example`. Required: `DATABASE_URL`, `JWT_SECRET`, `AZURE_OPENAI_API_KEY`, `AZURE_OPENAI_ENDPOINT`.

## Stack

| Component | Technology |
|---|---|
| Language | Go 1.26+ |
| ORM | GORM |
| Database | PostgreSQL 16 |
| AI | Azure OpenAI (gpt-5-mini) |
| Auth | JWT |
| Migrations | golang-migrate |
| Testing | testify |
| Deploy | Railway (Docker) |

## File Structure

```
alizia-be/
├── cmd/                    # Entry point + DI
│   ├── main.go
│   ├── app.go
│   ├── repositories.go
│   ├── usecases.go
│   ├── handlers.go
│   └── routes.go
├── src/
│   ├── core/
│   │   ├── entities/       # Pure data structs
│   │   ├── providers/      # Interfaces (contracts)
│   │   └── usecases/       # Business logic
│   │       ├── admin/
│   │       ├── coordination/
│   │       ├── teaching/
│   │       ├── resources/
│   │       └── ai/
│   ├── entrypoints/        # HTTP handlers
│   │   └── rest/
│   ├── repositories/       # GORM implementations
│   │   ├── admin/
│   │   ├── coordination/
│   │   ├── teaching/
│   │   ├── resources/
│   │   └── ai/
│   ├── mocks/              # Test mocks
│   ├── app/web/            # Route mapping
│   └── utils/
├── config/
├── db/migrations/
├── scripts/
├── docs/rfc-alizia/        # Technical RFC documentation
├── .github/workflows/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .air.toml
└── .golangci.yml
```

## Skills

- `/rfc-docs` - **ALWAYS USE THIS** when the user asks to create, desglosar, or modify epicas, historias de usuario (HU), or tareas (T) in the RFC documentation
- `/audit-rfc` - **ALWAYS USE THIS** when the user asks to audit, review, or check the RFC documentation
