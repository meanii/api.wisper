package configs

import (
	"log"

	"github.com/spf13/viper"
)

var Env *envs

type envs struct {
	Port         string `mapstructure:"PORT"`
	MongoUrl     string `mapstructure:"MONGO_URI"`
	Database     string `mapstructure:"DATABASE"`
	RedisUrl     string `mapstructure:"REDIS_URI"`
	SecretToken  string `mapstructure:"SECRET_TOKEN"`
	RefreshToken string `mapstructure:"REFRESH_TOKEN"`
}

func InitConfig() {
	Env = loadConfig()
}

func loadConfig() (e *envs) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Error reading config file, %s", err)
	}
	if err = viper.Unmarshal(&e); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return
}
