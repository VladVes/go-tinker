package models

type Actor struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
