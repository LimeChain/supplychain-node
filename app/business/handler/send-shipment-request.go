package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type SendShipmentRequestHandler struct {
	sendShipmentRepo    repository.SendShipmentRepository
	sendShipmentService *service.SendShipmentService
	p2pClient           common.Messenger
}

func (h *SendShipmentRequestHandler) Handle(msg *common.Message) error {
	remotePeerAddressCtx := msg.Ctx.Value("remotePeerAddress")

	if remotePeerAddressCtx == nil {
		return errors.New("The remote peer address is missing")
	}

	remotePeerAddress := fmt.Sprintf("%v", remotePeerAddressCtx)

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

	signatureCorrect, err := h.sendShipmentService.VerifyBuyer(&sendShipment)
	if err != nil {
		return err
	}

	if !signatureCorrect {
		return errors.New("Invalid signature by the buyer")
	}

	sendShipmentSignature, err := h.sendShipmentService.Sign(&sendShipment.UnsignedSendShipment)
	if err != nil {
		return err
	}

	sendShipment.SupplierSignature = sendShipmentSignature

	shipmentId, err := h.sendShipmentRepo.Save(&sendShipment)
	if err != nil {
		return err
	}

	sendShipmentAcceptedMsg := messages.CreateSendShipmentAcceptedMessage(&sendShipment)

	p2pBytes, err := json.Marshal(sendShipmentAcceptedMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}
	h.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, remotePeerAddress)

	log.Infof("Verified and saved shipment with id: %s\n", shipmentId)
	return nil
}

func NewSendShipmentRequestHandler(
	sendShipmentRepo repository.SendShipmentRepository,
	sendShipmentService *service.SendShipmentService,
	p2pClient common.Messenger) *SendShipmentRequestHandler {
	return &SendShipmentRequestHandler{sendShipmentRepo: sendShipmentRepo, sendShipmentService: sendShipmentService, p2pClient: p2pClient}
}
