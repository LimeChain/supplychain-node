package apiservices

import (
	"context"
	"encoding/json"

	"github.com/Limechain/pwc-bat-node/app/business/messages"
	sendShipmentModel "github.com/Limechain/pwc-bat-node/app/domain/send-shipment/model"
	sendShipmentRepo "github.com/Limechain/pwc-bat-node/app/domain/send-shipment/repository"
	sendShipmentService "github.com/Limechain/pwc-bat-node/app/domain/send-shipment/service"
	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
)

type SendShipmentService struct {
	ssr       sendShipmentRepo.SendShipmentRepository
	sss       *sendShipmentService.SendShipmentService
	p2pClient common.Messenger
}

func (ss *SendShipmentService) GetAllSentShipments() ([]*sendShipmentModel.SendShipment, error) {
	return ss.ssr.GetAll()
}

func (ss *SendShipmentService) GetSentShipment(shipmentId int) (*sendShipmentModel.SendShipment, error) {
	return ss.ssr.GetByID(shipmentId)
}

func (ss *SendShipmentService) GetSentShipmentByDLTMessage(dltMessage string) (*sendShipmentModel.SendShipment, error) {
	return ss.ssr.GetByDLTMessage(dltMessage)
}

func (ss *SendShipmentService) SaveAndSendSendShipment(unsignedSendShipment *sendShipmentModel.UnsignedSendShipment) (sendShipmentId int, sendShipmentHash, sendShipmentSignature string, err error) {
	sendShipmentHash, err = ss.sss.Hash(unsignedSendShipment)
	if err != nil {
		return 0, "", "", err
	}
	sendShipmentSignature, err = ss.sss.Sign(unsignedSendShipment)
	if err != nil {
		return 0, "", "", err
	}
	signedSendShipment := &sendShipmentModel.SendShipment{UnsignedSendShipment: *unsignedSendShipment, BuyerSignature: sendShipmentSignature, DLTAnchored: false}
	sendShipmentId, err = ss.ssr.Save(signedSendShipment)
	if err != nil {
		return 0, "", "", err
	}
	p2pMsg := messages.CreateSendShipmentRequestMessage(signedSendShipment)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return 0, "", "", err
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
