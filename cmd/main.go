package main

import (
	"github.com/MetaDandy/cuent-ai-core/cmd/api"
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()

	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	c := src.SetupContainer()
	api.SetupApi(app, c)

	app.Listen(":" + config.Port)
}
