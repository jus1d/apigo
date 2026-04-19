package v1

import (
	"github.com/labstack/echo/v4"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Register(group *echo.Group) {
	group.GET("/liveness", liveness)
	group.GET("/readiness", readiness)
}
