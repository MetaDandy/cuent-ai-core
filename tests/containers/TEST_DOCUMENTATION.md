# Testcontainers - Documentación Completa

## Introducción

Este documento describe la suite completa de **41 tests de integración** que utiliza **Testcontainers** para probar las tres capas de la aplicación cuent-ai-core: **Repositorio**, **Servicio** y **Handler (HTTP)**.

Los tests se ejecutan con: `go test -tags=containers ./tests/containers/... -v`

---

## 🏗️ Arquitectura de Tests

La suite está organizada en **3 capas de prueba**:

```
┌─────────────────────────────────────────────┐
│         Handler Layer (13 tests)            │  ← HTTP Endpoints
├─────────────────────────────────────────────┤
│         Service Layer (13 tests)            │  ← Business Logic
├─────────────────────────────────────────────┤
│        Repository Layer (15 tests)          │  ← Database Access
├─────────────────────────────────────────────┤
│   Testcontainers (PostgreSQL 15-alpine)    │  ← Database Container
└─────────────────────────────────────────────┘
```

**Cada test**:
- ✅ Crea un contenedor PostgreSQL aislado
- ✅ Ejecuta las migraciones de BD
- ✅ Prepara datos de prueba con fixtures
- ✅ Ejecuta la operación a probar
- ✅ Verifica los resultados
- ✅ Limpia el contenedor automáticamente

---

## 📦 Handler Tests (Capa HTTP - 13 tests)

Los tests de **Handler** prueban los **endpoints HTTP** de la API Fiber v2.

### Ubicación
`tests/containers/integration/handler/`

### Tests de Project

#### 1. **TestProjectHandler_CreateProject**
- **Propósito**: Verificar creación de un proyecto a través del endpoint POST
- **Endpoint**: `POST /api/v1/projects`
- **Lo que prueba**:
  - ✅ Autenticación con JWT válido
  - ✅ Validación de datos requeridos (name, description, user_id)
  - ✅ Respuesta HTTP 201 (Created)
  - ✅ Proyecto guardado correctamente en BD
  - ✅ Relación usuario-proyecto creada

**Datos de entrada**:
```json
{
  "name": "New Project",
  "description": "Test project",
  "user_id": "uuid-del-usuario"
}
```

---

#### 2. **TestProjectHandler_GetProject**
- **Propósito**: Obtener detalles de un proyecto específico
- **Endpoint**: `GET /api/v1/projects/:id`
- **Lo que prueba**:
  - ✅ Autenticación requerida
  - ✅ Respuesta HTTP 200 (OK)
  - ✅ Proyecto recuperado correctamente
  - ✅ Datos del proyecto son precisos

---

#### 3. **TestProjectHandler_GetProjects**
- **Propósito**: Listar todos los proyectos del usuario autenticado
- **Endpoint**: `GET /api/v1/projects`
- **Lo que prueba**:
  - ✅ Autenticación requerida
  - ✅ Paginación funciona
  - ✅ Respuesta HTTP 200
  - ✅ Solo se devuelven proyectos no eliminados (soft delete)

---

#### 4. **TestProjectHandler_UpdateProject**
- **Propósito**: Actualizar un proyecto existente
- **Endpoint**: `PATCH /api/v1/projects/:id`
- **Lo que prueba**:
  - ✅ Autenticación requerida
  - ✅ Cambios se persisten en BD
  - ✅ Respuesta HTTP 200 con datos actualizados
  - ✅ Validación de campos opcionales

---

#### 5. **TestProjectHandler_DeleteProject**
- **Propósito**: Eliminar un proyecto (soft delete)
- **Endpoint**: `DELETE /api/v1/projects/:id`
- **Lo que prueba**:
  - ✅ Autenticación requerida
  - ✅ Respuesta HTTP 204 (No Content)
  - ✅ Proyecto no aparece en búsquedas normales
  - ✅ Soft delete: registro permanece en BD con DeletedAt

---

#### 6. **TestProjectHandler_CreateProject_MissingToken**
- **Propósito**: Validar rechazo sin autenticación
- **Endpoint**: `POST /api/v1/projects`
- **Lo que prueba**:
  - ✅ Middleware JWT rechaza requests sin token
  - ✅ Respuesta HTTP 401 (Unauthorized)
  - ✅ Endpoint protegido correctamente

---

### Tests de User

#### 7. **TestUserHandler_SignUp**
- **Propósito**: Registrar un nuevo usuario
- **Endpoint**: `POST /api/v1/users/sign-up`
- **Lo que prueba**:
  - ✅ Validación de email único
  - ✅ Hash de contraseña correcto
  - ✅ Respuesta HTTP 201 (Created)
  - ✅ Usuario guardado en BD
  - ✅ JWT generado

