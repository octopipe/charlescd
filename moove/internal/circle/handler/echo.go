package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/circle"
	"github.com/octopipe/charlescd/moove/internal/core/customvalidator"
)

type EchoHandler struct {
	circleUseCase circle.CircleUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, circleUseCase circle.CircleUseCase) {
	handler := EchoHandler{
		circleUseCase: circleUseCase,
		validator:     customvalidator.NewCustomValidator(),
	}
	s := e.Group("/circles")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.GET("/:name", handler.FindById)
	s.PUT("/:name", handler.Update)
	s.DELETE("/:name", handler.Delete)
	s.GET("/:name/diagram", handler.Diagram)
	s.GET("/:name/resources/:resource", handler.Resource)
	s.GET("/:name/resources/:resource/logs", handler.ResourceLogs)
	s.GET("/:name/resources/:resource/events", handler.ResourceEvents)
}

func (h EchoHandler) FindAll(c echo.Context) error {
	circles, err := h.circleUseCase.FindAll()
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, circles)
}

func (h EchoHandler) Create(c echo.Context) error {
	w := new(circle.Circle)
	if err := c.Bind(w); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Validate(w); err != nil {
		return err
	}

	newCircle, err := h.circleUseCase.Create(*w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(201, newCircle)
}

func (h EchoHandler) FindById(c echo.Context) error {
	return c.JSON(200, circle.Circle{})
}

func (h EchoHandler) Update(c echo.Context) error {
	return c.JSON(200, circle.Circle{})
}

func (h EchoHandler) Delete(c echo.Context) error {
	return c.JSON(204, circle.Circle{})
}

func (h EchoHandler) Diagram(c echo.Context) error {
	circleName := c.Param("name")

	diagram, err := h.circleUseCase.GetDiagram(circleName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(200, diagram)
}

func (h EchoHandler) Resource(c echo.Context) error {
	circleName := c.Param("name")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	resource, err := h.circleUseCase.GetResource(circleName, resourceName, group, kind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(200, resource)
}

func (h EchoHandler) ResourceLogs(c echo.Context) error {
	circleName := c.Param("name")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	logs, err := h.circleUseCase.GetLogs(circleName, resourceName, group, kind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(200, logs)
}

func (h EchoHandler) ResourceEvents(c echo.Context) error {
	circleName := c.Param("name")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	events, err := h.circleUseCase.GetEvents(circleName, resourceName, group, kind)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(200, events)
}
