package dtos

type RequestElevenTTS struct {
	Text    string `json:"text" validate:"required"`
	VoiceID string `json:"voice_id"`
}
