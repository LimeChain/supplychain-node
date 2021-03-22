package repository

import "github.com/Limechain/pwc-bat-node/app/domain/product/model"

type ProductRepository interface {
	GetAll() ([]*model.Product, error)
	GetByID(id int) (*model.Product, error)
	Save(*model.Product) (id int, err error)
}
