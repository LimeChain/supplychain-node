package router

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/go-chi/chi"
)

func getWebServiceTest(webService *apiservices.WebPlatformService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// values := map[string]string{"ac": "a", "occupation": "gardener"}
		// json_data, err := json.Marshal(values)

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// resp, err := http.Post("http://hedera-web-service:11136/api/shipment", "application/json",
		// 	bytes.NewBuffer(json_data))

		// if err != nil {
		// 	log.Println("Error:")
		// 	log.Fatal(err)
		// }

		// var res map[string]interface{}

		// json.NewDecoder(resp.Body).Decode(&res)

		// log.Println("Response:")
		// log.Println(res["json"])
		w.Write([]byte("PING"))
	}
}

func postCreditShipment(webService *apiservices.WebPlatformService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		resp := webService.CreditShipment(bytes)
		w.Write([]byte(resp))
	}
}

func postCreditProduct(webService *apiservices.WebPlatformService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		webService.CreditProduct(bytes)
		w.Write([]byte("Done"))
	}
}

func NewWebPlatformRouter(webService *apiservices.WebPlatformService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getWebServiceTest(webService))
	r.Post("/credit-shipment", postCreditShipment(webService))
	r.Post("/credit-product", postCreditProduct(webService))
	return r
}
