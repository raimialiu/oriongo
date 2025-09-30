package infrastructure

import (
	"errors"
	"fmt"
	"oriongo/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type (
	MigrationFunc func(db *gorm.DB) error
	MigrationStep struct {
		Name string
		Up   MigrationFunc
		Down MigrationFunc
	}

	GormMigrator struct {
		_steps []MigrationStep
	}
)

func NewMigrator() *GormMigrator {
	return &GormMigrator{
		_steps: []MigrationStep{},
	}
}

func (m *GormMigrator) AddStep(step MigrationStep) *GormMigrator {
	m._steps = append(m._steps, step)
	return m
}

func (m *GormMigrator) AddSteps(steps ...MigrationStep) *GormMigrator {
	m._steps = append(m._steps, steps...)
	return m
}

func (m *GormMigrator) Up(db *gorm.DB) error {
	timestamp := time.Now().Unix()
	migrations := make([]entities.Migration, len(m._steps))
	tx := db.Begin()

	for _, step := range m._steps {
		var migration entities.Migration
		tx := tx.Where("name = ?", step.Name).First(&migration)
		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			println(fmt.Sprintf("migration %s is already applied, skipping....", step.Name))
			continue
		}

		if mError := step.Up(tx); mError != nil {
			fmt.Println(mError)
			tx.Rollback()
		}

		migrations = append(migrations, entities.Migration{
			Name:      step.Name,
			Timestamp: timestamp,
		})
	}

	tx.Create(migrations)
	tx.Commit()
	return nil
}
