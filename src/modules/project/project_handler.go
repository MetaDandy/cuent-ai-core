package project

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
	grp := router.Group("/projects")
	grp.Get("", h.FindAll)
	grp.Get("/:id", h.FindById)
	grp.Post("", h.Create)
	grp.Patch("/:id", h.Update)
	grp.Delete("/:id", h.SoftDelete)
	grp.Post("/:id", h.Restore)
}

func (h *Handler) FindAll(c *fiber.Ctx) error {
	opts := helper.NewFindAllOptionsFromQuery(c)
	project, err := h.svc.FindAll(opts)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo projectos", err.Error())
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
			"Projecto no encontrado")
	}

	return c.JSON(helper.Response{
		Data:    dto,
		Message: "Projecto obtenido",
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var input ProjectCreate
	if err := c.BodyParser(&input); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Input inválido", err.Error())
	}
	project, err := h.svc.Create(&input)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error creando projecto", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    project,
		Message: "Projecto creado",
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var input ProjectUpdate
	if err := c.BodyParser(&input); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Cuerpo inválido", err.Error())
	}
	project, err := h.svc.Update(c.Params("id"), &input)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error actualizando proejcto", err.Error())
	}
	if project == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"Proyecto no encontrado")
	}
	return c.JSON(helper.Response{
		Data:    project,
		Message: "Proyecto actualizado",
	})
}

func (h *Handler) SoftDelete(c *fiber.Ctx) error {
	ok, err := h.svc.SoftDelete(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error eliminando proyecto", err.Error())
	}
	if !ok {
		return helper.JSONError(c, http.StatusNotFound, "Projecto no encontrado")
	}
	return c.SendStatus(http.StatusNoContent)
}

func (h *Handler) Restore(c *fiber.Ctx) error {
	project, err := h.svc.Restore(c.Params("id"))
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error restaurando projecto", err.Error())
	}
	if project == nil {
		return helper.JSONError(c, http.StatusNotFound,
			"Proyecto no encontrado")
	}
	return c.JSON(helper.Response{
		Data:    project,
		Message: "Proyecto restaurado",
	})
}
