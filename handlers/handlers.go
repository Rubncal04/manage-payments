package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloWorld muestra un mensaje simple.
func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
