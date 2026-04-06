# Testing Guide

## Requisitos

- Go 1.26+
- PostgreSQL 16 (para integration tests)
- golangci-lint (para linting)

## Correr tests

```bash
# Todos
go test ./...

# Con verbose
go test -v ./...

# Con coverage
go test -cover ./...

# Coverage report HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Un paquete especifico
go test -v ./src/core/usecases/coordination/...
go test -v ./src/entrypoints/rest/coordination/...
go test -v ./src/repositories/coordination/...
```

## Estructura de tests

```
src/core/usecases/
    coordination/
        create_document.go
        create_document_test.go      # Unit test: mock providers

src/entrypoints/rest/
    coordination/
        create.go
        create_test.go               # Test: mock usecase

src/repositories/
    coordination/
        get_document.go
        get_document_test.go         # Integration test: PostgreSQL real
```

## Convenciones

### Patron AAA (Arrange, Act, Assert)

Separar las tres secciones con lineas en blanco. No poner comentarios `// Arrange`, `// Act`, `// Assert`.

```go
func TestCreateDocument_Execute_Success(t *testing.T) {
    repo := new(mocks.CoordinationProvider)
    repo.On("CreateDocument", mock.Anything, mock.Anything).Return(int64(1), nil)
    uc := NewCreateDocument(repo)

    id, err := uc.Execute(ctx, CreateDocumentRequest{Name: "Test", OrgID: 1, AreaID: 1})

    require.NoError(t, err)
    assert.Equal(t, int64(1), id)
}
```

### Assertions

Siempre usar testify. Nunca `if + t.Fatal`.

```go
require.NoError(t, err)              // Precondiciones (falla inmediatamente)
assert.Equal(t, expected, actual)    // Verificaciones (continua ejecutando)
assert.True(t, condition)
assert.Len(t, slice, 3)
```

### Comparacion de errores

```go
// Bien
assert.True(t, errors.Is(err, providers.ErrValidation))
assert.Equal(t, "validation error: name is required", err.Error())

// Mal
assert.Contains(t, err.Error(), "name")
```

### Naming

```go
// Formato: Test<Struct>_<Method>_<Scenario>
func TestCreateDocument_Execute_Success(t *testing.T) {}
func TestCreateDocument_Execute_ValidationError(t *testing.T) {}
func TestCoordinationRepo_GetDocument_NotFound(t *testing.T) {}
```

### Table-driven tests

Para escenarios repetitivos:

```go
func TestCreateDocumentRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateDocumentRequest
        wantErr bool
    }{
        {"valid", CreateDocumentRequest{Name: "Test", OrgID: 1, AreaID: 1}, false},
        {"empty name", CreateDocumentRequest{Name: "", OrgID: 1, AreaID: 1}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Testing por capa

| Capa | Que se testea | Como | DB |
|------|--------------|------|-----|
| usecases | Logica de negocio | Mock providers | No |
| handlers | Parsing HTTP, error mapping | Mock usecases | No |
| repositories | Queries, mapeo | PostgreSQL real | Si |

## Coverage target: 90%
