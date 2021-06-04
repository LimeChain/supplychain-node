package hcs

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

type (
	// Message struct used by the Hedera Mirror node REST API to represent Topic Message
	Message struct {
		ConsensusTimestamp string `json:"consensus_timestamp"`
		TopicId            string `json:"topic_id"`
		Contents           string `json:"message"`
		RunningHash        string `json:"running_hash"`
		SequenceNumber     int    `json:"sequence_number"`
	}
	// Messages struct used by the Hedera Mirror node REST API and returned once
	// Topic Messages are queried
	Messages struct {
		Messages []Message
	}
)

type (
	// Transaction struct used by the Hedera Mirror node REST API
	Transaction struct {
		ConsensusTimestamp   string     `json:"consensus_timestamp"`
		EntityId             string     `json:"entity_id"`
		TransactionHash      string     `json:"transaction_hash"`
		ValidStartTimestamp  string     `json:"valid_start_timestamp"`
		ChargedTxFee         int        `json:"charged_tx_fee"`
		MemoBase64           string     `json:"memo_base64"`
		Result               string     `json:"result"`
		Name                 string     `json:"name"`
		MaxFee               string     `json:"max_fee"`
		ValidDurationSeconds string     `json:"valid_duration_seconds"`
		Node                 string     `json:"node"`
		Scheduled            bool       `json:"scheduled"`
		TransactionID        string     `json:"transaction_id"`
		Transfers            []Transfer `json:"transfers"`
		TokenTransfers       []Transfer `json:"token_transfers"`
	}
	// Transfer struct used by the Hedera Mirror node REST API
	Transfer struct {
		Account string `json:"account"`
		Amount  int64  `json:"amount"`
		// When retrieving ordinary hbar transfers, this field does not get populated
		Token string `json:"token_id"`
	}
	// Response struct used by the Hedera Mirror node REST API and returned once
	// account transactions are queried
	Response struct {
		Transactions []Transaction
		Status       `json:"_status"`
	}
)

type (
	ErrorMessage struct {
		Message string `json:"message"`
	}
	Status struct {
		Messages []ErrorMessage
	}
)

type HCSClient struct {
	client  *hedera.Client
	topicID hedera.TopicID
	http    http.Client
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
	initialTimeStamp := time.Now()
	go c.beginWatching(receiver, string(initialTimeStamp.Second()))

	return nil
}

func (c *HCSClient) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

func (c *HCSClient) getMessagesAfterTimestamp(initialTimeStamp string) ([]Message, error) {
	query := fmt.Sprintf("/%s/messages?timestamp=gt:%s",
		c.topicID.String(),
		initialTimeStamp)

	messagesQuery := fmt.Sprintf("%s%s%s", "https://testnet.mirrornode.hedera.com/api/v1/", "topics", query)
	response, e := c.http.Get(messagesQuery)
	if e != nil {
		return nil, e
	}

	bodyBytes, e := readResponseBody(response)
	if e != nil {
		return nil, e
	}

	var messages *Messages
	e = json.Unmarshal(bodyBytes, &messages)
	if e != nil {
		return nil, e
	}
	return messages.Messages, nil
}

func (c *HCSClient) getTransaction(timestamp string) (*Response, error) {
	query := fmt.Sprintf("/%s/transactions?timestamp=gt:%s",
		c.topicID.String(),
		timestamp)

	messagesQuery := fmt.Sprintf("%s%s%s", "https://testnet.mirrornode.hedera.com/api/v1/", "topics", query)
	r, e := c.http.Get(messagesQuery)
	if e != nil {
		return nil, e
	}

	bodyBytes, e := readResponseBody(r)
	if e != nil {
		return nil, e
	}

	var response *Response
	e = json.Unmarshal(bodyBytes, &response)
	if e != nil {
		return nil, e
	}
	if r.StatusCode >= 400 {
		return response, errors.New(fmt.Sprintf(`Failed to execute query: [%s]. Error: [%s]`, query, response.Status))
	}

	return response, nil
}

func (c *HCSClient) beginWatching(receiver common.MessageReceiver, currentTimeStamp string) {
	for {
		messages, err := c.getMessagesAfterTimestamp(currentTimeStamp)
		if err != nil {
			log.Infof("Error while retrieving messages from mirror node. Error [%s]", err)
			go c.beginWatching(receiver, currentTimeStamp)
			return
		}

		log.Infof("Polling found [%d] Messages", len(messages))

		for _, msg := range messages {
			log.Infof("[HCS] The topic response received %s\n", msg)

			response, err := c.getTransaction(msg.ConsensusTimestamp)
			if err != nil {
				// handle error
			}
			var txID string
			for _, t := range response.Transactions {
				txID = t.TransactionID
			}

			txId := prepareTxId(fmt.Sprintf("%v", txID))
			sequenceNumber := strconv.FormatUint(uint64(msg.SequenceNumber), 10)

			log.Infof("[HCS] The topic response received - TransactionID: %s\n", txId)
			log.Infof("[HCS] The topic response received - SequenceNumber ID: %s\n", sequenceNumber)

			dltContextValues := DLTValues{map[string]string{
				SequenceNumberKey: sequenceNumber,
				TransactionIdKey:  txId,
			}}

			ctx := context.WithValue(context.Background(), DLTValuesKey, dltContextValues)
			receiver.Receive(&common.Message{Msg: []byte(msg.Contents), Ctx: ctx})
		}
		time.Sleep(5 * time.Second)
	}
}

func readResponseBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
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

	return &HCSClient{client: client, topicID: hcsTopicId, http: http.Client{}}
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
