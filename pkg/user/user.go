package user

import (
	"context"
	"errors"
	"log"
	"time"

	models "wallet/pkg/models"
	mongodb "wallet/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Assuming models package and User struct are defined elsewhere
// import "path/to/models"

func CreateUser(userName string, email string, password string) error {

	exists, err := doesUserExist(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exists")
	}

	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return err
	}

	defer mongodb.DisconnectClient(connection)

	userData := models.User{
		ID:        primitive.NewObjectID(),
		Username:  userName,
		Email:     email,
		Password:  password, // Consider hashing the password before storing
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
	}

	collection := connection.Database("wallet").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(userID primitive.ObjectID) (models.User, error) {
	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return models.User{}, err
	}

	defer mongodb.DisconnectClient(connection)

	var userData models.User
	collection := connection.Database("wallet").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&userData)
	if err != nil {
		return models.User{}, err
	}

	return userData, nil
}

func UpdateUser(userID primitive.ObjectID, updatedData map[string]interface{}) error {
	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return err
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("users")

	update := bson.M{}

	for key, value := range updatedData {
		update[key] = value
	}

	update["updated_at"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": update},
		options.Update().SetUpsert(false),
	)
	if err != nil {
		return err
	}

	return nil
}

func doesUserExist(email string) (bool, error) {
	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return false, err
	}

	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeactivateUser(userID primitive.ObjectID) error {
	connection, err := mongodb.Connect()
	if err != nil {
		return err
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the "active" field to false
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"active": false}},
	)
	if err != nil {
		return err
	}

	return nil
}
