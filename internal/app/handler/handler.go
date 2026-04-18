package handler

import (
	"net/http"

	"api/internal/config"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	config *config.Config
}

func New(c *config.Config) *Handler {
	return &Handler{
		config: c,
	}
}

func (h *Handler) Liveness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (h *Handler) Readiness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

type APIError struct {
	Status  int
	Message string
}

func Error(code int, message string) error {
	return &APIError{
		Status:  code,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return e.Message
}
