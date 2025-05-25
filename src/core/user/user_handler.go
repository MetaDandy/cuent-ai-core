package user

import (
	"log"
	"net/http"
	"os"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82/webhook"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	grp := router.Group("/users")
	grp.Post("/sign-up", h.SignUp)
	grp.Post("/sign-in", h.SignIn)
	grp.Post("/stripe-web-hook", h.StripeWebhook)

	grp.Use(middleware.JwtMiddleware())

	grp.Get("", h.FindAll)
	grp.Get("/profile", h.GetProfile)
	grp.Get("/subscription", h.GetActiveSubscription)
	grp.Get("/:id", h.FindById)
	grp.Post("/:id/add-subscription", h.AddSubscription)
	grp.Post("/payment-subscription", h.PaymentSubscription)
	grp.Patch("/change-password", h.ChangePassoword)
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

func (h *Handler) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.svc.FindById(id)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo usuario", err.Error())
	}

	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    user,
		Message: "User encontrado",
	})

}
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	user, err := h.svc.FindById(id)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obteniendo perfil", err.Error())
	}

	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    user,
		Message: "Perfil encontrado",
	})
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var singup Singup
	if err := c.BodyParser(&singup); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Sign Up inválido", err.Error())
	}

	users, token, err := h.svc.SignUp(&singup)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error al registrar usuario", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    users,
		Token:   token,
		Message: "User registrado",
	})
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	var signin Signin
	if err := c.BodyParser(&signin); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Login inválido", err.Error())
	}

	users, token, err := h.svc.Signin(&signin)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error al autenticar usuario", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    users,
		Token:   token,
		Message: "User autenticado",
	})
}

func (h *Handler) ChangePassoword(c *fiber.Ctx) error {
	var changepass ChangePassoword
	if err := c.BodyParser(&changepass); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Input inválido", err.Error())
	}

	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	data, err := h.svc.ChangePassword(id, &changepass)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error  cambiando la contraseña", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    data,
		Message: "Contraseña cambiada",
	})
}

func (h *Handler) AddSubscription(c *fiber.Ctx) error {
	subID := c.Params("id")
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	data, err := h.svc.AddSubscription(id, subID)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error añadiendo la suscripción", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    data,
		Message: "Suscripción añadida",
	})
}

func (h *Handler) StripeWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	log.Println(secret)

	evt, err := webhook.ConstructEvent(payload, sigHeader, secret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("signature error")
	}

	// Procesa sólo lo que te interesa
	if resp, err := h.svc.StripeWebhook(evt); err != nil {
		log.Printf("[stripe] %v", err) // ⬅️ Log interno
		return c.Status(fiber.StatusInternalServerError).SendString("internal error")
	} else if resp != nil {
		return c.Status(fiber.StatusOK).JSON(helper.Response{
			Data:    resp,
			Message: "Suscripción activada",
		})
	}

	// Evento ignorado o sin DTO que devolver
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) PaymentSubscription(c *fiber.Ctx) error {
	var payment Payment
	if err := c.BodyParser(&payment); err != nil {
		return helper.JSONError(c, http.StatusBadRequest,
			"Pago invalido", err.Error())
	}

	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	data, err := h.svc.PaySubscription(id, payment)
	if err != nil {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Error en el pago", "")
	}

	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    data,
		Message: "Procesando pago",
	})
}

func (h *Handler) GetActiveSubscription(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return helper.JSONError(c, http.StatusUnauthorized,
			"Token sin user_id", "")
	}

	data, err := h.svc.GetActiveSubscription(id)
	if err != nil {
		return helper.JSONError(c, http.StatusInternalServerError,
			"Error obtiendo la suscripción actual", err.Error())
	}
	return c.Status(http.StatusCreated).JSON(helper.Response{
		Data:    data,
		Message: "Suscripción actual",
	})
}
