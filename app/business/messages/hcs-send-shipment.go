package messages

type dltSendShipment struct {
	ShipmentHash string `json:"shipmentHash"`
}

type DLTSendShipmentMessage struct {
	BusinessMessage
	Data dltSendShipment `json:"data"`
}

func CreateDLTSendShipmentMessage(shipmentHash string) *DLTSendShipmentMessage {
	return &DLTSendShipmentMessage{
		BusinessMessage: BusinessMessage{Type: DLTMessageTypeSendShipment},
		Data: dltSendShipment{
			ShipmentHash: shipmentHash,
		},
	}
}
