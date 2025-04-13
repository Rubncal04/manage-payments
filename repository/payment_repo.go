package repository

import (
	"context"
	"fmt"
	"github/Rubncal04/youtube-premium/cache"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository struct {
	Mongo *db.MongoRepo
	cache cache.Cache
}

func NewPaymentRepository(mongo *db.MongoRepo, cache cache.Cache) *PaymentRepository {
	return &PaymentRepository{
		Mongo: mongo,
		cache: cache,
	}
}

// CreatePayment creates a new payment in processing state
func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	payment.Status = models.PaymentStatusProcessing

	result, err := r.Mongo.Create("payments", payment)
	if err != nil {
		return err
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)

	// Invalidar caché si está disponible
	if r.cache != nil {
		cache.InvalidateCache(context.Background(), r.cache,
			fmt.Sprintf("payment:%s", payment.ID.Hex()),
			fmt.Sprintf("payments:client:%s", payment.ClientID.Hex()),
			"payments:all")
	}

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

	err := r.Mongo.UpdateOne("payments", filter, update)
	if err != nil {
		return err
	}

	// Invalidar caché si está disponible
	if r.cache != nil {
		cache.InvalidateCache(context.Background(), r.cache,
			fmt.Sprintf("payment:%s", paymentID.Hex()))
	}

	return nil
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

	err := r.Mongo.UpdateOne("payments", filter, update)
	if err != nil {
		return err
	}

	// Invalidar caché si está disponible
	if r.cache != nil {
		cache.InvalidateCache(context.Background(), r.cache,
			fmt.Sprintf("payment:%s", paymentID.Hex()))
	}

	return nil
}

// GetPaymentsByClientID retrieves all payments for a specific client
func (r *PaymentRepository) GetPaymentsByClientID(clientID primitive.ObjectID) ([]models.Payment, error) {
	if r.cache != nil {
		key := fmt.Sprintf("payments:client:%s", clientID.Hex())
		var payments []models.Payment
		result, err := cache.WithCache(context.Background(), r.cache, key, payments, 1*time.Hour, func() ([]models.Payment, error) {
			filter := bson.M{"client_id": clientID}
			var dbPayments []models.Payment
			err := r.Mongo.FindAll("payments", filter, &dbPayments)
			return dbPayments, err
		})
		return result, err
	}

	// Si no hay caché, usar MongoDB directamente
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
	if r.cache != nil {
		key := "payments:all"
		var payments []models.Payment
		result, err := cache.WithCache(context.Background(), r.cache, key, payments, 1*time.Hour, func() ([]models.Payment, error) {
			var dbPayments []models.Payment
			err := r.Mongo.FindAll("payments", bson.M{}, &dbPayments)
			return dbPayments, err
		})
		return result, err
	}

	// Si no hay caché, usar MongoDB directamente
	var payments []models.Payment
	err := r.Mongo.FindAll("payments", bson.M{}, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// GetByID retrieves a payment by its ID
func (r *PaymentRepository) GetByID(id primitive.ObjectID) (*models.Payment, error) {
	if r.cache != nil {
		key := fmt.Sprintf("payment:%s", id.Hex())
		var payment models.Payment
		result, err := cache.WithCache(context.Background(), r.cache, key, payment, 1*time.Hour, func() (models.Payment, error) {
			filter := bson.M{"_id": id}
			var dbPayment models.Payment
			_, err := r.Mongo.FindOne("payments", filter, &dbPayment)
			return dbPayment, err
		})
		if err != nil {
			return nil, err
		}
		return &result, nil
	}

	// Si no hay caché, usar MongoDB directamente
	filter := bson.M{"_id": id}
	var payment models.Payment
	_, err := r.Mongo.FindOne("payments", filter, &payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
