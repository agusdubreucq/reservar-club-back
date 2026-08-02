package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type WhatsAppService struct {
	twilioAccountSID  string
	twilioAuthToken   string
	twilioPhoneNumber string
}

type MessagePayload struct {
	From  string
	To    string
	Body  string
	Sid   string
}

type WhatsAppMessageRequest struct {
	PhoneNumber string
	Message     string
}

type WhatsAppMessageResponse struct {
	MessageSID string
	Status     string
}

func NewWhatsAppService(
	accountSID string,
	authToken string,
	phoneNumber string,
) *WhatsAppService {
	return &WhatsAppService{
		twilioAccountSID:  accountSID,
		twilioAuthToken:   authToken,
		twilioPhoneNumber: phoneNumber,
	}
}

// ReceiveMessage procesa un mensaje entrante y responde automaticamente
func (w *WhatsAppService) ReceiveMessage(payload MessagePayload) (string, error) {
	log.Printf("📨 Mensaje recibido de %s: %s", payload.From, payload.Body)

	phone := w.ExtractPhoneNumber(payload.From)

	// Respuesta automática para demo
	response := "¡Hola! 👋 Soy el bot de reservas del club.\n\nDisponibilidades y reservas próximamente... ⏳"

	// Enviar respuesta
	_, err := w.SendMessage(WhatsAppMessageRequest{
		PhoneNumber: phone,
		Message:     response,
	})
	if err != nil {
		log.Printf("❌ Error enviando respuesta: %v", err)
		return "", err
	}

	log.Printf("✅ Respuesta enviada a %s", phone)
	return response, nil
}

// SendMessage envía un mensaje a un número de WhatsApp usando Twilio REST API
func (w *WhatsAppService) SendMessage(req WhatsAppMessageRequest) (*WhatsAppMessageResponse, error) {
	if w.twilioAccountSID == "" || w.twilioAuthToken == "" {
		return nil, fmt.Errorf("Twilio credentials not configured")
	}

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", w.twilioAccountSID)

	formData := url.Values{}
	formData.Set("From", "whatsapp:"+w.twilioPhoneNumber)
	formData.Set("To", "whatsapp:"+req.PhoneNumber)
	formData.Set("Body", req.Message)

	httpReq, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(w.twilioAccountSID + ":" + w.twilioAuthToken))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Twilio error: %s", string(body))
		return nil, fmt.Errorf("twilio API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &WhatsAppMessageResponse{
		MessageSID: result["sid"].(string),
		Status:     result["status"].(string),
	}, nil
}

// ExtractPhoneNumber limpia el formato del número de teléfono de Twilio
func (w *WhatsAppService) ExtractPhoneNumber(twilioPhone string) string {
	if strings.HasPrefix(twilioPhone, "whatsapp:") {
		return twilioPhone[9:]
	}
	return twilioPhone
}
