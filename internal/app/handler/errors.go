package handler

import (
	"log/slog"
	"net/http"

	"api/internal/lib/log/sl"
	"api/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	requestID := requestid.Get(c)
	if he, ok := err.(*echo.HTTPError); ok && (he.Code == http.StatusNotFound || he.Code == http.StatusMethodNotAllowed) {
		c.JSON(http.StatusNotFound, map[string]string{
			"message":    "resource not found",
			"request_id": requestID,
		})
		return
	}

	if ae, ok := err.(*APIError); ok {
		slog.Debug("responded with API error", sl.Err(err), slog.String("request_id", requestid.Get(c)))
		c.JSON(ae.Status, map[string]any{
			"message":    ae.Message,
			"request_id": requestID,
		})
		return
	}

	slog.Error("something went wrong", sl.Err(err), slog.String("request_id", requestid.Get(c)))
	c.JSON(http.StatusInternalServerError, map[string]any{
		"message":    "internal server error",
		"request_id": requestID,
	})
}
