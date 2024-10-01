package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Wallet    WalletKey          `json:"wallet,omitempty" bson:"wallet,omitempty"`
	Balances  []Balance          `json:"balances,omitempty" bson:"balances,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Active    bool               `json:"active" bson:"active"`
}

type UserResponse struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Balances []Balance          `json:"balances,omitempty" bson:"balances,omitempty"`
}

type Balance struct {
	Name    string  `json:"name" bson:"name"`
	Symbol  string  `json:"currency" bson:"currency"`
	Address string  `json:"address" bson:"address"`
	Amount  float64 `json:"amount" bson:"amount"`
}

type WalletKey struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Mnemonic   string `json:"mnemonic"`
}

type Config struct {
	MongoURI  string `json:"mongo_uri"`
	JWTSecret string
}

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
	IsActive  bool               `bson:"is_active"`
}
