package hcs

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
	"github.com/hashgraph/hedera-sdk-go/v2"
	log "github.com/sirupsen/logrus"
)

const SequenceNumberKey = "SequenceNumber"

type HCSClient struct {
	client  *hedera.Client
	topicID hedera.TopicID
}

func (c *HCSClient) Send(msg *common.Message) error {
	id, err := hedera.NewTopicMessageSubmitTransaction().
		SetTopicID(c.topicID).
		SetMessage(msg.Msg).
		Execute(c.client)

	if err != nil {
		return err
	}

	_, err = id.GetReceipt(c.client)

	if err != nil {
		return err
	}

	log.Infof("Sent message to HCS with Id :%s\n", id.TransactionID.String())

	return nil
}

func (c *HCSClient) Listen(receiver common.MessageReceiver) error {
	_, err := hedera.NewTopicMessageQuery().
		SetTopicID(c.topicID).
		Subscribe(
			c.client,
			func(resp hedera.TopicMessage) {
				fmt.Println("UDri resp.TransactionID")
				fmt.Println(resp.TransactionID)
				ctx := context.WithValue(context.Background(), SequenceNumberKey, resp.SequenceNumber)
				receiver.Receive(&common.Message{Msg: resp.Contents, Ctx: ctx})
			})

	if err != nil {
		return err
	}
	return nil
}

func (c *HCSClient) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

func NewHCSClient(account string, key ed25519.PrivateKey, topicID string, mainnet bool) *HCSClient {

	hcsPrvKey, err := hedera.PrivateKeyFromBytes(key)
	if err != nil {
		panic(err)
	}

	acc, err := hedera.AccountIDFromString(account)
	if err != nil {
		panic(err)
	}

	var client *hedera.Client

	if mainnet {
		client = hedera.ClientForMainnet()
		log.Infoln("[HCS] HCS Client connecting to mainnet")
	} else {
		client = hedera.ClientForTestnet()
		log.Infoln("[HCS] HCS Client connecting to testnet")
	}

	client = client.SetOperator(acc, hcsPrvKey)

	hcsTopicId, err := hedera.TopicIDFromString(topicID)
	if err != nil {
		panic(err)
	}

	log.Infof("[HCS] HCS Client started with account ID: %s\n", account)

	return &HCSClient{client: client, topicID: hcsTopicId}
}
