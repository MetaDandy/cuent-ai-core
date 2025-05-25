package asset

import (
	"net/http"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	grp := router.Group("/assets").Use(middleware.JwtMiddleware())
	grp.Get("", h.FindAll)
	grp.Get("/:id", h.FindById)
	grp.Get("/:id/script", h.FindByScriptID)
	grp.Post("/:id", h.GenerateOne)
	grp.Post("/:id/generate_all", h.GenerateAll)
	grp.Post("/:id/regenerate_all", h.RegenerateAll)
	grp.Post("/:id/generate_video", h.GenerateOneVideo)
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

func (h *Handler) FindByScriptID(c *fiber.Ctx) error {
	dto, err := h.svc.FindByScriptID(c.Params("id"))
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
		Message: "Todos los assets obtenidos",
	})
}

func (h *Handler) GenerateOne(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	dto, err := h.svc.GenerateOne(c.Params("id"), id)
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
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	dto, err := h.svc.GenerateAll(c.Params("id"), id, false)
	if err != nil {
		return c.JSON(helper.Response{
			Data:    dto,
			Message: err.Error(),
		})
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "assets generados",
	})
}

func (h *Handler) RegenerateAll(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	dto, err := h.svc.GenerateAll(c.Params("id"), id, true)
	if err != nil {
		return c.JSON(helper.Response{
			Data:    dto,
			Message: err.Error(),
		})
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "assets regenerados",
	})
}

func (h *Handler) GenerateOneVideo(c *fiber.Ctx) error {
	var key_words GenerateVideo
	if err := c.BodyParser(&key_words); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Input inválido", err.Error())
	}

	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	dto, err := h.svc.GenerateVideo(c.Params("id"), id, key_words)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error generando el asset video", err.Error())
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "video generado",
	})
}
