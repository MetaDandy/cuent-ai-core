package supabase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

var (
	baseURL string = os.Getenv("SUPABASE_PROJECT_URL")
	apiKey  string = os.Getenv("SUPABASE_API_KEY_SERVICE_ROLE")
)

func (s *Service) Upload(
	ctx context.Context,
	bucket string,
	dirPath string, // carpeta(s) dentro del bucket (puede ir vacío)
	fileName string, // nombre del objeto (obligatorio)
	body io.Reader,
	mime string,
	upsert bool,
) (string, error) {

	// 1. Normalizar la ruta: eliminar / inicial o final para evitar "//"
	cleanDir := strings.Trim(dirPath, "/")
	objectPath := path.Join(cleanDir, fileName) // usa el paquete path para OS-agnostic

	// 2. Construir URL Supabase
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", baseURL, bucket, objectPath)

	// 3. Crear petición HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", mime)
	if upsert {
		req.Header.Set("x-upsert", "true")
	}

	// 4. Ejecutar
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("supabase: %s – %s", res.Status, string(b))
	}
	return url, nil
}
