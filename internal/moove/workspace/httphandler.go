package workspace

import (
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/customvalidator"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"go.uber.org/zap"
)

type EchoHandler struct {
	logger           *zap.Logger
	workspaceUseCase WorkspaceUseCase
	validator        customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, logger *zap.Logger, workspaceUseCase WorkspaceUseCase) EchoHandler {
	handler := EchoHandler{
		logger:           logger,
		workspaceUseCase: workspaceUseCase,
		validator:        customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.GET("/:workspaceId", handler.FindById)
	s.PUT("/:workspaceId", handler.Update)
	s.DELETE("/:workspaceId", handler.Delete)

	return handler
}

func (h EchoHandler) FindAll(c echo.Context) error {
	workspaces, err := h.workspaceUseCase.FindAll()
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, workspaces)
}

func (h EchoHandler) bindAndValidateBody(c echo.Context) (Workspace, error) {
	body := Workspace{}
	if err := c.Bind(&body); err != nil {
		return Workspace{}, err
	}

	if err := h.validator.Validate(body); err != nil {
		validateErr := errs.E(errs.Validation, errs.Code("WORKSPACE_HTTP_VALIDATIONS"), err)
		return Workspace{}, validateErr
	}

	return body, nil
}

func (h EchoHandler) Create(c echo.Context) error {
	w, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	newWorkspace, err := h.workspaceUseCase.Create(w)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(201, newWorkspace)
}

func (h EchoHandler) FindById(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	workspace, err := h.workspaceUseCase.FindById(workspaceId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, workspace)
}

func (h EchoHandler) Update(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	body, err := h.bindAndValidateBody(c)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	workspace, err := h.workspaceUseCase.Update(workspaceId, body)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}
	return c.JSON(200, workspace)
}

func (h EchoHandler) Delete(c echo.Context) error {
	workspaceId := c.Param("workspaceId")
	err := h.workspaceUseCase.Delete(workspaceId)
	if err != nil {
		return errs.NewHTTPResponse(c, h.logger, err)
	}

	return c.JSON(204, nil)
}
