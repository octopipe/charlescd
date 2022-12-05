package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/customvalidator"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"github.com/octopipe/charlescd/internal/moove/resource"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger          *zap.Logger
	resourceUseCase resource.ResourceUseCase
	validator       customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, resourceUseCase resource.ResourceUseCase) {
	handler := EchoHandler{
		logger:          logger,
		resourceUseCase: resourceUseCase,
		validator:       customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces/:workspaceId/circles/:circleName/resources")
	s.GET("/tree", handler.Tree)
	s.GET("/:resource", handler.Resource)
	s.GET("/:resource/manifest", handler.ResourceManifest)
	s.GET("/:resource/logs", handler.ResourceLogs)
	s.GET("/:resource/events", handler.ResourceEvents)
}

func (h EchoHandler) Tree(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")

	diagram, err := h.resourceUseCase.GetTree(c.Request().Context(), workspaceId, circleName)
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

	resource, err := h.resourceUseCase.GetResource(c.Request().Context(), workspaceId, resourceName, group, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, map[string]interface{}{
		"metadata": resource,
	})
}

func (h EchoHandler) ResourceManifest(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	manifest, err := h.resourceUseCase.GetManifest(c.Request().Context(), workspaceId, resourceName, group, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, map[string]interface{}{
		"content": manifest,
	})
}

func (h EchoHandler) ResourceLogs(c echo.Context) error {
	circleName := c.Param("name")
	resourceName := c.Param("resource")
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")

	logs, err := h.resourceUseCase.GetLogs(c.Request().Context(), circleName, resourceName, group, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, logs)
}

func (h EchoHandler) ResourceEvents(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	resourceName := c.Param("resource")
	kind := c.QueryParam("kind")

	events, err := h.resourceUseCase.GetEvents(c.Request().Context(), workspaceId, resourceName, kind)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(200, events)
}
