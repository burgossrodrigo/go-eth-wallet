package user

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"log"
	"time"

	config "wallet/pkg/config"
)

func Connect() (*mongo.Client, error) {
	// Load configuration
	cfg := config.LoadEnv()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var client *mongo.Client
	var err error

	// Retry logic with backoff
	for i := 0; i < 5; i++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI).SetConnectTimeout(10*time.Second))
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				log.Println("Successfully connected to MongoDB")
				return client, nil
			}
		}
		log.Printf("Failed to connect to MongoDB (attempt %d): %v", i+1, err)
		time.Sleep(time.Duration(2*i) * time.Second)
	}

	// If still failing after retries, return error
	log.Fatal("Failed to connect to MongoDB after 5 attempts")
	return nil, err
}

func DisconnectClient(client *mongo.Client) error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Disconnect the client
	err := client.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect from MongoDB: %v", err)
		return err
	}

	log.Println("Successfully disconnected from MongoDB")
	return nil
}
