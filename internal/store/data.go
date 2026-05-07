package store

import (
	"github.com/VladVes/go-tinker/v2/internal/store/entities"
	"github.com/VladVes/go-tinker/v2/internal/store/schemas"
)

var (
	PostLikes = map[string]int64{}
	Logs      []schemas.LogEntry
	Orders    = make(map[string]entities.Order)
	Employee  = make(map[string]entities.Employee)
	Users     = make(map[string]entities.User)
)

// ==================================== HTML Templates ================================================
// Для примера html шаблонизации списка
// Структура с информацией о фильме для примера шаблонизации co списком
type Film struct {
	Title    string
	IsViewed bool
}

// Для простоты описываем хранилище фильмов в коде
var ListExample = []Film{
	{
		Title:    "The Shawshank Redemption",
		IsViewed: true,
	},
	{
		Title:    "The Godfather",
		IsViewed: true,
	},
	{
		Title:    "The Godfather: Part II",
		IsViewed: false,
	},
}

//=====================================================================================================
