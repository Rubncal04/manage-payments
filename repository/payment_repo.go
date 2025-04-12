package repository

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository struct {
	Mongo *db.MongoRepo
}

func NewPaymentRepository(mongo *db.MongoRepo) *PaymentRepository {
	return &PaymentRepository{Mongo: mongo}
}

// CreatePayment creates a new payment in processing state
func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	payment.Status = models.PaymentStatusProcessing

	result, err := r.Mongo.Create("payments", payment)
	if err != nil {
		return err
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// CompletePayment updates a payment to completed state
func (r *PaymentRepository) CompletePayment(paymentID primitive.ObjectID) error {
	filter := bson.M{"_id": paymentID}
	update := bson.M{
		"$set": bson.M{
			"status":     models.PaymentStatusCompleted,
			"updated_at": time.Now(),
		},
	}

	return r.Mongo.UpdateOne("payments", filter, update)
}

// RejectPayment updates a payment to rejected state with an error message
func (r *PaymentRepository) RejectPayment(paymentID primitive.ObjectID, errorMsg string) error {
	filter := bson.M{"_id": paymentID}
	update := bson.M{
		"$set": bson.M{
			"status":     models.PaymentStatusRejected,
			"error":      errorMsg,
			"updated_at": time.Now(),
		},
	}

	return r.Mongo.UpdateOne("payments", filter, update)
}

// GetPaymentsByClientID retrieves all payments for a specific client
func (r *PaymentRepository) GetPaymentsByClientID(clientID primitive.ObjectID) ([]models.Payment, error) {
	filter := bson.M{"client_id": clientID}
	var payments []models.Payment
	err := r.Mongo.FindAll("payments", filter, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// GetAllPayments retrieves all payments
func (r *PaymentRepository) GetAllPayments() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.Mongo.FindAll("payments", bson.M{}, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// GetByID retrieves a payment by its ID
func (r *PaymentRepository) GetByID(id primitive.ObjectID) (*models.Payment, error) {
	filter := bson.M{"_id": id}
	var payment models.Payment
	_, err := r.Mongo.FindOne("payments", filter, &payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
