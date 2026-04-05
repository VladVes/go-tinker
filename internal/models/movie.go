package models

import (
	"time"
	// "gorm.io/gorm"
)

type Movie struct {
	// gorm.Model
	ID           uint
	Title        string `gorm:"size:150;not null;"`
	Genre        string
	ReleasedAt   time.Time
	Descrtiption string
}
