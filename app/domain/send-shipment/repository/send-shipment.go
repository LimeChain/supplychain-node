package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"

type SendShipmentRepository interface {
	GetAll() ([]*model.SendShipment, error)
	GetByID(id string) (*model.SendShipment, error)
	Save(*model.SendShipment) (id string, err error)
	Update(*model.SendShipment) error
}
