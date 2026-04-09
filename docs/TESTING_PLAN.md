# Plan de Testing — Flujo Completo

Backend local: `http://localhost:8080`

## Setup

```bash
# 1. Base de datos
docker compose up -d

# 2. Migraciones
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5480/alizia?sslmode=disable" up

# 3. Seed (opcional, ya hay datos base)
./scripts/seed.sh

# 4. Levantar server
set -a && source .env && set +a && air
```

## Datos disponibles (seed)

| Recurso | Valor |
|---|---|
| Org ID | `a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11` (Provincia de Neuquén) |
| Admin | id=1, Ana Admin, `admin@neuquen.edu.ar` |
| Coordinator | id=2, Carlos Coordinador, `coord@neuquen.edu.ar` |
| Teacher | id=3, María Docente, `teacher1@neuquen.edu.ar` |
| Multi-rol | id=4, Pedro Multirol (teacher + coordinator) |
| Áreas | 1=Ciencias, 2=Humanidades |

Tokens JWT (10 años de vida, ya en `.env.example`):
- `TEST_TOKEN_ADMIN`
- `TEST_TOKEN_COORDINATOR`
- `TEST_TOKEN_TEACHER`

En Postman: usar `{{admin_token}}` / `{{coord_token}}` / `{{teacher_token}}`.

## Flujos de testing

### Flujo 1 — Happy path del onboarding (teacher)

Simula un docente nuevo que entra por primera vez.

| # | Endpoint | Expected | Qué validar |
|---|---|---|---|
| 1 | `GET /health` | 200 `{"status":"ok"}` | Server arriba |
| 2 | `GET /api/v1/onboarding-config` con teacher | 200, devuelve `skip_allowed`, `profile_fields`, `tour_steps` | Config del org se lee bien |
| 3 | `GET /api/v1/users/me/onboarding-status` | 200 `{"completed": false, "completed_at": null}` | Estado inicial sin completar |
| 4 | `GET /api/v1/users/me/profile` | 200 `{}` (vacío) | Perfil sin datos |
| 5 | `PUT /api/v1/users/me/profile` con `{"disciplines":["Math","Physics"],"experience_years":5}` | 200/204 | Guardar perfil |
| 6 | `GET /api/v1/users/me/profile` | 200 con disciplines + experience_years | Lectura consistente post-guardado |
| 7 | `GET /api/v1/users/me/onboarding/tour-steps` | 200 lista filtrada por rol `teacher` | Solo steps del rol teacher, orden ascendente |
| 8 | `POST /api/v1/users/me/onboarding/complete` | 200/204 | Marca completado |
| 9 | `GET /api/v1/users/me/onboarding-status` | 200 `{"completed": true, "completed_at": "2026-..."}` | Timestamp RFC3339, no null |

**Verificación DB:**
```sql
docker exec alizia-postgres psql -U postgres -d alizia -c \
  "SELECT id, onboarding_completed_at, profile_data FROM users WHERE id=3;"
```

### Flujo 2 — Multi-rol (Pedro, id=4)

Token: hay que generar uno para id=4 o usar teacher/coordinator y cambiar manualmente.
El tour debe devolver steps de **ambos roles** deduplicados.

| # | Acción | Expected |
|---|---|---|
| 1 | `GET /users/me/onboarding/tour-steps` | Combina steps de teacher + coordinator, sin duplicados por `key`, ordenado por `order` |

### Flujo 3 — Admin: asignar/quitar coordinador

| # | Endpoint | Expected |
|---|---|---|
| 1 | `POST /api/v1/areas/1/coordinators` con `{"user_id": 3}` (admin token) | 200/201, asigna María al área Ciencias |
| 2 | `POST /api/v1/areas/1/coordinators` con `{"user_id": 3}` de nuevo | 409 Conflict (ya asignado) |
| 3 | `DELETE /api/v1/areas/1/coordinators/3` | 200/204 |
| 4 | `DELETE /api/v1/areas/1/coordinators/3` de nuevo | 404 Not Found |

**Verificación DB:**
```sql
docker exec alizia-postgres psql -U postgres -d alizia -c \
  "SELECT area_id, user_id FROM area_coordinators;"
```

### Flujo 4 — Seguridad y errores

