package subscription

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
	grp := router.Group("/subscription")
	grp.Get("", h.FindAll)
	grp.Get("/:id", h.FindById)
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	opts := helper.NewFindAllOptionsFromQuery(c)
	subscription, err := h.svc.FindAll(opts)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo subscripciones", err.Error())
	}
	return c.JSON(subscription)
}

func (h *Handler) FindById(c *fiber.Ctx) error {
	dto, err := h.svc.FindByID(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo subscripciones", err.Error())
	}
	if dto == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"Subscripción no encontrado")
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "Subscripción obtenido",
	})
}
