package tts

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterTTSRoutes(router fiber.Router) {
	grp := router.Group("/elevenlabs")
	grp.Post("/TTS", h.GenerateTTS)
	grp.Post("/SFX", h.GenerateSFX)
}

func (h *Handler) GenerateTTS(c *fiber.Ctx) error {
	var req RequestElevenTTS
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido")
	}

	audio, err := h.svc.TextToSpeechElevenlabs(req.Text, req.VoiceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Content-Type", "audio/mpeg")
	c.Set("Content-Disposition", "attachment; filename=voz.mp3")
	fmt.Println("Bytes recibidos:", len(audio))
	return c.Send(audio)
}

func (h *Handler) GenerateSFX(c *fiber.Ctx) error {
	var req RequestElevenSFX
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido")
	}

	sfx, err := h.svc.TextToSoundEffects(req.Description, req.DurationSeconds,
		req.PromptInfluence, req.OutputFormat)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Content-Type", "audio/mpeg")
	c.Set("Content-Disposition", "attachment; filename=sfx.mp3")
	fmt.Println("Bytes recibidos:", len(sfx))
	return c.Send(sfx)
}
