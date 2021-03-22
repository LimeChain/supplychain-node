package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Limechain/pwc-bat-node/app/business/messages"
	"github.com/Limechain/pwc-bat-node/app/domain/send-shipment/repository"
	"github.com/Limechain/pwc-bat-node/app/domain/send-shipment/service"
	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type SendShipmentAcceptedHandler struct {
	sendShipmentRepo    repository.SendShipmentRepository
	dltClient           common.DLTMessenger
	sendShipmentService *service.SendShipmentService
}

func (h *SendShipmentAcceptedHandler) Handle(msg *common.Message) error {

	var sendShipmentMsg messages.SendShipmentMessage
	err := json.Unmarshal(msg.Msg, &sendShipmentMsg)
	if err != nil {
		return err
	}

	sendShipment := sendShipmentMsg.Data

	// TODO add more validation

	if len(sendShipment.BuyerSignature) == 0 {
		return errors.New("The sent shipment was not signed by the buyer")
	}

	if len(sendShipment.SupplierSignature) == 0 {
		return errors.New("The sent shipment was not signed by the supplir")
	}

	savedSendShipment, err := h.sendShipmentRepo.GetByID(sendShipment.Obj.ShipmentModel.ShipmentId)
	if err != nil {
		return err
	}

	if savedSendShipment.BuyerSignature != savedSendShipment.BuyerSignature {
		return errors.New("The sent shipment buyer signature was not the one stored. The supplier has tried to cheat you")
	}

	signatureCorrect, err := h.sendShipmentService.VerifySupplier(&sendShipment)
	if err != nil {
		return err
	}

	if !signatureCorrect {
		return errors.New("Invalid signature by the supplier")
	}

	dataAndSignaturesHash := h.sendShipmentService.HashDataAndSignatures(&sendShipment.UnsignedSendShipment, sendShipment.BuyerSignature, sendShipment.SupplierSignature)
	sendShipment.SignedDataHash = dataAndSignaturesHash

	err = h.sendShipmentRepo.Update(&sendShipment)
	if err != nil {
		return err
	}

	dltMessage := messages.CreateDLTSendShipmentMessage(dataAndSignaturesHash)

	dltBytes, err := json.Marshal(dltMessage)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}

	err = h.dltClient.Send(&common.Message{Ctx: context.TODO(), Msg: dltBytes})
	if err != nil {
		return err
	}

	log.Infof("Verified and saved accepted sent shipment with id: %d\n", sendShipment.Obj.ShipmentModel.ShipmentId)
	return nil
}

func NewSendShipmentAcceptedHandler(sendShipmentRepo repository.SendShipmentRepository, sendShipmentService *service.SendShipmentService, dltClient common.DLTMessenger) *SendShipmentAcceptedHandler {
	return &SendShipmentAcceptedHandler{sendShipmentRepo: sendShipmentRepo, sendShipmentService: sendShipmentService, dltClient: dltClient}
}
