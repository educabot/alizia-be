# HU-3.3: Topics — jerarquía curricular

> Como admin, necesito cargar la estructura curricular de la provincia como una jerarquía dinámica de topics para que coordinadores y docentes trabajen con los temas oficiales.

**Fase:** 2 — Admin/Integration
**Prioridad:** Alta (bloqueante para Épica 4)
**Estimación:** —

---

## Criterios de aceptación

- [ ] Tabla `topics` con: id, organization_id, parent_id (self-ref, nullable), name, description, level (precalculado), created_at
- [ ] `parent_id IS NULL` → `level = 1` (raíz)
- [ ] `parent_id IS NOT NULL` → `level = parent.level + 1`
- [ ] No se puede crear un topic que exceda `config.topic_max_levels` de la org
- [ ] Al mover un topic (cambiar parent_id), se re-computa el level de toda la rama
- [ ] Si re-computar causaría que algún descendiente exceda `topic_max_levels` → operación rechazada
- [ ] Endpoint `POST /api/v1/topics` crea topic (acepta parent_id nullable)
- [ ] Endpoint `GET /api/v1/topics` retorna árbol completo de la org (tree structure)
- [ ] Endpoint `GET /api/v1/topics?level=N` filtra por nivel
- [ ] Seed con jerarquía de ejemplo (3 niveles: núcleos → áreas de conocimiento → categorías)

## Tareas

| # | Tarea | Archivo | Estado |
|---|-------|---------|--------|
| 3.3.1 | [Migración: topics](./tareas/T-3.3.1-migracion-topics.md) | db/migrations/ | ⬜ |
| 3.3.2 | [Entity y provider](./tareas/T-3.3.2-entity-provider.md) | src/core/ | ⬜ |
| 3.3.3 | [Repository GORM + tree query](./tareas/T-3.3.3-repository-tree.md) | src/repositories/ | ⬜ |
| 3.3.4 | [Usecase: validación de niveles](./tareas/T-3.3.4-usecase-validacion-niveles.md) | src/core/usecases/ | ⬜ |
| 3.3.5 | [Endpoints y handler](./tareas/T-3.3.5-endpoints-handler.md) | src/entrypoints/ | ⬜ |
| 3.3.6 | [Seed y tests](./tareas/T-3.3.6-seed-tests.md) | db/seeds/, tests/ | ⬜ |

## Dependencias

- [HU-3.1: Organizaciones](../HU-3.1-organizaciones-configuracion/HU-3.1-organizaciones-configuracion.md) — config con `topic_max_levels` y `topic_level_names`

## Diseño técnico

### Tabla auto-referencial

```sql
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id),
    parent_id INTEGER REFERENCES topics(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Re-cómputo de levels

```sql
WITH RECURSIVE tree AS (
    SELECT id, parent_id,
           COALESCE((SELECT level FROM topics WHERE id = NEW.parent_id), 0) + 1 AS level
    FROM topics WHERE id = NEW.id
    UNION ALL
    SELECT t.id, t.parent_id, tree.level + 1
    FROM topics t JOIN tree ON t.parent_id = tree.id
)
UPDATE topics SET level = tree.level
FROM tree WHERE topics.id = tree.id;
```

### Config de la org

```jsonc
{
  "topic_max_levels": 3,
  "topic_level_names": ["Núcleos problemáticos", "Áreas de conocimiento", "Categorías"],
  "topic_selection_level": 3  // nivel al que se seleccionan topics en el wizard del doc
}
```

### Ejemplo de árbol

```
Pensamiento Lógico-Matemático          (level 1 — Núcleo)
├── Aritmética Básica                   (level 2 — Área de conocimiento)
│   ├── Suma y resta                    (level 3 — Categoría)
│   └── Multiplicación                  (level 3 — Categoría)
└── Geometría                           (level 2 — Área de conocimiento)
    ├── Figuras planas                  (level 3 — Categoría)
    └── Cuerpos geométricos             (level 3 — Categoría)
```

## Test cases

- 3.9: POST topic sin parent → level = 1
- 3.10: POST topic con parent level 2 → level = 3
- 3.11: POST topic que excede topic_max_levels → 422
- 3.12: GET topics → retorna árbol anidado
- 3.13: GET topics?level=3 → solo categorías
- 3.14: Mover topic cambiando parent_id → levels re-computados
