package notifications

import (
	"fmt"
	"github/Rubncal04/youtube-premium/models"
	"regexp"
	"strings"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioService implements WhatsApp message sending using Twilio's API
type TwilioService struct {
	client       *twilio.RestClient
	fromWhatsApp string
}

// NewTwilioService creates a new instance of TwilioService
func NewTwilioService(accountSID, authToken, fromWhatsApp string) *TwilioService {
	// Ensure the fromWhatsApp number is properly formatted
	if !strings.HasPrefix(fromWhatsApp, "whatsapp:") {
		fromWhatsApp = "whatsapp:" + fromWhatsApp
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioService{
		client:       client,
		fromWhatsApp: fromWhatsApp,
	}
}

// validatePhoneNumber ensures the phone number is in the correct format
func (ts *TwilioService) validatePhoneNumber(phone string) error {
	// Remove any non-digit characters
	cleanPhone := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Check if the phone number is between 10 and 15 digits
	if len(cleanPhone) < 10 || len(cleanPhone) > 15 {
		return fmt.Errorf("invalid phone number length: %s", phone)
	}

	return nil
}

// formatWhatsAppNumber formats a phone number for WhatsApp messaging
func (ts *TwilioService) formatWhatsAppNumber(phone string) string {
	// Remove any non-digit characters
	cleanPhone := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Ensure the number starts with a plus sign
	if !strings.HasPrefix(cleanPhone, "+") {
		cleanPhone = "+57" + cleanPhone
	}

	return "whatsapp:" + cleanPhone
}

// SendReminder sends a WhatsApp reminder using Twilio's API
func (ts *TwilioService) SendReminder(client models.Client, message string) error {
	// Validate phone number
	if err := ts.validatePhoneNumber(client.CellPhone); err != nil {
		return fmt.Errorf("invalid phone number for client %s: %w", client.ID.Hex(), err)
	}

	// Format the WhatsApp number
	toNumber := ts.formatWhatsAppNumber(client.CellPhone)

	fmt.Println("toNumber", toNumber)

	// Create message parameters
	params := &api.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(ts.fromWhatsApp)
	params.SetBody(message)

	// Send message using Twilio SDK
	messageResponse, err := ts.client.Api.CreateMessage(params)
	if err != nil {
		// Check for specific Twilio error codes
		if strings.Contains(err.Error(), "63007") {
			return fmt.Errorf("invalid WhatsApp configuration: %w\nPlease ensure:\n1. Your WhatsApp number is verified in Twilio\n2. The number format is correct (whatsapp:+1234567890)\n3. WhatsApp is enabled for your account", err)
		}
		return fmt.Errorf("error sending WhatsApp message: %w", err)
	}

	// Log successful message
	fmt.Printf("Message sent successfully to %s. SID: %s, Status: %s\n",
		client.CellPhone, *messageResponse.Sid, *messageResponse.Status)

	return nil
}
