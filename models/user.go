package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// En models/user.go
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name"`
	CellPhone       string             `bson:"cell_phone" json:"cell_phone"`
	DateToPay       string             `bson:"date_to_pay" json:"date_to_pay"`
	Paid            bool               `bson:"paid" json:"paid"`
	Status          string             `bson:"status" json:"status"` // 'active', 'inactive'
	LastPaymentDate time.Time          `bson:"last_payment_date" json:"last_payment_date"`
}
