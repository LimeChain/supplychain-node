package messages

type BusinessMessage struct {
	Type string `json:"type"`
}

const (
	P2PMessageTypeProduct              = "product"
	P2PMessageTypeSendShipmentRequest  = "sendshipmentrequest"
	P2PMessageTypeSendShipmentAccepted = "sendshipmentaccepted"
)

const (
	DLTMessageTypeSendShipment = "send-shipment"
	// DLTMessageTypeReceivedShipment = "receive-shipment"
)
