package handlers

import (
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	paymentRepo *repository.PaymentRepository
	clientRepo  *repository.ClientRepository
}

func NewPaymentHandler(paymentRepo *repository.PaymentRepository, clientRepo *repository.ClientRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo: paymentRepo,
		clientRepo:  clientRepo,
	}
}

// GetAllPayments handles getting all payments
func (h *PaymentHandler) GetAllPayments(c echo.Context) error {
	payments, err := h.paymentRepo.GetAllPayments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get payments"})
	}
	return c.JSON(http.StatusOK, payments)
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

	var payment models.Payment
	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	payment.ClientID = clientID
	newPayment, err := h.paymentRepo.CreatePayment(&payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create payment"})
	}
	// Update client's last payment date
	if err := h.clientRepo.UpdateLastPaymentDate(clientID, primitive.NewDateTimeFromTime(newPayment.PaymentDate)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update client's last payment date"})
	}

	return c.JSON(http.StatusCreated, newPayment)
}
