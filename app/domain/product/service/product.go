package service

import (
	"github.com/Limechain/pwc-bat-node/app/domain/product/repository"
)

type ProductService struct {
	r repository.ProductRepository
}

func (s *ProductService) Exists(ID int) (bool, error) {
	product, err := s.r.GetByID(ID)
	if err != nil {
		return false, err
	}

	return (product != nil), nil

}

func New(repo repository.ProductRepository) *ProductService {
	return &ProductService{r: repo}
}
