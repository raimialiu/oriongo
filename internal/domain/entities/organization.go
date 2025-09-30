package entities

type OrganizationUser struct {
	BaseModel
	Username        string `gorm:"uniqueIndex"`
	Password        string
	Email           string
	FirstName       string
	LastName        string
	IsAdmin         bool
	IsAccountOwner  bool   `gorm:"index"`
	ParentKey       string `gorm:"index;column:parent_key"`
	OrganizationKey string `gorm:"index;column:organization_key"`
	//Organization    Organization `gorm:"foreignKey:OrganizationKey;references:EncodedKey"`
}

type Organization struct {
	BaseModel
	Name        string               `gorm:"uniqueIndex"`
	Description string               `gorm:"type:text;null"`
	Users       []OrganizationUser   `gorm:"foreignKey:OrganizationKey;references:EncodedKey;constraint:OnDelete:CASCADE"`
	Workspaces  []Workspace          `gorm:"foreignKey:OrganizationKey;references:EncodedKey;constraint:OnDelete:CASCADE"`
	Settings    OrganizationSettings `gorm:"foreignKey:OrganizationKey;references;EncodedKey;constraint:OnDelete:CASCADE"`
}

type OrganizationSettings struct {
	EnableVersioning bool
	OrganizationKey  string `gorm:"index"`
}
