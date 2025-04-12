package handlers

import (
	"net/http"

	"github/Rubncal04/youtube-premium/auth"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

func Register(c echo.Context, mongoRepo *db.MongoRepo) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Check if user already exists
	filter := bson.M{"$or": []bson.M{
		{"email": req.Email},
		{"username": req.Username},
	}}
	existingUser := &models.User{}
	_, err := mongoRepo.FindOne("users", filter, existingUser)
	if err == nil {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	// Create new user
	user, err := models.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	// Save user to database
	_, err = mongoRepo.Create("users", user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save user")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "user created successfully",
	})
}

func Login(c echo.Context, mongoRepo *db.MongoRepo, secretKey string) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Find user by email
	filter := bson.M{"email": req.Email}
	user := &models.User{}
	_, err := mongoRepo.FindOne("users", filter, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Generate tokens
	tokenPair, err := auth.GenerateTokenPair(user.ID.Hex(), secretKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate tokens")
	}

	// Prepare response
	response := AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func RefreshToken(c echo.Context, secretKey string) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate refresh token
	claims, err := auth.ValidateToken(req.RefreshToken, secretKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	// Generate new token pair
	tokenPair, err := auth.GenerateTokenPair(claims.UserID, secretKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate tokens")
	}

	return c.JSON(http.StatusOK, tokenPair)
}
