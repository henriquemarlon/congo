package configs

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoCollection() (*mongo.Collection, error) {
	uri, ok := os.LookupEnv("MONGO_DB_URL")
	if !ok {
		slog.Error("MONGO_DB_URL not set")
		os.Exit(1)
	}

	dbName, ok := os.LookupEnv("MONGO_DB_NAME")
	if !ok {
		slog.Error("MONGO_DB_NAME not set")
		os.Exit(1)
	}

	collectionName, ok := os.LookupEnv("MONGO_COLLECTION_NAME")
	if !ok {
		slog.Error("MONGO_COLLECTION_NAME not set")
		os.Exit(1)
	}

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return collection, nil
}
