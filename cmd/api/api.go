package api

import (
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/gofiber/fiber/v2"
)

func SetupApi(app *fiber.App, c *src.Container) {
	v1 := app.Group("/api/v1")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Aloha")
	})

	handlers := []func(fiber.Router){
		c.TtsHandler.RegisterTTSRoutes,
		c.CuentHandler.RegisterRoutes,
		c.SupaHandler.RegisterRoutes,
		c.ProjectHdl.RegisterRoutes,
		c.UserHdl.RegisterRoutes,
		c.ScriptHdl.RegisterRoutes,
		c.AssetHdl.RegisterRoutes,
		c.SubsHdl.RegisterRoutes,
	}

	for _, register := range handlers {
		register(v1)
	}
}
