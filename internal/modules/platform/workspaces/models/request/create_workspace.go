package request

type CreateWorkspaceRequest struct {
	Name        string `json: "name"`
	Description string `json: "description"`
}
