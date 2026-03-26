# Architecture Decision: Auth Microservice + Monolith Strategy

> **NOTA: Este documento describe el plan futuro de un auth-service propio. Alizia v2 arranca con Auth0 (mismo sistema que tich-cronos). El auth-service es un plan a futuro para reemplazar Auth0 cuando sea conveniente.**

## Context

The v2 proposal defines ~26 tables for the educational planning system. TiCh Cronos currently uses Auth0 for authentication, which Alizia v2 will also use for launch. A custom auth-service is planned for the future to eventually replace Auth0. This document defines the architectural strategy for both concerns.

---

## Decision: Modular Monolith + Auth Microservice

### Alizia v2 remains a monolith

The 24 domain tables (everything except auth) stay in a single service. The tables are heavily coupled — rendering a single coordination document view requires joining 8+ tables across topics, subjects, areas, classes, and organization config. Splitting these into microservices would turn simple DB joins into cascading HTTP calls with no real benefit.

**Key reasons:**

| Factor | Assessment |
|--------|-----------|
| Domain coupling | High — coordination docs reference topics, subjects, areas, org config in every operation |
| Transactionality | AI generation touches 3+ tables atomically (docs, classes, class_topics) |
| Team size | Small — microservices solve organizational problems, not technical ones |
| Multi-tenancy | Handled via `organization_id` filter middleware, not separate services |
| Operational cost | 1 deploy, 1 DB, 1 log stream vs N pipelines, service discovery, distributed tracing |

**Internal structure (not microservices, but modular):**

```
backend/
├── main.py                    # App setup + middleware only
├── routes/
│   ├── admin.py               # orgs, users, areas, subjects, courses
│   ├── coordination.py        # docs, wizard, sections
│   ├── teaching.py            # lesson plans, moments
│   ├── resources.py           # fonts, resource_types, resources
│   └── ai.py                  # /generate, /chat endpoints
├── services/
│   ├── ai_service.py          # Azure OpenAI calls, tool execution
│   ├── shared_classes.py      # Shared class calculation
│   └── topic_service.py       # Dynamic hierarchy, level validation
├── models/                    # Pydantic schemas
├── db/                        # Connection, queries
└── middleware/
    └── tenant.py              # Injects organization_id from JWT
```

This gives separation of concerns, independent testing, and clear ownership — without the operational overhead of distributed systems.

### Auth extracted as a shared microservice

Authentication is the one bounded context that justifies extraction because it meets all three conditions:

1. **Multiple consumers** — TiCh Cronos and Alizia (and future products) need authentication
2. **Clear bounded context** — Auth has its own domain: users, roles, orgs, login, tokens
3. **Concrete driver** — Deprecating Auth0 is a real cost reduction, not speculation

---

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────┐
│   TiCh Cronos   │     │   Alizia v2     │     │  Future API │
│   (existing)    │     │   (monolith)    │     │             │
└────────┬────────┘     └────────┬────────┘     └──────┬──────┘
         │                       │                      │
         │  JWT validate         │  JWT validate        │  JWT validate
         │  (public key,         │  (public key,        │  (public key,
         │   no HTTP call)       │   no HTTP call)      │   no HTTP call)
         │                       │                      │
         └───────────┬───────────┴──────────────────────┘
                     │
                     │  Login / register / refresh
                     │  (HTTP calls only at auth time)
                     │
                     ▼
            ┌─────────────────┐
            │  Auth Service   │
            │                 │
            │  - organizations│
            │  - users        │
            │  - user_roles   │
            │  - sessions     │
            └─────────────────┘
```

---

## Auth Service Scope

### What goes IN the Auth Service

| Table / Feature | Description |
|----------------|-------------|
| `organizations` | Tenant definition (id, name, slug, config) |
| `users` | Core identity (id, org_id, email, password_hash, name, avatar_url) |
| `user_roles` | Role assignments (teacher, coordinator, admin) |
| `refresh_tokens` / `sessions` | Token lifecycle management |
| Login / Register | Email + password authentication |
| JWT issuance | Signs tokens with private key |
| Token refresh | Rotates refresh tokens |
| Password reset | Reset flow via email |

### What stays OUT (in each product API)

| Concern | Reason |
|---------|--------|
| `area_coordinators` | Domain logic (which user coordinates which area) — not an auth concern |
| `course_subjects.teacher_id` | Domain assignment, not a role |
| Extended user profiles | Product-specific preferences stay in product DB |
| Domain-specific permissions | Each API resolves "can this coordinator edit this document?" locally |

---

## Communication Pattern: JWT with Public Key Validation

The Auth Service signs JWTs with an **asymmetric key pair** (RS256 or EdDSA). Consumer APIs only need the public key to validate tokens — **zero HTTP calls to Auth Service at runtime**.

### Auth Service (token issuance)

```python
# Auth Service signs the token
token = jwt.encode(
    {
        "sub": user.id,
        "org_id": user.organization_id,
        "roles": ["teacher", "coordinator"],
        "exp": datetime.utcnow() + timedelta(hours=1)
    },
    AUTH_PRIVATE_KEY,
    algorithm="RS256"
)
```

### Consumer APIs (token validation)

```python
# Alizia v2 / TiCh Cronos — local validation, no HTTP call
def get_current_user(token: str):
    payload = jwt.decode(token, AUTH_PUBLIC_KEY, algorithms=["RS256"])
    return {
        "user_id": payload["sub"],
        "org_id": payload["org_id"],
        "roles": payload["roles"]
    }
