package script

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
	grp := router.Group("/scripts").Use(middleware.JwtMiddleware())
	grp.Get("", h.FindAll)
	grp.Get("/:id", h.FindById)
	grp.Post("", h.Create)
	grp.Post("/manual-create", h.ManualCreate)
	grp.Post("/:id/mixed", h.MixAudio)
	grp.Patch("/:id/regenerate", h.Regenerate)
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	opts := helper.NewFindAllOptionsFromQuery(c)
	project, err := h.svc.FindAll(opts)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo scripts", err.Error())
	}
	return c.JSON(project)
}

func (h *Handler) FindById(c *fiber.Ctx) error {
	dto, err := h.svc.FindByID(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo script", err.Error())
	}
	if dto == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"Script no encontrado")
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "Script obtenido",
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	var input ScriptCreate
	if err := c.BodyParser(&input); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Input inválido", err.Error())
	}

	project, err := h.svc.Create(id, &input)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error creando script", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    project,
		Message: "Script creado",
	})
}

func (h *Handler) ManualCreate(c *fiber.Ctx) error {
	var input ScriptManualCreate
	if err := c.BodyParser(&input); err != nil {
		return helper.JSONError(c, fiber.StatusBadRequest,
			"Input inválido", err.Error())
	}

	script, err := h.svc.ManualCreate(&input)
	if err != nil {
		return helper.JSONError(c, fiber.StatusInternalServerError,
			"Error creando script", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(helper.Response{
		Data:    script,
		Message: "Script creado",
	})
}

func (h *Handler) Regenerate(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	project, err := h.svc.Regenerate(id, c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error creando script", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    project,
		Message: "Script creado",
	})
}

func (h *Handler) MixAudio(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	dto, err := h.svc.MixAudio(c.Params("id"), id)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error al mixear assets", err.Error())
	}
	if dto == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"Script no encontrado")
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "Mixeo de audio de assets generado",
	})
}
