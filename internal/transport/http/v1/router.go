package v1

import (
	"api/internal/config"
	"api/internal/transport/http/middleware"
	"api/pkg/requestid"
	"api/pkg/requestlog"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	config *config.Config
}

func NewRouter(c *config.Config) *Router {
	return &Router{config: c}
}

func (r *Router) InitRoutes() *echo.Echo {
	router := echo.New()

	router.HTTPErrorHandler = httpErrorHandler

	router.Use(requestid.New)
	router.Use(requestlog.Completed)
	router.Pre(echomw.RemoveTrailingSlash())

	switch r.config.Env {
	case config.EnvLocal, config.EnvDevelopment:
		router.Use(middleware.CORS)
	}

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/liveness", liveness)
	v1.GET("/readiness", readiness)

	return router
}
