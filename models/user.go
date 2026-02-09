package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null;type:varchar(50)"`
	Password  string         `json:"password" gorm:"not null;type:varchar(255)"`
	Email     string         `json:"email" gorm:"unique;not null;type:varchar(100)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Books     []Book         `json:"books" gorm:"foreignKey:UserID"`    // 用户上传的书籍
	Comments  []Comment      `json:"comments" gorm:"foreignKey:UserID"` // 用户的评论
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
