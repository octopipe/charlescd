package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/circle"
	"github.com/octopipe/charlescd/moove/internal/core/customvalidator"
	"github.com/octopipe/charlescd/moove/internal/errs"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger        *zap.Logger
	circleUseCase circle.CircleUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, circleUseCase circle.CircleUseCase) {
	handler := EchoHandler{
		logger:        logger,
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
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circles)
}

func (h EchoHandler) Create(c echo.Context) error {
	w := new(pbv1.CreateCircleRequest)
	if err := c.Bind(w); err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	newCircle, err := h.circleUseCase.Create(w)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(201, newCircle)
}

func (h EchoHandler) FindById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	circle, err := h.circleUseCase.FindByName(workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circle)
}

func (h EchoHandler) Update(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")

	newCircle := new(pbv1.CreateCircleRequest)
	if err := c.Bind(newCircle); err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	err := h.circleUseCase.Update(workspaceId, circleName, newCircle)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, nil)
}

func (h EchoHandler) Delete(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	err := h.circleUseCase.Delete(workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, circle.Circle{})
}

func (h EchoHandler) Diagram(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")

	diagram, err := h.circleUseCase.GetDiagram(workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, diagram)
}

func (h EchoHandler) Resource(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	resource, manifest, err := h.circleUseCase.GetResource(workspaceId, resourceName, group, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, map[string]interface{}{
		"metadata": resource,
		"manifest": manifest,
	})
}

func (h EchoHandler) ResourceLogs(c echo.Context) error {
	circleName := c.Param("name")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	logs, err := h.circleUseCase.GetLogs(circleName, resourceName, group, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, logs)
}

func (h EchoHandler) ResourceEvents(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	resourceName := c.Param("resource")
	kind := c.QueryParam("kind")

	events, err := h.circleUseCase.GetEvents(workspaceId, resourceName, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(200, events)
}
