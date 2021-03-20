package router

import (
	"errors"
	"net/http"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type CreateRFPRequest struct {
	Type        int        `json:"type" bson:"type"`
	Obj         productObj `json:"obj" bson:"obj"`
	Destination string     `json:"destination" bson:"destination"`
}

type productObj struct {
	ProductModel productModel `json:"productModel" bson:"productModel"`
}

type productModel struct {
	ProductId          int    `json:"productId" bson:"productId"`
	ProductName        string `json:"productName" bson:"productName"`
	ProductUnit        int    `json:"productUnit" bson:"productUnit"`
	ProductDescription string `json:"productDescription" bson:"productDescription"`
	ProductDeleted     int    `json:"productDeleted" bson:"productDeleted"`
}

type storedRFPsResponse struct {
	api.IntegrationNodeAPIResponse
	RFPs []*model.Product `json:"rfps"`
}

type storedRFPResponse struct {
	api.IntegrationNodeAPIResponse
	RFP *model.Product `json:"rfp"`
}

type createRFPResponse struct {
	api.IntegrationNodeAPIResponse
	RFPId int `json:"rfpId,omitempty"`
}

func (rfpRequestModel *CreateRFPRequest) toRFP() *model.Product {
	productModel := model.ProductModel{
		ProductId:          rfpRequestModel.Obj.ProductModel.ProductId,
		ProductName:        rfpRequestModel.Obj.ProductModel.ProductName,
		ProductUnit:        rfpRequestModel.Obj.ProductModel.ProductUnit,
		ProductDescription: rfpRequestModel.Obj.ProductModel.ProductDescription,
		ProductDeleted:     rfpRequestModel.Obj.ProductModel.ProductDeleted,
	}

	productObj := model.ProductObj{
		ProductModel: productModel,
	}

	return &model.Product{
		Type:        rfpRequestModel.Type,
		Obj:         productObj,
		Destination: rfpRequestModel.Destination,
	}
}

func getAllStoredRFPs(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedRFPs, err := rfpService.GetAllRFPs()
		if err != nil {
			render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedRFPs})
	}
}

func getRFPById(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rfpId := chi.URLParam(r, "rfpId")
		rfp, err := rfpService.GetRFP(rfpId)
		if err != nil {
			render.JSON(w, r, storedRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedRFPResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, rfp})
	}
}

func createRFP(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rfpRequest *CreateRFPRequest

		err := parser.DecodeJSONBody(w, r, &rfpRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, 0})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0})
			return
		}

		// ToDo: Validate decoded struct

		rfp := rfpRequest.toRFP()

		storedRFPId, err := rfpService.CreateRFP(rfp)
		if err != nil {
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0})
			return
		}

		render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedRFPId})
	}
}

func NewRFPRouter(rfpService *apiservices.RFPService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredRFPs(rfpService))
	r.Get("/{rfpId}", getRFPById(rfpService))
	r.Post("/", createRFP(rfpService))
	return r
}
