package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;type:varchar(255)"`
	Author      string         `json:"author" gorm:"not null;type:varchar(100)"`
	Description string         `json:"description" gorm:"type:text"`
	CoverImage  string         `json:"cover_image" gorm:"type:varchar(255)"`
	Category    string         `json:"category" gorm:"type:varchar(50)"`
	UserID      uint           `json:"user_id" gorm:"not null"` // 上传书籍的用户ID
	User        User           `json:"user"`                     // 关联用户
	Comments    []Comment      `json:"comments" gorm:"foreignKey:BookID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}