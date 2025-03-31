package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/Rubncal04/youtube-premium/models"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TwilioService implementa el envío de mensajes usando la API de Twilio para WhatsApp.
type TwilioService struct {
	AccountSID   string
	AuthToken    string
	FromWhatsApp string // Número de WhatsApp de Twilio, en formato: "whatsapp:+14155238886"
}

// NewTwilioService crea una nueva instancia de TwilioService.
func NewTwilioService(accountSID, authToken, fromWhatsApp string) *TwilioService {
	return &TwilioService{
		AccountSID:   accountSID,
		AuthToken:    authToken,
		FromWhatsApp: fromWhatsApp,
	}
}

// SendReminder envía un recordatorio a través de WhatsApp usando la API de Twilio.
// Se espera que el modelo User tenga un campo CellPhone que contenga el número de teléfono
// del usuario en formato internacional, sin símbolos, por ejemplo: "573001234567".
func (ts *TwilioService) SendReminder(user models.User, message string) error {
	// Construir la URL del endpoint de Twilio.
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", ts.AccountSID)

	// Construir los datos del formulario.
	formData := url.Values{}
	formData.Set("To", fmt.Sprintf("whatsapp:+%s", user.CellPhone))
	formData.Set("From", ts.FromWhatsApp)
	formData.Set("Body", message)

	// Crear el request con timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating Twilio request: %w", err)
	}

	// Configurar headers y autenticación básica.
	req.SetBasicAuth(ts.AccountSID, ts.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Twilio request: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta para debug (puedes parsearla si lo deseas)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Twilio response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio API error, status: %s, response: %s", resp.Status, string(body))
	}

	// Opcional: parsear la respuesta JSON si se requiere.
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error parsing Twilio response: %w", err)
	}

	return nil
}
