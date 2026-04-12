package models

// ************** Настройка внешних ключей и поведения при удалении **********************************
// При нескольких внешних ключах в одной структуре каждый настраивается отдельно:

type Payment struct {
	ID       uint
	BuyerID  uint
	SellerID uint

	Buyer  User `gorm:"foreignKey:BuyerID;constraint:OnDelete:CASCADE"`   // удаляем покупателя — удаляются платежи
	Seller User `gorm:"foreignKey:SellerID;constraint:OnDelete:SET NULL"` // удаляем продавца — ссылка обнуляется
}
