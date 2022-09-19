package modules

import (
	"github.com/labstack/echo/v4"
)

type Controller struct {
}

func NewController(e *echo.Echo) *echo.Echo {
	c := new(Controller)
	g := e.Group("/modules")
	g.GET("/", c.list)
	return e
}

func (controller Controller) list(c echo.Context) error {
	return c.JSON(200, []string{})
}
