package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string     `gorm:"size:100" json:"name"`
	Email      string     `gorm:"size:100;uniqueIndex" json:"email"`
	Password   string     `json:"-"`
	Role       string     `gorm:"size:20" json:"role"`
	NIM        string     `gorm:"size:20" json:"nim"`
	NIDN       string     `gorm:"size:20" json:"nidn"`
	DosenPAID  *uuid.UUID `gorm:"type:uuid" json:"dosen_pa_id"`
	CreatedAt  time.Time  `gorm:"type:timestamp without time zone" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"type:timestamp without time zone" json:"updated_at"`
}

// Auto-generate UUID sebelum insert
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
