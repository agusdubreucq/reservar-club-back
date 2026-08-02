package handlers

import (
	"log"
	"net/http"

	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type WhatsAppHandler struct {
	whatsappService *services.WhatsAppService
}

func NewWhatsAppHandler(whatsappService *services.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService: whatsappService,
	}
}

// WebhookRequest estructura del webhook de Twilio
type WebhookRequest struct {
	From       string `form:"From" binding:"required"`
	To         string `form:"To" binding:"required"`
	Body       string `form:"Body" binding:"required"`
	MessageSid string `form:"MessageSid" binding:"required"`
}

// HandleIncomingMessage recibe mensajes de WhatsApp desde Twilio
// POST /api/whatsapp/webhook
func (h *WhatsAppHandler) HandleIncomingMessage(c *gin.Context) {
	var req WebhookRequest

	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Error binding webhook request: %v", err)
		c.Header("Content-Type", "application/xml")
		c.String(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><Response></Response>`)
		return
	}

	payload := services.MessagePayload{
		From: req.From,
		To:   req.To,
		Body: req.Body,
		Sid:  req.MessageSid,
	}

	_, err := h.whatsappService.ReceiveMessage(payload)
	if err != nil {
		log.Printf("Error processing message: %v", err)
	}

	// Twilio espera respuesta vacía (XML)
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><Response></Response>`)
}

// SendTestMessage envía un mensaje de prueba a un número
// POST /api/whatsapp/send
type SendMessageRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

func (h *WhatsAppHandler) SendTestMessage(c *gin.Context) {
	var req SendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.whatsappService.SendMessage(services.WhatsAppMessageRequest{
		PhoneNumber: req.PhoneNumber,
		Message:     req.Message,
	})

	if err != nil {
		log.Printf("Error sending message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_sid": response.MessageSID,
		"status":      response.Status,
	})
}

// HealthCheck endpoint para verificar que el webhook está activo
// GET /api/whatsapp/health
func (h *WhatsAppHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "WhatsApp webhook is active",
	})
}
