package messages

type dltSendShipment struct {
	ShipmentId        string `json:"shipmentId"`
	SendShipmentHash  string `json:"sendShipmentHash"`
	BuyerSignature    string `json:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature"`
}

type DLTSendShipmentMessage struct {
	BusinessMessage
	Data dltSendShipment `json:"data"`
}

func CreateDLTSendShipmentMessage(shipmentId, sendShipmentHash, buyerSignature, supplierSignature string) *DLTSendShipmentMessage {
	return &DLTSendShipmentMessage{
		BusinessMessage: BusinessMessage{Type: DLTMessageTypeSendShipment},
		Data: dltSendShipment{
			ShipmentId:        shipmentId,
			SendShipmentHash:  sendShipmentHash,
			BuyerSignature:    buyerSignature,
			SupplierSignature: supplierSignature,
		},
	}
}
