package src

import (
	tts "github.com/MetaDandy/cuent-ai-core/src/modules/Tts"
	"github.com/MetaDandy/cuent-ai-core/src/modules/cuentai"
	"github.com/MetaDandy/cuent-ai-core/src/modules/supabase"
)

type Container struct {
	// TTS
	TtsSvc     *tts.Service
	TtsHandler *tts.Handler

	// CuentAI
	CuentSvc     *cuentai.Service
	CuentHandler *cuentai.Handler

	// Supabase
	SupaSvc     *supabase.Service
	SupaHandler *supabase.Handler
}

func SetupContainer() *Container {
	// TTS
	ttsSvc := tts.NewService()
	ttsHandler := tts.NewHandler(ttsSvc)

	// Supabase
	supaSvc := supabase.NewService()
	supaHandler := supabase.NewHandler(supaSvc)

	// CuentAI
	cuentSvc := cuentai.NewService(ttsSvc, supaSvc)
	cuentHandler := cuentai.NewHandler(cuentSvc)

	return &Container{
		// TTS
		TtsSvc:     ttsSvc,
		TtsHandler: ttsHandler,

		// CuentAI
		CuentSvc:     cuentSvc,
		CuentHandler: cuentHandler,

		// Supabase
		SupaSvc:     supaSvc,
		SupaHandler: supaHandler,
	}
}
