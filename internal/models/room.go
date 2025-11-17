package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nama    string    `gorm:"size:50" json:"nama"`
	Classes []Class   `gorm:"foreignKey:RoomID" json:"classes,omitempty"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
