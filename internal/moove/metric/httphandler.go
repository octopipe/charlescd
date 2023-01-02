package metric

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
	metricUseCase MetricUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, metricUseCase MetricUseCase) EchoHandler {
	handler := EchoHandler{
		logger:        logger,
		metricUseCase: metricUseCase,
		validator:     customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces/:workspaceId/circles/:circleId/metrics")
	s.GET("", handler.FindAll)
	s.GET("/:metricId/query", handler.Query)

	return handler
}

func (h EchoHandler) FindAll(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
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

	circles, err := h.metricUseCase.FindAll(c.Request().Context(), workspaceId, circleId, listOptions)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circles)
}

func (h EchoHandler) Query(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	circleId := c.Param("circleId")
	metricId := c.Param("metricId")
	rangeTime := c.QueryParam("range")

	if rangeTime == "" {
		rangeTime = ThirdyMinutes
	}

	circles, err := h.metricUseCase.Query(c.Request().Context(), workspaceId, circleId, metricId, rangeTime)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, circles)
}
