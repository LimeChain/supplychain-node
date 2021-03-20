package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/router"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	productRepository "github.com/Limechain/HCS-Integration-Node/app/domain/product/repository"
	sendShipmentRepository "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/repository"
	sendShipmentService "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	apiRouter "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common/queue"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/dlt/hcs"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
	productMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/product"
	sendShipmentMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/send-shipment"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func setupP2PClient(
	prvKey ed25519.PrivateKey,
	hcsClient common.DLTMessenger,
	productRepo productRepository.ProductRepository,
	sendShipmentRepo sendShipmentRepository.SendShipmentRepository,
	sendShipmentService *sendShipmentService.SendShipmentService) common.Messenger {

	listenPort := os.Getenv("P2P_PORT")
	listenIp := os.Getenv("P2P_IP")

	p2pClient := libp2p.NewLibP2PClient(prvKey, listenIp, listenPort)

	// TODO get some env variables
	// TODO add more handlers
	productHandler := handler.NewProductHandler(productRepo)
	sendShipmentRequestHandler := handler.NewSendShipmentRequestHandler(sendShipmentRepo, sendShipmentService, p2pClient)
	sendShipmentAcceptedHandler := handler.NewSendShipmentAcceptedHandler(sendShipmentRepo, sendShipmentService, hcsClient)

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	r.AddHandler(messages.P2PMessageTypeProduct, productHandler)
	r.AddHandler(messages.P2PMessageTypeSendShipmentRequest, sendShipmentRequestHandler)
	r.AddHandler(messages.P2PMessageTypeSendShipmentAccepted, sendShipmentAcceptedHandler)

	p2pChannel := make(chan *common.Message)

	p2pQueue := queue.New(p2pChannel, r)

	p2pClient.Listen(p2pQueue)

	return p2pClient
}

func setupDLTClient(
	prvKey ed25519.PrivateKey,
	sendShipmentRepo sendShipmentRepository.SendShipmentRepository,
	sendShipmentService *sendShipmentService.SendShipmentService) common.DLTMessenger {

	shouldConnectToMainnet := (os.Getenv("HCS_MAINNET") == "true")
	hcsClientID := os.Getenv("HCS_CLIENT_ID")
	hcsMirrorNodeID := os.Getenv("HCS_MIRROR_NODE_ADDRESS")
	topicID := os.Getenv("HCS_TOPIC_ID")

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	sendShipmentHandler := handler.NewDLTSendShipmentHandler(sendShipmentRepo, sendShipmentService)

	r.AddHandler(messages.DLTMessageTypeSendShipment, sendShipmentHandler)

	ch := make(chan *common.Message)

	q := queue.New(ch, r)

	hcsClient := hcs.NewHCSClient(hcsClientID, prvKey, hcsMirrorNodeID, topicID, shouldConnectToMainnet)

	err := hcsClient.Listen(q)
	if err != nil {
		panic(err)
	}

	return hcsClient

}

func main() {

	args := os.Args[1:]
	if len(args) > 0 {
		godotenv.Load(args[0])
	} else {
		godotenv.Load()
	}

	logFilePath := os.Getenv("LOG_FILE")

	setupLogger()

	if len(logFilePath) > 0 {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Panic(err)
		}

		defer file.Close()

		setupFileLogger(file)
	}

	prvKey := getPrivateKey()

	pubKey, ok := prvKey.Public().(ed25519.PublicKey)
	if !ok {
		log.Errorf("Could not cast the public key to ed25519 public key")
	}

	log.Infof("Started node with public key: %s\n", hex.EncodeToString(pubKey))

	peerPubKey := getPeerPublicKey()

	client, db := connectToDb()

	defer client.Disconnect(context.Background())

	productRepo := productMongo.NewProductRepository(db)
	sendShipmentRepo := sendShipmentMongo.NewSendShipmentRepository(db)

	sss := sendShipmentService.New(prvKey, peerPubKey)

	hcsClient := setupDLTClient(prvKey, sendShipmentRepo, sss)

	defer hcsClient.Close()

	p2pClient := setupP2PClient(prvKey, hcsClient, productRepo, sendShipmentRepo, sss)

	defer p2pClient.Close()

	apiPort := os.Getenv("API_PORT")

	a := api.NewIntegrationNodeAPI()

	productApiService := apiservices.NewProductService(productRepo, p2pClient)

	sendShipmentApiService := apiservices.NewSendShipmentService(sendShipmentRepo, sss, p2pClient)

	nodeApiService := apiservices.NewNodeService(p2pClient)

	webPlatformApiService := apiservices.NewWebPlatformService()

	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RouteProduct), apiRouter.NewProductRouter(productApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RouteWebPlatform), apiRouter.NewWebPlatformRouter(webPlatformApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.Swagger), apiRouter.NewSwaggerRouter())
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.Node), apiRouter.NewNodeRouter(nodeApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.SendShipment), apiRouter.NewSendShipmentRouter(sendShipmentApiService))

	if err := a.Start(apiPort); err != nil {
		panic(err)
	}

}
