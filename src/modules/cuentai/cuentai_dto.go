package cuentai

type CuentAIFlowResult struct {
	Lines    []string
	TTSClips [][]byte
	SFXClips [][]byte
	Combined []byte
}
