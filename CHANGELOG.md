# Changelog

All notable changes to this project will be documented in this file.

## [0.0.2] - 2026-04-08

### Added

- JWT authentication middleware with HS256 validation via team-ai-toolkit/tokens
- Tenant middleware for multi-tenant org isolation
- Token refresh and CORS support
- User and role model with database migration (users, roles, user_roles, organizations)
- User entities, providers (interfaces), and GORM repository implementation
- Seed data for initial users and roles
- RequireRole authorization middleware with interceptor chain
- Standardized error handling for auth/authz (401/403 responses)
- Area coordinator assignment: admin endpoints to assign/remove coordinators from areas
- Migration for area_coordinators join table
- Admin-only route group with role-based access control
- Unit tests for auth, authorization, and assignment flows

## [0.0.1] - 2026-03-15

### Added

- Repository setup with Clean Architecture structure (`src/core/`, `src/entrypoints/`, `src/repositories/`)
- Go module with team-ai-toolkit, GORM, golang-migrate, testify, and Azure OpenAI SDK
- Manual dependency injection wiring in `cmd/` (main, app, repositories, usecases, handlers)
- Config struct embedding team-ai-toolkit BaseConfig
- Centralized route mapping in `src/app/web/mapping.go`
- GitHub Actions CI pipeline (test + lint with golangci-lint, 80% coverage target)
- Railway deployment with managed PostgreSQL
- Multi-stage Dockerfile and docker-compose for local development
- Makefile with run, build, test, lint, and migrate commands
- Air hot reload for development
- Health endpoint (`GET /health` → 200)
- Auto-deploy from main branch to Railway
- Database migrations directory structure
