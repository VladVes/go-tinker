package models

type Comment struct {
	ID       uint
	Body     string
	Approved bool
	PostID   uint
	UserID   uint // внешний ключ на users.id
	User     User `gorm:"foreignKey:UserID"` // belongs to: пост принадлежит пользователю
}
