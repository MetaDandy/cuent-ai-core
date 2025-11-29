# Quick Start - Selenium Screenshots

## 📸 Para empezar rápidamente

### 1. **Hacer ejecutable el script de gestión:**

```bash
chmod +x tests/selenium/manage-screenshots.sh
```

### 2. **Ejecutar pruebas exploratorias (recomendado para empezar):**

```bash
# En terminal 1: Iniciar WebDriver
msedgedriver --port=4444 --verbose

# En terminal 2: Ejecutar tu frontend
npm start

# En terminal 3: Ejecutar las pruebas
cd /ruta/del/proyecto
go test -tags=selenium ./tests/selenium/... -v -run Exploratory
```

### 3. **Ver los screenshots capturados:**

```bash
# Abrir automáticamente en el explorador de archivos
./tests/selenium/manage-screenshots.sh view-latest

# O listar todos los directorios
./tests/selenium/manage-screenshots.sh list
```

---

## 🎯 Casos de Uso Comunes

### A. Prueba simple con screenshots

```go
func TestSimple(t *testing.T) {
    driver, err := NewEdgeDriver("http://localhost:3000")
    require.NoError(t, err)
    defer driver.Close()

    // Navegar
    driver.NavigateTo("/dashboard")
    driver.Screenshot("dashboard_page")

    // Interactuar
    driver.Click(".menu-projects")
    driver.Screenshot("projects_clicked")

    // Verificar
    _, err = driver.FindElement(".project-item")
    require.NoError(t, err)
    driver.Screenshot("projects_loaded")
}
```

**Ejecutar:**
```bash
go test -tags=selenium ./tests/selenium/... -v -run TestSimple
```

---

### B. Prueba exploratoria libre

Usa los tests en `exploratory_test.go` como base y modifícalos:

```go
func TestMiExploracion(t *testing.T) {
    driver, err := NewEdgeDriverWithName("http://localhost:3000", "mi_test")
    require.NoError(t, err)
    defer driver.Close()

    // El directorio se crea automáticamente en:
    // tests/selenium/screenshots/YYYY-MM-DD_HH-MM-SS_mi_test/

    driver.NavigateTo("/")
    driver.Screenshot("01_home")

    // Puedes navegar libremente sin que falle el test
    _ = driver.Click(".some-button")  // Ignora si falla
    driver.Screenshot("02_after_click")

    // Las screenshots se numeran automáticamente
    // 01_home.png, 02_after_click.png, etc
}
```

---

### C. Capturar flujo completo

```go
func TestCompleteFlow(t *testing.T) {
    driver, err := NewEdgeDriver("http://localhost:3000")
    require.NoError(t, err)
    defer driver.Close()

    steps := []string{
        "login",
        "fill_email",
        "fill_password",
        "submit",
        "dashboard",
        "projects",
        "create_project",
        "success",
    }

    for i, step := range steps {
        // ... hacer la acción ...
        driver.Screenshot(fmt.Sprintf("%02d_%s", i+1, step))
    }

    t.Logf("Screenshots en: %s", driver.GetScreenshotDir())
}
```

---

## 🛠️ Gestionar Screenshots

### Ver estadísticas

```bash
./tests/selenium/manage-screenshots.sh size
```

Output:
```
Tamaño total de screenshots: 45M

Desglose por directorio:
12M	2024-11-29_15-30-45_login/
15M	2024-11-29_15-31-20_project_creation/
18M	2024-11-29_15-32-10_full_journey/
```

### Comprimir para reportes

```bash
./tests/selenium/manage-screenshots.sh zip
# Crea: screenshots_2024-11-29_15-32-45.zip
```

### Limpiar antiguos

```bash
# Elimina screenshots más antiguos que 7 días
./tests/selenium/manage-screenshots.sh clean-old 7

# Eliminar todos
./tests/selenium/manage-screenshots.sh clean
```

---

## 💡 Best Practices

1. **Usa descripciones claras:**
   ```go
   driver.Screenshot("email_validation_error")  // ✓ Claro
   driver.Screenshot("step_3")                   // ✗ Vago
   ```

2. **Captura estados importantes:**
   ```go
   driver.Screenshot("before_action")
   driver.Click(".button")
   driver.Screenshot("after_action")
   driver.Screenshot("final_state")
   ```

3. **Usa `NewEdgeDriverWithName` para tests largos:**
   ```go
   driver, _ := NewEdgeDriverWithName("http://localhost:3000", "user_signup_complete_flow")
   // Crea directorio: .../2024-11-29_15-30-45_user_signup_complete_flow/
   ```

4. **Agrupa screenshots relacionados:**
   ```go
   driver.Screenshot("01_form_empty")
   driver.Screenshot("02_form_validation")
   driver.Screenshot("03_form_filled")
   driver.Screenshot("04_form_submitted")
   driver.Screenshot("05_success_page")
   ```

---

## 🔍 Debugging con Screenshots

Si un test falla, las screenshots te muestran exactamente qué pasó:

1. **Ejecuta el test:**
   ```bash
   go test -tags=selenium ./tests/selenium/... -v -run TestLoginFlow
   ```

2. **Ve los screenshots:**
   ```bash
   ./tests/selenium/manage-screenshots.sh view-latest
   ```

3. **Revisa la secuencia de imágenes** para ver dónde se rompió

---

## 📊 Integración con CI/CD

### GitHub Actions ejemplo

```yaml
- name: Run Selenium Tests
  run: |
    msedgedriver --port=4444 &
    npm start &
    sleep 3
    go test -tags=selenium ./tests/selenium/... -v

- name: Upload Screenshots
  if: failure()
  uses: actions/upload-artifact@v2
  with:
    name: selenium-screenshots
    path: tests/selenium/screenshots/
```

---

## ❓ FAQs

**P: ¿Por qué los tests son lentos?**
R: Porque usan un navegador real. Es normal y necesario para E2E testing.

**P: ¿Puedo correr múltiples tests en paralelo?**
R: Sí, pero usa diferentes puertos de WebDriver.

**P: ¿Cómo cambio selectores sin tocar el código?**
R: Crea variables de configuración o lee de un archivo JSON.

**P: ¿Las screenshots ocupan mucho espacio?**
R: Sí. Comprime con `manage-screenshots.sh zip` o borra con `clean-old`.
