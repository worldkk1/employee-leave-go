package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"type:varchar;not null"`
	Email     string    `json:"email" gorm:"type:varchar;not null;unique"`
	Password  string    `json:"-" gorm:"type:varchar;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"default:now()"`
}
