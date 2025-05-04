package supabase

import "github.com/gofiber/fiber/v2"

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	grp := router.Group("/supabase")
	grp.Post("/", h.Upload)
}

type uploadForm struct {
	Bucket     string `form:"bucket"`
	ObjectPath string `form:"path"`
	FileName   string `form:"path"`
	Mime       string `form:"mime"`
}

func (h *Handler) Upload(c *fiber.Ctx) error {
	// 1. Parsear los campos del multipart/form-data
	var form uploadForm
	if err := c.BodyParser(&form); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "formato multipart inválido")
	}
	if form.Bucket == "" || form.ObjectPath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "bucket y path son requeridos")
	}

	// 2. Obtener el archivo
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	f, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	defer f.Close()

	// 3. Subir a Supabase
	url, err := h.svc.Upload(
		c.Context(),
		form.Bucket,
		form.ObjectPath,
		form.FileName,
		f,
		fileHeader.Header.Get("Content-Type"), // usa MIME real del archivo
		false,
	)
	if err != nil {
		// Devuelve el mensaje real del servidor
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"publicURL": url, // la URL a la que subiste (no firmada)
	})
}
