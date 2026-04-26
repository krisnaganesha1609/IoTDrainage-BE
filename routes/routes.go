package routes

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/handlers"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

type Routes struct {
	Handler *handlers.Handler
}

func InitializeRoutes(handler *handlers.Handler) *Routes {
	return &Routes{
		Handler: handler,
	}
}

func (r *Routes) Setup(app *fiber.App) {
	app.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		r.Handler.ConnectWebsocket(c)
	}))

	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return utils.RespondWithOK(c, "API is healthy", nil)
	})

	api.Post("/upload-image", func(c fiber.Ctx) error {
		return r.Handler.UploadImage(c)
	})

	api.Get("/sensor/history/:device_id", func(c fiber.Ctx) error {
		return r.Handler.GetSensorHistory(c)
	})

	api.Get("/sensor/latest/:device_id", func(c fiber.Ctx) error {
		return r.Handler.GetLatestSensorData(c)
	})
}
