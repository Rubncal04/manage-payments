package repository

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository struct {
	Mongo *db.MongoRepo
}

func NewPaymentRepository(mongo *db.MongoRepo) *PaymentRepository {
	return &PaymentRepository{Mongo: mongo}
}

func (r *PaymentRepository) CreatePayment(payment *models.Payment) (models.Payment, error) {
	// Crear el pago
	coll, err := r.Mongo.Create("payments", payment)
	if err != nil {
		return models.Payment{}, err
	}

	// Actualizar el estado del usuario
	update := bson.M{
		"$set": bson.M{
			"paid":              true,
			"status":            "active",
			"last_payment_date": payment.PaymentDate,
		},
	}

	err = r.Mongo.UpdateOne(
		"users",
		bson.M{"_id": payment.UserID},
		update,
	)

	if err != nil {
		return models.Payment{}, err
	}

	result := models.Payment{
		ID:          coll.InsertedID.(primitive.ObjectID),
		Amount:      payment.Amount,
		PaymentDate: payment.PaymentDate,
		UserID:      payment.UserID,
		Status:      "completed",
	}

	return result, nil
}

func (r *PaymentRepository) GetPaymentsByUser(userID primitive.ObjectID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.Mongo.FindAll("payments", bson.M{"user_id": userID}, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) GetAllPayments() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.Mongo.FindAll("payments", bson.M{}, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
