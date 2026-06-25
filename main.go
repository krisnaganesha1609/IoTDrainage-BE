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

	firebase := utils.InitFirebase()

	repo = repositories.InitializeRepository(influx, cloudinary, firebase, config.INFLUX_BUCKET)
	service = services.InitializeService(repo, firebase)

	// Step 1: build handler first (needed as OnConnectHandler closure).
	handler = handlers.InitializeHandler(service)
	route = routes.InitializeRoutes(handler)

	// Step 2: init MQTT — passes handler.SubscribeAllMQTT as OnConnectHandler
	// so subscriptions are (re)established on every connect/reconnect.
	mqttcnf := utils.LoadMQTTConfig(config.MQTT_BROKER, config.MQTT_BASE_TOPIC)
	mqttConfig = mqttcnf
	mqttcl, err := mqttConfig.InitMQTT(handler.SubscribeAllMQTT(mqttConfig))
	if err != nil {
		log.Fatalf("%s", "Failed to initialize MQTT: "+err.Error())
	}
	mqttClient = mqttcl

	// Step 3: inject MQTTClient back into handler so message callbacks
	// can call MarkMessageReceived() for the connection watchdog.
	handler.SetMQTTClient(mqttClient)
}

func main() {
	// Watchdog: marks devices OFFLINE when they miss their wakeup window.
	go service.RunWatchdog()

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

	app.Use("/mqtt-docs/*", static.New("./docs/mqtt-docs"))

	log.Fatal(app.Listen(":" + config.PORT))
}
