package main

import (
	"github.com/MetaDandy/cuent-ai-core/cmd/api"
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/middleware"
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()

	app := fiber.New()
	app.Use(middleware.Logger())

	c := src.SetupContainer()
	api.SetupApi(app, c)

	app.Listen(":" + config.Port)
}
