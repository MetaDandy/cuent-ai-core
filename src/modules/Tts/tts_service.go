package tts

import (
	"fmt"
	"io"
	"os"

	"github.com/go-resty/resty/v2"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) TextToSpeechElevenlabs(text, voice_id string) ([]byte, error) {
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key no encontrada en variables de entorno")
	}

	client := resty.New()
	if voice_id == "" {
		voice_id = "29vD33N1CtxCmqQRPOHJ"
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voice_id)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("xi-api-key", apiKey).
		SetBody(map[string]interface{}{
			"text":     text,
			"model_id": "eleven_monolingual_v1",
			"voice_settings": map[string]interface{}{
				"stability":        0.5,
				"similarity_boost": 0.75,
			},
		}).
		SetDoNotParseResponse(true).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error ElevenLabs: %s", resp.String())
	}

	defer resp.RawBody().Close()
	audio, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return nil, err
	}
	return audio, nil
}

// TextToSoundEffects convierte una descripción en un efecto de sonido.
// durationSeconds: duración en segundos (0.1–22.0), o 0 para que la API la estime.
// promptInfluence: [0.0–1.0], cuánto se ajusta al prompt (nil = valor por defecto).
// outputFormat: ej. "mp3_44100_128" (vacio = mp3_44100_128).
func (s *Service) TextToSoundEffects(
	description string,
	durationSeconds float64,
	promptInfluence float64,
	outputFormat string,
) ([]byte, error) {
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key no encontrada en variables de entorno")
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
		SetBody(body).
		SetDoNotParseResponse(true).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error ElevenLabs SFX: %s", resp.String())
	}

	defer resp.RawBody().Close()
	return io.ReadAll(resp.RawBody())
}
