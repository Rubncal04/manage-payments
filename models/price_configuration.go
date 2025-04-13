package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceConfiguration struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewPriceConfiguration(userID primitive.ObjectID, amount float64) *PriceConfiguration {
	now := time.Now()
	return &PriceConfiguration{
		UserID:    userID,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
