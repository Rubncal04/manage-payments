package handlers

import (
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"
	"net/http"

	"fmt"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	paymentRepo *repository.PaymentRepository
	clientRepo  *repository.ClientRepository
}

type PaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

func NewPaymentHandler(paymentRepo *repository.PaymentRepository, clientRepo *repository.ClientRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo: paymentRepo,
		clientRepo:  clientRepo,
	}
}

// GetAllPayments handles getting all payments for the authenticated user
func (h *PaymentHandler) GetAllPayments(c echo.Context) error {
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get all clients for the user
	clients, err := h.clientRepo.GetAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user's clients"})
	}

	// Get all payments for these clients
	var allPayments []models.Payment
	for _, client := range clients {
		payments, err := h.paymentRepo.GetPaymentsByClientID(client.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get payments"})
		}
		allPayments = append(allPayments, payments...)
	}

	return c.JSON(http.StatusOK, allPayments)
}

// GetOnePayment handles getting a single payment by ID
func (h *PaymentHandler) GetOnePayment(c echo.Context) error {

	clientID, err := primitive.ObjectIDFromHex(c.Param("clientId"))
	if err != nil {
		fmt.Printf("GetOnePayment: Error on convert clientId: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	paymentID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		fmt.Printf("GetOnePayment: Error on convert paymentId: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment ID"})
	}

	// Verify that the client belongs to the authenticated user
	client, err := h.clientRepo.GetByID(clientID.Hex())
	if err != nil {
		fmt.Printf("GetOnePayment: Error on get client: %v\n", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok {
		fmt.Println("GetOnePayment: Error on get user_id from context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if client.UserID != userID {
		fmt.Printf("GetOnePayment: Error on verify user_id: %v with client_id: %v\n", userID, client.UserID)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get the payment
	payment, err := h.paymentRepo.GetByID(paymentID)
	if err != nil {
		fmt.Printf("GetOnePayment: Error on get payment: %v\n", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Payment not found"})
	}

	// Verify that the payment belongs to the client
	if payment.ClientID != clientID {
		fmt.Printf("GetOnePayment: Error on verify payment_id: %v with client_id: %v\n", paymentID, clientID)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Payment not found for this client"})
	}

	return c.JSON(http.StatusOK, payment)
}

// GetPaymentsByClient handles getting all payments for a specific client
func (h *PaymentHandler) GetPaymentsByClient(c echo.Context) error {
	clientID, err := primitive.ObjectIDFromHex(c.Param("clientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	// Verify that the client belongs to the authenticated user
	client, err := h.clientRepo.GetByID(clientID.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok || client.UserID != userID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	payments, err := h.paymentRepo.GetPaymentsByClientID(clientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get payments"})
	}

	return c.JSON(http.StatusOK, payments)
}

// CreatePayment handles creating a new payment
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	clientID, err := primitive.ObjectIDFromHex(c.Param("clientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	// Verify that the client belongs to the authenticated user
	client, err := h.clientRepo.GetByID(clientID.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok || client.UserID != userID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var paymentRequest PaymentRequest
	if err := c.Bind(&paymentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Create new payment in processing state
	payment := models.NewPayment(clientID, paymentRequest.Amount)

	// Save payment in processing state
	if err := h.paymentRepo.CreatePayment(payment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create payment"})
	}

	// Try to process the payment
	// This is where you would integrate with your payment processor
	// For now, we'll simulate a successful payment
	if err := h.processPayment(payment); err != nil {
		// If payment processing fails, update status to rejected
		if updateErr := h.paymentRepo.RejectPayment(payment.ID, err.Error()); updateErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update payment status"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Payment processing failed"})
	}

	// Update client's last payment date
	if err := h.clientRepo.UpdateLastPaymentDate(clientID, primitive.NewDateTimeFromTime(payment.PaymentDate)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update client's last payment date"})
	}

	return c.JSON(http.StatusCreated, payment)
}

// processPayment simulates payment processing
// In a real application, this would integrate with a payment processor
func (h *PaymentHandler) processPayment(payment *models.Payment) error {
	// Simulate payment processing
	// In a real application, this would call your payment processor API
	// and handle the response

	// For demonstration, we'll just complete the payment
	return h.paymentRepo.CompletePayment(payment.ID)
}
