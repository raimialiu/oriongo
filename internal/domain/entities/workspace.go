package entities

import "database/sql"

type Workspace struct {
	BaseModel
	Name        string
	Description sql.NullString `gorm:"type:text"`
}

type WorkspaceSettings struct {
	SupportVersioning bool
	EnableVersioning  bool
	EnablePrefix      bool
	PrefixName        sql.NullString `gorm:"type:varchar(10)"`
}
