# Pruebas Funcionales (E2E)

Estas pruebas E2E están diseñadas para validar la API del proyecto en un entorno en ejecución. Todas las pruebas usan la etiqueta de build `e2e` y se ejecutan a través de `go test`.

Requisitos:
- Un servidor API corriendo (local o en docker) que responda a `TEST_API_URL` (por defecto `http://localhost:8000/api/v1`).
- Go instalado para ejecutar las pruebas.

Variables de entorno útiles:
- `TEST_API_URL` - URL base de la API (por ejemplo: `http://localhost:8000/api/v1`).
- `TEST_ADMIN_EMAIL` - Email del usuario administrador.
- `TEST_ADMIN_PASSWORD` - Contraseña del usuario administrador.

Cómo ejecutar las pruebas:
```powershell
# Ejecuta todas las pruebas E2E y genera un JSON de salida
.\n+\tests\functional\run_tests.ps1
```

Compilar solo los tests (sin ejecutar):
```powershell
go test -c -tags=e2e ./tests/functional/...
```

Ejecutar una prueba específica:
```powershell
go test -run TestAuthFlow_LoginSuccess -tags=e2e ./tests/functional -v
```

Notas de diseño:
- Las pruebas generan datos únicos para evitar colisiones.
- Evita usar entornos compartidos sin limpieza; se recomienda ejecutar en un entorno de pruebas o contenedores efímeros.
# Pruebas Funcionales E2E

Pruebas End-to-End para validar flujos completos de negocio de Cuent-AI Core.

## Estructura

```
tests/functional/
├── config/test_config.go          # Configuración
├── helpers/
│   ├── http_client.go             # Cliente HTTP
│   └── test_data.go               # Estructuras de datos
├── auth_flow_test.go              # Autenticación (4 casos)
├── project_flow_test.go           # Proyectos (5 casos)
├── script_generation_test.go      # Scripts (4 casos)
└── authorization_test.go          # Autorización (5 casos)
```

## Ejecución

```powershell
# Ejecutar todas las pruebas
go test -v -tags=e2e ./tests/functional/...

# Ejecutar módulo específico
go test -v -tags=e2e ./tests/functional/ -run TestAuthFlow
```

## Generar Reporte HTML

```powershell
# Instalar herramienta (primera vez)
go install github.com/vakenbolt/go-test-report@latest

# Generar reporte
go test -v -tags=e2e ./tests/functional/... -json | go-test-report -o tests/functional/reports/report.html
```

## Prerequisitos

1. Servidor API corriendo: `docker-compose -f docker-compose.dev.yml up`
2. Variables de entorno configuradas en `.env`

## Casos de Prueba

Si necesitas una lista de casos o una matriz de pruebas, añade un fichero `CASOS_PRUEBA.md` con los detalles.
