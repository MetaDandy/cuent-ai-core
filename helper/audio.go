package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hajimehoshi/go-mp3"
)

// ? Ver si enviar el context como parametro
func AudioOutput(line, id, bucket, dirPath, audioType string) (
	url, historyID string,
	duration time.Duration,
	err error,
) {
	var (
		audio    []byte
		fileName string
	)
	duration = 3 * time.Second

	if audioType == "SFX" {
		prompt := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		audio, historyID, err = TextToSoundEffects(
			prompt,
			duration.Seconds(),
			1.0,
			"mp3_44100_128",
		)

		fileName = fmt.Sprintf("sfx_%v.mp3", id)
	} else {
		// Esto es TTS normal
		audio, historyID, err = TextToSpeechElevenlabs(line, "")
		if err != nil {
			return "", historyID, 0, err
		}
		fileName = fmt.Sprintf("tts_%v.mp3", id)
		duration, err = Mp3Duration(audio)
		if err != nil {
			return "", historyID, 0, err
		}
	}

	if err != nil {
		return "", historyID, 0, err
	}

	if url, err = Upload(
		context.TODO(),
		bucket,
		dirPath,
		fileName,
		bytes.NewReader(audio),
		"audio/mpeg",
		true,
	); err != nil {
		return "", historyID, 0, err
	}

	return url, historyID, duration, nil
}

func TextToSpeechElevenlabs(text, voice_id string) ([]byte, string, error) {
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		return nil, "", fmt.Errorf("API key no encontrada en variables de entorno")
	}

	client := resty.New()
	if voice_id == "" {
		voice_id = "VR6AewLTigWG4xSOukaG" // 29vD33N1CtxCmqQRPOHJ
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voice_id)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("xi-api-key", apiKey).
		SetHeader("User-Agent", "Cuent-ai/1.0 (Go; +https://github.com/MetaDandy/cuent-ai-core)").
		SetHeader("Accept", "audio/mpeg").
		SetBody(map[string]interface{}{
			"text":     text,
			"model_id": "eleven_multilingual_v2",
			"voice_settings": map[string]interface{}{
				"stability":        0.5,
				"similarity_boost": 0.75,
			},
		}).
		SetDoNotParseResponse(true).
		Post(url)
	fmt.Printf("ElevenLabs status=%d body=%q\n", resp.StatusCode(), resp.String())
	if err != nil {
		return nil, "", err
	}
	defer resp.RawBody().Close()
	if resp.StatusCode() != 200 {
		return nil, "", fmt.Errorf("error ElevenLabs: %s", resp.String())
	}

	historyID := resp.Header().Get("history_item_id")
	audio, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return nil, historyID, err
	}
	return audio, historyID, nil
}

// TextToSoundEffects convierte una descripción en un efecto de sonido.
// durationSeconds: duración en segundos (0.1–22.0), o 0 para que la API la estime.
// promptInfluence: [0.0–1.0], cuánto se ajusta al prompt (nil = valor por defecto).
// outputFormat: ej. "mp3_44100_128" (vacio = mp3_44100_128).
func TextToSoundEffects(
	description string,
	durationSeconds float64,
	promptInfluence float64,
	outputFormat string,
) ([]byte, string, error) {
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		return nil, "", fmt.Errorf("API key no encontrada en variables de entorno")
	}

	client := resty.New()
	url := "https://api.elevenlabs.io/v1/sound-generation" // :contentReference[oaicite:0]{index=0}

	// Preparamos el body
	body := map[string]interface{}{"text": description}
	if durationSeconds > 0 {
		body["duration_seconds"] = durationSeconds
	}
	if promptInfluence >= 0 {
		body["prompt_influence"] = promptInfluence
	}
	if outputFormat != "" {
		body["output_format"] = outputFormat
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("xi-api-key", apiKey).
		SetHeader("User-Agent", "Cuent-ai/1.0 (Go; +https://github.com/MetaDandy/cuent-ai-core)").
		SetHeader("Accept", "audio/mpeg").
		SetBody(body).
		SetDoNotParseResponse(true).
		Post(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.RawBody().Close()

	historyID := resp.Header().Get("history_item_id")
	if resp.StatusCode() != 200 {
		return nil, historyID, fmt.Errorf("error ElevenLabs SFX: %s", resp.String())
	}

	audio, err := io.ReadAll(resp.RawBody())
	return audio, historyID, err
}

// ! Revisar método
// CharactersUsed obtiene el número de caracteres efectivos usados para generar
// un audio en ElevenLabs, consultando el historial hasta tres veces si es necesario.
//
// historyID es el identificador del elemento de historial en ElevenLabs.
//
// Devuelve el número de caracteres usados (To – From) o un error si falla la
// petición o si el historial no está disponible tras los reintentos.
func CharactersUsed(historyID string) (int, error) {
	client := resty.New()
	apiKey := os.Getenv("ELEVEN_API_KEY")

	// Poll hasta 3 veces porque puede haber retardo de propagación.
	for i := 0; i < 3; i++ {
		resp, err := client.R().
			SetHeader("xi-api-key", apiKey).
			Get(fmt.Sprintf("https://api.elevenlabs.io/v1/history/%s", historyID))
		if err != nil {
			return 0, err
		}
		if resp.StatusCode() == 200 {
			var data struct {
				From int `json:"character_count_change_from"`
				To   int `json:"character_count_change_to"`
			}
			if err := json.Unmarshal(resp.Body(), &data); err != nil {
				return 0, err
			}
			return data.To - data.From, nil
		}
		time.Sleep(2 * time.Second) // pequeño back-off
	}
	return 0, fmt.Errorf("history item %s no disponible tras reintentos", historyID)
}

func Mp3Duration(b []byte) (time.Duration, error) {
	d, err := mp3.NewDecoder(bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	// d.Length() = bytes de PCM (16-bit stereo) → 4 bytes por sample
	samples := d.Length() / 4
	return time.Duration(samples) * time.Second / time.Duration(d.SampleRate()), nil
}
