package libp2p

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const p2pStreamName = "/hcs-int-p2p-nodes/1.0.0"

type LibP2PClient struct {
	h           host.Host
	receiver    common.MessageReceiver
	streamPairs map[peer.ID]*bufio.ReadWriter
}

func handleIncommingMessage(c *LibP2PClient, connection network.Conn, receiver common.MessageReceiver) {
	go func() {
		for {
			msg, err := c.streamPairs[connection.RemotePeer()].ReadBytes('\n')
			if err != nil {
				c.streamPairs[connection.RemotePeer()] = nil
				return
			}

			rawMessage, validData := c.verifySignedMessage(msg)

			if !validData {
				log.Error("The data integrity is suspicious")
				continue
			}

			remotePeerInfo := peer.AddrInfo{
				ID:    connection.RemotePeer(),
				Addrs: []multiaddr.Multiaddr{connection.RemoteMultiaddr()},
			}

			addrs, err := peer.AddrInfoToP2pAddrs(&remotePeerInfo)
			if err != nil {
				panic(err)
			}

			receiver.Receive(&common.Message{Ctx: context.WithValue(context.Background(), "remotePeerAddress", addrs[0]), Msg: rawMessage})
		}
	}()
}

func (c *LibP2PClient) Listen(receiver common.MessageReceiver) error {
	c.receiver = receiver

	c.h.SetStreamHandler(p2pStreamName, func(s network.Stream) { // I'm waiting for incomming connection
		log.Infof("%s connected with you\n", s.Conn().RemotePeer())
		c.streamPairs[s.Conn().RemotePeer()] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		handleIncommingMessage(c, s.Conn(), receiver)
	})
	return nil
}

func (c *LibP2PClient) Send(msg *common.Message, peerAddress string) error {
	signedMessage, err := c.prepareSignedMessage(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	peerInfo, err := c.multiAddrToPeerInfo(peerAddress)
	if err != nil {
		log.Errorln(err)
		return err
	}

	c.streamPairs[peerInfo.ID].Write(signedMessage)
	c.streamPairs[peerInfo.ID].Flush()
	return nil
}

func (c *LibP2PClient) verifySignedMessage(msg []byte) ([]byte, bool) {
	var envelope common.Envelope
	err := json.Unmarshal(msg, &envelope)
	if err != nil {
		log.Error(err)
		return nil, false
	}

	isValid := c.verifyData(envelope.Payload, envelope.Signature, envelope.PeerId, envelope.PubKeyData)

	return envelope.Payload, isValid
}

func (c *LibP2PClient) prepareSignedMessage(msg *common.Message) ([]byte, error) {
	signature, err := c.signData(msg.Msg)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	nodePubKey, _ := c.h.Peerstore().PubKey(c.h.ID()).Bytes()
	peerId := c.h.ID().Pretty()

	p2pMessage := &common.Envelope{Signature: signature, PubKeyData: nodePubKey, PeerId: peerId, Payload: msg.Msg}

	signedMessage := EncodeToBytes(p2pMessage)

	return signedMessage, nil
}

// sign binary data using the local node's private key
func (c *LibP2PClient) signData(data []byte) ([]byte, error) {
	key := c.h.Peerstore().PrivKey(c.h.ID())
	res, err := key.Sign(data)
	return res, err
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerId: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (c *LibP2PClient) verifyData(data []byte, signature []byte, peerId string, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Error(err, "Failed to extract key from message key data")
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		log.Error(err, "Failed to extract peer id from public key")
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey.Pretty() != peerId {
		log.Error(err, "Node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Error(err, "Error authenticating data")
		return false
	}

	return res
}

func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Error(err)
	}

	return buf.Bytes()
}

func (c *LibP2PClient) Close() error {
	return nil
}

func (c *LibP2PClient) multiAddrToPeerInfo(peerMultiAddr string) (*peer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(peerMultiAddr)
	if err != nil {
		log.Errorln(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Errorln(err)
	}

	return info, err
}

func (c *LibP2PClient) Connect(peerAddress string) (bool, error) {
	if c.receiver == nil {
		errMsg := "Missing a client receiver. Listen function should be executed, firstly."
		log.Errorln(errMsg)
		return false, errors.New(errMsg)
	}

	ai, err := c.multiAddrToPeerInfo(peerAddress)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	c.h.Peerstore().AddAddrs(ai.ID, ai.Addrs, peerstore.TempAddrTTL)

	log.Printf("This is a conversation between %s and %s\n", c.h.ID(), ai.ID)

	s, err := c.h.NewStream(context.Background(), ai.ID, p2pStreamName)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	handleIncommingMessage(c, s.Conn(), c.receiver)

	c.streamPairs[ai.ID] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return true, nil
}

func NewLibP2PClient(key ed25519.PrivateKey, listenIp string, listenPort string) *LibP2PClient {
	libp2pKey, err := crypto.UnmarshalEd25519PrivateKey(key)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", listenIp, listenPort)) //

	h, err := libp2p.New(context.Background(), libp2p.ListenAddrs(sourceMultiAddr), libp2p.Identity(libp2pKey))
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	log.Infof("[LIBP2P] Started libp2p host and listening on: %s \n", addrs[0])

	client := &LibP2PClient{h: h}
	client.streamPairs = make(map[peer.ID]*bufio.ReadWriter)

	return client
}
