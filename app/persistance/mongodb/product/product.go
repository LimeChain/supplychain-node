package product

import (
	"context"

	"github.com/Limechain/pwc-bat-node/app/domain/product/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	db *mongo.Database
}

func (r *ProductRepository) GetAll() ([]*model.Product, error) {
	collection := r.db.Collection("products")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*model.Product, 0)

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.Product
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var result model.Product
	collection := r.db.Collection("products")
	if err := collection.FindOne(context.TODO(), bson.M{"obj.productModel.productId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
func (r *ProductRepository) Save(product *model.Product) (int, error) {
	collection := r.db.Collection("products")
	if product.Obj.ProductModel.ProductId == 0 {
		product.Obj.ProductModel.ProductId = 1
	}
	_, err := collection.InsertOne(context.TODO(), product)
	if err != nil {
		return 0, err
	}

	return product.Obj.ProductModel.ProductId, nil
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{db: db}
}
