package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Envs struct {
	Port     string
	MongoUrl string
	RedisUrl string
}

// LoadEnv loads the .env file
func (e *Envs) LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("loading env variables...")
	e.MongoUrl = os.Getenv("MONGO_URI")
	e.RedisUrl = os.Getenv("REDIS_URI")
	e.Port = os.Getenv("PORT")
}

var EnvInstance *Envs

func init() {
	// initialize static instance on load
	EnvInstance = &Envs{}
	EnvInstance.LoadEnv()
}

func GetConfig() *Envs {
	return EnvInstance
}
