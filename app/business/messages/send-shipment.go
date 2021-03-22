package messages

import "github.com/Limechain/pwc-bat-node/app/domain/send-shipment/model"

type SendShipmentMessage struct {
	BusinessMessage
	Data model.SendShipment `json:"data"`
}

func CreateSendShipmentRequestMessage(sendShipment *model.SendShipment) *SendShipmentMessage {
	return &SendShipmentMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeSendShipmentRequest}, Data: *sendShipment}
}

func CreateSendShipmentAcceptedMessage(sendShipment *model.SendShipment) *SendShipmentMessage {
	return &SendShipmentMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeSendShipmentAccepted}, Data: *sendShipment}
}
