package sendShipment

import (
	"context"
	"errors"

	"github.com/Limechain/HCS-Integration-Node/app/domain/send-shipment/model"
	"github.com/google/uuid"
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

func (r *SendShipmentRepository) GetByID(id string) (*model.SendShipment, error) {

	var result model.SendShipment
	collection := r.db.Collection("sent-shipments")
	if err := collection.FindOne(context.TODO(), bson.M{"unsignedSendShipment.shipmentId": id}).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SendShipmentRepository) Save(sendShipment *model.SendShipment) (string, error) {
	collection := r.db.Collection("sent-shipments")
	if len(sendShipment.ShipmentId) == 0 {
		sendShipment.ShipmentId = uuid.New().String()
	}
	_, err := collection.InsertOne(context.TODO(), sendShipment)
	if err != nil {
		return "", err
	}

	return sendShipment.ShipmentId, nil
}

func (r *SendShipmentRepository) Update(sendShipment *model.SendShipment) error {
	collection := r.db.Collection("sent-shipments")
	if len(sendShipment.ShipmentId) == 0 {
		return errors.New("Shipment sent without Id cannot be updated")
	}
	ur, err := collection.ReplaceOne(context.TODO(), bson.M{"unsignedSendShipment.shipmentId": sendShipment.ShipmentId}, sendShipment)
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
