package messages

type BusinessMessage struct {
	Type string `json:"type"`
}

const (
	P2PMessageTypeRFP              = "rfp"
	P2PMessageTypeProposal         = "proposal"
	P2PMessageTypeContractRequest  = "contractrequest"
	P2PMessageTypeContractAccepted = "contractaccepted"
	P2PMessageTypePORequest        = "porequest"
	P2PMessageTypePOAccepted       = "poaccepted"

	P2PMessageTypeSendShipmentRequest  = "sendshipmentrequest"
	P2PMessageTypeSendShipmentAccepted = "sendshipmentaccepted"
)

const (
	DLTMessageTypeContract         = "contract"
	DLTMessageTypePO               = "po"
	DLTMessageTypeSendShipment     = "send-shipment"
	DLTMessageTypeReceivedShipment = "receive-shipment"
)
