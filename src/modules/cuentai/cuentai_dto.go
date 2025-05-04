package cuentai

type CuentAIFlowResult struct {
	Lines    []string `json:"lines"`
	TTSURLS  []string `json:"tts_url"`
	SFXURLS  []string `json:"sfx_url"`
	MixedURL string   `json:"mixed_url"`
}