```

**Key property: if the Auth Service goes down, already-logged-in users continue working.** Only login/register/refresh require the Auth Service to be available.

### JWT Payload Structure

```json
{
  "sub": 42,
  "org_id": 1,
  "roles": ["teacher", "coordinator"],
  "email": "carlos@school.edu",
  "name": "Carlos Coordinador",
  "iat": 1742486400,
  "exp": 1742490000
}
```

Consumer APIs extract `user_id` and `org_id` from the token and use them to filter all queries (multi-tenancy) and resolve domain-specific permissions locally.

---

## Database Split

### Auth Service DB

```sql
-- Tenant
CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Identity
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Roles
CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role member_role NOT NULL,
    UNIQUE(user_id, role)
);

-- Sessions
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Alizia v2 DB

All 24 remaining domain tables. References to users are **logical** (no FK constraint):

```sql
-- No FK to auth DB — just an integer reference
CREATE TABLE area_coordinators (
    id SERIAL PRIMARY KEY,
    area_id INTEGER NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL,  -- References auth.users.id (logical)
    UNIQUE(area_id, user_id)
);

CREATE TABLE course_subjects (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id INTEGER NOT NULL,  -- References auth.users.id (logical)
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    school_year INTEGER NOT NULL
);
```

`organization_id` is injected from JWT via middleware — not stored redundantly on every row (or stored for query performance, depending on access patterns).

### TiCh Cronos DB

Same pattern — its own domain tables with `user_id` as logical reference from JWT.

---

## When to Extract Further

The monolith should only be split further when **real pain appears**, not speculatively:

| Signal | Action |
|--------|--------|
| AI generation needs 10x more scale than CRUD | Extract AI worker service with message queue |
| File storage (fonts/PDFs) grows significantly | Extract file service backed by S3 |
| Two teams of 5+ devs block each other on deploys | Split by bounded context (coordination vs teaching) |
| Another product needs the topic taxonomy | Extract topics as shared service |

**Rule: monolith first, extract when the pain is real.**

---

## Migration Path from Auth0 (PLAN FUTURO)

> Alizia v2 lanza con Auth0. Esta migración se ejecutará en el futuro cuando sea conveniente.

### Phase 0: Launch with Auth0 (CURRENT)
- Alizia v2 and tich-cronos use Auth0 for authentication
- team-ai-toolkit/tokens validates JWT via Auth0 JWKS
- No custom auth-service needed for launch

### Phase 1: Build Auth Service (FUTURE)
- Implement user/org/role tables
- JWT issuance with RS256
- Login, register, password reset endpoints
- Import existing Auth0 users (email + metadata)

### Phase 2: Integrate with Alizia v2 (FUTURE)
- Alizia v2 validates JWTs from Auth Service
- Remove Auth0 SDK dependency from Alizia

### Phase 3: Migrate TiCh Cronos (FUTURE)
- Replace Auth0 calls with Auth Service calls
- Validate JWTs with same public key
- Remove Auth0 SDK dependency

### Phase 4: Deprecate Auth0 (FUTURE)
- Disable Auth0 tenant
- Cancel Auth0 subscription

---

## Summary

| Component | Strategy | Database | Reason |
|-----------|----------|----------|--------|
| Auth | Shared microservice | Own DB (users, orgs, roles, sessions) | Multiple consumers + Auth0 deprecation |
| Alizia v2 | Modular monolith | Own DB (24 domain tables) | High coupling, transactionality, small team |
| TiCh Cronos | Existing service | Own DB | Consumes Auth Service via JWT |

The architecture optimizes for **simplicity where possible** (monolith for tightly coupled domain) and **separation where justified** (auth as shared infrastructure across products).
