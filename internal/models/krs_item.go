package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KRSItem struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	KRSID           uuid.UUID  `gorm:"type:uuid" json:"krs_id"`
	KRS             KRS        `gorm:"foreignKey:KRSID" json:"-"`
	ClassID         uuid.UUID  `gorm:"type:uuid" json:"class_id"`
	Class           Class      `gorm:"foreignKey:ClassID" json:"class"`
	Status          string     `gorm:"size:20" json:"status"`
	CreatedAt       time.Time  `gorm:"type:timestamp without time zone" json:"created_at"`
	DiajukanBatalAt *time.Time `gorm:"type:timestamp without time zone" json:"diajukan_batal_at"`
	DibatalkanAt    *time.Time `gorm:"type:timestamp without time zone" json:"dibatalkan_at"`
}

func (KRSItem) TableName() string {
	return "krs_items"
}

func (ki *KRSItem) BeforeCreate(tx *gorm.DB) (err error) {
	if ki.ID == uuid.Nil {
		ki.ID = uuid.New()
	}
	return nil
}
