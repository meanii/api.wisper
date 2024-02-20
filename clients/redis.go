package clients

import (
	"runtime"

	"github.com/gofiber/storage/redis"

	"github.com/meanii/api.wisper/configs"
)

var RedisClient = *&redisClient{}

type redisClient struct {
	Storage *redis.Storage
}

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

func RedisInit() {
	RedisClient.Storage = NewRedisClient(configs.Env.RedisUrl)
}
