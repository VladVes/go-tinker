package models

type Review struct {
	ID    uint
	Text  string
	Score int
	// ----
	MovieID uint
	Movie   Movie `gorm:"foreignKey:MovieID;references:ID"`
}
