package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"wallet/pkg/auth"
	"wallet/pkg/models"
	user "wallet/pkg/user"
)

func main() {
	r := gin.Default()

	// Apply authentication middleware
	r.POST("/login", loginUser) // Handle user login

	eg := r.Group("/api/v1", AuthMiddleware()) // Protect these routes with AuthMiddleware
	{
		eg.POST("/create", createUser)
		eg.PATCH("/update/:id", updateUser)
		eg.DELETE("/users/:id", deleteUser)
		eg.GET("/users/:id", getUser) // Optionally, protect the "get user" endpoint too
	}

	r.Run(":8080")
}

// AuthMiddleware validates JWT token from the request
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.ValidateToken(tokenString) // Use auth.ValidateToken for validation
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is expired and handle refresh
		if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}

			// Set the user email in the context
			c.Set("email", claims.Email)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// loginUser authenticates a user and returns a JWT token
func loginUser(c *gin.Context) {
	var userModel models.User
	if err := c.ShouldBindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate the user and generate a JWT token
	token, err := auth.UserLogin(userModel.Email, userModel.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// createUser creates a new user (protected by AuthMiddleware)
func createUser(c *gin.Context) {
	var userModel models.User
	if err := c.ShouldBindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := user.CreateUser(userModel.Username, userModel.Email, userModel.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// getUser retrieves a user by ID (protected by AuthMiddleware)
func getUser(c *gin.Context) {
	// Get user by ID
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userData, err := user.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map userData to UserResponse
	userResponse := models.UserResponse{
		ID:       userData.ID,
		Username: userData.Username,
		Email:    userData.Email,
		Balances: userData.Balances,
	}

	c.JSON(http.StatusOK, userResponse)
}

// updateUser updates a user's data (protected by AuthMiddleware)
func updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedData map[string]interface{}
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.UpdateUser(userID, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// deleteUser deactivates a user (protected by AuthMiddleware)
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = user.DeactivateUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}
