package handler

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/octopipe/charlescd/internal/moove/core/customvalidator"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger        *zap.Logger
	circleUseCase circle.CircleUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, circleUseCase circle.CircleUseCase) EchoHandler {
	handler := EchoHandler{
		logger:        logger,
		circleUseCase: circleUseCase,
		validator:     customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces/:workspaceId/circles")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.POST("/:circleName/sync", handler.Sync)
	s.GET("/:circleName", handler.FindById)
	s.PUT("/:circleName", handler.Update)
	s.DELETE("/:circleName", handler.Delete)

	return handler
}

func (h EchoHandler) FindAll(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	continueParam := c.QueryParam("continue")
	limitParam := c.QueryParam("limit")

	listOptions := listoptions.Request{
		Limit:    10,
		Continue: continueParam,
	}

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return errs.NewHTTPResponse(c, h.logger, errors.New("limit param invalid"))
		}

		listOptions.Limit = int64(limit)
	}

	circles, err := h.circleUseCase.FindAll(c.Request().Context(), workspaceId, listOptions)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circles)
}

func (h EchoHandler) Create(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	w := circle.Circle{}
	if err := c.Bind(&w); err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	newCircle, err := h.circleUseCase.Create(c.Request().Context(), workspaceId, w)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(201, newCircle)
}

func (h EchoHandler) FindById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	circle, err := h.circleUseCase.FindByName(c.Request().Context(), workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circle)
}

func (h EchoHandler) Sync(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	err := h.circleUseCase.Sync(c.Request().Context(), workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, nil)
}

func (h EchoHandler) Update(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")

	newCircle := circle.Circle{}
	if err := c.Bind(&newCircle); err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	updatedCircle, err := h.circleUseCase.Update(c.Request().Context(), workspaceId, circleName, newCircle)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, updatedCircle)
}

func (h EchoHandler) Delete(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleName := c.Param("circleName")
	err := h.circleUseCase.Delete(c.Request().Context(), workspaceId, circleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, circle.Circle{})
}
