package model

type UnsignedSendShipment struct {
	Type        int         `json:"type" bson:"type"`
	Obj         ShipmentObj `json:"obj" bson:"obj"`
	Destination string      `json:"destination" bson:"destination"`
}

type ShipmentObj struct {
	ShipmentModel          ShipmentModel           `json:"shipmentModel" bson:"shipmentModel"`
	SkuModels              []SkuModel              `json:"skuModels" bson:"skuModels"`
	SkuOriginModels        []SkuOriginModel        `json:"skuOriginModels" bson:"skuOriginModels"`
	ShipmentDocumentModels []ShipmentDocumentModel `json:"shipmentDocumentModels" bson:"shipmentDocumentModels"`
}

type ShipmentModel struct {
	ShipmentId                int    `json:"shipmentId" bson:"shipmentId"`
	ShipmentConsignmentNumber string `json:"shipmentConsignmentNumber" bson:"shipmentConsignmentNumber"`
	ShipmentName              string `json:"shipmentName" bson:"shipmentName"`
	ShipmentStatus            int    `json:"shipmentStatus" bson:"shipmentStatus"`
	ShipmentOriginSiteId      int    `json:"shipmentOriginSiteId" bson:"shipmentOriginSiteId"`
	ShipmentDestinationSiteId int    `json:"shipmentDestinationSiteId" bson:"shipmentDestinationSiteId"`
	ShipmentDateOfShipment    int    `json:"shipmentDateOfShipment" bson:"shipmentDateOfShipment"`
	ShipmentDateOfArrival     int    `json:"shipmentDateOfArrival" bson:"shipmentDateOfArrival"`
	ShipmentDltAnchored       int    `json:"shipmentDltAnchored" bson:"shipmentDltAnchored"`
	ShipmentDltProof          string `json:"shipmentDltProof" bson:"shipmentDltProof"`
	ShipmentDeleted           int    `json:"shipmentDeleted" bson:"shipmentDeleted"`
}

type SkuModel struct {
	SkuId        int     `json:"skuId" bson:"skuId"`
	ShipmentId   int     `json:"shipmentId" bson:"shipmentId"`
	ProductId    int     `json:"productId" bson:"productId"`
	Quantity     int     `json:"quantity" bson:"quantity"`
	PricePerUnit float32 `json:"pricePerUnit" bson:"pricePerUnit"`
	Currency     int     `json:"currency" bson:"currency"`
}

type SkuOriginModel struct {
	SkuOriginId int `json:"skuOriginId" bson:"skuOriginId"`
	SkuId       int `json:"skuId" bson:"skuId"`
	ShipmentId  int `json:"shipmentId" bson:"shipmentId"`
}

type ShipmentDocumentModel struct {
	ShipmentDocumentId  int    `json:"shipmentDocumentId" bson:"shipmentDocumentId"`
	ShipmentId          int    `json:"shipmentId" bson:"shipmentId"`
	DocumentType        int    `json:"documentType" bson:"documentType"`
	MimeType            string `json:"mimeType" bson:"mimeType"`
	ShipmentDocumentUrl string `json:"shipmentDocumentUrl" bson:"shipmentDocumentUrl"`
	SizeInBytes         int    `json:"sizeInBytes" bson:"sizeInBytes"`
	Name                string `json:"name" bson:"name"`
}

type SendShipment struct {
	UnsignedSendShipment `json:"unsignedSendShipment" bson:"unsignedSendShipment"`
	BuyerSignature       string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature    string `json:"supplierSignature" bson:"supplierSignature"`
	DLTAnchored          bool   `json:"DLTAnchored" bson:"DLTAnchored"`
	DLTProof             string `json:"DLTProof" bson:"DLTProof"`
	DLTMessage           string `json:"DLTMessage" bson:"DLTMessage"`
	DLTTransactionId     string `json:"DLTTransactionId" bson:"DLTTransactionId"`
	SignedDataHash       string `json:"signedDataHash" bson:"signedDataHash"`
}
