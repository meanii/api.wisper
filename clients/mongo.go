package clients

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meanii/api.wisper/configs"
)

var MongoClient *mongo.Client

func connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configs.Env.MongoUrl))
	if err != nil {
		log.Panicf("Failed to connect to MongoDB %s.\n", configs.Env.MongoUrl)
	}

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panicf("Failed to ping MongoDB %s.\n", configs.Env.MongoUrl)
	}
	log.Fatalf("Connected to MongoDB %s.\n", configs.Env.MongoUrl)
	return client
}

func MongoInit() {
	MongoClient = connect()
}

func GetDatabase() *mongo.Database {
	return MongoClient.Database(configs.Env.Database)
}
