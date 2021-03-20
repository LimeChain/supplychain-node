package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/domain/product/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type CreateProductRequest struct {
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

type storedProductsResponse struct {
	api.IntegrationNodeAPIResponse
	Products []*model.Product `json:"products"`
}

type storedProductResponse struct {
	api.IntegrationNodeAPIResponse
	Product *model.Product `json:"product"`
}

type createProductResponse struct {
	api.IntegrationNodeAPIResponse
	ProductId int `json:"productId,omitempty"`
}

func (productRequestModel *CreateProductRequest) toProduct() *model.Product {
	productModel := model.ProductModel{
		ProductId:          productRequestModel.Obj.ProductModel.ProductId,
		ProductName:        productRequestModel.Obj.ProductModel.ProductName,
		ProductUnit:        productRequestModel.Obj.ProductModel.ProductUnit,
		ProductDescription: productRequestModel.Obj.ProductModel.ProductDescription,
		ProductDeleted:     productRequestModel.Obj.ProductModel.ProductDeleted,
	}

	productObj := model.ProductObj{
		ProductModel: productModel,
	}

	return &model.Product{
		Type:        productRequestModel.Type,
		Obj:         productObj,
		Destination: productRequestModel.Destination,
	}
}

func getAllStoredProducts(productService *apiservices.ProductService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedProducts, err := productService.GetAllProducts()
		if err != nil {
			render.JSON(w, r, storedProductsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedProductsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedProducts})
	}
}

func getProductById(productService *apiservices.ProductService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		productIdParam := chi.URLParam(r, "productId")
		productId, _ := strconv.Atoi(productIdParam)
		product, err := productService.GetProduct(productId)
		if err != nil {
			render.JSON(w, r, storedProductResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedProductResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, product})
	}
}

func createProduct(productService *apiservices.ProductService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var productRequest *CreateProductRequest

		err := parser.DecodeJSONBody(w, r, &productRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, createProductResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, 0})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, createProductResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0})
			return
		}

		// ToDo: Validate decoded struct

		product := productRequest.toProduct()

		storedProductId, err := productService.CreateProduct(product)
		if err != nil {
			render.JSON(w, r, createProductResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, 0})
			return
		}

		render.JSON(w, r, createProductResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedProductId})
	}
}

func NewProductRouter(productService *apiservices.ProductService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredProducts(productService))
	r.Get("/{productId}", getProductById(productService))
	r.Post("/", createProduct(productService))
	return r
}
