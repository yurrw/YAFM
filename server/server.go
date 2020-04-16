package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	//Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//Mudar o codigo abaixo pra interface

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
	// Server

	// Setup route group for the API
	api := e.Group("/api")
	{
		api.GET("/", func(c *echo.Context) {
			c.JSON(http.StatusOK,{
				"message": "pong",
			})
		})
	}
	func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
		  return
		}
		return c.JSON(http.StatusOK, u)
	  }
	//Server
	e.Run(standard.New(":7777"))
}
