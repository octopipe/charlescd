package handler

import "github.com/labstack/echo/v4"

func getNamespace(c echo.Context) string {
	namespace := "default"
	if c.QueryParam("namespace") != "" {
		namespace = c.QueryParam("namespace")
	}

	return namespace
}
