package notifications

import "github/Rubncal04/youtube-premium/models"

// NotificationService define la interfaz para enviar notificaciones.
type NotificationService interface {
	// SendReminder envía un recordatorio al usuario indicado.
	SendReminder(user models.Client, message string) error
}
