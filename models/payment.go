package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Amount      float64            `bson:"amount" json:"amount"`
	PaymentDate time.Time          `bson:"payment_date" json:"payment_date"`
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`
	Status      string             `bson:"status" json:"status"` // 'completed', 'failed'
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewPayment creates a new payment with the provided information
func NewPayment(clientID primitive.ObjectID, amount float64, paymentDate time.Time) *Payment {
	return &Payment{
		ClientID:    clientID,
		Amount:      amount,
		PaymentDate: paymentDate,
		Status:      "completed",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
