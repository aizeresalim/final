package structures

import "time"

// Comment represents a comment made by a user on a blog post.
type Comment struct {
	ID       uint      `json:"id"`
	UserID   string    `json:"user_id"`
	PostID   uint      `json:"post_id"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"datetime"`
	User     User      `json:"user";gorm:"foreignkey:UserID"`
}
