package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	_ "github.com/joho/godotenv/autoload"
)

func generateEd25519Keypair() ed25519.PrivateKey {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	pubKeyString := hex.EncodeToString(pub)
	fmt.Printf("Ed25519 public key hex: %s\n", pubKeyString)

	return priv
}

func createHederaAccount(key ed25519.PrivateKey) hedera.AccountID {

	shouldConnectToMainnet := (os.Getenv("HCS_MAINNET") == "true")

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("HCS_OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("HCS_OPERATOR_PRV_KEY"))
	if err != nil {
		panic(err)
	}

	newKey, err := hedera.PrivateKeyFromBytes(key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hedera Account Public Key: %v\n", newKey.PublicKey().String())

	var client *hedera.Client

	if shouldConnectToMainnet {
		client = hedera.ClientForMainnet()
	} else {
		client = hedera.ClientForTestnet()
	}

	client = client.SetOperator(operatorAccountID, operatorPrivateKey)

	transactionID, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	newAccountID := *transactionReceipt.AccountID

	transactionID, err = hedera.NewTransferTransaction().
		AddHbarTransfer(operatorAccountID, hedera.NewHbar(-100)).
		AddHbarTransfer(newAccountID, hedera.NewHbar(100)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	_, err = transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	return newAccountID
}

func main() {

	key := generateEd25519Keypair()

	hcsAccId := createHederaAccount(key)

	fmt.Printf("Hedera Account Id %v\n", hcsAccId.String())

	x509Encoded, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}

	pemEncodedPrv := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	fmt.Println("Private key in PEM format")
	fmt.Println(string(pemEncodedPrv))

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/config/key.pem", path), pemEncodedPrv, 0644)
	if err != nil {
		panic(err)
	}

}
