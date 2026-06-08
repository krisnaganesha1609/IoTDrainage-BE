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
	return &Routes{Handler: handler}
}

func (r *Routes) Setup(app *fiber.App) {
	// WebSocket upgrade middleware
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

	// ── Web UI ────────────────────────────────────────────────────────────────
	// GET /logs — server-rendered page showing 7-day device system logs.
	// Optional query param: ?device_id=<id>
	app.Get("/logs", func(c fiber.Ctx) error {
		return r.Handler.ShowLogsPage(c)
	})

	// ── REST API ──────────────────────────────────────────────────────────────
	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return utils.RespondWithOK(c, "API is healthy", nil)
	})

	// Sensor telemetry (read-only for mobile/dashboard consumers)
	api.Get("/sensor/history/:device_id", func(c fiber.Ctx) error {
		return r.Handler.GetSensorHistory(c)
	})
	api.Get("/sensor/latest/:device_id", func(c fiber.Ctx) error {
		return r.Handler.GetLatestSensorData(c)
	})

	// Image endpoints
	// POST /api/image              — flood event photo (BAHAYA trigger)
	// POST /api/devices/:id/snapshot — scheduled daily snapshot at 07:00
	api.Post("/image", func(c fiber.Ctx) error {
		return r.Handler.UploadImage(c)
	})
	api.Get("/image", func(c fiber.Ctx) error {
		return r.Handler.GetLatestImage(c)
	})
	api.Post("/devices/:device_id/snapshot", func(c fiber.Ctx) error {
		return r.Handler.UploadSnapshot(c)
	})

	// FCM token registration (called by mobile app on first launch / token refresh)
	api.Post("/register-token", func(c fiber.Ctx) error {
		return r.Handler.RegisterToken(c)
	})
}
