package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return e
}

/*
func New2() *echo.Echo {
	e := echo.New()

	// router groups
	// adminGroup := e.Group("/admin")
	// jwtGroup := e.Group("/jwt")

	// set all middlewares
	middlewares.SetMainMiddlewares(e)
	middlewares.SetCompleteLogMiddlware(e)

	middlewares.SetAdminMiddlewares(adminGroup)
	middlewares.SetJwtMiddlewares(jwtGroup)

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.AdminGroup(adminGroup)
	api.JwtGroup(jwtGroup)

	return e
}
*/
