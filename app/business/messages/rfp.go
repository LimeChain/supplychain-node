package messages

import (
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
)

type RFPMessage struct {
	BusinessMessage
	Data model.Product `json:"data"`
}

func CreateRFPMessage(rfp *model.Product) *RFPMessage {
	return &RFPMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeRFP}, Data: *rfp}
}
