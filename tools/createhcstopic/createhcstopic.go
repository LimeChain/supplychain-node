package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	_ "github.com/joho/godotenv/autoload"
)

func createHCSTopic() hedera.TopicID {
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

	topicID := createHCSTopic()

	fmt.Printf("topicID: %v\n", topicID)

}
