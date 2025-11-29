# Selenium Testing con Edge en Go

Esta carpeta contiene tests de frontend usando Selenium WebDriver con Microsoft Edge.

## 📋 Requisitos Previos

### 1. Descargar msedgedriver

1. Ve a: https://developer.microsoft.com/en-us/microsoft-edge/tools/webdriver/
2. Descarga la versión que coincida con tu Edge:
   ```bash
   edge --version  # o microsoft-edge --version
   ```
3. Guarda `msedgedriver` en `/usr/local/bin/` o agrega su carpeta al PATH

Verificar instalación:
```bash
msedgedriver --version
```

### 2. Verificar que Edge está instalado

```bash
which edge  # o microsoft-edge
```

---

## 🚀 Cómo Ejecutar

### Paso 1: Iniciar WebDriver (en una terminal)

```bash
# En la carpeta donde está msedgedriver
msedgedriver --port=4444 --verbose
```

Deberías ver:
```
Starting WebDriver server...
Server running on port 4444
```

### Paso 2: Ejecutar tu aplicación frontend (en otra terminal)

```bash
# En la carpeta de tu frontend (React, Vue, etc)
npm start  # o tu comando de dev
```

Espera a que esté en `http://localhost:3000`

### Paso 3: Ejecutar los tests (en otra terminal)

```bash
# En la raíz del proyecto Go
go test -tags=selenium ./tests/selenium/... -v

# O un test específico
go test -tags=selenium ./tests/selenium/... -v -run TestLoginFlow

# Tests exploratorios con screenshots automáticos
go test -tags=selenium ./tests/selenium/... -v -run Exploratory
```

---

## 📸 Screenshots Automáticos

Los tests ahora capturan automáticamente screenshots en cada paso importante. Las imágenes se guardan en:

```
tests/selenium/screenshots/YYYY-MM-DD_HH-MM-SS_test_name/
```

### Cómo funcionan los screenshots

**Tests con screenshots automáticos:**
- Todos los tests en `frontend_test.go` capturan screenshots
- Los tests exploratorios en `exploratory_test.go` capturan más detalles

**Ejemplo de uso:**

```go
// Los screenshots se toman automáticamente
driver, err := NewEdgeDriver("http://localhost:3000")
defer driver.Close()

driver.NavigateTo("/login")
driver.Screenshot("login_page_loaded")  // Captura screenshot

driver.SendKeys("input[type='email']", "test@example.com")
driver.Screenshot("email_entered")      // Captura screenshot
```

### Gestionar Screenshots

Usa el script auxiliar `manage-screenshots.sh`:

```bash
# Listar todos los directorios de screenshots
./manage-screenshots.sh list

# Ver el último directorio (abre en explorer)
./manage-screenshots.sh view-latest
./manage-screenshots.sh open-last

# Ver un test específico
./manage-screenshots.sh view login

# Ver el tamaño de todos los screenshots
./manage-screenshots.sh size

# Comprimir todos los screenshots
./manage-screenshots.sh zip

# Limpiar screenshots antiguos (más de 7 días)
./manage-screenshots.sh clean-old 7

# Eliminar todos los screenshots
./manage-screenshots.sh clean
```

Hacer el script ejecutable:
```bash
chmod +x tests/selenium/manage-screenshots.sh
```

---

## 📝 Estructura de Tests

```
tests/selenium/
├── setup.go                  # Configuración de EdgeDriver con screenshots
├── frontend_test.go         # Tests funcionales (login, signup, etc)
├── exploratory_test.go      # Tests exploratorios con screenshots detallados
├── manage-screenshots.sh    # Script para gestionar screenshots
└── README.md               # Este archivo
```

---

## 🎯 Tests Incluidos

### Tests Funcionales (frontend_test.go)

#### 1. **TestLoginFlow**
- Navega a `/login`
- Completa el formulario de login
- Captura screenshots en cada paso
- Verifica que se redirige al dashboard

#### 2. **TestSignupFlow**
- Navega a `/signup`
- Completa el formulario de registro
- Captura screenshots en cada paso
- Verifica que se redirige correctamente

#### 3. **TestCreateProject**
- Hace login
- Navega a proyectos
- Crea un nuevo proyecto
- Captura 8 screenshots del flujo completo
- Verifica que aparece en la lista

#### 4. **TestNavigationMenu**
- Verifica que el menú lateral existe
- Prueba la navegación a diferentes secciones
- Captura screenshots de cada sección

### Tests Exploratorios (exploratory_test.go)

Estos tests son más libres y capturan más detalles:

#### 1. **TestExploratoryLoginFlow**
- Prueba libre del flujo de login
- Captura screenshot en cada campo completado

#### 2. **TestExploratoryDashboard**
- Explora todas las secciones del dashboard
- Captura screenshots de cada área

#### 3. **TestExploratoryProjectCreation**
- Prueba libre de creación de proyectos
- Maneja errores sin fallar el test
- Captura screenshots de cada intento

#### 4. **TestExploratorySignup**
- Prueba libre de registro
- Explora el formulario completo

#### 5. **TestExploratoryFullUserJourney**
- Recorre todo el flujo del usuario
- Desde home hasta todas las páginas principales
- Captura screenshots de todo el viaje

