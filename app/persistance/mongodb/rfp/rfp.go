package rfp

import (
	"context"

	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RFPRepository struct {
	db *mongo.Database
}

func (r *RFPRepository) GetAll() ([]*model.Product, error) {
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

func (r *RFPRepository) GetByID(id string) (*model.Product, error) {

	var result model.Product
	collection := r.db.Collection("rfps")
	if err := collection.FindOne(context.TODO(), bson.M{"rfpId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
func (r *RFPRepository) Save(rfp *model.Product) (int, error) {
	collection := r.db.Collection("products")
	if rfp.Obj.ProductModel.ProductId == 0 {
		rfp.Obj.ProductModel.ProductId = 1
	}
	_, err := collection.InsertOne(context.TODO(), rfp)
	if err != nil {
		return 0, err
	}

	return rfp.Obj.ProductModel.ProductId, nil
}

func NewRFPRepository(db *mongo.Database) *RFPRepository {
	return &RFPRepository{db: db}
}
