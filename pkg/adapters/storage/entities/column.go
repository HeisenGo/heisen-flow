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
	BoardID   uuid.UUID      `gorm:"index:idx_together_order_board_id,unique"`
	Board     Board          `gorm:"foreignKey:BoardID"`
	Order     uint           `gorm:"index:idx_together_order_board_id,unique"`
}
