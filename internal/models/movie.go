package models

import (
	"time"
)

type Movie struct {
	// gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:150;not null;"`
	Genre       string
	ReleasedAt  time.Time
	Description string
	// для теста db.Automigrate
	AdditionalInfo string

	// Тип данных в Go: float64 подходит для хранения чисел с дробной частью.
	// Соответствует типу NUMERIC в PostgreSQL по семантике (хранение дробных чисел).
	// Обеспечивает достаточную точность для формата NUMERIC(3,1) (диапазон от −99.9 до 99.9).
	// Тег GORM: gorm:"type:numeric(3,1)"
	// Директива type: явно задаёт тип столбца в БД.
	// numeric(3,1) — прямое указание создать столбец типа NUMERIC с:
	// precision = 3 — общее количество цифр (включая целую и дробную части);
	// scale = 1 — количество цифр после запятой.
	Rating float64 `gorm:"type:numeric(3,1)"`

	DirectorID uint
	Director   *Director `gorm:"foreignKey:DirectorID;references:ID"`
	// пример: `gorm:"many2many:post_tags"`
	Actor       []Actor `gorm:"many2many:movie_actor"`
	Review      []Review
	ReviewCount int
}
