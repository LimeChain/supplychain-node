package router

import (
	"net/http"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func getWebServiceTest(webService *apiservices.WebService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("wegwegweg")
	}
}

func NewWebServiceRouter(webService *apiservices.WebService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getWebServiceTest(webService))
	return r
}
