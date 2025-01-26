package models

import "time"

type User struct {
	UserID       int       `json:"userID"`
	UserAlias    string    `json:"userAlias" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
	FirstName    string    `json:"firstName" binding:"required,alpha"`
	Surname      string    `json:"surname" binding:"required,alpha"`
	Phone        *string   `json:"phone"`
	PasswordHash *string   `json:"passwordHash"`
	StatusId     int       `json:"statusId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	RoleID       int       `json:"roleId"`
}
type AddPostRequest struct {
	Post Post     `json:"post"`
	Tags []string `json:"tags"`
}
type UserToken struct {
	Email      string `json:"email"`
	Givenname  string `json:"givenname"`
	Familyname string `json:"familyname"`
}

type PostWithTagsResponse struct {
	Post Post  `json:"post"`
	Tags []Tag `json:"tags"`
}
type Post struct {
	PostID       int       `json:"postID"`
	ParentPostID *int      `json:"parentPostID"`
	UserID       int       `json:"userID" binding:"required,numeric"`
	CategoryID   int       `json:"categoryID" binding:"required,numeric"`
	Title        string    `json:"title"`
	Content      string    `json:"content" binding:"required"`
	StatusId     int       `json:"statusId"`
	IsTopic      bool      `json:"isTopic"`
	Owner        bool      `json:"owner"`
	CreatedAt    time.Time `json:"createdAt" binding:"-"`
	UpdatedAt    time.Time `json:"updatedAt" binding:"-"`
}

type Tag struct {
	TagID     int       `json:"tagID" binding:"required,numeric"`
	TagName   string    `json:"tagName"`
	CreatedAt time.Time `json:"createdAt" binding:"-"`
}

type Category struct {
	CategoryID  int       `json:"categoryID" binding:"required,numeric"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" binding:"-"`
}
