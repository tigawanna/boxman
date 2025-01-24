package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tigawanna/boxman/systemd"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/services", func(c echo.Context) error {
		partialName := c.QueryParam("name")
		services := systemd.GetSystemDServices(partialName)
		return c.JSON(http.StatusOK, services)
	})
	e.GET("/new", func(c echo.Context) error {
		config := systemd.NewSystemdServiceConfig(
			"pocketbase",
			"~/pb",
			"pocketbase serve yourdomain.com",
			&systemd.ConfigOptions{
				User:  "pocketbase",
				Group: "pocketbase",
			},
		)
		// fmt.Println(config.ToString())
		return c.String(http.StatusOK, config.ToString())
	})
	e.POST("/service/new", func(c echo.Context) error {
		name := c.FormValue("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "name is required")
		}
		path := c.FormValue("path")
		if path == "" {
			return c.String(http.StatusBadRequest, "path is required")
		}
		path = strings.TrimSpace(path)

		if len(path) >= 2 && path[:2] != "~/" {
			return c.String(http.StatusBadRequest, "path must be absolute, try ~/path/to/service")
		}
		config := systemd.NewSystemdServiceConfig(
			"pocketbase",
			"~/pb",
			"pocketbase serve yourdomain.com",
			&systemd.ConfigOptions{
				User:  "pocketbase",
				Group: "pocketbase",
			},
		)
		// fmt.Println(config.ToString())
		return c.String(http.StatusOK, config.ToString())
	})
	e.Logger.Fatal(e.Start(":1323"))
}