---

#### 8. **TestUserHandler_SignUp_InvalidEmail**
- **Propósito**: Rechazar emails inválidos
- **Endpoint**: `POST /api/v1/users/sign-up`
- **Lo que prueba**:
  - ✅ Validación de formato de email
  - ✅ Respuesta HTTP 400 (Bad Request)
  - ✅ Usuario NO es creado

---

#### 9. **TestUserHandler_SignIn**
- **Propósito**: Autenticar usuario existente
- **Endpoint**: `POST /api/v1/users/sign-in`
- **Lo que prueba**:
  - ✅ Validación de credenciales correctas
  - ✅ Respuesta HTTP 201 con JWT
  - ✅ Token válido para requests posteriores

---

#### 10. **TestUserHandler_SignIn_InvalidPassword**
- **Propósito**: Rechazar contraseña incorrecta
- **Endpoint**: `POST /api/v1/users/sign-in`
- **Lo que prueba**:
  - ✅ Comparación segura de hashes
  - ✅ Respuesta HTTP 401 (Unauthorized)
  - ✅ No se devuelve JWT

---

#### 11. **TestUserHandler_GetProfile**
- **Propósito**: Obtener perfil del usuario autenticado
- **Endpoint**: `GET /api/v1/users/profile`
- **Lo que prueba**:
  - ✅ Extrae user_id del JWT
  - ✅ Respuesta HTTP 201 (OK)
  - ✅ Datos del perfil completos

---

#### 12. **TestUserHandler_GetProfile_MissingToken**
- **Propósito**: Rechazar acceso sin token
- **Endpoint**: `GET /api/v1/users/profile`
- **Lo que prueba**:
  - ✅ Middleware JWT requerido
  - ✅ Respuesta HTTP 401 (Unauthorized)

---

#### 13. **TestUserHandler_FindById**
- **Propósito**: Obtener usuario por ID
- **Endpoint**: `GET /api/v1/users/:id`
- **Lo que prueba**:
  - ✅ Autenticación requerida
  - ✅ Respuesta HTTP 201 (Created)
  - ✅ Usuario encontrado correctamente

---

## 🔧 Service Tests (Lógica de Negocio - 13 tests)

Los tests de **Service** prueban la **lógica de negocio** sin HTTP.

### Ubicación
`tests/containers/integration/service/`

### Tests de Project Service

#### 1. **TestProjectService_FindByID**
- **Propósito**: Buscar proyecto por ID
- **Lo que prueba**:
  - ✅ Consulta GORM funciona
  - ✅ Relaciones cargadas (scripts, etc)
  - ✅ Soft delete respetado
  - ✅ No hay errores de tipo

---

#### 2. **TestProjectService_Create**
- **Propósito**: Crear nuevo proyecto con validaciones
- **Lo que prueba**:
  - ✅ Validación de datos requeridos
  - ✅ Generación de UUID
  - ✅ Estado inicial (PENDING)
  - ✅ Timestamps correctos
  - ✅ Relación con usuario creada
  - ✅ Proyecto en BD con datos correctos

---

#### 3. **TestProjectService_Update**
- **Propósito**: Actualizar campos de proyecto
- **Lo que prueba**:
  - ✅ Actualización parcial (null fields ignorados)
  - ✅ UpdatedAt actualizado
  - ✅ Cambios persistidos en BD
  - ✅ Validación de punteros

---

#### 4. **TestProjectService_Delete**
- **Propósito**: Soft delete de proyecto
- **Lo que prueba**:
  - ✅ DeletedAt es establecido
  - ✅ Proyecto oculto de queries normales
  - ✅ Recuperable con Unscoped

---

#### 5. **TestProjectService_GetAll**
- **Propósito**: Listar proyectos con paginación
- **Lo que prueba**:
  - ✅ Paginación (limit, offset)
  - ✅ Total correcto
  - ✅ Solo proyectos no eliminados
  - ✅ Orden por usuario (si aplica)

---

### Tests de User Service

#### 6. **TestUserService_FindById**
- **Propósito**: Buscar usuario por ID
- **Lo que prueba**:
  - ✅ Consulta con relaciones (subscriptions)
  - ✅ Usuario completo recuperado

---

#### 7. **TestUserService_SignUp**
- **Propósito**: Registrar usuario con validaciones
- **Lo que prueba**:
  - ✅ Email único (error si existe)
  - ✅ Contraseña hasheada (nunca en plain)
  - ✅ Validación de formato email
  - ✅ Usuario creado en BD

