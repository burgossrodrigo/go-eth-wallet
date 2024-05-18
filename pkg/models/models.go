type user struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email	string `json:"email" bson:"email"`
	WalletAddress string `json:"wallet_address" bson:"wallet_address"`
	Balances []Balance `json:"balances" bson:"balances"`
}

type Balance struct {
	Name string `json:"name" bson:"name"`
	Symbol string `json:"currency" bson:"currency"`
	Address string `json:"address" bson:"address"`
	Amount   float64 `json:"amount" bson:"amount"`
}

