package relationsexample

// У задачи должны быть два пользователя: исполнитель и заказчик.
// Используйте внешние ключи и поля связей, чтобы GORM корректно мигрировал и загружал данные.
// В модели пользователя добавьте списки задач для ролей исполнителя и заказчика.

// User represents a user entity.
type User struct {
	ID   uint
	Name string
	// BEGIN
	TasksAssigned  []Task `gorm:"foreignKey:AssigneeID"`
	TasksRequested []Task `gorm:"foreignKey:AuthorID"`
	// END
}

type Task struct {
	ID    uint
	Title string
	Done  bool
	// BEGIN
	AuthorID   uint
	AssigneeID uint
	Assignee   User `gorm:"foreignKey:AssigneeID"`
	Author     User `gorm:"foreignKey:AuthorID"`
	// END
}
