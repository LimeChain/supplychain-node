package handler

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/pwc-bat-node/app/business/messages"
	"github.com/Limechain/pwc-bat-node/app/domain/send-shipment/repository"
	"github.com/Limechain/pwc-bat-node/app/domain/send-shipment/service"
	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
	"github.com/Limechain/pwc-bat-node/app/interfaces/dlt/hcs"
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

	savedSendShipment, err := h.sendShipmentRepo.GetByHash(sendShipment.ShipmentHash)

	if savedSendShipment == nil {
		return errors.New("Saved hash is different then the DLT anchored")
	}

	if err != nil {
		return err
	}

	sn := msg.Ctx.Value(hcs.SequenceNumberKey)

	sequenceNumber, ok := sn.(uint64)
	if !ok {
		return errors.New("Could not get the proof sequence number")
	}

	savedSendShipment.DLTAnchored = true
	savedSendShipment.DLTProof = fmt.Sprintf("%d", sequenceNumber)
	savedSendShipment.DLTMessage = hex.EncodeToString(msg.Msg)

	err = h.sendShipmentRepo.Update(savedSendShipment)
	if err != nil {
		return err
	}

	log.Infof("Sent shipment with Id: %d seen in the dlt and verified\n", savedSendShipment.Obj.ShipmentModel.ShipmentId)
	return nil
}

func NewDLTSendShipmentHandler(sendShipmentRepo repository.SendShipmentRepository, sendShipmentService *service.SendShipmentService) *DLTSendShipmentHandler {
	return &DLTSendShipmentHandler{sendShipmentRepo: sendShipmentRepo, sendShipmentService: sendShipmentService}
}
