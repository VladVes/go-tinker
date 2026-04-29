// package service
package main

import (
	// BEGIN
	"fmt"
	// END
	"testing"

	"github.com/stretchr/testify/require"

	"exercise/internal/config"
	"exercise/internal/db"
	"exercise/internal/models"

	"gorm.io/gorm"
)

// BEGIN
func TestCreatePost(t *testing.T) {
	conn := setupPostDB(t)

	post := postFactory(1)
	require.NoError(t, CreatePost(conn, &post))
	require.NotZero(t, post.ID)
	require.Equal(t, "post-1", post.Slug)
}

func TestUpdatePost(t *testing.T) {
	conn := setupPostDB(t)

	post := postFactory(1)
	require.NoError(t, CreatePost(conn, &post))

	post.Title = "Post 1 Updated"
	require.NoError(t, UpdatePost(conn, &post))
	require.Equal(t, "post-1-updated", post.Slug)
}

func TestPublishPost(t *testing.T) {
	conn := setupPostDB(t)

	post := postFactory(1)
	require.NoError(t, CreatePost(conn, &post))

	require.NoError(t, PublishPost(conn, post.ID))
	var loaded models.Post
	require.NoError(t, conn.First(&loaded, post.ID).Error)
	require.True(t, loaded.Published)
}

func TestDeletePost(t *testing.T) {
	conn := setupPostDB(t)

	post := postFactory(1)
	require.NoError(t, CreatePost(conn, &post))

	require.NoError(t, DeletePost(conn, post.ID))
	var count int64
	require.NoError(t, conn.Model(&models.Post{}).Count(&count).Error)
	require.Equal(t, int64(0), count)
}

func TestCreatePostValidation(t *testing.T) {
	conn := setupPostDB(t)

	badTitle := postFactory(1)
	badTitle.Title = ""
	require.ErrorIs(t, CreatePost(conn, &badTitle), errEmptyTitle)

	badBody := postFactory(1)
	badBody.Body = ""
	require.ErrorIs(t, CreatePost(conn, &badBody), errEmptyBody)
}

func TestGetPost(t *testing.T) {
	conn := setupPostDB(t)

	post := postFactory(1)
	require.NoError(t, CreatePost(conn, &post))

	loaded, err := GetPost(conn, post.ID)
	require.NoError(t, err)
	require.Equal(t, post.Title, loaded.Title)
}

func TestListPosts(t *testing.T) {
	conn := setupPostDB(t)

	first := postFactory(1)
	second := postFactory(2)
	require.NoError(t, CreatePost(conn, &first))
	require.NoError(t, CreatePost(conn, &second))

	items, err := ListPosts(conn)
	require.NoError(t, err)
	require.Len(t, items, 2)
}

func postFactory(i int) models.Post {
	return models.Post{
		Title:     fmt.Sprintf("Post %d", i),
		Body:      fmt.Sprintf("Body %d", i),
		Published: false,
	}
}

// END

func setupPostDB(t *testing.T) *gorm.DB {
	cfg := config.FromDefaults()
	conn, err := db.Open(cfg)
	require.NoError(t, err)
	require.NoError(t, conn.AutoMigrate(&models.Post{}))

	tx := conn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })

	return tx
}
