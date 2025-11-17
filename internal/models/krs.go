package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KRS struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	MahasiswaID  uuid.UUID  `gorm:"type:uuid" json:"mahasiswa_id"`
	Mahasiswa    User       `gorm:"foreignKey:MahasiswaID" json:"mahasiswa"`
	Semester     string     `gorm:"size:10" json:"semester"`
	Status       string     `gorm:"size:20" json:"status"`
	CatatanDosen string     `gorm:"type:text" json:"catatan_dosen"`
	CreatedAt    time.Time  `gorm:"type:timestamp without time zone" json:"created_at"`
	VerifiedAt   *time.Time `gorm:"type:timestamp without time zone" json:"verified_at"`
	Items        []KRSItem  `gorm:"foreignKey:KRSID" json:"items,omitempty"`
}

func (KRS) TableName() string {
	return "krs"
}

func (k *KRS) BeforeCreate(tx *gorm.DB) (err error) {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}
