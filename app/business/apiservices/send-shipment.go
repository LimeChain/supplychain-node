package apiservices

import (
	"context"
	"encoding/json"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	sendShipmentModel "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
	sendShipmentRepo "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/repository"
	sendShipmentService "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type SendShipmentService struct {
	ssr       sendShipmentRepo.SendShipmentRepository
	sss       *sendShipmentService.SendShipmentService
	p2pClient common.Messenger
}

func (ss *SendShipmentService) GetAllSentShipments() ([]*sendShipmentModel.SendShipment, error) {
	return ss.ssr.GetAll()
}

func (ss *SendShipmentService) GetSentShipment(shipmentId string) (*sendShipmentModel.SendShipment, error) {
	return ss.ssr.GetByID(shipmentId)
}

func (ss *SendShipmentService) SaveAndSendSendShipment(unsignedSendShipment *sendShipmentModel.UnsignedSendShipment) (sendShipmentId, sendShipmentHash, sendShipmentSignature string, err error) {
	sendShipmentHash, err = ss.sss.Hash(unsignedSendShipment)
	if err != nil {
		return "", "", "", err
	}
	sendShipmentSignature, err = ss.sss.Sign(unsignedSendShipment)
	if err != nil {
		return "", "", "", err
	}
	signedSendShipment := &sendShipmentModel.SendShipment{UnsignedSendShipment: *unsignedSendShipment, BuyerSignature: sendShipmentSignature, DLTAnchored: false}
	sendShipmentId, err = ss.ssr.Save(signedSendShipment)
	if err != nil {
		return "", "", "", err
	}
	p2pMsg := messages.CreateSendShipmentRequestMessage(signedSendShipment)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return "", "", "", err
	}
	ss.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, signedSendShipment.Destination)
	return sendShipmentId, sendShipmentHash, sendShipmentSignature, nil
}

func NewSendShipmentService(
	ssr sendShipmentRepo.SendShipmentRepository,
	sss *sendShipmentService.SendShipmentService,
	p2pClient common.Messenger) *SendShipmentService {
	return &SendShipmentService{ssr: ssr, sss: sss, p2pClient: p2pClient}
}
