package handlers

import (
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	Repo *repository.PaymentRepository
}

func NewPaymentHandler(repo *repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{Repo: repo}
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	payment := new(models.Payment)
	userID := c.Param("userId")

	if err := c.Bind(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid userID"})
	}
	payment.UserID = objectID
	payment.PaymentDate = time.Now()
	payment.Status = "completed"

	newPayment, err := h.Repo.CreatePayment(payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newPayment)
}

func (h *PaymentHandler) GetPaymentsByUser(c echo.Context) error {
	userID := c.Param("userId")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid userID"})
	}

	payments, err := h.Repo.GetPaymentsByUser(objectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, payments)
}
