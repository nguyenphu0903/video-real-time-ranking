package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		MongoDB
		Redis
	}
	MongoDB struct {
		DatabaseName string `env:"MONGO_DATABASE_NAME"`
		URLString    string `env:"MONGO_URL_STRING"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD"`
	}
)

var C Config

func LoadConfig() {
	err := godotenv.Load("/Users/nguyenphu/Documents/nguyenphu/video-real-time-ranking/.env")
	if err != nil {
		log.Fatalln("Application config parsing failed: " + err.Error() + " => Exit!")
		return
	}

	cfg := &C
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalln("Application config parsing failed: " + err.Error() + " => Exit!")
		return
	}

	log.Println("Load Config Successfully!")
}
