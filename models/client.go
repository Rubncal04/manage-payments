package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name            string             `bson:"name" json:"name"`
	CellPhone       string             `bson:"cell_phone" json:"cell_phone"`
	DayToPay        int                `bson:"day_to_pay" json:"day_to_pay"`
	Status          string             `bson:"status" json:"status"` // 'active', 'inactive'
	LastPaymentDate time.Time          `bson:"last_payment_date" json:"last_payment_date"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewClient creates a new client with the provided information
func NewClient(userID primitive.ObjectID, name, cellPhone string, dayToPay int) *Client {
	return &Client{
		UserID:          userID,
		Name:            name,
		CellPhone:       cellPhone,
		DayToPay:        dayToPay,
		Status:          "inactive",
		LastPaymentDate: time.Time{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
