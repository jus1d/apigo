package router

import (
	"net/http"

	"api/internal/app/handler"
	"api/internal/config"
	"api/pkg/requestid"
	"api/pkg/requestlog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	config  *config.Config
	handler *handler.Handler
}

func New(c *config.Config) *Router {
	h := handler.New(c)
	return &Router{config: c, handler: h}
}

func (r *Router) InitRoutes() *echo.Echo {
	router := echo.New()

	router.HTTPErrorHandler = handler.HTTPErrorHandler

	router.Use(requestid.New)
	router.Use(requestlog.Completed)
	router.Pre(middleware.RemoveTrailingSlash())

	switch r.config.Env {
	case config.EnvLocal, config.EnvDevelopment:
		router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
				c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

				if c.Request().Method == "OPTIONS" {
					return c.NoContent(http.StatusNoContent)
				}

				return next(c)
			}
		})
	}

	// TODO: update rate limiting logic:
	//   Current:
	//     - request -> wait Ns -> request
	//
	//   Expected:
	//     - [request -> request -> request] - in such window, forbid to make more than M requests
	//       ^ 0s                       Ns ^

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/liveness", r.handler.Liveness)
	v1.GET("/readiness", r.handler.Readiness)

	return router
}
