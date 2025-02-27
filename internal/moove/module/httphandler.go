package module

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/customvalidator"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger        *zap.Logger
	moduleUseCase ModuleUseCase
	validator     customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, moduleUseCase ModuleUseCase) {
	handler := EchoHandler{
		logger:        logger,
		moduleUseCase: moduleUseCase,
		validator:     customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces/:workspaceId/modules")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.GET("/:moduleName", handler.FindById)
	s.PUT("/:moduleName", handler.Update)
	s.DELETE("/:moduleName", handler.Delete)
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

	modules, err := h.moduleUseCase.FindAll(c.Request().Context(), workspaceId, listOptions)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, modules)
}

func (h EchoHandler) bindAndValidateBody(c echo.Context) (Module, error) {
	body := Module{}
	if err := c.Bind(&body); err != nil {
		return Module{}, err
	}

	if err := h.validator.Validate(body); err != nil {
		validateErr := errs.E(errs.Validation, errs.Code("CIRCLE_HTTP_VALIDATIONS"), err)
		return Module{}, validateErr
	}

	fmt.Println(body)
	if body.Visibility == PrivateModule && body.Auth == nil && body.SecretRef == nil {
		validateErr := errs.E(errs.Validation, errs.Code("CIRCLE_HTTP_VALIDATIONS"), errors.New("private repository without authentication"))
		return Module{}, validateErr
	}

	return body, nil
}

func (h EchoHandler) Create(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	body, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	newModule, err := h.moduleUseCase.Create(c.Request().Context(), workspaceId, body)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(201, newModule)
}

func (h EchoHandler) FindById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	moduleName := c.Param("moduleName")
	module, err := h.moduleUseCase.FindById(c.Request().Context(), workspaceId, moduleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, module)
}

func (h EchoHandler) Update(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	moduleName := c.Param("moduleName")

	body, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	updatedModule, err := h.moduleUseCase.Update(c.Request().Context(), workspaceId, moduleName, body)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, updatedModule)
}

func (h EchoHandler) Delete(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	moduleName := c.Param("moduleName")
	err := h.moduleUseCase.Delete(c.Request().Context(), workspaceId, moduleName)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(204, Module{})
}
