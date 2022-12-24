package circle

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/customvalidator"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger        *zap.Logger
	circleUseCase CircleUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, circleUseCase CircleUseCase) EchoHandler {
	handler := EchoHandler{
		logger:        logger,
		circleUseCase: circleUseCase,
		validator:     customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces/:workspaceId/circles")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.POST("/:circleId/sync", handler.Sync)
	s.GET("/:circleId/status", handler.Status)
	s.GET("/:circleId", handler.FindById)
	s.PUT("/:circleId", handler.Update)
	s.DELETE("/:circleId", handler.Delete)

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

func (h EchoHandler) bindAndValidateBody(c echo.Context) (Circle, error) {
	body := Circle{}
	if err := c.Bind(&body); err != nil {
		return Circle{}, err
	}

	if err := h.validator.Validate(body); err != nil {
		validateErr := errs.E(errs.Validation, errs.Code("CIRCLE_HTTP_VALIDATIONS"), err)
		return Circle{}, validateErr
	}

	return body, nil
}

func (h EchoHandler) Create(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	body, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	newCircle, err := h.circleUseCase.Create(c.Request().Context(), workspaceId, body)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(201, newCircle)
}

func (h EchoHandler) FindById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
	circle, err := h.circleUseCase.FindById(c.Request().Context(), workspaceId, circleId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circle)
}

func (h EchoHandler) Sync(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
	err := h.circleUseCase.Sync(c.Request().Context(), workspaceId, circleId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, nil)
}

func (h EchoHandler) Update(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")

	body, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	updatedCircle, err := h.circleUseCase.Update(c.Request().Context(), workspaceId, circleId, body)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, updatedCircle)
}

func (h EchoHandler) Delete(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
	err := h.circleUseCase.Delete(c.Request().Context(), workspaceId, circleId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, Circle{})
}

func (h EchoHandler) Status(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
	status, err := h.circleUseCase.Status(c.Request().Context(), workspaceId, circleId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, status)
}
