package repository

import (
	"fmt"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientRepository struct {
	Mongo *db.MongoRepo
}

func NewClientRepository(mongo *db.MongoRepo) *ClientRepository {
	return &ClientRepository{Mongo: mongo}
}

func (r *ClientRepository) Create(client models.Client) (models.Client, error) {
	coll, err := r.Mongo.Create("clients", client)

	if err != nil {
		return models.Client{}, err
	}

	newClient := models.Client{
		ID:              coll.InsertedID.(primitive.ObjectID),
		UserID:          client.UserID,
		Name:            client.Name,
		CellPhone:       client.CellPhone,
		DayToPay:        client.DayToPay,
		Status:          client.Status,
		LastPaymentDate: client.LastPaymentDate,
	}

	return newClient, nil
}

func (r *ClientRepository) GetByID(id string) (models.Client, error) {
	var client models.Client
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Client{}, fmt.Errorf("invalid client ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	_, err = r.Mongo.FindOne("clients", filter, &client)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

func (r *ClientRepository) Update(id string, client bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	update := bson.M{
		"$set": client,
	}

	return r.Mongo.UpdateOne("clients", bson.M{"_id": objID}, update)
}

func (r *ClientRepository) GetAll(userID primitive.ObjectID) ([]models.Client, error) {
	var clients []models.Client
	filter := bson.M{"user_id": userID}
	err := r.Mongo.FindAll("clients", filter, &clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) UpdateLastPaymentDate(clientID primitive.ObjectID, paymentDate primitive.DateTime) error {
	return r.Mongo.UpdateOne("clients", bson.M{"_id": clientID}, bson.M{"$set": bson.M{"last_payment_date": paymentDate}})
}

func (r *ClientRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	return r.Mongo.DeleteOne("clients", filter)
}
