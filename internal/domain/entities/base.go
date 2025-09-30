package entities

import (
	"oriongo/internal/common/constants"
	"time"
)

type BaseModel struct {
	EncodedKey     string                 `gorm:"primaryKey;default:(UUID())"`
	CreationDate   time.Time              `gorm:"autoCreateTime"`
	CreatedBy      string                 `gorm:"index"`
	LastModified   *time.Time             `gorm:"column:last_modified;autoUpdateTime"`
	LastModifiedBy string                 `gorm:"index"`
	Status         constants.EntityStatus `gorm:"index"`
}
