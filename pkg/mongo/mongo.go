package mongo

import {
	config "wallet/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
}

func connect () (*mongo.Client, error) {
	cfg := config.LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return mongo.Client{}, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting to MongoDB: %v", err)
		}
	}()

	return client, nil
}