package tts

type RequestElevenTTS struct {
	Text    string `json:"text" validate:"required"`
	VoiceID string `json:"voice_id"`
}

type RequestElevenSFX struct {
	Description     string  `json:"description"`
	DurationSeconds float64 `json:"duration_seconds"`
	PromptInfluence float64 `json:"prompt_influence"`
	OutputFormat    string  `json:"output_format"`
}
