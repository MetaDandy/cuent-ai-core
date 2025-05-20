package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

func MixAudio(id string, assets []model.Asset) (string, error) {
	var (
		bucket  = "audio"
		dirPath = id
		apiKey  = os.Getenv("SUPABASE_API_KEY_SERVICE_ROLE")
		ctx     = context.TODO()
	)

	// 2. Preparar temporales
	var tmpFiles []string
	defer func() {
		for _, f := range tmpFiles {
			_ = os.Remove(f)
		}
	}()

	// 3. Descargar cada URL a un mp3 temporal
	for _, a := range assets {
		// crea fichero tmp
		tmp, err := os.CreateTemp("", fmt.Sprintf("asset_%d_*.mp3", a.Position))
		if err != nil {
			return "", err
		}
		tmpFiles = append(tmpFiles, tmp.Name())
		tmp.Close()

		// descarga
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, a.Audio_URL, nil)
		req.Header.Set("apikey", apiKey)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		out, err := os.OpenFile(tmp.Name(), os.O_WRONLY, 0)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, res.Body); err != nil {
			out.Close()
			return "", err
		}
		out.Close()
	}

	// 4. Crear archivo de lista para concat
	listFile, err := os.CreateTemp("", "ffmpeg_list_*.txt")
	if err != nil {
		return "", err
	}
	for _, p := range tmpFiles {
		fmt.Fprintf(listFile, "file '%s'\n", filepath.ToSlash(p))
	}
	listFile.Close()
	defer os.Remove(listFile.Name())

	// 5. Ejecutar ffmpeg concat
	mixPath := filepath.Join(os.TempDir(), fmt.Sprintf("mix_%d.mp3", time.Now().UnixNano()))
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "concat", "-safe", "0",
		"-i", listFile.Name(),
		"-acodec", "libmp3lame", "-b:a", "192k",
		"-ar", "44100", "-ac", "2",
		mixPath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %v – %s", err, string(out))
	}
	defer os.Remove(mixPath)

	// 6. Subir mix a Supabase
	mixBytes, err := os.ReadFile(mixPath)
	if err != nil {
		return "", err
	}
	mixName := fmt.Sprintf("mix_%s.mp3", id)
	mixedURL, err := Upload(ctx, bucket, dirPath, mixName,
		bytes.NewReader(mixBytes), "audio/mpeg", true)
	if err != nil {
		return "", err
	}

	return mixedURL, nil
}
