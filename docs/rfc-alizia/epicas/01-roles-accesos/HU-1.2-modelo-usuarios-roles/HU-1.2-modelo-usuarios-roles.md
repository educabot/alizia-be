# HU-1.2: Modelo de usuarios y roles

> Como admin, necesito poder crear usuarios con roles asignados para que cada persona tenga los permisos correctos en la plataforma.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [x] Migración crea tablas: `organizations`, `users`, `user_roles` + enum `member_role`
- [x] Entity Go para Organization, User, UserRole
- [x] Repository GORM con CRUD básico
- [x] Un usuario puede tener múltiples roles (teacher + coordinator)
- [x] `UNIQUE(user_id, role)` impide duplicados
- [x] Seed de datos iniciales para testing

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.2.1 | [Migración: organizations + users + user_roles](./tareas/T-1.2.1-migracion.md) | db/migrations/000001_create_base_tables.up.sql | ✅ |
| 1.2.2 | [Entities Go](./tareas/T-1.2.2-entities.md) | src/core/entities/user.go, organization.go | ✅ |
| 1.2.3 | [Provider interfaces](./tareas/T-1.2.3-providers.md) | src/core/providers/admin.go | ✅ |
| 1.2.4 | [Repository GORM](./tareas/T-1.2.4-repository.md) | src/repositories/admin/ | ✅ |
| 1.2.5 | [Seed de datos iniciales](./tareas/T-1.2.5-seed.md) | db/seeds/seed.sql | ✅ |
| 1.2.6 | [Tests del repository](./tareas/T-1.2.6-tests.md) | src/repositories/admin/*_test.go | ✅ |

## Modelo de datos

```
organizations (id, name, slug, config JSONB, created_at)
users (id, organization_id FK, email, name, password_hash, avatar_url, created_at)
user_roles (id, user_id FK, role member_role, UNIQUE(user_id, role))

member_role ENUM: 'teacher', 'coordinator', 'admin'
```

## Dependencias

- Épica 0 completada (PostgreSQL corriendo)
