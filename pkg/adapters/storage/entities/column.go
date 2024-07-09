package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Column struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"index"`
	BoardID   uuid.UUID      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Order     uint
}
