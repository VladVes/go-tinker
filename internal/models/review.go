package models

import "gorm.io/gorm"

type Review struct {
	ID    uint
	Text  string
	Score int
	// ----
	MovieID uint
	Movie   Movie `gorm:"foreignKey:MovieID;references:ID"`
}

// пример хука
// Рейтинг пересчитывается автоматически после создания отзыва.
// При добавлении нескольких отзывов среднее значение обновляется корректно.

func (r *Review) AfterCreate(tx *gorm.DB) error {
	var movie Movie
	if err := tx.First(&movie, r.MovieID).Error; err != nil {
		return err
	}
	newCount := movie.ReviewCount + 1
	newRating := (movie.Rating*float64(movie.ReviewCount) + float64(r.Score)/float64(newCount))

	return tx.Model(&Movie{}).
		Where("id = ?", r.MovieID).
		Updates(map[string]any{
			"rating":        newRating,
			"reviews_count": newCount,
		}).Error
}
