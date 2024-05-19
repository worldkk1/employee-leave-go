package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	Id           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string    `gorm:"type:varchar;not null"`
	IsActive     bool      `gorm:"default:false"`
	NeedApproval bool      `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
