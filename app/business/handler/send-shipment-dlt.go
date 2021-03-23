package handler

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

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

type NodeJsDltRequest struct {
	ShipmentId int    `json:"shipmentId"`
	Dlt        string `json:"dlt"`
}

type NodeJsRequestWrapper struct {
	Ac string `json:"ac"`
	Pl string `json:"pl"`
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

	sn := msg.Ctx.Value(hcs.DLTValuesKey).(hcs.DLTValues).Get(hcs.SequenceNumberKey)
	txId := msg.Ctx.Value(hcs.DLTValuesKey).(hcs.DLTValues).Get(hcs.TransactionIdKey)

	savedSendShipment.DLTAnchored = true
	savedSendShipment.DLTProof = sn
	savedSendShipment.DLTTransactionId = txId
	savedSendShipment.DLTMessage = hex.EncodeToString(msg.Msg)
	log.Println("RECEIVED")
	log.Println(savedSendShipment.Obj.ShipmentModel.ShipmentId)
	log.Println(savedSendShipment.DLTTransactionId)
	log.Println(savedSendShipment.DLTMessage)

	err = h.sendShipmentRepo.Update(savedSendShipment)
	if err != nil {
		return err
	}

	// values := map[string]string{"ac": "c", "pl": {"shipmentId": savedSendShipment.Obj.ShipmentModel.ShipmentId, "dlt": savedSendShipment.DLTTransactionId}}
	nodeJsRequest := NodeJsDltRequest{
		ShipmentId: savedSendShipment.Obj.ShipmentModel.ShipmentId,
		Dlt:        savedSendShipment.DLTTransactionId,
	}
	nodeJsRequestString, err := json.Marshal(nodeJsRequest)
	if err != nil {
		log.Fatal(err)
	}

	values := NodeJsRequestWrapper{
		Ac: "c",
		Pl: string(nodeJsRequestString),
	}
	json_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.Post(os.Getenv("NODEJS_SERVER"), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Println("Error:")
		log.Fatal(err)
	}

	log.Infof("Sent shipment with Id: %d seen in the dlt and verified\n", savedSendShipment.Obj.ShipmentModel.ShipmentId)
	return nil
}

func NewDLTSendShipmentHandler(sendShipmentRepo repository.SendShipmentRepository, sendShipmentService *service.SendShipmentService) *DLTSendShipmentHandler {
	return &DLTSendShipmentHandler{sendShipmentRepo: sendShipmentRepo, sendShipmentService: sendShipmentService}
}
