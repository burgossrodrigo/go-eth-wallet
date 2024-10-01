package ethereum

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"

	models "wallet/pkg/models"
)

func generateAddress() (models.WalletKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return models.WalletKey{}, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
		return models.WalletKey{}, err
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])

	// Generate mnemonic
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return models.WalletKey{}, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return models.WalletKey{}, err
	}

	return models.WalletKey{
		PublicKey:  address,
		PrivateKey: hexutil.Encode(privateKeyBytes)[2:],
		Mnemonic:   mnemonic,
	}, nil
}
