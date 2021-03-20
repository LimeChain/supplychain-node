package messages

import (
	"github.com/Limechain/HCS-Integration-Node/app/domain/product/model"
)

type ProductMessage struct {
	BusinessMessage
	Data model.Product `json:"data"`
}

func CreateProductMessage(product *model.Product) *ProductMessage {
	return &ProductMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeProduct}, Data: *product}
}
