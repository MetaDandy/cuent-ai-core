package api

import (
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/gofiber/fiber/v2"
)

func SetupApi(app *fiber.App, c *src.Container) {
	v1 := app.Group("/api/v1")

	handlers := []func(fiber.Router){
		c.TTSHandler.RegisterTTSRoutes,
	}

	for _, register := range handlers {
		register(v1)
	}
}
