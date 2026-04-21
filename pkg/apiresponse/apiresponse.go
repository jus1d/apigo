package apiresponse

import (
	"apigo/pkg/apierror"
	"apigo/pkg/requestid"

	"github.com/labstack/echo/v4"
)

type ErrorBody struct {
	Code    apierror.Code `json:"code"`
	Message string        `json:"message"`
	Hint    string        `json:"hint"`
}

type ErrorResponse struct {
	Kind      string    `json:"kind"`
	RequestID string    `json:"request_id"`
	Error     ErrorBody `json:"error"`
}

type SuccessResponse struct {
	Kind      string `json:"kind"`
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

type CollectionMeta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type CollectionResponse[T any] struct {
	Kind      string         `json:"kind"`
	RequestID string         `json:"request_id"`
	Data      []T            `json:"data"`
	Meta      CollectionMeta `json:"meta"`
}

func Success(c echo.Context, status int, data any) error {
	return c.JSON(status, SuccessResponse{
		Kind:      "success",
		RequestID: requestid.Get(c),
		Data:      data,
	})
}

func Collection[T any](c echo.Context, status int, data []T, meta CollectionMeta) error {
	return c.JSON(status, CollectionResponse[T]{
		Kind:      "collection",
		RequestID: requestid.Get(c),
		Data:      data,
		Meta:      meta,
	})
}

func Error(c echo.Context, status int, code apierror.Code, message string, hint string) error {
	return c.JSON(status, ErrorResponse{
		Kind:      "error",
		RequestID: requestid.Get(c),
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Hint:    hint,
		},
	})
}
