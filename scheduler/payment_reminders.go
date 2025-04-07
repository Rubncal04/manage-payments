// scheduler/payment_reminders.go
package scheduler

import (
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"
	"github/Rubncal04/youtube-premium/notifications"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// SendPaymentReminders consulta los usuarios y, si es el día de pago (o dentro de la ventana de 5 días)
// y no han pagado, envía un recordatorio usando el servicio de notificaciones proporcionado.
func SendPaymentReminders(mongoRepo *db.MongoRepo, notifier notifications.NotificationService) {
	now := time.Now()
	currentDay := now.Day()

	var users []models.User
	// Traer solo usuarios que no han pagado.
	err := mongoRepo.FindAll("users", bson.M{"paid": false}, &users)
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		return
	}

	for _, user := range users {
		dueDay, err := strconv.Atoi(user.DateToPay)
		if err != nil {
			log.Printf("Error converting date_to_pay for user %s: %v", user.ID.Hex(), err)
			continue
		}

		// Si el día actual está dentro de la ventana de 5 días a partir del día de pago.
		if currentDay >= dueDay && currentDay <= dueDay+4 {
			message := "Hola, te recuerdo el compromiso que tienes con YouTube Premium. ¡Quédate al día con tu pago! 😉"
			log.Printf("Sending reminder to user %s via Whatsapp...", user.Name)

			// Enviar recordatorio usando el servicio de notificaciones.
			err := notifier.SendReminder(user, message)
			if err != nil {
				log.Printf("Error sending reminder to user %s: %v", user.Name, err)
			}
		}
	}
}
