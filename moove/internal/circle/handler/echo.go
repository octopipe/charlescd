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
	s := e.Group("/workspaces/:workspaceId/circles")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.GET("/:circleName", handler.FindById)
	s.PUT("/:circleName", handler.Update)
	s.DELETE("/:circleName", handler.Delete)
	s.GET("/:circleName/diagram", handler.Diagram)
	s.GET("/:circleName/resources/:resource", handler.Resource)
	s.GET("/:circleName/resources/:resource/logs", handler.ResourceLogs)
	s.GET("/:circleName/resources/:resource/events", handler.ResourceEvents)
}

func (h EchoHandler) FindAll(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circles, err := h.circleUseCase.FindAll(workspaceId)
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
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	circle, err := h.circleUseCase.FindByName(workspaceId, circleName)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, circle)
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
