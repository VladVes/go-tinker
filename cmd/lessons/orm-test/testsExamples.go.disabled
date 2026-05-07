package main

import (
	"testing"

	"gorm.io/gorm"

	"github.com/VladVes/go-tinker/v2/internal/models"
	"github.com/stretchr/testify/require"
)

// Пример кода для тестирования (CRUD‑операции)
type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) FindByID(id uint) (models.Movie, error) {
	var movie models.Movie
	if err := r.db.First(&movie, id).Error; err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

func (r *MovieRepository) UpdateTitle(id uint, title string) error {
	return r.db.Model(&models.Movie{}).
		Where("id = ?", id).
		Update("title", title).Error
}

func (r *MovieRepository) Delete(id uint) error {
	return r.db.Delete(&models.Movie{}, id).Error
}

// Пример тестов:
func movieFactory(title string) *models.Movie {
	return &models.Movie{
		Title: title,
		Genre: "sci-fi",
	}
}

func TestMovieRepository_CreateAndFind(t *testing.T) {
	db := testDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	repo := NewMovieRepository(tx)
	movie := movieFactory("Inception")
	require.NoError(t, repo.Create(movie))

	loaded, err := repo.FindByID(movie.ID)
	require.NoError(t, err)
	require.Equal(t, movie.Title, loaded.Title)
}

func TestMovieRepository_UpdateTitle(t *testing.T) {
	db := testDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	repo := NewMovieRepository(tx)
	movie := movieFactory("Old title")
	require.NoError(t, repo.Create(movie))

	require.NoError(t, repo.UpdateTitle(movie.ID, "New title"))
	loaded, err := repo.FindByID(movie.ID)
	require.NoError(t, err)
	require.Equal(t, "New title", loaded.Title)
}
