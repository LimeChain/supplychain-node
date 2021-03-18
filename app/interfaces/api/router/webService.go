package router

import (
    "bytes"
    "encoding/json"
    "net/http"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func getWebServiceTest(webService *apiservices.WebService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		
		values := map[string]string{"ac": "a", "occupation": "gardener"}
		json_data, err := json.Marshal(values)
	
		if err != nil {
			log.Fatal(err)
		}
	
		resp, err := http.Post("http://hedera-web-service:11136/api/shipment", "application/json",
			bytes.NewBuffer(json_data))
	
		if err != nil {
			log.Println("Error:")
			log.Fatal(err)
		}
	
		var res map[string]interface{}
	
		json.NewDecoder(resp.Body).Decode(&res)
		
		log.Println("Response:")
		log.Println(res["json"])	}
}

func NewWebServiceRouter(webService *apiservices.WebService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getWebServiceTest(webService))
	return r
}
