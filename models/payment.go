package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentStatus represents the possible states of a payment
type PaymentStatus string

const (
	PaymentStatusProcessing PaymentStatus = "processing" // Initial state when payment is being processed
	PaymentStatusCompleted  PaymentStatus = "completed"  // Final state when payment is successful
	PaymentStatusRejected   PaymentStatus = "rejected"   // Final state when payment fails
)

// Payment represents a payment transaction
type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Amount      float64            `bson:"amount" json:"amount"`
	PaymentDate time.Time          `bson:"payment_date" json:"payment_date"`
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`
	Status      PaymentStatus      `bson:"status" json:"status"`
	Error       string             `bson:"error,omitempty" json:"error,omitempty"` // Stores error message if status is rejected
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewPayment creates a new payment with the provided information
func NewPayment(clientID primitive.ObjectID, amount float64) *Payment {
	return &Payment{
		ClientID:    clientID,
		Amount:      amount,
		PaymentDate: time.Now(),
		Status:      PaymentStatusProcessing, // Initial state
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ValidateStateTransition checks if a state transition is valid
func (p *Payment) ValidateStateTransition(newStatus PaymentStatus) error {
	switch p.Status {
	case PaymentStatusProcessing:
		// From processing, can only go to completed or rejected
		if newStatus != PaymentStatusCompleted && newStatus != PaymentStatusRejected {
			return errors.New("invalid state transition: processing can only transition to completed or rejected")
		}
	case PaymentStatusCompleted, PaymentStatusRejected:
		// Once in a final state, cannot change
		return errors.New("invalid state transition: cannot change from a final state")
	}
	return nil
}

// SetStatus updates the payment status with validation
func (p *Payment) SetStatus(newStatus PaymentStatus, errorMsg string) error {
	if err := p.ValidateStateTransition(newStatus); err != nil {
		return err
	}
	p.Status = newStatus
	p.Error = errorMsg
	p.UpdatedAt = time.Now()
	return nil
}
