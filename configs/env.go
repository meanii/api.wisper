package configs

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

type Envs struct {
	Port         string
	MongoUrl     string
	RedisUrl     string
	SecretToken  string
	RefreshToken string
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
	e.SecretToken = os.Getenv("SECRET_TOKEN")
	e.RefreshToken = os.Getenv("REFRESH_TOKEN")
}

var EnvInstance *Envs

// GetConfig returns the singleton instance of Envs
func GetConfig() *Envs {
	if EnvInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if EnvInstance == nil {
			EnvInstance = &Envs{}
			EnvInstance.LoadEnv()
			fmt.Printf("loading env from sington: %v\n", EnvInstance)
		}
	}
	return EnvInstance
}
