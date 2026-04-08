# Comparativa: GORM vs sqlx vs sqlw

> **Decisión: GORM** — Es el estándar de la empresa y la elección definitiva para Alizia. Este documento se mantiene como registro de la comparación realizada.

## Resumen ejecutivo

| Aspecto | GORM | sqlx | sqlw |
|---------|------|------|------|
| **Tipo** | ORM completo | Extensión de database/sql | Wrapper de alto rendimiento |
| **Filosofía** | Convención sobre configuración | SQL explícito con comodidades | SQL crudo con máxima velocidad |
| **Popularidad** | ★★★★★ (36k+ stars) | ★★★★★ (16k+ stars) | ★★☆☆☆ (nicho) |
| **Curva de aprendizaje** | Media | Baja | Baja |
| **Performance** | Menor (overhead de reflexión) | Alta | Máxima |
| **Control SQL** | Bajo-medio | Total | Total |
| **Migraciones** | Integradas (AutoMigrate) | No incluidas | No incluidas |
| **Transacciones** | API propia (`db.Transaction()`) | `database/sql` estándar | `database/sql` estándar |

---

## GORM

**Repositorio**: github.com/go-gorm/gorm

### Qué es
ORM completo para Go. Abstrae el SQL detrás de una API chainable orientada a objetos.

### Ejemplo CRUD

```go
// Modelo
type User struct {
    ID    string `gorm:"primaryKey"`
    Name  string
    Email string
}

// Create
db.Create(&User{ID: "1", Name: "Juan", Email: "juan@test.com"})

// Read
var user User
db.Where("email = ?", "juan@test.com").First(&user)

// Read con scopes (composable)
db.Scopes(pagination.Scope(p)).Where("active = ?", true).Find(&users)

// Update
db.Model(&user).Update("name", "Juan Carlos")

// Delete
db.Delete(&user)

// Preload (relaciones)
db.Preload("Orders").Find(&users)
```

### Ventajas
- Migraciones automáticas con `AutoMigrate`
- Hooks (BeforeCreate, AfterUpdate, etc.)
- Preload de relaciones (1:1, 1:N, N:M)
- Scopes reutilizables para queries comunes
- Soporte nativo para soft deletes
- Plugin ecosystem (sharding, caching, etc.)
- Gran comunidad y documentación

### Desventajas
- Overhead de reflexión en cada query
- Queries complejas requieren `db.Raw()` (pierde la ventaja del ORM)
- Debugging difícil: el SQL generado no siempre es obvio
- Puede generar N+1 queries si no se usa Preload correctamente
- Las migraciones automáticas no cubren ALTER de columnas existentes

### Cuándo usarlo
- Desarrollo rápido de APIs CRUD
- Proyectos donde la velocidad de desarrollo importa más que la de ejecución
- Equipos que prefieren no escribir SQL
- Cuando necesitás migraciones y relaciones out-of-the-box

---

## sqlx

**Repositorio**: github.com/jmoiron/sqlx

### Qué es
Extensión de `database/sql` estándar. No genera SQL — vos escribís las queries, sqlx te ayuda a mapear resultados a structs.

### Ejemplo CRUD

```go
// Modelo
type User struct {
    ID    string `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

// Create
db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", "1", "Juan", "juan@test.com")

// Read uno
var user User
db.Get(&user, "SELECT * FROM users WHERE email = $1", "juan@test.com")

// Read lista
var users []User
db.Select(&users, "SELECT * FROM users WHERE active = $1 LIMIT $2 OFFSET $3", true, 10, 0)

// Named queries
db.NamedExec("INSERT INTO users (id, name, email) VALUES (:id, :name, :email)", user)

// Update
db.Exec("UPDATE users SET name = $1 WHERE id = $2", "Juan Carlos", "1")

