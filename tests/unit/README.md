# Pruebas Unitarias

Este directorio contiene las pruebas unitarias del proyecto, separadas de las pruebas de integración y E2E.

## Estructura

```
tests/unit/
├── user/
│   └── user_service_test.go
├── subscription/
│   └── subscription_service_test.go
├── project/
│   └── project_service_test.go
└── validation/
    └── validation_test.go
```

## Ejecutar las Pruebas

### Ejecutar todas las pruebas unitarias

```bash
go test -tags=unit ./tests/unit/...
```

### Ejecutar pruebas de un módulo específico

```bash
# Pruebas de User
go test -tags=unit ./tests/unit/user/...

# Pruebas de Subscription
go test -tags=unit ./tests/unit/subscription/...

# Pruebas de Project
go test -tags=unit ./tests/unit/project/...

# Pruebas de Validation
go test -tags=unit ./tests/unit/validation/...
```

### Ejecutar con cobertura

```bash
go test -tags=unit -cover ./tests/unit/...
```

### Ejecutar con cobertura detallada

```bash
go test -tags=unit -coverprofile=coverage.out ./tests/unit/...
go tool cover -html=coverage.out
```

### Ejecutar con verbose (ver detalles)

```bash
go test -tags=unit -v ./tests/unit/...
```

### Ejecutar un test específico

```bash
go test -tags=unit -v ./tests/unit/user/... -run TestEmailValidation
```

## Módulos Cubiertos

1. **User Service**: 
   - Validación de email (formato RFC 5322)
   - Validación de contraseña (mínimo 8 caracteres)
   - Normalización de email (trim y lowercase)
   - Manejo de errores del servicio

2. **Subscription Service**: 
   - Búsqueda de suscripciones (FindAll, FindByID)
   - Paginación de resultados
   - Manejo de errores del repositorio

3. **Project Service**: 
   - CRUD de proyectos (FindByID, Update, SoftDelete, Restore, FindAll)
   - Validación de operaciones
   - Manejo de errores

4. **Validation Service**: 
   - Validación de líneas TTS con límite de caracteres (200)
   - Validación de líneas SFX (sin límite)
   - Conteo de caracteres Unicode
   - Casos límite y condiciones de borde

## Características

- ✅ Usan build tags `//go:build unit` para separación
- ✅ Mocks de repositorios para aislamiento completo
- ✅ Cobertura de casos exitosos y de error
- ✅ Tests independientes y rápidos
- ✅ Sin dependencias externas (DB, APIs)
- ✅ Validación de condiciones límite

## Ejemplo de Salida

Al ejecutar `go test -tags=unit -v ./tests/unit/...`, verás algo como:

```
=== RUN   TestEmailValidation
=== RUN   TestEmailValidation/Valid_email_-_standard
=== RUN   TestEmailValidation/Invalid_-_no_@
...
PASS
ok      github.com/MetaDandy/cuent-ai-core/tests/unit/user    0.123s

=== RUN   TestService_FindAll
=== RUN   TestService_FindAll/Success_-_Find_All_Subscriptions
...
PASS
ok      github.com/MetaDandy/cuent-ai-core/tests/unit/subscription    0.045s
```

## Notas

- Las pruebas unitarias están completamente aisladas y no requieren base de datos ni servicios externos
- Todos los repositorios están mockeados para garantizar velocidad y confiabilidad
- Los tests pueden ejecutarse en paralelo sin problemas
- El build tag `unit` permite ejecutar solo estas pruebas sin incluir las de integración o E2E

