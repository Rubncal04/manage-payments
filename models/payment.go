package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Amount      float64            `bson:"amount" json:"amount"`
	PaymentDate time.Time          `bson:"payment_date" json:"payment_date"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status      string             `bson:"status" json:"status"` // 'completed', 'failed'
}
