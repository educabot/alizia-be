# HU-1.2: Modelo de usuarios y roles

> Como admin, necesito poder crear usuarios con roles asignados para que cada persona tenga los permisos correctos en la plataforma.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta
**Estimación:** —

---

## Criterios de aceptación

- [ ] Migración crea tablas: `organizations`, `users`, `user_roles` + enum `member_role`
- [ ] Entity Go para Organization, User, UserRole
- [ ] Repository GORM con CRUD básico
- [ ] Un usuario puede tener múltiples roles (teacher + coordinator)
- [ ] `UNIQUE(user_id, role)` impide duplicados
- [ ] Seed de datos iniciales para testing

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 1.2.1 | [Migración: organizations + users + user_roles](./tareas/T-1.2.1-migracion.md) | db/migrations/000001_init.up.sql | ⬜ |
| 1.2.2 | [Entities Go](./tareas/T-1.2.2-entities.md) | internal/admin/entities/ | ⬜ |
| 1.2.3 | [Provider interfaces](./tareas/T-1.2.3-providers.md) | internal/admin/providers/ | ⬜ |
| 1.2.4 | [Repository GORM](./tareas/T-1.2.4-repository.md) | internal/admin/repositories/ | ⬜ |
| 1.2.5 | [Seed de datos iniciales](./tareas/T-1.2.5-seed.md) | db/seeds/ | ⬜ |
| 1.2.6 | [Tests del repository](./tareas/T-1.2.6-tests.md) | internal/admin/repositories/*_test.go | ⬜ |

## Modelo de datos

```
organizations (id, name, slug, config JSONB, created_at)
users (id, organization_id FK, email, name, password_hash, avatar_url, created_at)
user_roles (id, user_id FK, role member_role, UNIQUE(user_id, role))

member_role ENUM: 'teacher', 'coordinator', 'admin'
```

## Dependencias

- Épica 0 completada (PostgreSQL corriendo)