// Delete
db.Exec("DELETE FROM users WHERE id = $1", "1")
```

### Ventajas
- Control total del SQL generado
- Performance cercana a database/sql puro
- `.Get()` y `.Select()` eliminan el boilerplate de `rows.Scan()`
- Named queries (`:campo`) para mejor legibilidad
- Compatible con cualquier herramienta que use `database/sql`
- Mapeo a structs, maps y slices
- Maduro y extremadamente estable

### Desventajas
- Sin migraciones (necesitás herramientas externas: goose, migrate, etc.)
- Sin manejo de relaciones (JOINs manuales)
- Más código repetitivo que un ORM
- Sin hooks ni lifecycle management
- Errores de SQL se detectan en runtime, no en compilación

### Cuándo usarlo
- Queries complejas donde necesitás control total del SQL
- Proyectos con esquemas de DB estables
- Equipos cómodos escribiendo SQL
- Cuando la performance importa pero querés comodidad de mapeo

---

## sqlw

**Repositorio**: github.com/prongbang/sqlw

### Qué es
Wrapper ultra-liviano sobre `database/sql` diseñado para máximo rendimiento. Minimiza allocations y ofrece CRUD directo con structs.

### Ejemplo CRUD

```go
// Modelo
type User struct {
    ID    string `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

// Create
db.Insert(&User{ID: "1", Name: "Juan", Email: "juan@test.com"})

// Read
var user User
db.QueryRow("SELECT * FROM users WHERE email = $1", "juan@test.com").Scan(&user)

// Update
db.Update(&User{ID: "1", Name: "Juan Carlos", Email: "juan@test.com"})

// Delete
db.Delete(&User{ID: "1"})
```

### Ventajas
- Mínimas allocations de memoria (optimizado para GC)
- CRUD automático con punteros a structs
- API extremadamente simple y predecible
- Etiquetas de struct configurables en runtime (`db.SetTag()`)
- Overhead casi nulo sobre database/sql puro
- Ideal para hot paths y microservicios de alta carga

### Desventajas
- Comunidad muy pequeña — menor soporte y documentación
- Sin named queries avanzadas
- Sin migraciones
- Sin manejo de relaciones
- Menos battle-tested que sqlx en producción a escala
- Riesgo de abandono del proyecto (bajo mantenimiento)

### Cuándo usarlo
- Microservicios donde cada microsegundo cuenta
- Hot paths con miles de queries por segundo
- Cuando ya tenés el SQL optimizado y solo necesitás mapeo rápido
- Proyectos pequeños donde la simplicidad extrema es una ventaja

---

## Benchmark orientativo

| Operación | GORM | sqlx | sqlw |
|-----------|------|------|------|
| INSERT simple | ~150μs | ~50μs | ~45μs |
| SELECT 1 row | ~120μs | ~40μs | ~35μs |
| SELECT 100 rows | ~800μs | ~200μs | ~180μs |
| Allocations/op | Alto | Medio | Bajo |

> **Nota**: Estos valores son orientativos basados en benchmarks públicos. Los números reales dependen del hardware, driver, esquema y complejidad de las queries. En la mayoría de APIs, la latencia de red a la DB domina sobre el overhead de la librería.

---

## Decisión para team-ai-toolkit

| Criterio | Elección |
|----------|----------|
| **Default del toolkit** | GORM — alineado con tich-cronos, desarrollo rápido |
| **Si necesitás performance** | sqlw — para servicios de alta carga |
| **Si necesitás control SQL** | sqlx — estándar de la industria, máxima estabilidad |

### Arquitectura de soporte

```
dbconn/
├── postgres.go          # Interfaz común de conexión
├── gorm/                # Implementación GORM (actual)
│   └── connector.go
├── sqlx/                # Futuro: si se necesita
│   └── connector.go
└── sqlw/                # Futuro: si se necesita
    └── connector.go

transactions/
├── transactor.go        # Interfaz Persistor/Transactor
├── gorm/                # Implementación actual
│   └── transactor.go
├── sqlx/                # Futuro
│   └── transactor.go
└── sqlw/                # Futuro
    └── transactor.go
```

Los repositorios son específicos a cada librería — no se abstraen queries detrás de interfaces genéricas. Cada servicio elige su librería según sus necesidades.
