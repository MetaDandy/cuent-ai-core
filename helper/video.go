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
		return nil, fmt.Errorf("necesitas al menos una imagen, prueba con otras keywords")
	}

	// 2. Crear carpeta temporal
	tmpDir, err := os.MkdirTemp("", "video-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// 5. Descargar audio
	client := http.Client{Timeout: 30 * time.Second}
	audioPath := filepath.Join(tmpDir, "audio.mp3")
	if err := downloadFile(client, audioURL, audioPath); err != nil {
		return nil, err
	}

	// 1. Determinar cuántas imágenes usar y tiempo por imagen (2 s fijo)
	perImg := 2.0
	maxImgs := int(duration / perImg)
	useImgs := max(min(maxImgs, nImgs), 1)

	if useImgs < 2 {
		imgPath := filepath.Join(tmpDir, "img.jpg")
		if err := downloadFile(client, images[0].Url, imgPath); err != nil {
			return nil, err
		}
		out := filepath.Join(tmpDir, "final.mp4")
		// loop 1 imagen, t = duration, luego mezclar audio sin recortar
		cmd := exec.Command("ffmpeg", "-y",
			"-loop", "1", "-i", imgPath,
			"-i", audioPath,
			"-c:v", "libx264", "-t", fmt.Sprintf("%.2f", duration),
			"-c:a", "aac", "-b:a", "192k",
			"-map", "0:v", "-map", "1:a",
			out,
		)
		if outp, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("ffmpeg estático: %s / %w", outp, err)
		}
		return os.ReadFile(out)
	}

	// 3. Descargar las primeras useImgs imágenes
	imgPaths := make([]string, useImgs)
	for i := 0; i < useImgs; i++ {
		imgPaths[i] = filepath.Join(tmpDir, fmt.Sprintf("img-%02d.jpg", i))
		if err := downloadFile(client, images[i].Url, imgPaths[i]); err != nil {
			return nil, err
		}
	}

	// Construir filter_complex: primero escalado+pad, luego xfade en cadena
	// Definimos el tamaño de salida deseado:
	width, height := 1080, 608
	// Construir filter_complex dinámico de xfade
	// crossFadeDur = 1 segundo
	crossFadeDur := 1.0
	var fc bytes.Buffer
	// Escalado de cada stream de entrada
	// ! NO usar range en los for
	for i := 0; i < useImgs; i++ {
		// [i:v]scale=1080:608:force_original_aspect_ratio=decrease,pad=1080:608:(ow-iw)/2:(oh-ih)/2[si{i}];
		fmt.Fprintf(&fc, "[%d:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2[si%d];",
			i, width, height, width, height, i)
	}

	// Para el primer par: [0][1]xfade...
	for i := 0; i < useImgs-1; i++ {
		offset := perImg*float64(i+1) - crossFadeDur
		if i == 0 {
			fmt.Fprintf(&fc, "[si0][si1]xfade=transition=fade:duration=%.2f:offset=%.2f[v0];",
				crossFadeDur, offset)
		} else {
			fmt.Fprintf(&fc, "[v%d][si%d]xfade=transition=fade:duration=%.2f:offset=%.2f[v%d];",
				i-1, i+1, crossFadeDur, offset, i)
		}
	}

	lastLabel := fmt.Sprintf("[v%d]", useImgs-2)
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

	args = append(args,
		"-filter_complex", fc.String(),
		"-map", lastLabel,
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		videoNoAudio,
	)

	if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg vídeo sin audio: %s / %w", out, err)
	}

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
func downloadFile(client http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("código HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
