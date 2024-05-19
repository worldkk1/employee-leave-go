package models

import (
	"time"

	"github.com/google/uuid"
)

type UserLeaveRecord struct {
	Id            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId        uuid.UUID `gorm:"type:uuid;not null"`
	User          User      `gorm:"foreignKey:UserId"`
	LeaveTypesId  uuid.UUID `gorm:"type:uuid;not null"`
	LeaveType     LeaveType `gorm:"foreignKey:LeaveTypesId"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
	TotalLeaveDay int       `gorm:"not null"`
	Reason        *string
	AttachmentURL *string
	ApproveBy     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
