package repository

import (
	"github.com/tedbearr/react-go/entity"
	"gorm.io/gorm"
)

type PostsRepository interface {
	InsertPosts(post entity.Posts) entity.Posts
	UpdatePosts(post entity.Posts) entity.Posts
	DeletePosts(post entity.Posts)
	AllPosts() []entity.Posts
	FindPostsById(postsID uint64) entity.Posts
}

type postsConnection struct {
	connection *gorm.DB
}

func NewPostsRepository(dbConnection *gorm.DB) PostsRepository {
	return &postsConnection{
		connection: dbConnection,
	}
}

func (db *postsConnection) InsertPosts(post entity.Posts) entity.Posts {
	db.connection.Save(&post)
	db.connection.Preload("User").Find(&post)
	return post
}

func (db *postsConnection) UpdatePosts(post entity.Posts) entity.Posts {
	db.connection.Save(&post)
	db.connection.Preload("User").Find(&post)
	return post
}

func (db *postsConnection) DeletePosts(post entity.Posts) {
	db.connection.Delete(&post)
}

func (db *postsConnection) FindPostsById(postsID uint64) entity.Posts {
	var post entity.Posts
	db.connection.Preload("User").Find(&post, postsID)
	return post
}

func (db *postsConnection) AllPosts() []entity.Posts {
	var posts []entity.Posts
	db.connection.Preload("User").Find(&posts)
	return posts
}
