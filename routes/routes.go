package routes

import (
	"github/Rubncal04/youtube-premium/cache"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/handlers"
	"github/Rubncal04/youtube-premium/middleware"
	"github/Rubncal04/youtube-premium/repository"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes define las rutas principales de la aplicación.
func RegisterRoutes(e *echo.Echo, mongoRepo *db.MongoRepo, redisCache *cache.RedisCache, secretKey string) {
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

	// Initialize repositories
	clientRepo := repository.NewClientRepository(mongoRepo, redisCache)
	paymentRepo := repository.NewPaymentRepository(mongoRepo, redisCache)

	// Initialize handlers
	clientHandler := handlers.NewClientHandler(clientRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentRepo, clientRepo)

	// Payment routes - specific routes first
	api.GET("/clients/:clientId/payments/:id", paymentHandler.GetOnePayment)
	api.GET("/clients/:clientId/payments", paymentHandler.GetPaymentsByClient)
	api.POST("/clients/:clientId/payments", paymentHandler.CreatePayment)
	api.GET("/payments", paymentHandler.GetAllPayments)

	// Client routes
	api.POST("/clients", clientHandler.CreateClient)
	api.GET("/clients", clientHandler.GetClients)
	api.GET("/clients/:id", clientHandler.GetClient)
	api.PUT("/clients/:id", clientHandler.UpdateClient)
	api.DELETE("/clients/:id", clientHandler.DeleteClient)
}
