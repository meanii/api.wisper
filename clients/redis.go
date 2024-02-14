package clients

import (
	"runtime"

	"github.com/gofiber/storage/redis/v2"

	"github.com/meanii/api.wisper/configs"
)

func NewRedisClient(host string) *redis.Storage {
	// setting up redis
	storage := redis.New(redis.Config{
		Host:      host,
		Port:      6379,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	return storage
}

var redisClient = NewRedisClient(configs.Env.RedisUrl)

type Redis struct {
	Storage *redis.Storage
}

var RedisClient = &Redis{
	Storage: redisClient,
}
