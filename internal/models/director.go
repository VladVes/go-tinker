package models

type Director struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Movie []Movie
}
