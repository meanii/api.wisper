package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/meanii/api.wisper/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct{}

const (
	DATABASE = "wisper-local"
)

func (m *Mongo) Connect() *mongo.Client {
	mongodbUrl := configs.GetConfig().MongoUrl
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbUrl))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected to MongoDB %s.\n", mongodbUrl)
	return client
}

var MongoClient *Mongo

func init() {
	MongoClient.Connect()
}

func GetClient() *mongo.Client {
	return MongoClient.Connect()
}

func GetDatabase() *mongo.Database {
	return GetClient().Database(DATABASE)
}
