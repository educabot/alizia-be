# Épica 0: Setup e infraestructura

> Repositorio, CI/CD, infraestructura de deploy y entorno de desarrollo local.

**Estado:** MVP
**Fase de implementación:** Fase 1

---

## Problema

Antes de construir cualquier feature, el equipo necesita un proyecto base funcional con estructura definida, CI/CD corriendo, base de datos provisionada y deploy automático a staging.

## Objetivos

- Crear el repositorio con estructura Clean Architecture
- Configurar CI/CD con tests y linting automáticos
- Provisionar infraestructura de staging (Railway + PostgreSQL)
- Tener un endpoint `/health` respondiendo en producción
- Entorno de desarrollo local reproducible (Docker)

## Alcance MVP

**Incluye:**

- Repo con estructura de directorios estándar
- `go.mod` con team-ai-toolkit como dependencia base
- GitHub Actions (test + lint en cada PR)
- Railway con PostgreSQL
- Dockerfile multi-stage + docker-compose local
- Endpoint /health
- Deploy automático desde main

**No incluye:**

- Monitoring avanzado (Grafana, Prometheus) → horizonte
- Environments múltiples (preview per-PR) → por definir

---

## Historias de usuario

| # | Historia | Descripción | Fase | Tareas |
|---|---------|-------------|------|--------|
| HU-0.1 | [Setup del proyecto e infraestructura](./HU-0.1-setup-proyecto/HU-0.1-setup-proyecto.md) | Repo, CI/CD, Railway, PostgreSQL, Docker, /health | Fase 1 | 7 |

---

## Principios de diseño

- **Zero vendor lock-in:** Dockerfile portable a cualquier plataforma (Render, Fly.io, VPS)
- **Dev-prod parity:** Mismo PostgreSQL version local y en staging
- **CI rápido:** Pipeline < 3 min para no bloquear PRs

## Épicas relacionadas

- **Roles y accesos** — Necesita repo y DB funcionando para empezar con auth
- **Integración** — Necesita la estructura de módulos para crear entities

## Test cases asociados

- Fase 1: Test 1.1 (GET /health → 200)

Ver [testing.md](../../operaciones/testing.md) para la matriz completa.
