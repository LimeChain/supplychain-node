package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	_ "github.com/joho/godotenv/autoload"
)

func createHCSTopic(a, b hedera.PublicKey) hedera.TopicID {
	shouldConnectToMainnet := (os.Getenv("HCS_MAINNET") == "true")

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("HCS_OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("HCS_OPERATOR_PRV_KEY"))
	if err != nil {
		panic(err)
	}

	var client *hedera.Client

	if shouldConnectToMainnet {
		client = hedera.ClientForMainnet()
	} else {
		client = hedera.ClientForTestnet()
	}

	client = client.SetOperator(operatorAccountID, operatorPrivateKey)

	transactionID, err := hedera.NewTopicCreateTransaction().
		SetAdminKey(operatorPrivateKey.PublicKey()).
		SetAutoRenewAccountID(operatorAccountID).
		SetSubmitKey(a).
		SetSubmitKey(b).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	topicID := *transactionReceipt.TopicID

	return topicID

}

func main() {

	a := os.Getenv("A_PUB_KEY")
	b := os.Getenv("B_PUB_KEY")

	aPubKey, err := hedera.PublicKeyFromString(a)
	if err != nil {
		panic(err)
	}

	bPubKey, err := hedera.PublicKeyFromString(b)
	if err != nil {
		panic(err)
	}

	topicID := createHCSTopic(aPubKey, bPubKey)

	fmt.Printf("topicID: %v\n", topicID)

}
