package script

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
	grp := router.Group("/scripts")
	grp.Get("", h.FindAll)
	grp.Get("/:id", h.FindById)
	grp.Post("", h.Create)
	// grp.Patch("/:id", h.Update)
	// grp.Delete("/:id", h.SoftDelete)
	// grp.Post("/:id", h.Restore)
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
			"Error obteniendo projecto", err.Error())
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
	var input ScriptCreate
	if err := c.BodyParser(&input); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Input inválido", err.Error())
	}
	project, err := h.svc.Create(&input)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error creando script", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    project,
		Message: "Script creado",
	})
}
