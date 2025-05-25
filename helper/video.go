package helper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GenerateVideo crea un MP4 que muestra cada imagen durante un tiempo
// calculado según la duración total y el número de imágenes.
// Si duration < len(images), mostrará sólo las primeras duration imágenes
// a 1 s cada una. Además aplica crossfade de 1 s entre cada par.
// Luego superpone el audio (audioURL) y recorta al menor de vídeo o audio.
func GenerateVideo(images []Image, audioURL string, duration float64) ([]byte, error) {
	nImgs := len(images)
	if nImgs == 0 {
		return nil, fmt.Errorf("necesitas al menos una imagen")
	}

	// 1. Determinar cuántas imágenes usar y el tiempo por imagen
	useImgs := nImgs
	perImg := float64(duration) / float64(nImgs)
	if int(duration) < nImgs {
		useImgs = int(duration)
		perImg = 1.0
	}

	// 2. Crear carpeta temporal
	tmpDir, err := os.MkdirTemp("", "video-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// 3. Descargar las primeras useImgs imágenes
	imgPaths := make([]string, useImgs)
	client := http.Client{Timeout: 30 * time.Second}
	for i := range useImgs {
		resp, err := client.Get(images[i].Url)
		if err != nil {
			return nil, fmt.Errorf("descargando imagen %d: %w", i, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("imagen %d: código HTTP %d", i, resp.StatusCode)
		}

		path := filepath.Join(tmpDir, fmt.Sprintf("img-%02d.jpg", i))
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		io.Copy(f, resp.Body)
		f.Close()
		imgPaths[i] = path
	}

	// 4. Generar vídeo sin audio con crossfade
	videoNoAudio := filepath.Join(tmpDir, "video_no_audio.mp4")
	args := []string{"-y"}

	// Inputs: loop cada imagen el tiempo perImg+1 para cubrir el fade
	for _, p := range imgPaths {
		args = append(args,
			"-loop", "1",
			"-t", fmt.Sprintf("%.2f", perImg+1),
			"-i", p,
		)
	}

	// Construir filter_complex dinámico de xfade
	// crossFadeDur = 1 segundo
	crossFadeDur := 1.0
	var fc bytes.Buffer
	// Para el primer par: [0][1]xfade...
	for i := range useImgs - 1 {
		offset := perImg*float64(i+1) - crossFadeDur
		if i == 0 {
			fmt.Fprintf(&fc, "[0:v][1:v]xfade=transition=fade:duration=%.2f:offset=%.2f[v0];", crossFadeDur, offset)
		} else {
			fmt.Fprintf(&fc, "[v%d][%d:v]xfade=transition=fade:duration=%.2f:offset=%.2f[v%d];",
				i-1, i+1, crossFadeDur, offset, i)
		}
	}
	// El mapa final será [vN] donde N = useImgs-2
	finalLabel := fmt.Sprintf("[v%d]", useImgs-2)
	args = append(args,
		"-filter_complex", fc.String(),
		"-map", finalLabel,
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		videoNoAudio,
	)

	if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg vídeo sin audio: %s / %w", out, err)
	}

	// 5. Descargar audio
	audioPath := filepath.Join(tmpDir, "audio.mp3")
	resp, err := client.Get(audioURL)
	if err != nil {
		return nil, fmt.Errorf("descargando audio: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("audio: código HTTP %d", resp.StatusCode)
	}
	af, err := os.Create(audioPath)
	if err != nil {
		return nil, err
	}
	io.Copy(af, resp.Body)
	af.Close()

	// 6. Superponer audio y recortar al menor
	finalVid := filepath.Join(tmpDir, "final.mp4")
	mixCmd := exec.Command("ffmpeg",
		"-y",
		"-i", videoNoAudio,
		"-i", audioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-b:a", "192k",
		"-map", "0:v",
		"-map", "1:a",
		"-t", fmt.Sprintf("%.2f", duration),
		finalVid,
	)
	if out, err := mixCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg mezcla audio: %s / %w", out, err)
	}

	// 7. Leer y devolver bytes
	return os.ReadFile(finalVid)
}
