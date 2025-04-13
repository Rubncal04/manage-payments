package handlers

import (
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceConfigurationHandler struct {
	priceConfigRepo *repository.PriceConfigurationRepository
}

func NewPriceConfigurationHandler(repo *repository.PriceConfigurationRepository) *PriceConfigurationHandler {
	return &PriceConfigurationHandler{priceConfigRepo: repo}
}

type PriceConfigRequest struct {
	Amount float64 `json:"amount"`
}

func (h *PriceConfigurationHandler) CreatePriceConfig(c echo.Context) error {
	userID := c.Get("user_id").(primitive.ObjectID)

	var request PriceConfigRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if request.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount must be greater than 0"})
	}

	config := models.NewPriceConfiguration(userID, request.Amount)

	if err := h.priceConfigRepo.Create(config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, config)
}

func (h *PriceConfigurationHandler) GetPriceConfig(c echo.Context) error {
	userID := c.Get("user_id").(primitive.ObjectID)

	config, err := h.priceConfigRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Price configuration not found"})
	}

	return c.JSON(http.StatusOK, config)
}

func (h *PriceConfigurationHandler) UpdatePriceConfig(c echo.Context) error {
	userID := c.Get("user_id").(primitive.ObjectID)

	var request PriceConfigRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if request.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Amount must be greater than 0"})
	}

	// Verify if the configuration exists
	_, err := h.priceConfigRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Price configuration not found"})
	}

	if err := h.priceConfigRepo.Update(userID, request.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Price configuration updated successfully"})
}

func (h *PriceConfigurationHandler) DeletePriceConfig(c echo.Context) error {
	userID := c.Get("user_id").(primitive.ObjectID)

	// Verify if the configuration exists
	_, err := h.priceConfigRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Price configuration not found"})
	}

	if err := h.priceConfigRepo.Delete(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Price configuration deleted successfully"})
}
