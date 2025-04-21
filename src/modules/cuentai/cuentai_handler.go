package cuentai

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	grp := router.Group("/cuentai")
	grp.Post("/flow", h.CuentAIFlow)
}

func (h *Handler) CuentAIFlow(c *fiber.Ctx) error {
	var req struct {
		Text string `json:"text"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido")
	}

	res, err := h.svc.CuentAIFlow(req.Text)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// preparamos el JSON de respuesta con Base64
	type flowResp struct {
		Lines    []string `json:"lines"`
		TTSClips []string `json:"tts_clips"`
		SFXClips []string `json:"sfx_clips"`
		Combined string   `json:"combined"`
	}
	out := flowResp{Lines: res.Lines}
	for _, clip := range res.TTSClips {
		out.TTSClips = append(out.TTSClips, base64.StdEncoding.EncodeToString(clip))
	}
	for _, clip := range res.SFXClips {
		out.SFXClips = append(out.SFXClips, base64.StdEncoding.EncodeToString(clip))
	}
	out.Combined = base64.StdEncoding.EncodeToString(res.Combined)

	return c.JSON(out)
}
