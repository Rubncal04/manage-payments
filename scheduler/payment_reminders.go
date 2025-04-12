// scheduler/payment_reminders.go
package scheduler

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/notifications"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// SendPaymentReminders consulta los usuarios y, si es el día de pago (o dentro de la ventana de 5 días)
// y no han pagado, envía un recordatorio usando el servicio de notificaciones proporcionado.
func SendPaymentReminders(mongoRepo *db.MongoRepo, notifier notifications.NotificationService) {
	now := time.Now()
	currentDay := now.Day()

	var clients []models.Client
	// Traer solo usuarios que no han pagado.
	err := mongoRepo.FindAll("clients", bson.M{"status": "inactive"}, &clients)
	if err != nil {
		log.Printf("Error retrieving clients: %v", err)
		return
	}

	for _, client := range clients {
		dueDay := client.DayToPay

		// Si el día actual está dentro de la ventana de 5 días a partir del día de pago.
		if currentDay >= dueDay && currentDay <= dueDay+4 {
			message := "Hola, te recuerdo el compromiso que tienes con YouTube Premium. ¡Quédate al día con tu pago! 😉"
			log.Printf("Sending reminder to client %s via Whatsapp...", client.Name)

			// Enviar recordatorio usando el servicio de notificaciones.
			err := notifier.SendReminder(client, message)
			if err != nil {
				log.Printf("Error sending reminder to client %s: %v", client.Name, err)
			}
		}
	}
}
