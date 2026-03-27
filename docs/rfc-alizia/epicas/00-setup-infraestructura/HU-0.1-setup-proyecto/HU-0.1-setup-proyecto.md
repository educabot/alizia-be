# HU-0.1: Setup del proyecto e infraestructura

> Como equipo de desarrollo, necesito tener el proyecto base funcionando con arquitectura definida, CI/CD y deploy para poder empezar a construir features.

**Fase:** 1 — Setup
**Prioridad:** Alta (bloqueante para todo lo demás)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Repo `alizia-api` creado en GitHub con estructura `src/` Clean Architecture por capas
- [ ] `go.mod` configurado con team-ai-toolkit, GORM, golang-migrate, testify, Azure OpenAI SDK
- [ ] GitHub Actions corre tests y linting en cada PR (coverage target 80%)
- [ ] Railway configurado con proyecto + PostgreSQL managed
- [ ] Dockerfile multi-stage + docker-compose local con PostgreSQL
- [ ] Makefile con targets: build, test, test-cover, vet, lint, docker, migrate
- [ ] Air configurado para hot reload en desarrollo
- [ ] `/health` responde 200 `{"status": "ok"}` en staging
- [ ] `cmd/` scaffolding completo: main.go, app.go, repositories.go, usecases.go, handlers.go
- [ ] `config/config.go` embebiendo BaseConfig con campos propios (Azure OpenAI)
- [ ] Route mapping en `src/app/web/mapping.go` con middleware auth + tenant
- [ ] Error handling base: `src/core/providers/errors.go` + `src/entrypoints/rest/rest.go`
- [ ] Handler containers por feature (Admin, Coordination, Teaching, Resources)
- [ ] Deploy automático desde branch `main`
- [ ] README.md con stack, arquitectura, comandos, API overview
- [ ] TESTING.md con convenciones, estructura, CI

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 0.1.1 | [Crear repo con estructura de directorios](./tareas/T-0.1.1-crear-repo.md) | — | ⬜ |
| 0.1.2 | [Configurar go.mod con dependencias](./tareas/T-0.1.2-go-mod.md) | go.mod | ⬜ |
| 0.1.3 | [Configurar GitHub Actions (test + lint)](./tareas/T-0.1.3-ci.md) | .github/workflows/ | ⬜ |
| 0.1.4 | [Provisionar Railway + PostgreSQL](./tareas/T-0.1.4-railway.md) | — | ⬜ |
| 0.1.5 | [Docker, Makefile y dev tools](./tareas/T-0.1.5-docker.md) | Dockerfile, docker-compose.yml, Makefile | ⬜ |
| 0.1.6 | [Scaffolding cmd/ + config](./tareas/T-0.1.6-scaffolding-cmd.md) | cmd/, config/ | ⬜ |
| 0.1.7 | [Deploy inicial a staging](./tareas/T-0.1.7-deploy.md) | — | ⬜ |
| 0.1.8 | [Patterns base: errors, REST, containers, routes](./tareas/T-0.1.8-patterns-base.md) | src/ | ⬜ |
| 0.1.9 | [README.md + TESTING.md](./tareas/T-0.1.9-readme-testing.md) | README.md, TESTING.md | ⬜ |

## Dependencias

- team-ai-toolkit compilando y publicado en GitHub (v0.1.0+)
- Cuenta de Railway con permisos
- Cuenta de Bugsnag (o desactivar tracker)

## Test cases

- 1.1: GET /health → 200 `{"status": "ok"}`
