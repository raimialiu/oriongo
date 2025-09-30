package repositories

import (
	"errors"
	"oriongo/internal/domain/entities"
	"oriongo/internal/infrastructure"

	"gorm.io/gorm"
)

type (
	WorkspaceRepository struct {
		_base BaseRepository
	}
)

func NewWorkspaceRepository(db infrastructure.DbContext) WorkspaceRepository {
	return WorkspaceRepository{
		_base: *NewBaseRepository(db),
	}
}

func (repo WorkspaceRepository) CreateNew(workspace entities.Workspace) (bool, error) {
	db := repo._base.DB()
	createResponse := *db.Create(&workspace)
	if createResponse.Error != nil {
		return false, createResponse.Error
	}

	return createResponse.RowsAffected > 0, nil
}

/*
func (repo WorkspaceRepository) FindByName(name string) (*entities.Workspace, error) {
	var workspace *entities.Workspace
	_, findError := repo._base.Table("workspaces").FindOne(&workspace, fmt.Sprintf(`name = '%s'`, name), "")
	if findError != nil {
		return nil, findError
	}

	return workspace, nil
}

*/

func (repo WorkspaceRepository) FindByName(name string) (*entities.Workspace, error) {
	var workspace entities.Workspace
	db := repo._base.DB()
	findResult := *db.Where([]map[string]interface{}{{"username": name}}).First(&workspace)
	if error := findResult.Error; error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			return nil, errors.New("workspace not found")
		}
		return nil, error
	}

	return &workspace, nil
}

func (repo WorkspaceRepository) FindByUserKey(organizationKey string) ([]entities.Workspace, error) {
	var workspaces []entities.Workspace
	rows := *repo._base.DB().Scopes(ListWorkspaces(organizationKey)).Find(&workspaces)
	if error := rows.Error; error != nil {
		return nil, error
	}

	return workspaces, nil
}

// SCOPES

func ListWorkspaces(organizationKey string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var workspaces []entities.Workspace
		return db.Model(&entities.Workspace{}).Where("OrganizationKey = ?", organizationKey).Find(&workspaces)
	}
}

// END OF SCOPES
