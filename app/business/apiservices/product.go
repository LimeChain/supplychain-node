package apiservices

import (
	"context"
	"encoding/json"

	"github.com/Limechain/pwc-bat-node/app/business/messages"
	productModel "github.com/Limechain/pwc-bat-node/app/domain/product/model"
	productRepository "github.com/Limechain/pwc-bat-node/app/domain/product/repository"
	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
)

type ProductService struct {
	repo      productRepository.ProductRepository
	p2pClient common.Messenger
}

func (s *ProductService) GetAllProducts() ([]*productModel.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProduct(productId int) (*productModel.Product, error) {
	return s.repo.GetByID(productId)
}

func (s *ProductService) CreateProduct(product *productModel.Product) (id int, err error) {
	productId, err := s.repo.Save(product)
	if err != nil {
		return 0, err
	}
	p2pMsg := messages.CreateProductMessage(product)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return 0, err
	}
	s.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, product.Destination)
	return productId, nil
}

func NewProductService(repo productRepository.ProductRepository, p2pClient common.Messenger) *ProductService {
	return &ProductService{repo: repo, p2pClient: p2pClient}
}
