package main

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a blog post.
type TestPost struct {
	// BEGIN (write your solution here)
	ID            uint      `gorm:"index:idx_posts_created_at_id,priority:2"`
	UserID        uint      `gorm:"not null"`
	Title         string    `gorm:"type:VARCHAR(255);not null"`
	Body          string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"index:idx_posts_created_at_id,priority:1"`
	ShortBody     string    `gorm:"type:VARCHAR(103);not null"`
	CommentsCount int64     `gorm:"default:0"`
	Author        User      `gorm:"foreignKey:UserID"`
	// END
}

type PostListItem struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	ShortBody     string    `json:"short_body"`
	AuthorName    string    `json:"author_name"`
	CreatedAt     time.Time `json:"created_at"`
	CommentsCount int64     `json:"comments_count"`
}

// ListPosts returns a list of posts for the feed.
func ListPosts(conn *gorm.DB) ([]PostListItem, error) {
	// BEGIN (write your solution here)
	items := []PostListItem{}
	query := conn.Table("posts").
		Select("posts.id, posts.title, posts.short_body AS short_body, users.name AS author_name, posts.created_at, posts.comments_count AS comments_count").
		Joins("JOIN users ON users.id = posts.user_id").
		Order("posts.created_at DESC, posts.id DESC")

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
	// END
}
