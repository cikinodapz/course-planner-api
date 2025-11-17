package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Class struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID          uuid.UUID `gorm:"type:uuid" json:"course_id"`
	Course            Course    `gorm:"foreignKey:CourseID" json:"course"`
	DosenID           uuid.UUID `gorm:"type:uuid" json:"dosen_id"`
	Dosen             User      `gorm:"foreignKey:DosenID" json:"dosen"`
	NamaKelas         string    `gorm:"size:10" json:"nama_kelas"`
	Hari              string    `gorm:"size:10" json:"hari"`
	JamMulai          time.Time `gorm:"type:timestamp without time zone" json:"jam_mulai"`
	JamSelesai        time.Time `gorm:"type:timestamp without time zone" json:"jam_selesai"`
	RoomID            uuid.UUID `gorm:"type:uuid" json:"room_id"`
	Room              Room      `gorm:"foreignKey:RoomID" json:"room"`
	Kuota             int       `json:"kuota"`
	SemesterPenawaran string    `gorm:"size:10" json:"semester_penawaran"`
	KRSItems          []KRSItem `gorm:"foreignKey:ClassID" json:"krs_items,omitempty"`
}

func (c *Class) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
