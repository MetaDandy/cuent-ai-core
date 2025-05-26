package api

import (
	"github.com/MetaDandy/cuent-ai-core/src"
	ws "github.com/MetaDandy/cuent-ai-core/src/modules/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupApi(app *fiber.App, c *src.Container) {
	v1 := app.Group("/api/v1")
	hub := ws.NewHub(4)

	app.Get("/ws", websocket.New(ws.ServeWs(hub)))

	v1.Get("/ws/disconnect", func(c *fiber.Ctx) error {
		userID := c.Query("user")
		roomID := c.Query("room")
		hub.Unregister(roomID, userID)
		return c.SendStatus(fiber.StatusOK)
	})

	v1.Get("/ws/close-room", func(c *fiber.Ctx) error {
		roomID := c.Query("room")
		hub.CloseRoom(roomID)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Aloha")
	})

	handlers := []func(fiber.Router){
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
