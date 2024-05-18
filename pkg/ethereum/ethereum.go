package ethereum

import (
	"crypto/ecdsa"
	"log"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"

	models "wallet/pkg/models"
)

func generateAddress () (string, error){
	privateKey, err	:= crypto.GenerateKey()
	if err != nil {
		Println("Error generating private key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
		return "", err
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
}