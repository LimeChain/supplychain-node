package handler

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/dlt/hcs"
	log "github.com/sirupsen/logrus"
)

type DLTSendShipmentHandler struct {
	sendShipmentRepo    repository.SendShipmentRepository
	sendShipmentService *service.SendShipmentService
}

func (h *DLTSendShipmentHandler) Handle(msg *common.Message) error {

	var sendShipmentMsg messages.DLTSendShipmentMessage
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
		return errors.New("The sent shipment was not signed by the buyer")
	}

	savedSendShipment, err := h.sendShipmentRepo.GetByID(sendShipment.ShipmentId)
	if err != nil {
		return err
	}

	if savedSendShipment.BuyerSignature != sendShipment.BuyerSignature {
		return errors.New("The sent shipment buyer signature was not the one stored")
	}

	if savedSendShipment.SupplierSignature != sendShipment.SupplierSignature {
		return errors.New("The shipment sent supplier signature was not the one stored")
	}

	savedHash, err := h.sendShipmentService.Hash(&savedSendShipment.UnsignedSendShipment)
	if err != nil {
		return err
	}

	if savedHash != sendShipment.SendShipmentHash {
		return errors.New("The send shipment hash was not the one stored")
	}

	sn := msg.Ctx.Value(hcs.SequenceNumberKey)

	sequenceNumber, ok := sn.(uint64)
	if !ok {
		return errors.New("Could not get the proof sequence number")
	}

	savedSendShipment.DLTAnchored = true
	savedSendShipment.DLTProof = fmt.Sprintf("%d", sequenceNumber)

	err = h.sendShipmentRepo.Update(savedSendShipment)
	if err != nil {
		return err
	}

	log.Infof("Sent shipment with Id: %s seen in the dlt and verified\n", sendShipment.ShipmentId)
	return nil
}

func NewDLTSendShipmentHandler(sendShipmentRepo repository.SendShipmentRepository, sendShipmentService *service.SendShipmentService) *DLTSendShipmentHandler {
	return &DLTSendShipmentHandler{sendShipmentRepo: sendShipmentRepo, sendShipmentService: sendShipmentService}
}
