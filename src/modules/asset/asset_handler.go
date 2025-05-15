package asset

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
	grp := router.Group("/assets")
	grp.Get("", h.FindAll)
	grp.Post("/:id", h.GenerateOne)
	grp.Post("/:id/generate_all", h.GenerateAll)
	grp.Get("/:id", h.FindById)
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	opts := helper.NewFindAllOptionsFromQuery(c)
	project, err := h.svc.FindAll(opts)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo assets", err.Error())
	}
	return c.JSON(project)
}

func (h *Handler) FindById(c *fiber.Ctx) error {
	dto, err := h.svc.FindByID(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo asset", err.Error())
	}
	if dto == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"asset no encontrado")
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "asset obtenido",
	})
}

func (h *Handler) GenerateOne(c *fiber.Ctx) error {
	dto, err := h.svc.GenerateOne(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error generando el asset", err.Error())
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "asset generado",
	})
}

func (h *Handler) GenerateAll(c *fiber.Ctx) error {
	dto, err := h.svc.GenerateAll(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error todos los assets", err.Error())
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "assets generados",
	})
}
