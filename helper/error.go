package helper

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func JSONError(c *fiber.Ctx, status int, msg string, details ...any) error {
	res := ErrorResponse{
		Error: msg,
	}
	if len(details) > 0 {
		res.Details = details[0]
	}
	return c.Status(status).JSON(res)
}
