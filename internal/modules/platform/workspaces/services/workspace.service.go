package services

import (
	"oriongo/internal/domain/entities"
	"oriongo/internal/infrastructure"
	"oriongo/internal/infrastructure/repositories"
	"oriongo/internal/modules/platform/workspaces/models/response"
	"strings"
)

type (
	WorkspaceService interface {
		CreateWorkspace(w entities.Workspace) *entities.Workspace
	}
	WorkspaceServiceImpl struct {
		_repo repositories.WorkspaceRepository
	}
)

func NewWorkspaceService(db infrastructure.DbContext) WorkspaceServiceImpl {
	return WorkspaceServiceImpl{
		_repo: repositories.NewWorkspaceRepository(db),
	}
}

func (ws WorkspaceServiceImpl) CreateWorkspace(w entities.Workspace) response.CreateWorkspaceResult {
	workspace, _ := ws._repo.FindByName(strings.ToLower(w.Name))
	if workspace != nil {
		return response.CreateWorkspaceResult{
			Workspace: nil,
			Errors: []string{
				"Workspace already exists",
			},
		}
	}

	success, createError := ws._repo.CreateNew(w)
	if !success {
		return response.CreateWorkspaceResult{
			Workspace: nil,
			Errors:    []string{createError.Error()},
		}
	}

	return response.CreateWorkspaceResult{
		Workspace: &entities.Workspace{},
		Errors:    []string{},
	}
}
