package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT           string `env:"PORT"`
	INFLUX_URL     string `env:"INFLUX_URL"`
	INFLUX_TOKEN   string `env:"INFLUX_TOKEN"`
	INFLUX_ORG     string `env:"INFLUX_ORG"`
	INFLUX_BUCKET  string `env:"INFLUX_BUCKET"`
	CLOUDINARY_URL string `env:"CLOUDINARY_URL"`
	MQTT_BROKER    string `env:"MQTT_BROKER"`
	MQTT_TOPIC     string `env:"MQTT_TOPIC"`
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found, relying on environment variables")
	}
	config := &Config{
		PORT:           LoadEnv("PORT", "8080"),
		INFLUX_URL:     LoadEnv("INFLUX_URL", ""),
		INFLUX_TOKEN:   LoadEnv("INFLUX_TOKEN", ""),
		INFLUX_ORG:     LoadEnv("INFLUX_ORG", ""),
		INFLUX_BUCKET:  LoadEnv("INFLUX_BUCKET", ""),
		CLOUDINARY_URL: LoadEnv("CLOUDINARY_URL", ""),
		MQTT_BROKER:    LoadEnv("MQTT_BROKER", ""),
		MQTT_TOPIC:     LoadEnv("MQTT_TOPIC", ""),
	}
	return config, nil
}

func LoadEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
