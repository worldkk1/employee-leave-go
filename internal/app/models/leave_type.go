package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	Id              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name            string    `gorm:"type:varchar;not null"`
	IsActive        bool      `gorm:"default:false"`
	NeedApproval    bool      `gorm:"default:false"`
	TotalLeaves     int       `gorm:"not null"`
	IsRatioAllocate bool      `gorm:"default:false"`
	IsCarryOver     bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"default:now()"`
	UpdatedAt       time.Time `gorm:"default:now()"`
}
