package handler

import (
	"encoding/json"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/product/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type ProductHandler struct {
	productRepo repository.ProductRepository
}

func (h *ProductHandler) Handle(msg *common.Message) error {

	var productMsg messages.ProductMessage
	err := json.Unmarshal(msg.Msg, &productMsg)
	if err != nil {
		return err
	}
	productId, err := h.productRepo.Save(&productMsg.Data)
	if err != nil {
		return err
	}
	log.Infof("Saved product with id: %d\n", productId)
	return nil
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo}
}
