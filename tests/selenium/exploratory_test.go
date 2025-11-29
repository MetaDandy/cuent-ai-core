//go:build selenium

package selenium

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestExploratoryLoginFlow es una prueba exploratoria que toma screenshots en cada paso
func TestExploratoryLoginFlow(t *testing.T) {
	// Crear driver con nombre personalizado
	driver, err := NewEdgeDriverWithName("http://localhost:3000", "login")
	require.NoError(t, err, "Failed to create Edge driver")
	defer driver.Close()

	// 1. Navegar a login
	err = driver.NavigateTo("/login")
	require.NoError(t, err)
	driver.AutoScreenshot("01_login_page_loaded")

	time.Sleep(1 * time.Second)

	// 2. Completar email
	err = driver.SendKeys("input[type='email']", "alo@gmail.com")
	require.NoError(t, err)
	driver.AutoScreenshot("02_email_entered")

	time.Sleep(500 * time.Millisecond)

	// 3. Completar password
	err = driver.SendKeys("input[type='password']", "holamundo")
	require.NoError(t, err)
	driver.AutoScreenshot("03_password_entered")

	time.Sleep(500 * time.Millisecond)

	// 4. Click en submit
	err = driver.Click("button[type='submit']")
	require.NoError(t, err)
	driver.AutoScreenshot("04_submit_clicked")

	time.Sleep(3 * time.Second)

	// 5. Verificar dashboard
	driver.AutoScreenshot("05_after_login_redirect")

	t.Logf("✓ Exploratory login test completed. Screenshots in: %s", driver.GetScreenshotDir())
}

// TestExploratorySignup prueba el flujo de registro
func TestExploratorySignup(t *testing.T) {
	driver, err := NewEdgeDriverWithName("http://localhost:3000", "signup")
	require.NoError(t, err)
	defer driver.Close()

	// 1. Ir a signup
	err = driver.NavigateTo("/register")
	require.NoError(t, err)
	time.Sleep(2 * time.Second)
	driver.AutoScreenshot("01_signup_page_loaded")

	// 2. Llenar nombre
	err = driver.SendKeys("input[name='name']", "New User")
	if err != nil {
		t.Logf("⚠ Could not fill name: %v", err)
		driver.AutoScreenshot("02_error_name_field")
		return
	}
	driver.AutoScreenshot("02_name_filled")

	// 3. Llenar email
	err = driver.SendKeys("input[type='email']", "newuser@example.com")
	if err != nil {
		t.Logf("⚠ Could not fill email: %v", err)
		driver.AutoScreenshot("03_error_email_field")
		return
	}
	driver.AutoScreenshot("03_email_filled")

	// 4. Llenar password
	err = driver.SendKeys("input[type='password']", "SecurePassword123!")
	if err != nil {
		t.Logf("⚠ Could not fill password: %v", err)
		driver.AutoScreenshot("04_error_password_field")
		return
	}
	driver.AutoScreenshot("04_password_filled")

	// 5. Confirmar password
	err = driver.SendKeys("input[name='confirmPassword']", "SecurePassword123!")
	if err != nil {
		t.Logf("⚠ Could not fill confirm password: %v", err)
		driver.AutoScreenshot("05_error_confirm_password_field")
		return
	}
	driver.AutoScreenshot("05_confirm_password_filled")

	// 6. Submit
	err = driver.Click("button[type='submit']")
	if err != nil {
		t.Logf("⚠ Could not submit form: %v", err)
		driver.AutoScreenshot("06_error_submitting")
		return
	}
	time.Sleep(3 * time.Second)
	driver.AutoScreenshot("06_after_signup_submit")

	t.Logf("✓ Exploratory signup test completed. Screenshots in: %s", driver.GetScreenshotDir())
}

// TestExploratoryFullUserJourney es un test que sigue el flujo completo del usuario
func TestExploratoryFullUserJourney(t *testing.T) {
	driver, err := NewEdgeDriverWithName("http://localhost:3000", "full_journey")
	require.NoError(t, err)
	defer driver.Close()

	// 1. Home page
	err = driver.NavigateTo("/")
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	driver.AutoScreenshot("01_home_page")
	time.Sleep(3 * time.Second)

	// 2. Go to login
	err = driver.NavigateTo("/login")
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	driver.AutoScreenshot("02_login_page")
	time.Sleep(3 * time.Second)

	// 3. Try login (puede fallar pero captura el estado)
	_ = driver.SendKeys("input[type='email']", "alo@gmail.com")
	_ = driver.SendKeys("input[type='password']", "holamundo")
	_ = driver.Click("button[type='submit']")
	time.Sleep(2 * time.Second)
	driver.AutoScreenshot("03_after_login_attempt")
	time.Sleep(3 * time.Second)

	// 4. Go to dashboard
	err = driver.NavigateTo("/dashboard")
	if err == nil {
		time.Sleep(1 * time.Second)
		driver.AutoScreenshot("04_dashboard_loaded")
		time.Sleep(3 * time.Second)
	}

	// 5. Go to projects
	err = driver.NavigateTo("/projects")
	if err == nil {
		time.Sleep(1 * time.Second)
		driver.AutoScreenshot("05_projects_loaded")
		time.Sleep(3 * time.Second)
	}

	// 6. Go to profile
	err = driver.NavigateTo("/profile")
	if err == nil {
		time.Sleep(1 * time.Second)
		driver.AutoScreenshot("06_profile_loaded")
		time.Sleep(3 * time.Second)
	}

	// 7. Go to subscription
	err = driver.NavigateTo("/subscription")
	if err == nil {
		time.Sleep(1 * time.Second)
		driver.AutoScreenshot("07_subscription_loaded")
	}

	t.Logf("✓ Full user journey test completed. Screenshots in: %s", driver.GetScreenshotDir())
}