---

#### 8. **TestUserService_SignUp_InvalidEmail**
- **Propósito**: Rechazar emails inválidos
- **Lo que prueba**:
  - ✅ Regex de validación email
  - ✅ Error devuelto correctamente

---

#### 9. **TestUserService_SignUp_WeakPassword**
- **Propósito**: Rechazar contraseñas débiles
- **Lo que prueba**:
  - ✅ Validación de fuerza (min length, caracteres)
  - ✅ Error informativo

---

#### 10. **TestUserService_SignUp_EmailTaken**
- **Propósito**: Rechazar email duplicado
- **Lo que prueba**:
  - ✅ Constraint UNIQUE en BD respetado
  - ✅ Error apropiado devuelto

---

#### 11. **TestUserService_SignIn**
- **Propósito**: Autenticar usuario
- **Lo que prueba**:
  - ✅ Búsqueda por email
  - ✅ Hash verificado correctamente
  - ✅ JWT generado
  - ✅ Claims incluyen sub (user_id) y email

---

#### 12. **TestUserService_SignIn_InvalidPassword**
- **Propósito**: Rechazar contraseña incorrecta
- **Lo que prueba**:
  - ✅ Comparación segura (bcrypt)
  - ✅ No se devuelve info del usuario

---

#### 13. **TestUserService_FindAll**
- **Propósito**: Listar todos los usuarios
- **Lo que prueba**:
  - ✅ Paginación funciona
  - ✅ Total correcto
  - ✅ Soft delete respetado
  - ✅ Usuarios con emails únicos

---

## 💾 Repository Tests (Acceso a Datos - 15 tests)

Los tests de **Repository** prueban la **capa de persistencia** (GORM + PostgreSQL).

### Ubicación
`tests/containers/integration/repository/`

### Tests de Project Repository

#### 1. **TestProjectRepository_Create**
- **Propósito**: Guardar proyecto en BD
- **Lo que prueba**:
  - ✅ INSERT correcto
  - ✅ UUID generado
  - ✅ Timestamps establecidos
  - ✅ Relaciones guardadas (user_id)

---

#### 2. **TestProjectRepository_FindById**
- **Propósito**: Recuperar proyecto por ID
- **Lo que prueba**:
  - ✅ SELECT con WHERE id
  - ✅ Soft delete respetado (deleted_at IS NULL)
  - ✅ Relaciones cargadas

---

#### 3. **TestProjectRepository_Update**
- **Propósito**: Actualizar proyecto
- **Lo que prueba**:
  - ✅ UPDATE correcto
  - ✅ UpdatedAt modificado
  - ✅ Otros campos intactos

---

#### 4. **TestProjectRepository_SoftDelete**
- **Propósito**: Marcar proyecto como eliminado
- **Lo que prueba**:
  - ✅ DeletedAt establecido a NOW()
  - ✅ Proyecto ocultado de queries normales
  - ✅ No es eliminado físicamente

---

#### 5. **TestProjectRepository_Restore**
- **Propósito**: Recuperar proyecto eliminado
- **Lo que prueba**:
  - ✅ DeletedAt puesto a NULL
  - ✅ Visible nuevamente en queries normales
  - ✅ DeletedAt.Valid es false

---

#### 6. **TestProjectRepository_FindAll**
- **Propósito**: Listar proyectos sin paginación
- **Lo que prueba**:
  - ✅ SELECT sin WHERE
  - ✅ Total correcto
  - ✅ Soft delete respetado

---

#### 7. **TestProjectRepository_FindAll_Pagination**
- **Propósito**: Listar con LIMIT y OFFSET
- **Lo que prueba**:
  - ✅ LIMIT respetado
  - ✅ OFFSET correcto
  - ✅ Total general vs items devueltos

---

### Tests de User Repository

#### 8. **TestUserRepository_FindByEmail**
- **Propósito**: Buscar usuario por email
- **Lo que prueba**:
  - ✅ WHERE email = ?
  - ✅ Soft delete respetado
  - ✅ Relaciones cargadas (subscriptions)

---

#### 9. **TestUserRepository_FindByEmail_NotFound**
- **Propósito**: Manejar email inexistente
- **Lo que prueba**:
  - ✅ Error "record not found" devuelto
  - ✅ Tipo correcto de error GORM

---

#### 10. **TestUserRepository_FindById**
- **Propósito**: Buscar usuario por ID
- **Lo que prueba**:
  - ✅ SELECT WHERE id = ?
  - ✅ Relaciones preloaded

