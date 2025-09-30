package entities

import "database/sql"

type Workspace struct {
	BaseModel
	Name            string `gorm:"uniqueIndex"`
	Description     string
	UserKey         string            `gorm:"index"`
	Settings        WorkspaceSettings `gorm:"foreignKey:WorkspaceEncodedKey;references:EncodedKey;constraint:OnDelete:CASCADE"`
	OrganizationKey string            `gorm:"index"`
	Organization    Organization      `gorm:"foreignKey:OrganizationKey;references:EncodedKey;constraint:OnDelete:CASCADE"`
}

type WorkspaceSettings struct {
	EnableVersioning    bool
	EnablePrefix        bool
	PrefixName          sql.NullString `gorm:"type:varchar(10)"`
	WorkspaceEncodedKey string
}
