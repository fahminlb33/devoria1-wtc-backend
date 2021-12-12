package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"localhost"`
		Port int    `envconfig:"SERVER_PORT"`
	}
	Authentication struct {
		BasicUsername string `envconfig:"AUTH_BASIC_USERNAME"`
		BasicPassword string `envconfig:"AUTH_BASIC_PASSWORD"`
		PrivateKey    string `envconfig:"AUTH_JWT_PRIVATE_KEY"`
		PublicKey     string `envconfig:"AUTH_JWT_PUBLIC_KEY"`
	}
	Database struct {
		URI string `envconfig:"DB_URI"`
	}
}

var GlobalConfig Config = Config{}

func LoadConfig() {
	godotenv.Load()
	envconfig.Process("", &GlobalConfig)
}
