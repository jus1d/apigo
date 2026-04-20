package apiresponse

import (
	"api/pkg/apierror"
	"api/pkg/requestid"

	"github.com/labstack/echo/v4"
)

type ErrorBody struct {
	Type    apierror.Type `json:"type"`
	Message string        `json:"message"`
	Hint    string        `json:"hint"`
}

type ErrorResponse struct {
	Type      string    `json:"type"`
	RequestID string    `json:"request_id"`
	Error     ErrorBody `json:"error"`
}

type SuccessResponse struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

type CollectionMeta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type CollectionResponse struct {
	Type      string         `json:"type"`
	RequestID string         `json:"request_id"`
	Data      any            `json:"data"`
	Meta      CollectionMeta `json:"meta"`
}

func Success(c echo.Context, status int, data any) error {
	return c.JSON(status, SuccessResponse{
		Type:      "success",
		RequestID: requestid.Get(c),
		Data:      data,
	})
}

func Collection(c echo.Context, status int, data any, meta CollectionMeta) error {
	return c.JSON(status, CollectionResponse{
		Type:      "collection",
		RequestID: requestid.Get(c),
		Data:      data,
		Meta:      meta,
	})
}

func Error(c echo.Context, status int, errType apierror.Type, message string, hint string) error {
	return c.JSON(status, ErrorResponse{
		Type:      "error",
		RequestID: requestid.Get(c),
		Error: ErrorBody{
			Type:    errType,
			Message: message,
			Hint:    hint,
		},
	})
}
