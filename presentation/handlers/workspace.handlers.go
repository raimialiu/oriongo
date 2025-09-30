package handlers

import (
	"net/http"
	"oriongo/internal/domain/entities"
	"oriongo/internal/infrastructure"
	"oriongo/internal/modules/platform/workspaces/models/request"
	"oriongo/internal/modules/platform/workspaces/services"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type WorkspaceHandler struct {
	_base BaseHandler
	_svc  services.WorkspaceServiceImpl
}

func NewWorkspaceHandler(db infrastructure.DbContext) *WorkspaceHandler {

	return &WorkspaceHandler{
		_base: *NewBaseHandler(db),
		_svc:  services.NewWorkspaceService(db),
	}
}

func (h *WorkspaceHandler) CreateWorkspace(e echo.Context) error {
	var payload = &request.CreateWorkspaceRequest{}
	if err := e.Bind(payload); err != nil {
		return e.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	validationError := validation.ValidateStruct(payload,
		validation.Field(&payload.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&payload.Description, validation.Max(100)),
	)

	if validationError != nil {
		return e.JSON(http.StatusBadRequest, ErrorResponse(validationError.Error()))
	}

	response := h._svc.CreateWorkspace(entities.Workspace{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if len(response.Errors) > 0 {
		return e.JSON(http.StatusBadRequest, ErrorResponse(response.Errors[0]))
	}

	return e.JSON(http.StatusCreated, SuccessResponse(response.Workspace, "Workspace created successfully"))
}
