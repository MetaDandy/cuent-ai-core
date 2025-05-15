package helper

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

type HTTPError struct {
	StatusCode int
	Header     http.Header
	Message    string
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

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}
