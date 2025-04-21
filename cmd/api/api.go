package api

import (
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/gofiber/fiber/v2"
)

func SetupApi(app *fiber.App, c *src.Container) {
	v1 := app.Group("/api/v1")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	handlers := []func(fiber.Router){
		c.TtsHandler.RegisterTTSRoutes,
		c.CuentHandler.RegisterRoutes,
	}

	for _, register := range handlers {
		register(v1)
	}
}
