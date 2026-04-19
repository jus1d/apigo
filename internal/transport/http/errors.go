package http

import (
	"log/slog"
	"net/http"

	"api/pkg/apierror"
	"api/pkg/log/sl"
	"api/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	reqID := requestid.Get(c)

	if he, ok := err.(*echo.HTTPError); ok && (he.Code == http.StatusNotFound || he.Code == http.StatusMethodNotAllowed) {
		_ = c.JSON(http.StatusNotFound, map[string]string{
			"message":    "resource not found",
			"request_id": reqID,
		})
		return
	}

	if ae, ok := err.(*apierror.Error); ok {
		slog.Debug("responded with API error", sl.Err(err), slog.String("request_id", reqID))
		_ = c.JSON(ae.Status, map[string]any{
			"message":    ae.Message,
			"request_id": reqID,
		})
		return
	}

	slog.Error("something went wrong", sl.Err(err), slog.String("request_id", reqID))
	_ = c.JSON(http.StatusInternalServerError, map[string]any{
		"message":    "internal server error",
		"request_id": reqID,
	})
}
