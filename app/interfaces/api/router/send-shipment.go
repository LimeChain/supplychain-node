package router

import (
	"errors"
	"net/http"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	sendShipmentModel "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type SendSendShipmentRequest struct {
	ShipmentId  string `json:"shipmentId" bson:"contractId"`
	SupplierId  string `json:"supplierId" bson:"supplierId"`
	BuyerId     string `json:"buyerId" bson:"buyerId"`
	Destination string `json:"destination" bson:"destination"`
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
	ShipmentId            string `json:"shipmentId, omitempty" bson:"shipmentId"`
	SendShipmentHash      string `json:"sendShipmentHash, omitempty" bson:"sendShipmentHash"`
	SendShipmentSignature string `json:"sendShipmentSignature, omitempty" bson:"sendShipmentSignature"`
}

func (req *SendSendShipmentRequest) toUnsignedSendShipment() *sendShipmentModel.UnsignedSendShipment {
	return &sendShipmentModel.UnsignedSendShipment{
		ShipmentId:  req.ShipmentId,
		SupplierId:  req.SupplierId,
		BuyerId:     req.BuyerId,
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
		shipmentId := chi.URLParam(r, "shipmentId")
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
				render.JSON(w, r, sendSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, "", "", ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, sendSendShipmentResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
			return
		}

		// ToDo: Validate decoded struct

		unsignedSendShipment := sendShipmentRequest.toUnsignedSendShipment()

		shipmentId, sendShipmentHash, sendShipmentSignature, err := sendShipmentService.SaveAndSendSendShipment(unsignedSendShipment)
		if err != nil {
			render.JSON(w, r, sendContractResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
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

func NewSendShipmentRouter(sendShipmentService *apiservices.SendShipmentService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredSentShipments(sendShipmentService))
	r.Get("/{shipmentId}", getSendShipmentById(sendShipmentService))
	r.Post("/", sendSendShipment(sendShipmentService))
	return r
}
