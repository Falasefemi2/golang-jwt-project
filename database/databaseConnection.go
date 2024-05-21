package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func DBinstance() *mongo.Client {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	MongoDb := os.Getenv("MONGODB_URL")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return client
// }

func DBinstance() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MongoDb := os.Getenv("MONGODB_URL")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MongoDb).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ClusterAuth").Collection(collectionName)
	return collection
}
