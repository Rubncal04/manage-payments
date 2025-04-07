package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github/Rubncal04/youtube-premium/config"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/notifications"
	"github/Rubncal04/youtube-premium/routes"
	"github/Rubncal04/youtube-premium/scheduler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

func StartServer() {
	e := echo.New()

	// Configurar las variables de entorno
	envVariables := config.GetVariables()
	port := envVariables.PORT
	if port == "" {
		port = ":8080"
	}

	// Set CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // URL de tu aplicación React
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Inicializar conexión a MongoDB
	mongoRepo, err := db.NewMongoRepo(envVariables)
	if err != nil {
		log.Fatalf("Error initializing MongoDB: %v", err)
	}
	defer mongoRepo.Close()

	// Configurar rutas
	routes.RegisterRoutes(e, mongoRepo)

	// Iniciar el servidor
	go func() {
		log.Printf("🚀 Server running on http://localhost%s", port)
		if err := e.Start(port); err != nil {
			log.Fatalf("Error starting server: %v", err)
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
	_, err = c.AddFunc("0 17 * * *", func() {
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

	// Capturar señales para cierre ordenado
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	// Cierre ordenado con contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	} else {
		log.Println("✅ Server shut down cleanly")
	}
}
