package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"

type RFPRepository interface {
	GetAll() ([]*model.Product, error)
	GetByID(id string) (*model.Product, error)
	Save(*model.Product) (id int, err error)
}
