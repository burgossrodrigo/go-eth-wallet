type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email	string `json:"email" bson:"email"`
	Wallet WalletKey `json:"wallet" bson:"wallet"`
	Balances []Balance `json:"balances" bson:"balances"`
}

type Balance struct {
	Name string `json:"name" bson:"name"`
	Symbol string `json:"currency" bson:"currency"`
	Address string `json:"address" bson:"address"`
	Amount   float64 `json:"amount" bson:"amount"`
}

type WalletKey struct {
	PrivateKey string `json:"private_key
	PublicKey string `json:"public_key"`
}

