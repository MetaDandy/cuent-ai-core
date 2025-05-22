package helper

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"
)

// GenerateVideo genera un vídeo a partir de un prompt usando Veo 2 y devuelve los bytes del MP4.
// La duración es entre 5 a 8 segundos
func GenerateVideo(prompt string, duration int32) ([]byte, error) {
	ctx := context.TODO()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:      os.Getenv("GEMINI_API_KEY"),
		Backend:     genai.BackendGeminiAPI,
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GenAI client: %w", err)
	}

	videoConfig := &genai.GenerateVideosConfig{
		AspectRatio:      "16:9",
		PersonGeneration: "allow_adult",
		DurationSeconds:  &duration,
	}
	op, err := client.Models.GenerateVideos(ctx, "models/veo-2.0-generate-001", prompt, nil, videoConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start video generation: %w", err)
	}

	// Poll until done
	for !op.Done {
		time.Sleep(5 * time.Second)
		op, err = client.Operations.GetVideosOperation(ctx, op, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operation status: %w", err)
		}
	}

	if len(op.Response.GeneratedVideos) == 0 {
		return nil, fmt.Errorf("no videos generated")
	}

	return op.Response.GeneratedVideos[0].Video.VideoBytes, nil
}
