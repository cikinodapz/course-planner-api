package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Kode    string    `gorm:"size:20" json:"kode"`
	Nama    string    `gorm:"size:100" json:"nama"`
	SKS     int       `json:"sks"`
	Classes []Class   `gorm:"foreignKey:CourseID" json:"classes,omitempty"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
