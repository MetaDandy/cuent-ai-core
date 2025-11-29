//go:build selenium

package selenium

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tebeka/selenium"
)

// EdgeDriver maneja la configuración de WebDriver para Edge
type EdgeDriver struct {
	driver         selenium.WebDriver
	url            string
	screenshotDir  string
	screenshotStep int
}

// getEdgeBinary encuentra la ruta del ejecutable de Microsoft Edge
func getEdgeBinary() string {
	// Ruta estándar en Arch Linux con microsoft-edge-stable-bin
	binary := "/usr/bin/microsoft-edge-stable"

	// Verificar si existe
	if _, err := os.Stat(binary); err == nil {
		log.Printf("✓ Found Edge binary at: %s", binary)
		return binary
	}

	// Rutas alternativas
	paths := []string{
		os.Getenv("HOME") + "/.local/share/microsoft-edge-bin/microsoft-edge",
		"/opt/microsoft/edge/microsoft-edge",
		"/usr/bin/microsoft-edge",
		"/snap/bin/microsoft-edge",
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		"C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
		"C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe",
	}

	// Buscar en rutas conocidas
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("✓ Found Edge binary at: %s", path)
			return path
		}
	}

	// Intentar encontrar con 'which'
	if cmd, err := exec.LookPath("microsoft-edge"); err == nil {
		log.Printf("✓ Found Edge binary at: %s", cmd)
		return cmd
	}

	if cmd, err := exec.LookPath("edge"); err == nil {
		log.Printf("✓ Found Edge binary at: %s", cmd)
		return cmd
	}

	log.Println("⚠ Warning: Could not find Microsoft Edge binary.")
	log.Println("   Using default: /usr/bin/microsoft-edge-stable")
	return "/usr/bin/microsoft-edge-stable"
}

// NewEdgeDriver crea una nueva instancia del driver Edge
// baseURL es la URL de tu aplicación (ej: http://localhost:3000)
func NewEdgeDriver(baseURL string) (*EdgeDriver, error) {
	caps := selenium.Capabilities{
		"browserName": "MicrosoftEdge",
		"ms:edgeOptions": map[string]interface{}{
			"binary": getEdgeBinary(),
			"args": []string{
				"--start-maximized",
				"--disable-notifications",
			},
		},
	}

	// WebDriver se comunica con msedgedriver en puerto 4444
	driver, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Edge WebDriver: %w", err)
	}

	// Timeout de 10 segundos para encontrar elementos
	driver.SetImplicitWaitTimeout(10 * time.Second)

	// Crear directorio de screenshots
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	screenshotDir := fmt.Sprintf("tests/selenium/screenshots/%s", timestamp)
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		driver.Quit()
		return nil, fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	return &EdgeDriver{
		driver:         driver,
		url:            baseURL,
		screenshotDir:  screenshotDir,
		screenshotStep: 0,
	}, nil
}

// NewEdgeDriverWithName crea driver con nombre personalizado para screenshots
func NewEdgeDriverWithName(baseURL, testName string) (*EdgeDriver, error) {
	caps := selenium.Capabilities{
		"browserName": "MicrosoftEdge",
		"ms:edgeOptions": map[string]interface{}{
			"binary": getEdgeBinary(),
			"args": []string{
				"--start-maximized",
				"--disable-notifications",
			},
		},
	}

	driver, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Edge WebDriver: %w", err)
	}

	driver.SetImplicitWaitTimeout(10 * time.Second)

	// Crear directorio de screenshots con nombre del test
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	screenshotDir := fmt.Sprintf("tests/selenium/screenshots/%s_%s", timestamp, testName)
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		driver.Quit()
		return nil, fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	return &EdgeDriver{
		driver:         driver,
		url:            baseURL,
		screenshotDir:  screenshotDir,
		screenshotStep: 0,
	}, nil
}

// Close cierra la sesión del navegador
func (ed *EdgeDriver) Close() error {
	if ed.driver == nil {
		return nil
	}
	return ed.driver.Quit()
}

// GetDriver retorna el driver interno
func (ed *EdgeDriver) GetDriver() selenium.WebDriver {
	return ed.driver
}

// BaseURL retorna la URL base
func (ed *EdgeDriver) BaseURL() string {
	return ed.url
}

