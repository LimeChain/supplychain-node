package apiservices

import (
	contractService "github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
)

type WebService struct {
	ws        *contractService.ContractService
}

func (w *ContractService) testWebService() ([]*string) {
	var a[]*string
	return a
}

func NewWebService(ws *contractService.ContractService) *WebService{
	return &WebService{ws}
}
