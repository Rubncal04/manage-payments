package routes

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/handlers"
	"github/Rubncal04/youtube-premium/middleware"
	"github/Rubncal04/youtube-premium/repository"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes define las rutas principales de la aplicación.
func RegisterRoutes(e *echo.Echo, mongoRepo *db.MongoRepo, secretKey string) {
	// Public routes
	e.POST("/register", func(c echo.Context) error {
		return handlers.Register(c, mongoRepo)
	})
	e.POST("/login", func(c echo.Context) error {
		return handlers.Login(c, mongoRepo, secretKey)
	})
	e.POST("/refresh", func(c echo.Context) error {
		return handlers.RefreshToken(c, secretKey)
	})

	// Protected routes
	api := e.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(secretKey))

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(mongoRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	paymentRepo := repository.NewPaymentRepository(mongoRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentRepo)

	// User routes
	api.GET("/users", userHandler.GetAllUsers)
	// api.GET("/users/:id", userHandler.)
	api.PUT("/users/:id", userHandler.UpdateUser)
	// api.DELETE("/users/:id", userHandler.DeleteUser)

	// Payment routes
	api.GET("/payments", paymentHandler.GetAllPayments)
	api.GET("/:userId/payments", paymentHandler.GetPaymentsByUser)
	api.POST("/:userId/payments", paymentHandler.CreatePayment)
}
