package handlers

import (
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientHandler struct {
	clientRepo *repository.ClientRepository
}

type ClientRequest struct {
	Name      string `json:"name"`
	CellPhone string `json:"cell_phone"`
	DayToPay  int    `json:"day_to_pay"`
}

func NewClientHandler(clientRepo *repository.ClientRepository) *ClientHandler {
	return &ClientHandler{
		clientRepo: clientRepo,
	}
}

// CreateClient handles the creation of a new client
func (h *ClientHandler) CreateClient(c echo.Context) error {
	var clientRequest ClientRequest
	if err := c.Bind(&clientRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	client := models.NewClient(userID, clientRequest.Name, clientRequest.CellPhone, clientRequest.DayToPay)

	// Save client to database
	newClient, err := h.clientRepo.Create(*client)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create client"})
	}

	return c.JSON(http.StatusCreated, newClient)
}

// GetClients handles getting all clients for the authenticated user
func (h *ClientHandler) GetClients(c echo.Context) error {
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	clients, err := h.clientRepo.GetAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get clients"})
	}

	return c.JSON(http.StatusOK, clients)
}

// GetClient handles getting a specific client by ID
func (h *ClientHandler) GetClient(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	client, err := h.clientRepo.GetByID(id.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	// Verify that the client belongs to the authenticated user
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok || client.UserID != userID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	return c.JSON(http.StatusOK, client)
}

// UpdateClient handles updating a client
func (h *ClientHandler) UpdateClient(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	// Get existing client
	client, err := h.clientRepo.GetByID(id.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	// Verify ownership
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok || client.UserID != userID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Bind the update request to a new struct
	var updateRequest ClientRequest
	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Update client fields
	updateData := bson.M{
		"name":       updateRequest.Name,
		"cell_phone": updateRequest.CellPhone,
		"day_to_pay": updateRequest.DayToPay,
		"updated_at": time.Now(),
	}

	if err := h.clientRepo.Update(id.Hex(), updateData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update client"})
	}

	// Get the updated client
	updatedClient, err := h.clientRepo.GetByID(id.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get updated client"})
	}

	return c.JSON(http.StatusOK, updatedClient)
}

func (h *ClientHandler) DeleteClient(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid client ID"})
	}

	client, err := h.clientRepo.GetByID(id.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok || client.UserID != userID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	err = h.clientRepo.Delete(id.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete client"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Client deleted successfully"})
}
