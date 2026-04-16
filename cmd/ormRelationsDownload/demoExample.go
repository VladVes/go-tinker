package main

import (
	"gorm.io/gorm"

	"github.com/VladVes/go-tinker/v2/internal/models"
)

// Реализуйте функцию GetTasks, которая строит один запрос с учетом Query.
// Используйте Preload для связей Author и Assignee.
// Поля связей могут быть NULL, отчет должен работать и для таких задач.
// Фильтры по связям и статусу могут быть NULL, в этом случае условие не применяется.

// Query defines filters for report generation.
type Query struct {
	AuthorID   *uint
	AssigneeID *uint
	Done       *bool
}

// GetTasks loads tasks with optional filters.
func GetTasks(conn *gorm.DB, q Query) ([]models.Task, error) {
	// BEGIN
	tx := conn.Model(&models.Task{}).Preload("Author").Preload("Assignee")
	if q.AuthorID != nil {
		tx = tx.Where("author_id = ?", *q.AuthorID)
	}
	if q.AssigneeID != nil {
		tx = tx.Where("assignee_id = ?", *q.AssigneeID)
	}
	if q.Done != nil {
		tx = tx.Where("done = ?", *q.Done)
	}

	var tasks []models.Task
	err := tx.Find(&tasks).Error
	return tasks, err
	// END
}
