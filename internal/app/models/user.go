package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"type:varchar;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
