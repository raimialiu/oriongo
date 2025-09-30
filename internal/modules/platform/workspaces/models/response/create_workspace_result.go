package response

import "oriongo/internal/domain/entities"

type CreateWorkspaceResult struct {
	Workspace *entities.Workspace
	Errors    []string
}