// NavigateTo navega a una URL
func (ed *EdgeDriver) NavigateTo(path string) error {
	fullURL := ed.url + path
	if err := ed.driver.Get(fullURL); err != nil {
		return fmt.Errorf("failed to navigate to %s: %w", fullURL, err)
	}
	log.Printf("✓ Navigated to: %s", fullURL)
	return nil
}

// FindElement busca un elemento por selector CSS
func (ed *EdgeDriver) FindElement(selector string) (selenium.WebElement, error) {
	element, err := ed.driver.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		return nil, fmt.Errorf("element not found: %s - %w", selector, err)
	}
	return element, nil
}

// FindElements busca múltiples elementos por selector CSS
func (ed *EdgeDriver) FindElements(selector string) ([]selenium.WebElement, error) {
	elements, err := ed.driver.FindElements(selenium.ByCSSSelector, selector)
	if err != nil {
		return nil, fmt.Errorf("elements not found: %s - %w", selector, err)
	}
	return elements, nil
}

// Click hace click en un elemento
func (ed *EdgeDriver) Click(selector string) error {
	element, err := ed.FindElement(selector)
	if err != nil {
		return err
	}
	if err := element.Click(); err != nil {
		return fmt.Errorf("failed to click element %s: %w", selector, err)
	}
	log.Printf("✓ Clicked: %s", selector)
	return nil
}

// SendKeys envía texto a un input
func (ed *EdgeDriver) SendKeys(selector, text string) error {
	element, err := ed.FindElement(selector)
	if err != nil {
		return err
	}
	if err := element.SendKeys(text); err != nil {
		return fmt.Errorf("failed to send keys to %s: %w", selector, err)
	}
	log.Printf("✓ Sent keys to %s: %s", selector, text)
	return nil
}

// GetText obtiene el texto de un elemento
func (ed *EdgeDriver) GetText(selector string) (string, error) {
	element, err := ed.FindElement(selector)
	if err != nil {
		return "", err
	}
	text, err := element.Text()
	if err != nil {
		return "", fmt.Errorf("failed to get text from %s: %w", selector, err)
	}
	return text, nil
}

// GetAttribute obtiene un atributo de un elemento
func (ed *EdgeDriver) GetAttribute(selector, attrName string) (string, error) {
	element, err := ed.FindElement(selector)
	if err != nil {
		return "", err
	}
	value, err := element.GetAttribute(attrName)
	if err != nil {
		return "", fmt.Errorf("failed to get attribute %s from %s: %w", attrName, selector, err)
	}
	return value, nil
}

// WaitFor espera a que una condición se cumpla
func (ed *EdgeDriver) WaitFor(condition func() bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for condition")
}

// TakeScreenshot toma una captura de pantalla
func (ed *EdgeDriver) TakeScreenshot(filename string) error {
	screenshot, err := ed.driver.Screenshot()
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %w", err)
	}

	// Aquí solo retornamos el error, el caller decide dónde guardar
	if len(screenshot) == 0 {
		return fmt.Errorf("screenshot is empty")
	}
	log.Printf("✓ Screenshot captured: %d bytes", len(screenshot))
	return nil
}

// AutoScreenshot toma automáticamente una screenshot y la guarda
func (ed *EdgeDriver) AutoScreenshot(description string) error {
	screenshot, err := ed.driver.Screenshot()
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %w", err)
	}

	if len(screenshot) == 0 {
		return fmt.Errorf("screenshot is empty")
	}

	ed.screenshotStep++
	filename := fmt.Sprintf("%02d_%s.png", ed.screenshotStep, description)
	filepath := filepath.Join(ed.screenshotDir, filename)

	if err := ioutil.WriteFile(filepath, screenshot, 0644); err != nil {
		return fmt.Errorf("failed to write screenshot: %w", err)
	}

	log.Printf("✓ Screenshot saved: %s", filepath)
	return nil
}

// Screenshot alias corto para AutoScreenshot
func (ed *EdgeDriver) Screenshot(description string) error {
	return ed.AutoScreenshot(description)
}

// GetScreenshotDir retorna el directorio de screenshots
func (ed *EdgeDriver) GetScreenshotDir() string {
	return ed.screenshotDir
}