---

## 🔧 Personalizar Selectores CSS

Los tests usan selectores CSS. Asegúrate de que coincidan con tu HTML:

```go
// Cambiar estos selectores según tu HTML
driver.SendKeys("input[type='email']", "test@example.com")
driver.Click("button[type='submit']")
driver.FindElement(".dashboard-container")
```

**Cómo encontrar los selectores correctos:**

1. Abre Edge Developer Tools (F12)
2. Inspecciona el elemento
3. Copia el selector CSS más específico
4. Reemplaza en el test

---

## 📁 Estructura de Screenshots

Cada ejecución de test crea un directorio con este formato:

```
screenshots/
├── 2024-11-29_15-30-45_login/
│   ├── 01_login_page_loaded.png
│   ├── 02_email_entered.png
│   ├── 03_password_entered.png
│   ├── 04_submit_button_clicked.png
│   └── 05_dashboard_loaded_after_login.png
├── 2024-11-29_15-31-20_project_creation/
│   ├── 01_projects_page_loaded.png
│   ├── 02_new_project_dialog_opened.png
│   ├── 03_project_name_filled.png
│   ├── 04_description_filled.png
│   └── 05_project_created.png
└── 2024-11-29_15-32-10_full_journey/
    ├── 01_home_page.png
    ├── 02_login_page.png
    ├── 03_after_login_attempt.png
    └── ...
```

---

## 🐛 Troubleshooting

### Error: "Failed to connect to Edge WebDriver"
- Verifica que `msedgedriver --port=4444` está corriendo
- Comprueba el puerto (por defecto 4444)

### Error: "element not found"
- Los selectores CSS no coinciden con tu HTML
- Usa las DevTools de Edge para inspeccionar elementos
- Aumenta los timeouts en `setup.go`
- Las screenshots te ayudarán a ver qué está pasando

### Edge se cierra rápidamente después de abrir
- Es normal en modo headless
- Agrega `--headless` solo si quieres modo sin interfaz
- Los tests se ejecutan en Edge real, no headless

### Screenshots no se guardan
- Verifica que tienes permisos de escritura en `tests/selenium/`
- El directorio se crea automáticamente
- Comprueba el path relativo desde donde ejecutas `go test`

---

## 💡 Tips

1. **Esperas**: Usa `time.Sleep()` entre acciones importantes
2. **Screenshots automáticos**: Llama a `driver.Screenshot("descripción")` después de acciones
3. **Pruebas exploratorias**: Usa los tests en `exploratory_test.go` como base para nuevas pruebas
4. **Debugging**: Las screenshots automáticas ayudan a ver exactamente qué pasó en cada paso
5. **Velocidad**: Los tests son lentos porque interactúan con el navegador real - es normal

---

## 📚 Recursos

- [Selenium WebDriver Go Binding](https://github.com/tebeka/selenium)
- [Edge WebDriver Documentation](https://learn.microsoft.com/en-us/microsoft-edge/webdriver-chromium/)
- [CSS Selectors Guide](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors)

---

## ⚙️ Configuración Avanzada

### Ejecutar en modo headless (sin GUI)

En `setup.go`, modifica `NewEdgeDriver`:

```go
caps := selenium.Capabilities{
    "browserName": "MicrosoftEdge",
    "ms:edgeOptions": map[string]interface{}{
        "args": []string{
            "--headless",                    // Agregar esta línea
            "--start-maximized",
            "--disable-notifications",
        },
    },
}
```

### Cambiar tamaño de ventana

```go
// Después de NewEdgeDriver
driver.GetDriver().SetWindowSize(1920, 1080)
```

### Agregar delay antes de cada acción

```go
time.Sleep(500 * time.Millisecond)
```

### Crear screenshots con nombres personalizados

```go
// Usar el driver con nombre personalizado
driver, _ := NewEdgeDriverWithName("http://localhost:3000", "mi_test")

// Las screenshots se guardarán en:
// tests/selenium/screenshots/YYYY-MM-DD_HH-MM-SS_mi_test/
```

---

## 🚦 Build Tag

Los tests requieren el build tag `selenium`:

```bash
go test -tags=selenium ./tests/selenium/...
```

Esto se define en la primera línea de cada archivo de test:
```go
//go:build selenium
```

---

## 🎓 Ejemplo Completo

```go
//go:build selenium

package selenium

import (
    "testing"
    "time"
    "github.com/stretchr/testify/require"
)

func TestMyFeature(t *testing.T) {
    // Crear driver con screenshots automáticos
    driver, err := NewEdgeDriver("http://localhost:3000")
    require.NoError(t, err)
    defer driver.Close()

    // 1. Navegar
    err = driver.NavigateTo("/login")
    require.NoError(t, err)
    driver.Screenshot("page_loaded")

    // 2. Interactuar
    err = driver.SendKeys("input[type='email']", "user@example.com")
    require.NoError(t, err)
    driver.Screenshot("email_filled")

    // 3. Verificar
    element, err := driver.FindElement(".welcome-message")
    require.NoError(t, err)
    driver.Screenshot("final_state")

    t.Logf("Screenshots guardadas en: %s", driver.GetScreenshotDir())
}
```
