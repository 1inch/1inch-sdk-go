package web3_provider

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

func (w Wallet) Nonce(ctx context.Context) (uint64, error) {
	nonce, err := w.ethClient.NonceAt(ctx, w.address, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %v", err)
	}
	return nonce, nil
}

func (w Wallet) Address() common.Address {
	privateKey, err := crypto.HexToECDSA(w.privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address
}
