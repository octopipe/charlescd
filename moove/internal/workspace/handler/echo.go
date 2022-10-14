package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/core/customvalidator"
	"github.com/octopipe/charlescd/moove/internal/workspace"
)

type EchoHandler struct {
	workspaceUseCase workspace.WorkspaceUseCase
	validator        customvalidator.CustomValidator
}

func NewEchohandler(e *echo.Echo, workspaceUseCase workspace.WorkspaceUseCase) {
	handler := EchoHandler{
		workspaceUseCase: workspaceUseCase,
		validator:        customvalidator.NewCustomValidator(),
	}
	s := e.Group("/workspaces")
	s.GET("", handler.FindAll)
	s.POST("", handler.Create)
	s.GET("/:id", handler.FindById)
	s.PUT("/:id", handler.Update)
	s.DELETE("/:id", handler.Delete)
}

func (h EchoHandler) FindAll(c echo.Context) error {
	workspaces, err := h.workspaceUseCase.FindAll()
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, workspaces)
}

func (h EchoHandler) Create(c echo.Context) error {
	w := new(workspace.Workspace)
	if err := c.Bind(w); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Validate(w); err != nil {
		return err
	}

	newWorkspace, err := h.workspaceUseCase.Create(*w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(201, newWorkspace)
}

func (h EchoHandler) FindById(c echo.Context) error {
	return c.JSON(200, workspace.Workspace{})
}

func (h EchoHandler) Update(c echo.Context) error {
	return c.JSON(200, workspace.Workspace{})
}

func (h EchoHandler) Delete(c echo.Context) error {
	return c.JSON(204, workspace.Workspace{})
}
