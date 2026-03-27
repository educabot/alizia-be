# HU-8.4: Biblioteca y exploración

> Como docente, necesito explorar los recursos existentes de mi organización antes de generar uno nuevo, para reutilizar contenido y reducir duplicación.

**Fase:** 6 — Contenido y recursos
**Prioridad:** Media
**Estimación:** —

---

## Criterios de aceptación

- [ ] El docente puede ver todos los recursos de su organización en una vista de biblioteca
- [ ] Se puede filtrar por tipo de recurso, materia y creador
- [ ] El filtro por materia es UX (restricción soft), no de permisos — un docente de matemáticas puede ver recursos de otras materias si quiere
- [ ] Se muestra información clave: título, tipo, materia, creador, fecha
- [ ] Click en un recurso abre el detalle con el contenido renderizado
- [ ] Se pueden buscar recursos por texto en el título

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 8.4.1 | [Endpoints de biblioteca](./tareas/T-8.4.1-endpoints-biblioteca.md) | handlers/ | ⬜ |
| 8.4.2 | [Frontend de biblioteca](./tareas/T-8.4.2-frontend-biblioteca.md) | frontend/ | ⬜ |
| 8.4.3 | [Tests](./tareas/T-8.4.3-tests.md) | tests/ | ⬜ |

## Dependencias

- [HU-8.1: Modelo de datos](../HU-8.1-modelo-datos-recursos/HU-8.1-modelo-datos-recursos.md) — Tablas y queries base

## Test cases

- 8.12: Listar recursos de la org → retorna todos los recursos
- 8.13: Filtrar por tipo → solo recursos del tipo seleccionado
- 8.14: Filtrar por materia → recursos vinculados a esa materia
- 8.15: Buscar por título → match parcial
