package auth

import (
	"context"
	"errors"
	"log"
	"time"
	config "wallet/pkg/config"
	"wallet/pkg/models"
	mongodb "wallet/pkg/mongo"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserClaims defines the JWT claims for the user
type UserClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// UserLogin authenticates a user and returns a JWT token
func UserLogin(email string, password string) (string, error) {
	connection, err := mongodb.Connect()
	if err != nil {
		return "", err
	}
	defer connection.Disconnect(context.Background())

	var userData models.User
	collection := connection.Database("wallet").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT token with a 15-minute expiration time
	expirationTime := time.Now().Add(15 * time.Minute)
	token, err := GenerateJWT(userData.Email, expirationTime)
	if err != nil {
		return "", err
	}

	return token, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateJWT generates a JWT token for a given email
func GenerateJWT(email string, expirationTime time.Time) (string, error) {
	claims := &UserClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Set the expiration time here
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key from the config
	tokenString, err := token.SignedString(config.LoadEnv().JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// storeJWTInDB stores the generated JWT in the database
func StoreJWTInDB(userID primitive.ObjectID, token string, expirationTime time.Time) error {
	connection, err := mongodb.Connect()
	if err != nil {
		return err
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("tokens")

	tokenData := models.Token{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expirationTime,
		IsActive:  true,
	}

	_, err = collection.InsertOne(context.Background(), tokenData)
	if err != nil {
		return err
	}

	return nil
}

// generateAndStoreToken generates a new token and stores it in the database
func GenerateAndStoreToken(userID primitive.ObjectID) (string, error) {
	// Generate a new JWT
	expirationTime := time.Now().Add(15 * time.Minute)
	token, err := GenerateJWT(userID.Hex(), expirationTime)
	if err != nil {
		return "", err
	}

	// Store the token in the database
	err = StoreJWTInDB(userID, token, expirationTime)
	if err != nil {
		return "", err
	}

	return token, nil
}

// validateAndRefreshToken validates the JWT and refreshes it if near expiration
func ValidateAndRefreshToken(userID primitive.ObjectID, token string) (string, error) {
	connection, err := mongodb.Connect()
	if err != nil {
		return "", err
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("tokens")

	var tokenData models.Token
	err = collection.FindOne(context.Background(), bson.M{"token": token, "user_id": userID, "is_active": true}).Decode(&tokenData)
	if err != nil {
		return "", errors.New("invalid token")
	}

	// If the token is near expiration, refresh it
	if time.Until(tokenData.ExpiresAt) < 5*time.Minute {
		newToken, err := GenerateAndStoreToken(userID)
		if err != nil {
			return "", err
		}

		// Mark old token as inactive
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": tokenData.ID}, bson.M{"$set": bson.M{"is_active": false}})
		if err != nil {
			return "", err
		}

		return newToken, nil
	}

	return token, nil
}

// validateToken validates the token by checking the database and its expiration
func ValidateToken(token string) (primitive.ObjectID, error) {
	connection, err := mongodb.Connect()
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("tokens")

	var tokenData models.Token
	err = collection.FindOne(context.Background(), bson.M{"token": token, "is_active": true}).Decode(&tokenData)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid token")
	}

	// Check if the token is expired
	if tokenData.ExpiresAt.Before(time.Now()) {
		return primitive.NilObjectID, errors.New("token expired")
	}

	return tokenData.UserID, nil
}

// expireOldTokens invalidates expired tokens in the database
func ExpireOldTokens() {
	connection, err := mongodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return
	}
	defer connection.Disconnect(context.Background())

	collection := connection.Database("wallet").Collection("tokens")

	// Find and invalidate expired tokens
	_, err = collection.UpdateMany(
		context.Background(),
		bson.M{"expires_at": bson.M{"$lt": time.Now()}},
		bson.M{"$set": bson.M{"is_active": false}},
	)
	if err != nil {
		log.Fatalf("Error expiring tokens: %v", err)
		return
	}

	log.Println("Expired tokens invalidated")
}
