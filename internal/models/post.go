package models

type Post struct {
	ID     uint
	Title  string
	Body   string
	UserID uint // внешний ключ на users.id
	User   User `gorm:"foreignKey:UserID"` // belongs to: пост принадлежит пользователю

	// ************ many2many - Связь многие ко многим ***********************
	// Связь многие ко многим: many2many

	// Связь многие ко многим возникает, когда сущности могутссылаться
	// друг на друга в обе стороны.
	// Например, у поста может быть несколько тегов, и
	// один и тот же тег встречается у разных постов.
	// В реляционной модели для этого заводят промежуточную таблицу —
	// таблицу связей
	Tags []Tag `gorm:"many2many:post_tags"` // связь многие ко многим через post_tags

	Comments  []Comment
	Published bool
}
