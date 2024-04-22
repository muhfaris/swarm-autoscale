package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhfaris/swarm-autoscale/http/v1/handler"
)

func Routers(app *fiber.App) {
	apiGroup := app.Group("/api")
	apiGroup.Post("/stats", handler.HandlerServiceStats())
}
