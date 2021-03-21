package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
)

type SendShipmentService struct {
	signingKey ed25519.PrivateKey
	peerPubKey ed25519.PublicKey
}

func (ss *SendShipmentService) Hash(unsignedSendShipment *model.UnsignedSendShipment) (string, error) {
	var sb strings.Builder

	unsignedSendShipmentBytes := []byte(fmt.Sprintf("%v", unsignedSendShipment))

	sb.Write(unsignedSendShipmentBytes)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String()))), nil
}

func (ss *SendShipmentService) HashDataAndSignatures(unsignedSendShipment *model.UnsignedSendShipment, buyerSignature, sellerSignature string) string {
	var sb strings.Builder

	unsignedSendShipmentBytes := []byte(fmt.Sprintf("%v", unsignedSendShipment))

	sb.Write(unsignedSendShipmentBytes)
	sb.WriteRune(',')
	sb.WriteString(buyerSignature)
	sb.WriteRune(',')
	sb.WriteString(sellerSignature)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String())))
}

func (ss *SendShipmentService) Sign(unsignedSendShipment *model.UnsignedSendShipment) (string, error) {
	sendShipmentHash, err := ss.Hash(unsignedSendShipment)
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(ss.signingKey, []byte(sendShipmentHash))
	signatureStr := hex.EncodeToString(signature)
	return signatureStr, nil
}

func (ss *SendShipmentService) verify(unsignedSendShipment *model.UnsignedSendShipment, signature string) (bool, error) {
	sendShipmentHash, err := ss.Hash(unsignedSendShipment)
	if err != nil {
		return false, err
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(ss.peerPubKey, []byte(sendShipmentHash), signatureBytes), nil
}

func (ss *SendShipmentService) VerifyBuyer(sendShipment *model.SendShipment) (bool, error) {
	return ss.verify(&sendShipment.UnsignedSendShipment, sendShipment.BuyerSignature)
}

func (ss *SendShipmentService) VerifySupplier(sendShipment *model.SendShipment) (bool, error) {
	return ss.verify(&sendShipment.UnsignedSendShipment, sendShipment.SupplierSignature)
}

func New(signingKey ed25519.PrivateKey, peerPubKey ed25519.PublicKey) *SendShipmentService {
	return &SendShipmentService{signingKey: signingKey, peerPubKey: peerPubKey}
}
