# 🚀 Quick Start - Testcontainers

## Inicio Rápido

### 1. Verificar que Docker está corriendo
```bash
docker ps
```

### 2. Ejecutar todos los tests
```bash
cd /home/metadandy/Projects/cuent_ai/cuent-ai-core
go test -tags=containers ./tests/containers/... -v
```

### 3. ¿Qué pasó?
- ✅ Se inició un contenedor PostgreSQL
- ✅ Se ejecutaron las migraciones
- ✅ Se corrieron 44 tests
- ✅ El contenedor se limpió automáticamente

## 📊 Archivos Creados

### Setup
- `tests/containers/setup/testcontainers.go` - Configuración de PostgreSQL container

### Fixtures (Datos de Prueba)
- `tests/containers/fixtures/user.go` - Crear usuarios de prueba
- `tests/containers/fixtures/project.go` - Crear proyectos de prueba
- `tests/containers/fixtures/subscription.go` - Crear suscripciones de prueba

### Tests de Repository (Acceso a Datos)
- `tests/containers/integration/repository/user_repository_test.go` - 7 tests
- `tests/containers/integration/repository/project_repository_test.go` - 7 tests

### Tests de Service (Lógica de Negocio)
- `tests/containers/integration/service/user_service_test.go` - 10 tests
- `tests/containers/integration/service/project_service_test.go` - 5 tests

### Tests de Handler (API HTTP)
- `tests/containers/integration/handler/user_handler_test.go` - 8 tests
- `tests/containers/integration/handler/project_handler_test.go` - 7 tests

### Documentación
- `tests/containers/README.md` - Documentación completa
- `tests/containers/IMPLEMENTATION_SUMMARY.md` - Este resumen

## 🎯 Casos de Uso

### Ejecutar solo tests de User
```bash
go test -tags=containers -run User ./tests/containers/... -v
```

### Ejecutar solo un test específico
```bash
go test -tags=containers -run TestUserRepository_FindByEmail ./tests/containers/... -v
```

### Ejecutar solo Repository tests
```bash
go test -tags=containers ./tests/containers/integration/repository/... -v
```

### Ejecutar solo Service tests
```bash
go test -tags=containers ./tests/containers/integration/service/... -v
```

### Ejecutar solo Handler tests
```bash
go test -tags=containers ./tests/containers/integration/handler/... -v
```

## 📈 44 Tests Implementados

### User Module (25 tests)
- **Repository:** FindByEmail, FindById, Create, Update, SoftDelete, Restore, FindAll
- **Service:** FindById, SignUp, SignUp validation, SignIn, FindAll
- **Handler:** SignUp, SignIn, GetProfile, GetProfile (sin auth), FindById

### Project Module (19 tests)
- **Repository:** Create, FindById, Update, SoftDelete, Restore, FindAll, Pagination
- **Service:** FindByID, Create, Update, Delete, GetAll
- **Handler:** CreateProject, GetProject, GetProjects, UpdateProject, DeleteProject, Auth checks

## ✨ Características

✅ **PostgreSQL Real** - No mocks, base de datos real en containers  
✅ **Aislado** - Cada test en su propio contenedor  
✅ **Limpio** - Fixtures reutilizables y datos de prueba limpios  
✅ **Completo** - 3 niveles: Repository, Service, Handler  
✅ **Documentado** - README con ejemplos y guía  
✅ **Rápido** - Containers efímeros sin cleanup manual  
✅ **CI/CD Ready** - Funciona en cualquier máquina con Docker  

## 🔧 Estructura del Proyecto

```
tests/
├── containers/               ← AQUÍ están los testcontainers
│   ├── setup/
│   ├── fixtures/
│   ├── integration/
│   │   ├── repository/
│   │   ├── service/
│   │   └── handler/
│   ├── README.md
│   └── IMPLEMENTATION_SUMMARY.md
├── unit/                     ← Tests unitarios antiguos (sin cambios)
└── README.md
```

## 🎯 Puntos Clave

1. **Sin cambios al código existente** - Todo en `/tests/containers/`
2. **Build tag `containers`** - Ejecuta solo con `-tags=containers`
3. **Cada test limpio** - Contenedor nuevo para cada test
4. **Fixtures reutilizables** - Crear datos en segundos
5. **3 niveles de tests** - Repository → Service → Handler

## 📞 Problemas Comunes

**Error: "docker: command not found"**
- Instala Docker: https://docs.docker.com/get-docker/

**Error: "permission denied"**
- Agrega tu usuario a docker: `sudo usermod -aG docker $USER`

**Tests muy lentos**
- Normal en la primera ejecución (descarga imagen PostgreSQL)
- Siguientes ejecuciones son más rápidas

**Tests fallan con timeout**
- Intenta ejecutar con menos paralelismo: `go test -parallel 1`

## 🚀 Próximas Mejoras

Cuando quieras, podemos:
- [ ] Agregar tests para Script module
- [ ] Agregar tests para Asset module
- [ ] Agregar tests de flujos end-to-end
- [ ] Integrar en GitHub Actions
- [ ] Agregar coverage reports

---

**Todo listo para usar!** 🎉

Ejecuta: `go test -tags=containers ./tests/containers/... -v`