| # | Caso | Endpoint | Expected |
|---|---|---|---|
| 1 | Sin token | `GET /users/me/onboarding-status` | 401 Unauthorized |
| 2 | Token malformado | `GET /users/me/onboarding-status` con `Bearer foo` | 401 |
| 3 | Role insuficiente | `POST /areas/1/coordinators` con teacher token | 403 Forbidden |
| 4 | Validación: profile sin body | `PUT /users/me/profile` con `{}` | 400 Validation |
| 5 | Validación: areas inexistente | `POST /areas/9999/coordinators` con admin | 404 Not Found |
| 6 | Validación: user_id inexistente | `POST /areas/1/coordinators` con `{"user_id": 9999}` | 404 Not Found |
| 7 | Cross-tenant | Token de org A accediendo a recurso de org B | 404 (o 403 según implementación) |

### Flujo 5 — Idempotencia

| # | Caso | Expected |
|---|---|---|
| 1 | `POST /onboarding/complete` dos veces | Ambas 200, `completed_at` **no cambia** en la segunda llamada |
| 2 | `PUT /users/me/profile` con mismo payload 2 veces | Ambas 200, sin duplicar datos |

## Postman: cómo correrlo

1. Abrir colección **Alizia-be** en Postman
2. Seleccionar environment **Local**
3. Verificar variables: `base_url=http://localhost:8080`, `admin_token=<jwt>`
4. Ejecutar requests en orden del Flujo 1
5. Para los flujos 3/4, cambiar el token en header (o crear variables `teacher_token`, `coord_token`)

## Comandos curl rápidos (si no querés Postman)

```bash
# Variables
export BASE=http://localhost:8080
export TOKEN_ADMIN="eyJhbGciOiJIUzI1NiIs...jWCFEvDftP_p..."  # TEST_TOKEN_ADMIN de .env.example
export TOKEN_TEACHER="eyJhbGciOiJIUzI1NiIs...DWfmInwjEY..."  # TEST_TOKEN_TEACHER

# 1. Health
curl -s $BASE/health

# 2. Config
curl -s $BASE/api/v1/onboarding-config -H "Authorization: Bearer $TOKEN_TEACHER"

# 3. Status inicial
curl -s $BASE/api/v1/users/me/onboarding-status -H "Authorization: Bearer $TOKEN_TEACHER"

# 4. Profile vacío
curl -s $BASE/api/v1/users/me/profile -H "Authorization: Bearer $TOKEN_TEACHER"

# 5. Guardar profile
curl -s -X PUT $BASE/api/v1/users/me/profile \
  -H "Authorization: Bearer $TOKEN_TEACHER" \
  -H "Content-Type: application/json" \
  -d '{"disciplines":["Math","Physics"],"experience_years":5}'

# 6. Tour steps
curl -s $BASE/api/v1/users/me/onboarding/tour-steps -H "Authorization: Bearer $TOKEN_TEACHER"

# 7. Complete
curl -s -X POST $BASE/api/v1/users/me/onboarding/complete -H "Authorization: Bearer $TOKEN_TEACHER"

# 8. Status final
curl -s $BASE/api/v1/users/me/onboarding-status -H "Authorization: Bearer $TOKEN_TEACHER"

# 9. Admin: asignar coordinador
curl -s -X POST $BASE/api/v1/areas/1/coordinators \
  -H "Authorization: Bearer $TOKEN_ADMIN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 3}'

# 10. Admin: quitar coordinador
curl -s -X DELETE $BASE/api/v1/areas/1/coordinators/3 \
  -H "Authorization: Bearer $TOKEN_ADMIN"
```

## Reset entre corridas

Para re-correr el flujo desde cero (teacher sin onboarding completado):

```sql
docker exec alizia-postgres psql -U postgres -d alizia -c \
  "UPDATE users SET onboarding_completed_at = NULL, profile_data = '{}' WHERE id IN (1,2,3,4);"
```

## Criterios de éxito

- [ ] Flujo 1 completo sin errores, `completed_at` en formato RFC3339
- [ ] Flujo 2: tour deduplicado y ordenado correctamente
- [ ] Flujo 3: alta y baja de coordinador persiste en DB
- [ ] Flujo 4: todos los errores con status code correcto
- [ ] Flujo 5: idempotencia preservada
- [ ] Logs del server sin panics ni errores 5xx inesperados
