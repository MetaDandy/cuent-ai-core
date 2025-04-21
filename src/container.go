package src

import (
	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
	"github.com/MetaDandy/cuent-ai-core/src/modules/cuentai"
)

type Container struct {
	// TTS
	TtsSvc     *tts.Service
	TtsHandler *tts.Handler

	// CuentAI
	CuentSvc     *cuentai.Service
	CuentHandler *cuentai.Handler
}

func SetupContainer() *Container {
	// TTS
	ttsSvc := tts.NewService()
	ttsHandler := tts.NewHandler(ttsSvc)

	// CuentAI
	cuentSvc := cuentai.NewService(ttsSvc)
	cuentHandler := cuentai.NewHandler(cuentSvc)

	return &Container{
		// TTS
		TtsSvc:     ttsSvc,
		TtsHandler: ttsHandler,

		// CuentAI
		CuentSvc:     cuentSvc,
		CuentHandler: cuentHandler,
	}
}
