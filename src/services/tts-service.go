package services

import (
	"fmt"
	"io"
	"os"

	"github.com/go-resty/resty/v2"
)

type TTSService struct{}

func NewTTSService() *TTSService {
	return &TTSService{}
}

func (s *TTSService) TextToSpeechElevenlabs(text, voice_id string) ([]byte, error) {
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
