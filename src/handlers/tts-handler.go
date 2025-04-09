package handlers

import (
	"fmt"

	"github.com/MetaDandy/cuent-ai-core/src/dtos"
	"github.com/MetaDandy/cuent-ai-core/src/services"
	"github.com/gofiber/fiber/v2"
)

type TTSHandler struct {
	service *services.TTSService
}

func NewTTSHandler(service *services.TTSService) *TTSHandler {
	return &TTSHandler{service}
}

func (h *TTSHandler) RegisterTTSRoutes(router fiber.Router) {
	router.Post("/generate/elevenlabs", h.GenerateElevenlabs)
}

func (h *TTSHandler) GenerateElevenlabs(c *fiber.Ctx) error {
	var req dtos.RequestElevenTTS
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido")
	}

	audio, err := h.service.TextToSpeechElevenlabs(req.Text, req.VoiceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Content-Type", "audio/mpeg")
	c.Set("Content-Disposition", "attachment; filename=voz.mp3")
	fmt.Println("Bytes recibidos:", len(audio))
	return c.Send(audio)
}
