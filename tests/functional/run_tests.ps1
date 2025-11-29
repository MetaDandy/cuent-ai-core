# Script para ejecutar pruebas funcionales E2E

Write-Host "Ejecutando Pruebas Funcionales E2E" -ForegroundColor Cyan

# Base URL: utiliza variable de entorno TEST_API_URL si existe
$baseUrl = $env:TEST_API_URL
if ([string]::IsNullOrEmpty($baseUrl)) {
    $baseUrl = "http://localhost:8000/api/v1"
}

# Verificar servidor
try {
    $uri = [Uri]$baseUrl
    $hostRoot = "{0}://{1}:{2}" -f $uri.Scheme, $uri.Host, $uri.Port
    $response = Invoke-WebRequest -Uri $hostRoot -Method GET -TimeoutSec 5 -ErrorAction Stop
    Write-Host "Servidor API activo en $hostRoot" -ForegroundColor Green
} catch {
    Write-Host "Error: Servidor API no está corriendo o $hostRoot no responde" -ForegroundColor Red
    Write-Host "Ejecuta: docker-compose -f docker-compose.dev.yml up" -ForegroundColor Yellow
    exit 1
}

# Crear directorio de reportes
$reportsDir = "tests\functional\reports"
if (-not (Test-Path $reportsDir)) {
    New-Item -ItemType Directory -Path $reportsDir | Out-Null
}

# Ejecutar pruebas (compilar y ejecutar con tag e2e)
Write-Host "Ejecutando pruebas..." -ForegroundColor Cyan
go test -v -tags=e2e ./tests/functional/... -json | Tee-Object -FilePath "$reportsDir\test_output.json"

$exitCode = $LASTEXITCODE

if ($exitCode -eq 0) {
    Write-Host "PRUEBAS EXITOSAS" -ForegroundColor Green
} else {
    Write-Host "ALGUNAS PRUEBAS FALLARON" -ForegroundColor Red
}

exit $exitCode
