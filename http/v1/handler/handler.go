package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/muhfaris/swarm-autoscale/core/controller"
	"github.com/muhfaris/swarm-autoscale/core/models"
)

func HandlerServiceStats() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var (
			ctx  = c.UserContext()
			body models.ReqContainer
		)

		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(map[string]string{"error": err.Error()})
		}

		err := controller.UpdateService(ctx, body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "ok"})
	}
}
