package v1

import (
	"api/pkg/log/sl"
	"api/pkg/validate"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func echoHandler(c echo.Context) error {
	var body map[string]any
	if err := validate.Bind(c, &body); err != nil {
		slog.Error("invalid request body", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
	}

	return c.JSON(http.StatusOK, body)
}
