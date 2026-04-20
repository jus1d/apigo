package http

import (
	"log/slog"
	"net/http"

	"apigo/pkg/apierror"
	"apigo/pkg/apiresponse"
	"apigo/pkg/log/sl"
	"apigo/pkg/requestid"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	reqID := requestid.Get(c)

	if he, ok := err.(*echo.HTTPError); ok && (he.Code == http.StatusNotFound || he.Code == http.StatusMethodNotAllowed) {
		_ = apiresponse.Error(c, http.StatusNotFound, apierror.TypeNotFound, "resource not found", "Check the URL and HTTP method")
		return
	}

	if ae, ok := err.(*apierror.Error); ok {
		slog.Debug("responded with API error", sl.Err(err), slog.String("request_id", reqID))
		_ = apiresponse.Error(c, ae.Status, ae.Type, ae.Message, ae.Hint)
		return
	}

	slog.Error("something went wrong", sl.Err(err), slog.String("request_id", reqID))
	_ = apiresponse.Error(c, http.StatusInternalServerError, apierror.TypeInternal, "internal server error", "Try again later or contact support")
}
