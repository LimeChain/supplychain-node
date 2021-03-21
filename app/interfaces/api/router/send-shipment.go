package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	sendShipmentModel "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type ValidateShipmentMessageRequest struct {
	Message string `json:"message" bson:"message"`
}

type validateShipmentMessageResponse struct {
	api.IntegrationNodeAPIResponse
	Shipment          sendShipmentModel.UnsignedSendShipment `json:"shipment, omitempty" bson:"shipment"`
	BuyerSignature    string                                 `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature string                                 `json:"supplierSignature" bson:"supplierSignature"`
}

type SendSendShipmentRequest struct {
	Type        int         `json:"type" bson:"type"`
	Obj         shipmentObj `json:"obj" bson:"obj"`
	Destination string      `json:"destination" bson:"destination"`
}

type shipmentObj struct {
	ShipmentModel          shipmentModel           `json:"shipmentModel" bson:"shipmentModel"`
	SkuModels              []skuModel              `json:"skuModels" bson:"skuModels"`
	SkuOriginModels        []skuOriginModel        `json:"skuOriginModels" bson:"skuOriginModels"`
	ShipmentDocumentModels []shipmentDocumentModel `json:"shipmentDocumentModels" bson:"shipmentDocumentModels"`
}

type shipmentModel struct {
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

type skuModel struct {
	SkuId        int `json:"skuId" bson:"skuId"`
	ShipmentId   int `json:"shipmentId" bson:"shipmentId"`
	ProductId    int `json:"productId" bson:"productId"`
	Quantity     int `json:"quantity" bson:"quantity"`
	PricePerUnit int `json:"pricePerUnit" bson:"pricePerUnit"`
	Currency     int `json:"currency" bson:"currency"`
}

type skuOriginModel struct {
	SkuOriginId int `json:"skuOriginId" bson:"skuOriginId"`
	SkuId       int `json:"skuId" bson:"skuId"`
	ShipmentId  int `json:"shipmentId" bson:"shipmentId"`
}

type shipmentDocumentModel struct {
	ShipmentDocumentId  int    `json:"shipmentDocumentId" bson:"shipmentDocumentId"`
	ShipmentId          int    `json:"shipmentId" bson:"shipmentId"`
	DocumentType        int    `json:"documentType" bson:"documentType"`
	MimeType            string `json:"mimeType" bson:"mimeType"`
	ShipmentDocumentUrl string `json:"shipmentDocumentUrl" bson:"shipmentDocumentUrl"`
	SizeInBytes         int    `json:"sizeInBytes" bson:"sizeInBytes"`
	Name                string `json:"name" bson:"name"`
}

type storedSentShipmentsResponse struct {
	api.IntegrationNodeAPIResponse
	SentShipments []*sendShipmentModel.SendShipment `json:"send-shipments"`
}

type storedSendShipmentResponse struct {
	api.IntegrationNodeAPIResponse
	SendShipment *sendShipmentModel.SendShipment `json:"send-shipment"`
}

type sendSendShipmentResponse struct {
	api.IntegrationNodeAPIResponse
	ShipmentId            int    `json:"shipmentId, omitempty" bson:"shipmentId"`
	SendShipmentHash      string `json:"sendShipmentHash, omitempty" bson:"sendShipmentHash"`
	SendShipmentSignature string `json:"sendShipmentSignature, omitempty" bson:"sendShipmentSignature"`
}

func (req *SendSendShipmentRequest) toUnsignedSendShipment() *sendShipmentModel.UnsignedSendShipment {
	shipmentDocuments := make([]sendShipmentModel.ShipmentDocumentModel, len(req.Obj.ShipmentDocumentModels))

	for i, doc := range req.Obj.ShipmentDocumentModels {
		shipmentDocuments[i] = sendShipmentModel.ShipmentDocumentModel{
			ShipmentDocumentId:  doc.ShipmentDocumentId,
			ShipmentId:          doc.ShipmentId,
			DocumentType:        doc.DocumentType,
			MimeType:            doc.MimeType,
			ShipmentDocumentUrl: doc.ShipmentDocumentUrl,
			SizeInBytes:         doc.SizeInBytes,
			Name:                doc.Name,
		}
	}

	skuOriginModels := make([]sendShipmentModel.SkuOriginModel, len(req.Obj.SkuOriginModels))

	for i, skuOrigin := range req.Obj.SkuOriginModels {
		skuOriginModels[i] = sendShipmentModel.SkuOriginModel{
			SkuOriginId: skuOrigin.SkuOriginId,
			SkuId:       skuOrigin.SkuId,
			ShipmentId:  skuOrigin.ShipmentId,
		}
	}

	skuModels := make([]sendShipmentModel.SkuModel, len(req.Obj.SkuModels))

	for i, sku := range req.Obj.SkuModels {
		skuModels[i] = sendShipmentModel.SkuModel{
			SkuId:        sku.SkuId,
			ShipmentId:   sku.ShipmentId,
			ProductId:    sku.ProductId,
			Quantity:     sku.Quantity,
			PricePerUnit: sku.PricePerUnit,
			Currency:     sku.Currency,
		}
	}

	shipmentModel := sendShipmentModel.ShipmentModel{
		ShipmentId:                req.Obj.ShipmentModel.ShipmentId,
		ShipmentConsignmentNumber: req.Obj.ShipmentModel.ShipmentConsignmentNumber,
		ShipmentName:              req.Obj.ShipmentModel.ShipmentName,
		ShipmentStatus:            req.Obj.ShipmentModel.ShipmentStatus,
		ShipmentOriginSiteId:      req.Obj.ShipmentModel.ShipmentOriginSiteId,
		ShipmentDestinationSiteId: req.Obj.ShipmentModel.ShipmentDestinationSiteId,
		ShipmentDateOfShipment:    req.Obj.ShipmentModel.ShipmentDateOfShipment,
		ShipmentDateOfArrival:     req.Obj.ShipmentModel.ShipmentDateOfArrival,
		ShipmentDltAnchored:       req.Obj.ShipmentModel.ShipmentDltAnchored,
		ShipmentDltProof:          req.Obj.ShipmentModel.ShipmentDltProof,
		ShipmentDeleted:           req.Obj.ShipmentModel.ShipmentDeleted,
	}

	shipmentObject := sendShipmentModel.ShipmentObj{
		ShipmentModel:          shipmentModel,
		SkuModels:              skuModels,
		SkuOriginModels:        skuOriginModels,
		ShipmentDocumentModels: shipmentDocuments,
	}

	return &sendShipmentModel.UnsignedSendShipment{
		Type:        req.Type,
		Obj:         shipmentObject,
		Destination: req.Destination,
	}
}

func getAllStoredSentShipments(sendShipmentService *apiservices.SendShipmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedSentShipments, err := sendShipmentService.GetAllSentShipments()
		if err != nil {
			render.JSON(w, r, storedSentShipmentsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedSentShipmentsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedSentShipments})
	}
}

func getSendShipmentById(sendShipmentService *apiservices.SendShipmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		shipmentIdParam := chi.URLParam(r, "shipmentId")
		shipmentId, _ := strconv.Atoi(shipmentIdParam)
		storedSendShipment, err := sendShipmentService.GetSentShipment(shipmentId)
		if err != nil {
			render.JSON(w, r, storedSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedSendShipment})
	}
}

func sendSendShipment(sendShipmentService *apiservices.SendShipmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendShipmentRequest *SendSendShipmentRequest

		err := parser.DecodeJSONBody(w, r, &sendShipmentRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, sendSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, 0, "", ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, sendSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0, "", ""})
			return
		}

		// ToDo: Validate decoded struct

		unsignedSendShipment := sendShipmentRequest.toUnsignedSendShipment()

		shipmentId, sendShipmentHash, sendShipmentSignature, err := sendShipmentService.SaveAndSendSendShipment(unsignedSendShipment)
		if err != nil {
			render.JSON(w, r, sendSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0, "", ""})
			return
		}

		render.JSON(w, r, sendSendShipmentResponse{
			IntegrationNodeAPIResponse: api.IntegrationNodeAPIResponse{Status: true, Error: ""},
			ShipmentId:                 shipmentId,
			SendShipmentHash:           sendShipmentHash,
			SendShipmentSignature:      sendShipmentSignature,
		})
	}
}

func validateShipmentMessage(sendShipmentService *apiservices.SendShipmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var validateShipment *ValidateShipmentMessageRequest

		err := parser.DecodeJSONBody(w, r, &validateShipment)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, validateShipmentMessageResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, sendShipmentModel.UnsignedSendShipment{}, "", ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, validateShipmentMessageResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, sendShipmentModel.UnsignedSendShipment{}, "", ""})
			return
		}

		storedSendShipment, err := sendShipmentService.GetSentShipmentByDLTMessage(validateShipment.Message)

		if storedSendShipment == nil {
			render.JSON(w, r, validateShipmentMessageResponse{api.IntegrationNodeAPIResponse{Status: false, Error: "Invalid DLT message <> stored shipment info"}, sendShipmentModel.UnsignedSendShipment{}, "", ""})
			return
		}

		render.JSON(w, r, validateShipmentMessageResponse{
			IntegrationNodeAPIResponse: api.IntegrationNodeAPIResponse{Status: true, Error: ""},
			Shipment:                   storedSendShipment.UnsignedSendShipment,
			BuyerSignature:             storedSendShipment.BuyerSignature,
			SupplierSignature:          storedSendShipment.SupplierSignature,
		})
	}
}

func NewSendShipmentRouter(sendShipmentService *apiservices.SendShipmentService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredSentShipments(sendShipmentService))
	r.Get("/{shipmentId}", getSendShipmentById(sendShipmentService))
	r.Post("/", sendSendShipment(sendShipmentService))
	r.Post("/validate", validateShipmentMessage(sendShipmentService))
	return r
}
