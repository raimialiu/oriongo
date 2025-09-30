package repositories

import (
	"oriongo/internal/domain/entities"
	"oriongo/internal/infrastructure"

	"gorm.io/gorm"
)

type (
	UserRepository struct {
		_base BaseRepository
	}
)

func NewUserRepository(db infrastructure.DbContext) UserRepository {
	return UserRepository{
		_base: *NewBaseRepository(db),
	}
}

func (repo UserRepository) FindByUsername(username string) (*entities.OrganizationUser, error) {
	var user entities.OrganizationUser
	record := repo._base.DB().Scopes(UserByUsername(username)).First(&user)
	if record.Error != nil {
		return nil, record.Error
	}

	return &user, nil
}

// ===== SCOPES ====

func UserByUsername(username string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("users").Where("username = ?", username)
	}
}

// == END OF SCOPES ==
