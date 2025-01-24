package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tigawanna/boxman/system"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
    e.GET("/services", func(c echo.Context) error {
        partialName := c.QueryParam("name")
        services := system.GetSystemDServices(partialName)
        return c.JSON(http.StatusOK, services)
    })
	e.Logger.Fatal(e.Start(":1323"))
}
