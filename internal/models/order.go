package models

type Order struct {
	ID     uint
	Number string
	UserID uint
	// ************** Настройка внешних ключей и поведения при удалении **********************************
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
