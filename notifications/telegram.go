// notifications/telegram.go
package notifications

import (
	"context"
	"fmt"
	"github/Rubncal04/youtube-premium/models"
	"log"
	"net/http"
	"net/url"
	"time"
)

// TelegramService implementa NotificationService para enviar mensajes usando Telegram.
type TelegramService struct {
	BotToken string
	ChatID   string // Si cada usuario tuviera su chat id, podrías omitir este campo y obtenerlo desde el modelo.
}

// NewTelegramService crea una nueva instancia de TelegramService.
func NewTelegramService(botToken, chatID string) *TelegramService {
	return &TelegramService{
		BotToken: botToken,
		ChatID:   chatID,
	}
}

// SendReminder envía un recordatorio a través de Telegram.
// En este ejemplo, se utiliza el chat id global configurado, pero podrías modificarlo para que sea específico del usuario.
func (ts *TelegramService) SendReminder(user models.Client, message string) error {
	// En este ejemplo, usamos el ChatID global; si el modelo tuviera un campo TelegramChatID, lo usarías aquí.
	encodedMessage := url.QueryEscape(message)
	urlStr := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", ts.BotToken, ts.ChatID, encodedMessage)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("error creating Telegram request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API returned status: %s", resp.Status)
	}

	log.Printf("Telegram reminder sent to user %s, status: %s", user.ID.Hex(), resp.Status)
	return nil
}
