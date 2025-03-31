package routes

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/handlers"
	"github/Rubncal04/youtube-premium/repository"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes define las rutas principales de la aplicación.
func RegisterRoutes(e *echo.Echo, mongoRepo *db.MongoRepo) {
	e.GET("/", handlers.HelloWorld)

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(mongoRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	paymentRepo := repository.NewPaymentRepository(mongoRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentRepo)

	private := e.Group("/api/v1")

	// User routes
	private.GET("/users", userHandler.GetAllUsers)
	private.POST("/users", userHandler.CreateUser)
	private.PUT("/users/:id", userHandler.UpdateUser)

	// Payment routes
	private.GET("/payments", paymentHandler.GetAllPayments)
	private.GET("/:userId/payments", paymentHandler.GetPaymentsByUser)
	private.POST("/:userId/payments", paymentHandler.CreatePayment)
}
