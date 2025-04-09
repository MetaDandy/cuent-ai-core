package src

import (
	"github.com/MetaDandy/cuent-ai-core/src/handlers"
	"github.com/MetaDandy/cuent-ai-core/src/services"
)

type Container struct {
	TTSService *services.TTSService
	TTSHandler *handlers.TTSHandler
}

func SetupContainer() *Container {
	ttsService := services.NewTTSService()
	ttsHandler := handlers.NewTTSHandler(ttsService)

	return &Container{
		TTSService: ttsService,
		TTSHandler: ttsHandler,
	}
}
