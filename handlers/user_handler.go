package handlers

import (
	"github/Rubncal04/youtube-premium/repository"
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
