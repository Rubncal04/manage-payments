package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github/Rubncal04/youtube-premium/cache"
	"github/Rubncal04/youtube-premium/config"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/notifications"
	"github/Rubncal04/youtube-premium/routes"
	"github/Rubncal04/youtube-premium/scheduler"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

func StartServer() {
	// Get environment variables
	envVariables := config.GetVariables()
	secretKey := envVariables.JWT_SECRET_KEY
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // URL de tu aplicación React
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Initialize MongoDB repository
	mongoRepo, err := db.NewMongoRepo(envVariables)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoRepo.Close()

	// Initialize Redis cache
	redisDB, err := strconv.Atoi(envVariables.REDIS_DATABASES)
	if err != nil {
		redisDB = 0
	}
	redisCache, err := cache.NewRedisCache(envVariables.REDIS_ADDRESS, envVariables.REDIS_PASSWORD, envVariables.REDIS_PORT, redisDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without cache...")
		redisCache = nil
	}

	// Root route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to YouTube Premium API")
	})

	// Register all routes
	routes.RegisterRoutes(e, mongoRepo, redisCache, secretKey)

	// Start server
	port := envVariables.PORT
	if port == "" {
		port = "9120"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Printf("Servidor iniciando en el puerto %s...\n", port)

	// Graceful shutdown
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	twilioAccountSID := envVariables.TWILIO_ACCOUNT_SID
	twilioAuthToken := envVariables.TWILIO_AUTH_TOKEN
	twilioFromWhatsApp := envVariables.TWILIO_FROM_WHATSAPP

	if twilioAccountSID == "" || twilioAuthToken == "" || twilioFromWhatsApp == "" {
		log.Fatalf("TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN o TWILIO_FROM_WHATSAPP no están configurados")
	}
	twilioService := notifications.NewTwilioService(twilioAccountSID, twilioAuthToken, twilioFromWhatsApp)

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Fatalf("Error loading location: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	// Add payment reminder task
	_, err = c.AddFunc("0 16 * * *", func() {
		log.Println("Running payment verification...")
		scheduler.SendPaymentReminders(mongoRepo, twilioService)
	})
	if err != nil {
		log.Fatalf("Error scheduling payment reminder task: %v", err)
	}

	// Add payment status update tasks
	// Run on the 13th of each month
	_, err = c.AddFunc("0 0 13 * *", func() {
		log.Println("Running payment status update for users with payment dates 15-20...")
		if err := scheduler.UpdatePaymentStatus(mongoRepo); err != nil {
			log.Printf("Error updating payment status: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling payment status update task (13th): %v", err)
	}

	// Run on the 25th of each month
	_, err = c.AddFunc("0 0 25 * *", func() {
		log.Println("Running payment status update for users with payment dates 28-30...")
		if err := scheduler.UpdatePaymentStatus(mongoRepo); err != nil {
			log.Printf("Error updating payment status: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling payment status update task (25th): %v", err)
	}

	c.Start()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
