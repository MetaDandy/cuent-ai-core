# Testing Strategy - Cuent-AI Core

Este proyecto utiliza una estrategia de testing dividida en tres niveles para asegurar la calidad del software.

## Estructura de Directorios

Las pruebas se encuentran separadas del código fuente principal (`src/`) para mantener limpio el dominio y facilitar la identificación de los tipos de prueba.

```text
tests/
├── unit/           # Pruebas Unitarias
│   └── project/    # Tests específicos del módulo Project
├── integration/    # Pruebas de Integración (API + DB)
└── e2e/            # Pruebas End-to-End (Flujos completos)
```

## 1. Pruebas Unitarias (`tests/unit/`)

**Objetivo:** Probar la lógica de negocio de forma aislada.

- **Dependencias:** Se usan **Mocks** para simular bases de datos o servicios externos.
- **Velocidad:** Muy rápidas.
- **Comando:** `go test ./tests/unit/...`

## 2. Pruebas de Integración (`tests/integration/`)

**Objetivo:** Probar cómo interactúan los componentes (ej: Service con Repository real).

- **Dependencias:** Requieren una base de datos de prueba real (Docker/Testcontainers).
- **Velocidad:** Media.
- **Comando:** `go test ./tests/integration/...`

## 3. Pruebas End-to-End (`tests/e2e/`)

**Objetivo:** Probar la aplicación desde la perspectiva del usuario final (HTTP Requests).

- **Dependencias:** Requieren la aplicación levantada completamente.
- **Velocidad:** Lenta.
- **Comando:** `go test ./tests/e2e/...`

## Identificación por Build Tags (Opcional pero Recomendado)

Para mayor control, podemos agregar "etiquetas" al inicio de los archivos de prueba:

- Unitarios: `//go:build unit`
- Integración: `//go:build integration`
- E2E: `//go:build e2e`

Esto permite ejecutar: `go test -tags=integration ./...`
