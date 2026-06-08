package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/krisnaganesha1609/IoTDrainage-BE/configs"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/handlers"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/repositories"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/services"
	"github.com/krisnaganesha1609/IoTDrainage-BE/routes"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var (
	config     *configs.Config
	influx     *configs.InfluxDB
	cloudinary *configs.CloudinaryConfig
	mqttConfig *utils.MQTTConfig
	mqttClient *utils.MQTTClient
	repo       *repositories.Repository
	service    *services.Service
	handler    *handlers.Handler
	route      *routes.Routes
)

func init() {
	conf, err := configs.InitConfig()
	if err != nil {
		log.Fatalf("%s", "Failed to load configuration: "+err.Error())
	}
	config = conf

	conn, err := configs.InitInfluxDB(config.INFLUX_URL, config.INFLUX_TOKEN, config.INFLUX_ORG, config.INFLUX_BUCKET)
	if err != nil {
		log.Fatalf("%s", "Failed to initialize InfluxDB: "+err.Error())
	}
	influx = conn

	cld, err := configs.InitCloudinary(config.CLOUDINARY_URL)
	if err != nil {
		log.Fatalf("%s", "Failed to initialize Cloudinary: "+err.Error())
	}
	cloudinary = cld

	// MQTT — now driven by MQTT_BASE_TOPIC; three sub-topics are derived automatically.
	mqttcnf := utils.LoadMQTTConfig(config.MQTT_BROKER, config.MQTT_BASE_TOPIC)
	mqttConfig = mqttcnf
	mqttcl, err := mqttConfig.InitMQTT()
	if err != nil {
		log.Fatalf("%s", "Failed to initialize MQTT: "+err.Error())
	}
	mqttClient = mqttcl

	firebase := utils.InitFirebase()

	// Pass INFLUX_BUCKET so Flux queries don't hardcode the bucket name.
	repo = repositories.InitializeRepository(influx, cloudinary, firebase, config.INFLUX_BUCKET)
	service = services.InitializeService(repo, firebase)
	handler = handlers.InitializeHandler(service)
	route = routes.InitializeRoutes(handler)
}

func main() {
	// ── MQTT Subscribers ────────────────────────────────────────────────────
	// Each runs in its own goroutine; they are pure subscribers (no publish).
	go route.Handler.ReceiveSensorFromMQTT(mqttClient, mqttConfig)
	go route.Handler.ReceiveSensorLogFromMQTT(mqttClient, mqttConfig)
	go route.Handler.ReceiveDeviceInfoFromMQTT(mqttClient, mqttConfig)

	// ── Watchdog Cron ────────────────────────────────────────────────────────
	// Runs every minute; marks devices OFFLINE when they miss their wakeup window.
	go service.RunWatchdog()

	// ── HTTP Server ──────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		CaseSensitive:      true,
		StrictRouting:      true,
		EnableIPValidation: true,
		StructValidator: &utils.Validator{
			Validator: validator.New(),
		},
		ServerHeader: "Backend",
		AppName:      "🔥 IoT Drainage API",
	})

	route.Setup(app)

	// Scalar API docs (HTTP)
	swaggerBytes, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Fatalf("Failed to read Swagger file: %v", err)
	}
	app.Get("/http-docs/*", scalar.New(scalar.Config{
		BasePath:          "/",
		FileContentString: string(swaggerBytes),
		Path:              "/http-docs",
		Title:             "IoT Drainage API Docs",
		Theme:             scalar.ThemeKepler,
	}))

	// AsyncAPI docs (MQTT)
	app.Use("/mqtt-docs/*", static.New("./docs/mqtt-docs"))

	log.Fatal(app.Listen(":" + config.PORT))
}
