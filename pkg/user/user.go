package user

import (
	mongodb "wallet/pkg/mongo"
	models "wallet/pkg/models"
)

func createUser(user) {
	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	collection := connection.Database("wallet").Collection("users")
}



