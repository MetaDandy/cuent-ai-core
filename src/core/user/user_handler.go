package user

import (
	"net/http"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	grp := router.Group("/users")
	grp.Get("", h.FindAll)
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	opts := helper.NewFindAllOptionsFromQuery(c)
	users, err := h.svc.FindAll(opts)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo usuarios", err.Error())
	}
	return c.JSON(users)
}
