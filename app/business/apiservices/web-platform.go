package apiservices

type WebPlatformService struct {
}

func (ss *WebPlatformService) CreditShipment(data []byte) (hash string) {
	return "0x111111"
}

func (ss *WebPlatformService) CreditProduct(data []byte) {
}

func NewWebPlatformService() *WebPlatformService {
	return &WebPlatformService{}
}
