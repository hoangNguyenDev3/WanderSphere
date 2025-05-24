package types

import (
	"time"
)

// Base contains common fields for all models
type Base struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// User represents a user in the system
type User struct {
	Base
	HashedPassword string    `json:"-" gorm:"column:hashed_password;size:1000;not null"`
	Salt           []byte    `json:"-" gorm:"column:salt;size:1000;not null"`
	FirstName      string    `json:"first_name" gorm:"column:first_name;size:50;not null"`
	LastName       string    `json:"last_name" gorm:"column:last_name;size:50;not null"`
	DateOfBirth    time.Time `json:"date_of_birth" gorm:"column:date_of_birth;not null"`
	Email          string    `json:"email" gorm:"column:email;size:100;not null"`
	UserName       string    `json:"user_name" gorm:"column:user_name;size:50;unique;not null"`
	Posts          []*Post   `json:"-" gorm:"foreignKey:UserID"`
	// Followers: Users who follow this user (this user's ID is user_id, followers' IDs are follower_id)
	Followers []*User `json:"-" gorm:"many2many:following;joinForeignKey:user_id;joinReferences:follower_id"`
	// Followings: Users that this user follows (this user's ID is follower_id, followed users' IDs are user_id)
	Followings []*User `json:"-" gorm:"many2many:following;joinForeignKey:follower_id;joinReferences:user_id"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// Following represents a follow relationship between users
type Following struct {
	UserID     int64 `json:"user_id" gorm:"column:user_id;primaryKey"`
	FollowerID int64 `json:"follower_id" gorm:"column:follower_id;primaryKey"`
}

// TableName returns the table name for Following
func (Following) TableName() string {
	return "following"
}

// Post represents a post in the system
type Post struct {
	Base
	UserID           int64      `json:"user_id" gorm:"column:user_id;not null"`
	ContentText      string     `json:"content_text" gorm:"column:content_text;type:text;not null"`
	ContentImagePath string     `json:"content_image_path" gorm:"column:content_image_path;size:1000"`
	User             *User      `json:"-" gorm:"foreignKey:UserID"`
	Comments         []*Comment `json:"-" gorm:"foreignKey:PostID"`
	LikedUsers       []*User    `json:"-" gorm:"many2many:likes;joinForeignKey:post_id;joinReferences:user_id"`
}

// TableName returns the table name for Post
func (Post) TableName() string {
	return "posts"
}

// Comment represents a comment on a post
type Comment struct {
	Base
	PostID      int64  `json:"post_id" gorm:"column:post_id;not null"`
	UserID      int64  `json:"user_id" gorm:"column:user_id;not null"`
	ContentText string `json:"content_text" gorm:"column:content_text;type:text;not null"`
	User        *User  `json:"-" gorm:"foreignKey:UserID"`
	Post        *Post  `json:"-" gorm:"foreignKey:PostID"`
}

// TableName returns the table name for Comment
func (Comment) TableName() string {
	return "comments"
}

// Like represents a like on a post
type Like struct {
	PostID    int64      `json:"post_id" gorm:"column:post_id;primaryKey"`
	UserID    int64      `json:"user_id" gorm:"column:user_id;primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	User      *User      `json:"-" gorm:"foreignKey:UserID"`
	Post      *Post      `json:"-" gorm:"foreignKey:PostID"`
}

// TableName returns the table name for Like
func (Like) TableName() string {
	return "likes"
}
