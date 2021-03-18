package model

type UnsignedSendShipment struct {
	ShipmentId  string `json:"shipmentId" bson:"shipmentId"`
	SupplierId  string `json:"supplierId" bson:"supplierId"`
	BuyerId     string `json:"buyerId" bson:"buyerId"`
	Destination string `json:"destination" bson:"destination"`
}

type SendShipment struct {
	UnsignedSendShipment `json:"unsignedSendShipment" bson:"unsignedSendShipment"`
	BuyerSignature       string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature    string `json:"supplierSignature" bson:"supplierSignature"`
	DLTAnchored          bool   `json:"DLTAnchored" bson:"DLTAnchored"`
	DLTProof             string `json:"DLTProof" bson:"DLTProof"`
}
