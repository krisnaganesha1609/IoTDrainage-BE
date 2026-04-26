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

	mqttcnf := utils.LoadMQTTConfig(config.MQTT_BROKER, config.MQTT_TOPIC)
	mqttConfig = mqttcnf
	mqttcl, err := mqttConfig.InitMQTT()
	if err != nil {
		log.Fatalf("%s", "Failed to initialize MQTT: "+err.Error())
	}
	mqttClient = mqttcl

	repo = repositories.InitializeRepository(influx, cloudinary)
	service = services.InitializeService(repo)
	handler = handlers.InitializeHandler(service)
	route = routes.InitializeRoutes(handler)
}

func main() {
	go route.Handler.ReceiveSensorFromMQTT(mqttClient, mqttConfig)

	app := fiber.New(fiber.Config{CaseSensitive: true,
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

	httpFileContentString := string(swaggerBytes)

	app.Get("/http-docs/*", scalar.New(scalar.Config{
		BasePath:          "/",
		FileContentString: httpFileContentString,
		Path:              "/http-docs",
		Title:             "IoT Drainage API Docs",
		Theme:             scalar.ThemeKepler,
	}))

	app.Get("/mqtt-docs/index.html", static.New("./docs/mqtt-docs/index.html"))

	log.Fatal(app.Listen(":" + config.PORT))
}
