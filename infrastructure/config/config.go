package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"localhost"`
		Port int    `envconfig:"SERVER_PORT"`
	}
	Authentication struct {
		BasicUsername  string `envconfig:"AUTH_BASIC_USERNAME"`
		BasicPassword  string `envconfig:"AUTH_BASIC_PASSWORD"`
		PrivateKeyPath string `envconfig:"AUTH_JWT_PRIVATE_KEY_PATH"`
		PublicKeyPath  string `envconfig:"AUTH_JWT_PUBLIC_KEY_PATH"`
	}
	Database struct {
		URI string `envconfig:"DB_URI"`
	}
}

var GlobalConfig Config = Config{}

func LoadConfig() {
	godotenvErr := godotenv.Load()
	if godotenvErr != nil {
		log.Fatal("Error loading .env file", godotenvErr)
	}

	envconfigErr := envconfig.Process("", &GlobalConfig)
	if envconfigErr != nil {
		log.Fatal("Error decoding config from environment variables", envconfigErr)
	}
}
