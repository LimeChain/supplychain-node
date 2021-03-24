package hcs

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
	"github.com/hashgraph/hedera-sdk-go/v2"
	log "github.com/sirupsen/logrus"
)

const SequenceNumberKey = "SequenceNumber"
const TransactionIdKey = "TransactionId"
const DLTValuesKey = "DLTValues"

type DLTValues struct {
	m map[string]string
}

func (v DLTValues) Get(key string) string {
	return v.m[key]
}

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
				log.Infof("[HCS] The topic response received %s\n", resp)

				txId := prepareTxId(fmt.Sprintf("%v", resp.TransactionID))
				sequenceNumber := strconv.FormatUint(resp.SequenceNumber, 10)

				log.Infof("[HCS] The topic response received - TransactionID: %s\n", txId)
				log.Infof("[HCS] The topic response received - SequenceNumber ID: %s\n", sequenceNumber)

				dltContextValues := DLTValues{map[string]string{
					SequenceNumberKey: sequenceNumber,
					TransactionIdKey:  txId,
				}}

				ctx := context.WithValue(context.Background(), DLTValuesKey, dltContextValues)
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

func prepareTxId(rawTxId string) string {
	// Make a Regex to say we only want numbers
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedTxId := reg.ReplaceAllString(rawTxId, "")

	return processedTxId
}
