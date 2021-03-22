package repository

import "github.com/Limechain/pwc-bat-node/app/domain/send-shipment/model"

type SendShipmentRepository interface {
	GetAll() ([]*model.SendShipment, error)
	GetByID(id int) (*model.SendShipment, error)
	GetByHash(hash string) (*model.SendShipment, error)
	GetByDLTMessage(dltMessage string) (*model.SendShipment, error)
	Save(*model.SendShipment) (id int, err error)
	Update(*model.SendShipment) error
}
