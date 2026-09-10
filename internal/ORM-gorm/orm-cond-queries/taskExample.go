package main

import (
	"strings"

	"gorm.io/gorm"

	"github.com/VladVes/go-tinker/v2/internal/models"
)

// BEGIN
// Query defines filters and ordering for tasks.
type Query struct {
	Prefix string
	Done   *bool
	Sort   string
}

// FindTasks returns tasks matching the query parameters.
func FindTasks(conn *gorm.DB, q Query) ([]models.Task, error) {
	tx := conn.Model(&models.Task{})
	if q.Prefix != "" {
		tx = tx.Where("title LIKE ?", q.Prefix+"%")
	}
	if q.Done != nil {
		tx = tx.Where("done = ?", *q.Done)
	}

	column, direction := parseSort(q.Sort)
	tx = tx.Order(column + " " + direction)

	var tasks []models.Task
	err := tx.Find(&tasks).Error
	return tasks, err
}

func parseSort(sort string) (string, string) {
	column := "id"
	direction := "DESC"
	parts := strings.Split(strings.ToLower(sort), ",")
	if len(parts) > 0 && parts[0] != "" {
		switch parts[0] {
		case "id", "title":
			column = parts[0]
		}
	}
	if len(parts) > 1 && parts[1] != "" {
		switch parts[1] {
		case "asc":
			direction = "ASC"
		case "desc":
			direction = "DESC"
		}
	}
	return column, direction
}

// END
