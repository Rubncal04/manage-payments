package handlers

import (
	"fmt"
	"net/http"

	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/repository"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

type UserRequest struct {
	Name      *string `json:"name"`
	CellPhone *string `json:"cellphone"`
	Paid      *bool   `json:"paid"`
	DateToPay *string `json:"date_to_pay"`
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.Repo.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	newUser, err := h.Repo.CreateUser(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID := c.Param("id")
	updateData := new(UserRequest)
	if err := c.Bind(updateData); err != nil {
		fmt.Println("Error binding data:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updateMap := bson.M{}

	if updateData.Name != nil {
		updateMap["name"] = *updateData.Name
	}
	if updateData.CellPhone != nil {
		updateMap["cellphone"] = *updateData.CellPhone
	}
	if updateData.Paid != nil {
		updateMap["paid"] = *updateData.Paid
	}
	if updateData.DateToPay != nil {
		updateMap["date_to_pay"] = *updateData.DateToPay
	}

	fmt.Println("Update map:", updateMap)

	if len(updateMap) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields provided for update"})
	}

	err := h.Repo.UpdateUser(userID, updateMap)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}
