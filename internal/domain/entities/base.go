package entities

import (
	"oriongo/internal/common/constants"
	"time"
)

type BaseModel struct {
	EncodedKey     string                 `gorm:"primary_key;default:uuid_generate_v4()"`
	CreationDate   time.Time              `gorm:"autoCreateTime"`
	CreatedBy      string                 `gorm:"index"`
	LastModified   *time.Time             `gorm:"autoUpdateTime"`
	LastModifiedBy string                 `gorm:"index"`
	Status         constants.EntityStatus `gorm:"index"`
}
