package models

import "time"

type User struct {
	ID             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	HashedPassword string    `gorm:"column:hashed_password;type:varchar(1000);not null"`
	Salt           []byte    `gorm:"column:salt;type:varbinary(1000);not null"`
	FirstName      string    `gorm:"column:first_name;type:varchar(50);not null"`
	LastName       string    `gorm:"column:last_name;type:varchar(50);not null"`
	DOB            time.Time `gorm:"column:dob;type:date;not null"`
	Email          string    `gorm:"column:email;type:varchar(100);not null"`
	UserName       string    `gorm:"column:user_name;type:varchar(50);unique;not null"`
}

func (User) TableName() string {
	return "user"
}

type Post struct {
	ID               int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	UserID           int64     `gorm:"column:user_id;type:bigint;not null"`
	ContentText      string    `gorm:"column:content_text;type:text(100000);not null"`
	ContentImagePath string    `gorm:"column:content_image_path;type:varchar(1000)"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:current_timestamp;not null"`
	Visible          bool      `gorm:"column:visible;type:boolean;not null"`

	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Post) TableName() string {
	return "post"
}

type Comment struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	PostID    int64     `gorm:"column:post_id;type:bigint;not null"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null"`
	Content   string    `gorm:"column:content;type:text(100000);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`

	Post Post `gorm:"foreign_key:post_id;references:id"`
	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Comment) TableName() string {
	return "comment"
}

type Like struct {
	PostID    int64     `gorm:"column:post_id;type:bigint;not null;index:unique_post_id_user_id,unique"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null;index:unique_post_id_user_id,unique"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`

	Post Post `gorm:"constraint:foreign_key:post_id;references:id"`
	User User `gorm:"foreign_key:user_id;references:id"`
}

func (Like) TableName() string {
	return "like"
}

type UserUser struct {
	UserID     int64 `gorm:"column:user_id;type:bigint;not null;index:unique_user_id_follower_id,unique"`
	FollowerID int64 `gorm:"column:follower_id;type:bigint;not null;index:unique_user_id_follower_id,unique"`

	User     User `gorm:"foreign_key:user_id;references:id"`
	Follower User `gorm:"foreign_key:follower_id;references:id"`
}

func (UserUser) TableName() string {
	return "user_user"
}
