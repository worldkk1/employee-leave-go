package models

import (
	"time"

	"github.com/google/uuid"
)

type UserLeave struct {
	Id           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId       uuid.UUID `gorm:"type:uuid;not null"`
	User         User      `gorm:"foreignKey:UserId"`
	LeaveTypesId uuid.UUID `gorm:"type:uuid;not null"`
	LeaveType    LeaveType `gorm:"foreignKey:LeaveTypesId"`
	Remaining    int       `gorm:"not null"`
	Used         int       `gorm:"not null"`
	StartDate    time.Time
	EndDate      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
