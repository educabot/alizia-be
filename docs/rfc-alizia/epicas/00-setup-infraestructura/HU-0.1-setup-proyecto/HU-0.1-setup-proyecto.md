# HU-0.1: Setup del proyecto e infraestructura

> Como equipo de desarrollo, necesito tener el proyecto base funcionando con CI/CD y deploy para poder empezar a construir features.

**Fase:** 1 — Setup
**Prioridad:** Alta (bloqueante para todo lo demás)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Repo `alizia-api` creado en GitHub con estructura de directorios Clean Architecture
- [ ] `go.mod` configurado con `team-ai-toolkit` como dependencia
- [ ] GitHub Actions corre tests y linting en cada PR
- [ ] Railway configurado con proyecto + PostgreSQL
- [ ] `/health` responde 200 `{"status": "ok"}` en staging
- [ ] Deploy automático desde branch `main`

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 0.1.1 | [Crear repo con estructura de directorios](./tareas/T-0.1.1-crear-repo.md) | — | ⬜ |
| 0.1.2 | [Configurar go.mod con team-ai-toolkit](./tareas/T-0.1.2-go-mod.md) | go.mod | ⬜ |
| 0.1.3 | [Configurar GitHub Actions (test + lint)](./tareas/T-0.1.3-ci.md) | .github/workflows/ | ⬜ |
| 0.1.4 | [Provisionar Railway + PostgreSQL](./tareas/T-0.1.4-railway.md) | — | ⬜ |
| 0.1.5 | [Dockerfile + docker-compose local](./tareas/T-0.1.5-docker.md) | Dockerfile, docker-compose.yml | ⬜ |
| 0.1.6 | [Endpoint /health](./tareas/T-0.1.6-health.md) | cmd/main.go | ⬜ |
| 0.1.7 | [Deploy inicial a staging](./tareas/T-0.1.7-deploy.md) | — | ⬜ |

## Dependencias

- team-ai-toolkit compilando y publicado en GitHub
- Cuenta de Railway con permisos

## Test cases

- 1.1: GET /health → 200 `{"status": "ok"}`