---

#### 11. **TestUserRepository_Create**
- **Propósito**: Guardar usuario en BD
- **Lo que prueba**:
  - ✅ INSERT correcto
  - ✅ Constraint UNIQUE en email respetado
  - ✅ Timestamps establecidos

---

#### 12. **TestUserRepository_Update**
- **Propósito**: Actualizar usuario
- **Lo que prueba**:
  - ✅ UPDATE correcto
  - ✅ UpdatedAt modificado
  - ✅ Email puede cambiar (si hay cambio)

---

#### 13. **TestUserRepository_SoftDelete**
- **Propósito**: Marcar usuario como eliminado
- **Lo que prueba**:
  - ✅ DeletedAt establecido
  - ✅ Usuario oculto de queries normales

---

#### 14. **TestUserRepository_Restore**
- **Propósito**: Recuperar usuario eliminado
- **Lo que prueba**:
  - ✅ DeletedAt puesto a NULL
  - ✅ DeletedAt.Valid es false después de restore

---

#### 15. **TestUserRepository_FindAll**
- **Propósito**: Listar todos los usuarios
- **Lo que prueba**:
  - ✅ SELECT correcto
  - ✅ Soft delete respetado
  - ✅ Total de registros

---

## 🔧 Infraestructura de Tests

### Setup (Testcontainers)
Archivo: `tests/containers/setup/testcontainers.go`

**Qué hace**:
```go
// 1. Crea contenedor PostgreSQL 15
container := testcontainers.Run(ctx, postgres.ContainerRequest{...})

// 2. Obtiene connection string
dsn := container.ConnectionString(ctx)

// 3. Conecta GORM
db.Open(dsn)

// 4. Ejecuta migraciones
config.Migrate(db)

// 5. Retorna testDB con cleanup automático
defer container.Terminate(ctx)
```

### Fixtures
Archivos: `tests/containers/fixtures/`

**Usuario**:
```go
func CreateTestUser() *model.User {
  return &model.User{
    ID:       uuid.New(),
    Name:     "Test User",
    Email:    "testuser@example.com",
    Password: "hashedpassword123",
  }
}
```

**Proyecto**:
```go
func CreateTestProject(userID uuid.UUID) *model.Project {
  return &model.Project{
    ID:          uuid.New(),
    Name:        "Test Project",
    Description: "Test",
    UserID:      userID,
    State:       model.PENDING,
  }
}
```

---

## 📊 Cobertura de Tests

| Capa | Tests | Cobertura |
|------|-------|-----------|
| **Handler (HTTP)** | 13 | POST, GET, PATCH, DELETE, Auth |
| **Service (Logic)** | 13 | Create, Read, Update, Delete, Validations |
| **Repository (DB)** | 15 | CRUD, Pagination, Soft Delete, Restore |
| **Total** | **41** | Integración completa de BD a HTTP |

---

## 🚀 Ejecución

### Todos los tests
```bash
go test -tags=containers ./tests/containers/... -v
```

### Solo handlers
```bash
go test -tags=containers ./tests/containers/integration/handler -v
```

### Solo services
```bash
go test -tags=containers ./tests/containers/integration/service -v
```

### Solo repositories
```bash
go test -tags=containers ./tests/containers/integration/repository -v
```

### Test específico
```bash
go test -tags=containers ./tests/containers/... -v -run TestProjectHandler_CreateProject
```

---

## ✅ Resultado Final

```
✓ 41 tests PASSING
✓ 0 tests FAILING
✓ Coverage: Repository, Service, Handler layers
✓ Database: PostgreSQL isolated per test
✓ Cleanup: Automatic container termination
```

---

## 🛠️ Tecnologías

- **Testcontainers for Go** v0.40.0 - Contenedores Docker para tests
- **PostgreSQL** 15-alpine - Base de datos
- **GORM** - ORM Go
- **Fiber** v2.52.6 - Framework HTTP
- **testify/assert** - Assertions
- **testify/require** - Validaciones fatales

---

## 📝 Notas Importantes

1. **Build Tag**: Los tests requieren `//go:build containers` para ejecutarse
2. **Docker**: Es necesario tener Docker corriendo
3. **Aislamiento**: Cada test obtiene su propio contenedor PostgreSQL
4. **Soft Delete**: Los tests respetan la lógica de soft delete con DeletedAt
5. **JWT**: Los tests incluyen generación y validación de tokens
6. **Transacciones**: Cada test es independiente sin estado compartido
