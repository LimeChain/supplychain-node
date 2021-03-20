package sendShipment

import (
	"context"
	"errors"

	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendShipmentRepository struct {
	db *mongo.Database
}

func (r *SendShipmentRepository) GetAll() ([]*model.SendShipment, error) {
	collection := r.db.Collection("sent-shipments")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*model.SendShipment, 0)

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem model.SendShipment
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (r *SendShipmentRepository) GetByID(id int) (*model.SendShipment, error) {

	var result model.SendShipment
	collection := r.db.Collection("sent-shipments")
	if err := collection.FindOne(context.TODO(), bson.M{"unsignedSendShipment.obj.shipmentModel.shipmentId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SendShipmentRepository) Save(sendShipment *model.SendShipment) (int, error) {
	collection := r.db.Collection("sent-shipments")
	if sendShipment.Obj.ShipmentModel.ShipmentId == 0 {
		sendShipment.Obj.ShipmentModel.ShipmentId = 1
	}
	_, err := collection.InsertOne(context.TODO(), sendShipment)
	if err != nil {
		return 0, err
	}

	return sendShipment.Obj.ShipmentModel.ShipmentId, nil
}

func (r *SendShipmentRepository) Update(sendShipment *model.SendShipment) error {
	collection := r.db.Collection("sent-shipments")
	if sendShipment.Obj.ShipmentModel.ShipmentId == 0 {
		return errors.New("Shipment sent without Id cannot be updated")
	}
	ur, err := collection.ReplaceOne(context.TODO(), bson.M{"unsignedSendShipment.obj.shipmentModel.shipmentId": sendShipment.Obj.ShipmentModel.ShipmentId}, sendShipment)
	if err != nil {
		return err
	}

	if ur.MatchedCount == 0 {
		return errors.New("No such contract found")
	}

	return nil
}

func NewSendShipmentRepository(db *mongo.Database) *SendShipmentRepository {
	return &SendShipmentRepository{db: db}
}
