package models

import "time"

type Artist struct {
	BaseModel
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	BirthDate   time.Time `gorm:"type:date"`
	Genre       string    `gorm:"size:100"`
	Active      bool      `gorm:"default:true"`
}
